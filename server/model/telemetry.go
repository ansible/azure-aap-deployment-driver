package model

import (
	"fmt"

	"server/config"

	"github.com/segmentio/analytics-go/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// event name - required for Amplitude
const event string = "Deployment Completed"

// This slice should have the same metrics as /server/model/enums.go DeploymentMetrics
var metrics = []DeploymentMetric{
	StartTime,
	EndTime,
	CustomerSubscriptionID,
	Region,
	AccessType,
	DeployStatus,
	Errors,
	Retries,
}

func BuildSegmentPropertiesMap(db *gorm.DB) analytics.Properties {

	var propertiesMap = analytics.Properties{}
	for _, metric := range metrics {
		propertiesMap[string(metric)] = fmt.Sprintf("%v", Metric(db, metric).MetricValue)
	}
	return propertiesMap
}

func PublishToSegment(db *gorm.DB) {

	log.Info("Starting to Publish")
	client := analytics.New(config.GetEnvironment().SEGMENT_WRITE_KEY)
	log.Info("Found client", client)
	// set metrics in DB that are not set yet
	// TODO : Is there a better place where subscriptionId can be set? Not possible in env.go because of circular imports
	SetMetric(db, CustomerSubscriptionID, config.GetEnvironment().SUBSCRIPTION)
	propertiesMap := BuildSegmentPropertiesMap(db)
	log.Info("Map details : \n", propertiesMap)
	client.Enqueue(analytics.Track{
		UserId:     propertiesMap[string(CustomerSubscriptionID)].(string),
		Event:      event,
		Properties: propertiesMap,
	})

	client.Close()
}
