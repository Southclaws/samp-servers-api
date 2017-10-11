package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApp_GetStatistics(t *testing.T) {
	wantStatistics := Statistics{
		Servers:          4,
		Players:          1002,
		PlayersPerServer: 250,
	}

	gotStatistics, err := app.GetStatistics()

	assert.NoError(t, err)
	assert.Equal(t, wantStatistics, gotStatistics)
}
