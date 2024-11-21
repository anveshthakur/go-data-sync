package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	SourceDb *sql.DB
	TargetDB *sql.DB
}

var webPort string = "8080"

func main() {
	app := Config{
		SourceDb: nil,
		TargetDB: nil,
	}

	go app.listenForShutdown()

	app.serve()
}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Started webserver at port: %v \n", webPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
func (c *Config) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	c.shutdown()
	os.Exit(0)
}

func (c *Config) shutdown() {
	log.Printf("Cleaning up resources...")
	if c.SourceDb != nil {
		if err := c.SourceDb.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed.")
		}
	}

	if c.SourceDb != nil {
		if err := c.SourceDb.Close(); err != nil {
			log.Printf("Error closing source database connection: %v", err)
		} else {
			log.Println("Source Database connection closed.")
		}
	}

	if c.TargetDB != nil {
		if err := c.TargetDB.Close(); err != nil {
			log.Printf("Error closing  target database connection: %v", err)
		} else {
			log.Println("Target Database connection closed.")
		}
	}

}
