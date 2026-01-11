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

// cmd_economy_test.go - Unit tests for Sprint 25.8 economy commands

package taygete

import (
	"math/rand/v2"
	"testing"

	"github.com/mdhender/prng"
)

func setupEconomyTest() {
	if teg.prng == nil {
		teg.prng = prng.New(rand.NewPCG(12345, 67890))
	}
	if teg.globals.inventories == nil {
		teg.globals.inventories = make(map[int][]item_ent)
	}
	if teg.globals.names == nil {
		teg.globals.names = make(map[int]string)
	}
	if teg.globals.banners == nil {
		teg.globals.banners = make(map[int]string)
	}

	for id := range teg.globals.inventories {
		delete(teg.globals.inventories, id)
	}
	for id := range teg.globals.names {
		delete(teg.globals.names, id)
	}
	for id := range teg.globals.banners {
		delete(teg.globals.banners, id)
	}

	for i := 0; i < MAX_BOXES; i++ {
		teg.globals.bx[i] = nil
	}
	for i := range teg.globals.box_head {
		teg.globals.box_head[i] = 0
	}
	for i := range teg.globals.sub_head {
		teg.globals.sub_head[i] = 0
	}

	alloc_box(item_gold, T_item, 0)
}

func TestHowMany(t *testing.T) {
	setupEconomyTest()

	charID := 1001
	testItem := 2001

	alloc_box(charID, T_char, 0)
	alloc_box(testItem, T_item, 0)

	t.Run("returns 0 when no items held", func(t *testing.T) {
		qty := how_many(charID, charID, testItem, 10, 0)
		if qty != 0 {
			t.Errorf("how_many with no items = %d, want 0", qty)
		}
	})

	t.Run("returns all when qty is 0", func(t *testing.T) {
		teg.globals.inventories[charID] = []item_ent{{item: testItem, qty: 50}}
		qty := how_many(charID, charID, testItem, 0, 0)
		if qty != 50 {
			t.Errorf("how_many with qty=0 = %d, want 50", qty)
		}
	})

	t.Run("returns requested qty when available", func(t *testing.T) {
		teg.globals.inventories[charID] = []item_ent{{item: testItem, qty: 50}}
		qty := how_many(charID, charID, testItem, 20, 0)
		if qty != 20 {
			t.Errorf("how_many with qty=20 = %d, want 20", qty)
		}
	})

	t.Run("respects have_left parameter", func(t *testing.T) {
		teg.globals.inventories[charID] = []item_ent{{item: testItem, qty: 50}}
		qty := how_many(charID, charID, testItem, 0, 10)
		if qty != 40 {
			t.Errorf("how_many with have_left=10 = %d, want 40", qty)
		}
	})

	t.Run("returns 0 when have_left >= num_has", func(t *testing.T) {
		teg.globals.inventories[charID] = []item_ent{{item: testItem, qty: 10}}
		qty := how_many(charID, charID, testItem, 5, 10)
		if qty != 0 {
			t.Errorf("how_many with have_left >= num_has = %d, want 0", qty)
		}
	})

	t.Run("caps qty at available minus have_left", func(t *testing.T) {
		teg.globals.inventories[charID] = []item_ent{{item: testItem, qty: 30}}
		qty := how_many(charID, charID, testItem, 100, 10)
		if qty != 20 {
			t.Errorf("how_many with excessive qty = %d, want 20", qty)
		}
	})
}

func TestVDiscard(t *testing.T) {
	setupEconomyTest()

	playerID := 100
	charID := 1001
	testItem := 2001
	locID := 5000

	alloc_box(playerID, T_player, 0)
	alloc_box(charID, T_char, 0)
	alloc_box(testItem, T_item, 0)
	alloc_box(locID, T_loc, sub_plain)

	p_char(charID).unit_lord = playerID
	set_where(charID, locID)

	t.Run("fails with invalid item", func(t *testing.T) {
		c := &command{who: charID, a: 9999}
		result := v_discard(c)
		if result != FALSE {
			t.Errorf("v_discard with invalid item = %d, want FALSE", result)
		}
	})

	t.Run("fails with non-item entity", func(t *testing.T) {
		c := &command{who: charID, a: locID}
		result := v_discard(c)
		if result != FALSE {
			t.Errorf("v_discard with non-item = %d, want FALSE", result)
		}
	})

	t.Run("fails when character has no items", func(t *testing.T) {
		c := &command{who: charID, a: testItem}
		result := v_discard(c)
		if result != FALSE {
			t.Errorf("v_discard with no items = %d, want FALSE", result)
		}
	})

	t.Run("succeeds and drops item", func(t *testing.T) {
		teg.globals.inventories[charID] = []item_ent{{item: testItem, qty: 10}}

		c := &command{who: charID, a: testItem, b: 5}
		result := v_discard(c)
		if result != TRUE {
			t.Errorf("v_discard = %d, want TRUE", result)
		}

		remaining := has_item(charID, testItem)
		if remaining != 5 {
			t.Errorf("remaining items = %d, want 5", remaining)
		}
	})

	t.Run("drops all when qty is 0", func(t *testing.T) {
		teg.globals.inventories[charID] = []item_ent{{item: testItem, qty: 10}}

		c := &command{who: charID, a: testItem, b: 0}
		result := v_discard(c)
		if result != TRUE {
			t.Errorf("v_discard all = %d, want TRUE", result)
		}

		remaining := has_item(charID, testItem)
		if remaining != 0 {
			t.Errorf("remaining items after drop all = %d, want 0", remaining)
		}
	})

	t.Run("respects have_left parameter", func(t *testing.T) {
		teg.globals.inventories[charID] = []item_ent{{item: testItem, qty: 20}}

		c := &command{who: charID, a: testItem, b: 0, c: 5}
		result := v_discard(c)
		if result != TRUE {
			t.Errorf("v_discard with have_left = %d, want TRUE", result)
		}

		remaining := has_item(charID, testItem)
		if remaining != 5 {
			t.Errorf("remaining with have_left = %d, want 5", remaining)
		}
	})
}

