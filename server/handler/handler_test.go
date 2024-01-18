package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"server/api"
	"server/config"
	"server/handler"
	"server/model"
	"server/persistence"
	"server/test"
	"testing"

	"github.com/coreos/go-oidc/v3/oidc"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/oauth2-proxy/mockoidc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

var database *persistence.Database

func TestSteps(t *testing.T) {
	rec := testHttpRoute(t, "GET", "/step", nil)
	assert.Equal(t, 200, rec.Code)

	var response []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 4, len(response))
	for _, step := range response {
		switch step["name"] {
		case "step1", "step2", "step3", "step4":
			continue
		default:
			t.Errorf("Unexpected step in get all steps: %s", step["name"])
		}
	}
}

func TestStep(t *testing.T) {
	rec := testHttpRoute(t, "GET", "/step/1", nil)

	assert.Equal(t, 200, rec.Code)

	var response model.Step
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "step1", response.Name)
	assert.Equal(t, uint(0), response.Priority)
	assert.Equal(t, 2, len(response.Executions))
}

func TestMissingStep(t *testing.T) {
	rec := testHttpRoute(t, "GET", "/step/10", nil)
	assert.Equal(t, 404, rec.Code)
}

func TestExecutions(t *testing.T) {
	rec := testHttpRoute(t, "GET", "/execution", nil)
	assert.Equal(t, 200, rec.Code)

	var response []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(response))
}

func TestExecution(t *testing.T) {
	rec := testHttpRoute(t, "GET", "/execution/1", nil)

	assert.Equal(t, 200, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(armresources.ProvisioningStateFailed), response["provisioningState"])
	assert.Equal(t, string(model.Failed), response["status"])
}

func TestGetExecutionByStep(t *testing.T) {
	rec := testHttpRoute(t, "GET", "/execution?stepId=1", nil)

	assert.Equal(t, 200, rec.Code)

	var response []model.Execution
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, uint(1), response[0].StepID)
	assert.Equal(t, 2, len(response))
}

func TestMissingExecution(t *testing.T) {
	rec := testHttpRoute(t, "GET", "/execution/10", nil)
	assert.Equal(t, 404, rec.Code)
}

func TestRestartExecution(t *testing.T) {
	rec := testHttpRoute(t, "POST", "/execution/1/restart", nil)
	assert.Equal(t, 200, rec.Code)
	execution := model.Execution{}
	database.Instance.First(&execution, 1)
	assert.Equal(t, model.Restart, execution.Status)
}

func TestRestartMissingExecution(t *testing.T) {
	rec := testHttpRoute(t, "POST", "/execution/11/restart", nil)
	assert.Equal(t, 404, rec.Code)
}

func TestStatus(t *testing.T) {
	rec := testHttpRoute(t, "GET", "/status", nil)
	assert.Equal(t, 200, rec.Code)
}

func TestEntitlement(t *testing.T) {
	rec := testHttpRoute(t, "GET", "/azmarketplaceentitlementscount", nil)
	assert.Equal(t, 200, rec.Code)
	var entitlementResponse map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &entitlementResponse); err != nil {
		t.Error(err.Error())
	}
	// json treats all numbers as floats, hence .(float64) to get the number and then convert to int
	assert.Equal(t, 3, int(entitlementResponse["count"].(float64)))
	assert.Equal(t, "", entitlementResponse["error"].(string))
}

func TestErrorMessageInEntitlements(t *testing.T) {
	t.Cleanup(func() {
		database.Instance.Delete("error_message != ''")
	})
	database.Instance.Save(&model.AzureMarketplaceEntitlement{
		ErrorMessage: "Failed to reach Red Hat API",
	})
	rec := testHttpRoute(t, "GET", "/azmarketplaceentitlementscount", nil)
	assert.Equal(t, 200, rec.Code)
	var entitlementResponse map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &entitlementResponse); err != nil {
		t.Error(err.Error())
	}
	// json treats all numbers as floats, hence .(float64) to get the number and then convert to int
	assert.Equal(t, 0, int(entitlementResponse["count"].(float64)))
	assert.Equal(t, "Failed to reach Red Hat API", entitlementResponse["error"].(string))
}

func testHttpRoute(t *testing.T, method string, path string, body io.Reader) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		t.Fatal(err)
	}
	handler.ConfigureAuthenticationForTesting(true)
	rec := httptest.NewRecorder()

	installer := api.NewApp(database, nil, handler.CredentialsHandler{})
	installer.GetRouter().ServeHTTP(rec, req)

	return rec
}

