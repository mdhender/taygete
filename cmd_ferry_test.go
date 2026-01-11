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

// cmd_ferry_test.go - Unit tests for ferry commands
// Sprint 26.8: Ships & Ferries

package taygete

import (
	"math/rand/v2"
	"testing"

	"github.com/mdhender/prng"
)

func setupFerryTest() {
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
	alloc_box(item_peasant, T_item, 0)
	p_item(item_peasant).weight = 100
}

func TestVFee(t *testing.T) {
	setupFerryTest()

	playerID := 100
	captainID := 2001
	shipID := 3001
	provinceID := 1001

	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(captainID, T_char, 0)
	alloc_box(shipID, T_ship, 0)
	alloc_box(provinceID, T_loc, sub_plain)

	set_where(shipID, provinceID)
	set_where(captainID, shipID)
	p_loc_info(shipID).here_list = []int{captainID}
	p_loc_info(provinceID).here_list = []int{shipID}

	p_char(captainID).unit_lord = playerID

	c := &command{who: captainID, a: 25}
	result := v_fee(c)

	if result != TRUE {
		t.Errorf("v_fee returned %d, want TRUE", result)
	}

	if got := board_fee(captainID); got != 25 {
		t.Errorf("board_fee(captain) = %d, want 25", got)
	}

	c.a = 0
	result = v_fee(c)
	if result != TRUE {
		t.Errorf("v_fee(0) returned %d, want TRUE", result)
	}

	if got := board_fee(captainID); got != 0 {
		t.Errorf("board_fee(captain) = %d, want 0", got)
	}
}

func TestVBoard_NotAShip(t *testing.T) {
	setupFerryTest()

	charID := 2001
	provinceID := 1001

	alloc_box(charID, T_char, 0)
	alloc_box(provinceID, T_loc, sub_plain)

	set_where(charID, provinceID)
	p_loc_info(provinceID).here_list = []int{charID}

	c := &command{who: charID, a: 9999}
	result := v_board(c)

	if result != FALSE {
		t.Errorf("v_board(non-ship) returned %d, want FALSE", result)
	}
}

func TestVBoard_NoFeeSet(t *testing.T) {
	setupFerryTest()

	playerID := 100
	charID := 2001
	captainID := 2002
	shipID := 3001
	provinceID := 1001

	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(charID, T_char, 0)
	alloc_box(captainID, T_char, 0)
	alloc_box(shipID, T_ship, 0)
	alloc_box(provinceID, T_loc, sub_plain)

	set_where(charID, provinceID)
	set_where(shipID, provinceID)
	set_where(captainID, shipID)

	p_char(charID).unit_lord = playerID
	p_char(captainID).unit_lord = playerID

	p_loc_info(provinceID).here_list = []int{charID, shipID}
	p_loc_info(shipID).here_list = []int{captainID}

	c := &command{who: charID, a: shipID}
	result := v_board(c)

	if result != FALSE {
		t.Errorf("v_board(no-fee) returned %d, want FALSE", result)
	}
}

func TestVBoard_ShipOverloaded(t *testing.T) {
	setupFerryTest()

	playerID := 100
	player2ID := 101
	charID := 2001
	captainID := 2002
	shipID := 3001
	provinceID := 1001

	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(player2ID, T_player, sub_pl_regular)
	alloc_box(charID, T_char, 0)
	alloc_box(captainID, T_char, 0)
	alloc_box(shipID, T_ship, 0)
	alloc_box(provinceID, T_loc, sub_plain)

	set_where(charID, provinceID)
	set_where(shipID, provinceID)
	set_where(captainID, shipID)

	p_char(charID).unit_lord = playerID
	p_char(captainID).unit_lord = player2ID

	p_loc_info(provinceID).here_list = []int{charID, shipID}
	p_loc_info(shipID).here_list = []int{captainID}

	p_magic(captainID).fee = 10
	p_subloc(shipID).capacity = 10

	gen_item(charID, item_gold, 1000)

	c := &command{who: charID, a: shipID}
	result := v_board(c)

	if result != FALSE {
		t.Errorf("v_board(overloaded) returned %d, want FALSE", result)
	}
}

