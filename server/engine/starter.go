package engine

import (
	"context"
	"encoding/json"

	"server/config"
	"server/model"
	"server/persistence"
	"server/templates"

	log "github.com/sirupsen/logrus"
)

func NewEngine(ctx context.Context, db *persistence.Database) *Engine {
	engine := &Engine{
		context:              ctx,
		database:             db,
		resolver:             NewResolver(config.GetEnvironment().SUBSCRIPTION, config.GetEnvironment().RESOURCE_GROUP_NAME),
		done:                 make(chan struct{}),
		maxExecutionRestarts: config.GetEnvironment().EXECUTION_MAX_RETRY,
	}
	engine.initialize()
	return engine
}

func (engine *Engine) initialize() {
	// Load status from DB to check whether the templates and main outputs need to be processed
	status := model.Status{}
	engine.database.Instance.Find(&status)

	if !status.TemplatesLoaded {
		// Load templates into database
		templatePath := config.GetEnvironment().TEMPLATE_PATH
		log.Infof("Starting deployment template discovery in location: %s", templatePath)
		templateOrderArray, err := templates.DiscoverTemplateOrder(templatePath)
		if err != nil {
			log.Fatalf("Unable to import ARM templates: %v", err)
		}
		stepCount := 0
		for i, templateBatch := range templateOrderArray {
			for _, templateName := range templateBatch {
				templateContent, err := templates.ReadJSONTemplate(config.GetEnvironment().TEMPLATE_PATH, templateName)
				if err != nil {
					log.Fatalf("Unable to read in template file for [%s]", templateName)
				}
				parametersContent, err := templates.ReadJSONTemplateParameters(config.GetEnvironment().TEMPLATE_PATH, templateName)
				if err != nil {
					log.Fatalf("Unable to read in template file for [%s]", templateName)
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
			status.TemplatesLoaded = true
			engine.database.Instance.Save(&status)
		}
		log.Infof("Finished deployment template discovery, stored %d steps in database.", stepCount)
	} else {
		log.Infof("Skipped discovery of templates, they are in database already.")
	}

	if !status.MainOutputsLoaded {
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
		status.MainOutputsLoaded = true
		engine.database.Instance.Save(&status)
		log.Infof("Finished parsing main outputs, stored %d outputs in database.", len(outputValues))
	} else {
		log.Infof("Skipped parsing and storing main outputs, they are in database already.")
	}

	// Allways read the main outputs from DB
	engine.mainOutputs = &model.Output{}
	engine.database.Instance.Find(engine.mainOutputs, model.Output{ModuleName: ""})
}
