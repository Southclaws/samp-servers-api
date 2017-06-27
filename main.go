package main

import (
	"encoding/json"
	"flag"
	"os"

	"go.uber.org/zap"
)

// Config stores app global configuration
type Config struct {
	Bind      string `json:"bind"`
	MongoUser string `json:"mongodb_user"`
	MongoPass string `json:"mongodb_pass"`
	MongoHost string `json:"mongodb_host"`
	MongoPort string `json:"mongodb_port"`
	MongoName string `json:"mongodb_name"`
}

var logger *zap.Logger

func init() {
	var err error
	var config zap.Config

	dyn := zap.NewAtomicLevel()
	dyn.SetLevel(zap.DebugLevel)
	config.Level = dyn
	config = zap.NewDevelopmentConfig()
	config.DisableCaller = true

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

	Start(cfg)
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
