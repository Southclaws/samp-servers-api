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
	QueryInterval time.Duration // interval between query attempts
	MaxFailed     int           // maximum number of failed query attempts before removing address

	ctx            context.Context
	app            *App
	failedAttempts map[string]int
	tp             *tickerpool.TickerPool
}

// NewQueryDaemon sets up the query daemon and starts the background process
func NewQueryDaemon(ctx context.Context, app *App, initial []string, interval time.Duration, maxFailed int) *QueryDaemon {
	qd := QueryDaemon{
		QueryInterval:  interval,
		MaxFailed:      maxFailed,
		ctx:            ctx,
		app:            app,
		failedAttempts: make(map[string]int),
	}

	var err error
	qd.tp, err = tickerpool.NewTickerPool(interval)
	if err != nil {
		logger.Fatal("failed to create new ticker pool",
			zap.Error(err))
	}

	for _, address := range initial {
		qd.Add(address)
	}

	return &qd
}

// Add will add a new address to the TickerPool and query it every
func (qd *QueryDaemon) Add(address string) {
	qd.tp.Add(address, func() {
		attempts, hasFailed := qd.failedAttempts[address]

		server, err := GetServerLegacyInfo(address)
		if err != nil {
			if hasFailed {
				if qd.failedAttempts[address] > qd.MaxFailed {
					qd.Remove(address)

					logger.Debug("failed query too many times",
						zap.String("address", address),
						zap.Error(err))
				} else {
					qd.failedAttempts[address] = attempts + 1
					logger.Debug("failed query",
						zap.String("address", address))
				}
			} else {
				qd.failedAttempts[address] = 1
			}
		} else {
			if hasFailed {
				if attempts <= 1 {
					delete(qd.failedAttempts, address)
				} else {
					qd.failedAttempts[address]--
				}
			}

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
	if qd.tp.Exists(address) {
		delete(qd.failedAttempts, address)
		qd.tp.Remove(address)

		err := qd.app.RemoveServer(address)
		if err != nil {
			logger.Warn("failed to remove server",
				zap.String("address", address),
				zap.Error(err))
		}
	}
}
