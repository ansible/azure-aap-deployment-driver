package modm

import (
	"context"
	"server/config"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
)

const MODM_READY_WITHIN_MINUTES time.Duration = 5 * time.Minute
const MODM_HEALTH_CHECK_FREQ_SECONDS time.Duration = 10 * time.Second

func NewModmClient(endpoint string, credential *azidentity.ManagedIdentityCredential, opts *sdk.ClientOptions) *sdk.Client {
	client, err := sdk.NewClient(endpoint, credential, opts)

	if err != nil {
		log.Fatalf("Failed to get deployments client: %v", err)
	}
	log.Trace("Got deployment client.")
	return client
}

func EnsureModmReady(ctx context.Context, client *sdk.Client) {
	var elapsed time.Duration = 0
	for elapsed < MODM_READY_WITHIN_MINUTES {
		log.Debug("Checking modm health status...")
		status, err := client.HealthStatus(ctx)
		if err != nil {
			log.Errorf("Unable to get health status from MODM: %v", err)
		}
		if *status.IsHealthy {
			log.Info("MODM is ready to go...")
			return
		} else {
			log.Infof("MODM not ready, will check again in %.0f seconds.", MODM_HEALTH_CHECK_FREQ_SECONDS.Seconds())
		}
		sleep := 10 * time.Second
		time.Sleep(sleep)
		elapsed += sleep
	}
	log.Fatalf("MODM did not become healthy within %.0f minutes. Exiting.", MODM_READY_WITHIN_MINUTES.Minutes())
}

func CreateEventHook(ctx context.Context, client *sdk.Client) {
	createEventRequest := sdk.CreateEventHookRequest{
		APIKey:   to.Ptr(config.GetEnvironment().WEB_HOOK_API_KEY),
		Callback: to.Ptr(config.GetEnvironment().WEB_HOOK_CALLBACK_URL),
		Name:     to.Ptr("deployment-driver-hook"),
	}

	_, err := client.CreateEventHook(ctx, createEventRequest)
	if err != nil {
		log.Fatalf("Unable to create MODM event hook, exiting: %v", err)
	}
}
