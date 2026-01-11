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

// cmd_transfer_test.go - Unit tests for Sprint 26.5 inventory transfer commands

package taygete

import (
	"math/rand/v2"
	"testing"

	"github.com/mdhender/prng"
)

func setupTransferTest() {
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
	alloc_box(item_soldier, T_item, 0)
	alloc_box(item_sailor, T_item, 0)
	alloc_box(item_pirate, T_item, 0)

	p_item(item_peasant).is_man_item = 1
	p_item(item_soldier).is_man_item = 1
	p_item(item_sailor).is_man_item = 1
	p_item(item_pirate).is_man_item = 1
}

func TestVAccept(t *testing.T) {
	setupTransferTest()

	charID := 1001
	playerID := 2001

	alloc_box(charID, T_char, 0)
	alloc_box(playerID, T_player, sub_pl_regular)
	p_char(charID).unit_lord = playerID

	c := &command{
		who: charID,
		a:   0,         // from anyone
		b:   item_gold, // gold
		c:   100,       // up to 100
	}

	result := v_accept(c)
	if result != TRUE {
		t.Errorf("v_accept returned %d, want %d", result, TRUE)
	}

	p := rp_char(charID)
	if p == nil {
		t.Fatal("rp_char returned nil")
	}
	if len(p.accept) != 1 {
		t.Fatalf("expected 1 accept entry, got %d", len(p.accept))
	}
	if p.accept[0].from_who != 0 {
		t.Errorf("from_who = %d, want 0", p.accept[0].from_who)
	}
	if p.accept[0].item != item_gold {
		t.Errorf("item = %d, want %d", p.accept[0].item, item_gold)
	}
	if p.accept[0].qty != 100 {
		t.Errorf("qty = %d, want 100", p.accept[0].qty)
	}
}

func TestWillAcceptSup(t *testing.T) {
	setupTransferTest()

	charID := 1001
	playerID := 2001

	alloc_box(charID, T_char, 0)
	alloc_box(playerID, T_player, sub_pl_regular)
	p_char(charID).unit_lord = playerID

	p := p_char(charID)
	p.accept = append(p.accept, &accept_ent{
		item:     item_gold,
		from_who: 0,
		qty:      0,
	})

	if !will_accept_sup(charID, item_gold, 999, 50) {
		t.Error("expected will_accept_sup to return true for matching item")
	}

	if will_accept_sup(charID, item_peasant, 999, 50) {
		t.Error("expected will_accept_sup to return false for non-matching item")
	}
}

func TestWillAccept(t *testing.T) {
	setupTransferTest()

	char1 := 1001
	char2 := 1002
	player1 := 2001
	player2 := 2002

	alloc_box(char1, T_char, 0)
	alloc_box(char2, T_char, 0)
	alloc_box(player1, T_player, sub_pl_regular)
	alloc_box(player2, T_player, sub_pl_regular)

	p_char(char1).unit_lord = player1
	p_char(char2).unit_lord = player2

	t.Run("gold is always accepted", func(t *testing.T) {
		if !will_accept(char1, item_gold, char2, 50) {
			t.Error("gold should always be accepted")
		}
	})

	t.Run("same faction always accepted", func(t *testing.T) {
		p_char(char2).unit_lord = player1
		if !will_accept(char1, item_peasant, char2, 10) {
			t.Error("same faction should always be accepted")
		}
		p_char(char2).unit_lord = player2 // reset
	})

	t.Run("different faction without accept rule refused", func(t *testing.T) {
		if will_accept(char1, item_peasant, char2, 10) {
			t.Error("different faction without accept rule should be refused")
		}
	})

	t.Run("with accept rule", func(t *testing.T) {
		p := p_char(char1)
		p.accept = append(p.accept, &accept_ent{
			item:     item_peasant,
			from_who: char2,
			qty:      0,
		})
		if !will_accept(char1, item_peasant, char2, 10) {
			t.Error("should accept with matching accept rule")
		}
	})
}

