package engine

import (
	"context"
	"encoding/json"

	"server/config"
	"server/model"
	"server/persistence"
	"server/templates"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
)

func NewEngine(ctx context.Context, db *persistence.Database, modmClient *sdk.Client) *Engine {
	engine := &Engine{
		context:              ctx,
		database:             db,
		done:                 make(chan struct{}),
		status:               &model.Status{},
		maxExecutionRestarts: config.GetEnvironment().EXECUTION_MAX_RETRY,
		modmClient:           modmClient,
		deploymentStarted:    false,
		deploymentComplete:   false,
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

		mainTemplate, mainParameters, err := templates.GetMainTemplateAndParameters(templatePath)
		if err != nil {
			engine.Fatalf("Unable to read in main template and parameters files")
		}

		// Store template and parameters
		engine.template = mainTemplate
		engine.parameters = mainParameters
		// start steps with dry run step at the beginning
		engine.addDryRunStep(mainTemplate)
		engine.createModmDeployment(mainTemplate)
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

func (engine *Engine) addDryRunStep(mainTemplate map[string]any) {
	engine.database.Instance.Create(&model.Step{
		Name: model.DryRunStepName,
	})
	log.Info("Added dry run step to database")
}

func (engine *Engine) createModmDeployment(mainTemplate map[string]any) {
	createDeployment := sdk.CreateDeployment{}
	createDeployment.ResourceGroup = to.Ptr(config.GetEnvironment().RESOURCE_GROUP_NAME)
	createDeployment.Location = to.Ptr(config.GetEnvironment().AZURE_LOCATION)
	createDeployment.Name = to.Ptr(model.ModmDeploymentName)
	createDeployment.Template = mainTemplate
	createDeployment.SubscriptionID = to.Ptr(config.GetEnvironment().SUBSCRIPTION)

	resp, err := engine.modmClient.Create(engine.context, createDeployment)
	if err != nil {
		log.Fatalf("Failed to create MODM deployment: %v", err)
	}
	for _, stage := range resp.Stages {
		step := model.Step{
			Name:    *stage.Name,
			StageId: *stage.ID,
		}
		engine.database.Instance.Save(&step)
	}
	engine.modmDeploymentId = int(*resp.ID)
	log.Infof("Created MODM deployment.  Added %d steps to database.", len(resp.Stages))
}
