package engine

import (
	"context"
	"server/azure"
	"server/config"
	"server/model"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var (
	dryRunInstance     *dryRunController
	dryRunInstanceOnce sync.Once
)

type ErrorHandler func(err error)
type NameOrKeyCreate func() string

type dryRunController struct {
	// the MODM deployment id
	deploymentId   int
	db             *gorm.DB
	execution      *model.Execution
	done           chan struct{}
	clientEndpoint string
	location       string
	resourceGroup  string
	subscription   string
	apiKey         string
	hookName       string
	deploymentName string

	// this is the url that will be called by MODM. It maps to /eventhook route for handler/eventhook
	eventHookCallbackUrl string
	HandleError          ErrorHandler
}

func (d *dryRunController) Execute(ctx context.Context, parameters datatypes.JSONMap) {
	time.Sleep(10 * time.Second)

	go func() {
		step, err := d.getStep()
		if err != nil {
			d.HandleError(err)
		}

		azureInfo := azure.GetAzureInfo()
		client, err := sdk.NewClient(d.clientEndpoint, azureInfo.Credentials, nil)
		if err != nil {
			d.HandleError(err)
		}

		deploymentName := d.deploymentName
		createDeployment := sdk.CreateDeployment{
			Name:           &deploymentName,
			Template:       step.Template,
			Location:       &d.location,
			ResourceGroup:  &d.resourceGroup,
			SubscriptionID: &d.subscription,
		}

		dep, err := client.Create(ctx, createDeployment)
		if err != nil {
			d.HandleError(err)
		}
		d.deploymentId = int(*dep.ID)

		createEventRequest := sdk.CreateEventHookRequest{
			APIKey:   &d.apiKey,
			Callback: &d.eventHookCallbackUrl,
			Name:     &d.hookName,
		}

		_, err = client.CreateEventHook(ctx, createEventRequest)
		if err != nil {
			d.HandleError(err)
		}

		executionInfo, err := client.DryRun(ctx, d.deploymentId, parameters)
		if err != nil {
			d.HandleError(err)
		}

		d.create(uint(d.deploymentId), executionInfo, err)
	}()

	<-d.done
}

func NewDryRunControllerInstance(db *gorm.DB, execution *model.Execution) *dryRunController {
	dryRunInstanceOnce.Do(func() {
		dryRunInstance = &dryRunController{
			db:                   db,
			execution:            execution,
			resourceGroup:        config.GetEnvironment().RESOURCE_GROUP_NAME,
			subscription:         config.GetEnvironment().SUBSCRIPTION,
			location:             config.GetEnvironment().AZURE_LOCATION,
			apiKey:               config.GetEnvironment().WEB_HOOK_API_KEY,
			hookName:             "deployment-driver-hook",
			deploymentName:       "deployment-driver-" + uuid.New().String(),
			eventHookCallbackUrl: config.GetEnvironment().WEB_HOOK_CALLBACK_URL,
			clientEndpoint:       "http://localhost:8080",
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
	c.update(message)
	c.done <- struct{}{}
}

// creates a new step execution to track the dry run
func (c *dryRunController) create(deploymentId uint, response *sdk.InvokeDryRunResponse, err error) error {
	tx := c.db.Begin()

	status := model.Started
	if response.Status != sdk.StatusScheduled.String() || err != nil {
		status = model.Failed
	}

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

func (c *dryRunController) update(message *sdk.EventHookMessage) error {
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
		var errString strings.Builder
		for _, error := range data.Errors {
			errString.WriteString(*error.Code + ": " + *error.Message + "\n")
		}
		c.execution.Error = errString.String()
	}
	duration := data.CompletedAt.Sub(data.StartedAt)
	c.execution.Timestamp = data.StartedAt
	c.execution.Duration = duration.String() // TODO match formatting with other code
	c.execution.DryRunExecution.Status = data.Status
	c.execution.DryRunExecution.Errors = data.Errors

	c.db.Save(c.execution)
	return nil
}

