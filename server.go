package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os/signal"
	"path"
	"time"

	_ "github.com/kameshsampath/go-fruits-api/docs"
	"github.com/kameshsampath/go-fruits-api/pkg/db"
	"github.com/kameshsampath/go-fruits-api/pkg/routes"
	"github.com/kameshsampath/go-fruits-api/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"os"

	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
)

var (
	log            *logrus.Logger
	httpListenPort = "8080"
	router         *echo.Echo
)

// @title Fruits API
// @version 1.0
// @description The Fruits API that defines few REST operations with Fruits used for demos

// @contact.name Kamesh Sampath
// @contact.email kamesh.sampath@hotmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @query.collection.format multi
// @schemes http https
func main() {
	var v, dbType, dbFile, dataDir string
	flag.StringVar(&dbType, "dbType", utils.LookupEnvOrString("FRUITS_DB_TYPE", "sqlite"), "The database to use. Valid values are sqlite, pgsql, mysql")
	flag.StringVar(&dbFile, "dbPath", utils.LookupEnvOrString("FRUITS_DB_FILE", "/data/db"), "Sqlite DB file")
	flag.StringVar(&dataDir, "dataDir", "", "The data dir that will have the 'data.yaml' that will be loaded on to the Fruits table.")
	flag.StringVar(&v, "level", utils.LookupEnvOrString("LOG_LEVEL", logrus.InfoLevel.String()), "The log level to use. Allowed values trace,debug,info,warn,fatal,panic.")
	flag.Parse()

	ctx := context.Background()
	log = utils.LogSetup(os.Stdout, v)
	dbc := db.New(
		db.WithLogger(log),
		db.WithDBType(dbType),
		db.WithDBFile(dbFile))
	dbc.Init(ctx)

	//marker file to ensure we don't preload the data again on each
	//update of the application
	_, err := os.Stat(path.Join("/data", "db", ".loaded"))
	if dataDir != "" && errors.Is(err, os.ErrNotExist) {
		log.Info("Attempting to preload data")
		fixtures := dbfixture.New(dbc.DB, dbfixture.WithTruncateTables())
		if err := fixtures.Load(ctx, os.DirFS(dataDir), "data.yaml"); err != nil {
			log.Warnf("unable to preload the data,%v", err)
		}
		_, err := os.Create(path.Join("/data", "db", ".loaded"))
		if err != nil {
			log.Errorf("Error creating marker file %v", err)
		}
	} else {
		log.Info("Data already loaded, skipping preload.")
	}

	router = echo.New()
	router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogError:  true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Debug("request")

			return nil
		},
	}))
	router.Use(middleware.Recover())
	addRoutes(dbc)

	// Start server
	go func() {
		if p, ok := os.LookupEnv("HTTP_LISTEN_PORT"); ok {
			httpListenPort = p
		}
		if err := router.Start(fmt.Sprintf(":%s", httpListenPort)); err != nil && err.Error() != http.ErrServerClosed.Error() {
			router.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := router.Shutdown(ctx); err != nil {
		router.Logger.Fatal(err)
	}
}

func addRoutes(dbc *db.Config) {
	endpoints := routes.NewEndpoints(dbc)

	v1 := router.Group("/api")
	{
		//Health Endpoints accessible via /api/health
		health := v1.Group("/health")
		{
			health.GET("/live", endpoints.Live)
			health.GET("/ready", endpoints.Ready)
		}

		//Fruits API endpoints /api/fruits
		fruits := v1.Group("/fruits")
		{
			fruits.POST("/add", endpoints.AddFruit)
			fruits.GET("/", endpoints.ListFruits)
			fruits.DELETE("/:id", endpoints.DeleteFruit)
			fruits.DELETE("/", endpoints.DeleteAll)
			fruits.GET("/search/:name", endpoints.GetFruitsByName)
			fruits.GET("/season/:season", endpoints.GetFruitsBySeason)
		}
	}

	router.GET("/swagger/*any", echoSwagger.WrapHandler)
}
