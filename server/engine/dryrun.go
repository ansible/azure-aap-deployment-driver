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
	dryRunInstance     *dryRunController
	dryRunInstanceOnce sync.Once
	dryRunInstanceErr  error
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

// func (d *dryRunController) save(model *model.DryRun) error {
// 	tx := d.db.Begin()
// 	tx.Save(&model)

// 	if tx.Error != nil {
// 		tx.Rollback()
// 		return tx.Error
// 	}
// 	tx.Commit()

// 	return nil
// }

func (d *dryRunController) getStep() (model.Step, error) {
	var step model.Step
	// TODO: get model.Step from DB
	return step, nil
}

func (d *dryRunController) Execute(ctx context.Context) error {
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

	return d.save(res)
}

func DryRunControllerInstance() (*dryRunController, error) {
	dryRunInstanceOnce.Do(func() {
		dryRunInstance = &dryRunController{
			done: make(chan struct{}),
		}
	})
	return dryRunInstance, dryRunInstanceErr
}

func DryRunDone(eventHook *events.EventHookMessage) {
	controller, _ := DryRunControllerInstance()
	controller.dryRunDone(eventHook)
}

func (c *dryRunController) dryRunDone(eventHook *events.EventHookMessage) {
	// TODO update execution
	c.done <- struct{}{}
}

func (c *dryRunController) save(response *sdk.DryRunResponse) error {
	// dryRun, err := c.getStep()

	// tx.Save(&model)

	// if tx.Error != nil {
	// 	tx.Rollback()
	// 	return tx.Error
	// }
	// tx.Commit()

	return nil
}

func (c *dryRunController) getDryRun() (*model.Step, error) {
	// db := c.db

	// // there could be more dry runs, so get the most recent one
	// dryRun := &model.Step{}
	// tx := db.Model(dryRun).Order("updated_at desc").First(dryRun)
	// if tx.Error != nil { // not found
	// 	return nil, tx.Error
	// }
	// return dryRun, nil
	return nil, nil
}

func (c *dryRunController) saveResult(message *events.EventHookMessage) error {
	// TODO: save the dry run record as an execution on the step
	// get the step by model.DryRunStepName
	//	engine.database.Instance.Model(&model.Step{}).Preload("Executions").Find(&steps)
	return nil
}
