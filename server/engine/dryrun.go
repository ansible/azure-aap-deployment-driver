package engine

import (
	"context"
	"fmt"
	"server/azure"
	"server/config"
	"server/model"
	"server/persistence"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
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

func (d *dryRunController) Execute(ctx context.Context) {
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

		executionInfo, err := client.DryRun(ctx, d.deploymentId, step.Parameters)
		if err != nil {
			d.HandleError(err)
		}

		d.create(uint(d.deploymentId), executionInfo, err)
	}()

	<-d.done
}

func DryRunControllerInstance() *dryRunController {
	dryRunInstanceOnce.Do(func() {
		dryRunInstance = &dryRunController{
			db:                   persistence.NewPersistentDB(config.GetEnvironment().DB_PATH).Instance,
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
	step, stepErr := c.getStep()
	if stepErr != nil {
		return stepErr
	}

	tx := c.db.Begin()

	status := model.Started
	if response.Status != sdk.StatusScheduled.String() || err != nil {
		status = model.Failed
	}

	execution := &model.Execution{
		StepID: step.ID,
		Status: status,
		DryRunExecution: &model.DryRunExecution{
			Id:           response.Id.String(),
			DeploymentId: deploymentId,
			Status:       "", //status is the status result of the dry run. which isn't set yet because the result hasn't been received
		},
	}

	if err != nil {
		execution.Status = model.Failed
		execution.DryRunExecution.Status = sdk.StatusFailed.String()
		execution.Error = err.Error()
	}

	tx.Save(&execution)

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

	step, err := c.getStep()
	if err != nil {
		return err
	}

	id := data.OperationId.String()
	execution, err := c.getDryRunExecution(data.OperationId.String())
	if err != nil {
		return fmt.Errorf("failed to find a dry run execution with id [%s] for stepID [%d]: %v", id, err, step.ID)
	}

	status := model.Succeeded
	if message.Status == sdk.StatusFailed.String() {
		status = model.Failed
	}

	execution.Status = status
	execution.Error = message.Error

	dryRunStatus := data.Status
	if dryRunStatus == nil || data.Error != nil {
		dryRunStatus = to.Ptr(sdk.StatusFailed.String())
	}

	execution.DryRunExecution.Status = *dryRunStatus
	execution.DryRunExecution.Error = data.Error

	c.db.Save(&step.Executions)
	return nil
}

// method that gets the step exeuction
func (c *dryRunController) getDryRunExecution(id string) (*model.Execution, error) {
	step, err := c.getStep()
	if err != nil {
		return nil, err
	}

	var execution *model.Execution

	for i := range step.Executions {
		if step.Executions[i].DryRunExecution != nil && step.Executions[i].DryRunExecution.Id == id {
			execution = &step.Executions[i]
			break
		}
	}

	if execution == nil {
		execution = &model.Execution{
			StepID: step.ID,
			DryRunExecution: &model.DryRunExecution{
				Id: id,
			},
		}
		step.Executions = append(step.Executions, *execution)
		c.db.Save(&step.Executions)
	}
	return execution, nil
}
