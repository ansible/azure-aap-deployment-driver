package telemetry

import "server/segment"

func (t *TelemetryHandler) TestSetClient(client *segment.SegmentClient) {
	t.client = client
}
