package azure_test

import (
	"bytes"
	"context"
	"os"
	"server/azure"
	"server/test"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/assert"
)

func TestEnsureAzureLogin(t *testing.T) {
	opts := azure.GetClientOptionsWithLogging()
	// Capture logging
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stdout)
	}()
	opts.Transport = test.MockGetResourceGroupFailThenPass()
	client := azure.NewResourceGroupsClient(opts)
	azure.EnsureAzureLogin(client)
	assert.Contains(t, buf.String(), "ResourceGroupNotFound")
}

func TestClientSuccess(t *testing.T) {
	opts := arm.ClientOptions{}
	opts.Transport = test.MockDeploymentResult("dummy", armresources.ProvisioningStateSucceeded, nil, nil)

	client := azure.NewDeploymentsClient(&opts)
	poller, _ := azure.StartDeployARMTemplate(context.Background(), client, "dummyDeploy", make(map[string]interface{}), make(map[string]interface{}), "")
	result, err := azure.WaitForDeployARMTemplate(context.Background(), "dummyDeploy", poller)
	assert.Equal(t, string(armresources.ProvisioningStateSucceeded), result.ProvisioningState, "Expected Succeeded provisioning state")
	assert.Nil(t, err, "No error expected, but was %v", err)
}

func TestClientFail(t *testing.T) {
	opts := arm.ClientOptions{}
	opts.Transport = test.MockDeploymentResult("dummy", armresources.ProvisioningStateFailed, nil, nil)

	client := azure.NewDeploymentsClient(&opts)
	poller, _ := azure.StartDeployARMTemplate(context.Background(), client, "dummyDeploy", make(map[string]interface{}), make(map[string]interface{}), "")
	result, err := azure.WaitForDeployARMTemplate(context.Background(), "dummyDeploy", poller)
	assert.Nil(t, result, "Deployment result expected to be nil, but was: %v", poller)
	assert.NotNil(t, err, "Error not expected to be nil but was")
}

func TestGetDeployment(t *testing.T) {
	opts := arm.ClientOptions{}
	opts.Transport = test.MockGetDeployment()

	client := azure.NewDeploymentsClient(&opts)
	deploy, err := azure.GetDeployment(context.Background(), client, "dummy")
	assert.Nil(t, err, "Error should be nil, but is: %v", err)
	assert.Equal(t, string(armresources.ProvisioningStateSucceeded), deploy.ProvisioningState)
}

func TestFailedTemplate(t *testing.T) {
	opts := arm.ClientOptions{}
	opts.Transport = test.MockTemplateFailed()
	client := azure.NewDeploymentsClient(&opts)
	deploy, err := azure.StartDeployARMTemplate(context.Background(), client, "dummy", nil, nil, "")
	assert.Nil(t, deploy, "Deploy should be nil, but is: %v", deploy)
	assert.Contains(t, err.Error(), "InvalidTemplate")
}

// TestMain wraps the tests.  Setup is done before the call to m.Run() and any
// needed teardown after that.
func TestMain(m *testing.M) {
	test.SetEnvironment()
	m.Run()
}
