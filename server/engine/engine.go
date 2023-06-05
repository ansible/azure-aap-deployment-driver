package engine

import (
	"server/azure"
	"server/config"
	"server/model"
	"server/modm"
	"sync"
	"time"

	"github.com/google/uuid"
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

func ProcessEvent() {
	// Trigger another loop of engine
	log.Debug("Processing event")
	eventReceived <- true
}

func (engine *Engine) startDeploymentExecutions() {
	log.Info("Starting engine execution...")

	log.Info("Ensuring that MODM is ready...")
	modm.EnsureModmReady(engine.context, engine.modmClient)

	var executionWaitGroup sync.WaitGroup

	for ; engine.context.Err() == nil; {
		stepsToRun := []model.Step{}
		res := engine.database.Instance.Order("id").Find(&stepsToRun)
		if res.RowsAffected == 0 {
			// No steps at this order level, get out of here
			log.Info("No deployment steps found.")
			break
		}

		stepNames := make([]string, len(stepsToRun))
		for n, step := range stepsToRun {
			stepNames[n] = step.Name
		}

		// with the slice being size of steps the elements can be null!
		currentExecutions := make([]*model.Execution, len(stepsToRun))

		for index, step := range stepsToRun {
			latestExecution := engine.GetLatestExecution(step)
			latestExecution.StepID = step.ID  // TODO not sure if needed
			currentExecutions[index] = &latestExecution
			switch latestExecution.Status {
			case model.Started:
				log.Debugf("Step %s in Started state.", step.Name)
				// TODO Nothing to do except to handle container restart case
			case "":
				log.Debugf("Step %s in New state.", step.Name)
				// Start dry run, first time only
				if step.Name == model.DryRunStepName {
					log.Info("Starting dry run...")
					engine.executeDryRun(step)
					// Get result
					latestExecution = engine.GetLatestExecution(step)
					currentExecutions[index] = &latestExecution
					if latestExecution.Status == model.Succeeded {
						// Kick off deployment
						log.Info("Starting deployment...")
						engine.executeModmDeployment()
					}
				}
			case model.Restart:
				log.Debugf("Step %s in Restart state.", step.Name)
				// Step to restart, mark as seen and start
				latestExecution.Status = model.Restarted
				engine.database.Instance.Save(&latestExecution)

				if step.Name == model.DryRunStepName {
					// Restart dry run
					log.Info("Restarting dry run...")
					engine.executeDryRun(step)
					// Get result
					latestExecution = engine.GetLatestExecution(step)
					currentExecutions[index] = &latestExecution
					if latestExecution.Status == model.Succeeded {
						// Kick off deployment
						log.Info("Starting deployment...")
						engine.executeModmDeployment()
					}
				} else {
					// Other step to restart
					controller := GetModmRunControllerInstance()
					controller.RestartStage(engine.context, *engine.modmClient, stringToUuid(step.StageId))
				}
			case model.Succeeded, model.Canceled:
				log.Debugf("Step %s in Succeeded or Canceled state.", step.Name)
				// Nothing to do
			}
		}

		restartRequired := false
		// if the context is not yet cancelled, check for failed executions
		if engine.context.Err() == nil {
			log.Info("Checking execution status of completed steps...")
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
					break
				}
			}
		}

		// TODO Need to break loop when all steps have success
		log.Debug("Engine waiting for event.")
		<-eventReceived
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

func (engine *Engine) executeDryRun(step model.Step) {
	dryRunController := NewDryRunControllerInstance(engine.database.Instance, engine.modmDeploymentId)
	dryRunController.Execute(engine.context, *engine.modmClient, engine.template, engine.getParamsMap())
}

func (engine *Engine) executeModmDeployment() {
	_, err := engine.modmClient.Start(engine.context, engine.modmDeploymentId, engine.getParamsMap(), nil)
	if err != nil {
		log.Errorf("Failed to start MODM deployment: %v", err)
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

func (engine *Engine) CancelStep(step model.Step) {
	// TODO any error handling or whatever when implemented
	GetModmRunControllerInstance().CancelDeployment(engine.context, *engine.modmClient)
}
