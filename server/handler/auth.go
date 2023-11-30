package handler

import "net/http"

var authDisabled bool

func ConfigureAuthenticationForTesting(disableFlag bool) {
	authDisabled = disableFlag
}

func EnsureAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sessionHelper, err := getSessionHelper()
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		//hasSession, err := sessionHelper.ValidSession(r)
		hasSession, err := sessionHelper.HasSession(r)
		if err != nil {
			// TODO add logging
			respondError(w, http.StatusInternalServerError, "Could not find or establish session.")
			return
		}
		if !authDisabled && !hasSession {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetAuthCheckHandler(checkSSO bool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionHelper, err := getSessionHelper()
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		var hasSession bool
		if checkSSO {
			hasSession, err = sessionHelper.ValidSession(r)
		} else {
			hasSession, err = sessionHelper.HasSession(r)
		}
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if hasSession {
			respondOk(w)
			return
		}
		respondError(w, http.StatusUnauthorized, "Not authenticated.")
	})
}
