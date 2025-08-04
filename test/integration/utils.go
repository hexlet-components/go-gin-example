package integration

import (
	"database/sql"
	"testing"

	"github.com/gin-gonic/gin"
	db "github.com/hexlet-components/go-gin-example/db/generated"
	"github.com/hexlet-components/go-gin-example/handlers"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

const migrationsDir = "../../db/migrations"

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test DB: %v", err)
		return nil
	}

	if err := applyMigrations(testDB, migrationsDir); err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
		return nil
	}

	return testDB
}

func applyMigrations(database *sql.DB, migrationsDir string) error {
	goose.SetDialect("sqlite3")
	return goose.Up(database, migrationsDir)
}

func setupTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	testDB := setupTestDB(t)
	if testDB == nil {
		t.Fatalf("failed to setup test DB")
	}

	return handlers.SetupRouter(testDB)
}

func setupTestQueries(t *testing.T) (*db.Queries, *sql.DB) {
	t.Helper()
	testDB := setupTestDB(t)
	if testDB == nil {
		t.Fatalf("failed to setup test DB")
	}

	queries := db.New(testDB)
	return queries, testDB
}

func setupTestRouterWithQueries(t *testing.T) (*gin.Engine, *db.Queries) {
	t.Helper()
	queries, testDB := setupTestQueries(t)
	router := handlers.SetupRouter(testDB)
	return router, queries
}
