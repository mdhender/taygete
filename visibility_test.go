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

// setupVisibilityTest initializes test state for visibility tests.
func setupVisibilityTest() {
	teg.globals.bx = [MAX_BOXES]*box{}
	teg.globals.sysclock = olytime{days_since_epoch: 100}
	teg.globals.evening = false
	teg.globals.garrison_magic = 999
}

// setupVisibilityTestCharacter creates a basic test character.
func setupVisibilityTestCharacter(who int) {
	if teg.globals.bx[who] == nil {
		teg.globals.bx[who] = &box{}
	}
	teg.globals.bx[who].kind = T_char
	teg.globals.bx[who].skind = 0
	teg.globals.bx[who].x_char = &entity_char{
		health:       100,
		melt_me:      FALSE,
		prisoner:     FALSE,
		unit_lord:    indep_player,
		x_char_magic: &char_magic{},
	}
	teg.globals.bx[who].x_loc_info = loc_info{where: 0}
	teg.setName(who, "Test Character")
}

// setupVisibilityTestLocation creates a basic test location.
func setupVisibilityTestLocation(loc int, sk schar) {
	if teg.globals.bx[loc] == nil {
		teg.globals.bx[loc] = &box{}
	}
	teg.globals.bx[loc].kind = T_loc
	teg.globals.bx[loc].skind = sk
	teg.globals.bx[loc].x_loc_info = loc_info{where: 0}
	teg.setName(loc, "Test Location")
}

// setupVisibilityTestPlayer creates a basic test player.
func setupVisibilityTestPlayer(pl int) {
	if teg.globals.bx[pl] == nil {
		teg.globals.bx[pl] = &box{}
	}
	teg.globals.bx[pl].kind = T_player
	teg.globals.bx[pl].skind = sub_pl_regular
	teg.globals.bx[pl].x_player = &entity_player{}
	teg.setName(pl, "Test Player")
}

// setupVisibilityTestStorm creates a test storm.
func setupVisibilityTestStorm(storm int, sk schar, strength short) {
	if teg.globals.bx[storm] == nil {
		teg.globals.bx[storm] = &box{}
	}
	teg.globals.bx[storm].kind = T_storm
	teg.globals.bx[storm].skind = sk
	teg.globals.bx[storm].x_loc_info = loc_info{where: 0}
	teg.globals.bx[storm].x_misc = &entity_misc{storm_str: strength}
	teg.setName(storm, "Test Storm")
}

// TestContacted tests the contacted function.
func TestContacted(t *testing.T) {
	setupVisibilityTest()

	charA := 2001
	charB := 2002
	playerB := 1001

	setupVisibilityTestCharacter(charA)
	setupVisibilityTestCharacter(charB)
	setupVisibilityTestPlayer(playerB)

	p_char(charB).unit_lord = playerB

	if contacted(charA, charB) {
		t.Error("contacted(charA, charB) = true, want false (no contact)")
	}

	p_char(charA).contact = append(p_char(charA).contact, charB)
	if !contacted(charA, charB) {
		t.Error("contacted(charA, charB) = false, want true (direct contact)")
	}

	p_char(charA).contact = nil
	p_char(charA).contact = append(p_char(charA).contact, playerB)
	if !contacted(charA, charB) {
		t.Error("contacted(charA, charB) = false, want true (player contact)")
	}

	if contacted(charB, charA) {
		t.Error("contacted(charB, charA) = true, want false (not mutual)")
	}
}

// TestCharAlone tests the char_alone function.
func TestCharAlone(t *testing.T) {
	setupVisibilityTest()

	province := 1000
	charA := 2001
	charB := 2002

	setupVisibilityTestLocation(province, sub_plain)
	setupVisibilityTestCharacter(charA)
	setupVisibilityTestCharacter(charB)

	set_where(charA, province)

	if !char_alone(charA) {
		t.Error("char_alone(charA) = false, want true (alone in province)")
	}

	set_where(charB, charA)
	if char_alone(charA) {
		t.Error("char_alone(charA) = true, want false (has stack member)")
	}

	if char_alone(charB) {
		t.Error("char_alone(charB) = true, want false (stacked under charA)")
	}
}

