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

// Allow GORM to store the string value of ExecutionStatus
func (u *ExecutionStatus) Scan(value interface{}) error {
	*u = ExecutionStatus(value.(string))
	return nil
}

func (u ExecutionStatus) Value() (driver.Value, error) {
	return string(u), nil
}
