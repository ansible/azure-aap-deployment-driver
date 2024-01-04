package telemetry_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"server/model"
	"server/persistence"
	"server/segment"
	"server/telemetry"
	"server/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTelemetryHandler(t *testing.T) {
	db := persistence.NewInMemoryDB()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	// Set env vars
	os.Setenv("START_TIME", "something")
	os.Setenv("SEGMENT_WRITE_KEY", "DUMMY")
	os.Setenv("AZURE_MARKETPLACE_FUNCTION_BASE_URL", server.URL)
	os.Setenv("AZURE_MARKETPLACE_FUNCTION_KEY", "DUMMY")
	os.Setenv("APPLICATION_ID", "DUMMY")
	test.SetEnvironment()

	// Create main outputs
	mainoutputs := "{\"access\":{\"type\":\"String\",\"value\":\"public\"},\"location\":{\"type\":\"String\",\"value\":\"eastus\"}}"
	outputValues := make(map[string]interface{})
	err := json.Unmarshal([]byte(mainoutputs), &outputValues)
	assert.Nil(t, err)

	outputs := model.Output{
		ModuleName: "",
		Values:     outputValues,
	}
	db.Instance.Save(&outputs)

	// Create steps in DB
	step := model.Step{
		Name: "first",
	}
	db.Instance.Create(&step)
	exec := model.Execution{
		Status:       model.PermanentlyFailed,
		StepID:       step.ID,
		ErrorDetails: "ERROR",
	}
	db.Instance.Create(&exec)
	th := telemetry.Init(db.Instance, context.Background())
	c := segment.Init("DUMMY", "subscription")
	th.TestSetClient(c)
	track, err := th.FinalizeAndPublish()
	assert.Nil(t, err)
	require.NotNil(t, track)
	assert.Equal(t, "aap.azure.installer-deploy-failed", track.Event)
	assert.Equal(t, "12345678-90ab-cdef-0123-4567890abcde", track.UserId)
	// TODO check properties
}

func TestIncrementLogins(t *testing.T) {
	db := persistence.NewInMemoryDB()
	telemetry.IncrementLogins(db.Instance)
	logins := model.Metric(db.Instance, model.UserLogins)
	assert.Equal(t, "1", logins.MetricValue)

	telemetry.IncrementLogins(db.Instance)
	logins = model.Metric(db.Instance, model.UserLogins)
	assert.Equal(t, "2", logins.MetricValue)
}
