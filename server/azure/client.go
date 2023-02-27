package azure

import (
	"context"
	"server/config"
	"server/model"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerinstance/armcontainerinstance"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

type azureDetails struct {
	Credentials  *azidentity.DefaultAzureCredential
	Subscription string
}

var azureInfo azureDetails

func NewDeploymentsClient(opts *arm.ClientOptions) *armresources.DeploymentsClient {
	client, err := armresources.NewDeploymentsClient(getAzureInfo().Subscription, getAzureInfo().Credentials, opts)
	if err != nil {
		log.Fatalf("Failed to get deployments client: %v", err)
	}
	log.Trace("Got deployment client.")
	return client
}

func NewResourceGroupsClient(opts *arm.ClientOptions) *armresources.ResourceGroupsClient {
	client, err := armresources.NewResourceGroupsClient(getAzureInfo().Subscription, getAzureInfo().Credentials, opts)
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

func getAzureInfo() azureDetails {
	if azureInfo.Credentials == nil {
		opts := azidentity.DefaultAzureCredentialOptions{}
		opts.Retry.MaxRetries = 10
		cred, err := azidentity.NewDefaultAzureCredential(&opts)
		if err != nil {
			log.Fatalf("Error: Unable to create Azure credential: %v", err)
		}
		azureInfo.Credentials = cred
		azureInfo.Subscription = config.GetEnvironment().SUBSCRIPTION
	}
	return azureInfo
}

// Returns a deployment poller, from which the caller can extract a resume token in case the deployment is interrupted
// The poller should be passed to CompleteDeployARMTemplate next to await the deployment
func StartDeployARMTemplate(ctx context.Context, client *armresources.DeploymentsClient, name string, template map[string]interface{}, parameters map[string]interface{}, resumeToken string) (*runtime.Poller[armresources.DeploymentsClientCreateOrUpdateResponse], error) {

	opts := armresources.DeploymentsClientBeginCreateOrUpdateOptions{}

	// Restart of interrupted deployment
	if resumeToken != "" {
		opts.ResumeToken = resumeToken
	}
	deploy, err := client.BeginCreateOrUpdate(
		ctx,
		config.GetEnvironment().RESOURCE_GROUP_NAME,
		name,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Template:   template,
				Mode:       to.Ptr(armresources.DeploymentModeIncremental),
				Parameters: parameters,
			},
		},
		&opts,
	)
	if err != nil {
		return nil, err
	}
	log.Tracef("Triggered deployment for [%s] (resume token present: %v)", name, resumeToken != "")
	return deploy, err
}

// Pass the deployment poller from StartDeployARMTemplate to await its completion and result
// Returns model.DeploymentResult and error if any
func WaitForDeployARMTemplate(ctx context.Context, name string, deployment *runtime.Poller[armresources.DeploymentsClientCreateOrUpdateResponse]) (*model.DeploymentResult, error) {
	log.Tracef("Starting polling until deployment of [%s] is done...", name)
	resp, err := deployment.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Duration(config.GetEnvironment().AZURE_POLLING_FREQ_SECONDS) * time.Second})
	if err != nil {
		return nil, err
	}
	log.Tracef("Finished polling, deployment of [%s] is done.", name)
	return model.NewDeploymentResult(resp.DeploymentExtended), nil
}

func GetDeployment(ctx context.Context, client *armresources.DeploymentsClient, name string) (*model.DeploymentResult, error) {
	opts := armresources.DeploymentsClientGetOptions{}

	details, err := client.Get(
		ctx,
		config.GetEnvironment().RESOURCE_GROUP_NAME,
		name,
		&opts,
	)
	if err != nil {
		return nil, err
	}
	return model.NewDeploymentResult(details.DeploymentExtended), err
}

func CancelDeployment(ctx context.Context, client *armresources.DeploymentsClient, name string) error {
	_, err := client.Cancel(ctx, config.GetEnvironment().RESOURCE_GROUP_NAME, name, &armresources.DeploymentsClientCancelOptions{})
	if err != nil {
		return err
	}
	return nil
}

// Delete Azure storage account
func DeleteStorageAccount(resourceGroupName string, storageAccountName string) error {
	storageClient, err := armstorage.NewAccountsClient(getAzureInfo().Subscription, getAzureInfo().Credentials, nil)
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
	containerClient, err := armcontainerinstance.NewContainerGroupsClient(getAzureInfo().Subscription, getAzureInfo().Credentials, nil)
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
