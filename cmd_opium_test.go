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

// cmd_opium_test.go - Tests for opium & misc commands
// Sprint 26.10: Opium & Misc Commands

package taygete

import (
	"math/rand/v2"
	"testing"

	"github.com/mdhender/prng"
)

func setupOpiumTest() {
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
	if teg.globals.charSkills == nil {
		teg.globals.charSkills = make(map[int][]*skill_ent)
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
	for id := range teg.globals.charSkills {
		delete(teg.globals.charSkills, id)
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
}

// TestVImproveOpiumInPoppyField tests that v_improve_opium succeeds in a poppy field.
func TestVImproveOpiumInPoppyField(t *testing.T) {
	setupOpiumTest()

	poppyField := 60001
	alloc_box(poppyField, T_loc, sub_poppy_field)

	who := 1001
	alloc_box(who, T_char, 0)
	set_where(who, poppyField)

	c := &command{who: who}
	result := v_improve_opium(c)

	if result != TRUE {
		t.Errorf("v_improve_opium = %d, want TRUE (in poppy field)", result)
	}
}

// TestVImproveOpiumNotInPoppyField tests that v_improve_opium fails outside poppy fields.
func TestVImproveOpiumNotInPoppyField(t *testing.T) {
	setupOpiumTest()

	province := 10000
	alloc_box(province, T_loc, sub_plain)

	who := 1001
	alloc_box(who, T_char, 0)
	set_where(who, province)

	c := &command{who: who}
	result := v_improve_opium(c)

	if result != FALSE {
		t.Errorf("v_improve_opium = %d, want FALSE (not in poppy field)", result)
	}
}

// TestDImproveOpiumSetsFlag tests that d_improve_opium sets the opium_double flag.
func TestDImproveOpiumSetsFlag(t *testing.T) {
	setupOpiumTest()

	poppyField := 60001
	alloc_box(poppyField, T_loc, sub_poppy_field)

	who := 1001
	alloc_box(who, T_char, 0)
	set_where(who, poppyField)

	// Clear any previous flag
	p_misc(poppyField).opium_double = FALSE

	c := &command{who: who}
	result := d_improve_opium(c)

	if result != TRUE {
		t.Errorf("d_improve_opium = %d, want TRUE", result)
	}

	if p_misc(poppyField).opium_double != TRUE {
		t.Errorf("opium_double = %d, want TRUE", p_misc(poppyField).opium_double)
	}
}

// TestDImproveOpiumNotInPoppyField tests that d_improve_opium fails if no longer in poppy field.
func TestDImproveOpiumNotInPoppyField(t *testing.T) {
	setupOpiumTest()

	province := 10000
	alloc_box(province, T_loc, sub_plain)

	who := 1001
	alloc_box(who, T_char, 0)
	set_where(who, province)

	c := &command{who: who}
	result := d_improve_opium(c)

	if result != FALSE {
		t.Errorf("d_improve_opium = %d, want FALSE (not in poppy field anymore)", result)
	}
}

// TestVDieTriggersDeath tests that v_die calls kill_char and triggers death.
func TestVDieTriggersDeath(t *testing.T) {
	setupOpiumTest()

	// Set sysclock to a non-zero value so death_time will be set
	teg.globals.sysclock = olytime{day: 5, turn: 10}

	loc := 10000
	alloc_box(loc, T_loc, sub_plain)

	who := 1001
	alloc_box(who, T_char, 0)
	set_where(who, loc)
	p_char(who).health = 100

	c := &command{who: who}
	result := v_die(c)

	if result != TRUE {
		t.Errorf("v_die = %d, want TRUE", result)
	}

	// After death, the character's death_time should be set to sysclock
	if p_char(who).death_time.day != 5 || p_char(who).death_time.turn != 10 {
		t.Errorf("v_die did not trigger death - death_time = {%d, %d}, want {5, 10}",
			p_char(who).death_time.day, p_char(who).death_time.turn)
	}
}

// TestVDieWithSurviveFatal tests that v_die respects survive_fatal skill.
func TestVDieWithSurviveFatal(t *testing.T) {
	setupOpiumTest()

	loc := 10000
	alloc_box(loc, T_loc, sub_plain)

	who := 1001
	alloc_box(who, T_char, 0)
	set_where(who, loc)
	p_char(who).health = 100

	// Give the character survive_fatal skill
	teg.globals.charSkills[who] = []*skill_ent{
		{skill: sk_survive_fatal, know: SKILL_know},
	}

	c := &command{who: who}
	result := v_die(c)

	if result != TRUE {
		t.Errorf("v_die = %d, want TRUE", result)
	}

	// With survive_fatal, the character should still be alive
	if kind(who) != T_char {
		t.Errorf("Character with survive_fatal should still be T_char, got kind=%d", kind(who))
	}

	// Health should be restored to 100
	if p_char(who).health != 100 {
		t.Errorf("Health after survive_fatal = %d, want 100", p_char(who).health)
	}

	// The skill should be forgotten
	if has_skill_check(who, sk_survive_fatal) {
		t.Errorf("survive_fatal skill should be forgotten after use")
	}
}
