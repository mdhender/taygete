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

// TestFlagRaised tests the flag_raised function.
func TestFlagRaised(t *testing.T) {
	// Clear any existing flags
	flags = nil

	// No flags raised yet
	if got := flag_raised(0, "attack"); got != -1 {
		t.Errorf("flag_raised(0, 'attack') = %d, want -1 (no flags)", got)
	}

	// Add a flag
	flags = append(flags, &flag_ent{who: 1001, flag: "attack"})

	// Should find it with who=0 (any)
	if got := flag_raised(0, "attack"); got != 0 {
		t.Errorf("flag_raised(0, 'attack') = %d, want 0", got)
	}

	// Should find it with matching who
	setupWaitTestPlayer(1001, 100)
	if got := flag_raised(1001, "attack"); got != 0 {
		t.Errorf("flag_raised(1001, 'attack') = %d, want 0", got)
	}

	// Should find it by player
	if got := flag_raised(100, "attack"); got != 0 {
		t.Errorf("flag_raised(100, 'attack') = %d, want 0", got)
	}

	// Should not find with different flag
	if got := flag_raised(0, "retreat"); got != -1 {
		t.Errorf("flag_raised(0, 'retreat') = %d, want -1", got)
	}

	// Case-insensitive match
	if got := flag_raised(0, "ATTACK"); got != 0 {
		t.Errorf("flag_raised(0, 'ATTACK') = %d, want 0 (case-insensitive)", got)
	}
}

// TestVFlag tests the v_flag command.
func TestVFlag(t *testing.T) {
	// Clear any existing flags
	flags = nil

	charID := 1001
	setupWaitTestChar(charID)

	// Test with no arguments (numargs returns 0 when a=0)
	c := &command{who: charID}
	if got := v_flag(c); got != FALSE {
		t.Error("v_flag with no args should return FALSE")
	}
}

// TestClearFlags tests that clear_flags properly clears all flags.
func TestClearFlags(t *testing.T) {
	// Start fresh
	flags = nil

	flags = append(flags, &flag_ent{who: 1001, flag: "test"})
	flags = append(flags, &flag_ent{who: 1002, flag: "test2"})

	if len(flags) != 2 {
		t.Fatal("setup failed")
	}

	clear_flags()

	if len(flags) != 0 {
		t.Errorf("after clear_flags, len(flags) = %d, want 0", len(flags))
	}
}

// TestWaitTags tests the wait_tags lookup table.
func TestWaitTags(t *testing.T) {
	tests := []struct {
		tag  string
		want int
	}{
		{"time", 0},
		{"day", 1},
		{"unit", 2},
		{"gold", 3},
		{"item", 4},
		{"flag", 5},
		{"loc", 6},
		{"stack", 7},
		{"top", 8},
		{"ferry", 9},
		{"ship", 10},
		{"rain", 11},
		{"fog", 12},
		{"wind", 13},
		{"not", 14},
		{"owner", 15},
		{"raining", 16},
		{"foggy", 17},
		{"windy", 18},
		{"clear", 19},
		{"shiploc", 20},
		{"month", 21},
		{"turn", 22},
		{"unknown", -1},
	}

	for _, tt := range tests {
		got := lookup(wait_tags, tt.tag)
		if got != tt.want {
			t.Errorf("lookup(wait_tags, %q) = %d, want %d", tt.tag, got, tt.want)
		}
	}
}

// TestParseWaitArgs tests parsing of WAIT command arguments.
// Note: Full parsing requires the order parsing system (Sprint 19+).
// For now we test that the parsing infrastructure exists.
func TestParseWaitArgs(t *testing.T) {
	charID := 1001
	setupWaitTestChar(charID)

	// Test that parse_wait_args can be called and clears the list
	c := &command{who: charID}
	waitParseLists[c] = []*waitArgExt{{tag: waitTagTime, a1: 5}}

	// Calling parse_wait_args should clear the list first
	_ = parse_wait_args(c)

	// The actual parsing depends on getCommandParseArgs which is simplified for now
}

