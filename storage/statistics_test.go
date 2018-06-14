package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/samp-servers-api/types"
)

func TestManager_GetStatistics(t *testing.T) {
	wantStatistics := types.Statistics{
		Servers:          4,
		Players:          1002,
		PlayersPerServer: 250,
	}

	gotStatistics, err := mgr.GetStatistics()

	assert.NoError(t, err)
	assert.Equal(t, wantStatistics, gotStatistics)
}
