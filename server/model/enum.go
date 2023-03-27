package model

import "database/sql/driver"

type ExecutionStatus string

const (
	Started           ExecutionStatus = "Started"
	Failed            ExecutionStatus = "Failed"
	PermanentlyFailed ExecutionStatus = "PermanentlyFailed"
	Succeeded         ExecutionStatus = "Succeeded"
	Restart           ExecutionStatus = "Restart"
	Restarted         ExecutionStatus = "Restarted"
	RestartTimedOut   ExecutionStatus = "RestartTimedOut"
	Canceled          ExecutionStatus = "Canceled"
)

type DeploymentMetric string

const (
	StartTime              DeploymentMetric = "StartTime"
	EndTime                DeploymentMetric = "EndTime"
	CustomerSubscriptionID DeploymentMetric = "CustomerSubscriptionID"
	Region                 DeploymentMetric = "Region"
	AccessType             DeploymentMetric = "AccessType"
	DeployStatus           DeploymentMetric = "DeployStatus"
	Errors                 DeploymentMetric = "Errors"
	Retries                DeploymentMetric = "Retries"
)

// Allow GORM to store the string value of ExecutionStatus
func (u *ExecutionStatus) Scan(value interface{}) error {
	*u = ExecutionStatus(value.(string))
	return nil
}

func (u ExecutionStatus) Value() (driver.Value, error) {
	return string(u), nil
}

// Allow GORM to store the string value of DeploymentMetric
func (u *DeploymentMetric) Scan(value interface{}) error {
	*u = DeploymentMetric(value.(string))
	return nil
}

func (u DeploymentMetric) Value() (driver.Value, error) {
	return string(u), nil
}