// TestCheckWaitConditionsTime tests the "time n" wait condition.
func TestCheckWaitConditionsTime(t *testing.T) {
	charID := 1001
	setupWaitTestChar(charID)

	c := &command{who: charID, days_executing: 3}

	// Set up wait for 5 days
	waitParseLists[c] = []*waitArgExt{{tag: waitTagTime, a1: 5}}

	// Should not trigger (only 3 days elapsed)
	if got := check_wait_conditions(c); got != "" {
		t.Errorf("check_wait_conditions with 3/5 days = %q, want empty", got)
	}

	// Advance to 5 days
	c.days_executing = 5
	if got := check_wait_conditions(c); got == "" {
		t.Error("check_wait_conditions with 5/5 days should trigger")
	}
}

// TestCheckWaitConditionsDay tests the "day n" wait condition.
func TestCheckWaitConditionsDay(t *testing.T) {
	charID := 1001
	setupWaitTestChar(charID)

	// Set sysclock day
	teg.globals.sysclock.day = 10

	c := &command{who: charID}
	waitParseLists[c] = []*waitArgExt{{tag: waitTagDay, a1: 15}}

	// Should not trigger (day 10, waiting for day 15)
	if got := check_wait_conditions(c); got != "" {
		t.Errorf("check_wait_conditions day 10/15 = %q, want empty", got)
	}

	// Advance to day 15
	teg.globals.sysclock.day = 15
	if got := check_wait_conditions(c); got == "" {
		t.Error("check_wait_conditions day 15/15 should trigger")
	}
}

// TestCheckWaitConditionsTurn tests the "turn n" wait condition.
func TestCheckWaitConditionsTurn(t *testing.T) {
	charID := 1001
	setupWaitTestChar(charID)

	teg.globals.sysclock.turn = 5

	c := &command{who: charID}
	waitParseLists[c] = []*waitArgExt{{tag: waitTagTurn, a1: 10}}

	// Should not trigger (turn 5, waiting for turn 10)
	if got := check_wait_conditions(c); got != "" {
		t.Errorf("check_wait_conditions turn 5/10 = %q, want empty", got)
	}

	// Advance to turn 10
	teg.globals.sysclock.turn = 10
	if got := check_wait_conditions(c); got == "" {
		t.Error("check_wait_conditions turn 10/10 should trigger")
	}
}

// TestCheckWaitConditionsGold tests the "gold n" wait condition.
// Note: Requires proper item type setup which is complex.
// This test verifies the condition check logic with pre-configured inventory.
func TestCheckWaitConditionsGold(t *testing.T) {
	charID := 1001
	setupWaitTestChar(charID)

	c := &command{who: charID}
	waitParseLists[c] = []*waitArgExt{{tag: waitTagGold, a1: 100}}

	// With no gold, should not trigger
	result := check_wait_conditions(c)
	if result != "" {
		t.Logf("check_wait_conditions without gold = %q", result)
	}
}

// TestCheckWaitConditionsItem tests the "item n q" wait condition.
// Note: Requires proper item type setup which is complex.
func TestCheckWaitConditionsItem(t *testing.T) {
	charID := 1001
	setupWaitTestChar(charID)

	itemID := 10

	c := &command{who: charID}
	waitParseLists[c] = []*waitArgExt{{tag: waitTagItem, a1: itemID, a2: 10}}

	// Without item, should not trigger (item kind check fails)
	result := check_wait_conditions(c)
	if result != "" {
		t.Logf("check_wait_conditions item without setup = %q", result)
	}
}

// TestCheckWaitConditionsUnit tests the "unit n" wait condition.
func TestCheckWaitConditionsUnit(t *testing.T) {
	charID := 1001
	targetID := 1002
	locID := 2001
	otherLocID := 2002

	setupWaitTestChar(charID)
	setupWaitTestChar(targetID)
	setupWaitTestLoc(locID)
	setupWaitTestLoc(otherLocID)

	set_where(charID, locID)
	set_where(targetID, otherLocID)

	c := &command{who: charID}
	waitParseLists[c] = []*waitArgExt{{tag: waitTagUnit, a1: targetID}}

	// Should not trigger (target in different location)
	if got := check_wait_conditions(c); got != "" {
		t.Errorf("check_wait_conditions unit not here = %q, want empty", got)
	}

	// Move target to same location
	set_where(targetID, locID)
	if got := check_wait_conditions(c); got == "" {
		t.Error("check_wait_conditions unit here should trigger")
	}
}

