package model

import (
	"github.com/segmentio/analytics-go/v3"
	"gorm.io/gorm"
)

const event string = "Deployment Completed"

// TODO: How to store the write key?
func PublishToSegment(db *gorm.DB) {
	client := analytics.New("N1Dfx7gVokm42dmkJOr5GkNKhRmLjF7i")
	defer client.Close()

	subscriptionID := Metric(db, CustomerSubscriptionID).MetricValue
	accessType := Metric(db, AccessType).MetricValue
	deployStatus := Metric(db, DeployStatus).MetricValue
	endTime := Metric(db, EndTime).MetricValue
	errors := Metric(db, Errors).MetricValue
	region := Metric(db, Region).MetricValue
	startTime := Metric(db, StartTime).MetricValue
	retries := Metric(db, Retries).MetricValue

	client.Enqueue(analytics.Track{
		UserId: subscriptionID,
		Event:  event,
		Properties: map[string]interface{}{
			string(DeployStatus): deployStatus,
			string(AccessType):   accessType,
			string(EndTime):      endTime,
			string(StartTime):    startTime,
			string(Region):       region,
			string(Errors):       errors,
			string(Retries):      retries,
		},
	})
}
