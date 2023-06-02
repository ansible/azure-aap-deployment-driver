package azure

import (
	"context"
	"server/config"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerinstance/armcontainerinstance"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

type azureDetails struct {
	Credentials  *azidentity.ManagedIdentityCredential
	Subscription string
}

var azureInfo azureDetails

func NewResourceGroupsClient(opts *arm.ClientOptions) *armresources.ResourceGroupsClient {
	client, err := armresources.NewResourceGroupsClient(GetAzureInfo().Subscription, GetAzureInfo().Credentials, opts)
	if err != nil {
		log.Fatalf("Failed to get resource groups client: %v", err)
	}
	log.Trace("Got resource groups client.")
	return client
}

// Ensure Azure login/token with retry and exponential backoff
func EnsureAzureLogin(client *armresources.ResourceGroupsClient) {
	// Avoid needing to instantiate client in the main
	if client == nil {
		client = NewResourceGroupsClient(nil)
	}
	const MAX_ATTEMPTS = 5
	for attempt := 1; true; attempt++ {
		_, err := client.Get(context.Background(), config.GetEnvironment().RESOURCE_GROUP_NAME, nil)
		if err == nil {
			// Successfully logged in and retrieved data
			log.Debug("Initialized Azure connection.")
			return
		} else if attempt > MAX_ATTEMPTS {
			// Give up!
			log.Fatalf("Unable to connect to Azure after %d tries: %v", MAX_ATTEMPTS, err)
		} else {
			timeToSleep := time.Duration(1 << attempt * time.Second)
			log.Errorf("Failed to log in/connect to Azure, retry in %.0f seconds.: %v", timeToSleep.Seconds(), err)
			time.Sleep(timeToSleep)
		}
	}
}

func GetAzureInfo() azureDetails {
	if azureInfo.Credentials == nil {
		opts := azidentity.ManagedIdentityCredentialOptions{}
		opts.Retry.MaxRetries = 10
		cred, err := azidentity.NewManagedIdentityCredential(&opts)
		if err != nil {
			log.Fatalf("Error: Unable to create Azure credential: %v", err)
		}
		azureInfo.Credentials = cred
		azureInfo.Subscription = config.GetEnvironment().SUBSCRIPTION
	}
	return azureInfo
}

// Delete Azure storage account
func DeleteStorageAccount(resourceGroupName string, storageAccountName string) error {
	storageClient, err := armstorage.NewAccountsClient(GetAzureInfo().Subscription, GetAzureInfo().Credentials, nil)
	if err != nil {
		return err
	}

	// This particular API call only returns a blank/empty response, so no need to check it.
	_, err = storageClient.Delete(context.Background(), resourceGroupName, storageAccountName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Delete Azure container group and its containers
func DeleteContainer(resourceGroupName string, containerGroupName string) error {
	containerClient, err := armcontainerinstance.NewContainerGroupsClient(GetAzureInfo().Subscription, GetAzureInfo().Credentials, nil)
	if err != nil {
		return err
	}

	// Don't really care about polling for completion, since we will go away!
	_, err = containerClient.BeginDelete(context.Background(), resourceGroupName, containerGroupName, nil)
	if err != nil {
		return err
	}
	return err
}
