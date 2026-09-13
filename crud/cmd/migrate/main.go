package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/config"
)

const (
	migrationsPath   = "file://migrations"
	migrationsDir    = "./migrations"
	migrationCommand = "./cmd/migrate"
	timestampFmt     = "20060102150405"
)

func main() {
	fmt.Println("MIGRATE CLI STARTED")
	fmt.Println("timestamp format:", timestampFmt)

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "create":
		createMigration(os.Args[2:])

	case "up":
		runUp()

	case "down":
		runDown()

	case "down-all":
		runDownAll()

	default:
		slog.Error(
			"Unknown migration command",
			"command", command,
		)
		usage()
		os.Exit(1)
	}
}

// createMigration creates timestamp-based up/down migration files.
//
// Usage:
//
//	go run ./cmd/migrate create create_users_table
func createMigration(args []string) {
	if len(args) != 1 {
		slog.Error(
			"Migration name is required",
			"usage", "migrate create <name>",
		)
		os.Exit(1)
	}

	name := sanitizeMigrationName(args[0])

	if name == "" {
		slog.Error("Invalid migration name")
		os.Exit(1)
	}

	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		slog.Error(
			"Failed to create migrations directory",
			"error", err,
		)
		os.Exit(1)
	}

	timestamp := time.Now().UTC().Format(timestampFmt)

	upFile := filepath.Join(
		migrationsDir,
		fmt.Sprintf("%s_%s.up.sql", timestamp, name),
	)

	downFile := filepath.Join(
		migrationsDir,
		fmt.Sprintf("%s_%s.down.sql", timestamp, name),
	)

	if err := createFile(upFile); err != nil {
		slog.Error(
			"Failed to create up migration",
			"file", upFile,
			"error", err,
		)
		os.Exit(1)
	}

	if err := createFile(downFile); err != nil {
		// Clean up the .up.sql file if .down.sql creation fails.
		_ = os.Remove(upFile)

		slog.Error(
			"Failed to create down migration",
			"file", downFile,
			"error", err,
		)
		os.Exit(1)
	}

	slog.Info(
		"Migration files created successfully",
		"version", timestamp,
		"name", name,
		"up", upFile,
		"down", downFile,
	)
}

// runUp applies all pending migrations.
func runUp() {
	m := newMigration()
	defer closeMigration(m)

	slog.Info("Starting database migration", "command", "up")

	err := m.Up()

	if errors.Is(err, migrate.ErrNoChange) {
		slog.Info(
			"No pending migrations",
			"command", "up",
		)
		return
	}

	if err != nil {
		slog.Error(
			"Failed to apply migrations",
			"command", "up",
			"error", err,
		)
		os.Exit(1)
	}

	slog.Info(
		"Database migrations applied successfully",
		"command", "up",
	)
}

// runDown rolls back the latest migration.
func runDown() {
	m := newMigration()
	defer closeMigration(m)

	slog.Info("Starting database migration", "command", "down")

	err := m.Steps(-1)

	if errors.Is(err, migrate.ErrNoChange) {
		slog.Info(
			"No migration available to rollback",
			"command", "down",
		)
		return
	}

	if err != nil {
		slog.Error(
			"Failed to rollback migration",
			"command", "down",
			"error", err,
		)
		os.Exit(1)
	}

	slog.Info(
		"Latest migration rolled back successfully",
		"command", "down",
	)
}

// runDownAll rolls back all migrations.
//
// This operation is destructive and should not happen accidentally.
func runDownAll() {
	if len(os.Args) < 3 || os.Args[2] != "--confirm" {
		slog.Error(
			"down-all requires explicit confirmation",
			"usage", "migrate down-all --confirm",
		)
		os.Exit(1)
	}

	m := newMigration()
	defer closeMigration(m)

	slog.Warn(
		"Rolling back all database migrations",
		"command", "down-all",
	)

	err := m.Down()

	if errors.Is(err, migrate.ErrNoChange) {
		slog.Info(
			"No migrations available to rollback",
			"command", "down-all",
		)
		return
	}

	if err != nil {
		slog.Error(
			"Failed to rollback all migrations",
			"command", "down-all",
			"error", err,
		)
		os.Exit(1)
	}

	slog.Warn(
		"All database migrations rolled back successfully",
		"command", "down-all",
	)
}

// newMigration loads configuration and initializes the migration engine.
func newMigration() *migrate.Migrate {
	cfg, err := config.Load()
	if err != nil {
		slog.Error(
			"Failed to load configuration",
			"error", err,
		)
		os.Exit(1)
	}

	m, err := migrate.New(
		migrationsPath,
		cfg.Database.ConnectionString,
	)
	if err != nil {
		slog.Error(
			"Failed to initialize migration",
			"error", err,
		)
		os.Exit(1)
	}

	return m
}

// closeMigration closes migration resources.
func closeMigration(m *migrate.Migrate) {
	sourceErr, databaseErr := m.Close()

	if sourceErr != nil {
		slog.Error(
			"Failed to close migration source",
			"error", sourceErr,
		)
	}

	if databaseErr != nil {
		slog.Error(
			"Failed to close migration database connection",
			"error", databaseErr,
		)
	}
}

// createFile creates an empty SQL migration file.
func createFile(path string) error {
	file, err := os.OpenFile(
		path,
		os.O_WRONLY|os.O_CREATE|os.O_EXCL,
		0644,
	)
	if err != nil {
		return err
	}

	return file.Close()
}

// sanitizeMigrationName converts the migration name into a safe filename.
func sanitizeMigrationName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	name = strings.ReplaceAll(name, " ", "_")

	re := regexp.MustCompile(`[^a-z0-9_-]+`)
	name = re.ReplaceAllString(name, "")

	name = strings.Trim(name, "_-")

	return name
}

func usage() {
	fmt.Printf(`
Migration CLI

Usage:
	go run %s create <name>
	go run %s up
	go run %s down
	go run %s down-all --confirm

Examples:
	go run %s create create_users_table
	go run %s create add_phone_to_users

	go run %s up
	go run %s down
	go run %s down-all --confirm
`, migrationCommand, migrationCommand, migrationCommand, migrationCommand,
		migrationCommand, migrationCommand, migrationCommand, migrationCommand,
		migrationCommand)
}