func TestVGive(t *testing.T) {
	setupTransferTest()

	char1 := 1001
	char2 := 1002
	player1 := 2001
	locID := 3001

	alloc_box(char1, T_char, 0)
	alloc_box(char2, T_char, 0)
	alloc_box(player1, T_player, sub_pl_regular)
	alloc_box(locID, T_loc, sub_plain)

	p_char(char1).unit_lord = player1
	p_char(char2).unit_lord = player1
	p_loc_info(char1).where = locID
	p_loc_info(char2).where = locID
	p_loc_info(locID).here_list = []int{char1, char2}

	add_item(char1, item_gold, 100)

	c := &command{
		who: char1,
		a:   char2,
		b:   item_gold,
		c:   50,
		d:   0,
	}

	result := v_give(c)
	if result != TRUE {
		t.Errorf("v_give returned %d, want %d", result, TRUE)
	}

	if has_item(char1, item_gold) != 50 {
		t.Errorf("char1 should have 50 gold, has %d", has_item(char1, item_gold))
	}
	if has_item(char2, item_gold) != 50 {
		t.Errorf("char2 should have 50 gold, has %d", has_item(char2, item_gold))
	}
}

func TestVGiveDifferentFaction(t *testing.T) {
	setupTransferTest()

	char1 := 1001
	char2 := 1002
	player1 := 2001
	player2 := 2002
	locID := 3001

	alloc_box(char1, T_char, 0)
	alloc_box(char2, T_char, 0)
	alloc_box(player1, T_player, sub_pl_regular)
	alloc_box(player2, T_player, sub_pl_regular)
	alloc_box(locID, T_loc, sub_plain)

	p_char(char1).unit_lord = player1
	p_char(char2).unit_lord = player2
	p_loc_info(char1).where = locID
	p_loc_info(char2).where = locID
	p_loc_info(locID).here_list = []int{char1, char2}

	add_item(char1, item_gold, 100)

	c := &command{
		who: char1,
		a:   char2,
		b:   item_gold,
		c:   50,
		d:   0,
	}

	t.Run("gold to different faction succeeds", func(t *testing.T) {
		result := v_give(c)
		if result != TRUE {
			t.Errorf("v_give gold to different faction should succeed, got %d", result)
		}
	})

	t.Run("non-gold to different faction without accept fails", func(t *testing.T) {
		add_item(char1, item_peasant, 10)
		c.b = item_peasant
		c.c = 5

		result := v_give(c)
		if result != FALSE {
			t.Errorf("v_give non-gold to different faction without accept should fail, got %d", result)
		}
	})

	t.Run("non-gold with accept rule succeeds", func(t *testing.T) {
		p := p_char(char2)
		p.accept = append(p.accept, &accept_ent{
			item:     item_peasant,
			from_who: 0,
			qty:      0,
		})

		add_item(char1, item_peasant, 10) // ensure we have some
		c.b = item_peasant
		c.c = 5

		result := v_give(c)
		if result != TRUE {
			t.Errorf("v_give with accept rule should succeed, got %d", result)
		}
	})
}

func TestVPay(t *testing.T) {
	setupTransferTest()

	char1 := 1001
	char2 := 1002
	player1 := 2001
	locID := 3001

	alloc_box(char1, T_char, 0)
	alloc_box(char2, T_char, 0)
	alloc_box(player1, T_player, sub_pl_regular)
	alloc_box(locID, T_loc, sub_plain)

	p_char(char1).unit_lord = player1
	p_char(char2).unit_lord = player1
	p_loc_info(char1).where = locID
	p_loc_info(char2).where = locID
	p_loc_info(locID).here_list = []int{char1, char2}

	add_item(char1, item_gold, 100)

	c := &command{
		who: char1,
		a:   char2,
		b:   30,
		c:   0,
	}

	result := v_pay(c)
	if result != TRUE {
		t.Errorf("v_pay returned %d, want %d", result, TRUE)
	}

	if has_item(char1, item_gold) != 70 {
		t.Errorf("char1 should have 70 gold, has %d", has_item(char1, item_gold))
	}
	if has_item(char2, item_gold) != 30 {
		t.Errorf("char2 should have 30 gold, has %d", has_item(char2, item_gold))
	}
}