func TestLoopUnits(t *testing.T) {
	setupEconomyTest()

	playerID := 100
	char1 := 1001
	char2 := 1002
	char3 := 1003
	otherPlayerID := 101
	otherChar := 1010

	alloc_box(playerID, T_player, 0)
	alloc_box(char1, T_char, 0)
	alloc_box(char2, T_char, 0)
	alloc_box(char3, T_char, 0)
	alloc_box(otherPlayerID, T_player, 0)
	alloc_box(otherChar, T_char, 0)

	p_char(char1).unit_lord = playerID
	p_char(char2).unit_lord = playerID
	p_char(char3).unit_lord = playerID
	p_char(otherChar).unit_lord = otherPlayerID

	t.Run("returns all units for player", func(t *testing.T) {
		units := loop_units(playerID)
		if len(units) != 3 {
			t.Errorf("loop_units returned %d units, want 3", len(units))
		}

		found := make(map[int]bool)
		for _, u := range units {
			found[u] = true
		}
		if !found[char1] || !found[char2] || !found[char3] {
			t.Error("loop_units did not return all player units")
		}
		if found[otherChar] {
			t.Error("loop_units returned other player's unit")
		}
	})

	t.Run("returns empty for player with no units", func(t *testing.T) {
		emptyPlayerID := 102
		alloc_box(emptyPlayerID, T_player, 0)
		units := loop_units(emptyPlayerID)
		if len(units) != 0 {
			t.Errorf("loop_units for empty player = %d, want 0", len(units))
		}
	})
}

func TestLoopDeadBody(t *testing.T) {
	setupEconomyTest()

	body1 := 3001
	body2 := 3002
	regularItem := 3003

	alloc_box(body1, T_item, sub_dead_body)
	alloc_box(body2, T_item, sub_dead_body)
	alloc_box(regularItem, T_item, 0)

	t.Run("returns only dead bodies", func(t *testing.T) {
		bodies := loop_dead_body()
		if len(bodies) != 2 {
			t.Errorf("loop_dead_body returned %d bodies, want 2", len(bodies))
		}

		found := make(map[int]bool)
		for _, b := range bodies {
			found[b] = true
		}
		if !found[body1] || !found[body2] {
			t.Error("loop_dead_body did not return all bodies")
		}
		if found[regularItem] {
			t.Error("loop_dead_body returned non-body item")
		}
	})
}

func TestVQuit(t *testing.T) {
	setupEconomyTest()

	playerID := 100
	charID := 1001
	otherPlayerID := 101
	locID := 5000

	alloc_box(playerID, T_player, 0)
	alloc_box(charID, T_char, 0)
	alloc_box(otherPlayerID, T_player, 0)
	alloc_box(locID, T_loc, sub_plain)

	p_char(charID).unit_lord = playerID
	set_where(charID, locID)

	t.Run("fails when non-GM tries to quit another player", func(t *testing.T) {
		c := &command{who: charID, a: otherPlayerID}
		result := v_quit(c)
		if result != FALSE {
			t.Errorf("v_quit other player = %d, want FALSE", result)
		}
	})

	t.Run("fails with invalid target", func(t *testing.T) {
		c := &command{who: charID, a: locID}
		result := v_quit(c)
		if result != FALSE {
			t.Errorf("v_quit non-player = %d, want FALSE", result)
		}
	})

	t.Run("succeeds and deletes player", func(t *testing.T) {
		quitPlayerID := 103
		quitCharID := 1010

		alloc_box(quitPlayerID, T_player, 0)
		alloc_box(quitCharID, T_char, 0)
		p_char(quitCharID).unit_lord = quitPlayerID
		set_where(quitCharID, locID)

		c := &command{who: quitCharID, a: quitPlayerID}
		result := v_quit(c)
		if result != FALSE {
			t.Errorf("v_quit = %d, want FALSE (special return)", result)
		}

		if valid_box(quitPlayerID) {
			t.Error("player should be deleted after quit")
		}
	})

	t.Run("defaults target to player when target is 0", func(t *testing.T) {
		defaultPlayerID := 104
		defaultCharID := 1020

		alloc_box(defaultPlayerID, T_player, 0)
		alloc_box(defaultCharID, T_char, 0)
		p_char(defaultCharID).unit_lord = defaultPlayerID
		set_where(defaultCharID, locID)

		c := &command{who: defaultCharID, a: 0}
		result := v_quit(c)
		if result != FALSE {
			t.Errorf("v_quit with default target = %d, want FALSE", result)
		}

		if valid_box(defaultPlayerID) {
			t.Error("player should be deleted when target defaults")
		}
	})
}

