package handler

import (
	"net/http"
	"server/model"

	"gorm.io/gorm"
)

// GET gets the most recent dry run instance
func DryRun(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	handler := newDryRunHandler(r, w, db)

	dryRun, err := handler.getCurrentDryRun()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	respondJSON(w, http.StatusOK, dryRun)
}

//region handler

type dryRunHandler struct {
	request  *http.Request
	response http.ResponseWriter
	db       *gorm.DB
}

func newDryRunHandler(r *http.Request, w http.ResponseWriter, db *gorm.DB) *dryRunHandler {
	return &dryRunHandler{
		request:  r,
		response: w,
		db:       db,
	}
}

// gets the dry run instance from persistent storage
func (h *dryRunHandler) getCurrentDryRun() (*model.DryRun, error) {
	db := h.db

	// there could be more dry runs, so get the most recent one
	dryRun := &model.DryRun{}
	tx := db.Model(dryRun).Order("updated_at desc").First(dryRun)
	if tx.Error != nil { // not found
		return nil, tx.Error
	}
	return dryRun, nil
}

//endregion handler
