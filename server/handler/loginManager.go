package handler

import (
	"net/http"
	"server/config"
)

type LoginManager interface {
	GetLoginHandler() HandleFuncWithDB
}

func AuthType(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	if config.IsSsoEnabled() {
		resp["authtype"] = "SSO"
	} else {
		resp["authtype"] = "CREDENTIALS"
	}
	respondJSON(w, http.StatusOK, resp)
}
