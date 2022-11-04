package handler

import (
	"net/http"

	"server/model"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetAllExecutions(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	executions := []model.Execution{}
	sess := db.Model(&model.Execution{})
	filterByStep := r.URL.Query().Get("stepId")
	if filterByStep != "" {
		sess.Where("step_id = ?", filterByStep)
	}
	sess.Find(&executions)
	respondJSON(w, http.StatusOK, executions)
}

func GetExecution(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	execution := getExecutionOr404(db, id, w, r)
	if execution == nil {
		return
	}
	respondJSON(w, http.StatusOK, execution)
}

func getExecutionOr404(db *gorm.DB, id string, w http.ResponseWriter, r *http.Request) *model.Execution {
	execution := model.Execution{}
	if err := db.First(&execution, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &execution
}

func Restart(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	execution := getExecutionOr404(db, id, w, r)
	if execution == nil {
		return
	}
	execution.Status = model.Restart
	db.Save(execution)
	RespondOk(w)
}
