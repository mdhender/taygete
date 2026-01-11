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

func setupMetaTest() {
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

func TestMayName(t *testing.T) {
	setupMetaTest()

	playerID := 100
	charID := 1001
	otherCharID := 1002
	otherPlayerID := 101

	alloc_box(playerID, T_player, 0)
	alloc_box(charID, T_char, 0)
	alloc_box(otherCharID, T_char, 0)
	alloc_box(otherPlayerID, T_player, 0)

	p_char(charID).unit_lord = playerID
	p_char(otherCharID).unit_lord = otherPlayerID

	t.Run("character can name itself", func(t *testing.T) {
		if !may_name(charID, charID) {
			t.Error("character should be able to name itself")
		}
	})

	t.Run("character can name own player", func(t *testing.T) {
		if !may_name(charID, playerID) {
			t.Error("character should be able to name its player")
		}
	})

	t.Run("character cannot name other player's character", func(t *testing.T) {
		if may_name(charID, otherCharID) {
			t.Error("character should not be able to name another player's character")
		}
	})

	t.Run("character cannot name other player", func(t *testing.T) {
		if may_name(charID, otherPlayerID) {
			t.Error("character should not be able to name another player")
		}
	})
}

func TestMayNameItem(t *testing.T) {
	setupMetaTest()

	playerID := 100
	charID := 1001
	itemID := 2000
	otherItemID := 2001

	alloc_box(playerID, T_player, 0)
	alloc_box(charID, T_char, 0)
	alloc_box(itemID, T_item, 0)
	alloc_box(otherItemID, T_item, 0)

	p_char(charID).unit_lord = playerID

	t.Run("character can name item they created", func(t *testing.T) {
		p_item_magic(itemID).creator = charID
		if !may_name(charID, itemID) {
			t.Error("character should be able to name item they created")
		}
	})

	t.Run("character cannot name item created by others", func(t *testing.T) {
		p_item_magic(otherItemID).creator = 9999
		if may_name(charID, otherItemID) {
			t.Error("character should not be able to name item created by others")
		}
	})

	t.Run("character can name potion items", func(t *testing.T) {
		p_item_magic(itemID).creator = 0
		p_item_magic(itemID).use_key = use_heal_potion
		if !may_name(charID, itemID) {
			t.Error("character should be able to name heal potion")
		}
	})
}

func TestMayNameStorm(t *testing.T) {
	setupMetaTest()

	charID := 1001
	stormID := 3000

	alloc_box(charID, T_char, 0)
	alloc_box(stormID, T_storm, 0)

	t.Run("character can name storm they summoned", func(t *testing.T) {
		p_misc(stormID).summoned_by = charID
		if !may_name(charID, stormID) {
			t.Error("character should be able to name storm they summoned")
		}
	})

	t.Run("character cannot name storm summoned by others", func(t *testing.T) {
		p_misc(stormID).summoned_by = 9999
		if may_name(charID, stormID) {
			t.Error("character should not be able to name storm summoned by others")
		}
	})
}

func TestCmdShift(t *testing.T) {
	c := &command{
		a: 1,
		b: 2,
		c: 3,
		d: 4,
		e: 5,
		f: 6,
		g: 7,
		h: 8,
	}

	cmd_shift(c)

	if c.a != 2 || c.b != 3 || c.c != 4 || c.d != 5 ||
		c.e != 6 || c.f != 7 || c.g != 8 || c.h != 0 {
		t.Errorf("cmd_shift did not shift correctly: got a=%d b=%d c=%d d=%d e=%d f=%d g=%d h=%d",
			c.a, c.b, c.c, c.d, c.e, c.f, c.g, c.h)
	}
}

func TestCmdNumargsFull(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		c := &command{}
		if got := cmd_numargs_full(c); got != 0 {
			t.Errorf("cmd_numargs_full = %d, want 0", got)
		}
	})

	t.Run("one arg", func(t *testing.T) {
		c := &command{a: 1}
		if got := cmd_numargs_full(c); got != 1 {
			t.Errorf("cmd_numargs_full = %d, want 1", got)
		}
	})

	t.Run("three args", func(t *testing.T) {
		c := &command{a: 1, b: 2, c: 3}
		if got := cmd_numargs_full(c); got != 3 {
			t.Errorf("cmd_numargs_full = %d, want 3", got)
		}
	})
}

