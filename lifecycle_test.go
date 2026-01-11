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

// setupTestCharacter creates a basic test character with the given ID.
func setupTestCharacter(who int, health schar) {
	if teg.globals.bx[who] == nil {
		teg.globals.bx[who] = &box{}
	}
	teg.globals.bx[who].kind = T_char
	teg.globals.bx[who].skind = 0
	teg.globals.bx[who].x_char = &entity_char{
		health:    health,
		melt_me:   FALSE,
		prisoner:  FALSE,
		unit_lord: indep_player,
	}
	teg.globals.bx[who].x_loc_info = loc_info{where: 0}
	teg.setName(who, "Test Character")
}

// setupTestLocation creates a basic test location.
func setupTestLocation(loc int, sk schar) {
	if teg.globals.bx[loc] == nil {
		teg.globals.bx[loc] = &box{}
	}
	teg.globals.bx[loc].kind = T_loc
	teg.globals.bx[loc].skind = sk
	teg.globals.bx[loc].x_loc_info = loc_info{where: 0}
	teg.setName(loc, "Test Location")
}

// setupTestPlayer creates a basic test player.
func setupTestPlayer(pl int) {
	if teg.globals.bx[pl] == nil {
		teg.globals.bx[pl] = &box{}
	}
	teg.globals.bx[pl].kind = T_player
	teg.globals.bx[pl].skind = sub_pl_regular
	teg.globals.bx[pl].x_player = &entity_player{}
	teg.setName(pl, "Test Player")
}

// TestSurviveFatalWithoutSkill tests that survive_fatal returns false when skill is missing.
func TestSurviveFatalWithoutSkill(t *testing.T) {
	who := 5001
	setupTestCharacter(who, 100)

	result := survive_fatal(who)

	if result {
		t.Error("survive_fatal should return false when character lacks sk_survive_fatal")
	}
}

// TestSurviveFatalWithSkill tests that survive_fatal returns true and heals when skill is present.
func TestSurviveFatalWithSkill(t *testing.T) {
	who := 5002
	setupTestCharacter(who, 0)

	if teg.globals.charSkills == nil {
		teg.globals.charSkills = make(map[int][]*skill_ent)
	}
	teg.globals.charSkills[who] = []*skill_ent{
		{skill: sk_survive_fatal, know: SKILL_know},
	}

	result := survive_fatal(who)

	if !result {
		t.Error("survive_fatal should return true when character has sk_survive_fatal")
	}

	if p_char(who).health != 100 {
		t.Errorf("health = %d, want 100", p_char(who).health)
	}

	if p_char(who).sick != FALSE {
		t.Error("sick should be FALSE after survive_fatal")
	}

	skills := teg.globals.charSkills[who]
	for _, s := range skills {
		if s.skill == sk_survive_fatal && s.know == SKILL_know {
			t.Error("sk_survive_fatal should be forgotten after use")
		}
	}
}

// TestCharReclaimLifecycle tests that char_reclaim sets melt_me and triggers kill_char.
func TestCharReclaimLifecycle(t *testing.T) {
	who := 5003
	loc := 10001
	setupTestLocation(loc, sub_plain)
	setupTestCharacter(who, 100)
	set_where(who, loc)

	char_reclaim(who)

	if kind(who) != T_deadchar {
		t.Errorf("kind = %d, want T_deadchar (%d)", kind(who), T_deadchar)
	}
}

// TestStackmateInheritorNoStackmates tests that stackmate_inheritor returns 0 when alone.
func TestStackmateInheritorNoStackmates(t *testing.T) {
	who := 5004
	setupTestCharacter(who, 100)

	result := stackmate_inheritor(who)

	if result != 0 {
		t.Errorf("stackmate_inheritor = %d, want 0 (no stackmates)", result)
	}
}

// TestStackmateInheritorWithStackmate tests finding an inheritor from stack.
func TestStackmateInheritorWithStackmate(t *testing.T) {
	who := 5005
	stackmate := 5006
	setupTestCharacter(who, 100)
	setupTestCharacter(stackmate, 100)

	// Manually set up the stack relationship without using set_where to avoid duplicate add
	if teg.globals.bx[who].x_loc_info.here_list == nil {
		teg.globals.bx[who].x_loc_info.here_list = []int{}
	}
	teg.globals.bx[who].x_loc_info.here_list = append(teg.globals.bx[who].x_loc_info.here_list, stackmate)
	teg.globals.bx[stackmate].x_loc_info.where = who

	result := stackmate_inheritor(who)

	if result != stackmate {
		t.Errorf("stackmate_inheritor = %d, want %d", result, stackmate)
	}
}

// TestAddCharDamageReducesHealth tests that damage reduces health.
func TestAddCharDamageReducesHealth(t *testing.T) {
	who := 5007
	loc := 10002
	setupTestLocation(loc, sub_plain)
	setupTestCharacter(who, 100)
	set_where(who, loc)

	add_char_damage(who, 30, MATES)

	if p_char(who).health != 70 {
		t.Errorf("health = %d, want 70", p_char(who).health)
	}
}