func TestSsoHandler(t *testing.T) {
	// Set up SSO db store
	db := persistence.NewInMemoryDB()
	ssoStore := model.InitSsoStore(db.Instance)
	// Turn on SSO
	config.EnableSso()
	// Start mocked OIDC server
	ssoServer, _ := mockoidc.Run()
	// Save autogenerated credentials
	err := ssoStore.SetSsoClientCredentials(ssoServer.ClientID, ssoServer.ClientSecret)
	assert.Nil(t, err, "Unable to store SSO credentials")

	// Start server to receive SSO callback and grab random code
	var code string
	redirServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code = r.FormValue("code")
	}))
	// Create authenticator with redirect to above server
	authCfg := handler.AuthenticatorConfig{
		Context:      context.Background(),
		SsoEndpoint:  ssoServer.Issuer(),
		RedirectUrl:  redirServer.URL,
		ClientId:     ssoServer.ClientID,
		ClientSecret: ssoServer.ClientSecret,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		Audience:     ssoServer.ClientID,
	}
	auth, err := handler.NewAuthenticator(authCfg)
	assert.Nil(t, err, "Error from authenticator initialization")

	// Create login request
	req, err := http.NewRequest(http.MethodGet, "/login", nil)
	assert.Nil(t, err)

	// Get SSO login handler and create test API router
	ssoHandler := handler.GetSsoHandler(auth)
	installer := api.NewApp(database, nil, ssoHandler)
	// Call test router with login request
	rec := httptest.NewRecorder()
	installer.GetRouter().ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Result().StatusCode)

	// Get auth redirect response
	redirect := rec.Header().Get("Location")
	authUrl, _ := url.PathUnescape(redirect)
	// Auth with fake server to intercept code
	resp, err := http.PostForm(authUrl, nil)
	assert.Nil(t, err, "Unable to post to auth URL")
	require.Equal(t, 200, resp.StatusCode)
	// Call our own SSO callback handler with intercepted code and proper state
	callbackUrl := fmt.Sprintf("/ssocallback?code=%s&state=%s&session_state=unused", code, url.QueryEscape(ssoHandler.State))
	req, _ = http.NewRequest(http.MethodGet, callbackUrl, nil)
	rec = httptest.NewRecorder()
	installer.GetRouter().ServeHTTP(rec, req)
	var body []byte
	n, err := rec.Result().Body.Read(body)
	assert.Nil(t, err)
	assert.Empty(t, string(body))
	assert.Equal(t, 0, n)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Result().StatusCode, "Expected callback to redirect")

	redirUrl := rec.Header().Get("Location")
	assert.Contains(t, redirUrl, "/")

	// Test SSO session checker
	require.NotEmpty(t, rec.Result().Cookies())
	authCookie := rec.Result().Cookies()[0]
	req, _ = http.NewRequest(http.MethodGet, "/authstatus", nil)
	req.AddCookie(authCookie)
	rec = httptest.NewRecorder()
	handler.AuthStatusHandler(rec, req)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode, "Expected OK from authcheck")
	assert.Equal(t, "true", rec.Result().Header.Get("X-Session-Creds"), "Expected true in X-Session-Creds response header")
	assert.Equal(t, "true", rec.Result().Header.Get("X-Session-SSO"), "Expected true in X-Session-SSO response header")

	req, _ = http.NewRequest(http.MethodGet, "/authstatus", nil)
	rec = httptest.NewRecorder()
	handler.AuthStatusHandler(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Result().StatusCode, "Expected Unauthorized from authcheck")
	assert.Equal(t, "false", rec.Result().Header.Get("X-Session-Creds"), "Expected false in X-Session-Creds response header")
	assert.Equal(t, "false", rec.Result().Header.Get("X-Session-SSO"), "Expected false in X-Session-SSO response header")

	// Test invalid state
	callbackUrl = fmt.Sprintf("/ssocallback?code=%s&state=%s&session_state=unused", code, "BADSTATE")
	req, _ = http.NewRequest(http.MethodGet, callbackUrl, nil)
	rec = httptest.NewRecorder()
	installer.GetRouter().ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Result().StatusCode, "Expected unauthorized due to bad state.")

	// Test invalid code
	callbackUrl = fmt.Sprintf("/ssocallback?code=%s&state=%s&session_state=unused", "BADCODE", url.QueryEscape(ssoHandler.State))
	req, _ = http.NewRequest(http.MethodGet, callbackUrl, nil)
	rec = httptest.NewRecorder()
	installer.GetRouter().ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode, "Expected error due to bad code.")
}

func TestRedirectGenerator(t *testing.T) {
	test.SetEnvironment()
	url := handler.GetRedirectUrl()
	assert.Equal(t, "https://localhost/ssocallback", url, "Expected proper redirect URL.")
}

func TestAuthType(t *testing.T) {
	config.DisableSso()
	rec := testHttpRoute(t, "GET", "/authtype", nil)
	assert.Equal(t, 200, rec.Code)
	var authResp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &authResp)
	assert.Nil(t, err)
	assert.Equal(t, "CREDENTIALS", authResp["authtype"])

	config.EnableSso()
	rec = testHttpRoute(t, "GET", "/authtype", nil)
	assert.Equal(t, 200, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), &authResp)
	assert.Nil(t, err)
	assert.Equal(t, "SSO", authResp["authtype"])
}

