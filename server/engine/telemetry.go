package engine

import (
	"fmt"
	"server/config"
	"server/model"
	"strings"
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
	retriesMap := make(map[string]interface{})
	errorsMap := make(map[string]interface{})
	var metricData []model.Telemetry
	db.Find(&metricData)
	// Error Details and Number of retries will be granular at the step level hence, each error and retry will be mapped to a step
	// in a nested JSON like format
	// Rest of the metrics are granular at the deployment level
	for _, data := range metricData {
		fmt.Printf("%+v\n", data)
		if data.Step == "" {
			propertiesMap[string(data.MetricName)] = fmt.Sprintf("%v", data.MetricValue)
		} else {
			switch metric := data.MetricName; metric {
			case model.Retries:
				retriesMap[data.Step] = fmt.Sprintf("%v", data.MetricValue)
			case model.Errors:
				errorsMap[data.Step] = fmt.Sprintf("%v", data.MetricValue)
			}
		}
	}
	propertiesMap[string(model.Retries)] = retriesMap
	propertiesMap[string(model.Errors)] = errorsMap

	return propertiesMap
}

func StoreMetricFromMainOutputs(db *gorm.DB) {

	var mainOutput model.Output
	db.Where("module_name = ?", "").Find(&mainOutput)

	if loc, exists := mainOutput.Values["location"]; exists {
		model.SetMetric(db, model.Region, loc.(map[string]interface{})["value"].(string), "")
	} else {
		log.Error("Location of deployment is missing : will not be included in telemetry")
	}
	if access, exists := mainOutput.Values["access"]; exists {
		model.SetMetric(db, model.AccessType, access.(map[string]interface{})["value"].(string), "")
	} else {
		log.Error("Access Type of deployment is missing : will not be included in telemetry")
	}
}

func StoreRetriesPerStep(db *gorm.DB, step model.Step) {

	retries := len(step.Executions) - 1
	if retries < 0 {
		retries = 0
	}
	model.SetMetric(db, model.Retries, fmt.Sprint(retries), step.Name)

}

func StoreErrorsPerStep(db *gorm.DB, step model.Step) {

	var errorsArr []string
	var allErrors string

	for _, execution := range step.Executions {
		if execution.ErrorDetails != "" {
			errorsArr = append(errorsArr, execution.ErrorDetails)
		}
	}
	allErrors = strings.Join(errorsArr, "\n")
	if allErrors != "" {
		model.SetMetric(db, model.Errors, fmt.Sprint(allErrors), step.Name)
	}
}

func StoreMetricsPerStep(db *gorm.DB) {

	// Store number of retries and error details for each step if present
	var allSteps []model.Step
	db.Find(&allSteps)
	db.Model(&model.Step{}).Preload("Executions").Find(&allSteps)

	for _, step := range allSteps {
		// len(step.Executions) includes the 1st attempt as well
		// which should not be considered as a "retry"
		// In cases where a step is not executed at all, retries should be set to 0.
		// for errors, loop through all executions and combine all errors found for each step and then merge them all
		StoreRetriesPerStep(db, step)
		StoreErrorsPerStep(db, step)

	}
}

func PublishToSegment(db *gorm.DB) {

	writeKey := config.GetEnvironment().SEGMENT_WRITE_KEY
	if writeKey == "" {
		log.Errorf("Segment Write Key is missing : Not sending telemetry to Segment")
		return
	}
	// set metrics in DB that are not set yet
	model.SetMetric(db, model.ApplicationId, config.GetEnvironment().APPLICATION_ID, "")
	// time.RFC3339 format is the Go equivalent to ISO 8601 format (minus the milliseconds)
	model.SetMetric(db, model.EndTime, time.Now().Format(time.RFC3339), "")
	StoreMetricFromMainOutputs(db)
	StoreMetricsPerStep(db)
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