func TestVBoard_RefuseFee(t *testing.T) {
	setupFerryTest()

	playerID := 100
	player2ID := 101
	charID := 2001
	captainID := 2002
	shipID := 3001
	provinceID := 1001

	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(player2ID, T_player, sub_pl_regular)
	alloc_box(charID, T_char, 0)
	alloc_box(captainID, T_char, 0)
	alloc_box(shipID, T_ship, 0)
	alloc_box(provinceID, T_loc, sub_plain)

	set_where(charID, provinceID)
	set_where(shipID, provinceID)
	set_where(captainID, shipID)

	p_char(charID).unit_lord = playerID
	p_char(captainID).unit_lord = player2ID

	p_loc_info(provinceID).here_list = []int{charID, shipID}
	p_loc_info(shipID).here_list = []int{captainID}

	p_magic(captainID).fee = 100
	p_subloc(shipID).capacity = 100000

	gen_item(charID, item_gold, 1000)

	c := &command{who: charID, a: shipID, b: 5}
	result := v_board(c)

	if result != FALSE {
		t.Errorf("v_board(refuse-fee) returned %d, want FALSE", result)
	}
}

func TestVBoard_CantAfford(t *testing.T) {
	setupFerryTest()

	playerID := 100
	player2ID := 101
	charID := 2001
	captainID := 2002
	shipID := 3001
	provinceID := 1001

	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(player2ID, T_player, sub_pl_regular)
	alloc_box(charID, T_char, 0)
	alloc_box(captainID, T_char, 0)
	alloc_box(shipID, T_ship, 0)
	alloc_box(provinceID, T_loc, sub_plain)

	set_where(charID, provinceID)
	set_where(shipID, provinceID)
	set_where(captainID, shipID)

	p_char(charID).unit_lord = playerID
	p_char(captainID).unit_lord = player2ID

	p_loc_info(provinceID).here_list = []int{charID, shipID}
	p_loc_info(shipID).here_list = []int{captainID}

	p_magic(captainID).fee = 100
	p_subloc(shipID).capacity = 100000

	c := &command{who: charID, a: shipID}
	result := v_board(c)

	if result != FALSE {
		t.Errorf("v_board(cant-afford) returned %d, want FALSE", result)
	}
}

func TestVUnload_NotCaptain(t *testing.T) {
	setupFerryTest()

	charID := 2001
	provinceID := 1001

	alloc_box(charID, T_char, 0)
	alloc_box(provinceID, T_loc, sub_plain)

	set_where(charID, provinceID)
	p_loc_info(provinceID).here_list = []int{charID}

	c := &command{who: charID}
	result := v_unload(c)

	if result != FALSE {
		t.Errorf("v_unload(not-captain) returned %d, want FALSE", result)
	}
}

func TestVUnload_AtSea(t *testing.T) {
	setupFerryTest()

	playerID := 100
	captainID := 2001
	shipID := 3001
	oceanID := 1001

	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(captainID, T_char, 0)
	alloc_box(shipID, T_ship, 0)
	alloc_box(oceanID, T_loc, sub_ocean)

	set_where(shipID, oceanID)
	set_where(captainID, shipID)

	p_char(captainID).unit_lord = playerID

	p_loc_info(oceanID).here_list = []int{shipID}
	p_loc_info(shipID).here_list = []int{captainID}

	c := &command{who: captainID}
	result := v_unload(c)

	if result != FALSE {
		t.Errorf("v_unload(at-sea) returned %d, want FALSE", result)
	}
}

func TestVUnload_NoPassengers(t *testing.T) {
	setupFerryTest()

	playerID := 100
	captainID := 2001
	shipID := 3001
	provinceID := 1001

	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(captainID, T_char, 0)
	alloc_box(shipID, T_ship, sub_galley)
	alloc_box(provinceID, T_loc, sub_plain)

	set_where(shipID, provinceID)
	set_where(captainID, shipID)

	p_char(captainID).unit_lord = playerID

	p_loc_info(provinceID).here_list = []int{shipID}
	p_loc_info(shipID).here_list = []int{captainID}

	c := &command{who: captainID}
	result := v_unload(c)

	if result != TRUE {
		t.Errorf("v_unload(no-passengers) returned %d, want TRUE", result)
	}
}

