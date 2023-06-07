package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/config"
	"server/engine"
	"strings"

	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Receive a POST request from the MODM webhook
func EventHook(db *gorm.DB, engine *engine.Engine, w http.ResponseWriter, r *http.Request) {
	apiKey := config.GetEnvironment().WEB_HOOK_API_KEY

	handler := newEventHookHandler(r, w, apiKey, db)
	message, err := handler.getMessage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Debugf("RAW event: %s - %s - %v", message.Type, message.Subject, message.Data)

	switch message.Type {
	case sdk.EventTypeDeploymentCompleted.String():
		// Overall deployment complete
		log.Debug("Received event: Deployment completed.")
		//deploymentController.Done(message)
	case "dryRunStarted":
		log.Debugf("Received event: Dry Run started: %s", message.Subject)
		if message.Status != string(sdk.StatusSuccess) {
			engine.CreateExecution(message)
		} else {
			// TODO not sure what this actually means and what to do in this case
			log.Errorf("Dry Run start was not successful but was: %s", message.Status)
		}
	case sdk.EventTypeStageCompleted.String():
		log.Debugf("Received event: Stage completed: %s", message.Subject)
		engine.UpdateExecution(message)
	case sdk.EventTypeDryRunCompleted.String():
		log.Debugf("Received event: Dry Run completed: %s", message.Subject)
		engine.UpdateExecution(message)
	}
	respondJSON(w, http.StatusOK, message)
}

//region handler

type eventHookHandler struct {
	request  *http.Request
	response http.ResponseWriter
	db       *gorm.DB
	// the API Key that was set when the event hook was created with the MODM client sdk
	apiKey string
}

// function that validates the API Key from MODM to protect against unauthorized requests
func (h *eventHookHandler) authorizeRequest() error {
	var err error
	apiKeyToAuthorize, err := h.getApiKeyFromAuthorizationHeader()

	if !strings.EqualFold(apiKeyToAuthorize, h.apiKey) {
		err = errors.New("invalid API Key")
	}

	if err != nil {
		respondJSON(h.response, http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return err
}

func newEventHookHandler(request *http.Request, response http.ResponseWriter, apiKey string, db *gorm.DB) *eventHookHandler {
	return &eventHookHandler{
		request:  request,
		response: response,
		apiKey:   apiKey,
		db:       db,
	}
}

func (h *eventHookHandler) getMessage() (*sdk.EventHookMessage, error) {
	err := h.authorizeRequest()
	if err != nil {
		return nil, err
	}

	message := &sdk.EventHookMessage{}
	err = json.NewDecoder(h.request.Body).Decode(message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (h *eventHookHandler) getApiKeyFromAuthorizationHeader() (string, error) {
	authorizationHeaderValue := h.request.Header.Get("Authorization")

	if authorizationHeaderValue == "" {
		return "", errors.New("no Authorization header value found")
	}
	values := strings.Split(authorizationHeaderValue, " ")
	if len(values) != 2 || values[0] != "ApiKey" {
		return "", errors.New("invalid Authorization header value")
	}
	encodedApiKey := values[1]

	if encodedApiKey == "" {
		return "", errors.New("no API Key found in Authorization header")
	}

	apiKey, err := base64.StdEncoding.DecodeString(encodedApiKey)
	if err != nil {
		return "", fmt.Errorf("unable to decode API Key: %v", err)
	}

	return string(apiKey), nil
}

//endregion handler
