package server

import (
	"context"
	"time"

	sampquery "github.com/Southclaws/go-samp-query"
	"github.com/Southclaws/tickerpool"
	"go.uber.org/zap"
	"golang.org/x/sync/syncmap"

	"github.com/Southclaws/samp-servers-api/types"
)

// QueryDaemon crawls through a list of server addresses and gathers information about them via the
// legacy query API, it then stores the results as standard Server objects, accessible via the API.
type QueryDaemon struct {
	QueryInterval time.Duration // interval between query attempts
	MaxFailed     int           // maximum number of failed query attempts before removing address
	Function      QueryFunction // function for querying servers

	ctx            context.Context
	app            *App
	failedAttempts *syncmap.Map
	active         *tickerpool.TickerPool
	failed         *tickerpool.TickerPool
}

// QueryFunction represents a function capable of retreiving server information via the server API
type QueryFunction func(context.Context, string, bool) (sampquery.Server, error)

// NewQueryDaemon sets up the query daemon and starts the background process
func NewQueryDaemon(ctx context.Context, app *App, initial []string, interval time.Duration, maxFailed int, queryFunction QueryFunction) *QueryDaemon {
	qd := QueryDaemon{
		QueryInterval:  interval,
		MaxFailed:      maxFailed,
		Function:       queryFunction,
		ctx:            ctx,
		app:            app,
		failedAttempts: &syncmap.Map{},
	}

	var err error
	qd.active, err = tickerpool.NewTickerPool(interval)
	if err != nil {
		logger.Fatal("failed to create active ticker pool",
			zap.Error(err))
	}

	qd.failed, err = tickerpool.NewTickerPool(interval * 10) // query failed servers less often
	if err != nil {
		logger.Fatal("failed to create failed ticker pool",
			zap.Error(err))
	}

	for _, address := range initial {
		qd.Add(address)
	}

	return &qd
}

// Add will add a new address to the TickerPool and query it every
func (qd *QueryDaemon) Add(address string) {
	qd.active.Add(address, func() { qd.add(address) })
}

func (qd *QueryDaemon) add(address string) {
	remove, err := qd.query(address)
	if err != nil {
		if remove {
			err = qd.app.db.MarkInactive(address)
			if err != nil {
				logger.Error("failed to mark address as inactive",
					zap.String("address", address),
					zap.Error(err))
			}
			qd.addFailed(address)

			logger.Debug("failed query too many times",
				zap.String("address", address),
				zap.Error(err))
		} else {
			logger.Debug("failed query",
				zap.String("address", address),
				zap.Error(err))
		}
	}
}

// Remove will remove an address from the query rotation
func (qd *QueryDaemon) Remove(address string) {
	if qd.active.Exists(address) {
		qd.failedAttempts.Delete(address)
		qd.active.Remove(address)

		err := qd.app.db.RemoveServer(address)
		if err != nil {
			logger.Warn("failed to remove server",
				zap.String("address", address),
				zap.Error(err))
		}
	}
}

// addFailed marks a server as "inactive" and queries it less often
func (qd *QueryDaemon) addFailed(address string) {
	qd.failedAttempts.Delete(address)
	qd.active.Remove(address)

	qd.failed.Add(address, func() {
		remove, err := qd.query(address)
		if err != nil {
			if remove {
				qd.Remove(address)
				logger.Debug("failed revival query too many times",
					zap.String("address", address),
					zap.Error(err))
			} else {
				logger.Debug("failed revival query",
					zap.String("address", address),
					zap.Error(err))
			}
		}
	})
}

// removeFailed is called when a server is "revived" so it can be added back to the regular rotation
func (qd *QueryDaemon) removeFailed(address string) {
	if qd.active.Exists(address) {
		qd.failedAttempts.Delete(address)
		qd.failed.Remove(address)
		qd.Add(address)
	}
}

func (qd *QueryDaemon) query(address string) (remove bool, err error) {
	tmp, hasFailed := qd.failedAttempts.Load(address)
	attempts, _ := tmp.(int)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	serverData, err := qd.Function(ctx, address, true)
	if err != nil {
		if hasFailed {
			if attempts > qd.MaxFailed {
				return true, err
			}
			qd.failedAttempts.Store(address, attempts+1)
			return false, err
		}
		qd.failedAttempts.Store(address, 1)
		return false, err
	}

	if hasFailed {
		qd.failedAttempts.Delete(address)
	}
	qd.removeFailed(address)

	server := types.Server{
		Core: types.ServerCore{
			Address:    serverData.Address,
			Hostname:   serverData.Hostname,
			Players:    serverData.Players,
			MaxPlayers: serverData.MaxPlayers,
			Gamemode:   serverData.Gamemode,
			Language:   serverData.Language,
			Password:   serverData.Password,
		},
		Rules: serverData.Rules,
	}

	if server.Core.Players > server.Core.MaxPlayers {
		return true, nil
	}
	if server.Core.MaxPlayers > 1000 {
		return true, nil
	}

	version, ok := serverData.Rules["version"]
	if ok {
		server.Core.Version = version
	}

	err = qd.app.db.UpsertServer(server)
	if err != nil {
		logger.Warn("QueryDaemon failed to upsert",
			zap.Error(err))
	}

	return false, nil
}
