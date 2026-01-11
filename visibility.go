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

// visibility.go - Character visibility and skill check functions ported from src/u.c
// Sprint 26.3: Character Visibility & Skill Checks

package taygete

// contacted checks if character a has been contacted by character b.
// Returns true if b (or b's player) is in a's contact list.
// Ported from src/u.c lines 700-714.
func contacted(a, b int) bool {
	p := rp_char(a)
	if p == nil {
		return false
	}

	if IListLookup(p.contact, b) >= 0 {
		return true
	}

	if IListLookup(p.contact, player(b)) >= 0 {
		return true
	}

	return false
}

// char_alone checks if a character is alone (not stacked under anyone and has no stack members).
// Ported from src/oly.h macro: #define char_alone(n) (stack_parent(n) == 0 && count_stack_any(n) == 1)
func char_alone(n int) bool {
	return stack_parent(n) == 0 && count_stack_any(n) == 1
}

// char_really_hidden checks if a character is truly hidden (hidden and alone).
// Ported from src/oly.h macro: #define char_really_hidden(n) (char_hidden(n) && char_alone(n))
func char_really_hidden(n int) bool {
	return char_hidden(n) != 0 && char_alone(n)
}

// garrison_here returns the garrison character at a location, or 0 if none.
// Garrisons are always first in the here_list if present.
// Ported from src/garr.c lines 41-52.
func garrison_here(where int) int {
	n := first_character(where)
	if n != 0 && subkind(n) == sub_garrison {
		return n
	}
	return 0
}

// weather_here returns the total weather strength of a given type at a location.
// Returns 0 for buildings (weather doesn't affect interiors).
// Ported from src/storm.c lines 330-349.
func weather_here(where int, sk schar) int {
	if loc_depth(where) == LOC_build {
		return 0
	}

	where = province(where)

	p := rp_loc_info(where)
	if p == nil {
		return 0
	}

	sum := 0
	for _, i := range p.here_list {
		if kind(i) == T_storm && subkind(i) == sk {
			sum += int(storm_strength(i))
		}
	}

	return sum
}

// char_where checks if a target character can be seen at a specific location.
// Hidden characters are only visible to:
//   - Same player
//   - Garrison at the location
//   - Characters who have contacted them
//   - Characters in the same stack
//
// Ported from src/u.c lines 717-746.
func char_where(where, who, target int) bool {
	if where != subloc(target) {
		return false
	}

	if char_really_hidden(target) ||
		(loc_depth(where) == LOC_province && weather_here(where, sub_fog) != 0) {
		pl := player(who)
		if pl == player(target) {
			return true
		}

		if target == garrison_here(where) {
			return true
		}

		if contacted(target, who) {
			return true
		}

		if stack_leader(target) == stack_leader(who) {
			return true
		}

		return false
	}

	return true
}

// char_here checks if a target character can be seen from the caller's location.
// Ported from src/u.c lines 749-755.
func char_here(who, target int) bool {
	where := subloc(who)
	return char_where(where, who, target)
}

// check_char_where validates a target character at a location with error messages.
// Returns false (with wout message) if:
//   - Target is the special garrison_magic value
//   - Target is not a character
//   - Target cannot be seen at the location
//
// Ported from src/u.c lines 758-781.
func check_char_where(where, who, target int) bool {
	if target == teg.globals.garrison_magic {
		wout(who, "There is no garrison here.")
		return false
	}

	if kind(target) != T_char {
		wout(who, "%s is not a character.", box_code(target))
		return false
	}

	if !char_where(where, who, target) {
		wout(who, "%s cannot be seen here.", box_code(target))
		return false
	}

	return true
}

// check_char_here validates a target character at the caller's location.
// Ported from src/u.c lines 784-790.
func check_char_here(who, target int) bool {
	where := subloc(who)
	return check_char_where(where, who, target)
}

// check_char_gone validates a target character is visible and hasn't left.
// Returns false (with wout message) if:
//   - Target is the special garrison_magic value
//   - Target is not a character
//   - Target cannot be seen
//   - Target has left (is moving)
//
// Ported from src/u.c lines 793-822.
func check_char_gone(who, target int) bool {
	if target == teg.globals.garrison_magic {
		wout(who, "There is no garrison here.")
		return false
	}

	if kind(target) != T_char {
		wout(who, "%s is not a character.", box_code(target))
		return false
	}

	if !char_here(who, target) {
		wout(who, "%s cannot be seen here.", box_code(target))
		return false
	}

	if char_gone(target) != 0 {
		wout(who, "%s has left.", box_name(target))
		return false
	}

	return true
}

// check_still_here validates a target character is still visible.
// Similar to check_char_here but with different error message.
// The C code has a commented-out char_gone check.
// Ported from src/u.c lines 825-856.
func check_still_here(who, target int) bool {
	if target == teg.globals.garrison_magic {
		wout(who, "There is no garrison here.")
		return false
	}

	if kind(target) != T_char {
		wout(who, "%s is not a character.", box_code(target))
		return false
	}

	if !char_here(who, target) {
		wout(who, "%s can no longer be seen here.", box_name(target))
		return false
	}

	return true
}

// check_skill checks if a character has a required skill.
// Returns false (with wout message) if the skill is not known.
// Ported from src/u.c lines 859-870.
func check_skill(who, skill int) bool {
	if !has_skill(who, skill) {
		wout(who, "Requires %s.", box_name(skill))
		return false
	}
	return true
}

// has_skill checks if a character knows a skill.
// Returns true if the character has the skill at know level.
func has_skill(who, skill int) bool {
	return has_skill_check(who, skill)
}

// count_stack_any returns the count of all units in a stack (including the leader).
// This counts all characters stacked under who, plus who itself.
func count_stack_any(who int) int {
	count := 1
	p := rp_loc_info(who)
	if p == nil {
		return count
	}
	for _, id := range p.here_list {
		if kind(id) == T_char {
			count += count_stack_any(id)
		}
	}
	return count
}

// first_char_here returns the first character in the here_list of where.
// Ported from src/u.c lines 2249-2265.
func first_char_here(where int) int {
	return first_character(where)
}

// bark_dogs causes hounds at a location to bark, alerting observers.
// Each hound has a 50% chance to bark. The message varies based on
// whether one hound barks, or multiple.
// Ported from src/u.c lines 2440-2469.
func bark_dogs(where int) {
	vector_char_here(where)
	vector_add(where)

	sum := 0
	bark := 0

	p := rp_loc_info(where)
	if p == nil {
		return
	}

	for _, who := range p.here_list {
		if kind(who) != T_char {
			continue
		}

		n := stack_has_item(who, item_hound)
		sum += n
		for i := 1; i <= n; i++ {
			if rnd(1, 2) == 1 {
				bark++
			}
		}
	}

	if bark == 1 && sum == 1 {
		wout(VECT, "The hound is barking.")
	} else if bark == 1 {
		wout(VECT, "A hound is barking.")
	} else if bark > 1 {
		wout(VECT, "The hounds are barking.")
	}
}

// print_dot is a debug/progress helper that prints a character to stderr.
// Deprecated: This is a console progress indicator not needed in Go port.
// Ported from src/u.c lines 2232-2245.
func print_dot(_ int) {
	// No-op: progress indicators handled differently in Go
}

// stage is a debug/timing helper that prints stage names and timing info.
// Deprecated: Debug timing is handled via Go's standard profiling/logging.
// Ported from src/u.c lines 2370-2400.
func stage(_ string) {
	// No-op: timing/stage logging handled via slog in Go
}
