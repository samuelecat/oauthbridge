// server
package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var ServerStart func()

func init() {
	ServerStart = serverStart

	// Output to stdout
	log.Out = os.Stdout

	// Log level
	//log.Level = logrus.ErrorLevel
	log.Level = logrus.InfoLevel
	//log.Level = logrus.DebugLevel
}

func main() {
	log.Info("loading configuration...")
	LoadConfig("")
	LoadProviders()
	LoadOAuthClients()

	log.Info("starting the server...")
	ServerStart()
}

func serverStart() {
	r := mux.NewRouter()
	r.HandleFunc("/{provider}-{method}/{path:.*}", ProviderHandler)
	http.Handle("/", &MyServer{r})
	log.Info("ready")
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Fatalln(err)
	}
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
