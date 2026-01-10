// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// OpenGameDB opens or creates a SQLite3 database for a game.
// Use ":memory:" for an in-memory database (useful for tests).
// Use a file path for a persistent database.
func OpenGameDB(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, errors.New("dsn is required")
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return db, nil
}

// OpenTestDB creates an in-memory database for testing.
// Each call returns a fresh, isolated database.
func OpenTestDB() (*sql.DB, error) {
	return OpenGameDB(":memory:")
}

// runMigrations applies all pending migrations to the database.
func runMigrations(db *sql.DB) error {
	// Create migrations tracking table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	// Get list of applied migrations
	applied := make(map[string]bool)
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("query migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return fmt.Errorf("scan migration version: %w", err)
		}
		applied[version] = true
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate migrations: %w", err)
	}

	// Get list of migration files
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var migrations []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			migrations = append(migrations, entry.Name())
		}
	}
	sort.Strings(migrations)

	// Apply pending migrations
	for _, name := range migrations {
		version := strings.TrimSuffix(name, ".sql")
		if applied[version] {
			continue
		}

		content, err := fs.ReadFile(migrationsFS, "migrations/"+name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin transaction for %s: %w", name, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("execute migration %s: %w", name, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("record migration %s: %w", name, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", name, err)
		}
	}

	return nil
}

// TurnTx wraps a database transaction for turn processing.
type TurnTx struct {
	tx         *sql.Tx
	turnNumber int
}

// BeginTurn starts a new transaction for processing a turn.
func BeginTurn(db *sql.DB, turnNumber int) (*TurnTx, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin turn %d: %w", turnNumber, err)
	}

	// Ensure the turn record exists
	_, err = tx.Exec(`
		INSERT INTO turns (turn_number, status, started_at)
		VALUES (?, 'processing', CURRENT_TIMESTAMP)
		ON CONFLICT(turn_number) DO UPDATE SET
			status = 'processing',
			started_at = CURRENT_TIMESTAMP
	`, turnNumber)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("insert turn record: %w", err)
	}

	return &TurnTx{tx: tx, turnNumber: turnNumber}, nil
}

// Tx returns the underlying transaction for executing queries.
func (t *TurnTx) Tx() *sql.Tx {
	return t.tx
}

// TurnNumber returns the turn number being processed.
func (t *TurnTx) TurnNumber() int {
	return t.turnNumber
}

// Commit commits the turn transaction and marks the turn as finished.
func (t *TurnTx) Commit() error {
	_, err := t.tx.Exec(`
		UPDATE turns
		SET status = 'finished', finished_at = CURRENT_TIMESTAMP
		WHERE turn_number = ?
	`, t.turnNumber)
	if err != nil {
		t.tx.Rollback()
		return fmt.Errorf("update turn status: %w", err)
	}

	if err := t.tx.Commit(); err != nil {
		return fmt.Errorf("commit turn %d: %w", t.turnNumber, err)
	}
	return nil
}

// Rollback aborts the turn transaction.
func (t *TurnTx) Rollback() error {
	return t.tx.Rollback()
}
