package main

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	cfg := loadConfig("config_test.json")

	logger.Info("initialising announce-backend", zap.Any("config", cfg))

	go Start(cfg)  // start the server in a goroutine
	ret := m.Run() // run the tests against the server
	os.Exit(ret)
}
