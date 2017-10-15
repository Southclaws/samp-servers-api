package main

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	mgo "gopkg.in/mgo.v2"
)

// ConnectDB simply provides a function to set up a MongoDB connection and perform some checks
// against the selected database/collection to ensure it's ready for use.
func ConnectDB(config Config) (collection *mgo.Collection) {
	session, err := mgo.Dial(fmt.Sprintf("%s:%s", config.MongoHost, config.MongoPort))
	if err != nil {
		logger.Fatal("failed to connect to mongodb",
			zap.Error(err))
	}

	if config.MongoPass != "" {
		err = session.Login(&mgo.Credential{
			Source:   config.MongoName,
			Username: config.MongoUser,
			Password: config.MongoPass,
		})
		if err != nil {
			logger.Fatal("failed to log in to mongodb",
				zap.Error(err))
		}
		logger.Info("logged in to mongodb server")
	}

	if !CollectionExists(session, config.MongoName, config.MongoCollection) {
		err = session.DB(config.MongoName).C(config.MongoCollection).Create(&mgo.CollectionInfo{})
		if err != nil {
			logger.Fatal("collection create failed",
				zap.String("collection", config.MongoCollection),
				zap.Error(err))
		}
	}

	collection = session.DB(config.MongoName).C(config.MongoCollection)

	err = collection.EnsureIndex(mgo.Index{
		Key:         []string{"core.address"},
		Unique:      true,
		DropDups:    true,
		ExpireAfter: time.Hour * 168,
	})
	if err != nil {
		logger.Fatal("index ensure failed",
			zap.Error(err))
	}

	return
}

// CollectionExists checks if a collection exists in MongoDB
func CollectionExists(session *mgo.Session, db, wantCollection string) bool {
	collections, err := session.DB(db).CollectionNames()
	if err != nil {
		logger.Fatal("failed to get collection names",
			zap.Error(err))
	}

	for _, collection := range collections {
		if collection == wantCollection {
			return true
		}
	}

	return false
}
