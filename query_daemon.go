package main

import (
	"context"
	"time"

	"github.com/Southclaws/tickerpool"
	"go.uber.org/zap"
)

// QueryDaemon crawls through a list of server addresses and gathers information about them via the
// legacy query API, it then stores the results as standard Server objects, accessible via the API.
type QueryDaemon struct {
	ctx context.Context
	app *App
	tp  *tickerpool.TickerPool
}

// ServerWrapper wraps the Server object to add an error field for reporting errors back to the
// Daemon so it can remove the errored address.
type ServerWrapper struct {
	Error   error
	Address string
	Server  Server
}

// NewQueryDaemon sets up the query daemon and starts the background process
func NewQueryDaemon(ctx context.Context, app *App) *QueryDaemon {
	qd := QueryDaemon{
		ctx: ctx,
		app: app,
	}

	var err error
	qd.tp, err = tickerpool.NewTickerPool(time.Second * 5)
	if err != nil {
		logger.Fatal("failed to create new ticker pool",
			zap.Error(err))
	}

	return &qd
}

// Add will add a new address to the query rotation
func (qd *QueryDaemon) Add(address string) {
	qd.tp.Add(address, func() {
		server, err := GetServerLegacyInfo(address)
		if err != nil {
			logger.Debug("QueryDaemon failed to query address, removing from pool",
				zap.String("address", address),
				zap.Error(err))
			qd.Remove(address)
		} else {
			err = qd.app.UpsertServer(server)
			if err != nil {
				logger.Warn("QueryDaemon failed to upsert",
					zap.Error(err))
			}
		}
	})
}

// Remove will remove an address from the query rotation
func (qd *QueryDaemon) Remove(address string) {
	qd.tp.Remove(address)
}
