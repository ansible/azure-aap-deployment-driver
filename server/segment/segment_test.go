package segment_test

import (
	"net/http"
	"net/http/httptest"
	"server/model"
	"server/segment"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSegmentClient(t *testing.T) {
	s := segment.Init("bogus", "subscription")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	s.SetEndpoint(server.URL)
	props := make([]model.Telemetry, 5)
	props[0] = model.Telemetry{
		MetricName:  "key",
		MetricValue: "val",
		Step:        "first",
	}
	props[1] = model.Telemetry{
		MetricName:  "mainprop",
		MetricValue: "mainval",
	}
	props[2] = model.Telemetry{
		MetricName:  "errors",
		MetricValue: "error1",
		Step:        "failed",
	}
	props[3] = model.Telemetry{
		MetricName:  "errors",
		MetricValue: "error2",
		Step:        "failed",
	}
	props[4] = model.Telemetry{
		MetricName:  "deploystatus",
		MetricValue: "failed",
		Step:        model.MAIN_MARKER,
	}
	s.AddProperties(props)
	track, err := s.Publish()
	assert.Nil(t, err)
	require.NotNil(t, track)
	assert.Equal(t, "aap.azure.installer-deploy-failed", track.Event)
	assert.Equal(t, "subscription", track.UserId)
	assert.Equal(t, []string{"error1", "error2"}, track.Properties["errors"].(map[string][]string)["failed"])
	assert.Nil(t, err)
	s.Close()
}
