package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
)

// App stores global state for routing
type App struct {
	config Config
	Mongo  *mgo.Session
	Router *mux.Router
}

// Start binds the routes and starts listening for requests, blocking until fatal error.
func Start(config Config) {
	app := App{
		config: config,
	}

	var err error

	app.Mongo, err = mgo.Dial(config.MongoHost)
	if err != nil {
		logger.Fatal("failed to connect to mongodb",
			zap.Error(err))
	}

	err = app.Mongo.Login(&mgo.Credential{
		Source:   config.MongoName,
		Username: config.MongoUser,
		Password: config.MongoPass,
	})
	if err != nil {
		logger.Fatal("failed to log in to mongodb",
			zap.Error(err))
	}

	if !app.CollectionExists("servers") {
		err = app.Mongo.DB(config.MongoName).C("servers").Create(&mgo.CollectionInfo{})
		if err != nil {
			logger.Fatal("collection create failed",
				zap.Error(err))
		}
	}

	app.Router = mux.NewRouter().StrictSlash(true)

	app.Router.HandleFunc("/server/{address}", app.Server).Methods("GET", "POST").Name("server")
	app.Router.HandleFunc("/servers", app.Servers).Methods("GET").Name("servers")
	app.Router.HandleFunc("/players/{address}", app.Players).Methods("GET").Name("players")

	err = http.ListenAndServe(config.Bind, app.Router)

	logger.Fatal("http server encountered fatal error",
		zap.Error(err))
}

// CollectionExists checks if a collection exists in MongoDB
func (app *App) CollectionExists(name string) bool {
	collections, err := app.Mongo.DB(app.config.MongoName).CollectionNames()
	if err != nil {
		logger.Fatal("failed to get collection names",
			zap.Error(err))
	}

	for _, collection := range collections {
		if collection == name {
			return true
		}
	}

	return false
}
