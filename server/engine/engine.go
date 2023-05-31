package engine

import (
	"server/azure"
	"server/config"
	"server/model"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

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
		engine.startDeploymentExecutions()
	} else {
		log.Errorln("Engine failed to start. In fatal state. Check logs.")
	}
	log.Info("Main engine loop ended.")
	engine.waitBeforeEnding()
}

func (engine *Engine) startDeploymentExecutions() {
	log.Info("Starting engine execution...")

	dryRunStep := model.Step{}
	tx := engine.database.Instance.Where("name = ?", model.DryRunStepName).Find(&dryRunStep)
	if tx.RowsAffected != 1 {
		log.Fatal("Unable to find dry run step in database.")
	}
	
	log.Info("Executing Dry Run...")
	engine.executeDryRun(dryRunStep)
	// TODO Check result and wait for retry or cancel

	log.Info("Executing deployment...")
	engine.executeModmDeployment()

	var executionWaitGroup sync.WaitGroup

	for ; engine.context.Err() == nil; {
		stepsToRun := []model.Step{}
		// TODO change this next block to check length of the array instead of looking at DB stuff
		res := engine.database.Instance.Order("id").Find(&stepsToRun)
		if res.RowsAffected == 0 {
			// No steps at this order level, get out of here
			log.Info("No more deployment steps found.")
			break
		}

		stepNames := make([]string, len(stepsToRun))
		for n, step := range stepsToRun {
			stepNames[n] = step.Name
		}
		log.Infof("Next deployment steps to execute: %v", stepNames)

		// with the slice being size of steps the elements can be null!
		currentExecutions := make([]*model.Execution, len(stepsToRun))

		for _, step := range stepsToRun {
			latestExecution := engine.GetLatestExecution(step)

			switch latestExecution.Status {
			case model.Started:
				// After container restart, we may have in-progress deployments to restart
				// TODO Handle container restart case w/ MODM
			case "":
				// Unexecuted step, just wait for MODM
			case model.Restart:
				// Step to restart, mark as seen and start
				// TODO Figure out how to restart failed steps with MODM
			case model.Succeeded, model.Canceled:
				continue
			}
		}
		// TODO Figure out how to handle loop
		time.Sleep(60 * time.Second)

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
				// check if step can be restarted
				if execution != nil && execution.Status != model.Succeeded && execution.Status != model.Canceled {
					restartRequired = true
					engine.startWaitingForRestart(execution, &executionWaitGroup)
				}
			}
			// wait until executions are restarted or timed out
			if restartRequired {
				log.Info("Found failed deployment step(s), waiting for those failed deployment step(s) to be restarted...")
				// wait for all go routines to finish again
				executionWaitGroup.Wait()
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
		// TODO continue this loop until all steps are Succeeded
		// OR better yet, wait at the end of each loop for a stage completion event
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

/* TODO Figure out how to implement cancel for MODM
func (engine *Engine) CancelStep(step model.Step) {
	err := azure.CancelDeployment(engine.context, engine.deploymentsClient, step.Name)
	if err != nil {
		log.Errorf("Couldn't cancel deployment: %v", err)
	}
}
 */
