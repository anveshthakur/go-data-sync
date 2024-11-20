package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
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

	app.serve()
}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	fmt.Printf("Started webserver at port: %v", webPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
