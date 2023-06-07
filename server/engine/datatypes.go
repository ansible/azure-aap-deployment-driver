package engine

import (
	"context"
	"server/model"
	"server/persistence"

	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

type Engine struct {
	context              context.Context
	database             *persistence.Database
	mainOutputs          *model.Output
	template             map[string]interface{}
	parameters           map[string]interface{}
	done                 chan struct{}
	status               *model.Status
	maxExecutionRestarts int
	modmClient           *sdk.Client
	modmDeploymentId     int
	deploymentStarted    bool
}
