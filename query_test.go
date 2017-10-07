package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetServerLegacyInfo(t *testing.T) {

	// Travis doesn't seem to like outgoing UDP
	if os.Getenv("TRAVIS") == "1" {
		return
	}

	type args struct {
		host string
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{"valid", args{"server.ls-rp.com:7777"}, ""},
		{"valid", args{"46.174.54.184:7777"}, ""},
		{"invalid", args{"18.251.83.150:80"}, "socket read timed out"},
		{"invalid", args{"not a valid url"}, "failed to resolve: address not a valid url: missing port in address"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := GetServerLegacyInfo(tt.args.host)
			if err != nil {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NotEmpty(t, server.Core.Address)
				assert.NotEmpty(t, server.Core.Hostname)
				assert.NotEmpty(t, server.Core.Gamemode)
				assert.NotZero(t, server.Core.MaxPlayers)
			}
		})
	}
}