func TestVLook(t *testing.T) {
	setupMetaTest()

	charID := 1001
	locID := 5000

	alloc_box(charID, T_char, 0)
	alloc_box(locID, T_loc, sub_plain)
	set_where(charID, locID)

	c := &command{who: charID}

	t.Run("character can look", func(t *testing.T) {
		result := v_look(c)
		if result != TRUE {
			t.Errorf("v_look = %d, want TRUE", result)
		}
	})

	t.Run("non-character cannot look", func(t *testing.T) {
		c.who = locID
		result := v_look(c)
		if result != FALSE {
			t.Errorf("v_look for non-char = %d, want FALSE", result)
		}
	})
}

func TestVStop(t *testing.T) {
	setupMetaTest()

	charID := 1001
	alloc_box(charID, T_char, 0)

	c := &command{who: charID}

	result := v_stop(c)
	if result != TRUE {
		t.Errorf("v_stop = %d, want TRUE", result)
	}
}

func TestVPublic(t *testing.T) {
	setupMetaTest()

	playerID := 100
	charID := 1001

	alloc_box(playerID, T_player, 0)
	alloc_box(charID, T_char, 0)
	p_char(charID).unit_lord = playerID

	c := &command{who: charID}

	t.Run("first public succeeds and grants gold", func(t *testing.T) {
		result := v_public(c)
		if result != TRUE {
			t.Errorf("v_public = %d, want TRUE", result)
		}

		if player_public_turn(playerID) == 0 {
			t.Error("player public_turn not set")
		}

		goldQty := has_item(playerID, item_gold)
		if goldQty != 100 {
			t.Errorf("gold = %d, want 100", goldQty)
		}
	})

	t.Run("second public fails", func(t *testing.T) {
		result := v_public(c)
		if result != FALSE {
			t.Errorf("second v_public = %d, want FALSE", result)
		}
	})
}

func TestVEmote(t *testing.T) {
	setupMetaTest()

	charID := 1001
	targetID := 1002

	alloc_box(charID, T_char, 0)
	alloc_box(targetID, T_char, 0)

	t.Run("emote with insufficient args fails", func(t *testing.T) {
		c := &command{who: charID, a: targetID}
		result := v_emote(c)
		if result != FALSE {
			t.Errorf("v_emote with 1 arg = %d, want FALSE", result)
		}
	})

	t.Run("emote with target and message succeeds", func(t *testing.T) {
		c := &command{who: charID, a: targetID, b: 1}
		result := v_emote(c)
		if result != TRUE {
			t.Errorf("v_emote = %d, want TRUE", result)
		}
	})
}

func TestSetNameAndBanner(t *testing.T) {
	setupMetaTest()

	charID := 1001
	alloc_box(charID, T_char, 0)

	t.Run("set name", func(t *testing.T) {
		set_name(charID, "Test Character")
		if got := teg.getName(charID); got != "Test Character" {
			t.Errorf("getName = %q, want %q", got, "Test Character")
		}
	})

	t.Run("clear name", func(t *testing.T) {
		set_name(charID, "")
		if got := teg.getName(charID); got != "" {
			t.Errorf("getName after clear = %q, want empty", got)
		}
	})

	t.Run("set banner", func(t *testing.T) {
		set_banner(charID, "A mighty warrior")
		if got := teg.getBanner(charID); got != "A mighty warrior" {
			t.Errorf("getBanner = %q, want %q", got, "A mighty warrior")
		}
	})

	t.Run("clear banner", func(t *testing.T) {
		set_banner(charID, "")
		if got := teg.getBanner(charID); got != "" {
			t.Errorf("getBanner after clear = %q, want empty", got)
		}
	})
}

func TestHasAuraculum(t *testing.T) {
	setupMetaTest()

	charID := 1001
	auraID := 2001

	alloc_box(charID, T_char, 0)
	alloc_box(auraID, T_item, 0)

	t.Run("no auraculum returns 0", func(t *testing.T) {
		if got := has_auraculum(charID); got != 0 {
			t.Errorf("has_auraculum = %d, want 0", got)
		}
	})

	t.Run("auraculum set but not held returns 0", func(t *testing.T) {
		p_magic(charID).auraculum = auraID
		if got := has_auraculum(charID); got != 0 {
			t.Errorf("has_auraculum without item = %d, want 0", got)
		}
	})

	t.Run("auraculum set and held returns ID", func(t *testing.T) {
		p_magic(charID).auraculum = auraID
		teg.globals.inventories[charID] = []item_ent{{item: auraID, qty: 1}}
		if got := has_auraculum(charID); got != auraID {
			t.Errorf("has_auraculum with item = %d, want %d", got, auraID)
		}
	})
}
