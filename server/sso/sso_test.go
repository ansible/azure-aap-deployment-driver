package sso_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"server/sso"
	"server/test"
	"strings"
	"testing"
)

func getAcsClient(t *testing.T) *sso.AcsClient {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "token") {
			// Token request
			if r.URL.Path != sso.TOKEN_API {
				t.Errorf("Expected path %s, got path %s", sso.TOKEN_API, r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"access_token": "token"}`))
			if err != nil {
				t.Error("Testing error, could not write server response.")
			}
		} else if strings.HasSuffix(r.URL.Path, "v1") {
			// Client create request
			if r.URL.Path != sso.REG_API {
				t.Errorf("Expected path %s, got path %s", sso.REG_API, r.URL.Path)
			}
			w.WriteHeader(http.StatusCreated)
			_, err := w.Write([]byte(`{"clientId": "clientid", "secret": "secret"}`))
			if err != nil {
				t.Error("Testing error, could not write server response.")
			}
		} else if strings.HasSuffix(r.URL.Path, "12345") {
			// Client delete request
			w.WriteHeader(http.StatusNoContent)
		} else {
			// Unexpected
			t.Errorf("HTTP request path not expected: %s", r.URL.Path)
		}
	}))
	//defer server.Close()

	// Set up environment
	test.SetEnvironment()
	os.Setenv("SSO_ENDPOINT", server.URL)

	// Initialize client
	acsClient := sso.GetAcsClient(context.Background())
	return acsClient
}
func TestGetToken(t *testing.T) {
	client := getAcsClient(t)
	if client.Token != "token" {
		t.Errorf("Expected token == token, but got %s", client.Token)
	}

	if client.ClientId != "" {
		t.Errorf("Expected blank client ID, got %s", client.ClientId)
	}
}

func TestGetClientCredentials(t *testing.T) {
	client := getAcsClient(t)
	creds, err := client.GetClientCredentials("https://127.0.0.1/ssocallback")
	if err != nil {
		t.Errorf("Expected no error from get client credentials, got %v", err)
	}
	if creds.ClientId != "clientid" {
		t.Errorf("Expected clientid of clientid, got %s", creds.ClientId)
	}
}

func TestDeleteClientCredentials(t *testing.T) {
	client := getAcsClient(t)
	resp, err := client.DeleteACSClient("12345")
	if err != nil {
		t.Errorf("Expected no error from delete client credentials, got %v", err)
	}
	if resp != nil {
		t.Errorf("Expected nil delete client response, was not nil: %v", resp)
	}
}
