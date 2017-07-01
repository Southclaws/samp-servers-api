package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Server stores the standard SA:MP query fields as well as an additional details type that stores
// additional details implemented by this API and modern server browsers.
// The json keys are short to cut down on network traffic.
type Server struct {
	Address    string            `json:"ip"`
	Hostname   string            `json:"hn"`
	Players    int               `json:"pc"`
	MaxPlayers int               `json:"pm"`
	Gamemode   string            `json:"gm"`
	Language   string            `json:"la"`
	Password   bool              `json:"pa"`
	Rules      map[string]string `json:"ru"`
	PlayerList []string          `json:"pl"`
}

// Server handles either posting a server object or requesting a server object
func (app *App) Server(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["address"]
	if !ok {
		logger.Fatal("no address specified in request",
			zap.String("request", r.URL.String()))
	}

	switch r.Method {
	case "GET":
		logger.Debug("getting server",
			zap.String("address", address))

	case "POST":
		logger.Debug("posting server",
			zap.String("address", address))
	}
}
