package main

import (
	// loads environment variables from .env
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"

	"github.com/Southclaws/samp-servers-api/server"
	"github.com/Southclaws/samp-servers-api/types"
)

var version = "master"

func main() {
	config := types.Config{}
	err := envconfig.Process("SAMPLIST", &config)
	if err != nil {
		panic(err)
	}

	config.Version = version

	app, err := server.Initialise(config)
	if err != nil {
		panic(err)
	}

	err = app.Start()
	panic(err)
}
