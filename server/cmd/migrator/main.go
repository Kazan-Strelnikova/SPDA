package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	dbURL := flag.String("db", "", "PostgreSQL connection string")
	migrationsPath := flag.String("path", "", "Path to migration files")
	action := flag.String("action", "up", "Migration action: up or down")

	flag.Parse()

	if *dbURL == "" || *migrationsPath == "" {
		log.Fatal("Database connection string and migration path are required")
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", *migrationsPath),
		*dbURL,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migration: %v", err)
	}

	switch *action {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migrations applied successfully.")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Migrations reverted successfully.")

	default:
		log.Fatalf("Invalid action: %s. Use 'up' or 'down'", *action)
	}
}
