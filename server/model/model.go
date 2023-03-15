package model

import (
	"encoding/json"

	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	_ "gorm.io/driver/sqlite"
)

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
	Template   datatypes.JSONMap `json:"-"`
	Parameters datatypes.JSONMap `json:"-"`
	Priority   uint              `json:"order"`
	Executions []Execution       `json:"executions" gorm:"constraint:OnUpdate:CASCADE;"`
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
	ExecutionCount    int             `json:"executionCount"`
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

func UpdateExecution(execution *Execution, result *DeploymentResult, errJson string) {
	execution.ExecutionCount++
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
