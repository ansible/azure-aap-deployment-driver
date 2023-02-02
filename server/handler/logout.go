package handler

import "net/http"

func Logout(w http.ResponseWriter, r *http.Request) {
	sessionHelper, err := getSessionHelper()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := sessionHelper.RemoveSession(r, w); err != nil {
		respondError(w, http.StatusInternalServerError, "Could not remove session")
	}
	respondJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
