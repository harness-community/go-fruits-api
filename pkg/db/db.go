package db

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Config configures the database to initialize
type Config struct {
	Ctx          context.Context
	Log          *logrus.Logger
	URI          string
	Client       *mongo.Client
	DB           *mongo.Database
	databaseName string
	Collection   string
}

type Option func(*Config)

func WithContext(ctx context.Context) Option {
	return func(c *Config) {
		if ctx == nil {
			ctx = context.Background()
		}
		c.Ctx = ctx
	}
}

func WithCollection(collection string) Option {
	return func(c *Config) {
		if collection == "" {
			c.Collection = fruitsCollection
		}
		c.Collection = collection
	}
}

func WithLogger(log *logrus.Logger) Option {
	return func(c *Config) {
		c.Log = log
	}
}

func WithDB(databaseName string) Option {
	return func(c *Config) {
		c.databaseName = databaseName
	}
}

//New creates a new instance of Config to create and initialize new database
func New(uri string, options ...Option) *Config {
	cfg := &Config{}
	for _, o := range options {
		o(cfg)
	}
	cfg.URI = uri
	return cfg
}

//Init initializes the database with the given configuration
func (c *Config) Init() error {
	var client *mongo.Client
	var err error
	if client, err = mongo.Connect(c.Ctx, options.Client().ApplyURI(c.URI)); err != nil {
		return err
	}
	c.Client = client
	c.DB = client.Database(c.databaseName)
	return nil
}
