package sso_test

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"server/controllers"
	"server/handler"
	"server/model"
	"server/persistence"
	"server/sso"
	"server/test"
	"strings"
	"testing"

	"github.com/oauth2-proxy/mockoidc"
	"github.com/stretchr/testify/assert"
)

var client *sso.AcsClient
var server *httptest.Server
var ssoServer *mockoidc.MockOIDC

func getAcsClient(m *testing.M, clientId string, clientSecret string) {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "token") {
			// Token request
			if r.URL.Path != sso.TOKEN_API {
				log.Printf("Expected path %s, got path %s", sso.TOKEN_API, r.URL.Path)
				os.Exit(1)
			}
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"access_token": "token"}`))
			if err != nil {
				log.Printf("Testing error, could not write server response.")
				os.Exit(1)
			}
		} else if strings.HasSuffix(r.URL.Path, "v1") {
			// Client create request
			if r.URL.Path != sso.REG_API {
				log.Printf("Expected path %s, got path %s", sso.REG_API, r.URL.Path)
				os.Exit(1)
			}
			w.WriteHeader(http.StatusCreated)
			resp := sso.ClientResponse{ClientId: clientId, Secret: clientSecret}
			err := json.NewEncoder(w).Encode(&resp)
			if err != nil {
				log.Printf("Testing error, could not write server response.")
				os.Exit(1)
			}
		} else if strings.HasSuffix(r.URL.Path, clientId) {
			// Client delete request
			w.WriteHeader(http.StatusNoContent)
		} else if strings.HasSuffix(r.URL.Path, "45678") {
			// Client delete (failure) request
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(`{"error": "error"}`))
			if err != nil {
				log.Printf("Testing error, could not write server response.")
				os.Exit(1)
			}
		} else {
			// Unexpected
			log.Printf("HTTP request not expected: %v", r)
			os.Exit(1)
		}
	}))
}

func TestGetToken(t *testing.T) {
	assert.NotNil(t, client, "Expected not to get nil client.")
	assert.Equal(t, "token", client.Token)
	assert.Equal(t, "", client.ClientId, "Expected blank client ID.")
}

func TestGetClientCredentials(t *testing.T) {
	assert.NotNil(t, client, "Expected not to get nil client.")
	creds, err := client.GetClientCredentials("https://127.0.0.1/ssocallback")
	assert.Nil(t, err, "Expected no error from get client credentials.")
	assert.Equal(t, ssoServer.ClientID, creds.ClientId)
}

func TestGetClientCredentialsDb(t *testing.T) {
	db := persistence.NewInMemoryDB()
	ssoStore := model.InitSsoStore(db.Instance)
	err := ssoStore.SetSsoClientCredentials("dbclientid", "dbsecret")
	assert.Nil(t, err, "Unable to set SSO credentials.")

	creds, err := client.GetClientCredentials("https://127.0.0.1/ssocallback")
	assert.Nil(t, err, "Expected no error from get client credentials.")
	assert.Equal(t, "dbclientid", creds.ClientId)
}

func TestDeleteClientCredentials(t *testing.T) {
	assert.NotNil(t, client, "Expected client to not be nil.")
	resp, err := client.DeleteACSClient(ssoServer.ClientID)
	assert.Nil(t, err, "Expected no error from delete client credentials.")
	assert.Nil(t, resp, "Expected nil delete client response.")
}

func TestFailedDeleteClientCredentials(t *testing.T) {
	resp, err := client.DeleteACSClient("45678")
	assert.NotNil(t, err, "Expected error from delete client credentials.")
	assert.NotNil(t, resp, "Expected non-nil delete client response.")
}

func TestSsoManager(t *testing.T) {
	db := persistence.NewInMemoryDB()
	ssoStore := model.InitSsoStore(db.Instance)

	// Separate endpoint for dynamic client registration
	controllers.NewExitController()
	var lm handler.LoginManager
	man := sso.NewSsoManager(db, &lm)
	assert.NotNil(t, lm, "Expected SSO manager to instantiate login manager.")

	creds, _ := ssoStore.GetSsoClientCredentials()
	assert.Equal(t, ssoServer.ClientID, creds.ClientId)

	man.DeleteAcsClient()
	assert.False(t, ssoStore.SsoCredentialsExist(), "Expected db credentials to be deleted.")
}

func TestMain(m *testing.M) {
	ssoServer, _ = mockoidc.Run()

	getAcsClient(m, ssoServer.ClientID, ssoServer.ClientSecret)

	// Set up environment
	test.SetEnvironment()
	os.Setenv("AUTH_TYPE", "SSO")
	os.Setenv("SSO_ENDPOINT", ssoServer.Issuer())
	os.Setenv("DYNAMIC_CLIENT_REG_ENDPOINT", server.URL)

	// Initialize client
	client = sso.GetAcsClient(context.Background())
	os.Exit(m.Run())
}