func TestMayTake(t *testing.T) {
	setupTransferTest()

	char1 := 1001
	char2 := 1002
	char3 := 1003
	player1 := 2001
	player2 := 2002
	locID := 3001

	alloc_box(char1, T_char, 0)
	alloc_box(char2, T_char, 0)
	alloc_box(char3, T_char, 0)
	alloc_box(player1, T_player, sub_pl_regular)
	alloc_box(player2, T_player, sub_pl_regular)
	alloc_box(locID, T_loc, sub_plain)

	p_char(char1).unit_lord = player1
	p_char(char2).unit_lord = player1
	p_char(char3).unit_lord = player2
	p_loc_info(char1).where = locID
	p_loc_info(char2).where = locID
	p_loc_info(char3).where = locID
	p_loc_info(locID).here_list = []int{char1, char2, char3}

	t.Run("same faction allowed", func(t *testing.T) {
		if !may_take(char1, char2) {
			t.Error("may_take should allow same faction")
		}
	})

	t.Run("different faction denied", func(t *testing.T) {
		if may_take(char1, char3) {
			t.Error("may_take should deny different faction")
		}
	})
}

func TestVGet(t *testing.T) {
	setupTransferTest()

	char1 := 1001
	char2 := 1002
	player1 := 2001
	locID := 3001

	alloc_box(char1, T_char, 0)
	alloc_box(char2, T_char, 0)
	alloc_box(player1, T_player, sub_pl_regular)
	alloc_box(locID, T_loc, sub_plain)

	p_char(char1).unit_lord = player1
	p_char(char2).unit_lord = player1
	p_loc_info(char1).where = locID
	p_loc_info(char2).where = locID
	p_loc_info(locID).here_list = []int{char1, char2}

	add_item(char2, item_gold, 100)

	c := &command{
		who: char1,
		a:   char2,
		b:   item_gold,
		c:   40,
		d:   0,
	}

	result := v_get(c)
	if result != TRUE {
		t.Errorf("v_get returned %d, want %d", result, TRUE)
	}

	if has_item(char1, item_gold) != 40 {
		t.Errorf("char1 should have 40 gold, has %d", has_item(char1, item_gold))
	}
	if has_item(char2, item_gold) != 60 {
		t.Errorf("char2 should have 60 gold, has %d", has_item(char2, item_gold))
	}
}

func TestVClaim(t *testing.T) {
	setupTransferTest()

	charID := 1001
	playerID := 2001
	locID := 3001
	regionID := 4001

	alloc_box(charID, T_char, 0)
	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(locID, T_loc, sub_plain)
	alloc_box(regionID, T_loc, sub_region)

	p_char(charID).unit_lord = playerID
	p_loc_info(charID).where = locID
	p_loc_info(locID).where = regionID

	// Set special regions to unique IDs so they won't match our test region
	cloud_region = 9001
	hades_region = 9002
	faery_region = 9003

	add_item(playerID, item_gold, 1000)

	c := &command{
		who: charID,
		a:   item_gold,
		b:   100,
		c:   0,
	}

	result := v_claim(c)
	if result != TRUE {
		t.Errorf("v_claim returned %d, want %d", result, TRUE)
	}

	if has_item(charID, item_gold) != 100 {
		t.Errorf("char should have 100 gold, has %d", has_item(charID, item_gold))
	}
	if has_item(playerID, item_gold) != 900 {
		t.Errorf("player should have 900 gold, has %d", has_item(playerID, item_gold))
	}
}

func TestVClaimAutoCorrect(t *testing.T) {
	setupTransferTest()

	charID := 1001
	playerID := 2001
	locID := 3001
	regionID := 4001

	alloc_box(charID, T_char, 0)
	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(locID, T_loc, sub_plain)
	alloc_box(regionID, T_loc, sub_region)

	p_char(charID).unit_lord = playerID
	p_loc_info(charID).where = locID
	p_loc_info(locID).where = regionID

	// Set special regions to unique IDs so they won't match our test region
	cloud_region = 9001
	hades_region = 9002
	faery_region = 9003

	add_item(playerID, item_gold, 1000)

	// CLAIM 500 should be auto-corrected to CLAIM 1 500 (claim 500 gold)
	c := &command{
		who: charID,
		a:   500,
		b:   0,
		c:   0,
	}

	result := v_claim(c)
	if result != TRUE {
		t.Errorf("v_claim with auto-correct returned %d, want %d", result, TRUE)
	}

	if has_item(charID, item_gold) != 500 {
		t.Errorf("char should have 500 gold after auto-correct, has %d", has_item(charID, item_gold))
	}
}

