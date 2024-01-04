package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"server/config"
	"server/model"
	"strings"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SsoHandler struct {
	Auth  *Authenticator
	State string
}

var handler *SsoHandler

func GetSsoHandler(auth *Authenticator) *SsoHandler {
	handler = &SsoHandler{Auth: auth}
	return handler
}

func (s *SsoHandler) GetLoginHandler() HandleFuncWithDB {
	return func(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
		// Redirect to SSO, setting new nonce/state value for each login
		err := s.rollState()
		if err != nil {
			respondError(w, http.StatusInternalServerError, fmt.Errorf("unable to generate single use SSO state value: %v", err).Error())
			return
		}
		http.Redirect(w, r, s.Auth.Config.AuthCodeURL(s.State), http.StatusTemporaryRedirect)
	}
}

func (s *SsoHandler) SsoRedirect(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if !strings.EqualFold(state, s.State) {
		log.Errorf("SSO state mismatch. Sent: %s, Received: %s", s.State, state)
		respondError(w, http.StatusUnauthorized, "SSO state values do not match.")
		return
	}
	sessionState := r.URL.Query().Get("session_state")

	log.Trace("Performing SSO exchange.")
	oauth2Token, err := s.Auth.Config.Exchange(context.Background(), code)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Load SSO client credentials from db
	ssoStore := model.GetSsoStore()
	ssoCredentials, err := ssoStore.GetSsoClientCredentials()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Trace("Verifying SSO token/extracting ID token.")
	_, err = s.Auth.VerifyToken(oauth2Token, ssoCredentials.ClientId)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sessionHelper, err := getSessionHelper()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Trace("Creating new SSO session.")
	err = model.GetSsoStore().CreateSession(sessionState, code)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Trace("Creating session cookie.")
	err = sessionHelper.SetupSession(r, w, sessionState)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Trace("Redirecting to deployment driver home.")
	http.Redirect(w, r, fmt.Sprintf("https://%s/", config.GetEnvironment().INSTALLER_DOMAIN_NAME), http.StatusTemporaryRedirect)
}

func (s *SsoHandler) rollState() error {
	log.Trace("Rolling SSO one time code.")
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	s.State = base64.StdEncoding.EncodeToString(b)
	return nil
}
