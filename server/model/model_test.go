package model_test

import (
	"bytes"
	"os"
	"server/model"
	"server/persistence"
	"server/test"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func makeStringPointer(str string) *string {
	return &str
}

func TestExtendedDeployment(t *testing.T) {
	props := armresources.DeploymentPropertiesExtended{}
	state := armresources.ProvisioningStateSucceeded
	timestamp := time.Now()
	output := map[string]string{"value": "val", "type": "string"}
	outputs := map[string]interface{}{"name": output}
	props.ProvisioningState = &state
	props.CorrelationID = makeStringPointer("12345")
	props.Duration = makeStringPointer("PT1M1S")
	props.Timestamp = &timestamp
	props.Outputs = outputs
	dep := armresources.DeploymentExtended{
		ID:         makeStringPointer("/subscription/blah/blah/dummy"),
		Properties: &props,
	}
	depResult := model.NewDeploymentResult(dep)
	assert.Equal(t, "Succeeded", string(depResult.Status))
	assert.Equal(t, "12345", depResult.CorrelationID)
	assert.Equal(t, string(state), depResult.ProvisioningState)
	assert.Equal(t, string(armresources.ProvisioningStateSucceeded), depResult.ProvisioningState)
	assert.Equal(t, "PT1M1S", depResult.Duration)
	assert.Equal(t, timestamp, depResult.Timestamp)
	assert.Equal(t, "/subscription/blah/blah/dummy", depResult.ID)
	assert.Equal(t, outputs, depResult.Outputs)

	// Test non-Succeeded result
	state = armresources.ProvisioningStateCanceled
	props.ProvisioningState = &state
	depResult = model.NewDeploymentResult(dep)
	assert.Equal(t, "Canceled", string(depResult.Status))

	// Test no outputs
	props.Outputs = nil
	depResult = model.NewDeploymentResult(dep)
	assert.Equal(t, make(map[string]interface{}), depResult.Outputs)
}

func TestDurationParsing(t *testing.T) {
	durations := map[string]string{
		"PT1S":      "1 seconds",
		"PT2M1S":    "2 minutes 1 seconds",
		"PT1H10M6S": "1 hours 10 minutes 6 seconds",
		"PT22.674S": "22.67 seconds",
		"bogus":     "bogus",
	}
	for k, v := range durations {
		assert.Equal(t, v, model.GetAzureTimeFormatted(k))
	}
}

func TestUpdateExecution(t *testing.T) {
	execution := model.Execution{
		ResumeToken: "token",
	}

	// No result to check
	model.UpdateExecution(&execution, nil, "")
	assert.Equal(t, model.Failed, execution.Status)
	assert.Equal(t, "", execution.ResumeToken, "Resume token should be removed")
	execution.ResumeToken = "token"
	now := time.Now()
	result := model.DeploymentResult{
		ID:                "ID",
		Status:            "Failed",
		Timestamp:         now,
		Duration:          "PT1S",
		CorrelationID:     "12345",
		ProvisioningState: string(armresources.ProvisioningStateFailed),
	}

	// With result
	model.UpdateExecution(&execution, &result, "")
	assert.Equal(t, "", execution.ResumeToken, "Resume token should be removed")
	assert.Equal(t, "ID", execution.DeploymentID)
	assert.Equal(t, model.Failed, execution.Status)
	assert.Equal(t, "1 seconds", execution.Duration)
	assert.Equal(t, now, execution.Timestamp)
	assert.Equal(t, "12345", execution.CorrelationID)

	// With error JSON
	errJson := `{
	  "status": "Failed",
	  "error": {
		"code": "DeploymentFailed",
		"message": "At least one resource deployment...",
		"details": [
		  {
			"code": "BadRequest",
			"message": "{\r\n  \"error\": {\r\n    \"code\": \"MultipleErrorsOccurred\",\r\n    \"message\": \"trimmed\",\r\n    \"details\": [\r\n      {\r\n        \"code\": \"trimmed\",\r\n        \"message\": \"trimmed\"\r\n      },\r\n      {\r\n        \"code\": \"trimmed\",\r\n        \"message\": \"trimmed\"\r\n      }\r\n    ]\r\n  }\r\n}"
		  }
		]
	  }
	}`
	model.UpdateExecution(&execution, &result, errJson)
	assert.Equal(t, "DeploymentFailed", execution.Code)
	assert.Contains(t, execution.Error, "At least one")
	assert.Contains(t, execution.ErrorDetails, "BadRequest")
	assert.Contains(t, execution.ErrorDetails, "MultipleErrorsOccurred")

	// Bad JSON
	errJson = `{this"won't:"parse}`
	var buf bytes.Buffer
	log.SetOutput(&buf)
	model.UpdateExecution(&execution, &result, errJson)
	log.SetOutput(os.Stdout)
	assert.Contains(t, execution.Error, "invalid character")
	assert.Contains(t, buf.String(), "Unable to parse")
}

func TestSegmentPublisher(t *testing.T) {

	db := persistence.NewInMemoryDB()
	model.SetMetric(db.Instance, model.DeployStatus, "SUCCESS")
	model.SetMetric(db.Instance, model.AccessType, "PRIVATE")
	model.SetMetric(db.Instance, model.ApplicationId, "XYZ123")
	model.PublishToSegment(db.Instance)
}

// TestMain wraps the tests.  Setup is done before the call to m.Run() and any
// needed teardown after that.
func TestMain(m *testing.M) {
	test.SetEnvironment()
	m.Run()
}
