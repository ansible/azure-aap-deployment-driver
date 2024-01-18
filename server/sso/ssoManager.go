package sso

import (
	"context"
	"errors"
	"fmt"
	"server/config"
	"server/controllers"
	"server/handler"
	"server/model"
	"server/persistence"

	"github.com/coreos/go-oidc/v3/oidc"
	log "github.com/sirupsen/logrus"
)

var ssoManager *SsoManager

const ID_SCOPE string = "api.ansible_on_cloud"

type SsoManager struct {
	Context       context.Context
	Authenticator *handler.Authenticator
	SsoHandler    *handler.SsoHandler
	Db            *persistence.Database
}

func NewSsoManager(db *persistence.Database, loginManager *handler.LoginManager) *SsoManager {
	if config.GetEnvironment().AUTH_TYPE != "SSO" {
		log.Infof("SSO not enabled, skipping setup.")
		return nil
	}
	model.InitSsoStore(db.Instance)
	ssoManager = &SsoManager{
		Context: context.Background(),
		Db:      db,
	}
	err := ssoManager.initialize()
	if err != nil {
		log.Errorf("Failed to initialize SSO manager, will use credentials login: %v", err)
		return nil
	}
	// TODO Add cancel handling for normal engine exit as well
	err = controllers.AddCancelHandler("SSO Client Cleanup", ssoManager.DeleteAcsClient)
	if err != nil {
		log.Errorf("Unable to add exit controller cancel handler for ACS client cleanup: %v", err)
	}
	log.Trace("SSO enabled.")
	config.EnableSso()
	*loginManager = ssoManager.SsoHandler
	return ssoManager
}

func (s *SsoManager) initialize() error {
	store := model.GetSsoStore()
	var credentials *model.SsoCredentials
	var err error
	if store.SsoCredentialsExist() {
		// Fetch from db
		credentials, err = model.GetSsoStore().GetSsoClientCredentials()
		if err != nil {
			return fmt.Errorf("unable to load existing SSO credentials from db: %v", err)
		}
		log.Trace("Existing SSO credentials loaded from database.")
	} else {
		// Create dynamic client, get credentials
		credentials, err = createAndStoreSsoCredentials(s.Context, s.Db)
		if err != nil {
			return fmt.Errorf("unable to create client and get credentials: %v", err)
		}
		log.Trace("Created new SSO client and credentials.")
	}

	authCfg := handler.AuthenticatorConfig{
		Context:      ssoManager.Context,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", ID_SCOPE},
		ClientId:     credentials.ClientId,
		ClientSecret: credentials.ClientSecret,
		RedirectUrl:  handler.GetRedirectUrl(),
		SsoEndpoint:  config.GetEnvironment().SSO_ENDPOINT,
		Audience:     ID_SCOPE,
	}

	auth, err := handler.NewAuthenticator(authCfg)
	if err != nil {
		return fmt.Errorf("unable to instantiate SSO authenticator: %v", err)
	}
	ssoManager.SsoHandler = handler.GetSsoHandler(auth)
	return nil
}

func (s *SsoManager) DeleteAcsClient() {
	acsClient := GetAcsClient(s.Context)
	credentials, _ := model.GetSsoStore().GetSsoClientCredentials()
	log.Trace("Deleting SSO client.")
	_, err := acsClient.DeleteACSClient(credentials.ClientId)
	if err != nil {
		log.Errorf("failed to delete ACS client: %v", err)
	}
}

func createAndStoreSsoCredentials(ctx context.Context, db *persistence.Database) (*model.SsoCredentials, error) {
	acsClient := GetAcsClient(ctx)
	if acsClient == nil {
		return nil, errors.New("unable to create ACS client")
	}
	credentials, err := acsClient.GetClientCredentials(handler.GetRedirectUrl())
	if err != nil {
		log.Errorf("Unable to create SSO client, fall back to credentials login: %v", err)
		return nil, err
	} else {
		err := model.InitSsoStore(db.Instance).SetSsoClientCredentials(credentials.ClientId, credentials.ClientSecret)
		if err != nil {
			log.Errorf("Unable to store SSO credentials in DB: %v", err)
			return nil, err
		}
	}
	return credentials, nil
}
