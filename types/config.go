package types

import "time"

// Config stores app global configuration
type Config struct {
	Version         string
	Bind            string        `split_words:"true" required:"true"`
	MongoHost       string        `split_words:"true" required:"true"`
	MongoPort       string        `split_words:"true" required:"true"`
	MongoName       string        `split_words:"true" required:"true"`
	MongoUser       string        `split_words:"true" required:"true"`
	MongoPass       string        `split_words:"true" required:"false"`
	MongoCollection string        `split_words:"true" required:"true"`
	QueryInterval   time.Duration `split_words:"true" required:"true"`
	MaxFailedQuery  int           `split_words:"true" required:"true"`
	VerifyByHost    bool          `split_words:"true" required:"true"`
	LegacyList      bool          `split_words:"true" required:"true"`
}