// TestAddCharDamageKillsWhenZero tests that fatal damage triggers kill_char.
func TestAddCharDamageKillsWhenZero(t *testing.T) {
	who := 5008
	loc := 10003
	setupTestLocation(loc, sub_plain)
	setupTestCharacter(who, 50)
	set_where(who, loc)

	add_char_damage(who, 50, MATES)

	// Characters owned by indep_player (not NPC subkind) create dead bodies
	if kind(who) != T_item || subkind(who) != sub_dead_body {
		t.Errorf("kind = %d, subkind = %d, want T_item (%d) with sub_dead_body (%d)",
			kind(who), subkind(who), T_item, sub_dead_body)
	}
}

// TestAddCharDamageZeroAmountDoesNothing tests that 0 damage does nothing.
func TestAddCharDamageZeroAmountDoesNothing(t *testing.T) {
	who := 5009
	setupTestCharacter(who, 100)

	add_char_damage(who, 0, MATES)

	if p_char(who).health != 100 {
		t.Errorf("health = %d, want 100", p_char(who).health)
	}
}

// TestKillCharWithMeltMe tests that melting characters don't leave bodies.
func TestKillCharWithMeltMe(t *testing.T) {
	who := 5010
	loc := 10004
	setupTestLocation(loc, sub_plain)
	setupTestCharacter(who, 100)
	p_char(who).melt_me = TRUE
	set_where(who, loc)

	kill_char(who, MATES)

	if kind(who) != T_deadchar {
		t.Errorf("kind = %d, want T_deadchar (%d)", kind(who), T_deadchar)
	}

	if subkind(who) == sub_dead_body {
		t.Error("melting characters should not become dead bodies")
	}
}

// TestKillCharCreatesDeadBody tests that normal death creates a dead body item.
func TestKillCharCreatesDeadBody(t *testing.T) {
	who := 5011
	loc := 10005
	pl := 1001
	setupTestLocation(loc, sub_plain)
	setupTestPlayer(pl)
	setupTestCharacter(who, 100)
	p_char(who).unit_lord = pl
	set_where(who, loc)

	kill_char(who, MATES)

	if kind(who) != T_item {
		t.Errorf("kind = %d, want T_item (%d)", kind(who), T_item)
	}

	if subkind(who) != sub_dead_body {
		t.Errorf("subkind = %d, want sub_dead_body (%d)", subkind(who), sub_dead_body)
	}
}

// TestPutBackCookie tests NPC cookie return.
func TestPutBackCookie(t *testing.T) {
	who := 5012
	home := 10006
	cookie := item_mob_cookie // use a valid item constant
	setupTestCharacter(who, 100)
	setupTestLocation(home, sub_lair)

	// Create a valid item box for the cookie
	if teg.globals.bx[cookie] == nil {
		teg.globals.bx[cookie] = &box{}
	}
	teg.globals.bx[cookie].kind = T_item
	teg.globals.bx[cookie].skind = 0
	teg.globals.bx[cookie].x_item = &entity_item{}

	teg.globals.bx[who].x_misc = &entity_misc{
		npc_home:   home,
		npc_cookie: cookie,
	}

	initialQty := has_item(home, cookie)

	put_back_cookie(who)

	finalQty := has_item(home, cookie)
	if finalQty != initialQty+1 {
		t.Errorf("cookie qty = %d, want %d", finalQty, initialQty+1)
	}
}

// TestDeadCharBodyAtSea tests that deaths at sea find nearest land.
func TestDeadCharBodyAtSea(t *testing.T) {
	who := 5013
	ocean := 10007
	setupTestLocation(ocean, sub_ocean)
	setupTestCharacter(who, 100)
	set_where(who, ocean)

	dead_char_body(indep_player, who)

	if kind(who) != T_deadchar {
		t.Errorf("kind = %d, want T_deadchar (body lost at sea becomes deadchar)", kind(who))
	}
}

// TestRestoreDeadBody tests reviving a character from a dead body.
func TestRestoreDeadBody(t *testing.T) {
	who := 5014
	owner := 5015
	loc := 10008
	pl := 1002

	setupTestLocation(loc, sub_plain)
	setupTestPlayer(pl)
	setupTestCharacter(owner, 100)
	set_where(owner, loc)
	p_char(owner).unit_lord = pl

	if teg.globals.bx[who] == nil {
		teg.globals.bx[who] = &box{}
	}
	teg.globals.bx[who].kind = T_item
	teg.globals.bx[who].skind = sub_dead_body
	teg.globals.bx[who].x_item = &entity_item{weight: 10}
	teg.globals.bx[who].x_char = &entity_char{}
	teg.globals.bx[who].x_misc = &entity_misc{old_lord: pl}
	teg.setName(who, "dead body")
	// Use savedNames map for restored name
	savedNames[who] = "Dead Hero"

	add_item(owner, who, 1)
	p_item(who).who_has = owner

	restore_dead_body(owner, who)

	if kind(who) != T_char {
		t.Errorf("kind = %d, want T_char (%d)", kind(who), T_char)
	}

	if p_char(who).health != 100 {
		t.Errorf("health = %d, want 100", p_char(who).health)
	}

	if teg.getName(who) != "Dead Hero" {
		t.Errorf("name = %q, want %q", teg.getName(who), "Dead Hero")
	}
}

