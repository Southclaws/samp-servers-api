package main

import (
	"encoding/json"
	"flag"
	"os"

	"go.uber.org/zap"
)

// Config stores app global configuration
type Config struct {
	Bind           string `json:"bind"`
	MongoHost      string `json:"mongodb_host"`
	MongoPort      string `json:"mongodb_port"`
	MongoName      string `json:"mongodb_name"`
	MongoUser      string `json:"mongodb_user"`
	MongoPass      string `json:"mongodb_pass"`
	QueryInterval  int    `json:"query_interval"`
	MaxFailedQuery int    `json:"max_failed_query"`
}

var logger *zap.Logger

func init() {
	var err error
	var config zap.Config

	config = zap.NewProductionConfig()

	// dyn := zap.NewAtomicLevel()
	// dyn.SetLevel(zap.DebugLevel)
	// config.Level = dyn

	logger, err = config.Build()
	if err != nil {
		panic(err)
	}
}

func main() {
	configFile := flag.String("config", "config.json", "path to config.json file")
	flag.Parse()

	cfg := loadConfig(*configFile)

	logger.Info("initialising announce-backend", zap.Any("config", cfg))

	app := Initialise(cfg)
	app.Start()
}

func loadConfig(filename string) Config {
	var (
		err    error
		file   *os.File
		config Config
	)

	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		file, err = os.Create(filename)
		if err != nil {
			logger.Fatal("failed to create default config.json",
				zap.Error(err))
		}

		config.Bind = "localhost:7790"
		config.MongoHost = "localhost"
		config.MongoPort = "27017"
		config.MongoName = "samplist"
		config.MongoUser = "samplist"
		config.MongoPass = "changeme"
		config.QueryInterval = 60
		config.MaxFailedQuery = 10

		enc := json.NewEncoder(file)
		enc.SetIndent("", "    ")
		enc.Encode(&config)
		if err != nil {
			logger.Fatal("failed to encode default config.json",
				zap.Error(err))
		}
		return config
	}

	file, err = os.Open(filename)
	if err != nil {
		logger.Fatal("failed to open config file",
			zap.Error(err))
	}

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		logger.Fatal("failed to decode config file",
			zap.Error(err))
	}

	err = file.Close()
	if err != nil {
		logger.Fatal("failed to close config file",
			zap.Error(err))
	}

	return config
}
