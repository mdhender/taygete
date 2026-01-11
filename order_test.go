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

func TestOrderQueue(t *testing.T) {
	e := teg

	// Create a test player and unit
	playerID := 1001
	unitID := 2001

	// Ensure the player box exists
	e.globals.bx[playerID] = &box{
		kind:     T_player,
		skind:    sub_pl_regular,
		x_player: &entity_player{},
	}

	// Ensure the unit box exists
	e.globals.bx[unitID] = &box{
		kind:   T_char,
		skind:  0,
		x_char: &entity_char{unit_lord: playerID},
	}

	// Clear any existing orders
	e.ClearOrders()

	// Test empty queue
	if got := e.top_order(playerID, unitID); got != "" {
		t.Errorf("top_order on empty queue: got %q, want %q", got, "")
	}

	// Test queue_order
	e.queue_order(playerID, unitID, "move north")
	e.queue_order(playerID, unitID, "wait 5")
	e.queue_order(playerID, unitID, "attack 1234")

	if got := e.CountOrders(playerID, unitID); got != 3 {
		t.Errorf("CountOrders: got %d, want %d", got, 3)
	}

	// Test top_order
	if got := e.top_order(playerID, unitID); got != "move north" {
		t.Errorf("top_order: got %q, want %q", got, "move north")
	}

	// Test pop_order
	e.pop_order(playerID, unitID)
	if got := e.top_order(playerID, unitID); got != "wait 5" {
		t.Errorf("top_order after pop: got %q, want %q", got, "wait 5")
	}

	// Test prepend_order
	e.prepend_order(playerID, unitID, "stop")
	if got := e.top_order(playerID, unitID); got != "stop" {
		t.Errorf("top_order after prepend: got %q, want %q", got, "stop")
	}

	// Test GetAllOrders
	orders := e.GetAllOrders(playerID, unitID)
	expected := []string{"stop", "wait 5", "attack 1234"}
	if len(orders) != len(expected) {
		t.Errorf("GetAllOrders length: got %d, want %d", len(orders), len(expected))
	} else {
		for i, o := range orders {
			if o != expected[i] {
				t.Errorf("GetAllOrders[%d]: got %q, want %q", i, o, expected[i])
			}
		}
	}

	// Cleanup
	e.globals.bx[playerID] = nil
	e.globals.bx[unitID] = nil
	e.ClearOrders()
}

func TestStopOrder(t *testing.T) {
	e := teg

	playerID := 1002
	unitID := 2002

	e.globals.bx[playerID] = &box{
		kind:     T_player,
		skind:    sub_pl_regular,
		x_player: &entity_player{},
	}
	e.globals.bx[unitID] = &box{
		kind:   T_char,
		skind:  0,
		x_char: &entity_char{unit_lord: playerID},
	}

	e.ClearOrders()

	// Test stop_order when no orders
	if e.stop_order(playerID, unitID) {
		t.Error("stop_order on empty queue should return false")
	}

	// Test with non-stop order
	e.queue_order(playerID, unitID, "move north")
	if e.stop_order(playerID, unitID) {
		t.Error("stop_order with 'move north' should return false")
	}

	// Test queue_stop
	e.queue_stop(playerID, unitID)
	if !e.stop_order(playerID, unitID) {
		t.Error("stop_order after queue_stop should return true")
	}

	// Test that queue_stop doesn't add duplicate
	e.queue_stop(playerID, unitID)
	if e.CountOrders(playerID, unitID) != 2 {
		t.Errorf("queue_stop should not add duplicate: got %d orders, want 2",
			e.CountOrders(playerID, unitID))
	}

	// Cleanup
	e.globals.bx[playerID] = nil
	e.globals.bx[unitID] = nil
	e.ClearOrders()
}

