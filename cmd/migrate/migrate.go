// cmd/dbmigrate/main.go

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	var direction string
	flag.StringVar(&direction, "direction", "up", "migrate cmd")
	flag.Parse()
	pgConnString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	m, err := migrate.New("file://migrations/", pgConnString)
	if err != nil {
		log.Printf("%s", err)
		return
	}

	if direction == "up" {
		if err := m.Up(); err != nil {
			log.Printf("failed migrate up: %s", err)
			return
		}
	}

	if direction == "down" {
		if err := m.Down(); err != nil {
			log.Printf("failed migrate down: %s", err)
			return
		}
	}
}
