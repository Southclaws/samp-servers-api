package scraper

import (
	"context"
	"net"
	"time"

	"github.com/Southclaws/go-samp-query"
	"github.com/Southclaws/tickerpool"
	"github.com/pkg/errors"
	"golang.org/x/sync/syncmap"

	"github.com/Southclaws/samp-servers-api/types"
)

// Config contains parameters to tweak the scraper performance
type Config struct {
	QueryInterval    time.Duration      // interval between query attempts
	MaxFailed        int                // maximum number of failed query attempts before removing address
	QueryFunction    QueryFunction      // function for querying servers
	OnRequestArchive func(string)       // called to archive an address
	OnRequestRemove  func(string)       // called to remove an address
	OnRequestUpdate  func(types.Server) // called to update an address
}

// Scraper crawls through a list of server addresses and gathers information about them via the
// legacy query API, it then stores the results as standard Server objects, accessible via the API.
type Scraper struct {
	config         Config
	ctx            context.Context
	failedAttempts *syncmap.Map
	active         *tickerpool.TickerPool
	failed         *tickerpool.TickerPool
	metrics        *metrics
}

// QueryFunction represents a function capable of retreiving server information via the server API
type QueryFunction func(context.Context, string, bool) (sampquery.Server, error)

// New sets up the query daemon and starts the background processes
func New(ctx context.Context, initial []string, config Config) (daemon *Scraper, err error) {
	daemon = &Scraper{
		config:         config,
		ctx:            ctx,
		failedAttempts: &syncmap.Map{},
		metrics:        newMetricsRecorder(),
	}

	daemon.active, err = tickerpool.NewTickerPool(config.QueryInterval)
	if err != nil {
		err = errors.Wrap(err, "failed to create active ticker pool")
		return
	}

	daemon.failed, err = tickerpool.NewTickerPool(config.QueryInterval * 10) // query failed servers less often
	if err != nil {
		err = errors.Wrap(err, "failed to create failed ticker pool")
		return
	}

	for _, address := range initial {
		daemon.Add(address)
	}

	return
}

// Add will add a new address to the TickerPool and query it periodically
func (daemon *Scraper) Add(address string) {
	daemon.active.Add(address, func() {
		queryStart := time.Now()
		remove, err := daemon.query(address)
		if err != nil {
			daemon.metrics.Failures.Inc()
			if remove {
				daemon.metrics.Archives.Inc()
				daemon.config.OnRequestArchive(address)
				daemon.addFailed(address)
			}
		} else {
			daemon.metrics.Successes.Inc()
		}
		daemon.metrics.QueryTime.Observe(time.Since(queryStart).Seconds())
		daemon.metrics.Queries.Inc()
	})
}

// Remove will remove an address from the query rotation
func (daemon *Scraper) Remove(address string) {
	if daemon.active.Exists(address) {
		daemon.failedAttempts.Delete(address)
		daemon.active.Remove(address)
		daemon.metrics.Removals.Inc()

		daemon.config.OnRequestRemove(address)
	}
}

// addFailed marks a server as "inactive" and queries it less often
func (daemon *Scraper) addFailed(address string) {
	daemon.failedAttempts.Delete(address)
	daemon.active.Remove(address)

	daemon.failed.Add(address, func() {
		remove, err := daemon.query(address)
		if err != nil {
			if remove {
				daemon.Remove(address)
			}
		}
	})
}

// removeFailed is called when a server is "revived" so it can be added back to the regular rotation
func (daemon *Scraper) removeFailed(address string) {
	if daemon.active.Exists(address) {
		daemon.failedAttempts.Delete(address)
		daemon.failed.Remove(address)
		daemon.Add(address)
	}
}

func (daemon *Scraper) query(address string) (remove bool, err error) {
	tmp, hasFailed := daemon.failedAttempts.Load(address)
	attempts, _ := tmp.(int)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	serverData, err := daemon.config.QueryFunction(ctx, address, true)
	if err != nil {
		if hasFailed {
			if attempts > daemon.config.MaxFailed {
				return true, err
			}
			daemon.failedAttempts.Store(address, attempts+1)
			return false, err
		}
		daemon.failedAttempts.Store(address, 1)
		return false, err
	}

	if hasFailed {
		daemon.failedAttempts.Delete(address)
	}
	daemon.removeFailed(address)

	var ip string
	addrs, err := net.LookupHost(serverData.Address)
	if err != nil {
		err = nil
	}
	if len(addrs) > 0 {
		ip = addrs[0]
	}
	if ip == "" {
		ip = serverData.Address
	}

	server := types.Server{
		IP: ip,
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
	daemon.config.OnRequestUpdate(server)

	return false, nil
}
