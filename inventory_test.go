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

func setupInventoryTest() {
	if teg.prng == nil {
		teg.prng = prng.New(rand.NewPCG(12345, 67890))
	}
	if teg.globals.inventories == nil {
		teg.globals.inventories = make(map[int][]item_ent)
	}
	for id := range teg.globals.inventories {
		delete(teg.globals.inventories, id)
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
	alloc_box(item_iron, T_item, 0)
}

func TestHasItem(t *testing.T) {
	setupInventoryTest()

	charID := 1000
	alloc_box(charID, T_char, 0)

	itemID := item_gold

	if got := has_item(charID, itemID); got != 0 {
		t.Errorf("has_item empty = %d, want 0", got)
	}

	teg.globals.inventories[charID] = []item_ent{{item: itemID, qty: 100}}

	if got := has_item(charID, itemID); got != 100 {
		t.Errorf("has_item with gold = %d, want 100", got)
	}

	if got := has_item(charID, item_iron); got != 0 {
		t.Errorf("has_item no iron = %d, want 0", got)
	}
}

func TestSubItem(t *testing.T) {
	setupInventoryTest()

	charID := 1000
	alloc_box(charID, T_char, 0)
	teg.globals.inventories[charID] = []item_ent{{item: item_gold, qty: 100}}

	if !sub_item(charID, item_gold, 30) {
		t.Error("sub_item should succeed")
	}
	if got := has_item(charID, item_gold); got != 70 {
		t.Errorf("after sub_item = %d, want 70", got)
	}

	if sub_item(charID, item_gold, 100) {
		t.Error("sub_item should fail (insufficient)")
	}
	if got := has_item(charID, item_gold); got != 70 {
		t.Errorf("after failed sub_item = %d, want 70 (unchanged)", got)
	}
}

func TestAddItem(t *testing.T) {
	setupInventoryTest()

	charID := 1000
	alloc_box(charID, T_char, 0)

	add_item(charID, item_gold, 50)
	if got := has_item(charID, item_gold); got != 50 {
		t.Errorf("add_item first = %d, want 50", got)
	}

	add_item(charID, item_gold, 25)
	if got := has_item(charID, item_gold); got != 75 {
		t.Errorf("add_item second = %d, want 75", got)
	}

	add_item(charID, item_iron, 10)
	if got := has_item(charID, item_iron); got != 10 {
		t.Errorf("add_item different item = %d, want 10", got)
	}
}

func TestGenItem(t *testing.T) {
	setupInventoryTest()

	charID := 1000
	alloc_box(charID, T_char, 0)

	gen_item(charID, item_gold, 100)
	if got := has_item(charID, item_gold); got != 100 {
		t.Errorf("gen_item = %d, want 100", got)
	}
}

func TestConsumeItem(t *testing.T) {
	setupInventoryTest()

	charID := 1000
	alloc_box(charID, T_char, 0)
	teg.globals.inventories[charID] = []item_ent{{item: item_gold, qty: 100}}

	if !consume_item(charID, item_gold, 50) {
		t.Error("consume_item should succeed")
	}
	if got := has_item(charID, item_gold); got != 50 {
		t.Errorf("after consume_item = %d, want 50", got)
	}
}

func TestMoveItem(t *testing.T) {
	setupInventoryTest()

	char1 := 1000
	char2 := 1001
	alloc_box(char1, T_char, 0)
	alloc_box(char2, T_char, 0)
	teg.globals.inventories[char1] = []item_ent{{item: item_gold, qty: 100}}

	if !move_item(char1, char2, item_gold, 40) {
		t.Error("move_item should succeed")
	}
	if got := has_item(char1, item_gold); got != 60 {
		t.Errorf("char1 after move = %d, want 60", got)
	}
	if got := has_item(char2, item_gold); got != 40 {
		t.Errorf("char2 after move = %d, want 40", got)
	}

	if move_item(char1, char2, item_gold, 100) {
		t.Error("move_item should fail (insufficient)")
	}
}

func TestCanPayAndCharge(t *testing.T) {
	setupInventoryTest()

	charID := 1000
	alloc_box(charID, T_char, 0)
	teg.globals.inventories[charID] = []item_ent{{item: item_gold, qty: 100}}

	if !can_pay(charID, 50) {
		t.Error("can_pay(50) should be true")
	}
	if can_pay(charID, 150) {
		t.Error("can_pay(150) should be false")
	}

	if !charge(charID, 30) {
		t.Error("charge(30) should succeed")
	}
	if got := has_item(charID, item_gold); got != 70 {
		t.Errorf("after charge = %d, want 70", got)
	}
}

func TestStackHasItem(t *testing.T) {
	setupInventoryTest()

	leader := 1000
	follower := 1001
	playerID := 100
	locID := 2000

	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(locID, T_loc, sub_plain)
	alloc_box(leader, T_char, 0)
	alloc_box(follower, T_char, 0)

	p_char(leader).unit_lord = playerID
	p_char(follower).unit_lord = playerID

	set_where(leader, locID)
	set_where(follower, leader)

	teg.globals.inventories[leader] = []item_ent{{item: item_gold, qty: 100}}
	teg.globals.inventories[follower] = []item_ent{{item: item_gold, qty: 50}}

	if got := stack_has_item(leader, item_gold); got != 150 {
		t.Errorf("stack_has_item = %d, want 150", got)
	}

	if got := stack_has_item(follower, item_gold); got != 150 {
		t.Errorf("stack_has_item from follower = %d, want 150", got)
	}
}

func TestStackSubItem(t *testing.T) {
	setupInventoryTest()

	leader := 1000
	follower := 1001
	playerID := 100
	locID := 2000

	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(locID, T_loc, sub_plain)
	alloc_box(leader, T_char, 0)
	alloc_box(follower, T_char, 0)

	p_char(leader).unit_lord = playerID
	p_char(follower).unit_lord = playerID

	set_where(leader, locID)
	set_where(follower, leader)

	teg.globals.inventories[leader] = []item_ent{{item: item_gold, qty: 100}}
	teg.globals.inventories[follower] = []item_ent{{item: item_gold, qty: 50}}

	if !stack_sub_item(leader, item_gold, 120) {
		t.Error("stack_sub_item should succeed with borrowing")
	}

	leaderGold := has_item(leader, item_gold)
	followerGold := has_item(follower, item_gold)
	if leaderGold+followerGold != 30 {
		t.Errorf("total after stack_sub_item = %d, want 30", leaderGold+followerGold)
	}
}

func TestAutocharge(t *testing.T) {
	setupInventoryTest()

	leader := 1000
	follower := 1001
	playerID := 100
	locID := 2000

	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(locID, T_loc, sub_plain)
	alloc_box(leader, T_char, 0)
	alloc_box(follower, T_char, 0)

	p_char(leader).unit_lord = playerID
	p_char(follower).unit_lord = playerID

	set_where(leader, locID)
	set_where(follower, leader)

	teg.globals.inventories[leader] = []item_ent{{item: item_gold, qty: 50}}
	teg.globals.inventories[follower] = []item_ent{{item: item_gold, qty: 50}}

	if !autocharge(leader, 75) {
		t.Error("autocharge should succeed with borrowing")
	}

	total := has_item(leader, item_gold) + has_item(follower, item_gold)
	if total != 25 {
		t.Errorf("total after autocharge = %d, want 25", total)
	}
}

func TestHasUseKey(t *testing.T) {
	setupInventoryTest()

	charID := 1000
	artifactID := 5000

	alloc_box(charID, T_char, 0)
	alloc_box(artifactID, T_item, sub_artifact)

	p_item_magic(artifactID).use_key = use_heal_potion

	teg.globals.inventories[charID] = []item_ent{{item: artifactID, qty: 1}}

	if got := has_use_key(charID, use_heal_potion); got != artifactID {
		t.Errorf("has_use_key = %d, want %d", got, artifactID)
	}

	if got := has_use_key(charID, use_death_potion); got != 0 {
		t.Errorf("has_use_key wrong key = %d, want 0", got)
	}
}

func TestAddNpDeductNp(t *testing.T) {
	setupInventoryTest()

	playerID := 100
	alloc_box(playerID, T_player, sub_pl_regular)
	p_player(playerID).noble_points = 10

	add_np(playerID, 5)
	if got := p_player(playerID).noble_points; got != 15 {
		t.Errorf("after add_np = %d, want 15", got)
	}
	if got := p_player(playerID).np_gained; got != 5 {
		t.Errorf("np_gained = %d, want 5", got)
	}

	if !deduct_np(playerID, 3) {
		t.Error("deduct_np should succeed")
	}
	if got := p_player(playerID).noble_points; got != 12 {
		t.Errorf("after deduct_np = %d, want 12", got)
	}
	if got := p_player(playerID).np_spent; got != 3 {
		t.Errorf("np_spent = %d, want 3", got)
	}

	if deduct_np(playerID, 100) {
		t.Error("deduct_np should fail (insufficient)")
	}
}

func TestCreateUniqueItem(t *testing.T) {
	setupInventoryTest()

	charID := 1000
	alloc_box(charID, T_char, 0)

	itemID := create_unique_item(charID, sub_artifact)
	if itemID < 0 {
		t.Fatal("create_unique_item failed")
	}

	if kind(itemID) != T_item {
		t.Errorf("created item kind = %d, want T_item", kind(itemID))
	}

	if got := item_unique(itemID); got != charID {
		t.Errorf("item_unique = %d, want %d", got, charID)
	}

	if got := has_item(charID, itemID); got != 1 {
		t.Errorf("has_item after create = %d, want 1", got)
	}
}

func TestHackUniqueItem(t *testing.T) {
	setupInventoryTest()

	charID := 1000
	itemID := 5000

	alloc_box(charID, T_char, 0)
	alloc_box(itemID, T_item, sub_artifact)

	hack_unique_item(itemID, charID)

	if got := item_unique(itemID); got != charID {
		t.Errorf("item_unique after hack = %d, want %d", got, charID)
	}

	if got := has_item(charID, itemID); got != 1 {
		t.Errorf("has_item after hack = %d, want 1", got)
	}
}

func TestMinMax(t *testing.T) {
	if got := min(5, 3); got != 3 {
		t.Errorf("min(5,3) = %d, want 3", got)
	}
	if got := min(3, 5); got != 3 {
		t.Errorf("min(3,5) = %d, want 3", got)
	}
	if got := max(5, 3); got != 5 {
		t.Errorf("max(5,3) = %d, want 5", got)
	}
	if got := max(3, 5); got != 5 {
		t.Errorf("max(3,5) = %d, want 5", got)
	}
}

func TestAddS(t *testing.T) {
	if got := add_s(1); got != "" {
		t.Errorf("add_s(1) = %q, want \"\"", got)
	}
	if got := add_s(0); got != "s" {
		t.Errorf("add_s(0) = %q, want \"s\"", got)
	}
	if got := add_s(2); got != "s" {
		t.Errorf("add_s(2) = %q, want \"s\"", got)
	}
}

func TestDeductAura(t *testing.T) {
	setupInventoryTest()

	charID := 1000
	alloc_box(charID, T_char, 0)
	p_magic(charID).cur_aura = 10

	if !deduct_aura(charID, 3) {
		t.Error("deduct_aura(3) should succeed")
	}
	if got := char_cur_aura(charID); got != 7 {
		t.Errorf("cur_aura after deduct = %d, want 7", got)
	}

	if !deduct_aura(charID, 7) {
		t.Error("deduct_aura(7) should succeed")
	}
	if got := char_cur_aura(charID); got != 0 {
		t.Errorf("cur_aura after second deduct = %d, want 0", got)
	}

	if deduct_aura(charID, 1) {
		t.Error("deduct_aura(1) should fail (insufficient)")
	}
}

func TestDeductAuraNilMagic(t *testing.T) {
	setupInventoryTest()

	charID := 1001
	alloc_box(charID, T_char, 0)

	if deduct_aura(charID, 1) {
		t.Error("deduct_aura should fail for nil magic")
	}
}

func TestChargeAura(t *testing.T) {
	setupInventoryTest()

	charID := 1000
	alloc_box(charID, T_char, 0)
	p_magic(charID).cur_aura = 5

	if !charge_aura(charID, 3) {
		t.Error("charge_aura(3) should succeed")
	}
	if got := char_cur_aura(charID); got != 2 {
		t.Errorf("cur_aura after charge = %d, want 2", got)
	}

	if charge_aura(charID, 10) {
		t.Error("charge_aura(10) should fail (insufficient)")
	}
	if got := char_cur_aura(charID); got != 2 {
		t.Errorf("cur_aura should be unchanged = %d, want 2", got)
	}
}

func TestCheckAura(t *testing.T) {
	setupInventoryTest()

	charID := 1000
	alloc_box(charID, T_char, 0)
	p_magic(charID).cur_aura = 5

	if !check_aura(charID, 5) {
		t.Error("check_aura(5) should succeed")
	}
	if got := char_cur_aura(charID); got != 5 {
		t.Errorf("cur_aura should be unchanged after check = %d, want 5", got)
	}

	if !check_aura(charID, 3) {
		t.Error("check_aura(3) should succeed")
	}

	if check_aura(charID, 10) {
		t.Error("check_aura(10) should fail (insufficient)")
	}
	if got := char_cur_aura(charID); got != 5 {
		t.Errorf("cur_aura should be unchanged = %d, want 5", got)
	}
}

// p_item_magic is defined in accessor.go
