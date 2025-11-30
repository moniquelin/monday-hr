package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/moniquelin/monday-hr/internal/api"
	"github.com/moniquelin/monday-hr/internal/data"
	"github.com/moniquelin/monday-hr/internal/dbconn"
)

func main() {
	// Declare an instance of the config struct.
	var cfg api.Config

	// Read the value of the port and env command-line flags into the config struct
	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.Db.Dsn, "db-dsn", os.Getenv("MONDAY_HR_DB_DSN"), "PostgreSQL DSN")
	flag.StringVar(&cfg.JWT.Secret, "jwt-secret", os.Getenv("MONDAY_HR_JWT_SECRET"), "JWT secret")
	flag.DurationVar(&cfg.JWT.Expiry, "jwt-expiry", time.Hour*24, "JWT expiry")
	flag.Parse()

	// Initialize a logger which writes messages to the standard out stream
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Call the openDB() helper function (see below) to create the connection pool,
	db, err := dbconn.OpenDB(cfg.Db.Dsn)
	if err != nil {
		logger.Fatal(err)
	}

	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer db.Close()

	logger.Printf("database connection pool established")

	// Declare an instance of the application struct
	app := &api.Application{
		Config: cfg,
		Logger: logger,
		Models: data.NewModels(db),
	}

	// Declare a HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the HTTP server.
	logger.Printf("starting %s server on %s", cfg.Env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}
