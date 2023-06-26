package handler

import (
	"fmt"
	"net/http"
	"server/azure"
	"server/config"
	"server/controllers"
	"strings"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func DeleteContainer(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	qParams := r.URL.Query()
	if strings.ToLower(qParams.Get("confirmation")) != "yes" {
		respondError(w, http.StatusBadRequest, "Container deletion was not confirmed.")
		return
	}
	err := azure.DeleteStorageAccount(config.GetEnvironment().RESOURCE_GROUP_NAME, config.GetEnvironment().STORAGE_ACCOUNT_NAME)
	if err != nil {
		msg := fmt.Sprintf("Failed to delete storage account: %v", err)
		log.Printf("Failed to delete storage account: %v", err)
		respondError(w, http.StatusInternalServerError, msg)
	}
	err = azure.DeleteContainer(config.GetEnvironment().RESOURCE_GROUP_NAME, config.GetEnvironment().CONTAINER_GROUP_NAME)
	if err != nil {
		msg := fmt.Sprintf("Failed to delete container: %v", err)
		log.Printf(msg)
		respondError(w, http.StatusInternalServerError, msg)
	}
	respondOk(w)
}

func Terminate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	qParams := r.URL.Query()
	if strings.ToLower(qParams.Get("confirmation")) != "yes" {
		respondError(w, http.StatusBadRequest, "Installer termination was not confirmed.")
		return
	}
	log.Printf("Terminating execution due to API request.  Bye!")
	respondOk(w)
	err := controllers.NewExitController().Stop()
	if err != nil {
		log.Errorf("Error while calling Stop on exit controller: %v", err)
	}
}
