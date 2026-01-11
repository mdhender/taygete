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

// move_test.go - Unit tests for move.go
// Sprint 27: Movement & World

package taygete

import "testing"

func TestExitOppositeDirections(t *testing.T) {
	tests := []struct {
		dir      int
		expected int
	}{
		{DIR_N, DIR_S},
		{DIR_S, DIR_N},
		{DIR_E, DIR_W},
		{DIR_W, DIR_E},
	}

	for _, tt := range tests {
		if exit_opposite[tt.dir] != tt.expected {
			t.Errorf("exit_opposite[%d] = %d, want %d", tt.dir, exit_opposite[tt.dir], tt.expected)
		}
	}
}

func TestSaveRestoreVArray(t *testing.T) {
	c := &command{}
	v := &exit_view{
		direction:   DIR_N,
		destination: 12345,
		road:        100,
		dest_hidden: 1,
		distance:    5,
		orig:        54321,
		orig_hidden: 0,
	}

	save_v_array(c, v)

	var restored exit_view
	restore_v_array(c, &restored)

	if restored.direction != v.direction {
		t.Errorf("direction: got %d, want %d", restored.direction, v.direction)
	}
	if restored.destination != v.destination {
		t.Errorf("destination: got %d, want %d", restored.destination, v.destination)
	}
	if restored.road != v.road {
		t.Errorf("road: got %d, want %d", restored.road, v.road)
	}
	if restored.dest_hidden != v.dest_hidden {
		t.Errorf("dest_hidden: got %d, want %d", restored.dest_hidden, v.dest_hidden)
	}
	if restored.distance != v.distance {
		t.Errorf("distance: got %d, want %d", restored.distance, v.distance)
	}
	if restored.orig != v.orig {
		t.Errorf("orig: got %d, want %d", restored.orig, v.orig)
	}
	if restored.orig_hidden != v.orig_hidden {
		t.Errorf("orig_hidden: got %d, want %d", restored.orig_hidden, v.orig_hidden)
	}
}

func TestPluralMan(t *testing.T) {
	tests := []struct {
		n        int
		expected string
	}{
		{0, "men"},
		{1, "man"},
		{2, "men"},
		{100, "men"},
	}

	for _, tt := range tests {
		result := plural_man(tt.n)
		if result != tt.expected {
			t.Errorf("plural_man(%d) = %q, want %q", tt.n, result, tt.expected)
		}
	}
}

func TestOceanCharsInit(t *testing.T) {
	oldTeg := teg
	defer func() { teg = oldTeg }()

	teg = &Engine{}
	teg.globals.bx = [MAX_BOXES]*box{}

	ocean_chars = []int{1, 2, 3}
	init_ocean_chars()

	// After init, ocean_chars is reset to empty slice (or nil)
	if len(ocean_chars) != 0 {
		t.Errorf("ocean_chars should be empty, got %d elements", len(ocean_chars))
	}
}

func TestLandCheckWater(t *testing.T) {
	c := &command{who: 1}
	v := &exit_view{water: 1}

	result := land_check(c, v, false)
	if result {
		t.Error("land_check should return false for water routes")
	}
}

func TestLandCheckImpassable(t *testing.T) {
	c := &command{who: 1}
	v := &exit_view{impassable: 1}

	result := land_check(c, v, false)
	if result {
		t.Error("land_check should return false for impassable routes")
	}
}

func TestLandCheckInTransit(t *testing.T) {
	c := &command{who: 1}
	v := &exit_view{in_transit: 1}

	result := land_check(c, v, false)
	if result {
		t.Error("land_check should return false for in-transit destinations")
	}
}

func TestLandCheckSuccess(t *testing.T) {
	c := &command{who: 1}
	v := &exit_view{
		water:      0,
		impassable: 0,
		in_transit: 0,
	}

	result := land_check(c, v, false)
	if !result {
		t.Error("land_check should return true for valid land routes")
	}
}

func TestSailCheckNoWater(t *testing.T) {
	c := &command{who: 1}
	v := &exit_view{water: 0}

	result := sail_check(c, v, false)
	if result {
		t.Error("sail_check should return false for non-water routes")
	}
}

func TestSailCheckImpassable(t *testing.T) {
	c := &command{who: 1}
	v := &exit_view{water: 1, impassable: 1}

	result := sail_check(c, v, false)
	if result {
		t.Error("sail_check should return false for impassable water routes")
	}
}

func TestSailCheckSuccess(t *testing.T) {
	c := &command{who: 1}
	v := &exit_view{water: 1, impassable: 0}

	result := sail_check(c, v, true)
	if !result {
		t.Error("sail_check should return true for valid water routes")
	}
}

func TestFlyCheckInTransit(t *testing.T) {
	c := &command{who: 1}
	v := &exit_view{in_transit: 1}

	result := fly_check(c, v)
	if result {
		t.Error("fly_check should return false for in-transit destinations")
	}
}

