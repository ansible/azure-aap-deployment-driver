package engine

import (
	"context"
	"server/model"
	"sync"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var (
	modmRunInstance     *modmRunController
	modmRunInstanceOnce sync.Once
)

type modmRunController struct {
	deploymentId   int
	db             *gorm.DB
	execution      *model.Execution
	done           chan struct{}
	HandleError    ErrorHandler
}

func (d *modmRunController) Execute(ctx context.Context, client sdk.Client, template datatypes.JSONMap, parameters datatypes.JSONMap) {
	_, err := client.Start(ctx, d.deploymentId, parameters, &sdk.StartOptions{})
	if err != nil {
		log.Errorf("Unable to start deployment: %v", err)
		d.HandleError(err)
	}
	d.createExecution(err)
}

func (d *modmRunController) RestartStage(ctx context.Context, client sdk.Client, stageId uuid.UUID) {
	opts := sdk.RetryOptions{}
	opts.StageId = stageId
	_, err := client.Retry(ctx, d.deploymentId, &opts)
	if err != nil {
		log.Errorf("Unable to restart stage %s: %v", stageId.String(), err)
		d.HandleError(err)
	}
}

func (d *modmRunController) CancelDeployment(ctx context.Context, client sdk.Client) {
	// TODO doesn't exist yet client.CancelDeployment(ctx, d.deploymentId)
}

func NewModmRunControllerInstance(db *gorm.DB, deploymentId int) *modmRunController {
	modmRunInstanceOnce.Do(func() {
		modmRunInstance = &modmRunController{
			db:                   db,
			deploymentId:         deploymentId,
			execution:            &model.Execution{},
			done:                 make(chan struct{}),
			HandleError: func(err error) {
				if err != nil {
					log.Error(err)
				}
			},
		}
	})
	return modmRunInstance
}

func GetModmRunControllerInstance() *modmRunController {
	return modmRunInstance
}

func (c *modmRunController) getStep() (*model.Step, error) {
	step := &model.Step{}

	join := "left join executions on executions.step_id = steps.id"
	tx := c.db.Model(step).Preload("Executions").Joins(join).Where("steps.name = ?", model.DryRunStepName).First(&step)

	if tx.Error != nil { // not found
		return nil, tx.Error
	}
	return step, nil
}

// updates the step execution (or inserts) and signals dry run is done
func (c *modmRunController) Done(message *sdk.EventHookMessage) {
	c.updateExecution(message)
	c.done <- struct{}{}
}

func(c *modmRunController) StageStarted(message *sdk.EventHookMessage) {
	// Use stage ID to create execution
	// Set to in progress
}

func(c *modmRunController) StageDone(message *sdk.EventHookMessage) {
	// Use stage ID to look up execution
	// Set result
}

// creates a new step execution to track the dry run
func (c *modmRunController) createExecution(err error) error {
	tx := c.db.Begin()

	// Set status
	status := model.Started

	if err != nil {
		status = model.Failed
		c.execution.Error = err.Error()
	}

	c.execution.Status = status

	// Set Step ID
	step, err := c.getStep()
	if err != nil {
		log.Infof("Unable to get step for dry run: %v", err)
	}
	c.execution.StepID = step.ID

	tx.Save(&c.execution)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()

	return nil
}

func (c *modmRunController) updateExecution(message *sdk.EventHookMessage) error {
	data, err := message.DeploymentEventData()
	if err != nil {
		log.Debugf("event hook message is [%s] not deployment. error: %v", message.Type, err)
		return err
	}

	// Start with success
	c.execution.Status = model.Succeeded
	if message.Status == sdk.StatusFailed.String() {
		c.execution.Status = model.Failed
		c.execution.Error = message.Error
	} 
	c.execution.CorrelationID = (*data.CorrelationId).String()

	c.db.Save(c.execution)
	return nil
}
