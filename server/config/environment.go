package config

import (
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/elsevierlabs-os/go-envs"
)

type envVars struct {
	SUBSCRIPTION               string
	RESOURCE_GROUP_NAME        string
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
}

var environment envVars

func GetEnvironment() envVars {
	if (environment != envVars{}) {
		return environment
	}

	// setting defaults here
	environment.ENGINE_END_WAIT = 900     // 15 minutes wait before server exits after its done
	environment.ENGINE_RETRY_WAIT = 1800  // 30 minutes wait for a step to be restarted
	environment.ENGINE_MAX_RUNTIME = 7200 // 2 hours max run time for everything (including restarts)
	environment.EXECUTION_MAX_RETRY = 3
	environment.DB_PATH = "/installerstore/installer.db"
	environment.TEMPLATE_PATH = "/installerstore/templates"
	environment.AZURE_POLLING_FREQ_SECONDS = 5
	environment.AUTO_RETRY = false
	environment.AUTO_RETRY_DELAY = 60 // Retry after 60 seconds if AUTO_RETRY set

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

	dbPath := env.Get("DB_PATH")
	if len(dbPath) > 0 {
		environment.DB_PATH = dbPath
	}

	templatePath := env.Get("TEMPLATE_PATH")
	if len(templatePath) > 0 {
		environment.TEMPLATE_PATH = templatePath
	}

	engineEndWait, err := strconv.ParseInt(env.Get("ENGINE_END_WAIT", "0"), 10, 64)
	if err != nil {
		log.Warnf("ENGINE_END_WAIT environment variable is not a number, will use default of %d", environment.ENGINE_END_WAIT)
	} else if engineEndWait != 0 {
		environment.ENGINE_END_WAIT = engineEndWait
	}

	engineMaxRunTime, err := strconv.ParseInt(env.Get("ENGINE_MAX_RUNTIME", "0"), 10, 64)
	if err != nil {
		log.Warnf("ENGINE_MAX_RUNTIME environment variable is not a number, will use default of %d", environment.ENGINE_MAX_RUNTIME)
	} else if engineMaxRunTime != 0 {
		environment.ENGINE_MAX_RUNTIME = engineMaxRunTime
	}

	engineRetryWait, err := strconv.ParseInt(env.Get("ENGINE_RETRY_WAIT", "0"), 10, 64)
	if err != nil {
		log.Warnf("ENGINE_RETRY_WAIT environment variable is not a number, will use default of %d", environment.ENGINE_RETRY_WAIT)
	} else if engineMaxRunTime != 0 {
		environment.ENGINE_RETRY_WAIT = engineRetryWait
	}

	executionMaxRetry, err := strconv.ParseInt(env.Get("EXECUTION_MAX_RETRY", "0"), 10, 32)
	if err != nil {
		log.Warnf("EXECUTION_MAX_RETRY environment variable is not a number, will use default of %d", environment.EXECUTION_MAX_RETRY)
	} else {
		environment.EXECUTION_MAX_RETRY = int(executionMaxRetry)
	}

	azurePollingFreq, err := strconv.ParseInt(env.Get("AZURE_POLLING_FREQ_SECONDS", "0"), 10, 32)
	if err != nil {
		log.Warnf("AZURE_POLLING_FREQ_SECONDS environment variable is not a number, will use default of %d", environment.AZURE_POLLING_FREQ_SECONDS)
	} else if azurePollingFreq > 1 {
		environment.AZURE_POLLING_FREQ_SECONDS = int(azurePollingFreq)
	}

	autoRetry, err := strconv.ParseBool(env.Get("AUTO_RETRY", "false"))
	if err != nil {
		log.Warnf("AUTO_RETRY has unrecognized value, please set to true or false.  Using default: %t", environment.AUTO_RETRY)
	} else {
		environment.AUTO_RETRY = autoRetry
	}

	autoRetryDelay, err := strconv.ParseInt(env.Get("AUTO_RETRY_DELAY", "0"), 10, 32)
	if err != nil {
		log.Warnf("AUTO_RETRY_DELAY environment variable is not a number, will use default of %d", environment.AUTO_RETRY_DELAY)
	} else if autoRetryDelay > 1 {
		environment.AUTO_RETRY_DELAY = int(autoRetryDelay)
	}

	return environment
}
