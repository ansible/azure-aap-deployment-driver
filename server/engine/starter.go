package engine

import (
	"context"
	"encoding/json"
	"path/filepath"

	"server/config"
	"server/controllers/entitlement"
	"server/model"
	"server/persistence"
	"server/telemetry"
	"server/templates"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	log "github.com/sirupsen/logrus"
)

func NewEngine(ctx context.Context, db *persistence.Database, client *armresources.DeploymentsClient, entitlement *entitlement.EntitlementAPIController) *Engine {
	engine := &Engine{
		context:                ctx,
		database:               db,
		resolver:               NewResolver(config.GetEnvironment().SUBSCRIPTION, config.GetEnvironment().RESOURCE_GROUP_NAME),
		done:                   make(chan struct{}),
		status:                 &model.Status{},
		maxExecutionRestarts:   config.GetEnvironment().EXECUTION_MAX_RETRY,
		deploymentsClient:      client,
		entitlementsController: entitlement,
	}
	engine.initialize()
	return engine
}

func (engine *Engine) initialize() {
	// Load status from DB to check whether the templates and main outputs need to be processed
	engine.database.Instance.Find(engine.status)

	if !engine.status.TemplatesLoaded {
		// Load templates into database
		templatePath := filepath.Join(config.GetEnvironment().BASE_PATH, config.GetEnvironment().TEMPLATE_REL_PATH)
		templateOrderArray, err := templates.DiscoverTemplateOrder(templatePath)

		if err != nil {
			engine.Fatalf("Unable to import ARM templates: %v", err)
		}

		stepCount := 0
		for i, templateBatch := range templateOrderArray {
			for _, templateName := range templateBatch {
				if engine.IsFatalState() {
					return
				}

				templateContent, err := templates.ReadJSONTemplate(templatePath, templateName)
				if err != nil {
					engine.Fatalf("Unable to read in template file for [%s]", templateName)
				}
				parametersContent, err := templates.ReadJSONTemplateParameters(templatePath, templateName)
				if err != nil {
					engine.Fatalf("Unable to read in template file for [%s]", templateName)
				}
				engine.database.Instance.Create(&model.Step{
					Priority:   uint(i),
					Name:       templateName,
					Template:   templateContent,
					Parameters: parametersContent,
				})
				stepCount++
			}
		}

		if stepCount > 0 {
			engine.status.TemplatesLoaded = true
			engine.database.Instance.Save(engine.status)
		}
		log.Infof("Finished deployment template discovery, stored %d steps in database.", stepCount)
	} else {
		log.Info("Skipped discovery of templates, they are in database already.")
	}

	if !engine.status.MainOutputsLoaded {
		// parse main outputs and store them in db
		outputValues := make(map[string]interface{})
		log.Info("Parsing main outputs from environment variable...")
		if err := json.Unmarshal([]byte(config.GetEnvironment().MAIN_OUTPUTS), &outputValues); err != nil {
			log.Fatalf("Couldn't parse main outputs: %v", err)
		}
		engine.database.Instance.Save(&model.Output{
			ModuleName: "", // outputs from main install part don't have module name
			Values:     outputValues,
		})
		engine.status.MainOutputsLoaded = true
		engine.database.Instance.Save(engine.status)
		log.Infof("Finished parsing main outputs, stored %d outputs in database.", len(outputValues))
	} else {
		log.Info("Skipped parsing and storing main outputs, they are in database already.")
	}

	log.Info("Initializing telemetry handler for Segment reporting.")
	engine.telemetryHandler = telemetry.Init(engine.database.Instance, engine.context)

	// Allways read the main outputs from DB
	engine.mainOutputs = &model.Output{}
	engine.database.Instance.Find(engine.mainOutputs, model.Output{ModuleName: ""})
}
