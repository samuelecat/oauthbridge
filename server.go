// server
package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("loading configuration...")
	loadConfig()
	loadProviders()
	loadOAuthClients()

	log.Println("starting the server...")
	r := mux.NewRouter()
	r.HandleFunc("/{provider}-{method}/{path:.*}", providerHandler)
	http.Handle("/", &MyServer{r})
	log.Println("ready")
	http.ListenAndServe(":9999", nil)
}

type MyServer struct {
	r *mux.Router
}

func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	s.r.ServeHTTP(rw, req)
}

func providerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	method := vars["method"]

	// remove the "/provider-method/" part from the path
	path := "/" + strings.TrimPrefix(r.URL.Path, "/"+provider+"-"+method+"/")
	data, err := getProviderInfo(provider, path, r.URL.RawQuery)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("error: ", err)
	} else {
		switch method {
		case "redirect":
			http.Redirect(w, r, data["url_full"], http.StatusTemporaryRedirect)
			log.Printf("redirect to: %s", data["redirect_to"])
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
			log.Println("error: ", errors.New("unknow method provided"))
		}
	}
}
