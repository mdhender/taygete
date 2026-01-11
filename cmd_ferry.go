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

// cmd_ferry.go - Ships & ferry commands ported from src/c2.c
// Sprint 26.8: Ships & Ferries

package taygete

import "strings"

// v_fee sets the boarding fee for a ferry.
// Usage: FEE <amount>
// Sets the fee per 100 weight units for boarding the ship.
// Ported from src/c2.c lines 822-832.
func v_fee(c *command) int {
	amount := c.a

	p_magic(c.who).fee = amount

	wout(c.who, "Ship boarding fee set to %s per 100 weight.", gold_s(amount))

	return TRUE
}

// board_message announces a character boarding a ship.
// Only announces if the character is not hidden and weather permits.
// Ported from src/c2.c lines 836-861.
func board_message(who, ship int) {
	where := subloc(ship)

	if char_really_hidden(who) {
		return
	}

	if weather_here(where, sub_fog) != 0 {
		return
	}

	with := display_with(who)
	desc := liner_desc(who)

	comma := ""
	if strings.Contains(desc, ",") {
		comma = ","
	}

	if with == "" {
		with = "."
	}

	wout(where, "%s%s boarded %s%s", desc, comma, box_name(ship), with)
	show_chars_below(where, who)
}

// v_board boards a ferry, paying the boarding fee.
// Usage: BOARD <ship> [max-fee]
// The ship must have a FEE set to be operated as a ferry.
// Ported from src/c2.c lines 865-971.
func v_board(c *command) int {
	ship := c.a
	maxFee := c.b

	if !is_ship(ship) {
		wout(c.who, "%s is not a ship.", box_code(ship))
		return FALSE
	}

	log_write(LOG_SPECIAL, "BOARD for %s", box_name(player(c.who)))

	v := parse_exit_dir(c, subloc(c.who), "board")
	if v == nil {
		return FALSE
	}

	if v.destination != ship {
		wout(c.who, "No visible route from %s to %s.", box_name(subloc(c.who)), box_code(ship))
		return FALSE
	}

	if v.in_transit != 0 {
		wout(c.who, "%s is underway. Boarding is not possible.", box_name(v.destination))
		return FALSE
	}

	owner := building_owner(ship)

	shipFee := 0
	if valid_box(owner) {
		shipFee = board_fee(owner)
	}
	if !valid_box(owner) || shipFee == 0 {
		wout(c.who, "%s is not being operated as a ferry (no boarding FEE is set).", box_name(ship))
		return FALSE
	}

	var w weights
	determine_stack_weights(c.who, &w)

	sc := ship_cap(ship)
	if sc != 0 {
		sw := ship_weight(ship)

		if sw > sc {
			wout(c.who, "%s is already overloaded. It can take no more passengers.", box_name(ship))
			wout(owner, "Refused to let %s board because we are overloaded.", box_name(c.who))
			return FALSE
		}

		if sw+w.total_weight > sc {
			wout(c.who, "%s would be overloaded with us. We can't board.", box_name(ship))
			wout(owner, "Refused to let %s board because then we would be overloaded.", box_name(c.who))
			return FALSE
		}
	}

	amount := w.total_weight * shipFee / 100

	if maxFee != 0 && amount > maxFee {
		wout(c.who, "Refused to pay a boarding fee of %s.", gold_s(amount))
		wout(owner, "%s refused to pay a boarding fee of %s.", box_name(c.who), gold_s(amount))
		return FALSE
	}

	if !charge(c.who, amount) {
		wout(c.who, "Can't afford a boarding fee of %s.", gold_s(amount))
		wout(owner, "%s couldn't afford a boarding fee of %s.", box_name(c.who), gold_s(amount))
		return FALSE
	}

	wout(c.who, "Paid %s to board %s.", gold_s(amount), box_name(ship))
	wout(owner, "%s paid %s to board.", box_name(c.who), gold_s(amount))
	board_message(c.who, ship)

	gen_item(owner, item_gold, amount)
	add_gold_ferry(amount)
	move_stack(c.who, ship)

	return TRUE
}

// unboard_message announces a character disembarking from a ship.
// Only announces if the character is not hidden and weather permits.
// Ported from src/c2.c lines 976-1002.
func unboard_message(who, ship int) {
	where := subloc(ship)

	if char_really_hidden(who) {
		return
	}

	if weather_here(where, sub_fog) != 0 {
		return
	}

	with := display_with(who)
	desc := liner_desc(who)

	comma := ""
	if strings.Contains(desc, ",") {
		comma = ","
	}

	if with == "" {
		with = "."
	}

	wout(where, "%s%s disembarked from %s%s", desc, comma, box_name(ship), with)
	show_chars_below(where, who)
}

