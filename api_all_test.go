package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Southclaws/samp-servers-api/server"
	"github.com/Southclaws/samp-servers-api/types"
)

var app *server.App

func TestMain(m *testing.M) {
	config := types.Config{
		Bind:            "localhost:8080",
		MongoHost:       "localhost",
		MongoPort:       "27017",
		MongoName:       "samplist",
		MongoUser:       "root",
		MongoPass:       "",
		MongoCollection: "servers",
		QueryInterval:   time.Hour, // don't query during tests
		MaxFailedQuery:  0,
		VerifyByHost:    false,
	}

	fmt.Println("initialising announce-backend testing mode", config)

	var err error
	app, err = server.Initialise(config)
	if err != nil {
		panic(err)
	}
	go app.Start() // start the server in a goroutine

	ret := m.Run() // run the tests against the server
	os.Exit(ret)
}
