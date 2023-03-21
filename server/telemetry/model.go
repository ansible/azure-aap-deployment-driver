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

// Getter Setters for each DeploymentMetric
func SetMetric(db *gorm.DB, telemetry *Telemetry, metric DeploymentMetric, status string) {
	telemetry.MetricName = metric
	telemetry.MetricValue = status
	db.Save(&telemetry)
}

func Metric(db *gorm.DB, metric DeploymentMetric) *gorm.DB {

	return db.Where("MetricName = ?", metric)
}