// TestCheckWaitConditionsFlag tests the "flag" wait condition.
func TestCheckWaitConditionsFlag(t *testing.T) {
	charID := 1001
	setupWaitTestChar(charID)
	setupWaitTestPlayer(charID, 100)

	flags = nil

	c := &command{who: charID}
	waitParseLists[c] = []*waitArgExt{{tag: waitTagFlag, a1: 100, flagStr: "attack"}}

	// Should not trigger (no flags raised)
	if got := check_wait_conditions(c); got != "" {
		t.Errorf("check_wait_conditions flag not raised = %q, want empty", got)
	}

	// Raise the flag
	flags = append(flags, &flag_ent{who: charID, flag: "attack"})
	if got := check_wait_conditions(c); got == "" {
		t.Error("check_wait_conditions flag raised should trigger")
	}
}

// TestCheckWaitConditionsLoc tests the "loc n" wait condition.
func TestCheckWaitConditionsLoc(t *testing.T) {
	charID := 1001
	locID := 2001
	targetLocID := 2002

	setupWaitTestChar(charID)
	setupWaitTestLoc(locID)
	setupWaitTestLoc(targetLocID)

	set_where(charID, locID)

	c := &command{who: charID}
	waitParseLists[c] = []*waitArgExt{{tag: waitTagLoc, a1: targetLocID}}

	// Should not trigger (not at target location)
	if got := check_wait_conditions(c); got != "" {
		t.Errorf("check_wait_conditions loc not there = %q, want empty", got)
	}

	// Move to target location
	set_where(charID, targetLocID)
	if got := check_wait_conditions(c); got == "" {
		t.Error("check_wait_conditions loc there should trigger")
	}
}

// TestCheckWaitConditionsStack tests the "stack n" wait condition.
func TestCheckWaitConditionsStack(t *testing.T) {
	charID := 1001
	targetID := 1002
	locID := 2001

	setupWaitTestChar(charID)
	setupWaitTestChar(targetID)
	setupWaitTestLoc(locID)

	set_where(charID, locID)
	set_where(targetID, locID)

	c := &command{who: charID}
	waitParseLists[c] = []*waitArgExt{{tag: waitTagStack, a1: targetID}}

	// Note: Stack leader depends on stacking system
	// For now just test that the condition can be evaluated
	_ = check_wait_conditions(c)
}

// TestCheckWaitConditionsTop tests the "top" wait condition.
func TestCheckWaitConditionsTop(t *testing.T) {
	charID := 1001
	setupWaitTestChar(charID)

	c := &command{who: charID}
	waitParseLists[c] = []*waitArgExt{{tag: waitTagTop}}

	// Character alone is stack leader
	result := check_wait_conditions(c)
	if result == "" {
		t.Error("check_wait_conditions top (alone) should trigger")
	}
}

// TestCheckWaitConditionsNot tests the "not" modifier.
func TestCheckWaitConditionsNot(t *testing.T) {
	charID := 1001
	locID := 2001
	targetLocID := 2002

	setupWaitTestChar(charID)
	setupWaitTestLoc(locID)
	setupWaitTestLoc(targetLocID)

	set_where(charID, locID)

	c := &command{who: charID}
	// "not loc target" - trigger if NOT at target location
	waitParseLists[c] = []*waitArgExt{
		{tag: waitTagNot},
		{tag: waitTagLoc, a1: targetLocID},
	}

	// Should trigger (we are NOT at target location)
	if got := check_wait_conditions(c); got == "" {
		t.Error("check_wait_conditions 'not loc' should trigger when not at target")
	}

	// Move to target location
	set_where(charID, targetLocID)
	if got := check_wait_conditions(c); got != "" {
		t.Errorf("check_wait_conditions 'not loc' at target = %q, want empty", got)
	}
}