// TestCharReallyHidden tests the char_really_hidden function.
func TestCharReallyHidden(t *testing.T) {
	setupVisibilityTest()

	province := 1000
	charA := 2001
	charB := 2002

	setupVisibilityTestLocation(province, sub_plain)
	setupVisibilityTestCharacter(charA)
	setupVisibilityTestCharacter(charB)

	set_where(charA, province)

	if char_really_hidden(charA) {
		t.Error("char_really_hidden(charA) = true, want false (not hidden)")
	}

	p_magic(charA).hide_self = 1
	if !char_really_hidden(charA) {
		t.Error("char_really_hidden(charA) = false, want true (hidden and alone)")
	}

	set_where(charB, charA)
	if char_really_hidden(charA) {
		t.Error("char_really_hidden(charA) = true, want false (hidden but not alone)")
	}
}

// TestGarrisonHere tests the garrison_here function.
func TestGarrisonHere(t *testing.T) {
	setupVisibilityTest()

	province := 1000
	garrison := 2001
	charA := 2002

	setupVisibilityTestLocation(province, sub_plain)

	if got := garrison_here(province); got != 0 {
		t.Errorf("garrison_here(empty province) = %d, want 0", got)
	}

	setupVisibilityTestCharacter(garrison)
	teg.globals.bx[garrison].skind = sub_garrison
	set_where(garrison, province)

	if got := garrison_here(province); got != garrison {
		t.Errorf("garrison_here(province with garrison) = %d, want %d", got, garrison)
	}

	setupVisibilityTestCharacter(charA)
	set_where(charA, province)
	if got := garrison_here(province); got != garrison {
		t.Errorf("garrison_here(province) = %d, want %d (garrison is first)", got, garrison)
	}
}

// TestWeatherHere tests the weather_here function.
func TestWeatherHere(t *testing.T) {
	setupVisibilityTest()

	region := 500
	province := 1000
	building := 1500
	storm := 3001

	setupVisibilityTestLocation(region, sub_region)
	setupVisibilityTestLocation(province, sub_plain)
	setupVisibilityTestLocation(building, sub_castle)
	setupVisibilityTestStorm(storm, sub_fog, 5)

	set_where(province, region)
	set_where(building, province)

	if got := weather_here(province, sub_fog); got != 0 {
		t.Errorf("weather_here(province, sub_fog) = %d, want 0 (no storm)", got)
	}

	set_where(storm, province)
	if got := weather_here(province, sub_fog); got != 5 {
		t.Errorf("weather_here(province, sub_fog) = %d, want 5", got)
	}

	if got := weather_here(province, sub_rain); got != 0 {
		t.Errorf("weather_here(province, sub_rain) = %d, want 0 (wrong type)", got)
	}

	if got := weather_here(building, sub_fog); got != 0 {
		t.Errorf("weather_here(building, sub_fog) = %d, want 0 (inside building)", got)
	}
}

// TestCharWhere tests the char_where function.
func TestCharWhere(t *testing.T) {
	setupVisibilityTest()

	region := 500
	province := 1000
	charA := 2001
	charB := 2002
	playerA := 1001
	playerB := 1002

	setupVisibilityTestLocation(region, sub_region)
	setupVisibilityTestLocation(province, sub_plain)
	setupVisibilityTestCharacter(charA)
	setupVisibilityTestCharacter(charB)
	setupVisibilityTestPlayer(playerA)
	setupVisibilityTestPlayer(playerB)

	set_where(province, region)
	set_where(charA, province)
	set_where(charB, province)

	p_char(charA).unit_lord = playerA
	p_char(charB).unit_lord = playerA

	if !char_where(province, charA, charB) {
		t.Error("char_where = false, want true (both visible)")
	}

	p_magic(charB).hide_self = 1
	if !char_where(province, charA, charB) {
		t.Error("char_where = false, want true (hidden but not alone - charB stacked with charA in province)")
	}

	p_char(charB).unit_lord = playerB
	if char_where(province, charA, charB) {
		t.Error("char_where = true, want false (different player, hidden alone)")
	}

	p_char(charB).contact = append(p_char(charB).contact, charA)
	if !char_where(province, charA, charB) {
		t.Error("char_where = false, want true (contacted)")
	}
}

