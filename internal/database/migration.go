package database

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string) {
	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()

	if err == migrate.ErrNoChange {
		log.Println("no new migrations")
	} else if err != nil {
		log.Fatal(err)
	}

	log.Println("database migrations applied successfully")
}
