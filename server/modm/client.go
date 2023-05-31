package modm

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
)

func NewModmClient(endpoint string, credential *azidentity.DefaultAzureCredential, opts *sdk.ClientOptions) *sdk.Client {
	client, err := sdk.NewClient(endpoint, credential, opts)

	if err != nil {
		log.Fatalf("Failed to get deployments client: %v", err)
	}
	log.Trace("Got deployment client.")
	return client
}
