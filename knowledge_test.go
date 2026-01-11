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

import "testing"

func TestTestBit(t *testing.T) {
	t.Run("nil map returns false", func(t *testing.T) {
		if test_bit(nil, 100) {
			t.Error("expected false for nil map")
		}
	})

	t.Run("empty map returns false", func(t *testing.T) {
		m := make(map[int]bool)
		if test_bit(m, 100) {
			t.Error("expected false for missing key")
		}
	})

	t.Run("returns true for existing key", func(t *testing.T) {
		m := make(map[int]bool)
		m[100] = true
		if !test_bit(m, 100) {
			t.Error("expected true for existing key")
		}
	})
}

func TestSetBit(t *testing.T) {
	t.Run("creates map if nil", func(t *testing.T) {
		var m map[int]bool
		m = set_bit(m, 100)
		if m == nil {
			t.Fatal("expected non-nil map")
		}
		if !m[100] {
			t.Error("expected key to be set")
		}
	})

	t.Run("idempotent add", func(t *testing.T) {
		m := make(map[int]bool)
		m = set_bit(m, 100)
		m = set_bit(m, 100)
		m = set_bit(m, 100)
		if !m[100] {
			t.Error("expected key to be set")
		}
		if len(m) != 1 {
			t.Errorf("expected 1 key, got %d", len(m))
		}
	})

	t.Run("adds multiple keys", func(t *testing.T) {
		m := make(map[int]bool)
		m = set_bit(m, 100)
		m = set_bit(m, 200)
		m = set_bit(m, 300)
		if len(m) != 3 {
			t.Errorf("expected 3 keys, got %d", len(m))
		}
	})
}

func TestClearKnowRec(t *testing.T) {
	t.Run("clears all entries", func(t *testing.T) {
		m := make(map[int]bool)
		m[100] = true
		m[200] = true
		m[300] = true
		clear_know_rec(m)
		if len(m) != 0 {
			t.Errorf("expected 0 keys after clear, got %d", len(m))
		}
	})

	t.Run("safe on nil", func(t *testing.T) {
		clear_know_rec(nil)
	})

	t.Run("safe on empty", func(t *testing.T) {
		m := make(map[int]bool)
		clear_know_rec(m)
		if len(m) != 0 {
			t.Errorf("expected 0 keys, got %d", len(m))
		}
	})
}

func TestSetKnownAndTestKnown(t *testing.T) {
	setupTestPlayer := func(playerID, charID int) func() {
		teg.globals.bx[playerID] = &box{kind: T_player}
		teg.globals.bx[playerID].x_player = &entity_player{}
		teg.globals.bx[charID] = &box{kind: T_char}
		teg.globals.bx[charID].x_char = &entity_char{unit_lord: playerID}

		return func() {
			teg.globals.bx[playerID] = nil
			teg.globals.bx[charID] = nil
			if teg.globals.playerKnowledge != nil {
				delete(teg.globals.playerKnowledge, playerID)
			}
		}
	}

	t.Run("test_known returns false for unknown", func(t *testing.T) {
		cleanup := setupTestPlayer(1001, 2001)
		defer cleanup()

		if test_known(2001, 3001) {
			t.Error("expected false for unknown entity")
		}
	})

	t.Run("set_known then test_known", func(t *testing.T) {
		cleanup := setupTestPlayer(1001, 2001)
		defer cleanup()

		teg.globals.bx[3001] = &box{kind: T_loc}
		defer func() { teg.globals.bx[3001] = nil }()

		set_known(2001, 3001)
		if !test_known(2001, 3001) {
			t.Error("expected true after set_known")
		}
	})

	t.Run("idempotent set_known", func(t *testing.T) {
		cleanup := setupTestPlayer(1001, 2001)
		defer cleanup()

		teg.globals.bx[3001] = &box{kind: T_loc}
		defer func() { teg.globals.bx[3001] = nil }()

		set_known(2001, 3001)
		set_known(2001, 3001)
		set_known(2001, 3001)
		if !test_known(2001, 3001) {
			t.Error("expected true after multiple set_known calls")
		}
	})

	t.Run("who=0 returns false", func(t *testing.T) {
		if test_known(0, 100) {
			t.Error("expected false for who=0")
		}
	})

	t.Run("invalid who returns false", func(t *testing.T) {
		if test_known(99999, 100) {
			t.Error("expected false for invalid who")
		}
	})

	t.Run("set_known ignores invalid who", func(t *testing.T) {
		set_known(99999, 100)
	})

	t.Run("set_known ignores invalid i", func(t *testing.T) {
		cleanup := setupTestPlayer(1001, 2001)
		defer cleanup()

		set_known(2001, 99999)
	})

	t.Run("clear knowledge", func(t *testing.T) {
		cleanup := setupTestPlayer(1001, 2001)
		defer cleanup()

		teg.globals.bx[3001] = &box{kind: T_loc}
		teg.globals.bx[3002] = &box{kind: T_loc}
		defer func() {
			teg.globals.bx[3001] = nil
			teg.globals.bx[3002] = nil
		}()

		set_known(2001, 3001)
		set_known(2001, 3002)

		if !test_known(2001, 3001) || !test_known(2001, 3002) {
			t.Fatal("setup failed: entities not known")
		}

		teg.clearPlayerKnowledge(1001)

		if test_known(2001, 3001) {
			t.Error("expected false after clear")
		}
		if test_known(2001, 3002) {
			t.Error("expected false after clear")
		}
	})
}
