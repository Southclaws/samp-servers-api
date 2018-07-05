package server

import (
	"go.uber.org/zap"

	"github.com/Southclaws/samp-servers-api/types"
)

func (app *App) onRequestArchive(address string) {
	logger.Debug("archiving server",
		zap.String("address", address))

	errInner := app.db.ArchiveServer(address)
	if errInner != nil {
		logger.Error("failed to archive server",
			zap.Error(errInner),
			zap.String("address", address))
		return
	}
}

func (app *App) onRequestRemove(address string) {
	logger.Debug("removing server",
		zap.String("address", address))

	errInner := app.db.RemoveServer(address)
	if errInner != nil {
		logger.Error("failed to remove server",
			zap.Error(errInner),
			zap.String("address", address))
		return
	}
}

func (app *App) onRequestUpdate(server types.Server) {
	logger.Debug("updating server",
		zap.String("address", server.Core.Address))

	errInner := app.db.UpsertServer(server)
	if errInner != nil {
		logger.Error("failed to upsert server",
			zap.Error(errInner),
			zap.String("address", server.Core.Address))
		return
	}
}
