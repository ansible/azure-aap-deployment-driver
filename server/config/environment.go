package config

import (
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/elsevierlabs-os/go-envs"
)

type envVars struct {
	SUBSCRIPTION                        string
	RESOURCE_GROUP_NAME                 string
	CONTAINER_GROUP_NAME                string
	STORAGE_ACCOUNT_NAME                string
	PASSWORD                            string
	BASE_PATH                           string
	DB_REL_PATH                         string
	TEMPLATE_REL_PATH                   string
	MAIN_OUTPUTS                        string
	ENGINE_END_WAIT                     int64
	ENGINE_MAX_RUNTIME                  int64
	ENGINE_RETRY_WAIT                   int64
	EXECUTION_MAX_RETRY                 int
	AZURE_POLLING_FREQ_SECONDS          int
	AZURE_DEPLOYMENT_STEP_TIMEOUT_MIN   int
	AUTO_RETRY                          bool
	AUTO_RETRY_DELAY                    int
	AUTH_TYPE                           string
	SESSION_COOKIE_NAME                 string
	SESSION_COOKIE_PATH                 string
	SESSION_COOKIE_DOMAIN               string
	SESSION_COOKIE_SECURE               bool
	SESSION_COOKIE_MAX_AGE              int
	SAVE_CONTAINER                      bool
	SEGMENT_WRITE_KEY                   string
	APPLICATION_ID                      string
	START_TIME                          string
	LOG_REL_PATH                        string
	LOG_LEVEL                           string
	AZURE_LOGIN_RETRIES                 int
	SW_SUB_API_PRIVATEKEY               string
	SW_SUB_API_CERTIFICATE              string
	SW_SUB_API_URL                      string
	SW_SUB_VENDOR_PRODUCT_CODE          string
	AZURE_TENANT_ID                     string
	AZURE_MARKETPLACE_FUNCTION_BASE_URL string
	AZURE_MARKETPLACE_FUNCTION_KEY      string
	INSTALLER_DOMAIN_NAME               string
	SSO_ENDPOINT                        string
	SSO_CLIENT_ID                       string
	SSO_CLIENT_SECRET                   string
}

var (
	environment envVars
)

