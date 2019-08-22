package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/resty.v1"

	"github.com/Southclaws/samp-servers-api/types"
)

func TestAPI_ServerPostAddress(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
	}{
		{"valid", args{"93.119.25.177:7777"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := resty.
				SetDebug(false).
				SetRedirectPolicy(resty.FlexibleRedirectPolicy(2)).
				R().
				SetFormData(map[string]string{"address": tt.args.address}).
				Post("http://localhost:8080/v2/server")
			if err != nil {
				t.Errorf("/server PATCH failed: %v", err)
			}
			if resp.StatusCode() != 200 {
				t.Errorf("/server PATCH non-200: %s", resp.Status())
			}
		})
	}

	time.Sleep(time.Second)
}

func TestAPI_ServerPost(t *testing.T) {
	type args struct {
		address string
		server  types.Server
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"valid 1",
			args{"ss.southcla.ws", types.Server{
				Core: types.ServerCore{
					Address:    "ss.southcla.ws",
					Hostname:   "Scavenge and Survive Official",
					Players:    4,
					MaxPlayers: 32,
					Gamemode:   "Scavenge & Survive by Southclaws",
					Language:   "English",
					Password:   false,
				},
				Rules:       map[string]string{"mapname": "San Androcalypse"},
				Description: "Scavenge and Survive is a very fun server!",
				Banner:      "https://i.imgur.com/o13jh8h",
			}},
		},
		{
			"valid 2",
			args{"s2.example.com", types.Server{
				Core: types.ServerCore{
					Address:    "s2.example.com",
					Hostname:   "test server 2",
					Players:    0,
					MaxPlayers: 100,
					Gamemode:   "Grand Larceny",
					Language:   "English",
					Password:   false,
				},
				Rules:       map[string]string{"mapname": "Los Santos"},
				Description: "Test gamemode!",
			}},
		},
		{
			"valid 3",
			args{"s3.example.com", types.Server{
				Core: types.ServerCore{
					Address:    "s3.example.com",
					Hostname:   "test server 3",
					Players:    948,
					MaxPlayers: 1000,
					Gamemode:   "Grand Larceny",
					Language:   "English",
					Password:   false,
				},
				Rules:       map[string]string{"mapname": "San Fierro"},
				Description: "Best gamemode!",
			}},
		},
		{
			"valid 4",
			args{"s4.example.com", types.Server{
				Core: types.ServerCore{
					Address:    "s4.example.com",
					Hostname:   "test server 4",
					Players:    50,
					MaxPlayers: 50,
					Gamemode:   "rivershell",
					Language:   "Polish",
					Password:   true,
				},
				Rules:       map[string]string{"mapname": "rivershell"},
				Description: "Rivershell 4 ever",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := resty.SetDebug(false).R().SetBody(tt.args.server).Patch(fmt.Sprintf("http://localhost:8080/v2/server"))
			if err != nil {
				t.Errorf("/server POST failed: %v", err)
			}
			if resp.StatusCode() != 200 {
				t.Errorf("/server POST non-200: %s, %s", resp.Status(), string(resp.Body()))
			}
		})
	}

	time.Sleep(time.Second)
}

func TestAPI_ServerGet(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name       string
		args       args
		wantServer types.Server
	}{
		{"valid", args{"ss.southcla.ws"}, types.Server{
			Core: types.ServerCore{
				Address:    "ss.southcla.ws",
				Hostname:   "Scavenge and Survive Official",
				Players:    4,
				MaxPlayers: 32,
				Gamemode:   "Scavenge & Survive by Southclaws",
				Language:   "English",
				Password:   false,
			},
			Rules:       map[string]string{"mapname": "San Androcalypse"},
			Description: "Scavenge and Survive is a very fun server!",
			Banner:      "https://i.imgur.com/o13jh8h",
			Active:      true,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotServer := types.Server{}
			resp, err := resty.SetDebug(false).R().SetResult(&gotServer).Get(fmt.Sprintf("http://localhost:8080/v2/server/%s", tt.args.address))
			assert.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode())
			assert.Equal(t, tt.wantServer, gotServer)
		})
	}
}
