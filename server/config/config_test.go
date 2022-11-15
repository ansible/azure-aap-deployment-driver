package config_test

import (
	"bytes"
	"os"
	"server/config"
	"server/test"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func TestArgs(t *testing.T) {
	config.ParseArgs()
	assert.Equal(t, "9090", config.Args.Port)
	assert.Equal(t, "0.0.0.0", config.Args.Host)
}

func TestEnvironment(t *testing.T) {
	envs := config.GetEnvironment()
	assert.Equal(t, "3f7e29ba-24e0-42f6-8d9c-5149a14bda37", envs.SUBSCRIPTION)
}

func TestLogging(t *testing.T) {
	config.ConfigureLogging()
	// Capture logging
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stdout)
	}()
	log.Info("hello")
	assert.Contains(t, buf.String(), "hello")
	assert.Contains(t, buf.String(), "config_test")
	assert.Contains(t, buf.String(), time.Now().UTC().Format("2006-01-02"))
}

// TestMain wraps the tests.  Setup is done before the call to m.Run() and any
// needed teardown after that.
func TestMain(m *testing.M) {
	test.SetEnvironment()
	m.Run()
}
