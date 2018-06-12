package main

import (
	// loads environment variables from .env
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"

	"github.com/Southclaws/samp-servers-api/server"
)

func main() {
	config := server.Config{}
	err := envconfig.Process("SAMPLIST", &config)
	if err != nil {
		panic(err)
	}

	app, err := server.Initialise(config)
	if err != nil {
		panic(err)
	}

	err = app.Start()
	panic(err)
}
