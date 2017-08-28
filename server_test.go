package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/resty.v0"
)

func TestServer_Validate(t *testing.T) {
	tests := []struct {
		name     string
		server   *Server
		wantErrs []error
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrs := tt.server.Validate(); !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("Server.Validate() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}

func TestValidateAddress(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name     string
		args     args
		wanrAddr string
		wantErrs []string
	}{
		{"valid", args{"192.168.1.2"}, "samp://192.168.1.2:7777", nil},
		{"valid.port", args{"192.168.1.2:7777"}, "samp://192.168.1.2:7777", nil},
		{"valid.scheme", args{"samp://192.168.1.2"}, "samp://192.168.1.2:7777", nil},
		{"invalid.empty", args{""}, "", []string{"address is empty"}},
		{"invalid.port", args{"192.168.1.2:port"}, "", []string{"invalid port 'port' specified"}},
		{"invalid.scheme", args{"http://192.168.1.2"}, "", []string{"address contains invalid scheme 'http', must be either empty or 'samp://'"}},
		{"invalid.user", args{"user:pass@192.168.1.2"}, "", []string{"address contains a user:password component"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErrs := ValidateAddress(tt.args.address)

			for i := range gotErrs {
				assert.Equal(t, errors.Cause(gotErrs[i]).Error(), tt.wantErrs[i])
			}
		})
	}
}

func TestApp_ServerSimple(t *testing.T) {
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
			resp, err := resty.SetDebug(true).R().SetBody(tt.args.address).Post("http://localhost:7790/v2/server")
			if err != nil {
				t.Errorf("/server POST failed: %v", err)
			}
			if resp.StatusCode() != 200 {
				t.Errorf("/server POST non-200: %s", resp.Status())
			}
		})
	}

	time.Sleep(time.Second)
}

func TestApp_ServerPOST(t *testing.T) {
	type args struct {
		address string
		server  Server
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"valid 1",
			args{"ss.southcla.ws", Server{
				Core: ServerCore{
					Address:    "ss.southcla.ws",
					Hostname:   "Scavenge and Survive Official",
					Players:    4,
					MaxPlayers: 32,
					Gamemode:   "Scavenge & Survive by Southclaws",
					Language:   "English",
					Password:   false,
				},
				Rules:       map[string]string{"mapname": "San Androcalypse"},
				PlayerList:  []string{"Southclaws", "Dogmeat", "Avariam", "VIRUXE"},
				Description: "Scavenge and Survive is a very fun server!",
				Banner:      "https://i.imgur.com/o13jh8h",
			}},
		},
		{
			"valid 2",
			args{"s2.example.com", Server{
				Core: ServerCore{
					Address:    "s2.example.com",
					Hostname:   "test server 2",
					Players:    0,
					MaxPlayers: 100,
					Gamemode:   "Grand Larceny",
					Language:   "English",
					Password:   false,
				},
				Rules:       map[string]string{"mapname": "Los Santos"},
				PlayerList:  []string{},
				Description: "Test gamemode!",
			}},
		},
		{
			"valid 3",
			args{"s3.example.com", Server{
				Core: ServerCore{
					Address:    "s3.example.com",
					Hostname:   "test server 3",
					Players:    948,
					MaxPlayers: 1000,
					Gamemode:   "Grand Larceny",
					Language:   "English",
					Password:   false,
				},
				Rules:       map[string]string{"mapname": "San Fierro"},
				PlayerList:  []string{},
				Description: "Best gamemode!",
			}},
		},
		{
			"valid 4",
			args{"s4.example.com", Server{
				Core: ServerCore{
					Address:    "s4.example.com",
					Hostname:   "test server 4",
					Players:    50,
					MaxPlayers: 50,
					Gamemode:   "rivershell",
					Language:   "Polish",
					Password:   true,
				},
				Rules:       map[string]string{"mapname": "rivershell"},
				PlayerList:  []string{},
				Description: "Rivershell 4 ever",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := resty.SetDebug(true).R().SetBody(tt.args.server).Post(fmt.Sprintf("http://localhost:7790/v2/server/%s", tt.args.address))
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

func TestApp_ServerGET(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name       string
		args       args
		wantServer Server
	}{
		{"valid", args{"ss.southcla.ws"}, Server{
			Core: ServerCore{
				Address:    "ss.southcla.ws",
				Hostname:   "Scavenge and Survive Official",
				Players:    4,
				MaxPlayers: 32,
				Gamemode:   "Scavenge & Survive by Southclaws",
				Language:   "English",
				Password:   false,
			},
			Rules:       map[string]string{"mapname": "San Androcalypse"},
			PlayerList:  []string{"Southclaws", "Dogmeat", "Avariam", "VIRUXE"},
			Description: "Scavenge and Survive is a very fun server!",
			Banner:      "https://i.imgur.com/o13jh8h",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotServer := Server{}
			resp, err := resty.SetDebug(true).R().SetResult(&gotServer).Get(fmt.Sprintf("http://localhost:7790/v2/server/%s", tt.args.address))
			assert.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode())
			assert.Equal(t, tt.wantServer, gotServer)
		})
	}
}

func TestApp_GetServer(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name       string
		args       args
		wantServer Server
		wantFound  bool
		wantErr    bool
	}{
		{"valid", args{"ss.southcla.ws"},
			Server{
				Core: ServerCore{
					Address:    "ss.southcla.ws",
					Hostname:   "Scavenge and Survive Official",
					Players:    4,
					MaxPlayers: 32,
					Gamemode:   "Scavenge & Survive by Southclaws",
					Language:   "English",
					Password:   false,
				},
				Rules:       map[string]string{"mapname": "San Androcalypse"},
				PlayerList:  []string{"Southclaws", "Dogmeat", "Avariam", "VIRUXE"},
				Description: "Scavenge and Survive is a very fun server!",
				Banner:      "https://i.imgur.com/o13jh8h",
			},
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotServer, gotFound, err := app.GetServer(tt.args.address)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantServer, gotServer)
			assert.Equal(t, tt.wantFound, gotFound)
		})
	}
}

func TestApp_UpsertServer(t *testing.T) {
	type args struct {
		server Server
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{Server{
			Core: ServerCore{
				Address:    "ss.southcla.ws",
				Hostname:   "Scavenge and Survive Official",
				Players:    4,
				MaxPlayers: 32,
				Gamemode:   "Scavenge & Survive by Southclaws",
				Language:   "English",
				Password:   false,
			},
			Rules:       map[string]string{"mapname": "San Androcalypse"},
			PlayerList:  []string{"Southclaws", "Dogmeat", "Avariam", "VIRUXE"},
			Description: "Scavenge and Survive is a very fun server!",
			Banner:      "https://i.imgur.com/o13jh8h",
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := app.UpsertServer(tt.args.server); (err != nil) != tt.wantErr {
				t.Errorf("App.UpsertServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
