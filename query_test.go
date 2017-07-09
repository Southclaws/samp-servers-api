package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetServerLegacyInfo(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"198.251.83.150:7777"}, false},
		{"invalid", args{"18.251.83.150:80"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := GetServerLegacyInfo(tt.args.host)
			assert.NoError(t, err)
			assert.NotEmpty(t, server.Core.Address)
			assert.NotEmpty(t, server.Core.Hostname)
			assert.NotEmpty(t, server.Core.Gamemode)
			assert.NotZero(t, server.Core.MaxPlayers)
		})
	}
}
