package engine

import (
	"server/model"
	"server/persistence"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operation"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type testDryRunController struct {
	db   *gorm.DB
	done chan struct{}
}

func newTestDryRunController() *testDryRunController {
	database := persistence.NewInMemoryDB()
	engine := &Engine{
		database: database,
	}
	engine.addDryRunStep(nil, nil, 0)

	return &testDryRunController{
		db:   database.Instance,
		done: make(chan struct{}),
	}
}

func Test_dryRunController_getStep(t *testing.T) {
	test := newTestDryRunController()

	controller := &dryRunController{
		db: test.db,
	}

	step, err := controller.getStep()

	assert.NoError(t, err)
	assert.Equal(t, model.DryRunStepName, step.Name)
}

func Test_dryRunController_getStep_fetches_association_with_operationId(t *testing.T) {
	test := newTestDryRunController()

	controller := &dryRunController{
		db: test.db,
	}

	step, _ := controller.getStep()

	// add execution step to ensure that's returned with the step
	operationId := uuid.New()
	step.Executions = append(step.Executions, model.Execution{
		DeploymentID:  "1",
		Status:        model.Started,
		CorrelationID: operationId.String(),
	})

	test.db.Save(&step)

	step, err := controller.getStep()

	assert.NoError(t, err)
	assert.Equal(t, step.Executions[0].CorrelationID, operationId.String())
	assert.Equal(t, step.Executions[0].ID, uint(1))
}

func Test_dryRunController_update_sets_execution(t *testing.T) {
	test := newTestDryRunController()

	controller := &dryRunController{
		db: test.db,
	}

	operationId := uuid.New()
	message := &events.EventHookMessage{
		Id:     uuid.New(),
		Status: operation.StatusSuccess.String(),
		Data: events.DeploymentEventData{
			DeploymentId: 1,
			OperationId:  operationId,
		},
	}

	err := controller.update(message)
	assert.NoError(t, err)

	step, err := controller.getStep()
	assert.NoError(t, err)

	assert.Equal(t, operationId.String(), step.Executions[0].CorrelationID)
}

func Test_dryRunController_dryRunDone(t *testing.T) {
	test := newTestDryRunController()

	controller := &dryRunController{
		db:   test.db,
		done: test.done,
	}

	go func() {
		operationId := uuid.New()
		message := &events.EventHookMessage{
			Id:     uuid.New(),
			Status: operation.StatusSuccess.String(),
			Data: events.DeploymentEventData{
				DeploymentId: 1,
				OperationId:  operationId,
			},
		}
		time.Sleep(2 * time.Second)
		controller.dryRunDone(message)
	}()
	log.Print("waiting for done")
	<-test.done
	log.Print("dryRunDone called")
	assert.True(t, true)
}
