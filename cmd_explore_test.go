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
	"math/rand/v2"
	"testing"

	"github.com/mdhender/prng"
)

// setupExploreTest initializes a test environment for exploration tests.
func setupExploreTest(t *testing.T) {
	t.Helper()

	teg = &Engine{
		prng: prng.New(rand.NewPCG(12345, 67890)),
	}
	teg.globals.inventories = make(map[int][]item_ent)
	teg.globals.playerKnowledge = make(map[int]map[int]bool)
}

func TestVExplore(t *testing.T) {
	setupExploreTest(t)

	c := &command{who: 100}
	result := v_explore(c)

	if result != TRUE {
		t.Errorf("v_explore returned %d, want TRUE (%d)", result, TRUE)
	}
}

func TestFindLostItems_NoItems(t *testing.T) {
	setupExploreTest(t)

	who := 100
	where := 200

	teg.globals.bx[who] = &box{kind: T_char}
	teg.globals.bx[where] = &box{kind: T_loc, skind: sub_plain}

	result := find_lost_items(who, where)

	if result {
		t.Error("find_lost_items returned true when no items present")
	}
}

func TestFindLostItems_NonUniqueItems(t *testing.T) {
	setupExploreTest(t)

	who := 100
	where := 200
	item := 300

	teg.globals.bx[who] = &box{kind: T_char}
	teg.globals.bx[where] = &box{kind: T_loc, skind: sub_plain}
	teg.globals.bx[item] = &box{kind: T_item}

	// Add a non-unique item (no who_has set)
	teg.globals.inventories[where] = []item_ent{{item: item, qty: 5}}

	result := find_lost_items(who, where)

	if result {
		t.Error("find_lost_items returned true for non-unique items")
	}
}

func TestFindLostItems_UniqueItemInSubloc(t *testing.T) {
	setupExploreTest(t)

	who := 100
	where := 200
	item := 300
	playerID := 50

	teg.globals.bx[who] = &box{kind: T_char}
	teg.globals.bx[where] = &box{kind: T_loc, skind: sub_city} // LOC_subloc depth = 100% chance
	teg.globals.bx[item] = &box{
		kind:   T_item,
		skind:  sub_scroll,
		x_item: &entity_item{who_has: where}, // Unique item owned by location
	}
	teg.globals.bx[playerID] = &box{kind: T_player}

	// Set character's location and player
	setupExploreLocation(where, sub_city)
	set_where(who, where)

	// Add unique item to location
	teg.globals.inventories[where] = []item_ent{{item: item, qty: 1}}

	result := find_lost_items(who, where)

	if !result {
		t.Error("find_lost_items returned false for unique item in subloc")
	}

	// Check item moved to character
	if has_item(who, item) != 1 {
		t.Errorf("character should have 1 of item, got %d", has_item(who, item))
	}
}

func TestFindLostItems_SkipsDeadBodiesInGraveyard(t *testing.T) {
	setupExploreTest(t)

	who := 100
	where := 200
	deadBody := 300

	teg.globals.bx[who] = &box{kind: T_char}
	teg.globals.bx[where] = &box{kind: T_loc, skind: sub_graveyard}
	teg.globals.bx[deadBody] = &box{
		kind:   T_item,
		skind:  sub_dead_body,
		x_item: &entity_item{who_has: where},
	}

	setupExploreLocation(where, sub_graveyard)
	set_where(who, where)

	teg.globals.inventories[where] = []item_ent{{item: deadBody, qty: 1}}

	result := find_lost_items(who, where)

	if result {
		t.Error("find_lost_items should skip dead bodies in graveyards")
	}
}

func TestFindLostItems_SkipsSuffuseRingsInCity(t *testing.T) {
	setupExploreTest(t)

	who := 100
	where := 200
	ring := 300

	teg.globals.bx[who] = &box{kind: T_char}
	teg.globals.bx[where] = &box{kind: T_loc, skind: sub_city}
	teg.globals.bx[ring] = &box{
		kind:   T_item,
		skind:  sub_suffuse_ring,
		x_item: &entity_item{who_has: where},
	}

	setupExploreLocation(where, sub_city)
	set_where(who, where)

	teg.globals.inventories[where] = []item_ent{{item: ring, qty: 1}}

	result := find_lost_items(who, where)

	if result {
		t.Error("find_lost_items should skip suffuse rings in cities")
	}
}

func TestDExplore_NoFeaturesMessage(t *testing.T) {
	setupExploreTest(t)

	who := 100
	where := 200

	teg.globals.bx[who] = &box{kind: T_char}
	teg.globals.bx[where] = &box{kind: T_loc, skind: sub_plain}

	setupExploreLocation(where, sub_plain)
	set_where(who, where)

	c := &command{who: who}

	// With no items and no hidden exits, d_explore should return FALSE
	// (at least 50% of the time due to RNG, but with no hidden exits always FALSE)
	result := d_explore(c)

	if result != FALSE {
		t.Errorf("d_explore returned %d, expected FALSE (%d) when no features", result, FALSE)
	}
}

