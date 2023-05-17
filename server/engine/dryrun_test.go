package engine

import (
	"server/model"
	"server/persistence"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type testDryRunController struct {
	db   *gorm.DB
	done chan struct{}
}

func newTestDryRunController() *testDryRunController {
	database := persistence.NewNoCacheInMemoryDb()
	database.Instance = database.Instance.Session(&gorm.Session{NewDB: true})

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

	assert.Equal(t, 0, len(step.Executions))
	assert.Equal(t, model.DryRunStepName, step.Name)
}

func Test_dryRunController_getStep_fetches_association_with_operationId(t *testing.T) {
	test := newTestDryRunController()

	controller := &dryRunController{
		db: test.db,
	}

	step, err := controller.getStep()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(step.Executions))

	// add execution step to ensure that's returned with the step
	operationId := uuid.MustParse("d42e3b3d-c54e-4408-aa65-b9c98d137a72")
	execution := model.Execution{
		StepID:        step.ID,
		DeploymentID:  "1",
		Status:        model.Started,
		CorrelationID: operationId.String(),
	}

	test.db.Save(&execution)

	step, err = controller.getStep()
	assert.Equal(t, 1, len(step.Executions))

	assert.NoError(t, err)
	assert.Equal(t, step.Executions[0].CorrelationID, operationId.String())
	assert.Equal(t, step.Executions[0].ID, uint(1))
}

func Test_dryRunController_update(t *testing.T) {
	test := newTestDryRunController()
	controller := &dryRunController{
		db: test.db,
	}

	operationId := uuid.MustParse("f181f551-5d17-4ab4-bbcb-407d47b63f77")
	message := &sdk.EventHookMessage{
		Id:     uuid.New(),
		Status: sdk.StatusSuccess.String(),
		Data: sdk.DeploymentEventData{
			DeploymentId: 1,
			OperationId:  operationId,
		},
	}

	err := controller.update(message)
	assert.NoError(t, err)

	step, err := controller.getStep()
	assert.NoError(t, err)

	assert.Equal(t, 1, len(step.Executions))
	assert.Equal(t, operationId.String(), step.Executions[0].CorrelationID)

	// update to be ide
}

func Test_dryRunController_update_should_be_idempotent(t *testing.T) {
	test := newTestDryRunController()
	controller := &dryRunController{
		db: test.db,
	}

	operationId := uuid.MustParse("f181f551-5d17-4ab4-bbcb-407d47b63f77")
	message := &sdk.EventHookMessage{
		Id:     uuid.New(),
		Status: sdk.StatusSuccess.String(),
		Data: sdk.DeploymentEventData{
			DeploymentId: 1,
			OperationId:  operationId,
		},
	}
	err := controller.update(message)
	assert.NoError(t, err)

	err = controller.update(message)
	assert.NoError(t, err)

	step, _ := controller.getStep()

	assert.Equal(t, 1, len(step.Executions))
	assert.Equal(t, operationId.String(), step.Executions[0].CorrelationID)

	// update to be ide
}

func Test_dryRunController_dryRunDone(t *testing.T) {
	test := newTestDryRunController()

	controller := &dryRunController{
		db:   test.db,
		done: test.done,
	}

	go func() {
		operationId := uuid.New()
		message := &sdk.EventHookMessage{
			Id:     uuid.New(),
			Status: sdk.StatusSuccess.String(),
			Data: sdk.DeploymentEventData{
				DeploymentId: 1,
				OperationId:  operationId,
			},
		}
		time.Sleep(2 * time.Second)
		controller.DryRunDone(message)
	}()
	log.Print("waiting for done")
	<-test.done
	log.Print("dryRunDone called")
	assert.True(t, true)
}
