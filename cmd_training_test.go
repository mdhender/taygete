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

// cmd_training_test.go - Tests for combat training commands
// Sprint 26.9: Combat Training

package taygete

import (
	"math/rand/v2"
	"testing"

	"github.com/mdhender/prng"
)

func setupTrainingTest() {
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
}

func TestVArchery(t *testing.T) {
	c := &command{who: 1001}
	result := v_archery(c)
	if result != TRUE {
		t.Errorf("v_archery() = %d, want %d", result, TRUE)
	}
}

func TestDArchery(t *testing.T) {
	setupTrainingTest()

	charID := 1001
	alloc_box(charID, T_char, 0)
	p_char(charID).missile = 50

	c := &command{who: charID}

	initialMissile := p_char(charID).missile
	result := d_archery(c)

	if result != TRUE {
		t.Errorf("d_archery() = %d, want %d", result, TRUE)
	}

	if p_char(charID).missile <= initialMissile {
		t.Errorf("d_archery() did not increase missile rating: was %d, now %d",
			initialMissile, p_char(charID).missile)
	}

	increase := p_char(charID).missile - initialMissile
	if increase < 1 || increase > 10 {
		t.Errorf("d_archery() missile increase %d out of expected range [1,10]", increase)
	}
}

func TestDArcheryHighSkill(t *testing.T) {
	setupTrainingTest()

	charID := 1002
	alloc_box(charID, T_char, 0)
	p_char(charID).missile = 120

	c := &command{who: charID}

	result := d_archery(c)

	if result != TRUE {
		t.Errorf("d_archery() = %d, want %d", result, TRUE)
	}

	if p_char(charID).missile <= 120 {
		t.Errorf("d_archery() should increase missile even at high skill")
	}
}

func TestVDefense(t *testing.T) {
	c := &command{who: 1001}
	result := v_defense(c)
	if result != TRUE {
		t.Errorf("v_defense() = %d, want %d", result, TRUE)
	}
}

func TestDDefense(t *testing.T) {
	setupTrainingTest()

	charID := 1003
	alloc_box(charID, T_char, 0)
	p_char(charID).defense = 40

	c := &command{who: charID}

	initialDefense := p_char(charID).defense
	result := d_defense(c)

	if result != TRUE {
		t.Errorf("d_defense() = %d, want %d", result, TRUE)
	}

	if p_char(charID).defense <= initialDefense {
		t.Errorf("d_defense() did not increase defense rating: was %d, now %d",
			initialDefense, p_char(charID).defense)
	}

	increase := p_char(charID).defense - initialDefense
	if increase < 1 || increase > 10 {
		t.Errorf("d_defense() defense increase %d out of expected range [1,10]", increase)
	}
}

func TestDDefenseHighSkill(t *testing.T) {
	setupTrainingTest()

	charID := 1004
	alloc_box(charID, T_char, 0)
	p_char(charID).defense = 150

	c := &command{who: charID}

	result := d_defense(c)

	if result != TRUE {
		t.Errorf("d_defense() = %d, want %d", result, TRUE)
	}

	if p_char(charID).defense <= 150 {
		t.Errorf("d_defense() should increase defense even at high skill")
	}
}

func TestVSwordplay(t *testing.T) {
	c := &command{who: 1001}
	result := v_swordplay(c)
	if result != TRUE {
		t.Errorf("v_swordplay() = %d, want %d", result, TRUE)
	}
}

func TestDSwordplay(t *testing.T) {
	setupTrainingTest()

	charID := 1005
	alloc_box(charID, T_char, 0)
	p_char(charID).attack = 60

	c := &command{who: charID}

	initialAttack := p_char(charID).attack
	result := d_swordplay(c)

	if result != TRUE {
		t.Errorf("d_swordplay() = %d, want %d", result, TRUE)
	}

	if p_char(charID).attack <= initialAttack {
		t.Errorf("d_swordplay() did not increase attack rating: was %d, now %d",
			initialAttack, p_char(charID).attack)
	}

	increase := p_char(charID).attack - initialAttack
	if increase < 1 || increase > 10 {
		t.Errorf("d_swordplay() attack increase %d out of expected range [1,10]", increase)
	}
}

func TestDSwordplayHighSkill(t *testing.T) {
	setupTrainingTest()

	charID := 1006
	alloc_box(charID, T_char, 0)
	p_char(charID).attack = 200

	c := &command{who: charID}

	result := d_swordplay(c)

	if result != TRUE {
		t.Errorf("d_swordplay() = %d, want %d", result, TRUE)
	}

	if p_char(charID).attack <= 200 {
		t.Errorf("d_swordplay() should increase attack even at high skill")
	}
}

func TestVFightToDeathDisabled(t *testing.T) {
	setupTrainingTest()

	charID := 1007
	alloc_box(charID, T_char, 0)
	p_char(charID).break_point = 50

	c := &command{who: charID, a: 0}

	result := v_fight_to_death(c)

	if result != TRUE {
		t.Errorf("v_fight_to_death() = %d, want %d", result, TRUE)
	}

	if p_char(charID).break_point != 0 {
		t.Errorf("v_fight_to_death(0) break_point = %d, want 0",
			p_char(charID).break_point)
	}
}

func TestVFightToDeathEnabled(t *testing.T) {
	setupTrainingTest()

	charID := 1008
	alloc_box(charID, T_char, 0)
	p_char(charID).break_point = 0

	c := &command{who: charID, a: 1}

	result := v_fight_to_death(c)

	if result != TRUE {
		t.Errorf("v_fight_to_death() = %d, want %d", result, TRUE)
	}

	if p_char(charID).break_point != 50 {
		t.Errorf("v_fight_to_death(1) break_point = %d, want 50",
			p_char(charID).break_point)
	}
}

func TestTrainingStatIncreasesAreAdditive(t *testing.T) {
	setupTrainingTest()

	charID := 1009
	alloc_box(charID, T_char, 0)
	p_char(charID).attack = 10
	p_char(charID).defense = 10
	p_char(charID).missile = 10

	c := &command{who: charID}

	for i := 0; i < 5; i++ {
		d_archery(c)
		d_defense(c)
		d_swordplay(c)
	}

	if p_char(charID).missile <= 10+5 {
		t.Errorf("Multiple archery trainings should accumulate: missile = %d", p_char(charID).missile)
	}
	if p_char(charID).defense <= 10+5 {
		t.Errorf("Multiple defense trainings should accumulate: defense = %d", p_char(charID).defense)
	}
	if p_char(charID).attack <= 10+5 {
		t.Errorf("Multiple swordplay trainings should accumulate: attack = %d", p_char(charID).attack)
	}
}