func TestCountHiddenExits_NilList(t *testing.T) {
	setupExploreTest(t)

	result := count_hidden_exits(nil)

	if result != 0 {
		t.Errorf("count_hidden_exits(nil) = %d, want 0", result)
	}
}

func TestCountHiddenExits_NoHidden(t *testing.T) {
	setupExploreTest(t)

	exits := []*exit_view{
		{destination: 100, hidden: 0},
		{destination: 200, hidden: 0},
	}

	result := count_hidden_exits(exits)

	if result != 0 {
		t.Errorf("count_hidden_exits = %d, want 0", result)
	}
}

func TestCountHiddenExits_SomeHidden(t *testing.T) {
	setupExploreTest(t)

	exits := []*exit_view{
		{destination: 100, hidden: 0},
		{destination: 200, hidden: 1},
		{destination: 300, hidden: 1},
		{destination: 400, hidden: 0},
	}

	result := count_hidden_exits(exits)

	if result != 2 {
		t.Errorf("count_hidden_exits = %d, want 2", result)
	}
}

func TestHiddenCountToIndex(t *testing.T) {
	setupExploreTest(t)

	exits := []*exit_view{
		{destination: 100, hidden: 0},
		{destination: 200, hidden: 1}, // index 1, 1st hidden
		{destination: 300, hidden: 0},
		{destination: 400, hidden: 1}, // index 3, 2nd hidden
		{destination: 500, hidden: 1}, // index 4, 3rd hidden
	}

	tests := []struct {
		which int
		want  int
	}{
		{1, 1}, // 1st hidden is at index 1
		{2, 3}, // 2nd hidden is at index 3
		{3, 4}, // 3rd hidden is at index 4
	}

	for _, tt := range tests {
		got := hidden_count_to_index(tt.which, exits)
		if got != tt.want {
			t.Errorf("hidden_count_to_index(%d) = %d, want %d", tt.which, got, tt.want)
		}
	}
}

func TestHiddenCountToIndex_NilList(t *testing.T) {
	setupExploreTest(t)

	result := hidden_count_to_index(1, nil)

	if result != 0 {
		t.Errorf("hidden_count_to_index with nil list = %d, want 0", result)
	}
}

func TestFindHiddenExit(t *testing.T) {
	setupExploreTest(t)

	who := 100
	dest := 500
	playerID := 50

	teg.globals.bx[who] = &box{kind: T_char, x_char: &entity_char{}}
	teg.globals.bx[dest] = &box{kind: T_loc, skind: sub_cave}
	teg.globals.bx[playerID] = &box{kind: T_player}

	// Set up player relationship properly: unit_lord points to player
	teg.globals.bx[who].x_char.unit_lord = playerID

	exits := []*exit_view{
		{destination: dest, hidden: 1},
	}

	find_hidden_exit(who, exits, 0)

	// Check that destination is now known
	if teg.globals.playerKnowledge[playerID] == nil {
		t.Fatal("playerKnowledge[playerID] is nil after find_hidden_exit")
	}
	if !teg.globals.playerKnowledge[playerID][dest] {
		t.Error("destination should be marked as known after finding hidden exit")
	}
}

func TestSetKnown(t *testing.T) {
	setupExploreTest(t)

	who := 100
	what := 500
	playerID := 50

	teg.globals.bx[who] = &box{kind: T_char, x_char: &entity_char{}}
	teg.globals.bx[what] = &box{kind: T_loc}
	teg.globals.bx[playerID] = &box{kind: T_player}

	// Set up player relationship properly: unit_lord points to player
	teg.globals.bx[who].x_char.unit_lord = playerID

	set_known(who, what)

	if teg.globals.playerKnowledge[playerID] == nil {
		t.Fatal("playerKnowledge[playerID] is nil after set_known")
	}

	if !teg.globals.playerKnowledge[playerID][what] {
		t.Error("set_known did not mark entity as known")
	}
}

func TestDExplore_ShipInOcean(t *testing.T) {
	setupExploreTest(t)

	who := 100
	ship := 200
	ocean := 300

	teg.globals.bx[who] = &box{kind: T_char}
	teg.globals.bx[ship] = &box{kind: T_loc, skind: sub_galley}
	teg.globals.bx[ocean] = &box{kind: T_loc, skind: sub_ocean}

	// Set up ship in ocean
	setupExploreLocation(ocean, sub_ocean)
	setupExploreLocation(ship, sub_galley)
	teg.globals.bx[ship].x_loc_info.where = ocean

	set_where(who, ship)

	c := &command{who: who}

	// Should explore ocean, not ship
	result := d_explore(c)

	// With no items or hidden exits, should fail
	if result != FALSE {
		t.Errorf("d_explore in ship on ocean = %d, want FALSE", result)
	}
}

// setupExploreLocation sets up a test location for exploration tests.
// Uses a different name to avoid conflict with lifecycle_test.go helper.
func setupExploreLocation(id int, sk schar) {
	if teg.globals.bx[id] == nil {
		teg.globals.bx[id] = &box{}
	}
	teg.globals.bx[id].kind = T_loc
	teg.globals.bx[id].skind = sk
	if teg.globals.bx[id].x_loc_info.where == 0 {
		teg.globals.bx[id].x_loc_info.where = id
	}
}
