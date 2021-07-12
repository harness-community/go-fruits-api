package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/kameshsampath/go-fruits-api/docs"
	"github.com/kameshsampath/go-fruits-api/pkg/data"
	"github.com/kameshsampath/go-fruits-api/pkg/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	db             *sql.DB
	err            error
	dbFile         string
	router         *gin.Engine
	httpListenPort = "8080"
	pgHost         = "localhost"
	pgPort         = "5432"
	pgUser         = "demo"
	pgPassword     = "pa55Word!"
	pgDatabase     = "demodb"
)

// @title Fruits API
// @version 1.0
// @description The Fruits API that defines few REST operations with Fruits used for demos

// @contact.name Kamesh Sampath
// @contact.email kamesh.sampath@solo.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1/api
// @query.collection.format multi
// @schemes http https
func main() {
	if h := os.Getenv("POSTGRES_HOST"); h != "" {
		pgHost = h
	}
	if p := os.Getenv("POSTGRES_PORT"); p != "" {
		pgPort = p
	}
	if h := os.Getenv("POSTGRES_USER"); h != "" {
		pgUser = h
	}
	if p := os.Getenv("POSTGRES_PASSWORD"); p != "" {
		pgPassword = p
	}
	if d := os.Getenv("POSTGRES_DB"); d != "" {
		pgDatabase = d
	}
	pgsqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgHost, pgPort, pgUser, pgPassword, pgDatabase)
	db, err = sql.Open("postgres", pgsqlInfo)
	if err != nil {
		log.Fatalf("Error opening DB %s, reason %s", dbFile, err)
	}
	//TODO Graceful shutdown
	defer db.Close()

	_, err = db.Exec(data.DDLFRUITSTABLE)
	if err != nil {
		log.Fatalf("Error initializing DB: %s", err)
	}
	//Load some data
	loadFruits()

	if mode := os.Getenv("GIN_MODE"); mode != "" {
		gin.SetMode(mode)
	}
	router = gin.Default()
	addRoutes()
	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%s", httpListenPort),
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listent: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server shutting down server ...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced shutdown", err)
	}

	log.Println("Server Exiting")
}

func addRoutes() {
	endpoints := routes.NewEndpoints()
	endpoints.DB = db
	v1 := router.Group("/v1/api")
	{
		//Health Endpoints accessible via /v1/api/health
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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func loadFruits() {
	log.Println("Loading data into fruits table")
	data := `
DELETE FROM fruits;
INSERT INTO fruits(name,season,emoji) VALUES ('Mango','Spring','U+1F96D');
INSERT INTO fruits(name,season,emoji) VALUES ('Strawberry','Spring','U+1F353');
INSERT INTO fruits(name,season,emoji) VALUES ('Orange','Winter','U+1F34A');
INSERT INTO fruits(name,season,emoji) VALUES ('Lemon','Winter','U+1F34B');
INSERT INTO fruits(name,season,emoji) VALUES ('Blueberry','Summer','U+1FAD0');
INSERT INTO fruits(name,season,emoji) VALUES ('Banana','Summer','U+1F34C');
INSERT INTO fruits(name,season,emoji) VALUES ('Watermelon','Summer','U+1F349');
INSERT INTO fruits(name,season,emoji) VALUES ('Apple','Fall','U+1F34E');
INSERT INTO fruits(name,season,emoji) VALUES ('Pear','Fall','U+1F350');
`
	if _, err = db.Exec(data); err != nil {
		log.Fatalf("Error loading data into fruits table %s", err)
	}
}