func TestSessionHelper(t *testing.T) {
	// First test regular authentication
	// Set up session helper
	key, _ := handler.GenerateSessionAuthKey()
	cfg := handler.SessionHelperConfiguration{
		AuthKey:      key,
		CookieName:   "cookie",
		CookiePath:   "/",
		CookieDomain: "localhost",
		Secure:       false,
		MaxAge:       10,
	}
	handler.ConfigureSessionHelper(cfg)

	// Create next handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
	// Ensure logged out
	resp := testHttpRoute(t, http.MethodPost, "/logout", nil)
	assert.Equal(t, http.StatusOK, resp.Result().StatusCode, "Expected OK from logout.")

	handler.ConfigureAuthenticationForTesting(false)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	underTest := handler.EnsureAuthenticated(nextHandler)
	underTest.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Result().StatusCode, "Expected unauthorized due to missing cookie.")

	loginBody := make(map[string]interface{})
	loginBody["uid"] = "admin"
	loginBody["pwd"] = config.GetEnvironment().PASSWORD
	jsonLoginBody, _ := json.Marshal(&loginBody)
	resp = testHttpRoute(t, http.MethodPost, "/login", bytes.NewReader(jsonLoginBody))
	assert.Equal(t, http.StatusOK, resp.Result().StatusCode, "Expected OK from login.")

	validCookie := http.Cookie{
		Name:   "cookie",
		Domain: "localhost",
		Path:   "/",
		Secure: false,
		MaxAge: 10,
	}
	req.AddCookie(&validCookie)
	rec = httptest.NewRecorder()
	underTest.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode, "Expected OK from auth with cookie.")
}

func TestEngineConfiguration(t *testing.T) {
	timeouts := model.EngineConfiguration{}
	resp := testHttpRoute(t, http.MethodGet, "/engineconfiguration", nil)
	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
	if err := json.Unmarshal(resp.Body.Bytes(), &timeouts); err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, config.GetEnvironment().AZURE_DEPLOYMENT_STEP_TIMEOUT, timeouts.StepDeploymentTimeout)
	assert.Equal(t, config.GetEnvironment().ENGINE_END_WAIT, timeouts.EngineExitDelay)
	assert.Equal(t, config.GetEnvironment().ENGINE_MAX_RUNTIME, timeouts.OverallTimeout)
	assert.Equal(t, config.GetEnvironment().EXECUTION_MAX_RETRY, timeouts.StepMaxRetries)
	assert.Equal(t, config.GetEnvironment().AUTO_RETRY_DELAY, timeouts.AutoRetryDelay)
	assert.Equal(t, config.GetEnvironment().ENGINE_RETRY_WAIT, timeouts.StepRestartTimeout)
}

// TestMain wraps the tests.  Setup is done before the call to m.Run() and any
// needed teardown after that.
func TestMain(m *testing.M) {
	test.SetEnvironment()
	database = persistence.NewInMemoryDB()

	execution1 := model.Execution{
		StepID:            1,
		Status:            model.Failed,
		ProvisioningState: string(armresources.ProvisioningStateFailed),
	}

	execution2 := model.Execution{
		StepID:            1,
		Status:            model.Succeeded,
		ProvisioningState: string(armresources.ProvisioningStateSucceeded),
	}

	step1 := model.Step{
		Name:     "step1",
		Template: datatypes.JSONMap{},
		Priority: 0,
	}

	step2 := model.Step{
		Name:     "step2",
		Template: datatypes.JSONMap{},
		Priority: 1,
	}

	step3 := model.Step{
		Name:     "step3",
		Template: datatypes.JSONMap{},
		Priority: 1,
	}

	step4 := model.Step{
		Name:     "step4",
		Template: datatypes.JSONMap{},
		Priority: 2,
	}

	database.Instance.Save(&step1)
	database.Instance.Save(&step2)
	database.Instance.Save(&step3)
	database.Instance.Save(&step4)

	database.Instance.Save(&execution1)
	database.Instance.Save(&execution2)

	database.Instance.Save(&model.AzureMarketplaceEntitlement{
		AzureSubscriptionId: "subscription1",
		AzureCustomerId:     "customer1",
		RHEntitlements:      make([]model.RedHatEntitlements, 0),
		RedHatAccountId:     "",
		Status:              "SUBSCRIBED",
	})

	database.Instance.Save(&model.AzureMarketplaceEntitlement{
		AzureSubscriptionId: "subscription2",
		AzureCustomerId:     "customer1",
		RHEntitlements:      make([]model.RedHatEntitlements, 0),
		RedHatAccountId:     "",
		Status:              "SUBSCRIBED",
	})

	database.Instance.Save(&model.AzureMarketplaceEntitlement{
		AzureSubscriptionId: "subscription3",
		AzureCustomerId:     "customer1",
		RHEntitlements:      make([]model.RedHatEntitlements, 0),
		RedHatAccountId:     "",
		Status:              "SUBSCRIBED",
	})

	m.Run()
}