func TestCountManItems(t *testing.T) {
	setupTransferTest()

	charID := 1001
	alloc_box(charID, T_char, 0)

	add_item(charID, item_peasant, 10)
	add_item(charID, item_soldier, 5)
	add_item(charID, item_gold, 100)

	count := count_man_items(charID)
	if count != 15 {
		t.Errorf("count_man_items should be 15, got %d", count)
	}
}

func TestMyPrisoner(t *testing.T) {
	setupTransferTest()

	char1 := 1001
	char2 := 1002
	char3 := 1003
	player1 := 2001
	locID := 3001

	alloc_box(char1, T_char, 0)
	alloc_box(char2, T_char, 0)
	alloc_box(char3, T_char, 0)
	alloc_box(player1, T_player, sub_pl_regular)
	alloc_box(locID, T_loc, sub_plain)

	p_char(char1).unit_lord = player1
	p_char(char2).unit_lord = player1
	p_char(char3).unit_lord = player1

	p_loc_info(char1).where = locID
	p_loc_info(char2).where = char1 // prisoner is located at the captor
	p_char(char2).prisoner = 1
	p_loc_info(char3).where = locID
	p_loc_info(locID).here_list = []int{char1, char3}

	t.Run("prisoner under captor", func(t *testing.T) {
		if !my_prisoner(char1, char2) {
			t.Error("my_prisoner should return true for char2 under char1")
		}
	})

	t.Run("non-prisoner", func(t *testing.T) {
		if my_prisoner(char1, char3) {
			t.Error("my_prisoner should return false for non-prisoner")
		}
	})
}

func TestVGetGarrisonMinimum(t *testing.T) {
	setupTransferTest()

	charID := 1001
	garrID := 1002
	playerID := 2001
	locID := 3001
	castleID := 3002
	regionID := 4001

	alloc_box(charID, T_char, 0)
	alloc_box(garrID, T_char, sub_garrison)
	alloc_box(playerID, T_player, sub_pl_regular)
	alloc_box(locID, T_loc, sub_plain)
	alloc_box(castleID, T_loc, sub_castle)
	alloc_box(regionID, T_loc, sub_region)

	p_char(charID).unit_lord = playerID
	p_char(garrID).unit_lord = playerID

	p_loc_info(charID).where = castleID
	p_loc_info(garrID).where = castleID
	p_loc_info(castleID).where = locID
	p_loc_info(locID).where = regionID
	// charID must be first to be building owner (first_character)
	p_loc_info(castleID).here_list = []int{charID, garrID}
	p_loc_info(locID).here_list = []int{castleID}

	// Set up the garrison-castle relationship
	p_misc(garrID).garr_castle = castleID

	// Give garrison 15 soldiers
	add_item(garrID, item_soldier, 15)

	t.Run("cannot take leaving less than 10 men", func(t *testing.T) {
		c := &command{
			who: charID,
			a:   garrID,
			b:   item_soldier,
			c:   10, // try to take 10, would leave 5 which is < 10
			d:   0,
		}

		result := v_get(c)
		if result != FALSE {
			t.Errorf("v_get from garrison leaving < 10 men should fail, got %d", result)
		}
	})

	t.Run("can take leaving 10 men", func(t *testing.T) {
		// Reset: give garrison 15 soldiers again
		teg.globals.inventories[garrID] = []item_ent{{item: item_soldier, qty: 15}}

		c := &command{
			who: charID,
			a:   garrID,
			b:   item_soldier,
			c:   5, // take 5, leaves 10
			d:   0,
		}

		result := v_get(c)
		if result != TRUE {
			t.Errorf("v_get from garrison leaving 10 men should succeed, got %d", result)
		}
	})
}


