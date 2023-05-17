package config

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/elsevierlabs-os/go-envs"
)

type envVars struct {
	SUBSCRIPTION               string
	RESOURCE_GROUP_NAME        string
	AZURE_LOCATION             string
	CONTAINER_GROUP_NAME       string
	STORAGE_ACCOUNT_NAME       string
	PASSWORD                   string
	DB_PATH                    string
	TEMPLATE_PATH              string
	MAIN_OUTPUTS               string
	ENGINE_END_WAIT            int64
	ENGINE_MAX_RUNTIME         int64
	ENGINE_RETRY_WAIT          int64
	EXECUTION_MAX_RETRY        int
	AZURE_POLLING_FREQ_SECONDS int
	AUTO_RETRY                 bool
	AUTO_RETRY_DELAY           int
	SESSION_COOKIE_NAME        string
	SESSION_COOKIE_PATH        string
	SESSION_COOKIE_DOMAIN      string
	SESSION_COOKIE_SECURE      bool
	SESSION_COOKIE_MAX_AGE     int
	SAVE_CONTAINER             bool
	SEGMENT_WRITE_KEY          string
	APPLICATION_ID             string
	START_TIME                 string
	WEB_HOOK_API_KEY           string
	WEB_HOOK_CALLBACK_URL      string
	LOG_PATH                   string
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
	environment.DB_PATH = "/installerstore/installer.db"
	environment.TEMPLATE_PATH = "/installerstore/templates"
	environment.AZURE_POLLING_FREQ_SECONDS = 5
	environment.AUTO_RETRY = false
	environment.AUTO_RETRY_DELAY = 60 // Retry after 60 seconds if AUTO_RETRY set
	environment.SESSION_COOKIE_NAME = "madd_session"
	environment.SESSION_COOKIE_PATH = "/"
	environment.SESSION_COOKIE_DOMAIN = ""
	environment.SESSION_COOKIE_SECURE = true
	environment.SESSION_COOKIE_MAX_AGE = 0 // 0 to make it a session cookie
	environment.SAVE_CONTAINER = false
	environment.START_TIME = time.Now().Format(time.RFC3339)
	environment.LOG_PATH = "/installerstore/engine.txt"

	// TODO: need to set this to a real value that's not hardcoded
	environment.WEB_HOOK_API_KEY = "6P7Q9SATBVDWEXGZH2J4M5N6Q8"
	environment.WEB_HOOK_CALLBACK_URL = fmt.Sprintf("http://localhost:%s/eventhook", Args.Port)

	env := envs.EnvConfig{}
	env.ReadEnvs()

	environment.SUBSCRIPTION = env.Get("AZURE_SUBSCRIPTION_ID")
	if environment.SUBSCRIPTION == "" {
		log.Fatal("AZURE_SUBSCRIPTION_ID environment variable must be set.")
	}

	environment.RESOURCE_GROUP_NAME = env.Get("RESOURCE_GROUP_NAME")
	if environment.RESOURCE_GROUP_NAME == "" {
		log.Fatal("RESOURCE_GROUP_NAME environment variable must be set.")
	}

	environment.AZURE_LOCATION = env.Get("AZURE_LOCATION")
	if environment.AZURE_LOCATION == "" {
		log.Fatal("AZURE_LOCATION environment variable must be set.")
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

	dbPath := env.Get("DB_PATH")
	if len(dbPath) > 0 {
		environment.DB_PATH = dbPath
	}

	logPath := env.Get("LOG_PATH")
	if len(logPath) > 0 {
		environment.LOG_PATH = logPath
	}

	templatePath := env.Get("TEMPLATE_PATH")
	if len(templatePath) > 0 {
		environment.TEMPLATE_PATH = templatePath
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
	environment.APPLICATION_ID = env.Get("APPLICATION_ID")
	if environment.APPLICATION_ID == "" {
		log.Warn("APPLICATION_ID environment variable is either unset or is an empty string, deployment telemetry will not contain the applicationid property")
	}
	environment.START_TIME = env.Get("START_TIME")
	if environment.START_TIME == "" {
		log.Warn("START_TIME environment variable is either unset or is an empty string, telemetry will contain start time of deployment driver engine")
	}

	return environment
}
