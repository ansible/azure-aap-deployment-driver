package handler

import (
	"encoding/json"
	"net/http"
	"server/engine"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HandleFuncWithDB func(db *gorm.DB, w http.ResponseWriter, r *http.Request)
type HandleFuncWithDBAndEngine func(db *gorm.DB, engine *engine.Engine, w http.ResponseWriter, r *http.Request)

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("Error: Unable to write error to http output stream: %v", err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write([]byte(response))
	if err != nil {
		log.Printf("Error: Unable to write response to http output stream: %v", err)
	}
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

func respondOk(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}
