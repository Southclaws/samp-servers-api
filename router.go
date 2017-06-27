package main

import (
	"net/http"

	"go.uber.org/zap"
	"github.com/gorilla/mux"
)

// App stores global state for routing
type App struct {
	config Config
	Router *mux.Router
}

// Start binds the routes and starts listening for requests, blocking until fatal error.
func Start(config Config) {
	app := App{
		config: config,
	}

	app.Router = mux.NewRouter().StrictSlash(true)

	app.Router.HandleFunc("/server/{address}", app.Server).Methods("GET", "POST").Name("server")
	app.Router.HandleFunc("/servers", app.Servers).Methods("GET").Name("servers")
	app.Router.HandleFunc("/players/{address}", app.Players).Methods("GET").Name("players")

	err := http.ListenAndServe(config.Bind, app.Router)

	logger.Fatal("http server encountered fatal error",
		zap.Error(err))
}
