package engine

import (
	"server/model"

	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

// func (d *dryRunController) save(model *model.DryRun) error {
// 	tx := d.db.Begin()
// 	tx.Save(&model)

// 	if tx.Error != nil {
// 		tx.Rollback()
// 		return tx.Error
// 	}
// 	tx.Commit()

// 	return nil
// }

func (d *dryRunController) getStep() (model.Step, error) {
	var step model.Step
	// TODO: get model.Step from DB
	return step, nil
}

func (c *dryRunController) dryRunDone(eventHook *events.EventHookMessage) {
	// TODO update execution
	c.done <- struct{}{}
}

func (c *dryRunController) save(response *sdk.DryRunResponse) error {
	// dryRun, err := c.getStep()

	// tx.Save(&model)

	// if tx.Error != nil {
	// 	tx.Rollback()
	// 	return tx.Error
	// }
	// tx.Commit()

	return nil
}

func (c *dryRunController) getDryRun() (*model.Step, error) {
	// db := c.db

	// // there could be more dry runs, so get the most recent one
	// dryRun := &model.Step{}
	// tx := db.Model(dryRun).Order("updated_at desc").First(dryRun)
	// if tx.Error != nil { // not found
	// 	return nil, tx.Error
	// }
	// return dryRun, nil
	return nil, nil
}

func (c *dryRunController) saveResult(message *events.EventHookMessage) error {
	// TODO: save the dry run record as an execution on the step
	// get the step by model.DryRunStepName
	//	engine.database.Instance.Model(&model.Step{}).Preload("Executions").Find(&steps)
	return nil
}