func GetEnvironment() envVars {
	if (environment != envVars{}) {
		return environment
	}

	// setting defaults here
	environment.ENGINE_END_WAIT = 900     // 15 minutes wait before server exits after its done
	environment.ENGINE_RETRY_WAIT = 1800  // 30 minutes wait for a step to be restarted
	environment.ENGINE_MAX_RUNTIME = 7200 // 2 hours max run time for everything (including restarts)
	environment.EXECUTION_MAX_RETRY = 10  // 10 executions in total allowed
	environment.BASE_PATH = "/installerstore"
	environment.DB_REL_PATH = "installer.db"    // on top of BASE_PATH
	environment.TEMPLATE_REL_PATH = "templates" // on top of BASE_PATH
	environment.AZURE_POLLING_FREQ_SECONDS = 5
	environment.AZURE_DEPLOYMENT_STEP_TIMEOUT_MIN = 30
	environment.AUTO_RETRY = false
	environment.AUTO_RETRY_DELAY = 60 // Retry after 60 seconds if AUTO_RETRY set
	environment.SESSION_COOKIE_NAME = "madd_session"
	environment.SESSION_COOKIE_PATH = "/"
	environment.SESSION_COOKIE_DOMAIN = ""
	environment.SESSION_COOKIE_SECURE = true
	environment.SESSION_COOKIE_MAX_AGE = 0 // 0 to make it a session cookie
	environment.SAVE_CONTAINER = false
	environment.START_TIME = time.Now().Format(time.RFC3339)
	environment.LOG_REL_PATH = "engine.log" // on top of BASE_PATH
	environment.LOG_LEVEL = "info"
	environment.AZURE_LOGIN_RETRIES = 10
	environment.SW_SUB_API_CERTIFICATE = ""
	environment.SW_SUB_API_PRIVATEKEY = ""
	environment.SW_SUB_API_URL = "https://ibm-entitlement-gateway.api.redhat.com/v1/partnerSubscriptions"
	environment.SW_SUB_VENDOR_PRODUCT_CODE = "rhaapomsa"
	environment.SSO_ENDPOINT = "https://sso.stage.redhat.com/auth/realms/redhat-external"
	environment.AZURE_MARKETPLACE_FUNCTION_BASE_URL = "https://marketplace-notification.azurewebsites.net/api/resource"

	env := envs.EnvConfig{}
	env.ReadEnvs()

	environment.SUBSCRIPTION = env.Get("AZURE_SUBSCRIPTION_ID")
	if environment.SUBSCRIPTION == "" {
		log.Fatal("AZURE_SUBSCRIPTION_ID environment variable must be set.")
	}

	environment.AZURE_TENANT_ID = env.Get("AZURE_TENANT_ID")
	if environment.AZURE_TENANT_ID == "" {
		log.Fatal("AZURE_TENANT_ID environment variable must be set.")
	}

	environment.RESOURCE_GROUP_NAME = env.Get("RESOURCE_GROUP_NAME")
	if environment.RESOURCE_GROUP_NAME == "" {
		log.Fatal("RESOURCE_GROUP_NAME environment variable must be set.")
	}

	environment.STORAGE_ACCOUNT_NAME = env.Get("STORAGE_ACCOUNT_NAME")
	if environment.STORAGE_ACCOUNT_NAME == "" {
		log.Fatal("STORAGE_ACCOUNT_NAME environment variable must be set.")
	}

	environment.CONTAINER_GROUP_NAME = env.Get("CONTAINER_GROUP_NAME")
	if environment.CONTAINER_GROUP_NAME == "" {
		log.Fatal("CONTAINER_GROUP_NAME environment variable must be set.")
	}

	environment.PASSWORD = env.Get("ADMIN_PASS")
	if environment.PASSWORD == "" {
		log.Fatal("ADMIN_PASS environment variable must be set.")
	}

	mainOutputsString := env.Get("MAIN_OUTPUTS")
	if mainOutputsString == "" {
		log.Fatal("MAIN_OUTPUTS environment variable must be set.")
	}
	environment.MAIN_OUTPUTS = mainOutputsString

	authType := env.Get("AUTH_TYPE")
	if strings.EqualFold(authType, "sso") {
		environment.AUTH_TYPE = "SSO"
	} else {
		environment.AUTH_TYPE = "CREDENTIALS"
	}

	sessionCookieName := env.Get("SESSION_COOKIE_NAME")
	if sessionCookieName != "" {
		environment.SESSION_COOKIE_NAME = sessionCookieName
	}

	sessionCookiePath := env.Get("SESSION_COOKIE_PATH")
	if sessionCookiePath != "" {
		environment.SESSION_COOKIE_PATH = sessionCookiePath
	}

	sessionCookieDomain := env.Get("SESSION_COOKIE_DOMAIN")
	if sessionCookieDomain != "" {
		environment.SESSION_COOKIE_DOMAIN = sessionCookieDomain
	}

	sessionCookieSecure, err := strconv.ParseBool(env.Get("SESSION_COOKIE_SECURE", "true"))
	if err != nil {
		log.Warnf("SESSION_COOKIE_SECURE has unrecognized value, please set to true or false.  Using default: %t", environment.SESSION_COOKIE_SECURE)
	} else {
		environment.SESSION_COOKIE_SECURE = sessionCookieSecure
	}

	sessionCookieMaxAge, err := strconv.ParseInt(env.Get("SESSION_COOKIE_MAX_AGE", "0"), 10, 32)
	if err != nil {
		log.Warnf("SESSION_COOKIE_MAX_AGE environment variable is not a number, will use default of %d", environment.SESSION_COOKIE_MAX_AGE)
	} else if sessionCookieMaxAge != 0 {
		environment.SESSION_COOKIE_MAX_AGE = int(sessionCookieMaxAge)
	}

	azureLoginRetries, err := strconv.ParseInt(env.Get("AZURE_LOGIN_RETRIES", "0"), 10, 32)
	if err != nil {
		log.Warnf("AZURE_LOGIN_RETRIES environment variable is not a number, will use default of %d", environment.AZURE_LOGIN_RETRIES)
	} else if azureLoginRetries != 0 {
		environment.AZURE_LOGIN_RETRIES = int(azureLoginRetries)
	}

	azureStepTimeoutMin, err := strconv.ParseInt(env.Get("AZURE_DEPLOYMENT_STEP_TIMEOUT_MIN", "0"), 10, 32)
	if err != nil {
		log.Warnf("AZURE_DEPLOYMENT_STEP_TIMEOUT_MIN environment variable is not a number, will use default of %d", environment.AZURE_DEPLOYMENT_STEP_TIMEOUT_MIN)
	} else if azureStepTimeoutMin != 0 {
		environment.AZURE_DEPLOYMENT_STEP_TIMEOUT_MIN = int(azureStepTimeoutMin)
	}

	basePath := env.Get("BASE_PATH")
	if len(basePath) > 0 {
		environment.BASE_PATH = basePath
	}

	dbPath := env.Get("DB_REL_PATH")
	if len(dbPath) > 0 {
		environment.DB_REL_PATH = dbPath
	}

	logPath := env.Get("LOG_REL_PATH")
	if len(logPath) > 0 {
		environment.LOG_REL_PATH = logPath
	}

	logLevel := env.Get("LOG_LEVEL")
	if len(logLevel) > 0 {
		environment.LOG_LEVEL = logLevel
	}

	templatePath := env.Get("TEMPLATE_REL_PATH")
	if len(templatePath) > 0 {
		environment.TEMPLATE_REL_PATH = templatePath
	}

	installerDomainName := env.Get("INSTALLER_DOMAIN_NAME")
	if len(installerDomainName) > 0 {
		environment.INSTALLER_DOMAIN_NAME = installerDomainName
	}

	ssoEndpoint := env.Get("SSO_ENDPOINT")
	if len(ssoEndpoint) > 0 {
		environment.SSO_ENDPOINT = ssoEndpoint
	}

	ssoClientId := env.Get("SSO_CLIENT_ID")
	if len(ssoClientId) > 0 {
		environment.SSO_CLIENT_ID = ssoClientId
	}

	ssoClientSecret := env.Get("SSO_CLIENT_SECRET")
	if len(ssoClientSecret) > 0 {
		environment.SSO_CLIENT_SECRET = ssoClientSecret
	}

	// using empty string as default to force error condition and use of default when env variable not set
	engineEndWait, err := strconv.ParseInt(env.Get("ENGINE_END_WAIT", ""), 10, 64)
	if err != nil {
		log.Warnf("ENGINE_END_WAIT environment variable is not set or is not a number, will use default: %d", environment.ENGINE_END_WAIT)
	} else if engineEndWait != 0 {
		environment.ENGINE_END_WAIT = engineEndWait
	}

	// using empty string as default to force error condition and use of default when env variable not set
	engineMaxRunTime, err := strconv.ParseInt(env.Get("ENGINE_MAX_RUNTIME", ""), 10, 64)
	if err != nil {
		log.Warnf("ENGINE_MAX_RUNTIME environment variable is not set or is not a number, will use default: %d", environment.ENGINE_MAX_RUNTIME)
	} else if engineMaxRunTime != 0 {
		environment.ENGINE_MAX_RUNTIME = engineMaxRunTime
	}

	// using empty string as default to force error condition and use of default when env variable not set
	engineRetryWait, err := strconv.ParseInt(env.Get("ENGINE_RETRY_WAIT", ""), 10, 64)
	if err != nil {
		log.Warnf("ENGINE_RETRY_WAIT environment variable is not set or is not a number, will use default: %d", environment.ENGINE_RETRY_WAIT)
	} else if engineMaxRunTime != 0 {
		environment.ENGINE_RETRY_WAIT = engineRetryWait
	}

	// using empty string as default to force error condition and use of default when env variable not set
	executionMaxRetry, err := strconv.ParseInt(env.Get("EXECUTION_MAX_RETRY", ""), 10, 32)
	if err != nil {
		log.Warnf("EXECUTION_MAX_RETRY environment variable is not set or is not a number, will use default: %d", environment.EXECUTION_MAX_RETRY)
	} else {
		environment.EXECUTION_MAX_RETRY = int(executionMaxRetry)
	}

	// using empty string as default to force error condition and use of default when env variable not set
	azurePollingFreq, err := strconv.ParseInt(env.Get("AZURE_POLLING_FREQ_SECONDS", ""), 10, 32)
	if err != nil {
		log.Warnf("AZURE_POLLING_FREQ_SECONDS environment variable is not set or is not a number, will use default: %d", environment.AZURE_POLLING_FREQ_SECONDS)
	} else if azurePollingFreq > 1 {
		environment.AZURE_POLLING_FREQ_SECONDS = int(azurePollingFreq)
	}

	autoRetry, err := strconv.ParseBool(env.Get("AUTO_RETRY", "false"))
	if err != nil {
		log.Warnf("AUTO_RETRY has unrecognized value, please set to true or false. Using default: %t", environment.AUTO_RETRY)
	} else {
		environment.AUTO_RETRY = autoRetry
	}

	// using empty string as default to force error condition and use of default when env variable not set
	autoRetryDelay, err := strconv.ParseInt(env.Get("AUTO_RETRY_DELAY", ""), 10, 32)
	if err != nil {
		log.Warnf("AUTO_RETRY_DELAY environment variable is not set or is not a number, will use default: %d", environment.AUTO_RETRY_DELAY)
	} else if autoRetryDelay > 1 {
		if autoRetryDelay > int64(environment.ENGINE_RETRY_WAIT) {
			maxAutoRetryDelay := environment.ENGINE_RETRY_WAIT / 2
			log.Warnf("AUTO_RETRY_DELAY cannot exceed ENGINE_RETRY_WAIT, setting to: %d", maxAutoRetryDelay)
			environment.AUTO_RETRY_DELAY = int(maxAutoRetryDelay)
		} else {
			environment.AUTO_RETRY_DELAY = int(autoRetryDelay)
		}
	}

	saveContainer, err := strconv.ParseBool(env.Get("SAVE_CONTAINER", "false"))
	if err != nil {
		log.Warnf("SAVE_CONTAINER has unrecognized value, please set to true or false. Using default: %t", environment.SAVE_CONTAINER)
	} else {
		environment.SAVE_CONTAINER = saveContainer
	}

	environment.SEGMENT_WRITE_KEY = env.Get("SEGMENT_WRITE_KEY")
	if environment.SEGMENT_WRITE_KEY == "" {
		log.Warn("SEGMENT_WRITE_KEY environment variable is either unset or is an empty string, deployment telemetry will not be published to Segment")
	}

	environment.AZURE_MARKETPLACE_FUNCTION_KEY = env.Get("AZURE_MARKETPLACE_FUNCTION_KEY")
	if environment.AZURE_MARKETPLACE_FUNCTION_KEY == "" {
		log.Warn("AZURE_MARKETPLACE_FUNCTION_KEY environment variable is either unset or is an empty string, deployment identification event will not be published")
	}

	environment.APPLICATION_ID = env.Get("APPLICATION_ID")
	if environment.APPLICATION_ID == "" {
		log.Warn("APPLICATION_ID environment variable is either unset or is an empty string, deployment telemetry will not contain the applicationid property")
	}
	environment.START_TIME = env.Get("START_TIME")
	if environment.START_TIME == "" {
		log.Warn("START_TIME environment variable is either unset or is an empty string, telemetry will contain start time of deployment driver engine")
	}

	environment.SW_SUB_API_CERTIFICATE = env.Get("SW_SUB_API_CERTIFICATE")
	if environment.SW_SUB_API_CERTIFICATE == "" {
		log.Warn("SW_SUB_API_CERTIFICATE environment variable is either unset or is an empty string, engine will not be able to make API call for SW subscriptions")
	}

	environment.SW_SUB_API_PRIVATEKEY = env.Get("SW_SUB_API_PRIVATEKEY")
	if environment.SW_SUB_API_PRIVATEKEY == "" {
		log.Warn("SW_SUB_API_PRIVATEKEY environment variable is either unset or is an empty string, engine will not be able to make API call for SW subscriptions")
	}

	swSubApiUrl := env.Get("SW_SUB_API_URL")
	if swSubApiUrl != "" {
		environment.SW_SUB_API_URL = swSubApiUrl
		log.Infof("SW subscription API calls will use URL: %s", swSubApiUrl)
	}

	vendorProductCode := env.Get("SW_SUB_VENDOR_PRODUCT_CODE")
	if vendorProductCode != "" {
		environment.SW_SUB_VENDOR_PRODUCT_CODE = vendorProductCode
		log.Infof("SW subscription API calls will use vendor product code: %s", vendorProductCode)
	}

	return environment
}
