package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

const (
	DefaultDBFile        = "app.db"
	DefaultMigrationsDir = "db/migrations"
)

// MigrationOptions contains configuration for database migrations
type MigrationOptions struct {
	DBFile        string
	MigrationsDir string
	Dialect       string
}

// DefaultMigrationOptions returns default migration options
func DefaultMigrationOptions() *MigrationOptions {
	return &MigrationOptions{
		DBFile:        DefaultDBFile,
		MigrationsDir: DefaultMigrationsDir,
		Dialect:       "sqlite3",
	}
}

// openDB opens a database connection with the given options
func openDB(opts *MigrationOptions) (*sql.DB, error) {
	goose.SetDialect(opts.Dialect)

	db, err := sql.Open(opts.Dialect, opts.DBFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	if _, err := goose.EnsureDBVersion(db); err != nil {
		// This is not fatal, the table will be created by the first migration
		fmt.Printf("Warning: failed to ensure goose_db_version table (will be created by first migration): %v\n", err)
	}

	return db, nil
}

// MigrateUp applies all available migrations
func MigrateUp(opts *MigrationOptions) error {
	if opts == nil {
		opts = DefaultMigrationOptions()
	}

	db, err := openDB(opts)
	if err != nil {
		return err
	}
	defer db.Close()

	return goose.Up(db, opts.MigrationsDir)
}

// MigrateDown rolls back the last migration
func MigrateDown(opts *MigrationOptions) error {
	if opts == nil {
		opts = DefaultMigrationOptions()
	}

	db, err := openDB(opts)
	if err != nil {
		return err
	}
	defer db.Close()

	return goose.Down(db, opts.MigrationsDir)
}

// MigrateStatus shows the current migration status
func MigrateStatus(opts *MigrationOptions) error {
	if opts == nil {
		opts = DefaultMigrationOptions()
	}

	db, err := openDB(opts)
	if err != nil {
		return err
	}
	defer db.Close()

	return goose.Status(db, opts.MigrationsDir)
}

// MigrateReset removes the database file and applies all migrations from scratch
func MigrateReset(opts *MigrationOptions) error {
	if opts == nil {
		opts = DefaultMigrationOptions()
	}

	// Remove the database file
	if err := os.Remove(opts.DBFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove DB file: %w", err)
	}

	// Recreate and migrate
	return MigrateUp(opts)
}

// SetupTestDB creates a test database with all migrations applied
func SetupTestDB(dbFile string) (*sql.DB, error) {
	opts := &MigrationOptions{
		DBFile:        dbFile,
		MigrationsDir: DefaultMigrationsDir,
		Dialect:       "sqlite3",
	}

	// Remove existing test DB if it exists
	if err := os.Remove(dbFile); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to remove test DB file: %w", err)
	}

	// Apply migrations
	if err := MigrateUp(opts); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Open and return the database connection
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open test DB: %w", err)
	}

	return db, nil
}

// CleanupTestDB removes the test database file
func CleanupTestDB(dbFile string) error {
	if err := os.Remove(dbFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove test DB file: %w", err)
	}
	return nil
}
