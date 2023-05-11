package engine

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

func (d *dryRunController) Execute(ctx context.Context) {
	go func() {
		step, err := d.getStep()
		if err != nil {
			d.HandleError(err)
		}

		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			d.HandleError(err)
		}

		client, err := sdk.NewClient(d.clientEndpoint, cred, nil)
		if err != nil {
			d.HandleError(err)
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
			d.HandleError(err)
		}
		d.deploymentId = int(*dep.ID)

		createEventRequest := api.CreateEventHookRequest{
			APIKey:   &d.apiKey,
			Callback: &d.callbackClientEndpoint,
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
