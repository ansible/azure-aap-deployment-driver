package engine

import (
	"context"
	"server/azure"
	"server/config"
	"server/model"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func (engine *Engine) Run() {
	log.Println("Starting main engine loop...")

	var executionWaitGroup sync.WaitGroup

	// Find lowest priority step(s) without successful execution and run
	for p := 0; engine.context.Err() == nil; {
		stepsToRun := []model.Step{}
		// TODO change this next block to check length of the array instead of looking at DB stuff
		res := engine.database.Instance.Where("priority = ?", p).Find(&stepsToRun)
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

		for stepIndex, step := range stepsToRun {
			latestExecution := engine.getLatestExecution(step)

			switch latestExecution.Status {
			case model.Started:
				// After container restart, we may have in-progress deployments to restart
				engine.startExecution(step, &latestExecution, &executionWaitGroup)
			case "":
				// Unexecuted step
				engine.startExecution(step, &latestExecution, &executionWaitGroup)
			case model.Restart:
				// Step to restart, mark as seen and start
				latestExecution.Status = model.Restarted
				engine.database.Instance.Save(&latestExecution)

				engine.startExecution(step, &model.Execution{}, &executionWaitGroup)
			case model.Succeeded:
				continue
			}
			currentExecutions[stepIndex] = &latestExecution
		}
		// wait for all go routines to finish
		log.Info("Waiting for execution of step(s) to finish...")
		executionWaitGroup.Wait()

		restartRequired := false
		// if the context is not yet cancelled, check for failed executions
		if engine.context.Err() == nil {
			log.Info("Checking execution status of completed steps...")
			// first check all executions for those that can't be restarted anymore
			foundPermanentlyFailedExecution := false
			for _, execution := range currentExecutions {
				if execution != nil && execution.Status != model.Succeeded && execution.ExecutionCount == engine.maxExecutionRestarts {
					log.Error("Found failed deployment step that can not be restarted again.")
					foundPermanentlyFailedExecution = true
					execution.Status = model.PermanentlyFailed
					engine.database.Instance.Save(execution)
					break
				}
			}
			if foundPermanentlyFailedExecution {
				log.Info("Will terminate main loop because at least one deployment step can not be restarted.")
				break
			}
			// check all executions for those can be restarted
			for _, execution := range currentExecutions {
				// check if step can be restarted
				if execution != nil && execution.Status != model.Succeeded {
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

		// if no executions need to be restarted, increment priority level to move to next level
		if !restartRequired {
			p++
		}
	}
	log.Info("Main engine loop ended.")
	engine.waitBeforeEnding()
}

func (engine *Engine) waitBeforeEnding() {
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
	// at this point its safe to close the "done" channel
	close(engine.done)
}

func (engine *Engine) Done() <-chan struct{} {
	return engine.done
}

func (engine *Engine) getLatestExecution(step model.Step) model.Execution {
	latestExecution := model.Execution{}
	// Avoid GORM error from Last() if no executions yet
	var count int64
	engine.database.Instance.Model(&model.Execution{}).Where("step_id = ?", step.ID).Count(&count)
	if count > 0 {
		engine.database.Instance.Last(&latestExecution, "step_id = ?", step.ID)
	}
	return latestExecution
}

func (engine *Engine) startExecution(step model.Step, execution *model.Execution, waitGroup *sync.WaitGroup) {
	execution.Status = model.Started
	execution.StepID = step.ID

	// Run in goroutine to allow parallel deployments
	log.Infof("Starting execution of deployment step [%s]...", step.Name)
	waitGroup.Add(1)
	go engine.runStep(step, execution, waitGroup)
}

func (engine *Engine) startWaitingForRestart(execution *model.Execution, waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)
	go engine.waitForStepRestart(execution, waitGroup)
}

func (engine *Engine) waitForStepRestart(execution *model.Execution, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	// create a timer and a ticker and release them when leaving this function
	waitTime := time.Duration(config.GetEnvironment().ENGINE_RETRY_WAIT) * time.Second
	waitTimer := time.NewTimer(waitTime)
	defer waitTimer.Stop()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
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

func (engine *Engine) runStep(step model.Step, execution *model.Execution, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	// Check if this is an interrupted/restarted deployment
	resumeToken := ""
	if execution.Status == model.Started && execution.ResumeToken != "" {
		resumeToken = execution.ResumeToken
	}

	engine.resolver.ResolveReferencesToParameters(step.Parameters, engine.mainOutputs.Values)

	// find all outputs, skip over those with no module names and build a map of them
	outputValues := make(map[string]map[string]interface{})
	var allOutputs []model.Output
	engine.database.Instance.Model(&model.Output{}).Find(&allOutputs)
	for _, v := range allOutputs {
		if v.ModuleName != "" {
			outputValues[v.ModuleName] = v.Values
		}
	}
	engine.resolver.ResolveReferencesToOutputs(step.Parameters, outputValues)

	// Create the deployment
	deployment, err := azure.StartDeployARMTemplate(engine.context, step.Name, step.Template, step.Parameters, resumeToken)
	if err != nil {
		if err == context.Canceled {
			log.Printf("Starting of step [%s] deployment interrupted by shutdown.", step.Name)
			return
		}
		log.Printf("Failed to start step [%s] deployment: %v", step.Name, err)
		model.UpdateExecution(execution, nil, model.GetAzureErrorJSONString(err))
		engine.database.Instance.Save(&execution)
		return
	}
	log.Printf("Started execution of step [%s]", step.Name)

	// Deployment started, grab resume token in case we get restarted
	token, err := deployment.ResumeToken()
	if err != nil {
		log.Printf("Failed to extract resume token from started deployment: %v", err)
	}
	execution.ResumeToken = token
	if err := engine.database.Instance.Save(&execution).Error; err != nil {
		log.Printf("Failed to update execution in DB with resume token: %v", err)
	}

	// Finish deployment and wait for result
	deployResponse, err := azure.WaitForDeployARMTemplate(engine.context, step.Name, deployment)
	if err != nil {
		if err == context.Canceled {
			log.Printf("Completion of step [%s] deployment interrupted by shutdown.", step.Name)
			return
		}
		log.Printf("Deployment of step [%s] failed: %v", step.Name, err)
		failedDeploymentResponse, getDeploymentErr := azure.GetDeployment(engine.context, step.Name)
		if getDeploymentErr != nil {
			log.Tracef("Unable to get failed deployment details: %v", getDeploymentErr)
		}
		model.UpdateExecution(execution, failedDeploymentResponse, model.GetAzureErrorJSONString(err))
		engine.database.Instance.Save(&execution)
		return
	}
	log.Printf("Deployment of step [%s] complete", step.Name)

	// store outputs
	engine.database.Instance.Create(model.CreateNewOutput(step.Name, deployResponse))
	// store execution
	model.UpdateExecution(execution, deployResponse, "")
	engine.database.Instance.Save(&execution)
}
