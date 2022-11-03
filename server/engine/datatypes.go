package engine

import (
	"context"
	"server/model"
	"server/persistence"
)

type Engine struct {
	context              context.Context
	database             *persistence.Database
	resolver             *Resolver
	mainOutputs          *model.Output
	done                 chan struct{}
	maxExecutionRestarts int
}
