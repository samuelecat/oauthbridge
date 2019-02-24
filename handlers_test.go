// handlers_test.go
package main

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestProviderHandler(t *testing.T) {
	assert := assert.New(t)
	// disable the real GetProviderInfo()
	origGetProviderInfo := GetProviderInfo
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

	// restore GetProviderInfo()
	GetProviderInfo = origGetProviderInfo
}