// v_unload unloads all passengers from a ferry.
// Usage: UNLOAD
// Only the ship's captain can unload passengers.
// Cannot unload at sea.
// Ported from src/c2.c lines 1010-1053.
func v_unload(c *command) int {
	ship := subloc(c.who)

	if !is_ship(ship) || building_owner(ship) != c.who {
		wout(c.who, "%s is not the captain of a ship.", box_name(c.who))
		return FALSE
	}

	where := subloc(ship)

	if subkind(where) == sub_ocean {
		wout(c.who, "Can't unload passengers at sea. They won't go.")
		return FALSE
	}

	any := false

	var chars []int
	loop_char_here(ship, &chars)

	for _, i := range chars {
		if stack_leader(i) == c.who {
			continue
		}

		wout(c.who, "%s disembarks.", box_name(i))
		wout(i, "%s disembarks.", box_name(i))
		unboard_message(i, ship)

		move_stack(i, where)
		any = true
	}

	if any {
		wout(c.who, "All passengers unloaded.")
	} else {
		wout(c.who, "No passengers to unload.")
	}

	return TRUE
}

// v_ferry sounds the ferry horn to wake up waiting units.
// Usage: FERRY
// Sets the ferry_flag on the ship, which wakes up any WAIT FERRY conditions.
// Ported from src/c2.c lines 1061-1081.
func v_ferry(c *command) int {
	ship := subloc(c.who)

	if !is_ship(ship) || building_owner(ship) != c.who {
		wout(c.who, "%s is not the captain of a ship.", box_name(c.who))
		return FALSE
	}

	where := subloc(ship)

	wout(where, "%s sounds a blast on its horn.", box_name(ship))
	log_write(LOG_SPECIAL, "FERRY for %s", box_name(player(c.who)))

	p_magic(ship).ferry_flag = TRUE

	return TRUE
}

// parse_exit_dir parses a direction or destination from a command.
// Returns nil if no valid exit is found.
// Ported from src/move.c lines 165-247.
func parse_exit_dir(c *command, where int, zeroArg string) *exit_view {
	l := exits_from_loc(c.who, where)

	if valid_box(c.a) {
		if where == c.a {
			if zeroArg != "" {
				wout(c.who, "Already in %s.", box_name(where))
			}
			return nil
		}

		var ret *exit_view
		var impassRet *exit_view

		for _, v := range l {
			if v.destination == c.a && (v.hidden == 0 || see_all(c.who)) {
				if v.impassable != 0 {
					impassRet = v
				} else {
					ret = v
				}
			}
		}

		if ret != nil {
			return ret
		}
		if impassRet != nil {
			return impassRet
		}

		if zeroArg != "" {
			wout(c.who, "No visible route from %s to %s.", box_name(where), box_code(c.a))
		}
		return nil
	}

	dir := lookup_dir(get_parse_arg(c, 1))
	if dir < 0 {
		if zeroArg != "" {
			wout(c.who, "Unknown direction or destination '%s'.", get_parse_arg(c, 1))
		}
		return nil
	}

	for _, v := range l {
		if v.direction == dir && (v.hidden == 0 || see_all(c.who)) {
			if dir == DIR_IN && zeroArg != "" {
				wout(c.who, "(assuming '%s %s')", zeroArg, box_code_less(v.destination))
			}
			return v
		}
	}

	if zeroArg != "" {
		wout(c.who, "No visible %s route from %s.", dir_name(dir), box_name(where))
	}
	return nil
}

// display_with returns a string describing what is stacked with a character.
// Returns empty string if nothing is stacked with them.
// Stub for now - will be fully implemented in display sprint.
func display_with(who int) string {
	return ""
}

// show_chars_below outputs characters stacked below a given character to a location.
// Stub for now - will be fully implemented in display sprint.
func show_chars_below(where, who int) {
}

// add_gold_ferry adds to the gold_ferry global tracking ferry income.
// Stub for now - gold_ferry is a global tracking stat.
func add_gold_ferry(amount int) {
}

// see_all returns true if the character can see all hidden things.
// Stub for now - will be fully implemented in visibility sprint.
func see_all(who int) bool {
	return teg.globals.immedSeeAll
}

// lookup_dir looks up a direction string and returns the direction constant.
// Returns -1 if not found.
func lookup_dir(s string) int {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return -1
	}

	for i, d := range full_dir_s {
		if strings.EqualFold(d, s) {
			return i
		}
	}
	for i, d := range short_dir_s {
		if strings.EqualFold(d, s) {
			return i
		}
	}
	return -1
}

// dir_name returns the full name for a direction constant.
func dir_name(dir int) string {
	if dir >= 0 && dir < len(full_dir_s) {
		return full_dir_s[dir]
	}
	return "unknown"
}

// Direction constants and string tables are defined in glob.go:
// DIR_N, DIR_E, DIR_S, DIR_W, DIR_UP, DIR_DOWN, DIR_IN, DIR_OUT
// full_dir_s, short_dir_s
