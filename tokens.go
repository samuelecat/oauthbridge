// tokens
package main

import (
	//"log"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var OAuthClientsConf map[string]*clientcredentials.Config

func loadOAuthClients() {
	OAuthClientsConf = make(map[string]*clientcredentials.Config)

	for key, value := range Providers {
		OAuthClientsConf[key] = &clientcredentials.Config{
			ClientID:     value.ClientId,
			ClientSecret: value.ClientSecret,
			Scopes:       value.Scopes,
			TokenURL:     value.TokenURL,
		}
	}
}

func getBitbucketToken() (string, error) {
	// https://developer.atlassian.com/cloud/bitbucket/oauth-2/

	var err error = nil
	current_t := int32(time.Now().Unix())

	if Providers["bitbucket"].Token == "" ||
		(current_t-Providers["bitbucket"].LastRefresh) > Providers["bitbucket"].ExpireTime {

		var token *oauth2.Token = nil
		c := OAuthClientsConf["bitbucket"]
		token, err = c.Token(oauth2.NoContext)
		if err == nil {
			Providers["bitbucket"].Token = token.AccessToken
			Providers["bitbucket"].LastRefresh = current_t
		}
	}

	//log.Println(Providers["bitbucket"].Token)
	return Providers["bitbucket"].Token, err
}
