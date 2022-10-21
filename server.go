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

	"os"

	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/uptrace/bun"
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
	var v, uri, dbName, dbCollection string
	flag.StringVar(&uri, "dbConnectionString", utils.LookupEnvOrString("QUARKUS_MONGODB_CONNECTION_STRING", ""), "The mongodb database connection string.")
	flag.StringVar(&dbName, "dbName", utils.LookupEnvOrString("FRUITS_DB", "demodb"), "The database to use.")
	flag.StringVar(&dbCollection, "dbCollectionName", utils.LookupEnvOrString("FRUITS_DB_COLLECTION", "fruits"), "The Fruits collection in Database")
	flag.StringVar(&v, "level", utils.LookupEnvOrString("LOG_LEVEL", logrus.InfoLevel.String()), "The log level to use. Allowed values trace,debug,info,warn,fatal,panic.")
	flag.Parse()

	log = utils.LogSetup(os.Stdout, v)
	ctx := context.Background()
	dbc := db.New(uri,
		db.WithContext(ctx),
		db.WithDB(dbName),
		db.WithLogger(log),
		db.WithCollection(dbCollection))
	err := dbc.Init()

	if err != nil {
		log.Fatal(err)
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
	defer dbc.Client.Disconnect(ctx)
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
			fruits.GET("/search/:name", endpoints.GetFruitsByName)
			fruits.GET("/season/:season", endpoints.GetFruitsBySeason)
		}
	}

	router.GET("/swagger/*any", echoSwagger.WrapHandler)
}
