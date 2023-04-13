package handler

import (
	"fmt"
	"net/http"

	"server/model"

	"gorm.io/gorm"
)

type InstallationStatus string

const (
	Deploying InstallationStatus = "DEPLOYING"
	Canceled  InstallationStatus = "CANCELED"
	Failed    InstallationStatus = "FAILED"
	Done      InstallationStatus = "DONE"
)

func (i InstallationStatus) toString() string {
	return fmt.Sprintf("%v", i)
}

func Status(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	steps := []model.Step{}
	db.Model(&model.Step{}).Preload("Executions").Find(&steps)
	status := Done
	for _, step := range steps {
		latestExecution := getLatestExecution(db, step)
		if latestExecution.Status == model.PermanentlyFailed || latestExecution.Status == model.RestartTimedOut {
			status = Failed
			break
		} else if latestExecution.Status == model.Canceled {
			status = Canceled
			break
		} else if latestExecution.Status != model.Succeeded {
			status = Deploying
		}
	}
	respondJSON(w, http.StatusOK, map[string]string{"status": status.toString()})
}

func getLatestExecution(db *gorm.DB, step model.Step) model.Execution {
	latestExecution := model.Execution{}
	// Avoid GORM error from Last() if no executions yet
	var count int64
	db.Model(&model.Execution{}).Where("step_id = ?", step.ID).Count(&count)
	if count > 0 {
		db.Last(&latestExecution, "step_id = ?", step.ID)
	}
	return latestExecution
}
