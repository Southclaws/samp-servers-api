package main

import (
	"context"
	"sync"
	"time"

	"github.com/Southclaws/tickerpool"
	"go.uber.org/zap"
)

// QueryDaemon crawls through a list of server addresses and gathers information about them via the
// legacy query API, it then stores the results as standard Server objects, accessible via the API.
type QueryDaemon struct {
	ctx      context.Context
	app      *App
	InputAdd chan string        // input channel for addresses to add to the peroidic query
	InputDel chan string        // input channel for addresses to remove from the periodic query
	Finished chan ServerWrapper // successfully queried servers get sent down here
	ToQuery  []string           // list of addresses to query periodically in a round-robin fashion
	Lookup   map[string]int     // maps from address to ToQuery index
	Next     int                // the next available index
	Total    int32              // total amount of addresses, because len(ToQuery) won't work
	Index    int32              // rotates through the ToQuery list of addresses on each Daemon tick
	Mutex    sync.Mutex         // locks when a new address is added or when one is being queried
	tp       *tickerpool.TickerPool
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
		ctx:      ctx,
		app:      app,
		InputAdd: make(chan string),
		InputDel: make(chan string),
		Finished: make(chan ServerWrapper),
		Lookup:   make(map[string]int),
		Next:     -1,
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