// TestIsNpc tests the is_npc helper.
func TestIsNpc(t *testing.T) {
	who := 5016
	setupTestCharacter(who, 100)

	if is_npc(who) {
		t.Error("character with subkind 0 and LOY_unsworn should not be NPC")
	}

	teg.globals.bx[who].skind = sub_ni
	if !is_npc(who) {
		t.Error("character with non-zero subkind should be NPC")
	}

	teg.globals.bx[who].skind = 0
	p_char(who).loy_kind = LOY_npc
	if !is_npc(who) {
		t.Error("character with LOY_npc should be NPC")
	}

	p_char(who).loy_kind = LOY_summon
	if !is_npc(who) {
		t.Error("character with LOY_summon should be NPC")
	}
}

// TestForgetSkill tests skill forgetting.
func TestForgetSkill(t *testing.T) {
	who := 5017
	setupTestCharacter(who, 100)

	if teg.globals.charSkills == nil {
		teg.globals.charSkills = make(map[int][]*skill_ent)
	}
	teg.globals.charSkills[who] = []*skill_ent{
		{skill: sk_combat, know: SKILL_know},
		{skill: sk_survive_fatal, know: SKILL_know},
	}

	result := forget_skill(who, sk_survive_fatal)

	if !result {
		t.Error("forget_skill should return true when skill is known")
	}

	skills := teg.globals.charSkills[who]
	for _, s := range skills {
		if s.skill == sk_survive_fatal && s.know == SKILL_know {
			t.Error("skill should be forgotten")
		}
	}

	result = forget_skill(who, sk_survive_fatal)
	if result {
		t.Error("forget_skill should return false when skill already forgotten")
	}
}

// TestTakeUnitItemsToStackmate tests item transfer to stackmates.
func TestTakeUnitItemsToStackmate(t *testing.T) {
	from := 5018
	to := 5019
	setupTestCharacter(from, 100)
	setupTestCharacter(to, 100)

	// Manually set up the stack relationship without using set_where
	teg.globals.bx[from].x_loc_info.here_list = []int{to}
	teg.globals.bx[to].x_loc_info.where = from

	add_item(from, item_gold, 100)

	take_unit_items(from, MATES, TAKE_ALL)

	toGold := has_item(to, item_gold)
	if toGold == 0 {
		t.Error("stackmate should have received gold")
	}
}

// TestTakeUnitItemsDiscard tests item discard when no inheritor.
func TestTakeUnitItemsDiscard(t *testing.T) {
	from := 5020
	setupTestCharacter(from, 100)

	add_item(from, item_gold, 100)

	take_unit_items(from, 0, TAKE_ALL)

	fromGold := has_item(from, item_gold)
	if fromGold != 0 {
		t.Errorf("from should have 0 gold after discard, got %d", fromGold)
	}
}

// TestLoopStackList tests recursive stack enumeration.
func TestLoopStackList(t *testing.T) {
	leader := 5021
	follower1 := 5022
	follower2 := 5023
	setupTestCharacter(leader, 100)
	setupTestCharacter(follower1, 100)
	setupTestCharacter(follower2, 100)

	// Manually set up the stack relationship without using set_where
	teg.globals.bx[leader].x_loc_info.here_list = []int{follower1}
	teg.globals.bx[follower1].x_loc_info.here_list = []int{follower2}
	teg.globals.bx[follower1].x_loc_info.where = leader
	teg.globals.bx[follower2].x_loc_info.where = follower1

	result := loop_stack_list(leader)

	if len(result) != 3 {
		t.Errorf("loop_stack_list returned %d members, want 3", len(result))
	}

	found := make(map[int]bool)
	for _, id := range result {
		found[id] = true
	}

	if !found[leader] || !found[follower1] || !found[follower2] {
		t.Error("loop_stack_list should return all stack members")
	}
}

// TestChangeBoxKindLifecycle tests kind changing.
func TestChangeBoxKindLifecycle(t *testing.T) {
	who := 5024
	setupTestCharacter(who, 100)

	if kind(who) != T_char {
		t.Errorf("initial kind = %d, want T_char", kind(who))
	}

	change_box_kind(who, T_deadchar)

	if kind(who) != T_deadchar {
		t.Errorf("kind = %d, want T_deadchar", kind(who))
	}
}

// TestChangeBoxSubkindLifecycle tests subkind changing.
func TestChangeBoxSubkindLifecycle(t *testing.T) {
	who := 5025
	if teg.globals.bx[who] == nil {
		teg.globals.bx[who] = &box{}
	}
	teg.globals.bx[who].kind = T_item
	teg.globals.bx[who].skind = 0

	change_box_subkind(who, sub_dead_body)

	if subkind(who) != sub_dead_body {
		t.Errorf("subkind = %d, want sub_dead_body", subkind(who))
	}
}
