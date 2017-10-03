package main

import (
	"os"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var version = "master"

// Config stores app global configuration
type Config struct {
	Version         string
	Bind            string
	MongoHost       string
	MongoPort       string
	MongoName       string
	MongoUser       string
	MongoPass       string
	MongoCollection string
	QueryInterval   int
	MaxFailedQuery  int
	VerifyByHost    bool
}

var logger *zap.Logger

func init() {
	var config zap.Config
	debug := os.Getenv("DEBUG")

	if os.Getenv("TESTING") != "" {
		config = zap.NewDevelopmentConfig()
		config.DisableCaller = true
	} else {
		config = zap.NewProductionConfig()
		config.EncoderConfig.MessageKey = "@message"
		config.EncoderConfig.TimeKey = "@timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		if debug != "0" && debug != "" {
			dyn := zap.NewAtomicLevel()
			dyn.SetLevel(zap.DebugLevel)
			config.Level = dyn
		}
	}
	_logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	logger = _logger.With(
		zap.String("@version", os.Getenv("GIT_HASH")),
		zap.Namespace("@fields"),
	)
}

func main() {
	config := Config{
		Version:         version,
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
	app := Initialise(config)
	app.Start()
}
func configStrFromEnv(name string) (value string) {
	value = os.Getenv(name)
	if value == "" {
		logger.Fatal("environment variable not set",
			zap.String("name", name))
	}
	return
}

func configIntFromEnv(name string) (value int) {
	valueStr := os.Getenv(name)
	if valueStr == "" {
		logger.Fatal("environment variable not set",
			zap.String("name", name))
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		logger.Fatal("failed to convert environment variable to int",
			zap.Error(err),
			zap.String("name", name))
	}
	return
}
