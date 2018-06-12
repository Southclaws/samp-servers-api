package server

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

	"github.com/Southclaws/samp-servers-api/types"
)

// Config stores app global configuration
type Config struct {
	Bind            string `split_words:"true" required:"true"`
	MongoHost       string `split_words:"true" required:"true"`
	MongoPort       string `split_words:"true" required:"true"`
	MongoName       string `split_words:"true" required:"true"`
	MongoUser       string `split_words:"true" required:"true"`
	MongoPass       string `split_words:"true" required:"false"`
	MongoCollection string `split_words:"true" required:"true"`
	QueryInterval   int    `split_words:"true" required:"true"`
	MaxFailedQuery  int    `split_words:"true" required:"true"`
	VerifyByHost    bool   `split_words:"true" required:"true"`
}

// App stores global state for routing
type App struct {
	ctx        context.Context
	cancel     context.CancelFunc
	config     Config
	collection *mgo.Collection
	qd         *QueryDaemon
	handlers   map[string][]Route
	httpServer *http.Server
}

// Initialise sets up a database connection, binds all the routes and prepares for Start
func Initialise(config Config) (app *App, err error) {
	logger.Debug("initialising samp-servers-api with debug logging", zap.Any("config", config))

	app = &App{
		config: config,
	}
	app.ctx, app.cancel = context.WithCancel(context.Background())

	// Connect to the database, receive a collection pointer
	// TODO: return error
	app.collection = ConnectDB(config)
	logger.Info("connected to mongodb server")

	// Grab existing addresses from database and pass to the Query Daemon
	// TODO: return errors
	addresses := app.LoadAllAddresses()
	app.qd = NewQueryDaemon(app.ctx, app, addresses, time.Second*time.Duration(config.QueryInterval), config.MaxFailedQuery, sampquery.GetServerInfo)

	// Start a periodic query against the SA:MP official internet list (if it's even online...)
	// TODO: errors?
	app.LegacyListQuery()

	// TODO: split off to versioned packages
	app.handlers = map[string][]Route{
		"v2": {
			{
				Name:    "serverAdd",
				Path:    "/v2/server",
				Method:  "POST",
				handler: app.serverAdd,
			},
			{
				Name:    "serverPost",
				Path:    "/v2/server/{address}",
				Method:  "POST",
				handler: app.serverPost,
			},
			{
				Name:    "serverGet",
				Path:    "/v2/server/{address}",
				Method:  "GET",
				handler: app.serverGet,
			},
			{
				Name:    "serverList",
				Path:    "/v2/servers",
				Method:  "GET",
				handler: app.serverList,
			},
			{
				Name:    "serverStats",
				Path:    "/v2/stats",
				Method:  "GET",
				handler: app.serverStats,
			},
		},
	}

	router := mux.NewRouter().StrictSlash(true)
	for name, routes := range app.handlers {
		logger.Debug("loaded handler",
			zap.String("name", name),
			zap.Int("routes", len(routes)))

		for _, route := range routes {
			router.Methods(route.Method).
				Path(route.Path).
				Name(route.Name).
				Handler(EndpointHandler(route.handler))

			logger.Debug("registered handler route",
				zap.String("name", route.Name),
				zap.String("method", route.Method),
				zap.String("path", route.Path))
		}
	}

	app.httpServer = &http.Server{
		Addr: app.config.Bind,
		Handler: handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"HEAD", "GET", "POST", "PUT", "OPTIONS"}),
		)(router),
	}

	return app, nil
}

// Start begins listening for requests and blocks until fatal error
func (app *App) Start() error {
	defer app.cancel()
	return app.httpServer.ListenAndServe()
}

// LoadAllAddresses loads all addresses from the database as a slice of strings for synchronisation
// with the QueryDaemon.
func (app *App) LoadAllAddresses() (result []string) {
	allServers := []types.Server{}
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
