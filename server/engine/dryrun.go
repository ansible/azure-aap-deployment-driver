package engine

import (
	"context"
	"server/model"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
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
	db *gorm.DB
	done chan struct{}
	clientEndpoint string
	dryRunCancelFunc context.CancelFunc
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

func (d *dryRunController) Execute(deploymentId int, paramsMap map[string]interface{}) (error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Error(err)
	}
	client, err := sdk.NewClient(d.clientEndpoint, cred, nil)
	if err != nil {
		log.Println(err)
	}

	ctx := context.Background()
	ctx, dryRunCancelFunc := context.WithTimeout(ctx, 60 * time.Minute)
	d.dryRunCancelFunc = dryRunCancelFunc

	res, err := client.DryRun(ctx, deploymentId, paramsMap)
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
	// save to db
	// call channel to proceed to steps
}
