package main

import (
	"fmt"
	"os"
	"testing"
)

var app *App

func TestMain(m *testing.M) {
	config := Config{
		Version:         "testing",
		Bind:            configStrFromEnv("BIND"),
		MongoHost:       configStrFromEnv("MONGO_HOST"),
		MongoPort:       configStrFromEnv("MONGO_PORT"),
		MongoName:       configStrFromEnv("MONGO_NAME"),
		MongoUser:       configStrFromEnv("MONGO_USER"),
		MongoPass:       configStrFromEnv("MONGO_PASS"),
		MongoCollection: configStrFromEnv("MONGO_COLLECTION"),
		QueryInterval:   configIntFromEnv("QUERY_INTERVAL"),
		MaxFailedQuery:  configIntFromEnv("MAX_FAILED_QUERY"),
		VerifyByHost:    configIntFromEnv("VERIFY_BY_HOST") == 1,
	}

	fmt.Println("initialising announce-backend testing mode", config)

	app = Initialise(config)
	go app.Start() // start the server in a goroutine

	ret := m.Run() // run the tests against the server
	os.Exit(ret)
}
