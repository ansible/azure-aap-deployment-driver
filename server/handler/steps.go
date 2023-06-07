package handler

import (
	"net/http"

	"server/engine"
	"server/model"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetAllSteps(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	steps := getAllSteps(db)
	respondJSON(w, http.StatusOK, steps)
}

func GetStep(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	step := getStepOr404(db, id, w, r)
	if step == nil {
		return
	}
	respondJSON(w, http.StatusOK, step)
}

func CancelAllSteps(db *gorm.DB, engine *engine.Engine, w http.ResponseWriter, r *http.Request) {
	steps := getAllSteps(db)
	// first mark all non executing steps as cancelled
	for _, aStep := range steps {
		if len(aStep.Executions) == 0 {
			db.Save(&model.Execution{
				Status: model.Canceled,
				StepID: aStep.ID,
			})
		}
	}
	engine.CancelDeployment()
}

func getAllSteps(db *gorm.DB) []model.Step {
	steps := []model.Step{}
	db.Model(&model.Step{}).Preload("Executions").Find(&steps)
	return steps
}

func getStepOr404(db *gorm.DB, id string, w http.ResponseWriter, r *http.Request) *model.Step {
	step := model.Step{}
	if err := db.Preload("Executions").First(&step, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &step
}
