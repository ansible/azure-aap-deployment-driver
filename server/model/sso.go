package model

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SsoStore interface {
	SetSsoClientCredentials(string, string) error
	GetSsoClientCredentials() (*SsoCredentials, error)
	RemoveSsoClientCredentials()
	SsoCredentialsExist() bool
	CreateSession(*SsoSession) error
	RemoveSession(string) error
	ValidSession(string) bool
	GetSessions() ([]SsoSession, error)
}

var once sync.Once
var store SsoStore

type SsoCredentials struct {
	ClientId     string
	ClientSecret string
}

type SsoSession struct {
	Code           string
	State          string `gorm:"primaryKey"`
	OrganizationId string
	Name           string
	Email          string
}

type ssoStore struct {
	db *gorm.DB
}

func InitSsoStore(db *gorm.DB) SsoStore {
	once.Do(func() {
		store = ssoStore{
			db: db,
		}
	})
	return store
}

func GetSsoStore() SsoStore {
	return store
}

func (s ssoStore) SetSsoClientCredentials(clientId string, clientSecret string) error {
	creds := &SsoCredentials{}
	err := s.db.Where(SsoCredentials{ClientId: clientId}).Assign(SsoCredentials{ClientSecret: clientSecret}).FirstOrCreate(creds).Error
	if err != nil {
		return fmt.Errorf("unable to store/update SSO credentials in DB: %v", err)
	}
	return nil
}

func (s ssoStore) GetSsoClientCredentials() (*SsoCredentials, error) {
	creds := &SsoCredentials{}
	if err := s.db.First(creds).Error; err != nil {
		return nil, fmt.Errorf("unable to load SSO credentials from DB: %v", err)
	}
	return creds, nil
}

func (s ssoStore) SsoCredentialsExist() bool {
	var count int64
	s.db.Table("sso_credentials").Count(&count)
	if count > 0 {
		return true
	} else {
		return false
	}
}

func (s ssoStore) RemoveSsoClientCredentials() {
	// Clear table
	s.db.Exec("DELETE FROM sso_credentials")
}

func (s ssoStore) CreateSession(session *SsoSession) error {
	if err := s.db.Create(session).Error; err != nil {
		log.Errorf("Unable to store SSO session: %v", err)
		return err
	}
	return nil
}

func (s ssoStore) RemoveSession(sessionState string) error {
	if err := s.db.Delete(SsoSession{State: sessionState}).Error; err != nil {
		log.Errorf("Unable to remove SSO session: %v", err)
		return err
	}
	return nil
}

func (s ssoStore) GetSessions() ([]SsoSession, error) {
	sessions := []SsoSession{}
	tx := s.db.Find(sessions)
	if tx.Error != nil {
		log.Errorf("Unable to load SSO sessions from DB: %v", tx.Error)
		return nil, tx.Error
	}
	return sessions, nil
}

func (s ssoStore) ValidSession(sessionStateHash string) bool {
	sessions := []SsoSession{}
	if err := s.db.Find(&sessions).Error; err != nil {
		log.Errorf("Unable to load SSO sessions from DB, rejecting validation: %v", err)
		return false
	}

	for _, sess := range sessions {
		if sess.State == sessionStateHash {
			log.Trace("Validated SSO session hash.")
			return true
		}
	}
	log.Trace("Rejecting SSO session validation.  No matching sessions found.")
	return false
}
