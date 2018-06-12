package storage

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

// Config describes db connection information
type Config struct {
	MongoHost       string `split_words:"true" required:"true"`
	MongoPort       string `split_words:"true" required:"true"`
	MongoName       string `split_words:"true" required:"true"`
	MongoUser       string `split_words:"true" required:"true"`
	MongoPass       string `split_words:"true" required:"false"`
	MongoCollection string `split_words:"true" required:"true"`
}

// Manager provides access to collections and predefined CRUD functionality.
type Manager struct {
	config     Config
	session    *mgo.Session
	db         *mgo.Database
	collection *mgo.Collection
}

// New sets up a MongoDB connection and ensures it is ready to use
func New(config Config) (mgr *Manager, err error) {
	mgr = &Manager{
		config: config,
	}

	mgr.session, err = mgo.Dial(fmt.Sprintf("%s:%s", config.MongoHost, config.MongoPort))
	if err != nil {
		return
	}

	if config.MongoPass != "" {
		err = mgr.session.Login(&mgo.Credential{
			Source:   config.MongoName,
			Username: config.MongoUser,
			Password: config.MongoPass,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to log in to mongodb")
		}
	}

	mgr.collection = mgr.session.DB(config.MongoName).C(config.MongoCollection)

	err = mgr.collection.EnsureIndex(mgo.Index{
		Key:         []string{"core.address"},
		Unique:      true,
		DropDups:    true,
		ExpireAfter: time.Hour * 168,
	})
	if err != nil {
		return nil, errors.Wrap(err, "index ensure failed")
	}

	return
}