func TestVFerry_NotCaptain(t *testing.T) {
	setupFerryTest()

	charID := 2001
	provinceID := 1001

	alloc_box(charID, T_char, 0)
	alloc_box(provinceID, T_loc, sub_plain)

	set_where(charID, provinceID)
	p_loc_info(provinceID).here_list = []int{charID}

	c := &command{who: charID}
	result := v_ferry(c)

	if result != FALSE {
		t.Errorf("v_ferry(not-captain) returned %d, want FALSE", result)
	}
}

func TestVFerry_Success(t *testing.T) {
	setupFerryTest()

	playerID := 100
	captainID := 2001
	shipID := 3001
	provinceID := 1001

	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(captainID, T_char, 0)
	alloc_box(shipID, T_ship, sub_galley)
	alloc_box(provinceID, T_loc, sub_plain)

	set_where(shipID, provinceID)
	set_where(captainID, shipID)

	p_char(captainID).unit_lord = playerID

	p_loc_info(provinceID).here_list = []int{shipID}
	p_loc_info(shipID).here_list = []int{captainID}

	c := &command{who: captainID}
	result := v_ferry(c)

	if result != TRUE {
		t.Errorf("v_ferry returned %d, want TRUE", result)
	}

	if p_magic(shipID).ferry_flag != TRUE {
		t.Errorf("ferry_flag = %d, want TRUE", p_magic(shipID).ferry_flag)
	}
}

func TestLookupDir(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"north", DIR_N},
		{"NORTH", DIR_N},
		{"n", DIR_N},
		{"N", DIR_N},
		{"east", DIR_E},
		{"e", DIR_E},
		{"south", DIR_S},
		{"s", DIR_S},
		{"west", DIR_W},
		{"w", DIR_W},
		{"up", DIR_UP},
		{"u", DIR_UP},
		{"down", DIR_DOWN},
		{"d", DIR_DOWN},
		{"in", DIR_IN},
		{"i", DIR_IN},
		{"out", DIR_OUT},
		{"o", DIR_OUT},
		{"", -1},
		{"invalid", -1},
	}

	for _, tt := range tests {
		got := lookup_dir(tt.input)
		if got != tt.expected {
			t.Errorf("lookup_dir(%q) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}

func TestDirName(t *testing.T) {
	tests := []struct {
		dir      int
		expected string
	}{
		{DIR_N, "north"},
		{DIR_E, "east"},
		{DIR_S, "south"},
		{DIR_W, "west"},
		{DIR_UP, "up"},
		{DIR_DOWN, "down"},
		{DIR_IN, "in"},
		{DIR_OUT, "out"},
		{-1, "unknown"},
		{100, "unknown"},
	}

	for _, tt := range tests {
		got := dir_name(tt.dir)
		if got != tt.expected {
			t.Errorf("dir_name(%d) = %q, want %q", tt.dir, got, tt.expected)
		}
	}
}

func TestBoardMessage_Hidden(t *testing.T) {
	setupFerryTest()

	charID := 2001
	shipID := 3001
	provinceID := 1001

	alloc_box(charID, T_char, 0)
	alloc_box(shipID, T_ship, 0)
	alloc_box(provinceID, T_loc, sub_plain)

	set_where(shipID, provinceID)
	p_magic(charID).hide_self = 1

	board_message(charID, shipID)
}

func TestUnboardMessage_Hidden(t *testing.T) {
	setupFerryTest()

	charID := 2001
	shipID := 3001
	provinceID := 1001

	alloc_box(charID, T_char, 0)
	alloc_box(shipID, T_ship, 0)
	alloc_box(provinceID, T_loc, sub_plain)

	set_where(shipID, provinceID)
	p_magic(charID).hide_self = 1

	unboard_message(charID, shipID)
}
