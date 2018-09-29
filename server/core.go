package server

import (
	"context"
	"net/http"
	"path"

	"github.com/Southclaws/go-samp-query"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/Southclaws/samp-servers-api/scraper"
	"github.com/Southclaws/samp-servers-api/server/v2"
	"github.com/Southclaws/samp-servers-api/storage"
	"github.com/Southclaws/samp-servers-api/types"
)

// App stores global state for routing
type App struct {
	ctx        context.Context
	cancel     context.CancelFunc
	config     types.Config
	db         *storage.Manager
	qd         *scraper.Scraper
	handlers   map[string]types.RouteHandler
	httpServer *http.Server
}

// Initialise sets up a database connection, binds all the routes and prepares for Start
func Initialise(config types.Config) (app *App, err error) {
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

	app.qd, err = scraper.New(
		app.ctx,
		addresses,
		scraper.Config{
			QueryInterval:    config.QueryInterval,
			MaxFailed:        config.MaxFailedQuery,
			QueryFunction:    sampquery.GetServerInfo,
			OnRequestArchive: app.onRequestArchive,
			OnRequestRemove:  app.onRequestRemove,
			OnRequestUpdate:  app.onRequestUpdate,
		})
	if err != nil {
		return
	}

	if config.LegacyList {
		// Start a periodic query against the SA:MP official internet list (if it's even online...)
		go app.LegacyListQuery()
	}

	app.handlers = map[string]types.RouteHandler{
		"v2": v2.Init(app.db, app.qd, app.qd.Metrics, config),
		// "v3": v3.Init(app.db, app.qd, app.qd.Metrics, config),
	}

	router := mux.NewRouter().StrictSlash(true)
	for name, handler := range app.handlers {
		routes := handler.Routes()

		logger.Debug("loaded handler",
			zap.String("name", name),
			zap.Int("routes", len(routes)))

		for _, route := range routes {
			router.Methods(route.Method).
				Path(path.Join("/", name, route.Path)).
				Name(route.Name).
				Handler(route.Handler)

			logger.Debug("registered handler route",
				zap.String("name", route.Name),
				zap.String("method", route.Method),
				zap.String("path", path.Join(name, route.Path)))
		}

		router.Methods("GET").
			Path(path.Join("/", name, "docs")).
			Name("docs").
			Handler(app.docsWrapper(handler))
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
