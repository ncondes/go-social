package testutils

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/ncondes/go/social/internal/env"
)

func NewTestDB(t *testing.T) *sql.DB {
	t.Helper() // Marks the function as a helper function

	addr := env.GetString("DB_TEST_ADDR", "postgres://postgres:password@localhost:5433/social_test?sslmode=disable")

	db, err := sql.Open("postgres", addr)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	return db
}

func TeardownTestDB(t *testing.T, db *sql.DB) {
	t.Helper() // Marks the function as a helper function

	if err := db.Close(); err != nil {
		t.Fatalf("failed to close database: %v", err)
	}
}

func TruncateTables(t *testing.T, db *sql.DB, tables ...string) {
	t.Helper() // Marks the function as a helper function

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("failed to truncate table %s: %v", table, err)
		}
	}
}

func CreateTestUser(t *testing.T, db *sql.DB) int64 {
	t.Helper() // Marks the function as a helper function

	timestamp := time.Now().UnixNano()
	username := fmt.Sprintf("testuser_%d", timestamp)
	email := fmt.Sprintf("testuser_%d@example.com", timestamp)

	// Get the default "user" role ID
	var roleID int64
	err := db.QueryRow(`SELECT id FROM roles WHERE name = 'user' LIMIT 1`).Scan(&roleID)
	if err != nil {
		t.Fatalf("failed to get user role: %v", err)
	}

	var userID int64
	err = db.QueryRow(`
        INSERT INTO users (first_name, last_name, username, email, password, role_id)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`,
		"Test", "User", username, email, "hashedpassword", roleID,
	).Scan(&userID)

	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	return userID
}
