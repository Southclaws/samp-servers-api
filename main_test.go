package main

import (
	"fmt"
	"os"
	"testing"
)

var app *App

func TestMain(m *testing.M) {
	config := Config{
		Bind:            "localhost:8080",
		MongoHost:       "localhost",
		MongoPort:       "27017",
		MongoName:       "samplist",
		MongoUser:       "root",
		MongoPass:       "",
		MongoCollection: "servers",
		QueryInterval:   100000,
		MaxFailedQuery:  0,
		VerifyByHost:    false,
	}

	fmt.Println("initialising announce-backend testing mode", config)

	app = Initialise(config)
	go app.Start() // start the server in a goroutine

	ret := m.Run() // run the tests against the server
	os.Exit(ret)
}
