package telemetry

import (
	"errors"
	"server/config"

	"github.com/segmentio/analytics-go/v3"
)

func SendEntitlementResult(customerRhOrgId string, registeredRhOrgId string) error {
	writeKey := config.GetEnvironment().SEGMENT_WRITE_KEY
	if writeKey == "" {
		return errors.New("segment write key not set")
	}

	props := analytics.Properties{}
	props.Set("azureSubscriptionId", config.GetEnvironment().SUBSCRIPTION)
	props.Set("azureTenantId", config.GetEnvironment().AZURE_TENANT_ID)
	props.Set("applicationId", config.GetEnvironment().APPLICATION_ID)
	props.Set("customerRedHatOrgId", customerRhOrgId)
	props.Set("registeredRedHatOrgId", registeredRhOrgId)

	var event string
	if registeredRhOrgId == "" {
		event = "aap.azure.entitlement-failure"
	} else if registeredRhOrgId == customerRhOrgId {
		event = "aap.azure.entitlement-new"
	} else {
		event = "aap.azure.entitlement-existing"
	}

	client := analytics.New(writeKey)
	defer client.Close()

	track := analytics.Track{
		UserId:     config.GetEnvironment().APPLICATION_ID,
		Event:      event,
		Properties: props,
	}
	return client.Enqueue(track)
}
