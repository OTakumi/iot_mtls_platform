package persistence_test

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var testDB *gorm.DB //nolint:gochecknoglobals

func TestMain(m *testing.M) {
	// Set up the entire test suite.
	err := setupTestDatabase()
	if err != nil {
		log.Fatalf("failed to set up test database: %v", err)
	}

	// Run all tests.
	code := m.Run()

	// Teardown.
	// Close the database connection.
	sqlDB, err := testDB.DB()
	if err == nil {
		err = sqlDB.Close()
		if err != nil {
			log.Printf("failed to close test database: %v", err)
		}
	}

	os.Exit(code)
}

var ErrDSNTestNotSet = errors.New("environment variable DSN_TEST is not set")

func setupTestDatabase() error {
	// Get the test database DSN from environment variables.
	dsnTest := os.Getenv("DSN_TEST")
	if dsnTest == "" {
		return ErrDSNTestNotSet
	}

	// Run migrations.
	migrationURL := "file://../../../../infra/db-auth/migrations"

	// Use dsnTest directly as the second argument for migrate.New.
	mi, err := migrate.New(migrationURL, dsnTest)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Clear the database state before running migrations.
	err = mi.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Printf("migrate down failed, but continuing test: %v", err)
	}

	err = mi.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up failed: %w", err)
	}

	// Connect to the database using GORM.
	// Simple GORM configuration with a silent logger.
	db, err := gorm.Open(postgres.Open(dsnTest), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %w", err)
	}

	testDB = db

	return nil
}

func cleanupTable(t *testing.T) {
	t.Helper()

	err := testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", "devices")).Error
	if err != nil {
		t.Fatalf("failed to cleanup table (%s): %v", "devices", err)
	}
}
