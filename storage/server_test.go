package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/samp-servers-api/types"
)

func TestManager_UpsertServer(t *testing.T) {
	type args struct {
		server types.Server
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{types.Server{
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
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mgr.UpsertServer(tt.args.server); (err != nil) != tt.wantErr {
				t.Errorf("App.UpsertServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_GetServer(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name       string
		args       args
		wantServer types.Server
		wantFound  bool
		wantErr    bool
	}{
		{"valid", args{"ss.southcla.ws"},
			types.Server{
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
			},
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotServer, gotFound, err := mgr.GetServer(tt.args.address)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantServer, gotServer)
			assert.Equal(t, tt.wantFound, gotFound)
		})
	}
}
