package storage

import (
	"os"
	"testing"
)

var mgr *Manager

func TestMain(m *testing.M) {
	var err error
	mgr, err = New(Config{
		MongoHost:       "localhost",
		MongoPort:       "27017",
		MongoName:       "samplist",
		MongoUser:       "root",
		MongoPass:       "",
		MongoCollection: "servers",
	})
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
