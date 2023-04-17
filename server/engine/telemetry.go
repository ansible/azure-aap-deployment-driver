package engine

import (
	"fmt"
	"server/config"
	"server/model"
	"time"

	"github.com/segmentio/analytics-go/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetEvent(propertiesMap analytics.Properties) string {

	switch status := propertiesMap[string(model.DeployStatus)]; status {
	case "succeeded":
		return "aap.azure.installer-deploy-success"
	case "failed":
		return "aap.azure.installer-deploy-failed"
	case "canceled":
		return "aap.azure.installer-deploy-cancel"
	}
	return ""
}

func BuildSegmentPropertiesMap(db *gorm.DB) analytics.Properties {

	var propertiesMap = analytics.Properties{}
	var metricData []model.Telemetry
	db.Find(&metricData)
	for _, data := range metricData {
		propertiesMap[string(data.MetricName)] = fmt.Sprintf("%v", data.MetricValue)
	}
	return propertiesMap
}

func GetMetricFromMainOutputs(db *gorm.DB) {

	//var outputsMap map[string]interface{}
	var outputs []model.Output
	db.Find(&outputs)
	for _, data := range outputs {
		if data.ModuleName == "" {
			location := data.Values["location"].(map[string]interface{})["value"]
			accessType := data.Values["access"].(map[string]interface{})["value"]
			model.SetMetric(db, model.Region, location.(string))
			model.SetMetric(db, model.AccessType, accessType.(string))
		}
	}
}

func PublishToSegment(db *gorm.DB) {

	writeKey := config.GetEnvironment().SEGMENT_WRITE_KEY
	if writeKey == "" {
		log.Errorf("Segment Write Key is missing : Not sending telemetry to Segment")
		return
	}
	// set metrics in DB that are not set yet
	model.SetMetric(db, model.ApplicationId, config.GetEnvironment().APPLICATION_ID)
	// time.RFC3339 format is the Go equivalent to ISO 8601 format (minus the milliseconds)
	model.SetMetric(db, model.EndTime, time.Now().Format(time.RFC3339))
	GetMetricFromMainOutputs(db)
	//gather all metrics in a property map
	propertiesMap := BuildSegmentPropertiesMap(db)
	eventName := GetEvent(propertiesMap)
	if eventName == "" {
		log.Errorf("Unexpected value for deploy status: [%v]. Not sending telemetry to Segment.", propertiesMap[string(model.DeployStatus)])
		return
	}

	client := analytics.New(writeKey)
	client.Enqueue(analytics.Track{
		UserId:     config.GetEnvironment().SUBSCRIPTION,
		Event:      eventName,
		Properties: propertiesMap,
	})
	client.Close()
}
