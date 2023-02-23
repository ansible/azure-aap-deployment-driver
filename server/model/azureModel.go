package model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

const (
	// AZURE_ERROR_WITH_JSON_REGEXP matches json string embedded in errors, surrounded by dashes
	AZURE_ERROR_WITH_JSON_REGEXP = `(?s)\-+\s*(\{.+\})\s*\-+` // the first () specifies flags, second one is capturing group
	// ISO8601_TIME_PERIOD_REGEXP matches time portion ISO 8601 standard
	ISO8601_TIME_PERIOD_REGEXP = `PT((\d+)?H)?((\d+)?M)?((\d+\.?\d{0,2})\d*?S)?` // the double groups ensure proper parsing when not all elements are present
)

var (
	azureErrorWithJSONRegexp = regexp.MustCompile(AZURE_ERROR_WITH_JSON_REGEXP)
	azureTimestampRegexp     = regexp.MustCompile(ISO8601_TIME_PERIOD_REGEXP)
)

type valueMap struct {
	Value string `json:"value"`
}

type deploymentParameters struct {
	Location valueMap `json:"location"`
	Name     valueMap `json:"name"`
}

type deploymentGenerator struct {
	Name         string `json:"name"`
	TemplateHash string `json:"templateHash"`
	Version      string `json:"version"`
}

type deploymentMetadata struct {
	Generator deploymentGenerator `json:"_generator"`
}

type deploymentSubTemplate struct {
	Schema         string                 `json:"$schema"`
	ContentVersion string                 `json:"contentVersion"`
	Metadata       deploymentMetadata     `json:"metadata"`
	Outputs        map[string]interface{} `json:"outputs"`    // Can be anything here
	Parameters     map[string]interface{} `json:"parameters"` // Complex structure, don't care now but flesh out if needed
	Resources      map[string]interface{} `json:"resources"`  // Complex structure, don't care now but flesh out if needed
}

type deploymentProperties struct {
	ExpressionEvaluationOptions valueMap             `json:"expressionEvaluationOptions"`
	Mode                        string               `json:"mode"`
	Parameters                  deploymentParameters `json:"parameters"`
}

type deploymentResources struct {
	ApiVersion string                `json:"apiVersion"`
	Name       string                `json:"name"`
	Properties deploymentProperties  `json:"properties"`
	Location   string                `json:"location"`
	Type       string                `json:"type"`
	Template   deploymentSubTemplate `json:"template"`
}

type DeploymentTemplate struct {
	Schema         string                `json:"$schema"`
	ContentVersion string                `json:"contentVersion"`
	Metadata       deploymentMetadata    `json:"metadata"`
	Resources      []deploymentResources `json:"resources"`
}

// Errors that come back from a deployment failure in the error message
type deploymentErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (d *deploymentError) DetailString() string {
	var out strings.Builder
	for _, detail := range d.Details {
		out.WriteString(fmt.Sprintf("Code: %s - Message: %s\n", detail.Code, detail.Message))
	}
	return out.String()
}

type deploymentError struct {
	Code    string                  `json:"code"`
	Message string                  `json:"message"`
	Details []deploymentErrorDetail `json:"details"`
}
type ErroredDeployment struct {
	Status string          `json:"status"`
	Error  deploymentError `json:"error"`
}

type DeploymentResult struct {
	ID                string                 `json:"id"`
	CorrelationID     string                 `json:"correlationId"`
	Duration          string                 `json:"duration"`
	Timestamp         time.Time              `json:"timestamp"`
	ProvisioningState string                 `json:"provisioningState"`
	Outputs           map[string]interface{} `json:"outputs"`
	Status            ExecutionStatus
}

func NewDeploymentResult(response armresources.DeploymentExtended) *DeploymentResult {
	var status ExecutionStatus
	switch *response.Properties.ProvisioningState {
	case armresources.ProvisioningStateSucceeded:
		status = Succeeded
	case armresources.ProvisioningStateCanceled:
		status = Canceled
	default:
		status = Failed
	}
	// make sure response outputs are always there, even if empty
	var responseOutputs map[string]interface{}
	if response.Properties.Outputs != nil {
		responseOutputs = response.Properties.Outputs.(map[string]interface{})
	} else {
		responseOutputs = make(map[string]interface{})
	}
	res := DeploymentResult{}
	if response.Properties.ProvisioningState != nil {
		res.ProvisioningState = string(*response.Properties.ProvisioningState)
	}
	if response.ID != nil {
		res.ID = *response.ID
	}
	if response.Properties.CorrelationID != nil {
		res.CorrelationID = *response.Properties.CorrelationID
	}
	if response.Properties.Duration != nil {
		res.Duration = *response.Properties.Duration
	}
	if response.Properties.Timestamp != nil {
		res.Timestamp = *response.Properties.Timestamp
	}
	res.Status = status
	res.Outputs = responseOutputs
	return &res
}

func GetAzureErrorJSONString(err error) string {
	errString := err.Error()
	jsonMatches := azureErrorWithJSONRegexp.FindStringSubmatch(errString)
	if len(jsonMatches) == 2 { // there should be two matches, one for full match and one for the group
		return jsonMatches[1]
	}
	cleanString := strconv.Quote(err.Error())
	return fmt.Sprintf("{\"status\":\"Failed\",\"error\":{\"message\":%s}}", cleanString)
}

// Returns duration as a human readable string
func GetAzureTimeFormatted(duration string) string {
	matches := azureTimestampRegexp.FindStringSubmatch(duration)
	if len(matches) == 7 { // checking for 7 because there are 6 groups + 1 full match
		var hours, minutes, seconds string
		if matches[2] != "" {
			hours = matches[2] + " hours "
		}
		if matches[4] != "" {
			minutes = matches[4] + " minutes "
		}
		if matches[6] != "" {
			seconds = matches[6] + " seconds"
		}
		return hours + minutes + seconds
	}
	return duration
}
