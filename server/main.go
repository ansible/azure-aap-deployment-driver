package main

import (
	"path/filepath"
	"server/api"
	"server/azure"
	"server/config"
	"server/controllers"
	"server/engine"
	"server/persistence"
)

func main() {
	config.ConfigureLogging()
	config.ParseArgs()

	db := persistence.NewPersistentDB(filepath.Join(config.GetEnvironment().BASE_PATH, config.GetEnvironment().DB_REL_PATH))
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
