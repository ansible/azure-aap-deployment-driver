package engine

import (
	"context"
	"server/model"
	"server/persistence"
	"server/telemetry"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type Engine struct {
	context              context.Context
	database             *persistence.Database
	resolver             *Resolver
	mainOutputs          *model.Output
	done                 chan struct{}
	status               *model.Status
	maxExecutionRestarts int
	deploymentsClient    *armresources.DeploymentsClient
	telemetryHandler     *telemetry.TelemetryHandler
}
