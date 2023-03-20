package telemetry

import "database/sql/driver"

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

// Allow GORM to store the string value of DeploymentMetric
func (u *DeploymentMetric) Scan(value interface{}) error {
	*u = DeploymentMetric(value.(string))
	return nil
}

func (u DeploymentMetric) Value() (driver.Value, error) {
	return string(u), nil
}
