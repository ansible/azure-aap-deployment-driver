package engine

import (
	"context"
	"encoding/json"

	"server/config"
	"server/model"
	"server/persistence"
	"server/templates"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	log "github.com/sirupsen/logrus"
)

func NewEngine(ctx context.Context, db *persistence.Database, client *armresources.DeploymentsClient) *Engine {
	engine := &Engine{
		context:              ctx,
		database:             db,
		resolver:             NewResolver(config.GetEnvironment().SUBSCRIPTION, config.GetEnvironment().RESOURCE_GROUP_NAME),
		done:                 make(chan struct{}),
		status:               &model.Status{},
		maxExecutionRestarts: config.GetEnvironment().EXECUTION_MAX_RETRY,
		deploymentsClient:    client,
	}
	engine.initialize()
	return engine
}

func (engine *Engine) initialize() {
	// Load status from DB to check whether the templates and main outputs need to be processed
	engine.database.Instance.Find(engine.status)

	if !engine.status.TemplatesLoaded {
		templatePath := config.GetEnvironment().TEMPLATE_PATH
		// Load main template
		mainTemplate, mainParameters, err := templates.GetMainTemplateAndParameters(templatePath)
		if err != nil {
			engine.Fatalf("Unable to read in main template and parameters files")
		}

		engine.createWhatIfStep(mainTemplate, mainParameters)

		// Load templates into database
		templateOrderArray, err := templates.DiscoverTemplateOrder(templatePath)

		if err != nil {
			engine.Fatalf("Unable to import ARM templates: %v", err)
		}

		stepCount := 1
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
					Priority:   uint(i + 1),
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
		log.Infof("Skipped discovery of templates, they are in database already.")
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
		log.Infof("Skipped parsing and storing main outputs, they are in database already.")
	}

	// Allways read the main outputs from DB
	engine.mainOutputs = &model.Output{}
	engine.database.Instance.Find(engine.mainOutputs, model.Output{ModuleName: ""})
}

func (engine *Engine) createWhatIfStep(mainTemplate map[string]interface{}, mainParams map[string]interface{}) {
	engine.database.Instance.Create(&model.Step{
		Name:       model.WHAT_IF_STEP_NAME,
		Template:   mainTemplate,
		Parameters: mainParams,
		Priority:   0,
	})
}
