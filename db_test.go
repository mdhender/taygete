// taygete - a game engine for a game.
// Copyright (c) 2026 Michael D Henderson.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package taygete

import (
	"testing"
)

func TestOpenTestDB(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	// Verify foreign keys are enabled
	var fkEnabled int
	err = db.QueryRow("PRAGMA foreign_keys").Scan(&fkEnabled)
	if err != nil {
		t.Fatalf("query foreign_keys: %v", err)
	}
	if fkEnabled != 1 {
		t.Errorf("foreign_keys = %d, want 1", fkEnabled)
	}
}

func TestMigrationsApplied(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	// Check that schema_migrations table exists and has entries
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count)
	if err != nil {
		t.Fatalf("query schema_migrations: %v", err)
	}
	if count == 0 {
		t.Error("no migrations recorded")
	}

	// Check that key tables exist
	tables := []string{
		"game_meta", "accounts", "players", "entities",
		"locations", "characters", "turns", "commands",
	}
	for _, table := range tables {
		var name string
		err = db.QueryRow(
			"SELECT name FROM sqlite_master WHERE type='table' AND name=?",
			table,
		).Scan(&name)
		if err != nil {
			t.Errorf("table %s not found: %v", table, err)
		}
	}
}

func TestMigrationsIdempotent(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	// Running migrations again should be a no-op
	err = runMigrations(db)
	if err != nil {
		t.Fatalf("runMigrations (second call): %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count)
	if err != nil {
		t.Fatalf("query schema_migrations: %v", err)
	}
	// Should still have exactly one migration
	if count != 1 {
		t.Errorf("migration count = %d, want 1", count)
	}
}

func TestBeginTurnCommit(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	ttx, err := BeginTurn(db, 1)
	if err != nil {
		t.Fatalf("BeginTurn: %v", err)
	}

	if ttx.TurnNumber() != 1 {
		t.Errorf("TurnNumber = %d, want 1", ttx.TurnNumber())
	}

	// Check turn is processing
	var status string
	err = ttx.Tx().QueryRow("SELECT status FROM turns WHERE turn_number = 1").Scan(&status)
	if err != nil {
		t.Fatalf("query turn status: %v", err)
	}
	if status != "processing" {
		t.Errorf("status = %q, want 'processing'", status)
	}

	// Commit
	err = ttx.Commit()
	if err != nil {
		t.Fatalf("Commit: %v", err)
	}

	// Check turn is finished
	err = db.QueryRow("SELECT status FROM turns WHERE turn_number = 1").Scan(&status)
	if err != nil {
		t.Fatalf("query turn status after commit: %v", err)
	}
	if status != "finished" {
		t.Errorf("status = %q, want 'finished'", status)
	}
}

func TestBeginTurnRollback(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	ttx, err := BeginTurn(db, 2)
	if err != nil {
		t.Fatalf("BeginTurn: %v", err)
	}

	// Insert some data in the transaction
	_, err = ttx.Tx().Exec("INSERT INTO game_meta (id, game_name) VALUES (1, 'test')")
	if err != nil {
		t.Fatalf("insert game_meta: %v", err)
	}

	// Rollback
	err = ttx.Rollback()
	if err != nil {
		t.Fatalf("Rollback: %v", err)
	}

	// Check that game_meta is empty (rolled back)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM game_meta").Scan(&count)
	if err != nil {
		t.Fatalf("query game_meta: %v", err)
	}
	if count != 0 {
		t.Errorf("game_meta count = %d, want 0 (rollback failed)", count)
	}

	// Check that turn record was also rolled back
	err = db.QueryRow("SELECT COUNT(*) FROM turns WHERE turn_number = 2").Scan(&count)
	if err != nil {
		t.Fatalf("query turns: %v", err)
	}
	if count != 0 {
		t.Errorf("turns count = %d, want 0 (rollback failed)", count)
	}
}

func TestIsolatedTestDBs(t *testing.T) {
	db1, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB 1: %v", err)
	}
	defer db1.Close()

	db2, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB 2: %v", err)
	}
	defer db2.Close()

	// Insert into db1
	_, err = db1.Exec("INSERT INTO game_meta (id, game_name) VALUES (1, 'game1')")
	if err != nil {
		t.Fatalf("insert into db1: %v", err)
	}

	// db2 should be empty
	var count int
	err = db2.QueryRow("SELECT COUNT(*) FROM game_meta").Scan(&count)
	if err != nil {
		t.Fatalf("query db2: %v", err)
	}
	if count != 0 {
		t.Errorf("db2 game_meta count = %d, want 0 (dbs not isolated)", count)
	}
}
