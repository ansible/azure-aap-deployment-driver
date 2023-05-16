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
		// Load templates into database
		templatePath := config.GetEnvironment().TEMPLATE_PATH
		templateOrderArray, err := templates.DiscoverTemplateOrder(templatePath)
		if err != nil {
			engine.Fatalf("Unable to import ARM templates: %v", err)
		}

		mainTemplate, mainParameters, err := templates.GetMainTemplateAndParameters(templatePath)
		if err != nil {
			engine.Fatalf("Unable to read in main template and parameters files")
		}

		// start steps with dry run step at the beginning
		insertDryRunAt := 0
		engine.addDryRunStep(mainTemplate, mainParameters, insertDryRunAt)
		engine.addSteps(templateOrderArray, insertDryRunAt+1, templatePath)

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

// explicity adds a dry run step
//
//	remarks: expects the full main template and parameters
func (engine *Engine) addDryRunStep(mainTemplate map[string]any, mainParameters map[string]any, priority int) {
	engine.database.Instance.Create(&model.Step{
		Priority:   uint(priority),
		Name:       model.DryRunStepName,
		Template:   mainTemplate,
		Parameters: mainParameters,
		Executions: []model.Execution{},
	})
}

func (engine *Engine) addSteps(templateOrderArray [][]string, startAt int, templatePath string) {
	stepCount := startAt
	for i, templateBatch := range templateOrderArray {
		priority := startAt + i
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
				Priority:   uint(priority),
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
}
