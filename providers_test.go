// providers_test.go
package main

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2/clientcredentials"
)

func TestIsValidProvider(t *testing.T) {
	assert := assert.New(t)

	test := IsValidProvider("not-valid")
	assert.False(test, "an invalid Provider passed the test IsValidProvider()")

	test = IsValidProvider("bitbucket")
	assert.True(test, "a valid Provider did not pass the test IsValidProvider()")
}

type mockProviders map[string]struct {
	BaseURI      string   `yaml:"base_uri"`
	QueryParams  []string `yaml:"query_params"`
	ClientId     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	Scopes       []string `yaml:"scopes"`
	TokenURL     string   `yaml:"token_url"`
	ExpireTime   int32    `yaml:"expire_time"`
}

func MockProviders() mockProviders {
	p := make(mockProviders)
	bb := p["bitbucket"]
	bb.BaseURI = "https://x-token-auth:{access_token}@bitbucket.org"
	bb.QueryParams = nil
	bb.ClientId = "ClientId"
	bb.ClientSecret = "ClientSecret"
	bb.Scopes = []string{"test"}
	bb.TokenURL = "TokenURL"
	bb.ExpireTime = 100
	p["bitbucket"] = bb
	return p
}

func TestLoadProviders(t *testing.T) {
	assert := assert.New(t)
	var hook *test.Hook
	log, hook = test.NewNullLogger()

	// mocking a pseudo valid structure
	p := MockProviders()
	Config = configuration{Providers: p}

	loadProviders()
	_, ok := Providers["bitbucket"]
	assert.True(ok, "a valid provider was not loaded")

	// mocking valid structure relying on default values
	p = MockProviders()
	bb := p["bitbucket"]
	bb.TokenURL = ""
	bb.ExpireTime = 0
	Config = configuration{Providers: p}

	loadProviders()
	data, ok := Providers["bitbucket"]
	assert.True(ok, "a valid provider was not loaded")
	assert.NotNil(data.TokenURL, "TokenURL was not set from default value")
	assert.NotEqual(data.ExpireTime, 0, "ExpireTime was not set from default value")

	p = MockProviders()
	bb = p["bitbucket"]

	bb.BaseURI = "not-valid-for-bitbucket"
	p["bitbucket"] = bb

	Config = configuration{Providers: p}

	log.ExitFunc = func(int) {}
	defer func() {
		// restore
		log.ExitFunc = os.Exit
	}()

	loadProviders()

	assert.Equal(logrus.FatalLevel, hook.LastEntry().Level)

	hook.Reset()
}

func TestGetProviderInfo(t *testing.T) {
	assert := assert.New(t)
	log, _ = test.NewNullLogger()

	// mocking a pseudo valid structure
	p := MockProviders()
	Config = configuration{Providers: p}
	loadProviders()

	// disable the real GetToken()
	origGetToken := GetToken
	defer func() {
		// restore
		GetToken = origGetToken
	}()

	GetToken = func(c *clientcredentials.Config) (string, error) {
		return "test-token", nil
	}

	_, err := getProviderInfo("", "", "")

	assert.NotNil(err, "GetProviderInfo working using an invalid Provider")

	data, err := getProviderInfo("bitbucket", "/test", "")

	assert.NotNil(data)

	token, ok := data["token"]
	assert.True(ok, "a valid provider was not loaded")
	assert.Equal("test-token", token, "getProviderInfo() did not return our valid test-token")

	_, ok = data["url_no_auth"]
	assert.True(ok, "provider is bitbucket but url_no_auth not found")
}
