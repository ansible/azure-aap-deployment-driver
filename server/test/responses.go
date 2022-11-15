package test

import (
	"encoding/json"
	"fmt"
	"server/config"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type SettableValue struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func makeStringPointer(str string) *string {
	return &str
}

// Timestamp format: "2022-11-03T19:09:17.0071829Z"
func MakeDeploymentResponse(name string, provisioningState armresources.ProvisioningState, timestamp time.Time, duration string, params map[string]SettableValue, outputs map[string]SettableValue) string {
	resp := armresources.DeploymentExtended{}
	props := armresources.DeploymentPropertiesExtended{}
	resp.Type = makeStringPointer(string(DEPLOYMENTS))
	props.CorrelationID = makeStringPointer("12345678-90ab-cdef-0123-4567890abcde")
	props.TemplateHash = makeStringPointer("1234567890123456789")
	mode := armresources.DeploymentModeIncremental
	props.Mode = &mode

	resp.ID = makeStringPointer(fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Resources/deployments/%s", config.GetEnvironment().SUBSCRIPTION, config.GetEnvironment().RESOURCE_GROUP_NAME, name))
	resp.Name = &name
	props.Parameters = &params
	props.Outputs = &outputs
	props.ProvisioningState = &provisioningState
	props.Timestamp = &timestamp
	props.Duration = &duration

	resp.Properties = &props
	jsonBytes, _ := json.Marshal(resp)
	return string(jsonBytes)
}

// If you want a "does not exist" error, set exists to false
func MakeGetResourceGroupResponse(name string, exists bool) string {
	var jsonBytes []byte
	if exists {
		resp := armresources.DeploymentExtended{}
		props := armresources.DeploymentPropertiesExtended{}
		resp.Location = makeStringPointer("eastus")
		resp.Type = makeStringPointer(string(RESOURCE_GROUPS))

		resp.ID = makeStringPointer(fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", config.GetEnvironment().SUBSCRIPTION, name))
		resp.Name = &name
		state := armresources.ProvisioningStateSucceeded
		props.ProvisioningState = &state
		resp.Properties = &props
		jsonBytes, _ = json.Marshal(resp)
	} else {
		resp := armresources.ErrorResponse{}
		resp.Code = makeStringPointer("ResourceGroupNotFound")
		resp.Message = makeStringPointer("Resource group 'dummy' could not be found.")
		jsonBytes, _ = json.Marshal(resp)
	}
	return string(jsonBytes)
}

func MakeTemplateFailure() string {
	templateFail := armresources.ErrorResponse{}
	templateFail.Code = makeStringPointer("InvalidTemplate")
	templateFail.Message = makeStringPointer("Deployment template validation failed: 'The template resource 'virtualNetworkDeploy' at line blah blah..")
	jsonBytes, _ := json.Marshal(templateFail)
	return string(jsonBytes)
}
