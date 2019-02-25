// config_test.go
package main

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestGetConf(t *testing.T) {
	assert := assert.New(t)
	Config.getConf("./docs/test.yml")

	assert.False(len(Config.Providers) == 0, "error parsing test configuration file")
}

func TestLoadConfig(t *testing.T) {
	assert := assert.New(t)
	var hook *test.Hook
	log, hook = test.NewNullLogger()
	log.ExitFunc = func(int) {}

	// disable the real IsValidProvider()
	origIsValidProvider := IsValidProvider

	// always valid
	IsValidProvider = func(string) bool {
		return true
	}
	Config = configuration{}
	loadConfig("./docs/test.yml")
	_, ok := Config.Providers["examplebucket"]

	assert.True(ok, "error parsing valid provider for test configuration file")

	// always invalid
	IsValidProvider = func(string) bool {
		return false
	}
	Config = configuration{}
	loadConfig("./docs/test.yml")
	_, ok = Config.Providers["examplebucket"]

	assert.Equal(logrus.FatalLevel, hook.LastEntry().Level, "parsing invalid provider for test configuration file did not throw a fatal")

	log.ExitFunc = os.Exit
	// restore GetProviderInfo()
	IsValidProvider = origIsValidProvider
}
