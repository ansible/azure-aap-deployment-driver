package engine

import (
	"context"
	"server/model"
	"sync"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	dryRunInstance *dryRunController
	dryRunInstanceOnce sync.Once
	dryRunInstanceErr error
)

type dryRunController struct {
	deploymentId 			int
	db 						*gorm.DB
	done 					chan struct{}
	clientEndpoint 			string
	location 				string
	resourceGroup 			string
	subscription			string
	apiKey 					string
	hookName 				string
	callbackClientEndpoint	string
}

func (d *dryRunController) save(model *model.DryRun) error {
	tx := d.db.Begin()
	tx.Save(&model)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()

	return nil
}

func (d *dryRunController) getStep() (model.Step, error) {
	var step model.Step
	// TODO: get model.Step from DB
	return step, nil
}

func (d *dryRunController) Execute(ctx context.Context) (error) {
	step, err := d.getStep()
	if err != nil {
		// TODO: handle error
		return err
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Error(err)
	}

	client, err := sdk.NewClient(d.clientEndpoint, cred, nil)
	if err != nil {
		log.Println(err)
	}

	deploymentName := "TaggedDeployment"
	request := api.CreateDeployment{
		Name:           &deploymentName,
		Template:       step.Template,
		Location:       &d.location,
		ResourceGroup:  &d.resourceGroup,
		SubscriptionID: &d.subscription,
	}

	dep, err := client.Create(ctx, request)
	if err != nil {
		// TODO: handle error
		return err
	}



	d.deploymentId = int(*dep.ID)

	res, err := client.DryRun(ctx, d.deploymentId, step.Parameters)
	if err != nil {
		return err
	}

	dryRun := model.DryRun{
		OperationId: res.Id,
		Status: res.Status,
		Result: "",
	}

	return d.save(&dryRun)
}

func DryRunControllerInstance() (*dryRunController, error) {
	dryRunInstanceOnce.Do(func() {
		dryRunInstance = &dryRunController{
			done: make(chan struct{}),
		}
	})
	return dryRunInstance, dryRunInstanceErr
}

func DryRunDone(eventHook *events.EventHookMessage)  {
	dryRunControllerInstance, err := DryRunControllerInstance()
	if err != nil {
		// TODO: handle error
		log.Error(err)
	}

	dryRunControllerInstance.done <- struct{}{}
}
