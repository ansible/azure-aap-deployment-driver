package engine

import (
	"fmt"
	"server/azure"
	"server/config"
	"server/model"
	"server/modm"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
)

var eventReceived chan bool

func (engine *Engine) Fatalf(format string, args ...interface{}) {
	log.Errorf(format, args...)

	if !engine.status.IsFatalState {
		engine.status.IsFatalState = true
		engine.database.Instance.Save(engine.status)
	}
}

func (engine *Engine) IsFatalState() bool {
	return engine.status.IsFatalState
}

func (engine *Engine) Run() {
	if !engine.IsFatalState() {
		eventReceived = make(chan bool)
		engine.startDeploymentExecutions()
	} else {
		log.Errorln("Engine failed to start. In fatal state. Check logs.")
	}
	log.Info("Main engine loop ended.")
	engine.waitBeforeEnding()
}

func (engine *Engine) startDeploymentExecutions() {
	log.Info("Starting engine execution...")

	log.Info("Ensuring that MODM is ready...")
	modm.EnsureModmReady(engine.context, engine.modmClient)

	log.Info("Registering MODM webhook...")
	modm.CreateEventHook(engine.context, engine.modmClient)

	for engine.context.Err() == nil {
		stepsToRun := []model.Step{}
		res := engine.database.Instance.Order("id").Find(&stepsToRun)
		if res.RowsAffected == 0 {
			log.Error("No deployment steps found.")
			break
		}
		for _, step := range stepsToRun {
			latestExecution := engine.GetLatestExecution(step)
			latestExecution.StepID = step.ID // TODO not sure if needed
			log.Debugf("Step %s execution in %s state.", step.Name, latestExecution.Status)
			switch latestExecution.Status {
			case model.Started, model.Failed, model.Restarted, model.Canceled, model.PermanentlyFailed, model.RestartTimedOut:
				// TODO Handle container restart case
			case "":
				if step.Name == model.DryRunStepName {
					engine.executeModmDryRun(step, &latestExecution)
				}
			case model.Restart:
				if step.Name == model.DryRunStepName {
					log.Info("Restarting dry run...")
					// Step to restart, mark as seen and start
					latestExecution.Status = model.Restarted
					engine.database.Instance.Save(&latestExecution)
					latestExecution = model.Execution{}
					engine.executeModmDryRun(step, &latestExecution)
				} else {
					// Other step to restart
					engine.restartModmStage(stringToUuid(step.StageId))
				}
			case model.Succeeded:
				if step.Name == model.DryRunStepName && !engine.deploymentStarted {
					// Dry run is succeded, start deployment (first time)
					if engine.executeModmDeployment() {
						engine.deploymentStarted = true
					}
				}
			}
			time.Sleep(10 * time.Second)
		}

		restartRequired := false
		var executionWaitGroup sync.WaitGroup
		// if the context is not yet cancelled, check for failed executions
		if engine.context.Err() == nil {
			log.Info("Checking execution status of completed steps...")

			stepsToCheck := []model.Step{}
			engine.database.Instance.Order("id").Find(&stepsToCheck)

			currentExecutions := make([]*model.Execution, len(stepsToCheck))
			for index, step := range stepsToCheck {
				latestExecution := engine.GetLatestExecution(step)
				currentExecutions[index] = &latestExecution
			}

			// first check all executions for those that can't be restarted anymore
			terminateMainLoop := false
			for _, execution := range currentExecutions {
				if execution == nil { // skip over null elements
					continue
				}
				if execution.Status == model.Canceled {
					// cancelled execution means the whole process will be cancelled
					log.Warn("Found cancelled execution.")
					terminateMainLoop = true
					break
				}
				executionsCount := engine.countStepExecutions(execution.StepID)
				if execution.Status != model.Succeeded && execution.Status != model.Canceled &&
					(executionsCount >= int64(engine.maxExecutionRestarts)) {
					log.Errorf("Found failed deployment step that can not be restarted because it had %d executions. Maximum is %d.", executionsCount, engine.maxExecutionRestarts)
					terminateMainLoop = true
					execution.Status = model.PermanentlyFailed
					engine.database.Instance.Save(execution)
					break
				}
			}
			if terminateMainLoop {
				log.Info("Will terminate main loop because steps can't be restarted or deployment is being cancelled.")
				engine.CancelDeployment()
				break
			}

			// check all executions for those can be restarted
			for _, execution := range currentExecutions {
				log.Debugf("Execution for %d is in %s state.", execution.StepID, execution.Status)
				// check if step can be restarted
				if execution != nil && execution.Status == model.Failed {
					log.Debugf("Setting restart required due to execution for step ID %d", execution.StepID)
					restartRequired = true
					engine.startWaitingForRestart(execution, &executionWaitGroup)
				}
			}
			// wait until executions are restarted or timed out
			if restartRequired {
				log.Info("Found failed deployment step(s), waiting for those failed deployment step(s) to be restarted...")

				// check if wait for restart timed out
				restartTimedOut := false
				for _, execution := range currentExecutions {
					if execution != nil && execution.Status == model.RestartTimedOut {
						log.Error("Found failed deployment step that was not restarted.")
						restartTimedOut = true
						break
					}
				}
				if restartTimedOut {
					log.Info("Will terminate main loop because at least one deployment step was not restarted.")
					engine.CancelDeployment()
					break
				}
			}
		}
	}
}

