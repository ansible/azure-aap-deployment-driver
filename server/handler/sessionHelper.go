package handler

import (
	"crypto/rand"
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/sessions"
)

var (
	once           sync.Once
	helperInstance *SessionHelper
	configuration  SessionHelperConfiguration
)

type SessionHelperConfiguration struct {
	AuthKey      []byte
	CookieName   string
	CookiePath   string
	CookieDomain string
	Secure       bool
	MaxAge       int
}

type SessionHelper struct {
	store       sessions.Store
	sessionName string
}

func ConfigureSessionHelper(c SessionHelperConfiguration) {
	configuration = c
}

func getSessionHelper() (*SessionHelper, error) {
	// checking the auth key and cookie name, that should be enough to detect lack of configuration
	if len(configuration.AuthKey) == 0 && len(configuration.CookieName) == 0 {
		return nil, errors.New("no session configuration provided")
	}
	// only initialize once
	once.Do(func() {
		aStore := sessions.NewCookieStore(configuration.AuthKey)
		aStore.Options.Domain = configuration.CookieDomain
		aStore.Options.Path = configuration.CookiePath
		aStore.Options.Secure = configuration.Secure
		aStore.Options.MaxAge = configuration.MaxAge
		aStore.Options.HttpOnly = true
		aStore.MaxAge(configuration.MaxAge)
		helperInstance = &SessionHelper{
			store:       aStore,
			sessionName: configuration.CookieName,
		}
	})
	return helperInstance, nil
}

func (s *SessionHelper) HasSession(r *http.Request) (bool, error) {
	aSession, err := s.store.Get(r, s.sessionName)
	if err != nil {
		return false, err
	}
	// only established session is considered
	return !aSession.IsNew, nil
}

func (s *SessionHelper) SetupSession(r *http.Request, w http.ResponseWriter) error {
	aSession, err := s.store.New(r, s.sessionName)
	if err != nil {
		return err
	}
	aSession.Options.HttpOnly = true
	aSession.Save(r, w)
	return nil
}

func (s *SessionHelper) RemoveSession(r *http.Request, w http.ResponseWriter) error {
	aSession, err := s.store.Get(r, s.sessionName)
	if err != nil {
		return err
	}
	aSession.Options.MaxAge = -1
	if err := aSession.Save(r, w); err != nil {
		return err
	}
	return nil
}

func GenerateSessionAuthKey() ([]byte, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
