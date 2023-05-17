package engine

import (
	"context"
	"server/azure"
	"server/config"
	"server/model"
	"server/persistence"
	"strconv"
	"sync"
	"time"

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
		request := sdk.CreateDeployment{
			Name:           &deploymentName,
			Template:       step.Template,
			Location:       &d.location,
			ResourceGroup:  &d.resourceGroup,
			SubscriptionID: &d.subscription,
		}

		dep, err := client.Create(ctx, request)
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

		res, err := client.DryRun(ctx, d.deploymentId, step.Parameters)
		if err != nil {
			d.HandleError(err)
		}

		d.create(res)
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
			deploymentName:       "deployment-driver/" + uuid.New().String(),
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
	tx := c.db.Model(step).Preload("Executions").Joins(join).Where("steps.name = ?", model.DryRunStepName).First(step)
	if tx.Error != nil { // not found
		return nil, tx.Error
	}
	return step, nil
}

// updates the step execution (or inserts) and signals dry run is done
func (c *dryRunController) DryRunDone(message *sdk.EventHookMessage) {
	c.update(message)
	c.done <- struct{}{}
}

// creates a new step execution to track the dry run
func (c *dryRunController) create(response *sdk.InvokeDryRunResponse) error {
	tx := c.db.Begin()
	step, err := c.getStep()
	if err != nil {
		return err
	}

	status := model.Started
	if response.Status != sdk.StatusScheduled.String() {
		status = model.Failed
	}

	execution := model.Execution{
		StepID:        step.ID,
		DeploymentID:  strconv.Itoa(c.deploymentId),
		Status:        status,
		CorrelationID: response.Id.String(),
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
	step, err := c.getStep()
	if err != nil {
		return err
	}
	data := message.Data.(sdk.DeploymentEventData)
	var execution *model.Execution

	for i := range step.Executions {
		if step.Executions[i].CorrelationID == data.OperationId.String() {
			execution = &step.Executions[i]
			break
		}
	}

	if execution == nil {
		execution = &model.Execution{StepID: step.ID, CorrelationID: data.OperationId.String()}
		step.Executions = append(step.Executions, *execution)
	}

	status := model.Succeeded
	if message.Status == sdk.StatusFailed.String() {
		status = model.Failed
	}
	execution.Status = status
	execution.Details = data.Message

	c.db.Save(&step.Executions)
	return nil
}
