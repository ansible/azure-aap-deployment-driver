package model

import (
	"fmt"
	"server/config"

	"github.com/segmentio/analytics-go/v3"
	"gorm.io/gorm"
)

func GetEvent(propertiesMap analytics.Properties) string {

	switch status := propertiesMap[string(DeployStatus)]; status {
	case "succeeded":
		return "aap.azure.installer-deploy-success"
	case "failed":
		return "aap.azure.installer-deploy-failed"
	case "cancelled":
		return "aap.azure.installer-deploy-cancel"
	}
	return ""
}

func BuildSegmentPropertiesMap(db *gorm.DB) analytics.Properties {

	var propertiesMap = analytics.Properties{}
	var metricData []Telemetry
	db.Find(&metricData)
	for _, data := range metricData {
		propertiesMap[string(data.MetricName)] = fmt.Sprintf("%v", data.MetricValue)
	}
	return propertiesMap
}

func PublishToSegment(db *gorm.DB) {

	client := analytics.New(config.GetEnvironment().SEGMENT_WRITE_KEY)
	// set metrics in DB that are not set yet
	SetMetric(db, ApplicationId, config.GetEnvironment().APPLICATION_ID)
	//gather all metrics in a property map
	propertiesMap := BuildSegmentPropertiesMap(db)
	client.Enqueue(analytics.Track{
		UserId:     config.GetEnvironment().SUBSCRIPTION,
		Event:      GetEvent(propertiesMap),
		Properties: propertiesMap,
	})

	client.Close()
}
