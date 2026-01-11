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

// cmd_form_test.go - Unit tests for Sprint 26.6 noble formation commands

package taygete

import (
	"math/rand/v2"
	"testing"

	"github.com/mdhender/prng"
)

func setupFormTest() {
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
	if teg.globals.playerUnits == nil {
		teg.globals.playerUnits = make(map[int][]int)
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
	for id := range teg.globals.playerUnits {
		delete(teg.globals.playerUnits, id)
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

	teg.globals.sysclock = olytime{turn: 1, day: 1}
}

func TestNobleCost(t *testing.T) {
	if cost := noble_cost(100); cost != 1 {
		t.Errorf("noble_cost(100) = %d, want 1", cost)
	}
}

func TestNextNpTurn(t *testing.T) {
	setupFormTest()

	// Formula: ct = (7 - (turn + 1) % 8), ft = 2, n = (ft + ct) % 8
	// Turn 1: ct = 7 - 2 = 5, n = (2 + 5) % 8 = 7
	// Turn 2: ct = 7 - 3 = 4, n = (2 + 4) % 8 = 6
	// Turn 6: ct = 7 - 7 = 0, n = (2 + 0) % 8 = 2
	// Turn 7: ct = 7 - 0 = 7, n = (2 + 7) % 8 = 1
	// Turn 8: ct = 7 - 1 = 6, n = (2 + 6) % 8 = 0  <- NP granted this turn
	// Turn 9: ct = 7 - 2 = 5, n = (2 + 5) % 8 = 7
	// Turn 16: ct = 7 - 1 = 6, n = (2 + 6) % 8 = 0 <- NP granted
	tests := []struct {
		turn int
		want int
	}{
		{1, 7},
		{2, 6},
		{6, 2},
		{7, 1},
		{8, 0},
		{9, 7},
		{16, 0},
	}

	for _, tc := range tests {
		teg.globals.sysclock.turn = short(tc.turn)
		got := next_np_turn(100)
		if got != tc.want {
			t.Errorf("next_np_turn at turn %d = %d, want %d", tc.turn, got, tc.want)
		}
	}
}

func TestPlayerUnformedHelpers(t *testing.T) {
	setupFormTest()

	playerID := 100
	alloc_box(playerID, T_player, sub_pl_regular)

	if got := getPlayerUnformed(playerID); got != nil {
		t.Errorf("getPlayerUnformed should return nil initially, got %v", got)
	}

	addPlayerUnformed(playerID, 5001)
	addPlayerUnformed(playerID, 5002)
	addPlayerUnformed(playerID, 5003)

	unformed := getPlayerUnformed(playerID)
	if len(unformed) != 3 {
		t.Fatalf("expected 3 unformed nobles, got %d", len(unformed))
	}
	if unformed[0] != 5001 || unformed[1] != 5002 || unformed[2] != 5003 {
		t.Errorf("unexpected unformed list: %v", unformed)
	}

	removePlayerUnformed(playerID, 5002)
	unformed = getPlayerUnformed(playerID)
	if len(unformed) != 2 {
		t.Fatalf("expected 2 unformed nobles after removal, got %d", len(unformed))
	}
	if unformed[0] != 5001 || unformed[1] != 5003 {
		t.Errorf("unexpected unformed list after removal: %v", unformed)
	}
}

func TestContainsInt(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	if !containsInt(slice, 3) {
		t.Error("containsInt should find 3")
	}
	if containsInt(slice, 10) {
		t.Error("containsInt should not find 10")
	}
	if containsInt(nil, 1) {
		t.Error("containsInt on nil should return false")
	}
}

func TestJoinStrings(t *testing.T) {
	if got := joinStrings([]string{"a", "b", "c"}, " "); got != "a b c" {
		t.Errorf("joinStrings = %q, want %q", got, "a b c")
	}
	if got := joinStrings([]string{}, " "); got != "" {
		t.Errorf("joinStrings empty = %q, want %q", got, "")
	}
	if got := joinStrings([]string{"only"}, ","); got != "only" {
		t.Errorf("joinStrings single = %q, want %q", got, "only")
	}
}

func TestVFormRequiresCity(t *testing.T) {
	setupFormTest()

	charID := 1001
	playerID := 100
	province := 10001

	alloc_box(charID, T_char, 0)
	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(province, T_loc, sub_plain)

	p_char(charID).unit_lord = playerID
	p_player(playerID).noble_points = 5

	set_where(charID, province)

	c := &command{who: charID}
	result := v_form(c)
	if result != FALSE {
		t.Error("v_form should fail when not in a city")
	}
}

func TestVFormInCity(t *testing.T) {
	setupFormTest()

	charID := 1001
	playerID := 100
	cityID := 59001

	alloc_box(charID, T_char, 0)
	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(cityID, T_loc, sub_city)

	p_char(charID).unit_lord = playerID
	p_player(playerID).noble_points = 5

	set_where(charID, cityID)

	c := &command{who: charID}
	result := v_form(c)
	if result != TRUE {
		t.Error("v_form should succeed when in a city with NP")
	}
}

func TestVFormInsufficientNP(t *testing.T) {
	setupFormTest()

	charID := 1001
	playerID := 100
	cityID := 59001

	alloc_box(charID, T_char, 0)
	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(cityID, T_loc, sub_city)

	p_char(charID).unit_lord = playerID
	p_player(playerID).noble_points = 0

	set_where(charID, cityID)

	c := &command{who: charID}
	result := v_form(c)
	if result != FALSE {
		t.Error("v_form should fail with 0 NP")
	}
}

func TestDFormCreatesNoble(t *testing.T) {
	setupFormTest()

	charID := 1001
	playerID := 100
	cityID := 59001
	unformedID := 5001

	alloc_box(charID, T_char, 0)
	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(cityID, T_loc, sub_city)
	alloc_box(unformedID, T_unform, 0)

	p_char(charID).unit_lord = playerID
	p_player(playerID).noble_points = 5

	set_where(charID, cityID)
	addPlayerUnformed(playerID, unformedID)

	c := &command{who: charID}
	result := d_form(c)
	if result != TRUE {
		t.Error("d_form should succeed")
	}

	if kind(unformedID) != T_char {
		t.Errorf("unformed noble should now be T_char, got %d", kind(unformedID))
	}

	newChar := rp_char(unformedID)
	if newChar == nil {
		t.Fatal("new noble should have character data")
	}
	if newChar.health != 100 {
		t.Errorf("new noble health = %d, want 100", newChar.health)
	}
	if newChar.attack != 80 {
		t.Errorf("new noble attack = %d, want 80", newChar.attack)
	}
	if newChar.defense != 80 {
		t.Errorf("new noble defense = %d, want 80", newChar.defense)
	}
	if newChar.break_point != 50 {
		t.Errorf("new noble break_point = %d, want 50", newChar.break_point)
	}

	if player_np(playerID) != 4 {
		t.Errorf("NP after formation = %d, want 4", player_np(playerID))
	}

	unformed := getPlayerUnformed(playerID)
	if len(unformed) != 0 {
		t.Errorf("unformed list should be empty after formation, got %v", unformed)
	}
}

func TestDFormWithSpecificID(t *testing.T) {
	setupFormTest()

	charID := 1001
	playerID := 100
	cityID := 59001
	unformedID1 := 5001
	unformedID2 := 5002
	unformedID3 := 5003

	alloc_box(charID, T_char, 0)
	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(cityID, T_loc, sub_city)
	alloc_box(unformedID1, T_unform, 0)
	alloc_box(unformedID2, T_unform, 0)
	alloc_box(unformedID3, T_unform, 0)

	p_char(charID).unit_lord = playerID
	p_player(playerID).noble_points = 5

	set_where(charID, cityID)
	addPlayerUnformed(playerID, unformedID1)
	addPlayerUnformed(playerID, unformedID2)
	addPlayerUnformed(playerID, unformedID3)

	c := &command{who: charID, a: unformedID2}
	result := d_form(c)
	if result != TRUE {
		t.Error("d_form with specific ID should succeed")
	}

	if kind(unformedID2) != T_char {
		t.Errorf("specified unformed noble should now be T_char, got %d", kind(unformedID2))
	}

	if kind(unformedID1) != T_unform {
		t.Error("other unformed nobles should remain T_unform")
	}
	if kind(unformedID3) != T_unform {
		t.Error("other unformed nobles should remain T_unform")
	}

	unformed := getPlayerUnformed(playerID)
	if len(unformed) != 2 {
		t.Errorf("expected 2 remaining unformed, got %d", len(unformed))
	}
}

func TestDFormNoUnformed(t *testing.T) {
	setupFormTest()

	charID := 1001
	playerID := 100
	cityID := 59001

	alloc_box(charID, T_char, 0)
	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(cityID, T_loc, sub_city)

	p_char(charID).unit_lord = playerID
	p_player(playerID).noble_points = 5

	set_where(charID, cityID)

	c := &command{who: charID}
	result := d_form(c)
	if result != FALSE {
		t.Error("d_form should fail with no unformed nobles")
	}
}

func TestDFormInvalidUnformedID(t *testing.T) {
	setupFormTest()

	charID := 1001
	playerID := 100
	cityID := 59001
	unformedID := 5001

	alloc_box(charID, T_char, 0)
	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(cityID, T_loc, sub_city)
	alloc_box(unformedID, T_unform, 0)

	p_char(charID).unit_lord = playerID
	p_player(playerID).noble_points = 5

	set_where(charID, cityID)
	addPlayerUnformed(playerID, unformedID)

	c := &command{who: charID, a: 9999}
	result := d_form(c)
	if result != TRUE {
		t.Error("d_form with invalid ID should fall back to first unformed")
	}

	if kind(unformedID) != T_char {
		t.Errorf("fallback to first unformed should work, got kind %d", kind(unformedID))
	}
}

func TestFormNewNobleStats(t *testing.T) {
	setupFormTest()

	charID := 1001
	playerID := 100
	cityID := 59001
	newID := 5001

	alloc_box(charID, T_char, 0)
	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(cityID, T_loc, sub_city)
	alloc_box(newID, T_unform, 0)

	p_char(charID).unit_lord = playerID
	p_char(charID).behind = 1

	set_where(charID, cityID)

	form_new_noble(charID, "Sir Test", newID)

	if kind(newID) != T_char {
		t.Errorf("new noble kind = %d, want T_char", kind(newID))
	}

	n := rp_char(newID)
	if n.behind != 1 {
		t.Errorf("new noble behind = %d, want 1 (inherited from forming char)", n.behind)
	}
	if n.fresh_hire != TRUE {
		t.Error("new noble should be fresh_hire")
	}
	if n.health != 100 {
		t.Errorf("new noble health = %d, want 100", n.health)
	}
	if n.attack != 80 {
		t.Errorf("new noble attack = %d, want 80", n.attack)
	}
	if n.defense != 80 {
		t.Errorf("new noble defense = %d, want 80", n.defense)
	}
	if n.break_point != 50 {
		t.Errorf("new noble break_point = %d, want 50", n.break_point)
	}
	if n.unit_lord != playerID {
		t.Errorf("new noble unit_lord = %d, want %d", n.unit_lord, playerID)
	}
	if n.loy_kind != LOY_contract {
		t.Errorf("new noble loy_kind = %d, want LOY_contract", n.loy_kind)
	}
	if n.loy_rate != 500 {
		t.Errorf("new noble loy_rate = %d, want 500", n.loy_rate)
	}
}

func TestPrintHiringStatus(t *testing.T) {
	setupFormTest()

	playerID := 100
	alloc_box(playerID, T_player, sub_pl_regular)

	print_hiring_status(playerID)
}

func TestPrintUnformed(t *testing.T) {
	setupFormTest()

	playerID := 100
	alloc_box(playerID, T_player, sub_pl_regular)

	addPlayerUnformed(playerID, 5001)
	addPlayerUnformed(playerID, 5002)

	print_unformed(playerID)
}
