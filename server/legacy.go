package server

import (
	"bufio"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/Southclaws/samp-servers-api/types"
)

// LegacyListQuery periodically hits the lists.sa-mp.com endpoint to update the new servers list.
func (app *App) LegacyListQuery() {
	app.getMasterlist()
	go func() {
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			err := app.getMasterlist()
			if err != nil {
				logger.Error("failed to get lists.sa-mp.com",
					zap.Error(err))
			}
		}
	}()
}

func (app *App) getMasterlist() (err error) {
	resp, err := http.Get("http://lists.sa-mp.com/0.3.7/servers")
	if err != nil {
		return errors.Wrap(err, "failed to get masterlist")
	}

	if resp.StatusCode != 200 {
		return errors.Errorf("unexpected masterlist status %s", resp.Status)
	}

	count := 0
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		address, errs := types.ValidateAddress(scanner.Text())
		if errs != nil {
			err = errs[0]
			return
		}

		logger.Debug("adding server from legacy masterlist",
			zap.String("address", address))

		app.qd.Add(address)
		count++
	}
	logger.Debug("added servers from masterlist", zap.Int("servers", count))

	return
}
