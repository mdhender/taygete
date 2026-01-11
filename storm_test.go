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

// TestShipMoving tests the ship_moving function.
func TestShipMoving(t *testing.T) {
	// Save and restore global state
	oldBx := teg.globals.bx
	defer func() { teg.globals.bx = oldBx }()

	// Clear bx array for testing
	teg.globals.bx = [MAX_BOXES]*box{}

	// Test 1: nil box returns 0
	if got := ship_moving(1000); got != 0 {
		t.Errorf("ship_moving(nil box) = %d, want 0", got)
	}

	// Test 2: box with nil subloc returns 0
	teg.globals.bx[1001] = &box{}
	if got := ship_moving(1001); got != 0 {
		t.Errorf("ship_moving(nil subloc) = %d, want 0", got)
	}

	// Test 3: ship not moving returns 0
	teg.globals.bx[1002] = &box{
		x_subloc: &entity_subloc{moving: 0},
	}
	if got := ship_moving(1002); got != 0 {
		t.Errorf("ship_moving(not moving) = %d, want 0", got)
	}

	// Test 4: ship moving returns daystamp
	teg.globals.bx[1003] = &box{
		x_subloc: &entity_subloc{moving: 42},
	}
	if got := ship_moving(1003); got != 42 {
		t.Errorf("ship_moving(moving) = %d, want 42", got)
	}
}

// TestShipGone tests the ship_gone function.
func TestShipGone(t *testing.T) {
	// Save and restore global state
	oldBx := teg.globals.bx
	oldSysclock := teg.globals.sysclock
	oldEvening := teg.globals.evening
	defer func() {
		teg.globals.bx = oldBx
		teg.globals.sysclock = oldSysclock
		teg.globals.evening = oldEvening
	}()

	// Clear state for testing
	teg.globals.bx = [MAX_BOXES]*box{}
	teg.globals.sysclock = olytime{days_since_epoch: 100}
	teg.globals.evening = false

	// Test 1: ship not moving returns 0
	teg.globals.bx[1000] = &box{
		x_subloc: &entity_subloc{moving: 0},
	}
	if got := ship_gone(1000); got != 0 {
		t.Errorf("ship_gone(not moving) = %d, want 0", got)
	}

	// Test 2: ship moving, not evening
	teg.globals.bx[1001] = &box{
		x_subloc: &entity_subloc{moving: 95},
	}
	// Expected: 100 - 95 + 0 = 5
	if got := ship_gone(1001); got != 5 {
		t.Errorf("ship_gone(moving, not evening) = %d, want 5", got)
	}

	// Test 3: ship moving, is evening
	teg.globals.evening = true
	// Expected: 100 - 95 + 1 = 6
	if got := ship_gone(1001); got != 6 {
		t.Errorf("ship_gone(moving, is evening) = %d, want 6", got)
	}

	// Test 4: nil box returns 0
	if got := ship_gone(9999); got != 0 {
		t.Errorf("ship_gone(nil box) = %d, want 0", got)
	}
}

// TestCharMoving tests the char_moving function.
func TestCharMoving(t *testing.T) {
	// Save and restore global state
	oldBx := teg.globals.bx
	defer func() { teg.globals.bx = oldBx }()

	// Clear bx array for testing
	teg.globals.bx = [MAX_BOXES]*box{}

	// Test 1: nil box returns 0
	if got := char_moving(2000); got != 0 {
		t.Errorf("char_moving(nil box) = %d, want 0", got)
	}

	// Test 2: box with nil char returns 0
	teg.globals.bx[2001] = &box{}
	if got := char_moving(2001); got != 0 {
		t.Errorf("char_moving(nil char) = %d, want 0", got)
	}

	// Test 3: character not moving returns 0
	teg.globals.bx[2002] = &box{
		x_char: &entity_char{moving: 0},
	}
	if got := char_moving(2002); got != 0 {
		t.Errorf("char_moving(not moving) = %d, want 0", got)
	}

	// Test 4: character moving returns daystamp
	teg.globals.bx[2003] = &box{
		x_char: &entity_char{moving: 77},
	}
	if got := char_moving(2003); got != 77 {
		t.Errorf("char_moving(moving) = %d, want 77", got)
	}
}

// TestCharGone tests the char_gone function.
// Note: The C code uses a simplified version that returns 1 if moving, 0 if not.
func TestCharGone(t *testing.T) {
	// Save and restore global state
	oldBx := teg.globals.bx
	defer func() { teg.globals.bx = oldBx }()

	// Clear bx array for testing
	teg.globals.bx = [MAX_BOXES]*box{}

	// Test 1: character not moving returns 0
	teg.globals.bx[2000] = &box{
		x_char: &entity_char{moving: 0},
	}
	if got := char_gone(2000); got != 0 {
		t.Errorf("char_gone(not moving) = %d, want 0", got)
	}

	// Test 2: character moving returns 1 (simplified version)
	teg.globals.bx[2001] = &box{
		x_char: &entity_char{moving: 50},
	}
	if got := char_gone(2001); got != 1 {
		t.Errorf("char_gone(moving) = %d, want 1", got)
	}

	// Test 3: nil box returns 0
	if got := char_gone(9999); got != 0 {
		t.Errorf("char_gone(nil box) = %d, want 0", got)
	}
}

// TestCharGoneFull tests the char_gone_full function.
// This is the "full" version that calculates actual days elapsed.
func TestCharGoneFull(t *testing.T) {
	// Save and restore global state
	oldBx := teg.globals.bx
	oldSysclock := teg.globals.sysclock
	oldEvening := teg.globals.evening
	defer func() {
		teg.globals.bx = oldBx
		teg.globals.sysclock = oldSysclock
		teg.globals.evening = oldEvening
	}()

	// Clear state for testing
	teg.globals.bx = [MAX_BOXES]*box{}
	teg.globals.sysclock = olytime{days_since_epoch: 200}
	teg.globals.evening = false

	// Test 1: character not moving returns 0
	teg.globals.bx[2000] = &box{
		x_char: &entity_char{moving: 0},
	}
	if got := char_gone_full(2000); got != 0 {
		t.Errorf("char_gone_full(not moving) = %d, want 0", got)
	}

	// Test 2: character moving, not evening
	teg.globals.bx[2001] = &box{
		x_char: &entity_char{moving: 193},
	}
	// Expected: 200 - 193 + 0 = 7
	if got := char_gone_full(2001); got != 7 {
		t.Errorf("char_gone_full(moving, not evening) = %d, want 7", got)
	}

	// Test 3: character moving, is evening
	teg.globals.evening = true
	// Expected: 200 - 193 + 1 = 8
	if got := char_gone_full(2001); got != 8 {
		t.Errorf("char_gone_full(moving, is evening) = %d, want 8", got)
	}
}

// TestBoolToInt tests the boolToInt helper function.
func TestBoolToInt(t *testing.T) {
	if got := boolToInt(false); got != 0 {
		t.Errorf("boolToInt(false) = %d, want 0", got)
	}
	if got := boolToInt(true); got != 1 {
		t.Errorf("boolToInt(true) = %d, want 1", got)
	}
}