func TestDropPlayer(t *testing.T) {
	setupEconomyTest()

	playerID := 100
	char1 := 1001
	char2 := 1002
	locID := 5000

	alloc_box(playerID, T_player, 0)
	alloc_box(char1, T_char, 0)
	alloc_box(char2, T_char, 0)
	alloc_box(locID, T_loc, sub_plain)

	p_char(char1).unit_lord = playerID
	p_char(char2).unit_lord = playerID
	set_where(char1, locID)
	set_where(char2, locID)

	t.Run("removes all player units", func(t *testing.T) {
		drop_player(playerID)

		if valid_box(playerID) {
			t.Error("player should be deleted")
		}
		// Characters are converted to T_deadchar via char_reclaim/kill_char, not deleted
		if kind(char1) != T_deadchar {
			t.Errorf("char1 should be T_deadchar, got kind %d", kind(char1))
		}
		if kind(char2) != T_deadchar {
			t.Errorf("char2 should be T_deadchar, got kind %d", kind(char2))
		}
	})
}

func TestCharReclaim(t *testing.T) {
	setupEconomyTest()

	charID := 1001
	locID := 5000

	alloc_box(charID, T_char, 0)
	alloc_box(locID, T_loc, sub_plain)
	set_where(charID, locID)

	p := p_loc_info(locID)
	p.here_list = append(p.here_list, charID)
	p_char(charID).unit_lord = indep_player

	t.Run("marks character for melting and converts to deadchar", func(t *testing.T) {
		char_reclaim(charID)

		// char_reclaim now calls kill_char, which converts to T_deadchar (melt_me = TRUE)
		if kind(charID) != T_deadchar {
			t.Errorf("character should be converted to T_deadchar, got kind %d", kind(charID))
		}
	})

	t.Run("handles non-character gracefully", func(t *testing.T) {
		// kill_char returns early if not a character
		char_reclaim(locID)
	})
}

func TestSetLord(t *testing.T) {
	setupEconomyTest()

	playerID := 100
	newPlayerID := 101
	charID := 1001

	alloc_box(playerID, T_player, 0)
	alloc_box(newPlayerID, T_player, 0)
	alloc_box(charID, T_char, 0)

	p_char(charID).unit_lord = playerID

	t.Run("changes lord and preserves prev_lord", func(t *testing.T) {
		set_lord(charID, newPlayerID, LOY_oath, 100)

		c := rp_char(charID)
		if c.unit_lord != newPlayerID {
			t.Errorf("unit_lord = %d, want %d", c.unit_lord, newPlayerID)
		}
		if c.prev_lord != playerID {
			t.Errorf("prev_lord = %d, want %d", c.prev_lord, playerID)
		}
		if c.loy_kind != schar(LOY_oath) {
			t.Errorf("loy_kind = %d, want %d", c.loy_kind, LOY_oath)
		}
		if c.loy_rate != 100 {
			t.Errorf("loy_rate = %d, want 100", c.loy_rate)
		}
	})
}

func TestUnitDeserts(t *testing.T) {
	setupEconomyTest()

	playerID := 100
	charID := 1001
	locID := 5000

	alloc_box(playerID, T_player, 0)
	alloc_box(charID, T_char, 0)
	alloc_box(locID, T_loc, sub_plain)

	p_char(charID).unit_lord = playerID
	set_where(charID, locID)

	t.Run("sets lord to 0 when to_who is 0", func(t *testing.T) {
		unit_deserts(charID, 0, true, LOY_unsworn, 0)

		// unit_deserts sets lord to 0, doesn't delete the character
		c := rp_char(charID)
		if c.unit_lord != 0 {
			t.Errorf("unit_lord after desert to 0 = %d, want 0", c.unit_lord)
		}
	})

	t.Run("changes lord when to_who is valid", func(t *testing.T) {
		charID2 := 1002
		newPlayerID := 101

		alloc_box(charID2, T_char, 0)
		alloc_box(newPlayerID, T_player, 0)
		p_char(charID2).unit_lord = playerID

		unit_deserts(charID2, newPlayerID, true, LOY_unsworn, 0)

		c := rp_char(charID2)
		if c.unit_lord != newPlayerID {
			t.Errorf("unit_lord after desert = %d, want %d", c.unit_lord, newPlayerID)
		}
	})
}
