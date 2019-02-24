// server
package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	// Output to stdout
	log.Out = os.Stdout

	// Log level
	//log.Level = logrus.ErrorLevel
	log.Level = logrus.InfoLevel
	//log.Level = logrus.DebugLevel
}

func main() {
	log.Info("loading configuration...")
	loadConfig("")
	loadProviders()
	loadOAuthClients()

	log.Info("starting the server...")
	r := mux.NewRouter()
	r.HandleFunc("/{provider}-{method}/{path:.*}", providerHandler)
	http.Handle("/", &MyServer{r})
	log.Info("ready")
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
