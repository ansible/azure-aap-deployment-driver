package segment

import (
	"github.com/segmentio/analytics-go/v3"
	"github.com/sirupsen/logrus"
)

func (s *SegmentClient) TestSetEndpoint(endpoint string) {
	cfg := analytics.Config{
		Endpoint: endpoint,
	}
	c, err := analytics.NewWithConfig(s.writeKey, cfg)
	if err != nil {
		logrus.Errorf("Error applying test endpoint to segment client: %v", err)
	}
	s.client = c
}