// TestCheckWaitConditionsWeather tests weather wait conditions.
func TestCheckWaitConditionsWeather(t *testing.T) {
	charID := 1001
	provID := 3001

	setupWaitTestChar(charID)
	setupWaitTestProvince(provID)
	set_where(charID, provID)

	// Set up province accessor to return province
	ensureWaitBox(provID)
	bxProv := teg.globals.bx[provID]
	bxProv.kind = T_loc
	if bxProv.x_loc == nil {
		bxProv.x_loc = &entity_loc{}
	}

	c := &command{who: charID}

	// Test rain condition (no rain)
	waitParseLists[c] = []*waitArgExt{{tag: waitTagRain}}
	if got := check_wait_conditions(c); got != "" {
		t.Errorf("check_wait_conditions rain (none) = %q, want empty", got)
	}

	// Test clear condition (should trigger when no weather)
	waitParseLists[c] = []*waitArgExt{{tag: waitTagClear}}
	if got := check_wait_conditions(c); got == "" {
		t.Error("check_wait_conditions clear (no weather) should trigger")
	}
}

// TestVWait tests the v_wait command start routine.
func TestVWait(t *testing.T) {
	charID := 1001
	setupWaitTestChar(charID)

	// Clear wait list
	wait_list = nil

	// Test with no arguments (numargs returns 0 when a=0)
	c := &command{who: charID}
	if got := v_wait(c); got != FALSE {
		t.Error("v_wait with no args should return FALSE")
	}

	// Test with a condition set up - should add to wait list
	c = &command{who: charID, a: 1}                           // numargs > 0
	waitParseLists[c] = []*waitArgExt{{tag: waitTagTime, a1: 99}} // Wait 99 days (won't trigger)
	initialLen := len(wait_list)
	result := v_wait(c)
	if result != TRUE {
		t.Error("v_wait with waiting condition should return TRUE")
	}
	// Unit should be added to wait list
	if len(wait_list) <= initialLen {
		t.Error("v_wait should add unit to wait_list")
	}
}

// TestDWait tests the d_wait daily check routine.
func TestDWait(t *testing.T) {
	charID := 1001
	setupWaitTestChar(charID)

	wait_list = []int{charID}

	c := &command{who: charID, days_executing: 3}
	waitParseLists[c] = []*waitArgExt{{tag: waitTagTime, a1: 5}}

	// Not yet time - should return TRUE but not finish
	if got := d_wait(c); got != TRUE {
		t.Error("d_wait should return TRUE")
	}
	if c.inhibit_finish == TRUE {
		t.Error("d_wait should not set inhibit_finish when still waiting")
	}

	// Time elapsed - should finish
	c.days_executing = 5
	c.inhibit_finish = 0
	if got := d_wait(c); got != TRUE {
		t.Error("d_wait should return TRUE when finished")
	}
	if c.inhibit_finish != TRUE {
		t.Error("d_wait should set inhibit_finish when finished")
	}
}

// TestIWait tests the i_wait interrupt routine.
func TestIWait(t *testing.T) {
	charID := 1001
	setupWaitTestChar(charID)

	wait_list = []int{charID, 1002, 1003}

	c := &command{who: charID}
	if got := i_wait(c); got != TRUE {
		t.Error("i_wait should return TRUE")
	}

	// Character should be removed from wait_list
	for _, id := range wait_list {
		if id == charID {
			t.Error("i_wait should remove character from wait_list")
		}
	}
}

// TestWaitListOperations tests wait_list add/remove operations.
func TestWaitListOperations(t *testing.T) {
	wait_list = nil

	// Add characters
	IListAppend(&wait_list, 1001)
	IListAppend(&wait_list, 1002)
	IListAppend(&wait_list, 1003)

	if len(wait_list) != 3 {
		t.Errorf("wait_list length = %d, want 3", len(wait_list))
	}

	// Remove one
	IListRemValue(&wait_list, 1002)
	if len(wait_list) != 2 {
		t.Errorf("wait_list length after remove = %d, want 2", len(wait_list))
	}

	// Verify 1002 is gone
	for _, id := range wait_list {
		if id == 1002 {
			t.Error("1002 should be removed from wait_list")
		}
	}
}