// TestCharHere tests the char_here function.
func TestCharHere(t *testing.T) {
	setupVisibilityTest()

	region := 500
	province := 1000
	province2 := 1001
	charA := 2001
	charB := 2002

	setupVisibilityTestLocation(region, sub_region)
	setupVisibilityTestLocation(province, sub_plain)
	setupVisibilityTestLocation(province2, sub_plain)
	setupVisibilityTestCharacter(charA)
	setupVisibilityTestCharacter(charB)

	set_where(province, region)
	set_where(province2, region)
	set_where(charA, province)
	set_where(charB, province)

	if !char_here(charA, charB) {
		t.Error("char_here = false, want true (same location)")
	}

	set_where(charB, province2)

	if char_here(charA, charB) {
		t.Error("char_here = true, want false (different location)")
	}
}

// TestCheckCharWhere tests the check_char_where function.
func TestCheckCharWhere(t *testing.T) {
	setupVisibilityTest()

	region := 500
	province := 1000
	charA := 2001
	charB := 2002
	item := 4001

	setupVisibilityTestLocation(region, sub_region)
	setupVisibilityTestLocation(province, sub_plain)
	setupVisibilityTestCharacter(charA)
	setupVisibilityTestCharacter(charB)

	teg.globals.bx[item] = &box{kind: T_item}

	set_where(province, region)
	set_where(charA, province)
	set_where(charB, province)

	if !check_char_where(province, charA, charB) {
		t.Error("check_char_where = false, want true")
	}

	if check_char_where(province, charA, teg.globals.garrison_magic) {
		t.Error("check_char_where(garrison_magic) = true, want false")
	}

	if check_char_where(province, charA, item) {
		t.Error("check_char_where(item) = true, want false (not a character)")
	}
}

// TestCheckCharHere tests the check_char_here function.
func TestCheckCharHere(t *testing.T) {
	setupVisibilityTest()

	region := 500
	province := 1000
	charA := 2001
	charB := 2002

	setupVisibilityTestLocation(region, sub_region)
	setupVisibilityTestLocation(province, sub_plain)
	setupVisibilityTestCharacter(charA)
	setupVisibilityTestCharacter(charB)

	set_where(province, region)
	set_where(charA, province)
	set_where(charB, province)

	if !check_char_here(charA, charB) {
		t.Error("check_char_here = false, want true")
	}
}

// TestCheckCharGone tests the check_char_gone function.
func TestCheckCharGone(t *testing.T) {
	setupVisibilityTest()

	region := 500
	province := 1000
	charA := 2001
	charB := 2002

	setupVisibilityTestLocation(region, sub_region)
	setupVisibilityTestLocation(province, sub_plain)
	setupVisibilityTestCharacter(charA)
	setupVisibilityTestCharacter(charB)

	set_where(province, region)
	set_where(charA, province)
	set_where(charB, province)

	if !check_char_gone(charA, charB) {
		t.Error("check_char_gone = false, want true (not moving)")
	}

	p_char(charB).moving = 1
	if check_char_gone(charA, charB) {
		t.Error("check_char_gone = true, want false (target is moving)")
	}

	if check_char_gone(charA, teg.globals.garrison_magic) {
		t.Error("check_char_gone(garrison_magic) = true, want false")
	}
}

// TestCheckStillHere tests the check_still_here function.
func TestCheckStillHere(t *testing.T) {
	setupVisibilityTest()

	region := 500
	province := 1000
	charA := 2001
	charB := 2002

	setupVisibilityTestLocation(region, sub_region)
	setupVisibilityTestLocation(province, sub_plain)
	setupVisibilityTestCharacter(charA)
	setupVisibilityTestCharacter(charB)

	set_where(province, region)
	set_where(charA, province)
	set_where(charB, province)

	if !check_still_here(charA, charB) {
		t.Error("check_still_here = false, want true")
	}

	if check_still_here(charA, teg.globals.garrison_magic) {
		t.Error("check_still_here(garrison_magic) = true, want false")
	}
}

