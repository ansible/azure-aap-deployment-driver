package handler

import (
	"encoding/json"
	"net/http"
	"server/config"
	"server/engine"

	"gorm.io/gorm"
)

type CredentialsHandler struct{}

type loginData struct {
	Username string `json:"uid"`
	Password string `json:"pwd"`
}

func (a CredentialsHandler) GetLoginHandler() HandleFuncWithDB {
	// initialization of the expected credentials
	const userName = "admin"
	userPassword := config.GetEnvironment().PASSWORD

	return func(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
		sessionHelper, err := getSessionHelper()
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		var data loginData

		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		if data.Username == userName && data.Password == userPassword {
			if err := sessionHelper.SetupSession(r, w, "credentials"); err != nil {
				respondError(w, http.StatusInternalServerError, "Could not setup session")
			}
			respondJSON(w, http.StatusOK, map[string]string{"status": "success"})
			// update metrics because user attempted logging in with correct credentials
			engine.UpdateLogInMetrics(db)
			return
		}

		respondError(w, http.StatusUnauthorized, "Invalid user id or password")
	}
}