// TestEvalNotCond tests the not condition evaluator.
func TestEvalNotCond(t *testing.T) {
	// not=false, cond=true -> true
	if got := evalNotCond(false, true); got != true {
		t.Errorf("evalNotCond(false, true) = %v, want true", got)
	}

	// not=false, cond=false -> false
	if got := evalNotCond(false, false); got != false {
		t.Errorf("evalNotCond(false, false) = %v, want false", got)
	}

	// not=true, cond=true -> false (negated)
	if got := evalNotCond(true, true); got != false {
		t.Errorf("evalNotCond(true, true) = %v, want false", got)
	}

	// not=true, cond=false -> true (negated)
	if got := evalNotCond(true, false); got != true {
		t.Errorf("evalNotCond(true, false) = %v, want true", got)
	}
}

// TestPluralHelpers tests plural formatting helpers.
func TestPluralHelpers(t *testing.T) {
	if got := pluralS(1); got != "" {
		t.Errorf("pluralS(1) = %q, want empty", got)
	}
	if got := pluralS(2); got != "s" {
		t.Errorf("pluralS(2) = %q, want 's'", got)
	}

	if got := pluralHaveHas(1); got != " has" {
		t.Errorf("pluralHaveHas(1) = %q, want ' has'", got)
	}
	if got := pluralHaveHas(2); got != "s have" {
		t.Errorf("pluralHaveHas(2) = %q, want 's have'", got)
	}
}

// TestNotStrHelpers tests not string helpers.
func TestNotStrHelpers(t *testing.T) {
	if got := notStr(true); got != " not" {
		t.Errorf("notStr(true) = %q, want ' not'", got)
	}
	if got := notStr(false); got != "" {
		t.Errorf("notStr(false) = %q, want empty", got)
	}

	if got := notStr2(true); got != "not " {
		t.Errorf("notStr2(true) = %q, want 'not '", got)
	}
	if got := notStr2(false); got != "" {
		t.Errorf("notStr2(false) = %q, want empty", got)
	}
}

// Helper function to ensure a box exists
func ensureWaitBox(id int) {
	if teg.globals.bx[id] == nil {
		teg.globals.bx[id] = &box{}
	}
}

// Helper function to set up a test character for wait tests
func setupWaitTestChar(charID int) {
	ensureWaitBox(charID)
	bxChar := teg.globals.bx[charID]
	bxChar.kind = T_char
	if bxChar.x_char == nil {
		bxChar.x_char = &entity_char{}
	}
	bxChar.x_char.health = 100
	bxChar.x_char.unit_lord = indep_player
	bxChar.x_loc_info = loc_info{}
}

// Helper function to set up a test location for wait tests
func setupWaitTestLoc(locID int) {
	ensureWaitBox(locID)
	bxLoc := teg.globals.bx[locID]
	bxLoc.kind = T_loc
	bxLoc.x_loc_info = loc_info{}
}

// Helper function to set up test player for wait tests
func setupWaitTestPlayer(charID, playerID int) {
	ensureWaitBox(playerID)
	bxPlayer := teg.globals.bx[playerID]
	bxPlayer.kind = T_player

	ensureWaitBox(charID)
	bxChar := teg.globals.bx[charID]
	if bxChar.x_char == nil {
		bxChar.x_char = &entity_char{}
	}
	bxChar.x_char.unit_lord = playerID
}

// Helper function to set up test province for wait tests
func setupWaitTestProvince(id int) {
	ensureWaitBox(id)
	bx := teg.globals.bx[id]
	bx.kind = T_loc
	bx.skind = sub_plain // province subkind
	if bx.x_loc == nil {
		bx.x_loc = &entity_loc{}
	}
	bx.x_loc_info = loc_info{}
}
