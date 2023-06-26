package test

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"server/config"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type ResourceType string

const (
	RESOURCE_GROUPS ResourceType = "Microsoft.Resources/resourceGroups"
	DEPLOYMENTS     ResourceType = "Microsoft.Resources/deployments"
)

func GetTimestampNow() string {
	return time.Now().Format("2006-01-02T15:04:05.000Z")
}

type DoFunc func(req *http.Request) *http.Response

func (f DoFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func SetEnvironment() {
	os.Setenv("AZURE_SUBSCRIPTION_ID", "3f7e29ba-24e0-42f6-8d9c-5149a14bda37")
	os.Setenv("RESOURCE_GROUP_NAME", "dummy")
	os.Setenv("CONTAINER_GROUP_NAME", "dummy")
	os.Setenv("STORAGE_ACCOUNT_NAME", "dummy")
	os.Setenv("ADMIN_PASS", "password1234")
	os.Setenv("MAIN_OUTPUTS", "{}")
}

func MockDeploymentResult(name string, provisioningState armresources.ProvisioningState, params map[string]SettableValue, outputs map[string]SettableValue) DoFunc {
	responseJson := MakeDeploymentResponse(name, provisioningState, time.Now(), "PT1M1S", nil, nil)
	status := map[string]string{"status": string(provisioningState)}
	statusJson, _ := json.Marshal(status)
	return DoFunc(func(req *http.Request) *http.Response {
		headers := http.Header{}
		headers.Add("Content-Type", "application/json")
		if strings.Contains(req.URL.String(), "operationStatuses") {
			// Sometimes the client wants to poll the status, so this branch responds to those requests
			return &http.Response{
				StatusCode: 200,
				Request:    req,
				Body:       io.NopCloser(strings.NewReader(string(statusJson))),
				Header:     headers,
			}
		} else {
			headers.Add("Azure-Asyncoperation", "https://management.azure.com/subscriptions/3f7e29ba-24e0-42f6-8d9c-5149a14bda37/resourcegroups/bhavenstGroup/providers/Microsoft.Resources/deployments/networkingAAPDeploy/operationStatuses/08585341014485256763?api-version=2021-04-01")
			return &http.Response{
				StatusCode: 200,
				Request:    req,
				Body:       io.NopCloser(strings.NewReader(responseJson)),
				Header:     headers,
			}
		}
	})
}

func MockGetDeployment() DoFunc {
	finalJson := MakeDeploymentResponse("dummy", armresources.ProvisioningStateSucceeded, time.Now(), "PT1M1S", nil, nil)
	return DoFunc(func(req *http.Request) *http.Response {
		headers := http.Header{}
		headers.Add("Content-Type", "application/json")
		return &http.Response{
			StatusCode: 200,
			Request:    req,
			Body:       io.NopCloser(strings.NewReader(finalJson)),
			Header:     headers,
		}
	})
}

type counter struct {
	val int
}

var count counter

func (c *counter) incrCounterVal() {
	c.val++
}

func MockGetResourceGroupFailThenPass() DoFunc {
	passed := MakeGetResourceGroupResponse(config.GetEnvironment().RESOURCE_GROUP_NAME, true)
	failed := MakeGetResourceGroupResponse(config.GetEnvironment().RESOURCE_GROUP_NAME, false)
	return DoFunc(func(req *http.Request) *http.Response {
		if count.val == 1 {
			return &http.Response{
				StatusCode: 200,
				Request:    req,
				Body:       io.NopCloser(strings.NewReader(passed)),
			}
		} else {
			count.incrCounterVal()
			return &http.Response{
				StatusCode: 404,
				Request:    req,
				Body:       io.NopCloser(strings.NewReader(failed)),
			}
		}
	})
}

func MockTemplateFailed() DoFunc {
	templateFail := MakeTemplateFailure()

	return DoFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 400,
			Request:    req,
			Body:       io.NopCloser(strings.NewReader(templateFail)),
		}
	})
}
