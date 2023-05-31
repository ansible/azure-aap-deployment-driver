package engine

import (
	"context"
	"server/config"
	"server/model"
	"sync"
	"time"

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
	// the MODM deployment id
	deploymentId   int
	db             *gorm.DB
	execution      *model.Execution
	done           chan struct{}
	apiKey         string
	hookName       string

	// this is the url that will be called by MODM. It maps to /eventhook route for handler/eventhook
	eventHookCallbackUrl string
	HandleError          ErrorHandler
}

func (d *modmRunController) Execute(ctx context.Context, client sdk.Client, template datatypes.JSONMap, parameters datatypes.JSONMap) {
	time.Sleep(10 * time.Second)

	go func() {
		createEventRequest := sdk.CreateEventHookRequest{
			APIKey:   &d.apiKey,
			Callback: &d.eventHookCallbackUrl,
			Name:     &d.hookName,
		}

		_, err := client.CreateEventHook(ctx, createEventRequest)
		if err != nil {
			d.HandleError(err)
		}

		_, err = client.Start(ctx, d.deploymentId, parameters, &sdk.StartOptions{})
		if err != nil {
			d.HandleError(err)
		}

		d.createExecution(err)
	}()

	<-d.done
}

func NewModmRunControllerInstance(db *gorm.DB, deploymentId int) *modmRunController {
	modmRunInstanceOnce.Do(func() {
		modmRunInstance = &modmRunController{
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

// creates a new step execution to track the dry run
func (c *modmRunController) createExecution(err error) error {
	tx := c.db.Begin()

	status := model.Started

	step, err := c.getStep()
	if err != nil {
		log.Infof("Unable to get step for dry run: %v", err)
	}
	c.execution.StepID = step.ID

	c.execution.Status = status

	if err != nil {
		c.execution.Status = model.Failed
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

func (c *modmRunController) updateExecution(message *sdk.EventHookMessage) error {
	data, err := message.DeploymentEventData()
	if err != nil {
		log.Debugf("event hook message is [%s] not deployment. error: %v", message.Type, err)
		return err
	}

	// Start with success
	c.execution.Status = model.Succeeded
	if message.Status == sdk.StatusFailed.String() {
		// Dry run failed to run
		c.execution.Status = model.Failed
		c.execution.Error = message.Error
	} 
	c.execution.CorrelationID = (*data.CorrelationId).String()

	c.db.Save(c.execution)
	return nil
}
