package main

import (
	"log"
	"net/http"
	"testing"
)

func TestApp_Servers(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		app  *App
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.app.Servers(tt.args.w, tt.args.r)
		})
	}
}

func TestApp_GetServers(t *testing.T) {
	servers, err := app.GetServers()
	if err != nil {
		t.Errorf("App.GetServers() error = %v", err)
	}

	log.Print(servers)
}
