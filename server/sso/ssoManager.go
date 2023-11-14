package sso

import (
	"context"
	"fmt"
	"server/config"
	"server/handler"
	"server/model"
	"server/persistence"

	log "github.com/sirupsen/logrus"
)

var ssoManager *SsoManager

type SsoManager struct {
	Context       context.Context
	Authenticator *handler.Authenticator
	SsoHandler    *handler.SsoHandler
	Db            *persistence.Database
}

func NewSsoManager(ctx context.Context, db *persistence.Database) (*SsoManager, error) {
	ssoManager = &SsoManager{
		Context: ctx,
		Db:      db,
	}
	err := ssoManager.initialize()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SSO manager: %v", err)
	}
	return ssoManager, nil
}

func (s *SsoManager) initialize() error {
	store := model.GetSsoStore()
	var credentials *model.SsoCredentials
	var err error
	if store == nil {
		// Create dynamic client, get credentials
		credentials, err = createAndStoreSsoCredentials(s.Context, s.Db)
		if err != nil {
			return fmt.Errorf("unable to create client and get credentials: %v", err)
		}
	} else {
		// Fetch from db
		credentials, err = model.GetSsoStore().GetSsoClientCredentials()
		if err != nil {
			return fmt.Errorf("unable to load existing SSO credentials from db: %v", err)
		}
	}

	auth, err := handler.NewAuthenticator(ssoManager.Context, config.GetEnvironment().SSO_ENDPOINT, handler.GetRedirectUrl(), credentials.ClientId, credentials.ClientSecret)
	if err != nil {
		return fmt.Errorf("unable to instantiate SSO authenticator: %v", err)
	}
	ssoManager.SsoHandler = handler.GetSsoHandler(auth)
	return nil
}

func (s *SsoManager) DeleteAcsClient() {
	acsClient := GetAcsClient(s.Context)
	credentials, _ := model.GetSsoStore().GetSsoClientCredentials()
	acsClient.DeleteACSClient(credentials.ClientId)
}

func createAndStoreSsoCredentials(ctx context.Context, db *persistence.Database) (*model.SsoCredentials, error) {
	acsClient := GetAcsClient(ctx)
	credentials, err := acsClient.GetClientCredentials(handler.GetRedirectUrl())
	if err != nil {
		log.Errorf("Unable to create SSO client, fall back to credentials login: %v", err)
		return nil, err
	} else {
		model.InitSsoStore(db.Instance).SetSsoClientCredentials(credentials.ClientId, credentials.ClientSecret)
	}
	return credentials, nil
}
