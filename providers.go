// providers
package main

import (
	"encoding/base64"
	"errors"
	"log"
	"net/url"
	"strings"

	oauth_bitbucket "golang.org/x/oauth2/bitbucket"
)

type Provider struct {
	BaseURI      string
	QueryParams  url.Values
	ClientId     string
	ClientSecret string
	ExpireTime   int32
	TokenURL     string
	Scopes       []string

	Token       string // private attribute
	LastRefresh int32  // private attribute
}

var Providers map[string]*Provider

type ProviderInfo map[string]string

func IsValidProvider(name string) bool {
	//GoLang way to check if something is in an array of string
	switch name {
	case
		"bitbucket":
		return true
	}
	return false
}

func loadProviders() {
	Providers = make(map[string]*Provider)

	// load providers from a configuration file
	if val, ok := Config.Providers["bitbucket"]; ok {
		if !strings.Contains(val.BaseURI, "{access_token}") {
			log.Fatalln("error on bitbucket invalid base_uri")
		}

		if val.TokenURL == "" {
			// "https://bitbucket.org/site/oauth2/access_token"
			val.TokenURL = oauth_bitbucket.Endpoint.TokenURL
		}

		if val.ExpireTime == 0 {
			// 2 hours for a bitbucket token to expire, but we set to 1 hour
			val.ExpireTime = 60 * 60
		}

		if len(val.Scopes) == 0 {
			val.Scopes = []string{"repository"}
		}

		Providers["bitbucket"] = &Provider{
			BaseURI:      val.BaseURI,
			QueryParams:  nil,
			ClientId:     val.ClientId,
			ClientSecret: val.ClientSecret,
			ExpireTime:   val.ExpireTime,
			TokenURL:     val.TokenURL,
			Scopes:       val.Scopes,
			Token:        "",
			LastRefresh:  0,
		}
	}
}

func getProviderInfo(provider string, url_path string, url_query string) (ProviderInfo, error) {
	var err error
	var token, url_full string
	data := make(ProviderInfo)

	switch provider {
	case "bitbucket":
		//log.Println("serving provider ", provider)
		bb := Providers[provider]
		token, err = getBitbucketToken()
		if err == nil {
			var u *url.URL
			base_url := strings.Replace(bb.BaseURI, "{access_token}", token, -1)

			u, err = url.Parse(base_url)
			if err == nil {
				credentials := u.User.String()
				// build the url for the redirect to
				url_full = base_url + url_path
				if len(url_query) > 0 {
					url_full += "?" + url_query
				}
				data["token"] = token
				data["url_full"] = url_full
				if credentials != "" {
					data["url_no_auth"] = strings.Replace(url_full, "//"+credentials+"@", "//", 1)
					data["auth_base64"] = base64.URLEncoding.EncodeToString([]byte(credentials))
				}
			}
		}
	default:
		err = errors.New("an invalid provider has been specified in the path")
	}
	return data, err
}
