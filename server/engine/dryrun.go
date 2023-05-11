package engine

import (
	"sync"

	"gorm.io/gorm"
)

var (
	dryRunInstance     *dryRunController
	dryRunInstanceOnce sync.Once
	dryRunInstanceErr  error
)

type ErrorHandler func(err error)

type dryRunController struct {
	deploymentId           int
	db                     *gorm.DB
	done                   chan struct{}
	clientEndpoint         string
	location               string
	resourceGroup          string
	subscription           string
	apiKey                 string
	hookName               string
	callbackClientEndpoint string
	HandleError				ErrorHandler
}
