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

// cmd_training.go - Combat training commands ported from src/c2.c
// Sprint 26.9: Combat Training

package taygete

// v_archery starts archery training.
// Training takes 7 days (see use.c skill table).
// Ported from src/c2.c lines 647-651.
func v_archery(c *command) int {
	return TRUE
}

// d_archery executes archery training, raising the missile rating.
// 5% chance of +10, otherwise +3-5 if missile < 100, else +1-3.
// Ported from src/c2.c lines 655-676.
func d_archery(c *command) int {
	p := p_char(c.who)

	var amount int
	if rnd(1, 100) <= 5 {
		amount = 10
	} else if p.missile < 100 {
		amount = rnd(3, 5)
	} else {
		amount = rnd(1, 3)
	}

	p.missile += short(amount)

	wout(c.who, "Missile rating raised %d to %d.", amount, p.missile)
	return TRUE
}

// v_defense starts defense training.
// Training takes 7 days (see use.c skill table).
// Ported from src/c2.c lines 680-684.
func v_defense(c *command) int {
	return TRUE
}

// d_defense executes defense training, raising the defense rating.
// 5% chance of +10, otherwise +3-5 if defense < 100, else +1-3.
// Ported from src/c2.c lines 688-706.
func d_defense(c *command) int {
	p := p_char(c.who)

	var amount int
	if rnd(1, 100) <= 5 {
		amount = 10
	} else if p.defense < 100 {
		amount = rnd(3, 5)
	} else {
		amount = rnd(1, 3)
	}

	p.defense += short(amount)

	wout(c.who, "Defense rating raised %d to %d.", amount, p.defense)
	return TRUE
}

// v_swordplay starts swordplay training.
// Training takes 7 days (see use.c skill table).
// Ported from src/c2.c lines 710-714.
func v_swordplay(c *command) int {
	return TRUE
}

// d_swordplay executes swordplay training, raising the attack rating.
// 5% chance of +10, otherwise +3-5 if attack < 100, else +1-3.
// Ported from src/c2.c lines 718-736.
func d_swordplay(c *command) int {
	p := p_char(c.who)

	var amount int
	if rnd(1, 100) <= 5 {
		amount = 10
	} else if p.attack < 100 {
		amount = rnd(3, 5)
	} else {
		amount = rnd(1, 3)
	}

	p.attack += short(amount)

	wout(c.who, "Attack rating raised %d to %d.", amount, p.attack)
	return TRUE
}

// v_fight_to_death sets the break point for troops.
// FIGHT 0 - troops fight to the death (break_point = 0)
// FIGHT 1 - troops break at 50% (break_point = 50)
// Ported from src/c2.c lines 801-818.
func v_fight_to_death(c *command) int {
	flag := c.a

	if flag != 0 {
		p_char(c.who).break_point = 50
		wout(c.who, "Troops led by %s will break at 50%%.", box_name(c.who))
	} else {
		p_char(c.who).break_point = 0
		wout(c.who, "Troops led by %s will fight to the death.", box_name(c.who))
	}

	return TRUE
}
