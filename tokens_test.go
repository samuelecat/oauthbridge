package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2/clientcredentials"
)

func TestLoadOAuthClients(t *testing.T) {
	assert := assert.New(t)

	Providers = make(map[string]*Provider)
	Providers["test"] = &Provider{
		BaseURI:      "BaseURI",
		QueryParams:  nil,
		ClientId:     "ClientId",
		ClientSecret: "ClientSecret",
		ExpireTime:   100,
		TokenURL:     "TokenURL",
		Scopes:       []string{"test"},
		Token:        "",
		LastRefresh:  0,
	}

	loadOAuthClients()

	_, ok := OAuthClientsConf["test"]
	assert.True(ok, "test provider was not loaded into oauth clients conf")
}

func TestGetBitbucketTokenNew(t *testing.T) {
	assert := assert.New(t)

	origGetToken := GetToken
	defer func() {
		// restore
		GetToken = origGetToken
	}()

	// disable the real GetToken()
	GetToken = func(c *clientcredentials.Config) (string, error) {
		return "test-token", nil
	}

	Providers = make(map[string]*Provider)
	Providers["bitbucket"] = &Provider{
		Token:       "",
		ExpireTime:  100,
		LastRefresh: 0,
	}
	OAuthClientsConf = make(map[string]*clientcredentials.Config)

	OAuthClientsConf["bitbucket"] = &clientcredentials.Config{
		ClientID:     "ClientId",
		ClientSecret: "ClientSecret",
		Scopes:       []string{"test"},
		TokenURL:     "TokenURL",
	}

	token, _ := getBitbucketToken()

	assert.Equal(token, "test-token", "bitbucket new test-token not retrieved")
}

func TestGetBitbucketTokenCache(t *testing.T) {
	assert := assert.New(t)

	Providers = make(map[string]*Provider)
	Providers["bitbucket"] = &Provider{
		Token:       "test-token",
		ExpireTime:  100,
		LastRefresh: int32(time.Now().Unix()),
	}
	OAuthClientsConf = make(map[string]*clientcredentials.Config)

	OAuthClientsConf["bitbucket"] = &clientcredentials.Config{
		ClientID:     "ClientId",
		ClientSecret: "ClientSecret",
		Scopes:       []string{"test"},
		TokenURL:     "TokenURL",
	}

	token, _ := getBitbucketToken()

	assert.Equal(token, "test-token", "bitbucket test token not found in cache")
}

func TestGetTokenError(t *testing.T) {
	assert := assert.New(t)

	invalidConfig := &clientcredentials.Config{
		ClientID:     "ClientId",
		ClientSecret: "ClientSecret",
		Scopes:       []string{"test"},
		TokenURL:     "",
	}
	_, err := getToken(invalidConfig)
	assert.NotNil(err, "getToken using invalid credential did not fail")
}
