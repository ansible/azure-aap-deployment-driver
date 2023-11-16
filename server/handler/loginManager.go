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
	resp["authtype"] = config.GetEnvironment().AUTH_TYPE
	respondJSON(w, http.StatusOK, resp)
}
