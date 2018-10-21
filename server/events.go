package server

import (
	"github.com/Southclaws/samp-servers-api/types"
	"go.uber.org/zap"
)

func (app *App) onRequestArchive(address string) {
	logger.Debug("archiving server",
		zap.String("address", address))

	err := app.db.ArchiveServer(address)
	if err != nil {
		logger.Error("failed to archive server",
			zap.Error(err),
			zap.String("address", address))
		return
	}

	app.updateIndexMetrics()
}

func (app *App) onRequestRemove(address string) {
	logger.Debug("removing server",
		zap.String("address", address))

	err := app.db.RemoveServer(address)
	if err != nil {
		logger.Error("failed to remove server",
			zap.Error(err),
			zap.String("address", address))
		return
	}

	app.updateIndexMetrics()
}

func (app *App) onRequestUpdate(server types.Server) {
	logger.Debug("updating server",
		zap.String("address", server.Core.Address))

	err := app.db.UpsertServer(server)
	if err != nil {
		logger.Error("failed to upsert server",
			zap.Error(err),
			zap.String("address", server.Core.Address))
		return
	}

	app.updateIndexMetrics()
}

func (app *App) updateIndexMetrics() {
	c, err := app.db.GetActiveServers()
	if err != nil {
		logger.Error("failed to get active servers metric",
			zap.Error(err))
	}
	app.metrics.Active.Set(float64(c))

	c, err = app.db.GetInactiveServers()
	if err != nil {
		logger.Error("failed to get inactive servers metric",
			zap.Error(err))
	}
	app.metrics.Inactive.Set(float64(c))

	c, err = app.db.GetTotalPlayers()
	if err != nil {
		logger.Error("failed to get total players metric",
			zap.Error(err))
	}
	app.metrics.Players.Set(float64(c))
}
