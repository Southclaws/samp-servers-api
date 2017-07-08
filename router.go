package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
)

// App stores global state for routing
type App struct {
	config Config
	Mongo  *mgo.Session
	Router *mux.Router
	db     *mgo.Collection
}

// Initialise sets up a database connection, binds all the routes and prepares for Start
func Initialise(config Config) *App {
	app := App{
		config: config,
	}

	var err error

	app.Mongo, err = mgo.Dial(fmt.Sprintf("%s:%s", config.MongoHost, config.MongoPort))
	if err != nil {
		logger.Fatal("failed to connect to mongodb",
			zap.Error(err))
	}
	logger.Info("connected to mongodb server")

	err = app.Mongo.Login(&mgo.Credential{
		Source:   config.MongoName,
		Username: config.MongoUser,
		Password: config.MongoPass,
	})
	if err != nil {
		logger.Fatal("failed to log in to mongodb",
			zap.Error(err))
	}
	logger.Info("logged in to mongodb server")

	if !app.CollectionExists("servers") {
		err = app.Mongo.DB(config.MongoName).C("servers").Create(&mgo.CollectionInfo{})
		if err != nil {
			logger.Fatal("collection create failed",
				zap.Error(err))
		}
	}
	app.db = app.Mongo.DB(config.MongoName).C("servers")

	err = app.db.EnsureIndex(mgo.Index{
		Key:         []string{"core.address"},
		Unique:      true,
		DropDups:    true,
		ExpireAfter: time.Hour,
	})
	if err != nil {
		logger.Fatal("index ensure failed",
			zap.Error(err))
	}

	app.Router = mux.NewRouter().StrictSlash(true)

	app.Router.HandleFunc("/server", app.ServerSimple).
		Methods("POST").
		Name("server")

	app.Router.HandleFunc("/server/{address}", app.Server).
		Methods("GET", "POST").
		Name("server")

	app.Router.HandleFunc("/servers", app.Servers).
		Methods("GET").
		Name("servers")

	app.Router.HandleFunc("/players/{address}", app.Players).
		Methods("GET").
		Name("players")

	return &app
}

// Start begins listening for requests and blocks until fatal error
func (app *App) Start() {
	err := http.ListenAndServe(app.config.Bind, app.Router)

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

// WriteError is a utility function for logging a request error and writing a response all in one.
func WriteError(w http.ResponseWriter, status int, err error) {
	logger.Debug("request error", zap.Error(err))
	w.WriteHeader(status)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		logger.Fatal("failed to write error to response", zap.Error(err))
	}
}

// WriteErrors does the same but for groups of errors
func WriteErrors(w http.ResponseWriter, status int, errs []error) {
	logger.Debug("request errors", zap.Errors("errors", errs))
	w.WriteHeader(status)
	for _, err := range errs {
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			logger.Fatal("failed to write error to response", zap.Error(err))
		}
	}
}
