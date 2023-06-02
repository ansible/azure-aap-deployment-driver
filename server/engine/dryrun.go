package engine

import (
	"context"
	"fmt"
	"server/config"
	"server/model"
	"strings"
	"sync"

	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var (
	dryRunInstance     *dryRunController
	dryRunInstanceOnce sync.Once
)

type dryRunController struct {
	deploymentId   int
	db             *gorm.DB
	execution      *model.Execution
	done           chan struct{}
	apiKey         string
	hookName       string
	eventHookCallbackUrl string
	HandleError          ErrorHandler
}

func (d *dryRunController) Execute(ctx context.Context, client sdk.Client, template datatypes.JSONMap, parameters datatypes.JSONMap) {
	createEventRequest := sdk.CreateEventHookRequest{
		APIKey:   &d.apiKey,
		Callback: &d.eventHookCallbackUrl,
		Name:     &d.hookName,
	}

	_, err := client.CreateEventHook(ctx, createEventRequest)
	if err != nil {
		d.HandleError(err)
	}

	executionInfo, err := client.DryRun(ctx, d.deploymentId, parameters)
	if err != nil {
		d.HandleError(err)
	}
	d.createExecution(uint(d.deploymentId), executionInfo, err)

	<-d.done
}

func NewDryRunControllerInstance(db *gorm.DB, deploymentId int) *dryRunController {
	dryRunInstanceOnce.Do(func() {
		dryRunInstance = &dryRunController{
			db:                   db,
			deploymentId:         deploymentId,
			execution:            &model.Execution{},
			apiKey:               config.GetEnvironment().WEB_HOOK_API_KEY,
			hookName:             "deployment-driver-hook",
			eventHookCallbackUrl: config.GetEnvironment().WEB_HOOK_CALLBACK_URL,
			done:                 make(chan struct{}),
			HandleError: func(err error) {
				if err != nil {
					log.Error(err)
				}
			},
		}
	})
	return dryRunInstance
}

func GetDryRunControllerInstance() *dryRunController {
	return dryRunInstance
}

func (c *dryRunController) getStep() (*model.Step, error) {
	step := &model.Step{}

	join := "left join executions on executions.step_id = steps.id"
	tx := c.db.Model(step).Preload("Executions").Joins(join).Where("steps.name = ?", model.DryRunStepName).First(&step)

	if tx.Error != nil { // not found
		return nil, tx.Error
	}
	return step, nil
}

// updates the step execution (or inserts) and signals dry run is done
func (c *dryRunController) Done(message *sdk.EventHookMessage) {
	c.updateExecution(message)
	c.done <- struct{}{}
}

// creates a new step execution to track the dry run
func (c *dryRunController) createExecution(deploymentId uint, response *sdk.InvokeDryRunResponse, err error) error {
	tx := c.db.Begin()

	status := model.Started
	if response.Status != sdk.StatusScheduled.String() || err != nil {
		status = model.Failed
	}

	step, err := c.getStep()
	if err != nil {
		log.Infof("Unable to get step for dry run: %v", err)
	}
	c.execution.StepID = step.ID
	c.execution.Status = status
	c.execution.DryRunExecution = &model.DryRunExecution{
			Id:           response.Id.String(),
			DeploymentId: deploymentId,
			Status:       "", //status is the status result of the dry run. which isn't set yet because the result hasn't been received
	}

	if err != nil {
		c.execution.Status = model.Failed
		c.execution.DryRunExecution.Status = sdk.StatusFailed.String()
		c.execution.Error = err.Error()
	}

	tx.Save(&c.execution)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()

	return nil
}

func (c *dryRunController) updateExecution(message *sdk.EventHookMessage) error {
	data, err := message.DryRunEventData()
	if err != nil {
		log.Debugf("event hook message is [%s] not dryrun. error: %v", message.Type, err)
		return err
	}

	// Start with success
	c.execution.Status = model.Succeeded

	if message.Status == sdk.StatusFailed.String() {
		// Dry run failed to run
		c.execution.Status = model.Failed
		c.execution.Error = message.Error
	} else if data.Status == sdk.StatusFailed.String() {
		// Dry run ran, but failed maybe with multiple errors
		// TODO must be a better way to do the error concatenation
		c.execution.Status = model.Failed
		var errString, errDetails strings.Builder
		for _, error := range data.Errors {
			errString.WriteString(*error.Code + ": " + *error.Message + "\n")
			for _, detail := range error.Details {
				errDetails.WriteString(*detail.Message + "\n")
			}
		}
		c.execution.Error = errString.String()
		c.execution.ErrorDetails = errDetails.String()
	}
	duration := data.CompletedAt.Sub(data.StartedAt)
	c.execution.Timestamp = data.StartedAt
	c.execution.Duration = fmt.Sprintf("%.2f seconds", duration.Seconds())
	//c.execution.DryRunExecution.Status = data.Status
	//c.execution.DryRunExecution.Errors = data.Errors
	c.execution.CorrelationID = "N/A"
	log.Infof("Dry run completed with status: %s", c.execution.Status)
	c.db.Save(c.execution)
	return nil
}
