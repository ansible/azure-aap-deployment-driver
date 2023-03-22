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
func SetMetric(db *gorm.DB, metric DeploymentMetric, value string) {

	row := Telemetry{
		MetricName:  metric,
		MetricValue: value,
	}
	db.Create(&row)
}

func Metric(db *gorm.DB, metric DeploymentMetric) Telemetry {

	telemetry := Telemetry{}
	db.Where("metric_name = ?", metric).Find(&telemetry)
	return telemetry
}
