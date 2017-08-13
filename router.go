package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// App stores global state for routing
type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	config Config
	Mongo  *mgo.Session
	db     *mgo.Collection
	qd     *QueryDaemon
	Router *mux.Router
}

// Initialise sets up a database connection, binds all the routes and prepares for Start
func Initialise(config Config) *App {
	app := App{
		config: config,
	}
	app.ctx, app.cancel = context.WithCancel(context.Background())

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

	addresses, err := app.LoadAllAddresses()
	if err != nil {
		logger.Fatal("failed to load current addresses for query daemon",
			zap.Error(err))
	}

	app.qd = NewQueryDaemon(app.ctx, &app, addresses, time.Second*10, 5)

	app.Router = mux.NewRouter().StrictSlash(true)

	app.Router.HandleFunc("/v1/server", app.ServerSimple).
		Methods("OPTIONS", "POST").
		Name("server")

	app.Router.HandleFunc("/v1/server/{address}", app.Server).
		Methods("OPTIONS", "GET", "POST").
		Name("server")

	app.Router.HandleFunc("/v1/servers", app.Servers).
		Methods("OPTIONS", "GET").
		Name("servers")

	app.Router.HandleFunc("/v1/players/{address}", app.Players).
		Methods("OPTIONS", "GET").
		Name("players")

	return &app
}

// Start begins listening for requests and blocks until fatal error
func (app *App) Start() {
	defer app.cancel()

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	err := http.ListenAndServe(app.config.Bind, handlers.CORS(headersOk, originsOk, methodsOk)(app.Router))

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

// LoadAllAddresses loads all addresses from the database as a slice of strings for synchronisation
// with the QueryDaemon.
func (app *App) LoadAllAddresses() (result []string, err error) {
	allServers := []Server{}
	err = app.db.Find(bson.M{}).All(&allServers)
	if err == nil {
		for i := range allServers {
			result = append(result, allServers[i].Core.Address)
		}
	}
	return
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
		_, err = w.Write([]byte(err.Error() + ", "))
		if err != nil {
			logger.Fatal("failed to write error to response", zap.Error(err))
		}
	}
}
