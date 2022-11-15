package engine

import (
	"context"
	"server/model"
	"server/persistence"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type Engine struct {
	context              context.Context
	database             *persistence.Database
	resolver             *Resolver
	mainOutputs          *model.Output
	done                 chan struct{}
	maxExecutionRestarts int
	deploymentsClient    *armresources.DeploymentsClient
}
