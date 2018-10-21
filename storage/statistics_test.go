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

	gotStatistics := types.Statistics{}

	var err error
	gotStatistics.Servers, err = mgr.GetActiveServers()
	if err != nil {
		t.Error(err)
	}
	gotStatistics.Players, err = mgr.GetActiveServers()
	if err != nil {
		t.Error(err)
	}
	gotStatistics.PlayersPerServer = float32(gotStatistics.Players) / float32(gotStatistics.Servers)

	assert.Equal(t, wantStatistics, gotStatistics)
}
