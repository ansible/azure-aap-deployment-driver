package handler

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

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

func AuthStatusHandler(w http.ResponseWriter, r *http.Request) {
	sessionHelper, err := getSessionHelper()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// check both types of session, credentials and SSO
	var hasSSOSession, hasCredSession bool
	if hasSSOSession, err = sessionHelper.ValidSession(r); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}
	if hasCredSession, err = sessionHelper.HasSession(r); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}
	// set response headers indicating what sessions are established
	w.Header().Add("X-Session-Creds", strconv.FormatBool(hasCredSession))
	w.Header().Add("X-Session-SSO", strconv.FormatBool(hasSSOSession))

	log.Tracef("Session status: creds: %t, sso:%t", hasCredSession, hasSSOSession)

	if hasCredSession && hasSSOSession {
		respondOk(w)
		return
	}
	respondError(w, http.StatusUnauthorized, "Not authenticated.")
}
