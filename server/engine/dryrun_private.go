package engine

import (
	"server/model"
	"strconv"

	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

func (c *dryRunController) getStep() (*model.Step, error) {
	step := &model.Step{}

	join := "left join executions on executions.step_id = steps.id"
	tx := c.db.Model(step).Preload("Executions").Joins(join).Where("steps.name = ?", model.DryRunStepName).First(step)
	if tx.Error != nil { // not found
		return nil, tx.Error
	}
	return step, nil
}

// updates the step execution (or inserts) and signals dry run is done
func (c *dryRunController) dryRunDone(message *events.EventHookMessage) {
	c.update(message)
	c.done <- struct{}{}
}

// creates a new step execution to track the dry run
func (c *dryRunController) create(response *sdk.DryRunResponse) error {
	tx := c.db.Begin()
	step, err := c.getStep()
	if err != nil {
		return err
	}

	status := model.Started
	if response.Status != operation.StatusScheduled.String() {
		status = model.Failed
	}

	execution := model.Execution{
		StepID:        step.ID,
		DeploymentID:  strconv.Itoa(c.deploymentId),
		Status:        status,
		CorrelationID: response.Id.String(),
	}

	tx.Save(&execution)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()

	return nil
}

func (c *dryRunController) update(message *events.EventHookMessage) error {
	step, err := c.getStep()
	if err != nil {
		return err
	}
	data := message.Data.(events.DeploymentEventData)
	var execution *model.Execution

	for i := range step.Executions {
		if step.Executions[i].CorrelationID == data.OperationId.String() {
			execution = &step.Executions[i]
			break
		}
	}

	if execution == nil {
		execution = &model.Execution{StepID: step.ID, CorrelationID: data.OperationId.String()}
		step.Executions = append(step.Executions, *execution)
	}

	status := model.Succeeded
	if message.Status == operation.StatusFailed.String() {
		status = model.Failed
	}
	execution.Status = status
	execution.Details = data.Message

	c.db.Save(&step.Executions)
	return nil
}
