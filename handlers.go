package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var (
	ProviderHandler func(http.ResponseWriter, *http.Request)
)

func init() {
	ProviderHandler = providerHandler
}

func providerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	method := vars["method"]

	// remove the "/provider-method/" part from the path
	path := "/" + strings.TrimPrefix(r.URL.Path, "/"+provider+"-"+method+"/")
	data, err := GetProviderInfo(provider, path, r.URL.RawQuery)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Error(err)
	} else {
		switch method {
		case "redirect":
			http.Redirect(w, r, data["url_full"], http.StatusTemporaryRedirect)
			log.Info("redirect to: " + data["redirect_to"])
		case "info":
			// headers
			w.Header().Set("Content-Type", "application/json")
			for key, value := range data {
				w.Header().Set(Config.HeadersPrefix+key, value)
			}
			w.WriteHeader(http.StatusUnauthorized)
			// body
			str, _ := json.Marshal(data)
			w.Write([]byte(str))
		default:
			http.Error(w, "Bad Request", http.StatusBadRequest)
			log.Error("unknow method provided")
		}
	}
}