func TestFlyCheckSuccess(t *testing.T) {
	c := &command{who: 1}
	v := &exit_view{in_transit: 0}

	result := fly_check(c, v)
	if !result {
		t.Error("fly_check should return true for valid destinations")
	}
}

func TestCountStackMoveNobles(t *testing.T) {
	oldTeg := teg
	defer func() { teg = oldTeg }()

	teg = &Engine{}
	teg.globals.bx = [MAX_BOXES]*box{}

	leaderID := 1000
	followerID := 1001
	playerID := 100

	teg.globals.bx[leaderID] = &box{
		kind:   T_char,
		x_char: &entity_char{unit_lord: playerID},
	}
	teg.globals.bx[leaderID].x_loc_info.here_list = []int{followerID}

	teg.globals.bx[followerID] = &box{
		kind:   T_char,
		x_char: &entity_char{unit_lord: playerID},
	}
	teg.globals.bx[followerID].x_loc_info.where = leaderID

	teg.globals.bx[playerID] = &box{
		kind: T_player,
	}

	count := count_stack_move_nobles(leaderID)
	if count < 1 {
		t.Errorf("count_stack_move_nobles should return at least 1, got %d", count)
	}
}

func TestMoveAlive(t *testing.T) {
	oldTeg := teg
	defer func() { teg = oldTeg }()

	teg = &Engine{}
	teg.globals.bx = [MAX_BOXES]*box{}

	charID := 1000
	teg.globals.bx[charID] = &box{
		kind:   T_char,
		x_char: &entity_char{health: 100},
	}

	// The alive() function in accessor.go just checks kind == T_char
	if !alive(charID) {
		t.Error("alive should return true for character")
	}

	locID := 2000
	teg.globals.bx[locID] = &box{
		kind: T_loc,
	}

	if alive(locID) {
		t.Error("alive should return false for non-character")
	}
}

func TestLoopStack(t *testing.T) {
	oldTeg := teg
	defer func() { teg = oldTeg }()

	teg = &Engine{}
	teg.globals.bx = [MAX_BOXES]*box{}

	leaderID := 1000
	follower1ID := 1001
	follower2ID := 1002
	nestedID := 1003

	teg.globals.bx[leaderID] = &box{
		kind:   T_char,
		x_char: &entity_char{},
	}
	teg.globals.bx[leaderID].x_loc_info.here_list = []int{follower1ID, follower2ID}

	teg.globals.bx[follower1ID] = &box{
		kind:   T_char,
		x_char: &entity_char{},
	}
	teg.globals.bx[follower1ID].x_loc_info.where = leaderID
	teg.globals.bx[follower1ID].x_loc_info.here_list = []int{nestedID}

	teg.globals.bx[follower2ID] = &box{
		kind:   T_char,
		x_char: &entity_char{},
	}
	teg.globals.bx[follower2ID].x_loc_info.where = leaderID

	teg.globals.bx[nestedID] = &box{
		kind:   T_char,
		x_char: &entity_char{},
	}
	teg.globals.bx[nestedID].x_loc_info.where = follower1ID

	var stack []int
	loop_stack(leaderID, &stack)

	if len(stack) != 4 {
		t.Errorf("loop_stack should return 4 members, got %d", len(stack))
	}

	found := make(map[int]bool)
	for _, id := range stack {
		found[id] = true
	}

	if !found[leaderID] {
		t.Error("stack should include leader")
	}
	if !found[follower1ID] {
		t.Error("stack should include follower1")
	}
	if !found[follower2ID] {
		t.Error("stack should include follower2")
	}
	if !found[nestedID] {
		t.Error("stack should include nested follower")
	}
}

func TestCountAny(t *testing.T) {
	oldTeg := teg
	defer func() { teg = oldTeg }()

	teg = &Engine{}
	teg.globals.bx = [MAX_BOXES]*box{}
	teg.globals.inventories = make(map[int][]item_ent)

	charID := 1000
	teg.globals.bx[charID] = &box{
		kind:   T_char,
		x_char: &entity_char{},
	}

	count := count_any(charID)
	if count != 1 {
		t.Errorf("count_any for char with no items should be 1, got %d", count)
	}
}

func TestIsManItem(t *testing.T) {
	oldTeg := teg
	defer func() { teg = oldTeg }()

	teg = &Engine{}
	teg.globals.bx = [MAX_BOXES]*box{}

	manItemID := 10
	teg.globals.bx[manItemID] = &box{
		kind:   T_item,
		x_item: &entity_item{is_man_item: 1},
	}

	if !is_man_item(manItemID) {
		t.Error("is_man_item should return true for man items")
	}

	nonManItemID := 20
	teg.globals.bx[nonManItemID] = &box{
		kind:   T_item,
		x_item: &entity_item{is_man_item: 0},
	}

	if is_man_item(nonManItemID) {
		t.Error("is_man_item should return false for non-man items")
	}
}
