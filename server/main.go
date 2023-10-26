package main

import (
	"context"
	"path/filepath"
	"server/api"
	"server/azure"
	"server/config"
	"server/controllers"
	"server/engine"
	"server/persistence"
	"server/util"

	log "github.com/sirupsen/logrus"
)

type APIFilter struct {
	VendorProductCode   string `json:"vendorProductCode,omitempty"`
	AzureSubscriptionId string `json:"azureSubscriptionId,omitempty"`
	AzureTenantId       string `json:"azureTenantId,omitempty"`
}

func main() {
	config.ConfigureLogging()
	config.ParseArgs()

	db := persistence.NewPersistentDB(filepath.Join(config.GetEnvironment().BASE_PATH, config.GetEnvironment().DB_REL_PATH))
	// TODO store first start up in DB so we can determine max allowed run time for installer

	cert := config.GetEnvironment().SW_SUB_API_CERTIFICATE
	key := config.GetEnvironment().SW_SUB_API_PRIVATEKEY

	requester, err := util.NewHttpRequesterWithCertificate(cert, key)
	if err != nil {
		log.Fatalf("Did not get client. %v", err)
	}

	ctx := context.Background()
	response, err := requester.MakeRequestWithJSONBody(
		ctx,
		"POST",
		"https://ibm-entitlement-gateway.api.redhat.com/v1/partnerSubscriptions",
		nil,
		APIFilter{
			VendorProductCode:   "rhaapomsa",
			AzureSubscriptionId: "5275bbe1-ea9f-4fca-a8ef-4ef79d5e5a0b",
			//AzureTenantId:       "97e9a597-d113-4e68-83d3-15cd863399c9",
		},
	)
	if err != nil {
		log.Fatalf("Failed to get response: %v", err)
	}
	log.Printf("Got response code: %d\n", response.StatusCode)
	log.Printf("%v", string(response.Body))

	// Instantiate Azure clients and session
	azure.EnsureAzureLogin(nil)
	deploymentsClient := azure.NewDeploymentsClient(nil)

	// Graceful exit handler
	exit := controllers.NewExitController()

	engine := engine.NewEngine(exit.Context(), db, deploymentsClient)

	app := api.NewApp(db, engine)

	// Start listening for shutdown signal
	exit.Start()

	// Start the engine
	go engine.Run()

	// Start the API server
	go app.Run()

	// Wait for either the engine being done or a signal received by exit controller
	select {
	case <-exit.Done():
	case <-engine.Done():
	}
}
