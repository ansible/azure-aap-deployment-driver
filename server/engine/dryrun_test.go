package engine

import (
	"server/model"
	"server/persistence"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type dryRunTest struct {
	db   *gorm.DB
	done chan struct{}
}

func newDryRunTest() *dryRunTest {
	database := persistence.NewNoCacheInMemoryDb()
	db, err := database.Instance.DB()
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	database.Instance = database.Instance.Session(&gorm.Session{NewDB: true})

	engine := &Engine{
		database: database,
	}
	engine.addDryRunStep(nil)

	return &dryRunTest{
		db:   database.Instance,
		done: make(chan struct{}),
	}
}

func Test_dryRunController_getStep(t *testing.T) {
	test := newDryRunTest()

	controller := &dryRunController{
		db: test.db,
	}

	step, err := controller.getStep()

	assert.NoError(t, err)

	assert.Equal(t, 0, len(step.Executions))
	assert.Equal(t, model.DryRunStepName, step.Name)
}
/*  TODO Fix tests
func Test_dryRunController_getStep_fetches_association_with_dryrun_id(t *testing.T) {
	test := newDryRunTest()

	controller := &dryRunController{
		db: test.db,
	}

	step, err := controller.getStep()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(step.Executions))

	// add execution step to ensure that's returned with the step
	id := uuid.MustParse("d42e3b3d-c54e-4408-aa65-b9c98d137a72")
	execution := model.Execution{
		StepID: step.ID,
		Status: model.Started,
		DryRunExecution: &model.DryRunExecution{
			Id:           id.String(),
			DeploymentId: 1,
		},
	}

	test.db.Save(&execution)

	step, err = controller.getStep()
	assert.Equal(t, 1, len(step.Executions))

	assert.NoError(t, err)
	assert.Equal(t, step.Executions[0].DryRunExecution.Id, id.String())
	assert.Equal(t, step.Executions[0].ID, uint(1))
}

func Test_dryRunController_update_succeeds_with_execution_by_dryrun_id(t *testing.T) {
	test := newDryRunTest()

	dryRunInstanceId := uuid.MustParse("f181f551-5d17-4ab4-bbcb-407d47b63f77")

	controller := &dryRunController{
		db: test.db,
		execution: &model.Execution{},
	}

	//setup execution with the dryrun id
	controller.createExecution(1, &sdk.InvokeDryRunResponse{
		Id:     dryRunInstanceId,
		Status: sdk.StatusScheduled.String(),
	}, nil)

	//simulate the event hook message
	message := &sdk.EventHookMessage{
		Id:     uuid.New(),
		Type:   string(sdk.EventTypeDryRunCompleted),
		Status: sdk.StatusSuccess.String(),
		Data: sdk.DryRunEventData{
			DeploymentId: 1,
			OperationId:  dryRunInstanceId,
		},
	}

	//should succeed to update
	err := controller.updateExecution(message)
	assert.NoError(t, err)

	step, _ := controller.getStep()
	assert.Equal(t, 1, len(step.Executions))
	assert.Equal(t, dryRunInstanceId.String(), step.Executions[0].DryRunExecution.Id)
}

func Test_dryRunController_update_should_be_idempotent(t *testing.T) {
	test := newDryRunTest()
	controller := &dryRunController{
		db: test.db,
		execution: &model.Execution{},
	}

	dryRunInstanceId := uuid.MustParse("f181f551-5d17-4ab4-bbcb-407d47b63f77")
	message := &sdk.EventHookMessage{
		Id:     uuid.New(),
		Type:   string(sdk.EventTypeDryRunCompleted),
		Status: sdk.StatusSuccess.String(),
		Data: sdk.DryRunEventData{
			DeploymentId: 1,
			OperationId:  dryRunInstanceId,
		},
	}
	err := controller.updateExecution(message)
	assert.NoError(t, err)

	step, _ := controller.getStep()

	assert.Equal(t, 1, len(step.Executions))
	assert.Equal(t, dryRunInstanceId.String(), step.Executions[0].DryRunExecution.Id)
}

func Test_dryRunController_Done(t *testing.T) {
	test := newDryRunTest()

	controller := &dryRunController{
		db:   test.db,
		done: test.done,
		execution: &model.Execution{},
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
		time.Sleep(1 * time.Second)
		controller.Done(message)
	}()

	log.Print("waiting for done")
	<-test.done
	log.Print("dryRunDone called")
	assert.True(t, true)
}
 */