func (engine *Engine) ReportFinalDeploymentStatusToTelemetry() {
	steps := []model.Step{}
	engine.database.Instance.Model(&model.Step{}).Preload("Executions").Find(&steps)
	status := model.DeploymentSucceeded
	for _, step := range steps {
		latestExecution := engine.GetLatestExecution(step)
		if latestExecution.Status == model.PermanentlyFailed {
			status = model.DeploymentFailed
			break
		} else if latestExecution.Status == model.Canceled {
			status = model.DeploymentCanceled
			break
		}
	}
	model.SetMetric(engine.database.Instance, model.DeployStatus, string(status), "")
}

func (engine *Engine) waitBeforeEnding() {
	// Add DeploymentMetric to Database
	engine.ReportFinalDeploymentStatusToTelemetry()
	// Publish telemetry for this deployment to Segment before starting wait time
	log.Info("Setting final metrics before sending telemetry to Segment")
	SetFinalMetrics(engine.database.Instance)
	log.Info("Sending telemetry for this deployment to Segment")
	PublishToSegment(engine.database.Instance)
	// if the context is not yet cancelled, check for failed executions
	if engine.context.Err() == nil {
		waitTime := time.Duration(config.GetEnvironment().ENGINE_END_WAIT) * time.Second
		log.Infof("Engine will wait %s before terminating...", waitTime)
		// wait for either either the timer to end or context being cancelled
		select {
		case <-time.After(waitTime): // time.After() is ok to use here because its one-time use
		case <-engine.context.Done():
		}
	}
	// Start the process to delete ourself
	if !config.GetEnvironment().SAVE_CONTAINER {
		log.Info("Engine starting storage account and container deletion and terminating...")
		// TODO delete service bus
		azure.DeleteStorageAccount(config.GetEnvironment().RESOURCE_GROUP_NAME, config.GetEnvironment().STORAGE_ACCOUNT_NAME)
		azure.DeleteContainer(config.GetEnvironment().RESOURCE_GROUP_NAME, config.GetEnvironment().CONTAINER_GROUP_NAME)
	} else {
		log.Info("Engine terminating...")
	}
	// at this point its safe to close the "done" channel
	close(engine.done)
}

func (engine *Engine) Done() <-chan struct{} {
	return engine.done
}

func (engine *Engine) countStepExecutions(stepId uint) int64 {
	var count int64
	engine.database.Instance.Model(&model.Execution{}).Where("step_id = ?", stepId).Count(&count)
	return count
}

func (engine *Engine) GetLatestExecution(step model.Step) model.Execution {
	latestExecution := model.Execution{}
	// Avoid GORM error from Last() if no executions yet
	count := engine.countStepExecutions(step.ID)
	if count > 0 {
		engine.database.Instance.Last(&latestExecution, "step_id = ?", step.ID)
	}
	return latestExecution
}

func (engine *Engine) executeModmDeployment() (started bool) {
	started = true
	_, err := engine.modmClient.Start(engine.context, engine.modmDeploymentId, engine.getParamsMap(), nil)
	if err != nil {
		log.Errorf("Engine could not start modm deployment: %v", err)
		started = false
		return
	}
	return
}

func (engine *Engine) executeModmDryRun(step model.Step, execution *model.Execution) {
	// TODO if modm sends a "start" event for dry run, then no need to create an execution here
	execution.Status = model.Started
	execution.StepID = step.ID
	engine.database.Instance.Save(&execution)
	opts := sdk.DryRunOptions{}
	opts.Retries = 0
	resp, err := engine.modmClient.DryRun(engine.context, engine.modmDeploymentId, engine.getParamsMap(), &opts)
	if err != nil {
		log.Errorf("Engine could not start dry run.")
	}
	if resp.Status != sdk.StatusScheduled.String() {
		log.Errorf("Dry run did not start.  Status: %s", resp.Status)
	}
}

func (engine *Engine) restartModmStage(stageId uuid.UUID) {
	opts := sdk.RetryOptions{}
	opts.StageId = stageId
	_, err := engine.modmClient.Retry(engine.context, engine.modmDeploymentId, &opts)
	if err != nil {
		log.Errorf("Unable to restart stage %s: %v", stageId.String(), err)
	}
}

func (engine *Engine) getParamsMap() map[string]interface{} {
	// MODM wants just a map of key/value pairs {"access": "public"} for instance
	outputs := engine.mainOutputs.Values
	outMap := make(map[string]interface{})
	for k, v := range engine.parameters {
		val, ok := outputs[k]
		if ok {
			outMap[k] = val.(map[string]interface{})["value"]
		} else {
			// Take default (empty) value since it will have correct type
			outMap[k] = v.(map[string]interface{})["value"]
		}
	}
	return outMap
}

