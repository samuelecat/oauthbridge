// server_test.go
package main

import (
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

var loadC, loadP, loadA, srvS bool

func TestMain(t *testing.T) {
	assert := assert.New(t)
	log, _ = test.NewNullLogger()

	// save original functions
	origLoadConfig := LoadConfig
	origLoadProviders := LoadProviders
	origLoadOAuthClients := LoadOAuthClients
	origServerStart := ServerStart
	defer func() {
		// restore
		LoadConfig = origLoadConfig
		LoadProviders = origLoadProviders
		LoadOAuthClients = origLoadOAuthClients
		ServerStart = origServerStart
	}()

	LoadConfig = func(string) {
		loadC = true
	}
	LoadProviders = func() {
		loadP = true
	}
	LoadOAuthClients = func() {
		loadA = true
	}
	ServerStart = func() {
		srvS = true
	}

	main()

	assert.True(loadC, "LoadConfig() was not execute in main()")
	assert.True(loadP, "LoadProviders() was not execute in main()")
	assert.True(loadA, "LoadOAuthClients() was not execute in main()")
	assert.True(loadA, "ServerStart() was not execute in main()")
}
