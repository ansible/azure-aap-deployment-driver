package model

import (
	"encoding/json"

	"time"

	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DryRunStepName = "Deployment__Readiness__Check"
const ModmDeploymentName = "ansible-on-azure"

// Replicate GORM base model, hiding times from json
type BaseModel struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-" gorm:"index"`
}

type Step struct {
	BaseModel
	Name       string            `gorm:"unique" json:"name"`
	Executions []Execution       `json:"executions" gorm:"constraint:OnUpdate:CASCADE;"`
	StageId    string            `json:"-"`
}

type Output struct {
	BaseModel
	ModuleName string `json:"moduleName"`
	Values     datatypes.JSONMap
}

type Execution struct {
	BaseModel
	Status            ExecutionStatus `json:"status" gorm:"type:string"`
	StepID            uint            `json:"stepId"`
	DeploymentID      string          `json:"-"`
	Error             string          `json:"error"`
	ErrorDetails      string          `json:"errorDetails"`
	Code              string          `json:"code"`
	ProvisioningState string          `json:"provisioningState"`
	Details           string          `json:"details"`
	Timestamp         time.Time       `json:"timestamp"`
	Duration          string          `json:"duration"`
	CorrelationID     string          `json:"correlationId"`
	ResumeToken       string          `json:"-"`

	// the execution data of a dry run (received from modm)
	DryRunExecution *DryRunExecution `json:"dryRunExecution" gorm:"json"`
}

type DryRunExecution struct { // the execution instance of the dry run (received from modm)
	Id string `json:"operationId"`
	// momd deploymentId, different than the DeploymentID in Execution
	DeploymentId uint   `json:"deploymentId"`
	Status       string `json:"status"`

	// the errors captured from the dry run (received from modm) if the status (on this struct) is failed
	Errors []sdk.DryRunError `json:"errors"`
}

type Status struct {
	BaseModel
	TemplatesLoaded   bool
	MainOutputsLoaded bool
	IsFatalState      bool
	FirstStart        time.Time
}

type SessionConfig struct {
	BaseModel
	SessionAuthKey []byte
}

type Telemetry struct {
	BaseModel
	MetricName  DeploymentMetric `gorm:"type:string"`
	MetricValue string
	Step        string
}

func UpdateExecution(execution *Execution, result *DeploymentResult, errJson string) {
	execution.ResumeToken = ""

	if result != nil {
		// Failed during deployment
		execution.Status = result.Status
		execution.DeploymentID = result.ID
		execution.CorrelationID = result.CorrelationID
		if result.Duration != "" {
			execution.Duration = GetAzureTimeFormatted(result.Duration)
		}
		execution.Timestamp = result.Timestamp
		execution.ProvisioningState = result.ProvisioningState
	} else {
		// Failed before deployment was created
		execution.Status = Failed
	}

	if errJson != "" {
		errorStruct := ErroredDeployment{}
		err := json.Unmarshal([]byte(errJson), &errorStruct)
		if err != nil {
			log.Warnf("Unable to parse Azure error: %v", err)
			execution.Error = err.Error()
			return
		}
		execution.Error = errorStruct.Error.Message
		execution.ErrorDetails = errorStruct.Error.DetailString()
		execution.Code = errorStruct.Error.Code
	}
}

func CreateNewOutput(name string, result *DeploymentResult) *Output {
	return &Output{
		ModuleName: name,
		Values:     result.Outputs,
	}
}

// Setter function for each deployment metric
func SetMetric(db *gorm.DB, metric DeploymentMetric, value string, step string) {
	db.Create(&Telemetry{
		MetricName:  metric,
		MetricValue: value,
		Step:        step,
	})
}

// Getter function for each deployment metric
func Metric(db *gorm.DB, metric DeploymentMetric) Telemetry {
	telemetry := Telemetry{}
	db.Where("metric_name = ?", metric).Find(&telemetry)
	return telemetry
}