// TestCheckSkill tests the check_skill function.
func TestCheckSkill(t *testing.T) {
	setupVisibilityTest()

	charA := 2001
	skillID := 5001

	setupVisibilityTestCharacter(charA)

	teg.globals.bx[skillID] = &box{kind: T_skill}
	teg.setName(skillID, "Test Skill")

	if teg.globals.charSkills == nil {
		teg.globals.charSkills = make(map[int][]*skill_ent)
	}
	teg.globals.charSkills[charA] = nil

	if check_skill(charA, skillID) {
		t.Error("check_skill = true, want false (no skill)")
	}

	teg.globals.charSkills[charA] = []*skill_ent{
		{skill: skillID, know: SKILL_know},
	}
	if !check_skill(charA, skillID) {
		t.Error("check_skill = false, want true (has skill)")
	}
}

// TestCountStackAny tests the count_stack_any function.
func TestCountStackAny(t *testing.T) {
	setupVisibilityTest()

	province := 1000
	charA := 2001
	charB := 2002
	charC := 2003

	setupVisibilityTestLocation(province, sub_plain)
	setupVisibilityTestCharacter(charA)
	setupVisibilityTestCharacter(charB)
	setupVisibilityTestCharacter(charC)

	set_where(charA, province)

	if got := count_stack_any(charA); got != 1 {
		t.Errorf("count_stack_any(alone) = %d, want 1", got)
	}

	set_where(charB, charA)
	if got := count_stack_any(charA); got != 2 {
		t.Errorf("count_stack_any(with one) = %d, want 2", got)
	}

	set_where(charC, charB)
	if got := count_stack_any(charA); got != 3 {
		t.Errorf("count_stack_any(nested) = %d, want 3", got)
	}
}

// TestFirstCharHere tests the first_char_here function.
func TestFirstCharHere(t *testing.T) {
	setupVisibilityTest()

	province := 1000
	charA := 2001

	setupVisibilityTestLocation(province, sub_plain)

	if got := first_char_here(province); got != 0 {
		t.Errorf("first_char_here(empty) = %d, want 0", got)
	}

	setupVisibilityTestCharacter(charA)
	set_where(charA, province)

	if got := first_char_here(province); got != charA {
		t.Errorf("first_char_here(with char) = %d, want %d", got, charA)
	}
}

// TestHasSkill tests the has_skill function.
func TestHasSkill(t *testing.T) {
	setupVisibilityTest()

	charA := 2001
	skillID := 5001

	setupVisibilityTestCharacter(charA)

	if teg.globals.charSkills == nil {
		teg.globals.charSkills = make(map[int][]*skill_ent)
	}
	teg.globals.charSkills[charA] = nil

	if has_skill(charA, skillID) {
		t.Error("has_skill = true, want false (no skill)")
	}

	teg.globals.charSkills[charA] = []*skill_ent{
		{skill: skillID, know: SKILL_know},
	}
	if !has_skill(charA, skillID) {
		t.Error("has_skill = false, want true (has skill)")
	}

	teg.globals.charSkills[charA] = []*skill_ent{
		{skill: skillID, know: SKILL_dont},
	}
	if has_skill(charA, skillID) {
		t.Error("has_skill = true, want false (skill forgotten)")
	}
}

// TestBarkDogs tests the bark_dogs function.
func TestBarkDogs(t *testing.T) {
	setupVisibilityTest()
	teg.globals.inventories = make(map[int][]item_ent)

	province := 1000
	charA := 2001

	setupVisibilityTestLocation(province, sub_plain)
	setupVisibilityTestCharacter(charA)
	set_where(charA, province)

	// No hounds - should produce no output
	bark_dogs(province)

	// Give character hounds via the inventory system
	teg.globals.inventories[charA] = []item_ent{{item: item_hound, qty: 2}}

	// With hounds - function runs without panic (output is via VECT)
	bark_dogs(province)
}

// TestPrintDot tests that print_dot is a no-op.
func TestPrintDot(t *testing.T) {
	// print_dot is deprecated, just verify it doesn't panic
	print_dot('.')
	print_dot('*')
}

// TestStage tests that stage is a no-op.
func TestStage(t *testing.T) {
	// stage is deprecated, just verify it doesn't panic
	stage("test stage")
	stage("")
	stage("another stage")
}
