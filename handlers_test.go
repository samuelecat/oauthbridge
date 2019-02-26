// handlers_test.go
package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestProviderHandler(t *testing.T) {
	assert := assert.New(t)

	origGetProviderInfo := GetProviderInfo
	defer func() {
		// restore
		GetProviderInfo = origGetProviderInfo
	}()

	// replace the real GetProviderInfo()
	GetProviderInfo = func(string, string, string) (ProviderInfo, error) {
		data := make(ProviderInfo)
		data["token"] = "test-token"
		return data, nil
	}

	// https://golang.org/pkg/net/http/httptest/
	r := httptest.NewRequest("GET", "http://localhost:9999/bitbucket-info/foo", nil)
	w := httptest.NewRecorder()
	// inject wanted variables in the request
	r = mux.SetURLVars(r, map[string]string{"provider": "bitbucket", "method": "info"})

	providerHandler(w, r)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(resp.StatusCode, 401, "response is not 401")
	assert.True(strings.Contains(string(body), "\"token\":\"test-token\""), "response does not have our test-token")
}

func TestProviderHandlerErrors(t *testing.T) {
	assert := assert.New(t)
	log, _ = test.NewNullLogger()

	origGetProviderInfo := GetProviderInfo
	defer func() {
		// restore
		GetProviderInfo = origGetProviderInfo
	}()

	GetProviderInfo = func(string, string, string) (ProviderInfo, error) {
		return nil, errors.New("something")
	}

	// https://golang.org/pkg/net/http/httptest/
	r := httptest.NewRequest("GET", "http://localhost:9999/bitbucket-info/foo", nil)
	w := httptest.NewRecorder()
	// inject wanted variables in the request
	r = mux.SetURLVars(r, map[string]string{"provider": "bitbucket", "method": "redirect"})

	providerHandler(w, r)
	resp := w.Result()

	assert.Equal(resp.StatusCode, http.StatusBadRequest, "an error was injected but the response is not 400")

	GetProviderInfo = func(string, string, string) (ProviderInfo, error) {
		data := make(ProviderInfo)
		data["token"] = "test-token"
		return data, nil
	}

	w = httptest.NewRecorder()
	providerHandler(w, r)
	resp = w.Result()

	assert.Equal(resp.StatusCode, http.StatusTemporaryRedirect, "an redirect was injected but the response is not 307")

	// test invalid method
	r = httptest.NewRequest("GET", "http://localhost:9999/bitbucket-info/foo", nil)
	w = httptest.NewRecorder()
	// inject wanted variables in the request
	r = mux.SetURLVars(r, map[string]string{"provider": "bitbucket", "method": "unknown"})

	providerHandler(w, r)
	resp = w.Result()

	assert.Equal(resp.StatusCode, http.StatusBadRequest, "an error was injected but the response is not 400")
}
