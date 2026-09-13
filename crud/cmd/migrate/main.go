package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/config"
)

const migrationsPath = "file://migrations"

func main() {
	// # Migratio statring
	
	// # Step 1 : Check is migration command is provided
	if len(os.Args) < 2 {
		log.Fatal("usage:migrate <up | down>")
	}

	// # Step 2 : Load configuration for database connection string
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	dbURL := cfg.Database.ConnectionString

	// # Step 3 : Initialize migration instance
	m, err := migrate.New(
		migrationsPath,
		dbURL,
	)
	if err != nil {
		log.Fatalf("failed to initialize migration: %v", err)
	}

	defer m.Close()

	command := os.Args[1]

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("failed to run migration up: %v", err)
		}

		fmt.Println("Migrations up applied successfully.")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("failed to run migration down: %v", err)
		}
		fmt.Println("Migrations down rolled back successfully.")

	default:
		log.Fatalf("Unknown migration command : %s", command)
	}

	fmt.Println("Successfully completed migration command:", command)
}
