package server

import (
	"testing"

	"github.com/Southclaws/samp-servers-api/types"
	"github.com/stretchr/testify/assert"
)

func TestApp_GetStatistics(t *testing.T) {
	wantStatistics := types.Statistics{
		Servers:          4,
		Players:          1002,
		PlayersPerServer: 250,
	}

	gotStatistics, err := app.db.GetStatistics()

	assert.NoError(t, err)
	assert.Equal(t, wantStatistics, gotStatistics)
}