func (engine *Engine) startWaitingForRestart(execution *model.Execution, waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)
	go engine.waitForStepRestart(execution, waitGroup)
}

func (engine *Engine) restartStepAfterDelay(delay time.Duration, execution *model.Execution) *time.Timer {
	if config.GetEnvironment().AUTO_RETRY {
		log.Tracef("Starting a timer to automatically restart step after: %s", delay)
		return time.AfterFunc(delay, func() {
			storedExecution := model.Execution{}
			engine.database.Instance.Last(&storedExecution, model.Execution{StepID: execution.StepID})
			storedExecution.Status = model.Restart
			engine.database.Instance.Save(&storedExecution)
			log.Trace("Automatically marked execution for restart.")
		})
	}
	return nil
}

func (engine *Engine) waitForStepRestart(execution *model.Execution, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	// create a timer and a ticker and release them when leaving this function
	waitTime := time.Duration(config.GetEnvironment().ENGINE_RETRY_WAIT) * time.Second
	waitTimer := time.NewTimer(waitTime)
	defer waitTimer.Stop()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	autoRestartTimer := engine.restartStepAfterDelay(time.Duration(config.GetEnvironment().AUTO_RETRY_DELAY)*time.Second, execution)
	defer func() {
		if autoRestartTimer != nil {
			autoRestartTimer.Stop()
		}
	}()
	log.Tracef("Engine will wait %s for deployment step(s) to be restarted...", waitTime)
	for keepChecking := true; keepChecking; {
		select {
		case <-engine.context.Done():
			log.Trace("Ending wait because context was cancelled.")
			keepChecking = false
		case <-waitTimer.C:
			log.Trace("Ending wait because wait time is up.")
			execution.Status = model.RestartTimedOut
			engine.database.Instance.Save(execution)
			keepChecking = false
		case <-ticker.C:
			storedExecution := model.Execution{}
			engine.database.Instance.Last(&storedExecution, model.Execution{StepID: execution.StepID})
			if storedExecution.Status == model.Restart {
				log.Trace("Ending wait because execution has been marked for restart.")
				keepChecking = false
			}
		}
	}
}

func stringToUuid(uuidAsString string) uuid.UUID {
	uid, err := uuid.Parse(uuidAsString)
	if err != nil {
		log.Errorf("Error while converting %s to UUID: %v", uuidAsString, err)
	}
	return uid
}

func (engine *Engine) CancelDeployment() {
	resp, err := engine.modmClient.Cancel(engine.context, engine.modmDeploymentId)
	if err != nil || !resp.IsCancelled {
		log.Errorf("Unable to cancel MODM deployment: %v", err)
	}
}

func (engine *Engine) CreateExecution(message *sdk.EventHookMessage) {
	execution := model.Execution{}
	switch message.Type {
	case "dryRunStarted":
		log.Debug("Creating Started execution for Dry Run.")
		step := model.Step{}
		engine.database.Instance.Where("name = ?", model.DryRunStepName).Find(&step)
		execution.Status = model.Started
		execution.StepID = step.ID
		engine.database.Instance.Save(&step)

		//TODO case sdk.EventTypeStageStarted.String():
		// Create in progress execution for stage
	}
}

func (engine *Engine) UpdateExecution(message *sdk.EventHookMessage) {
	log.Debugf("Data for message %s: %v", message.Type, message.Data)
	switch message.Type {
	case sdk.EventTypeDryRunCompleted.String():
		// Find in progress execution, set result
		step := model.Step{}
		engine.database.Instance.Where("name = ?", model.DryRunStepName).Find(&step)
		execution := engine.GetLatestExecution(step)
		data, err := message.DryRunEventData()
		// Handle duplicate messages
		if execution.Status != model.Started {
			log.Errorf("Execution for dry run already updated. Ignoring event. Status: %s", execution.Status)
		}
		// Check result
		if err != nil || message.Status != string(sdk.StatusSuccess) {
			// Failed
			execution.Status = model.Failed
			var errString, errDetails strings.Builder
			for _, error := range data.Errors {
				errString.WriteString(*error.Code + ": " + *error.Message + "\n")
				for _, detail := range error.Details {
					errDetails.WriteString(*detail.Message + "\n")
				}
			}
			execution.Error = errString.String()
			execution.ErrorDetails = errDetails.String()
		} else {
			// TODO Abstract this out if it also applies to stages.
			execution.Status = model.Succeeded
			execution.Timestamp = data.CompletedAt
			duration := data.CompletedAt.Sub(data.StartedAt)
			execution.Duration = fmt.Sprintf("%.2f seconds", duration.Seconds())
			execution.CorrelationID = "N/A"
			engine.database.Instance.Save(&execution)
		}
	case sdk.EventTypeStageCompleted.String():
		// Find in progress execution, set result
		log.Debugf("Would update stage completion for %v", message.Data)
	}
}
