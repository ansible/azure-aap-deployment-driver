package telemetry

import (
	"context"
	"fmt"
	"net/http"
	"server/config"
	"server/model"
	"server/segment"
	"server/util"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/analytics-go/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TelemetryHandler struct {
	client *segment.SegmentClient
	db     *gorm.DB
	ctx    context.Context
}

var once sync.Once
var handler TelemetryHandler

func Init(db *gorm.DB, ctx context.Context) *TelemetryHandler {
	// Get write key
	writeKey := config.GetEnvironment().SEGMENT_WRITE_KEY
	if writeKey == "" {
		log.Error("Segment Write Key is missing : Not sending telemetry to Segment")
		return nil
	}
	// Set metrics
	if config.GetEnvironment().APPLICATION_ID != "" {
		model.SetMetric(db, model.ApplicationId, config.GetEnvironment().APPLICATION_ID, model.MAIN_MARKER)
	}
	if config.GetEnvironment().START_TIME != "" {
		model.SetMetric(db, model.StartTime, config.GetEnvironment().START_TIME, model.MAIN_MARKER)
	}
	extractMainOutputs(db)

	once.Do(func() {
		handler = TelemetryHandler{
			client: segment.Init(writeKey, config.GetEnvironment().SUBSCRIPTION),
			db:     db,
			ctx:    ctx,
		}
	})
	return &handler
}

func (t *TelemetryHandler) FinalizeAndPublish() (*analytics.Track, error) {
	// Set final metrics
	model.SetMetric(t.db, model.EndTime, time.Now().Format(time.RFC3339), model.MAIN_MARKER)
	setStepMetrics(t.db)
	err := sendDeploymentIdentification(t.ctx)
	if err != nil {
		log.Errorf("Unable to send deployment identification event to Azure Function/Segment: %v", err)
	}

	// Set deployment status
	setDeploymentStatus(t.db)

	// Populate properties
	var metricData []model.Telemetry
	t.db.Find(&metricData)

	// Publish
	t.client.AddProperties(metricData)
	return t.client.Publish()
}

func setStepMetrics(db *gorm.DB) {
	var allSteps []model.Step
	db.Find(&allSteps)
	db.Model(&model.Step{}).Preload("Executions").Find(&allSteps)

	for _, step := range allSteps {
		storeRetriesPerStep(db, step)
		storeErrorsPerStep(db, step)
	}
}

func setDeploymentStatus(db *gorm.DB) {
	var steps []model.Step
	db.Find(&steps)

	status := model.DeploymentSucceeded
	for _, step := range steps {
		latestExecution := model.Execution{}

		var count int64
		db.Model(&model.Execution{}).Where("step_id = ?", step.ID).Count(&count)
		if count > 0 {
			db.Last(&latestExecution, "step_id = ?", step.ID)
		}
		if latestExecution.Status == model.PermanentlyFailed {
			status = model.DeploymentFailed
			break
		} else if latestExecution.Status == model.Canceled {
			status = model.DeploymentCanceled
			break
		}
	}
	model.SetMetric(db, model.DeployStatus, string(status), model.MAIN_MARKER)
}

func storeRetriesPerStep(db *gorm.DB, step model.Step) {
	// len(step.Executions) includes the 1st attempt as well
	// which should not be considered as a "retry"
	// In cases where a step is not executed at all, retries should be set to 0.
	retries := len(step.Executions) - 1
	if retries < 0 {
		retries = 0
	}
	model.SetMetric(db, model.Retries, fmt.Sprint(retries), step.Name)
}

func storeErrorsPerStep(db *gorm.DB, step model.Step) {
	// for errors, loop through all executions and combine all errors found for each step and then merge them all
	for _, execution := range step.Executions {
		if execution.ErrorDetails != "" {
			model.SetMetric(db, model.Errors, execution.ErrorDetails, step.Name)
		}
	}
}

func extractMainOutputs(db *gorm.DB) {
	var mainOutput model.Output
	db.Where("module_name = ?", "").Find(&mainOutput)

	if loc, exists := mainOutput.Values["location"]; exists {
		model.SetMetric(db, model.Region, loc.(map[string]interface{})["value"].(string), model.MAIN_MARKER)
	} else {
		log.Warn("Location of deployment is missing and will not be included in telemetry")
	}
	if access, exists := mainOutput.Values["access"]; exists {
		model.SetMetric(db, model.AccessType, access.(map[string]interface{})["value"].(string), model.MAIN_MARKER)
	} else {
		log.Warn("Access Type of deployment is missing and will not be included in telemetry")
	}
}

func sendDeploymentIdentification(ctx context.Context) error {
	if config.GetEnvironment().AZURE_MARKETPLACE_FUNCTION_KEY == "" {
		return nil
	}
	azureFunctionUrl := strings.Join([]string{config.GetEnvironment().AZURE_MARKETPLACE_FUNCTION_BASE_URL, config.GetEnvironment().AZURE_MARKETPLACE_FUNCTION_KEY}, "?code=")
	req := util.NewHttpRequester()
	body := make(map[string]interface{})
	body["subscriptionId"] = config.GetEnvironment().SUBSCRIPTION
	body["tenantId"] = config.GetEnvironment().AZURE_TENANT_ID
	body["applicationId"] = config.GetEnvironment().APPLICATION_ID
	body["eventType"] = "IDENTIFICATION"    // Needed for marketplace notification function
	body["provisioningState"] = "Succeeded" // ditto
	resp, err := req.MakeRequestWithJSONBody(ctx, http.MethodPost, azureFunctionUrl, nil, body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d, body: %s", resp.StatusCode, resp.Body)
	}
	return nil
}

func IncrementLogins(db *gorm.DB) {
	loginsCount := "1" // default value for case if no logins were recorded yet

	// get existing logins telemetry
	telemetry := model.Metric(db, model.UserLogins)
	// if any logins were already recorded, set logins value incremented by 1
	if telemetry.MetricValue != "" {
		n, err := strconv.Atoi(telemetry.MetricValue)
		if err == nil {
			loginsCount = strconv.Itoa(n + 1)
		}
	}
	// store new logins count
	model.SetMetric(db, model.UserLogins, loginsCount, model.MAIN_MARKER)
}
