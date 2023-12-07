package handler_test

import (
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
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/oauth2-proxy/mockoidc"
	"github.com/stretchr/testify/assert"
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
	if err != nil {
		t.Errorf("Unable to store SSO credentials: %v", err)
		return
	}
	// Start server to receive SSO callback and grab random code
	var code string
	redirServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code = r.FormValue("code")
	}))
	// Create authenticator with redirect to above server
	auth, err := handler.NewAuthenticator(context.Background(), ssoServer.Issuer(), redirServer.URL, ssoServer.ClientID, ssoServer.ClientSecret)
	if err != nil {
		t.Errorf("Error from authenticator initialization: %v", err)
		return
	}
	// Create login request
	req, err := http.NewRequest(http.MethodGet, "/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Get SSO login handler and create test API router
	ssoHandler := handler.GetSsoHandler(auth)
	installer := api.NewApp(database, nil, ssoHandler)
	// Call test router with login request
	rec := httptest.NewRecorder()
	installer.GetRouter().ServeHTTP(rec, req)
	if rec.Result().StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("Expected http status temp redirect, got %d", rec.Result().StatusCode)
		return
	}
	// Get auth redirect response
	redirect := rec.Header().Get("Location")
	authUrl, _ := url.PathUnescape(redirect)
	// Auth with fake server to intercept code
	_, err = http.PostForm(authUrl, nil)
	if err != nil {
		t.Errorf("Unable to post to auth URL: %v", err)
		return
	}
	// Call our own SSO callback handler with intercepted code and proper state
	callbackUrl := fmt.Sprintf("/ssocallback?code=%s&state=%s&session_state=unused", code, url.QueryEscape(ssoHandler.State))
	req, _ = http.NewRequest(http.MethodGet, callbackUrl, nil)
	rec = httptest.NewRecorder()
	installer.GetRouter().ServeHTTP(rec, req)
	if rec.Result().StatusCode != http.StatusTemporaryRedirect {
		body, _ := io.ReadAll(rec.Result().Body)
		t.Logf("Failure body: %s", string(body))
		t.Errorf("Expected callback to redirect, got %d", rec.Result().StatusCode)
		return
	}
	redirUrl := rec.Header().Get("Location")
	if !strings.HasSuffix(redirUrl, "deployment") {
		t.Errorf("Expected redirect URL to end in /deployment, was %s", redirUrl)
		return
	}
	// Test invalid state
	callbackUrl = fmt.Sprintf("/ssocallback?code=%s&state=%s&session_state=unused", code, "BADSTATE")
	req, _ = http.NewRequest(http.MethodGet, callbackUrl, nil)
	rec = httptest.NewRecorder()
	installer.GetRouter().ServeHTTP(rec, req)
	if rec.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected unauthorized due to bad state, got %d", rec.Result().StatusCode)
		return
	}
	// Test invalid code
	callbackUrl = fmt.Sprintf("/ssocallback?code=%s&state=%s&session_state=unused", "BADCODE", url.QueryEscape(ssoHandler.State))
	req, _ = http.NewRequest(http.MethodGet, callbackUrl, nil)
	rec = httptest.NewRecorder()
	installer.GetRouter().ServeHTTP(rec, req)
	if rec.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected error due to bad code, got %d", rec.Result().StatusCode)
	}
}

func TestRedirectGenerator(t *testing.T) {
	test.SetEnvironment()
	url := handler.GetRedirectUrl()
	if url != "https://localhost/ssocallback" {
		t.Errorf("Expected https://localhost/ssocallback redirect, got: %s", url)
	}
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
