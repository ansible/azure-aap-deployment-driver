package main

import (
	"server/api"
	"server/azure"
	"server/config"
	"server/controllers"
	"server/engine"
	"server/persistence"
)

func main() {
	config.ParseArgs()
	config.ConfigureLogging()

	db := persistence.NewPersistentDB(config.GetEnvironment().DB_PATH)
	// TODO store first start up in DB so we can determine max allowed run time for installer

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
