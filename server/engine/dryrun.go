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
	deploymentId   int
	db             *gorm.DB
	done           chan struct{}
	clientEndpoint string
	location       string
	resourceGroup  string
	subscription   string
	apiKey         string
	hookName       string

	// this is the url that will be called by MODM. It maps to /eventhook route for handler/eventhook
	eventHookCallbackUrl string
	HandleError          ErrorHandler
}