func TestIsStopOrder(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"stop", true},
		{"STOP", true},
		{"Stop", true},
		{"  stop", true},
		{"\tstop", true},
		{"stop now", true},
		{"stopper", false},
		{"move", false},
		{"", false},
		{"attack stop", false},
	}

	for _, tc := range tests {
		got := is_stop_order(tc.input)
		if got != tc.want {
			t.Errorf("is_stop_order(%q): got %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestFlushUnitOrders(t *testing.T) {
	e := teg

	playerID := 1003
	unitID := 2003

	e.globals.bx[playerID] = &box{
		kind:     T_player,
		skind:    sub_pl_regular,
		x_player: &entity_player{},
	}
	e.globals.bx[unitID] = &box{
		kind:   T_char,
		skind:  0,
		x_char: &entity_char{unit_lord: playerID},
	}

	e.ClearOrders()

	// Queue some orders
	e.queue_order(playerID, unitID, "move north")
	e.queue_order(playerID, unitID, "wait 5")
	e.queue_order(playerID, unitID, "attack 1234")

	if e.CountOrders(playerID, unitID) != 3 {
		t.Errorf("before flush: got %d orders, want 3", e.CountOrders(playerID, unitID))
	}

	// Flush orders
	e.flush_unit_orders(playerID, unitID)

	if e.CountOrders(playerID, unitID) != 0 {
		t.Errorf("after flush: got %d orders, want 0", e.CountOrders(playerID, unitID))
	}

	// Cleanup
	e.globals.bx[playerID] = nil
	e.globals.bx[unitID] = nil
	e.ClearOrders()
}

func TestQueueLimit(t *testing.T) {
	e := teg

	playerID := 1004
	unitID := 2004

	e.globals.bx[playerID] = &box{
		kind:     T_player,
		skind:    sub_pl_regular,
		x_player: &entity_player{},
	}
	e.globals.bx[unitID] = &box{
		kind:   T_char,
		skind:  0,
		x_char: &entity_char{unit_lord: playerID},
	}

	e.ClearOrders()

	// Queue 260 orders - should be limited to 250
	for i := 0; i < 260; i++ {
		e.queue_order(playerID, unitID, "wait 1")
	}

	if got := e.CountOrders(playerID, unitID); got != 250 {
		t.Errorf("queue limit: got %d orders, want 250", got)
	}

	// Cleanup
	e.globals.bx[playerID] = nil
	e.globals.bx[unitID] = nil
	e.ClearOrders()
}

func TestQueueConvenienceFunction(t *testing.T) {
	e := teg

	playerID := 1005
	unitID := 2005

	e.globals.bx[playerID] = &box{
		kind:     T_player,
		skind:    sub_pl_regular,
		x_player: &entity_player{},
	}
	e.globals.bx[unitID] = &box{
		kind:   T_char,
		skind:  0,
		x_char: &entity_char{unit_lord: playerID},
	}

	e.ClearOrders()

	// Test queue() convenience function with formatting
	e.queue(unitID, "move %s", "north")
	e.queue(unitID, "wait %d", 5)
	e.queue(unitID, "attack %d", 1234)

	orders := e.GetAllOrders(playerID, unitID)
	expected := []string{"move north", "wait 5", "attack 1234"}

	if len(orders) != len(expected) {
		t.Errorf("queue() length: got %d, want %d", len(orders), len(expected))
	} else {
		for i, o := range orders {
			if o != expected[i] {
				t.Errorf("queue()[%d]: got %q, want %q", i, o, expected[i])
			}
		}
	}

	// Cleanup
	e.globals.bx[playerID] = nil
	e.globals.bx[unitID] = nil
	e.ClearOrders()
}

func TestLoadSaveOrders(t *testing.T) {
	// Use isolated database to avoid interference with other tests
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	e := &Engine{db: db}

	playerID := 1006
	unitID := 2006
	turnNumber := 99

	// Ensure entities exist in the bx array
	e.globals.bx[playerID] = &box{
		kind:     T_player,
		skind:    sub_pl_regular,
		x_player: &entity_player{},
	}
	e.globals.bx[unitID] = &box{
		kind:   T_char,
		skind:  0,
		x_char: &entity_char{unit_lord: playerID},
	}

	// Create player in database if needed
	_, _ = e.db.Exec(`INSERT OR IGNORE INTO players (id, code, subkind) VALUES (?, ?, ?)`,
		playerID, "aa01", sub_pl_regular)

	// Create entity for character (required for foreign key)
	_, _ = e.db.Exec(`INSERT OR IGNORE INTO entities (id, kind, subkind) VALUES (?, ?, ?)`,
		unitID, T_char, 0)

	// Create character in database (required for foreign key in orders table)
	_, _ = e.db.Exec(`INSERT OR IGNORE INTO characters (id, player_id) VALUES (?, ?)`,
		unitID, playerID)

	// Create a turn record for foreign key constraint
	_, _ = e.db.Exec(`INSERT OR IGNORE INTO turns (turn_number, status) VALUES (?, 'pending')`,
		turnNumber)

	e.ClearOrders()

	// Queue some orders
	e.queue_order(playerID, unitID, "move north")
	e.queue_order(playerID, unitID, "wait 5")

	// Save orders
	if err := e.SaveOrders(turnNumber); err != nil {
		t.Fatalf("SaveOrders: %v", err)
	}

	// Clear in-memory orders
	e.ClearOrders()

	// Verify orders are gone from memory
	if e.CountOrders(playerID, unitID) != 0 {
		t.Error("orders should be cleared from memory")
	}

	// Load orders back
	if err := e.LoadOrders(turnNumber); err != nil {
		t.Fatalf("LoadOrders: %v", err)
	}

	// Verify orders are restored
	if e.CountOrders(playerID, unitID) != 2 {
		t.Errorf("after load: got %d orders, want 2", e.CountOrders(playerID, unitID))
	}

	orders := e.GetAllOrders(playerID, unitID)
	if len(orders) >= 1 && orders[0] != "move north" {
		t.Errorf("order[0]: got %q, want %q", orders[0], "move north")
	}
	if len(orders) >= 2 && orders[1] != "wait 5" {
		t.Errorf("order[1]: got %q, want %q", orders[1], "wait 5")
	}

	// Cleanup memory (database cleanup happens via defer db.Close())
	e.globals.bx[playerID] = nil
	e.globals.bx[unitID] = nil
	e.ClearOrders()
}
