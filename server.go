package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os/signal"
	"time"

	_ "github.com/kameshsampath/go-fruits-api/docs"
	"github.com/kameshsampath/go-fruits-api/pkg/db"
	"github.com/kameshsampath/go-fruits-api/pkg/routes"
	"github.com/kameshsampath/go-fruits-api/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	// _ "github.com/mattn/go-sqlite3"

	"os"

	_ "github.com/lib/pq"
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
// @BasePath /v1/api
// @query.collection.format multi
// @schemes http https
func main() {
	var v, dbType, dbFile string
	flag.StringVar(&dbType, "dbtype", utils.LookupEnvOrString("FRUITS_DB_TYPE", "sqlite"), "The database to use. Valid values are sqlite, pg, mysql")
	flag.StringVar(&dbFile, "dbPath", utils.LookupEnvOrString("FRUITS_DB_FILE", "/data/db"), "Sqlite DB file")
	flag.StringVar(&v, "level", utils.LookupEnvOrString("LOG_LEVEL", logrus.InfoLevel.String()), "The log level to use. Allowed values trace,debug,info,warn,fatal,panic.")
	flag.Parse()

	log = utils.LogSetup(os.Stdout, v)
	dbc := db.New(
		db.WithContext(context.Background()),
		db.WithLogger(log),
		db.WithDBType(dbType),
		db.WithDBFile(dbFile))
	dbc.Init()
	if err := dbc.DB.Ping(); err != nil {
		log.Fatal("Unable to ping the database")
	}

	router = echo.New()
	router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")

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

	v1 := router.Group("/api/v1")
	{
		//Health Endpoints accessible via /api/v1/health
		health := v1.Group("/health")
		{
			health.GET("/live", endpoints.Live)
			health.GET("/ready", endpoints.Ready)
		}

		//Fruits API endpoints
		fruits := v1.Group("/fruits")
		{
			fruits.POST("/add", endpoints.AddFruit)
			fruits.GET("/", endpoints.ListFruits)
			fruits.DELETE("/:id", endpoints.DeleteFruit)
			fruits.GET(":name", endpoints.GetFruitsByName)
			fruits.GET("/season/:season", endpoints.GetFruitsBySeason)
		}
	}

	// the default path to get swagger json is :8080/swagger/docs.json
	// TODO enable/disable based on ENV variable
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
