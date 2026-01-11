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

// TestImmediate tests the immediate command infrastructure.
func TestImmediate(t *testing.T) {
	e := teg

	t.Run("v_add_item adds items", func(t *testing.T) {
		// Create a test character
		charID := 1001
		e.globals.bx[charID] = &box{kind: T_char}
		e.globals.bx[charID].x_char = &entity_char{}

		// Set up item_gold as T_item so Kind check passes
		e.globals.bx[item_gold] = &box{kind: T_item}
		e.globals.bx[item_gold].x_item = &entity_item{}

		c := &command{
			who: charID,
			a:   item_gold,
			b:   100,
		}

		result := e.v_add_item(c)
		if result != TRUE {
			t.Errorf("v_add_item returned %d, expected TRUE", result)
		}

		// Check inventory
		inv := e.globals.inventories[charID]
		found := false
		for _, ie := range inv {
			if ie.item == item_gold && ie.qty == 100 {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected 100 gold in inventory, got %v", inv)
		}
	})

	t.Run("v_sub_item removes items", func(t *testing.T) {
		charID := 1002
		e.globals.bx[charID] = &box{kind: T_char}
		e.globals.bx[charID].x_char = &entity_char{}

		// Add initial items
		e.gen_item(charID, item_gold, 200)

		c := &command{
			who: charID,
			a:   item_gold,
			b:   50,
		}

		result := e.v_sub_item(c)
		if result != TRUE {
			t.Errorf("v_sub_item returned %d, expected TRUE", result)
		}

		// Check inventory - should have 150 gold left
		inv := e.globals.inventories[charID]
		for _, ie := range inv {
			if ie.item == item_gold {
				if ie.qty != 150 {
					t.Errorf("Expected 150 gold, got %d", ie.qty)
				}
				break
			}
		}
	})

	t.Run("v_see_all toggles visibility", func(t *testing.T) {
		charID := 1003
		e.globals.bx[charID] = &box{kind: T_char}

		c := &command{
			who: charID,
			a:   1,
		}

		result := e.v_see_all(c)
		if result != TRUE {
			t.Errorf("v_see_all returned %d, expected TRUE", result)
		}
		if !e.globals.immedSeeAll {
			t.Error("Expected immedSeeAll to be true")
		}

		// Toggle off
		c.a = 0
		e.v_see_all(c)
		// Note: The logic sets to true if a==0, matching C behavior
	})

	t.Run("v_be validates box", func(t *testing.T) {
		charID := 1004
		e.globals.bx[charID] = &box{kind: T_char}

		c := &command{
			who: charID,
			a:   charID, // valid box
		}

		result := e.v_be(c)
		if result != TRUE {
			t.Errorf("v_be returned %d for valid box, expected TRUE", result)
		}

		// Try with invalid box
		c.a = 999999
		result = e.v_be(c)
		if result != FALSE {
			t.Errorf("v_be returned %d for invalid box, expected FALSE", result)
		}
	})

	t.Run("v_xyzzy returns TRUE", func(t *testing.T) {
		c := &command{who: 1}
		result := e.v_xyzzy(c)
		if result != TRUE {
			t.Errorf("v_xyzzy returned %d, expected TRUE", result)
		}
	})

	t.Run("v_plugh returns TRUE", func(t *testing.T) {
		c := &command{who: 1}
		result := e.v_plugh(c)
		if result != TRUE {
			t.Errorf("v_plugh returned %d, expected TRUE", result)
		}
	})

	t.Run("gen_item and consume_item work correctly", func(t *testing.T) {
		charID := 1005
		e.globals.bx[charID] = &box{kind: T_char}
		e.globals.bx[charID].x_char = &entity_char{}

		// Add items
		e.gen_item(charID, item_soldier, 50)
		e.gen_item(charID, item_soldier, 25) // Should stack to 75

		inv := e.globals.inventories[charID]
		var soldierQty int
		for _, ie := range inv {
			if ie.item == item_soldier {
				soldierQty = ie.qty
				break
			}
		}
		if soldierQty != 75 {
			t.Errorf("Expected 75 soldiers, got %d", soldierQty)
		}

		// Consume some
		consumed := e.consume_item(charID, item_soldier, 30)
		if consumed != 30 {
			t.Errorf("Expected to consume 30, got %d", consumed)
		}

		// Check remaining
		inv = e.globals.inventories[charID]
		for _, ie := range inv {
			if ie.item == item_soldier {
				if ie.qty != 45 {
					t.Errorf("Expected 45 soldiers remaining, got %d", ie.qty)
				}
				break
			}
		}
	})

	t.Run("v_credit adds items to target", func(t *testing.T) {
		charID := 1006
		e.globals.bx[charID] = &box{kind: T_char}
		e.globals.bx[charID].x_char = &entity_char{}

		c := &command{
			who: 1,       // GM
			a:   charID,  // target
			b:   500,     // amount
			c:   item_gold, // item
		}

		result := e.v_credit(c)
		if result != TRUE {
			t.Errorf("v_credit returned %d, expected TRUE", result)
		}

		// Check that target received items
		inv := e.globals.inventories[charID]
		found := false
		for _, ie := range inv {
			if ie.item == item_gold && ie.qty >= 500 {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected target to have 500 gold, got %v", inv)
		}
	})

	t.Run("v_save returns TRUE", func(t *testing.T) {
		c := &command{who: 1}
		result := e.v_save(c)
		if result != TRUE {
			t.Errorf("v_save returned %d, expected TRUE", result)
		}
	})

	t.Run("v_listcmds returns TRUE", func(t *testing.T) {
		c := &command{who: 1}
		result := e.v_listcmds(c)
		if result != TRUE {
			t.Errorf("v_listcmds returned %d, expected TRUE", result)
		}
	})

	t.Run("v_dump returns FALSE for invalid box", func(t *testing.T) {
		c := &command{
			who: 1,
			a:   999999, // invalid
		}
		result := e.v_dump(c)
		if result != FALSE {
			t.Errorf("v_dump returned %d for invalid box, expected FALSE", result)
		}
	})

	t.Run("v_dump returns TRUE for valid box", func(t *testing.T) {
		charID := 1007
		e.globals.bx[charID] = &box{kind: T_char}

		c := &command{
			who: 1,
			a:   charID,
		}
		result := e.v_dump(c)
		if result != TRUE {
			t.Errorf("v_dump returned %d for valid box, expected TRUE", result)
		}
	})

	t.Run("v_know rejects non-skill", func(t *testing.T) {
		charID := 1008
		e.globals.bx[charID] = &box{kind: T_char}

		c := &command{
			who: charID,
			a:   charID, // not a skill
		}
		result := e.v_know(c)
		if result != FALSE {
			t.Errorf("v_know returned %d for non-skill, expected FALSE", result)
		}
	})

	t.Run("v_relore rejects non-skill", func(t *testing.T) {
		c := &command{
			who: 1,
			a:   item_gold, // not a skill
		}
		result := e.v_relore(c)
		if result != FALSE {
			t.Errorf("v_relore returned %d for non-skill, expected FALSE", result)
		}
	})
}

// TestImmediateMode tests the ImmediateMode method.
func TestImmediateMode(t *testing.T) {
	e := teg

	t.Run("returns false for nil command", func(t *testing.T) {
		// Use invalid who
		result := e.ImmediateMode(0, "test")
		if result {
			t.Error("Expected false for invalid who")
		}
	})

	t.Run("returns false for empty line", func(t *testing.T) {
		charID := 2001
		e.globals.bx[charID] = &box{kind: T_char}
		e.globals.bx[charID].x_char = &entity_char{}
		e.globals.bx[charID].cmd = nil

		result := e.ImmediateMode(charID, "")
		if result {
			t.Error("Expected false for empty command line")
		}
	})
}
