package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	routes "github.com/kameshsampath/go-fruits-api/pkg/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

//DDLFRUITSTABLE  creates the database
const DDLFRUITSTABLE = `
DROP TABLE IF EXISTS fruits;
CREATE TABLE IF NOT EXISTS fruits (
id INTEGER PRIMARY KEY AUTOINCREMENT,
name TEXT NOT NULL,
season TEXT NOT NULL,
emoji TEXT)`

var (
	dbDir  string
	db     *sql.DB
	err    error
	dbFile string
	router *gin.Engine
)

func init() {
	if dbDir = os.Getenv("FRUITS_DB_DIR"); dbDir == "" {
		homedir, _ := os.UserHomeDir()
		dbDir = filepath.Join(homedir, ".fruits-app")
	}
	if _, err := os.Stat(dbDir); err != nil {
		if err = os.Mkdir(dbDir, os.ModeDir); err != nil {
			panic(fmt.Sprintf("Error creating DB Dir %s", dbDir))
		}
	}
}

func addRoutes() {
	v1 := router.Group("/v1")
	routes.HealthResource(v1)
	routes.FruitsResource(v1, db)
}

func main() {

	if err == nil {
		dbFile = filepath.Join(dbDir, "fruits.db")
		db, err = sql.Open("sqlite3", dbFile)
		if err != nil {
			log.Fatalf("Error opening DB %s, reason %s", dbFile, err)
		}
		//TODO Graceful shutdown
		defer db.Close()

		_, err = db.Exec(DDLFRUITSTABLE)
		if err != nil {
			log.Fatalf("Error initializing DB: %s", err)
		}
		//Load some data
		loadFruits()
	}

	if mode := os.Getenv("GIN_MODE"); mode != "" {
		gin.SetMode(mode)
	}
	router = gin.Default()
	// this is liberal CORS settings only for demo
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	addRoutes()
	server := &http.Server{
		Handler: router,
		Addr:    ":8080",
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
