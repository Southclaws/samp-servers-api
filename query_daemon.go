package main

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

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

	go qd.Daemon()

	return &qd
}

// Add will add a new address to the query rotation
func (qd *QueryDaemon) Add(address string) {
	qd.InputAdd <- address
}

// Remove will remove an address from the query rotation
func (qd *QueryDaemon) Remove(address string) {
	qd.InputDel <- address
}

// GetTotal returns the total in a threadsafe way
func (qd *QueryDaemon) GetTotal() int32 {
	return atomic.LoadInt32(&qd.Total)
}

func (qd *QueryDaemon) add(address string) {
	_, exists := qd.Lookup[address]
	if exists {
		return
	}

	index := len(qd.ToQuery)              // first, grab the index of the back of the list
	if qd.Next > -1 && qd.Next != index { // if Next is valid and doesn't point to the back
		index = qd.Next // then use Next as the next insersion index
	}

	qd.ToQuery = append(qd.ToQuery, address)
	qd.Lookup[address] = index
	atomic.AddInt32(&qd.Total, 1)
	logger.Debug("added address to query daemon",
		zap.String("address", address),
		zap.Int("index", index))
}

func (qd *QueryDaemon) remove(address string) {
	index, exists := qd.Lookup[address]
	if !exists {
		return
	}

	if qd.Next == -1 || qd.Next > index { // if Next is valid and index is below Next
		qd.Next = index // then shift Next down to index so the next insersion goes here
	}

	delete(qd.Lookup, address)
	qd.ToQuery[index] = ""
	atomic.AddInt32(&qd.Total, -1)
	logger.Debug("removed address from query daemon",
		zap.String("address", address),
		zap.Int("next", qd.Next),
		zap.Int("index", index))
}

// Daemon runs in the background and periodically queries servers in the list round-robin style
func (qd *QueryDaemon) Daemon() {
	tick := time.NewTicker(time.Millisecond * 1000)
	logger.Debug("starting query daemon background process")
main:
	for {
		select {
		case <-qd.ctx.Done():
			break main

		// doing the add/remove inside the for-select keeps everything in sync
		case address := <-qd.InputAdd:
			qd.add(address)

		case address := <-qd.InputDel:
			qd.remove(address)

		case <-tick.C:
			if qd.Total == 0 {
				continue
			}

			if qd.ToQuery[qd.Index] != "" {
				logger.Debug("performing periodic query", zap.String("address", qd.ToQuery[qd.Index]), zap.Int32("index", qd.Index))
				go qd.query(qd.ToQuery[qd.Index])
			}

			qd.Index++
			if qd.Index >= qd.Total {
				qd.Index = 0
			}

		case result := <-qd.Finished:
			if result.Error != nil {
				logger.Debug("QueryDaemon failed to query address, removing from pool",
					zap.String("address", result.Address),
					zap.Error(result.Error))
				qd.Remove(result.Address)
			} else {
				err := qd.app.UpsertServer(result.Server)
				if err != nil {
					logger.Warn("QueryDaemon failed to upsert",
						zap.Error(err))
				}
			}
		}
	}
}

func (qd *QueryDaemon) query(address string) {
	result := ServerWrapper{
		Address: address,
	}

	server, err := GetServerLegacyInfo(address)
	if err != nil {
		result.Error = err
	}
	result.Server = server

	qd.Finished <- result
}
