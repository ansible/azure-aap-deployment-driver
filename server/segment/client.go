package segment

import (
	"fmt"
	"server/model"
	"sync"

	"github.com/segmentio/analytics-go/v3"
)

type SegmentClient struct {
	client       analytics.Client
	writeKey     string
	subscription string
	props        analytics.Properties
}

var once sync.Once
var client SegmentClient

func Init(writeKey string, subscription string) *SegmentClient {
	once.Do(func() {
		client = SegmentClient{
			client:       analytics.New(writeKey),
			writeKey:     writeKey,
			subscription: subscription,
			props:        analytics.Properties{},
		}
	})
	return &client
}

func (s *SegmentClient) AddProperties(metrics []model.Telemetry) {
	retriesMap := make(map[string]interface{})
	errorsMap := make(map[string][]string)

	// Error Details and Number of retries will be granular at the step level hence, each error and retry will be mapped to a step
	// in a nested JSON like format. Rest of the metrics are granular at the deployment level
	for _, data := range metrics {
		if data.Step == model.MAIN_MARKER {
			s.props[string(data.MetricName)] = data.MetricValue
		} else {
			switch data.MetricName {
			case model.Retries:
				retriesMap[data.Step] = data.MetricValue
			case model.Errors:
				errorsMap[data.Step] = append(errorsMap[data.Step], data.MetricValue)
			}
		}
	}
	s.props[string(model.Retries)] = retriesMap
	s.props[string(model.Errors)] = errorsMap
}

func (s *SegmentClient) createTrack(eventName string) analytics.Track {
	return analytics.Track{
		UserId:     s.subscription,
		Event:      eventName,
		Properties: s.props,
	}
}

func (s *SegmentClient) Publish() (*analytics.Track, error) {
	eventName := determineEvent(s.props)
	if eventName == "" {
		return nil, fmt.Errorf("unable to determine telemetry event, not publishing")
	}
	track := s.createTrack(eventName)
	err := track.Validate()
	if err != nil {
		return nil, err
	}
	return &track, s.client.Enqueue(track)
}

func (s *SegmentClient) Close() {
	s.client.Close()
}

func determineEvent(props analytics.Properties) string {
	switch status := props[string(model.DeployStatus)]; status {
	case "succeeded":
		return "aap.azure.installer-deploy-success"
	case "failed":
		return "aap.azure.installer-deploy-failed"
	case "canceled":
		return "aap.azure.installer-deploy-cancel"
	}
	return ""
}
