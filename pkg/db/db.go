package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"

	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"

	// used to import mysql drivers
	_ "github.com/go-sql-driver/mysql"
)

// Config configures the database to initialize
type Config struct {
	Log    *logrus.Logger
	DBFile string
	DB     *bun.DB
	DBType dialect.Name
}

type Option func(*Config)

func WithLogger(log *logrus.Logger) Option {
	return func(c *Config) {
		c.Log = log
	}
}

func WithDBFile(dbFile string) Option {
	return func(c *Config) {
		if dbFile == "" {
			dbFile = "/data/db"
		}
		c.DBFile = dbFile
	}
}

func WithDBType(dbType string) Option {
	return func(c *Config) {
		switch dbType {
		case "pgsql":
			c.DBType = dialect.PG
		case "mysql":
			c.DBType = dialect.MySQL
		case "sqlite":
			c.DBType = dialect.SQLite
		default:
			c.DBType = dialect.SQLite
		}
	}
}

// New creates a new instance of Config to create and initialize new database
func New(options ...Option) *Config {
	cfg := &Config{}
	for _, o := range options {
		o(cfg)
	}

	return cfg
}

// Init initializes the database with the given configuration
func (c *Config) Init() {
	log := c.Log
	log.Infof("Initializing DB of type %s", c.DBType)
	var db *bun.DB
	switch c.DBType {
	case dialect.PG:
		pgConn := buildPGConnector()
		sqldb := sql.OpenDB(pgConn)
		db = bun.NewDB(sqldb, pgdialect.New())
	case dialect.MySQL:
		sqldb, err := sql.Open("mysql", buildMYSQLDSN())
		if err != nil {
			log.Fatal(err)
		}
		db = bun.NewDB(sqldb, mysqldialect.New())
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
	default:
		sqlite, err := sql.Open(sqliteshim.ShimName, fmt.Sprintf("file:%s?cache=shared", c.DBFile))
		if err != nil {
			log.Fatal(err)
		}
		db = bun.NewDB(sqlite, sqlitedialect.New())
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	isVerbose := log.Level == logrus.DebugLevel || log.Level == logrus.TraceLevel
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(isVerbose),
		bundebug.WithVerbose(isVerbose),
	))

	c.DB = db

	//Setup Schema
	if err := c.createTables(); err != nil {
		log.Errorf("%s", err)
	}
}

func (c *Config) createTables() error {
	//Fruits
	if _, err := c.DB.NewCreateTable().
		Model((*Fruit)(nil)).
		IfNotExists().
		Exec(context.Background()); err != nil {
		return err
	}
	return nil
}

func buildPGConnector() *pgdriver.Connector {
	var (
		pgHost     = "localhost"
		pgPort     = "5432"
		pgUser     = "demo"
		pgPassword = "pa55Word!"
		pgDatabase = "demodb"
	)

	if h, ok := os.LookupEnv("POSTGRES_HOST"); ok {
		pgHost = h
	}
	if p, ok := os.LookupEnv("POSTGRES_PORT"); ok {
		pgPort = p
	}
	if h, ok := os.LookupEnv("POSTGRES_USER"); ok {
		pgUser = h
	}
	if p, ok := os.LookupEnv("POSTGRES_PASSWORD"); ok {
		pgPassword = p
	}
	if d, ok := os.LookupEnv("POSTGRES_DB"); ok {
		pgDatabase = d
	}

	pgConn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(fmt.Sprintf("%s:%s", pgHost, pgPort)),
		pgdriver.WithPassword(pgPassword),
		pgdriver.WithUser(pgUser),
		pgdriver.WithDatabase(pgDatabase),
		pgdriver.WithTLSConfig(nil),
		pgdriver.WithTimeout(5*time.Second),
		pgdriver.WithDialTimeout(5*time.Second),
		pgdriver.WithReadTimeout(5*time.Second),
		pgdriver.WithWriteTimeout(5*time.Second),
		pgdriver.WithApplicationName("fruits-api"),
	)

	log.Infof("PostgresSQL Address %s", pgConn.Config().Addr)

	return pgConn
}

func buildMYSQLDSN() string {
	var (
		mySQLHost     = "localhost"
		mySQLPort     = "3306"
		mySQLUser     = "demo"
		mySQLPassword = "pa55Word!"
		mySQLDatabase = "demodb"
		mySQLProtocol = "tcp"
	)

	if h, ok := os.LookupEnv("MYSQL_HOST"); ok {
		mySQLHost = h
	}
	if p, ok := os.LookupEnv("MYSQL_PORT"); ok {
		mySQLPort = p
	}
	if h, ok := os.LookupEnv("MYSQL_USER"); ok {
		mySQLUser = h
	}
	if p, ok := os.LookupEnv("MYSQL_PASSWORD"); ok {
		mySQLPassword = p
	}
	if d, ok := os.LookupEnv("MYSQL_DB"); ok {
		mySQLDatabase = d
	}
	if d, ok := os.LookupEnv("MYSQL_PROTOCOL"); ok {
		mySQLProtocol = d
	}

	return fmt.Sprintf("%s:%s@%s(%s:%s)/%s",
		mySQLUser, mySQLPassword, mySQLProtocol, mySQLHost, mySQLPort, mySQLDatabase)
}
