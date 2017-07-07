package main

import (
	"testing"
)

func TestGetServerLegacyInfo(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name       string
		args       args
		wantServer Server
		wantErr    bool
	}{
		{"valid", args{"198.251.83.150:7777"}, Server{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetServerLegacyInfo(tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetServerLegacyInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
