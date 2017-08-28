package main

import (
	"fmt"
	"os"
	"testing"
)

var app *App

func TestMain(m *testing.M) {
	cfg := loadConfig(os.Getenv("TEST_CONFIG"))

	fmt.Println("initialising announce-backend testing mode", cfg)

	app = Initialise(cfg)
	go app.Start() // start the server in a goroutine

	ret := m.Run() // run the tests against the server
	os.Exit(ret)
}
