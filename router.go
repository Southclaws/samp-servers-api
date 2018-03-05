package main

import (
	"context"
	"net/http"
	"time"

	"github.com/Southclaws/go-samp-query"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// App stores global state for routing
type App struct {
	ctx        context.Context
	cancel     context.CancelFunc
	config     Config
	collection *mgo.Collection
	qd         *QueryDaemon
	router     *mux.Router
}

// Initialise sets up a database connection, binds all the routes and prepares for Start
func Initialise(config Config) *App {
	logger.Debug("initialising samp-servers-api with debug logging", zap.Any("config", config))

	app := App{
		config: config,
	}
	app.ctx, app.cancel = context.WithCancel(context.Background())

	// Connect to the database, receive a collection pointer
	app.collection = ConnectDB(config)
	logger.Info("connected to mongodb server")

	// Grab existing addresses from database and pass to the Query Daemon
	addresses := app.LoadAllAddresses()
	app.qd = NewQueryDaemon(app.ctx, &app, addresses, time.Second*time.Duration(config.QueryInterval), config.MaxFailedQuery, sampquery.GetServerInfo)

	// Set up HTTP server
	app.router = mux.NewRouter().StrictSlash(true)

	app.router.HandleFunc("/v2/server", app.ServerSimple).
		Methods("OPTIONS", "POST").
		Name("server")

	app.router.HandleFunc("/v2/server/{address}", app.Server).
		Methods("OPTIONS", "GET", "POST").
		Name("server")

	app.router.HandleFunc("/v2/servers", app.Servers).
		Methods("OPTIONS", "GET").
		Name("servers")

	app.router.HandleFunc("/v2/stats", app.Statistics).
		Methods("OPTIONS", "GET").
		Name("stats")

	app.router.HandleFunc("/graphql/v1", app.GraphQL).
		Methods("OPTIONS", "GET", "POST").
		Name("stats")

	return &app
}

// Start begins listening for requests and blocks until fatal error
func (app *App) Start() {
	defer app.cancel()

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	err := http.ListenAndServe(app.config.Bind, handlers.CORS(headersOk, originsOk, methodsOk)(app.router))

	logger.Fatal("http server encountered fatal error",
		zap.Error(err))
}

// LoadAllAddresses loads all addresses from the database as a slice of strings for synchronisation
// with the QueryDaemon.
func (app *App) LoadAllAddresses() (result []string) {
	allServers := []Server{}
	err := app.collection.Find(bson.M{}).All(&allServers)
	if err != nil {
		logger.Fatal("failed to load current addresses for query daemon",
			zap.Error(err))
	}
	for i := range allServers {
		result = append(result, allServers[i].Core.Address)
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
