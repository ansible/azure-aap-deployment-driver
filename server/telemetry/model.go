package telemetry

import (
	"server/model"

	"gorm.io/gorm"
)

type Telemetry struct {
	model.BaseModel
	MetricName  DeploymentMetric `gorm:"type:string"`
	MetricValue string
	Step        string
}

func SetDeploymentStatus(db *gorm.DB, telemetry *Telemetry, status string) {
	telemetry.MetricName = DeployStatus
	telemetry.MetricValue = status
	db.Save(&telemetry)
}

func DeploymentStatus(db *gorm.DB) *gorm.DB {

	return db.Where("MetricName = ?", DeployStatus)
}
