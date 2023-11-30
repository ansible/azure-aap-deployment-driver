package sso

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"server/config"
	"server/model"
	"server/util"
	"sync"

	log "github.com/sirupsen/logrus"
)

const REG_API string = "apis/beta/acs/v1"

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Type        string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type ClientRequest struct {
	Name         string   `json:"name"`
	RedirectUris []string `json:"redirectUris"`
	OrgId        string   `json:"orgId"`
}

type ClientError struct {
	Name string `json:"name"`
}

type ClientResponse struct {
	ClientId         string        `json:"clientId"`
	Secret           string        `json:"secret"`
	Name             string        `json:"name"`
	CreatedAt        int           `json:"createdAt"`
	Error            string        `json:"error"`
	ErrorDescription string        `json:"error_description"`
	Fields           []ClientError `json:"fields"`
}

const TOKEN_SCOPE string = "api.iam.clients.aoc"
const TOKEN_API string = "protocol/openid-connect/token"
const SSO_ORG_ID string = "123" // TODO Make this one value for test, one for production
const SSO_CLIENT_NAME string = "deploymentdriver"

var acsClient *AcsClient
var once sync.Once

type AcsClient struct {
	Endpoint     string
	ClientId     string
	ClientSecret string
	Token        string
	Context      context.Context
}

func GetAcsClient(ctx context.Context) *AcsClient {
	once.Do(func() {
		acsClient = &AcsClient{
			Endpoint:     config.GetEnvironment().SSO_ENDPOINT,
			ClientId:     config.GetEnvironment().SSO_CLIENT_ID,
			ClientSecret: config.GetEnvironment().SSO_CLIENT_SECRET,
			Context:      ctx,
		}
	})
	log.Trace("Getting SSO client API access token.")
	token, err := acsClient.getToken()
	if err != nil {
		log.Errorf("Unable to get SSO client access token: %v", err)
		return nil
	}
	acsClient.Token = token
	log.Trace("Initialized SSO client.")
	return acsClient
}

func (c *AcsClient) getToken() (string, error) {
	apiUrl, _ := url.JoinPath(c.Endpoint, TOKEN_API)

	form := make(map[string]string)
	form["grant_type"] = "client_credentials"
	form["scope"] = TOKEN_SCOPE
	form["client_id"] = c.ClientId
	form["client_secret"] = c.ClientSecret

	req := util.NewHttpRequester()
	resp, err := req.MakeRequestWithWWWFormUrlEncodedBody(c.Context, http.MethodPost, apiUrl, nil, form)
	if err != nil {
		return "", fmt.Errorf("unable to post token request to ACS: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed.  Status code: %d, text: %s", resp.StatusCode, string(resp.Body))
	}
	bodyJson := TokenResponse{}
	err = json.Unmarshal(resp.Body, &bodyJson)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshal token response: %v", err)
	}
	return bodyJson.AccessToken, nil
}

func (c *AcsClient) GetClientCredentials(redirectUrl string) (*model.SsoCredentials, error) {
	// Use stored credentials if they exist
	ssoStore := model.GetSsoStore()
	if ssoStore.SsoCredentialsExist() {
		log.Trace("Returning existing SSO client credentials.")
		creds, err := ssoStore.GetSsoClientCredentials()
		if err != nil {
			return nil, err
		}
		return creds, nil
	}
	log.Trace("Creating new SSO client.")
	resp, err := c.createACSClient(SSO_CLIENT_NAME, redirectUrl, SSO_ORG_ID)
	if err != nil {
		log.Errorf("unable to create SSO client: %v", err)
		return nil, err
	}
	return &model.SsoCredentials{ClientId: resp.ClientId, ClientSecret: resp.Secret}, nil
}

func (c *AcsClient) createACSClient(name string, redirectUrl string, orgId string) (*ClientResponse, error) {
	api, _ := url.JoinPath(c.Endpoint, REG_API)

	headers := make(map[string]string)
	headers["Authorization"] = getAuthHeader(c.Token)

	urls := []string{redirectUrl}
	reqBody := ClientRequest{Name: name, RedirectUris: urls, OrgId: orgId}

	req := util.NewHttpRequester()
	resp, err := req.MakeRequestWithJSONBody(c.Context, http.MethodPost, api, headers, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create ACS client request failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create ACS client request failed: %s", string(resp.Body))
	}
	bodyJson := ClientResponse{}
	err = json.Unmarshal(resp.Body, &bodyJson)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal ACS client creation response to json: %v", err)
	}
	return &bodyJson, nil
}

func (c *AcsClient) DeleteACSClient(clientID string) (*ClientResponse, error) {
	api, _ := url.JoinPath(c.Endpoint, REG_API)
	deleteUrl, _ := url.JoinPath(api, clientID)

	headers := make(map[string]string)
	headers["Authorization"] = getAuthHeader(c.Token)

	req := util.NewHttpRequester()
	resp, err := req.MakeRequestWithJSONBody(c.Context, http.MethodDelete, deleteUrl, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("delete ACS client request failed: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		bodyJson := ClientResponse{}
		err = json.Unmarshal(resp.Body, &bodyJson)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal failed DELETE response to json: %v", err)
		}
		return &bodyJson, fmt.Errorf("request to delete ACS client failed: %s", bodyJson.ErrorDescription)
	}
	model.GetSsoStore().RemoveSsoClientCredentials()
	return nil, nil
}

func getAuthHeader(token string) string {
	return fmt.Sprintf("Bearer %s", token)
}
