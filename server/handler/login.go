package handler

import (
	"encoding/json"
	"net/http"
)

type loginData struct {
	Username string `json:"uid"`
	Password string `json:"pwd"`
}

func GetLoginHandler(userName, userPassword string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			if err := sessionHelper.SetupSession(r, w); err != nil {
				respondError(w, http.StatusInternalServerError, "Could not setup session")
			}
			respondJSON(w, http.StatusOK, map[string]string{"status": "success"})
			return
		}

		respondError(w, http.StatusUnauthorized, "Invalid user id or password")
	})
}
