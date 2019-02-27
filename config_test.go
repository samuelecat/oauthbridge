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
	var hook *test.Hook
	log, hook = test.NewNullLogger()

	log.ExitFunc = func(int) {}
	defer func() {
		// restore
		log.ExitFunc = os.Exit
	}()

	Config.getConf("./tests/test_valid.yml")

	assert.False(len(Config.Providers) == 0, "error parsing test configuration file")

	Config.getConf("")
	assert.Equal(logrus.FatalLevel, hook.LastEntry().Level, "parsing invalid configuration file did not throw a fatal")

	Config.getConf("./tests/test_invalid.yml")
	assert.Equal(logrus.FatalLevel, hook.LastEntry().Level, "parsing invalid yaml configuration file did not throw a fatal")
}

func TestLoadConfig(t *testing.T) {
	assert := assert.New(t)
	var hook *test.Hook
	log, hook = test.NewNullLogger()

	log.ExitFunc = func(int) {}
	origIsValidProvider := IsValidProvider
	defer func() {
		// restore
		log.ExitFunc = os.Exit
		IsValidProvider = origIsValidProvider
	}()

	// no providers in the configuration
	Config = configuration{}
	loadConfig("./tests/test_no_providers.yml")
	assert.Equal(logrus.FatalLevel, hook.LastEntry().Level, "parsing a test configuration without providers did not throw a fatal")

	// always valid
	IsValidProvider = func(string) bool {
		return true
	}
	Config = configuration{}
	loadConfig("./tests/test_valid.yml")
	_, ok := Config.Providers["examplebucket"]

	assert.True(ok, "error parsing valid provider for test configuration file")

	// always invalid
	IsValidProvider = func(string) bool {
		return false
	}
	Config = configuration{}
	loadConfig("./tests/test_valid.yml")
	_, ok = Config.Providers["examplebucket"]

	assert.Equal(logrus.FatalLevel, hook.LastEntry().Level, "parsing invalid provider for test configuration file did not throw a fatal")
}

func TestLoadConfigFileNotFound(t *testing.T) {
	assert := assert.New(t)
	var hook *test.Hook
	log, hook = test.NewNullLogger()

	log.ExitFunc = func(int) {}
	origOsStat := osStat
	defer func() {
		// restore
		log.ExitFunc = os.Exit
		osStat = origOsStat
	}()

	// mock, osStat() will return error
	osStat = func(string) (os.FileInfo, error) {
		return os.Stat("-file-not-found-")
	}

	Config = configuration{}
	loadConfig("irrelevant-file-path")

	assert.NotNil(hook.LastEntry())
	assert.Equal(logrus.FatalLevel, hook.LastEntry().Level, "parsing not found configuration file did not throw a fatal")
}
