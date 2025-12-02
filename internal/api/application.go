package api

import (
	"log"

	"github.com/moniquelin/monday-hr/internal/data"
)

// Version number
const Version = "1.0.0"

// Config struct holds all the configuration settings for our application
type Config struct {
	Port int
	Env  string
	Db   struct {
		Dsn string
	}
	Jwt struct {
		Secret string
	}
}

// Application struct holds the dependencies for our HTTP handlers, helpers,
// and middleware
type Application struct {
	Config Config
	Logger *log.Logger
	Models data.Models
}
