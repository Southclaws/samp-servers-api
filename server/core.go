package server

import (
	"context"
	"net/http"
	"time"

	"github.com/Southclaws/go-samp-query"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/Southclaws/samp-servers-api/storage"
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
	db         *storage.Manager
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

	app.db, err = storage.New(storage.Config{
		MongoHost:       config.MongoHost,
		MongoPort:       config.MongoPort,
		MongoName:       config.MongoName,
		MongoUser:       config.MongoUser,
		MongoPass:       config.MongoPass,
		MongoCollection: config.MongoCollection,
	})
	if err != nil {
		return
	}

	// Grab existing addresses from database and pass to the Query Daemon
	addresses, err := app.db.LoadAllAddresses()
	if err != nil {
		return
	}

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
