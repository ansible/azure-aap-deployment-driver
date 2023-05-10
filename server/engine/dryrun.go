package engine

import (
	"context"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	log "github.com/sirupsen/logrus"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

var (
	dryRunInstance *dryRunController
	dryRunInstanceOnce sync.Once
	dryRunInstanceErr error
)

type dryRunController struct {
	done chan struct{}
	clientEndpoint string
	dryRunCancelFunc context.CancelFunc
}

func (d *dryRunController) Execute(deploymentId int, paramsMap map[string]interface{}) (*sdk.DryRunResponse, error) {
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
		return nil, err
	}
	return res, nil
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
