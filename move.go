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

// move.go - Core movement system ported from src/move.c
// Sprint 27: Movement & World

package taygete

import "strings"

// Global state for ocean character tracking
var ocean_chars []int // characters flying over ocean

// exit_opposite is defined in glob.go

// show_to_garrison controls whether garrison units should see messages.
// This is toggled when movements pass by garrisons.
var show_to_garrison = false

// departure_message outputs a message when a character departs a location.
// Only shows if the character is not hidden, weather permits, and destination is not hidden.
// Ported from src/move.c lines 11-73.
func departure_message(who int, v *exit_view) {
	if !valid_box(who) {
		return
	}

	if char_really_hidden(who) {
		return
	}

	if loc_depth(v.orig) == LOC_province && weather_here(v.orig, sub_fog) != 0 {
		return
	}

	if v.dest_hidden != 0 {
		return
	}

	desc := liner_desc(who)

	to := ""
	if subloc(v.destination) == v.orig {
		to = sout(" entered %s", box_name(v.destination))
	} else if subloc(v.orig) == v.destination {
		to = sout(" exited %s", box_name(v.orig))
	} else if viewloc(v.orig) != viewloc(v.destination) {
		if v.direction >= DIR_N && v.direction <= DIR_W {
			to = sout(" went %s", full_dir_s[v.direction])
		} else {
			to = sout(" left for %s", box_name(v.destination))
		}
	} else {
		return
	}

	with := display_with(who)
	comma := ""
	if strings.Contains(desc, ",") {
		comma = ","
	}
	if with == "" {
		with = "."
	}

	if viewloc(v.orig) != viewloc(v.destination) {
		garr := garrison_here(v.orig)
		if garr != 0 && garrison_notices(garr, who) {
			show_to_garrison = true
		}
	}

	wout(v.orig, "%s%s%s%s", desc, comma, to, with)
	show_chars_below(v.orig, who)

	show_to_garrison = false
}

// arrival_message outputs a message when a character arrives at a location.
// Only shows if the character is not hidden and weather permits.
// Ported from src/move.c lines 77-133.
func arrival_message(who int, v *exit_view) {
	if char_really_hidden(who) {
		return
	}

	if loc_depth(v.destination) == LOC_province && weather_here(v.destination, sub_fog) != 0 {
		return
	}

	desc := liner_desc(who)

	from := ""
	if v.orig_hidden == 0 {
		if v.direction >= DIR_N && v.direction <= DIR_W {
			from = sout(" from the %s", full_dir_s[exit_opposite[v.direction]])
		} else {
			from = sout(" from %s", box_name(v.orig))
		}
	}

	with := display_with(who)
	comma := ""
	if strings.Contains(desc, ",") {
		comma = ","
	}
	if with == "" {
		with = "."
	}

	if viewloc(v.orig) != viewloc(v.destination) {
		garr := garrison_here(v.destination)
		if garr != 0 {
			if garrison_notices(garr, who) {
				show_to_garrison = true
			}
			if garrison_spot_check(garr, who) {
				indent += 3
				wout(garr, "%s%s", desc, with)
				show_chars_below(garr, who)
				indent -= 3
			}
		}
	}

	wout(v.destination, "%s%s arrived%s%s", desc, comma, from, with)
	show_chars_below(v.destination, who)

	show_to_garrison = false
}

// discover_road marks both ends of a hidden road as known to the character and stack.
// Ported from src/move.c lines 140-162.
func discover_road(who, where int, v *exit_view) {
	l := exits_from_loc(who, v.destination)

	for _, exit := range l {
		if exit.road != 0 && exit.destination == where {
			set_known(who, exit.road)
			set_known(who, v.road)

			var chars []int
			loop_char_here(who, &chars)
			for _, j := range chars {
				set_known(j, exit.road)
				set_known(j, v.road)
			}
		}
	}
}

// save_v_array saves exit_view data to command fields for later restoration.
// Ported from src/move.c lines 427-438.
func save_v_array(c *command, v *exit_view) {
	c.b = v.direction
	c.c = v.destination
	c.d = v.road
	c.e = v.dest_hidden
	c.f = v.distance
	c.g = v.orig
	c.h = v.orig_hidden
}

// restore_v_array restores exit_view data from command fields.
// Ported from src/move.c lines 441-454.
func restore_v_array(c *command, v *exit_view) {
	*v = exit_view{}
	v.direction = c.b
	v.destination = c.c
	v.road = c.d
	v.dest_hidden = c.e
	v.distance = c.f
	v.orig = c.g
	v.orig_hidden = c.h
}

// suspend_stack_actions marks all characters in a stack as moving.
// This prevents them from executing orders while in transit.
// Ported from src/move.c lines 457-467.
func suspend_stack_actions(who int) {
	var stackMembers []int
	loop_stack(who, &stackMembers)
	for _, i := range stackMembers {
		p_char(i).moving = teg.globals.sysclock.days_since_epoch
	}
}

// restore_stack_actions_impl clears the moving flag for all characters in a stack.
// This is called when movement completes.
// Ported from src/move.c lines 470-480.
func restore_stack_actions_impl(who int) {
	var stackMembers []int
	loop_stack(who, &stackMembers)
	for _, i := range stackMembers {
		p_char(i).moving = 0
	}
}

// clear_guard_flag clears the guard flag for a character and all stacked beneath them.
// Ported from src/move.c lines 483-497.
func clear_guard_flag(who int) {
	if kind(who) == T_char {
		p_char(who).guard = FALSE
	}

	var chars []int
	loop_char_here(who, &chars)
	for _, i := range chars {
		p_char(i).guard = FALSE
	}
}

// move_exit_land calculates the delay for land movement.
// Returns -1 if movement is not possible due to weight constraints.
// Ported from src/move.c lines 250-355.
func move_exit_land(c *command, v *exit_view, show bool) int {
	delay := v.distance
	if delay == 0 {
		return 0
	}

	terr := subkind(v.destination)
	swamp := (terr == sub_swamp || terr == sub_bog || terr == sub_pits)

	var w weights
	determine_stack_weights(c.who, &w)

	if delay > 1 && w.ride_cap >= w.ride_weight && !swamp {
		delay -= delay / 2
	} else {
		if w.land_weight > w.land_cap*2 {
			if show {
				wout(c.who, "%s is too overloaded to travel.", box_name(c.who))
			}
			return -1
		}

		if swamp && w.animals > 0 {
			if show {
				wout(c.who, "Difficult terrain slows the animals. Travel will take an extra day.")
			}
			delay += 1
		}

		if w.land_weight > w.land_cap {
			ratio := (w.land_weight - w.land_cap) * 100 / w.land_cap
			additional := delay * ratio / 100

			if show {
				if additional == 1 {
					wout(c.who, "Excess inventory slows movement. Travel will take an extra day.")
				} else if additional > 1 {
					wout(c.who, "Excess inventory slows movement. Travel will take an extra %s days.", nice_num(additional))
				}
			}
			delay += additional
		}
	}

	nobles := count_stack_move_nobles(c.who)
	men := count_stack_figures(c.who) - nobles
	if nobles < 1 {
		nobles = 1
	}

	extra := men / (nobles * army_slow_factor)
	if extra > v.distance*2 {
		extra = v.distance * 2
	}

	if extra == 1 {
		wout(c.who, "%s noble%s, %s %s: travel will take an extra day.",
			cap_str(nice_num(nobles)), add_s(nobles),
			nice_num(men),
			plural_man(men))
	} else if extra > 1 {
		wout(c.who, "%s noble%s, %s %s: travel will take an extra %d days.",
			cap_str(nice_num(nobles)), add_s(nobles),
			nice_num(men),
			plural_man(men),
			extra)
	}

	delay += extra
	return delay
}

// move_exit_fly calculates the delay for flying movement.
// Returns -1 if movement is not possible.
// Ported from src/move.c lines 358-424.
func move_exit_fly(c *command, v *exit_view, show bool) int {
	delay := v.distance

	if subkind(v.destination) == sub_under || subkind(v.destination) == sub_tunnel {
		if show {
			wout(c.who, "Cannot fly underground.")
		}
		return -1
	}

	if delay < 8 {
		if delay > 3 {
			delay = 3
		}
	} else {
		delay = 4
	}

	var w weights
	determine_stack_weights(c.who, &w)

	if w.fly_cap < w.fly_weight {
		if show {
			wout(c.who, "%s is too overloaded to fly.", box_name(c.who))
		}
		return -1
	}

	nobles := count_stack_move_nobles(c.who)
	men := count_stack_figures(c.who) - nobles
	if nobles < 1 {
		nobles = 1
	}

	extra := men / (nobles * army_slow_factor)
	if extra > v.distance*2 {
		extra = v.distance * 2
	}

	if extra == 1 {
		wout(c.who, "%s noble%s, %s %s: travel will take an extra day.",
			cap_str(nice_num(nobles)), add_s(nobles),
			nice_num(men),
			plural_man(men))
	} else if extra > 1 {
		wout(c.who, "%s noble%s, %s %s: travel will take an extra %d days.",
			cap_str(nice_num(nobles)), add_s(nobles),
			nice_num(men),
			plural_man(men),
			extra)
	}

	delay += extra
	return delay
}

// land_check validates if land movement to a destination is possible.
// Ported from src/move.c lines 500-545.
func land_check(c *command, v *exit_view, show bool) bool {
	if v.water != 0 {
		if show {
			wout(c.who, "A sea-worthy ship is required for travel across water.")
		}
		return false
	}

	if v.impassable != 0 {
		if show {
			wout(c.who, "That route is impassable.")
		}
		return false
	}

	if v.in_transit != 0 {
		if show {
			wout(c.who, "%s is underway. Boarding is not possible.", box_name(v.destination))
		}
		return false
	}

	if loc_depth(v.destination) == LOC_build &&
		subkind(v.destination) != sub_sewer {
		owner := building_owner(v.destination)
		if owner != 0 && !will_admit(owner, c.who, v.destination) && v.direction != DIR_OUT {
			if show {
				wout(c.who, "%s refused to let us enter.", box_name(owner))
				wout(owner, "Refused to let %s enter.", box_name(c.who))
			}
			return false
		}
	}

	return true
}

// can_move_here checks if a character can move in a direction from where.
// Ported from src/move.c lines 548-564.
func can_move_here(where int, c *command) bool {
	v := parse_exit_dir(c, where, "")
	if v != nil && v.direction != DIR_IN && land_check(c, v, false) && move_exit_land(c, v, false) >= 0 {
		return true
	}
	return false
}

// can_move_at_outer_level checks if movement is possible from an outer location.
// Returns the depth difference if movement is possible, 0 otherwise.
// Ported from src/move.c lines 567-581.
func can_move_at_outer_level(where int, c *command) int {
	outer := subloc(where)
	for loc_depth(outer) > LOC_region {
		if can_move_here(outer, c) {
			return loc_depth(outer) - loc_depth(where)
		}
		outer = subloc(outer)
	}
	return 0
}

// v_move is the MOVE command handler.
// Handles land movement to a destination or direction.
// Ported from src/move.c lines 584-682.
func v_move(c *command) int {
	where := subloc(c.who)
	checkOuter := true

	if numargs(c) < 1 {
		wout(c.who, "Specify direction or destination to MOVE.")
		return FALSE
	}

	var v *exit_view
	for numargs(c) > 0 {
		v = parse_exit_dir(c, where, "move")
		if v != nil {
			checkOuter = false
			if land_check(c, v, true) {
				break
			}
			v = nil
		}
		cmd_shift(c)
	}

	if v == nil && checkOuter && can_move_at_outer_level(where, c) != 0 {
		c.a = subloc(where)
		v = parse_exit_dir(c, where, "")
		if v != nil {
			if move_exit_land(c, v, false) >= 0 {
				wout(c.who, "(assuming 'move out' first)")
				prepend_order(player(c.who), c.who, cmd_to_string(c))
			}
		}
	}

	if v == nil {
		return FALSE
	}

	delay := move_exit_land(c, v, true)
	if delay < 0 {
		return FALSE
	}

	if v.hades_cost != 0 {
		n := count_stack_figures(c.who)
		cost := v.hades_cost * n

		log_write(LOG_SPECIAL, "%s (%s) tries to enter Hades",
			box_name(player(c.who)), box_name(c.who))

		if !autocharge(c.who, cost) {
			wout(c.who, "Can't afford %s to enter Hades.", gold_s(cost))
			return FALSE
		}

		wout(c.who, "The Gatekeeper Spirit of Hades took %s from us.", gold_s(cost))
	}

	v.distance = delay
	c.wait = delay
	save_v_array(c, v)
	leave_stack(c.who)

	if delay > 0 {
		vector_stack(c.who, true)
		wout(VECT, "Travel to %s will take %s day%s.",
			box_name(v.destination),
			nice_num(delay),
			add_s(delay))
	}

	suspend_stack_actions(c.who)
	clear_guard_flag(c.who)

	if delay > 1 {
		prisoner_movement_escape_check(c.who)
	}

	departure_message(c.who, v)
	return TRUE
}

// d_move is the completion handler for MOVE command.
// Called when land movement completes.
// Ported from src/move.c lines 769-825.
func d_move(c *command) int {
	var v exit_view
	restore_v_array(c, &v)

	if !valid_box(v.destination) {
		wout(c.who, "Your destination no longer exists!")
		return FALSE
	}

	if v.road != 0 {
		discover_road(c.who, subloc(c.who), &v)
	}

	vector_stack(c.who, true)
	wout(VECT, "Arrival at %s.", box_name(v.destination))

	if loc_depth(v.destination) == LOC_province &&
		viewloc(subloc(c.who)) != viewloc(v.destination) &&
		weather_here(v.destination, sub_fog) != 0 {
		wout(VECT, "The province is blanketed in fog.")
	}

	restore_stack_actions_impl(c.who)
	move_stack_impl(c.who, v.destination)

	if viewloc(v.orig) != viewloc(v.destination) {
		arrival_message(c.who, &v)
	}

	return TRUE
}

// touch_loc_after_move touches a location for all characters in a stack.
// This updates visibility and triggers location events.
// Ported from src/move.c lines 685-700.
func touch_loc_after_move(who, where int) {
	if kind(who) == T_char {
		touch_loc(who)
	}

	var chars []int
	loop_char_here(who, &chars)
	for _, i := range chars {
		if !is_prisoner(i) {
			touch_loc(i)
		}
	}
}

// move_stack_impl is the core implementation that moves a character and stack to a new location.
// Updates location, known information, weather, and checks for special region effects.
// Ported from src/move.c lines 703-766.
func move_stack_impl(who, where int) {
	if kind(who) != T_char {
		panic("move_stack_impl: who is not a character")
	}

	if !in_faery(subloc(who)) && in_faery(where) {
		log_write(LOG_SPECIAL, "%s enters Faery at %s.", box_name(who), box_name(where))
	}

	set_where(who, where)
	mark_loc_stack_known(who, where)
	touch_loc_after_move(who, where)
	update_weather_view_locs(who, where)
	clear_contacts(who)

	if subkind(where) == sub_city {
		var stackMembers []int
		loop_stack(who, &stackMembers)
		for _, i := range stackMembers {
			match_trades(i)
		}
	}

	if subkind(where) == sub_ocean {
		p := p_char(who)
		if p.time_flying == 0 {
			p.time_flying++
			ocean_chars = append(ocean_chars, who)
		}
	}

	if subkind(where) != sub_ocean {
		p := rp_char(who)
		if p != nil && p.time_flying != 0 {
			p.time_flying = 0
			for i, ch := range ocean_chars {
				if ch == who {
					ocean_chars = append(ocean_chars[:i], ocean_chars[i+1:]...)
					break
				}
			}
		}
	}

	if loc_depth(where) == LOC_province &&
		subkind(where) != sub_ocean &&
		in_faery(where) {
		faery_attack_check(who, where)
	}

	if loc_depth(where) == LOC_province &&
		subkind(where) != sub_ocean &&
		in_hades(where) {
		hades_attack_check(who, where)
	}
}

// init_ocean_chars initializes the ocean_chars list at game start.
// Ported from src/move.c lines 828-842.
func init_ocean_chars() {
	ocean_chars = nil
	for _, i := range teg.Characters() {
		where := subloc(i)
		if subkind(where) == sub_ocean {
			ocean_chars = append(ocean_chars, i)
		}
	}
}

// check_ocean_chars checks flying units over ocean for drowning.
// Units that have been flying over ocean for too long will plunge into the sea.
// Ported from src/move.c lines 845-888.
func check_ocean_chars() {
	charsCopy := make([]int, len(ocean_chars))
	copy(charsCopy, ocean_chars)

	for _, who := range charsCopy {
		where := subloc(who)
		p := p_char(who)

		if !alive(who) || subkind(where) != sub_ocean {
			p.time_flying = 0
			for i, ch := range ocean_chars {
				if ch == who {
					ocean_chars = append(ocean_chars[:i], ocean_chars[i+1:]...)
					break
				}
			}
			continue
		}

		p.time_flying++

		if p.time_flying <= 15 {
			continue
		}

		if stack_parent(who) != 0 {
			continue
		}

		vector_stack(who, true)
		wout(VECT, "Flight can no longer be maintained. %s plunges into the sea.", box_name(who))

		kill_stack_ocean(who)
	}
}

// fly_check validates if flying to a destination is possible.
// Ported from src/move.c lines 891-912.
func fly_check(c *command, v *exit_view) bool {
	if v.in_transit != 0 {
		wout(c.who, "%s is underway. Boarding is not possible.", box_name(v.destination))
		return false
	}

	if loc_depth(v.destination) == LOC_build &&
		subkind(v.destination) != sub_sewer {
		owner := building_owner(v.destination)
		if owner != 0 && !will_admit(owner, c.who, v.destination) && v.direction != DIR_OUT {
			wout(c.who, "%s refused to let us enter.", box_name(owner))
			wout(owner, "Refused to let %s enter.", box_name(c.who))
			return false
		}
	}

	return true
}

// can_fly_here checks if a character can fly in a direction from where.
// Ported from src/move.c lines 915-928.
func can_fly_here(where int, c *command) bool {
	v := parse_exit_dir(c, where, "")
	if v != nil && v.direction != DIR_IN && move_exit_fly(c, v, false) >= 0 {
		return true
	}
	return false
}

// can_fly_at_outer_level checks if flying is possible from an outer location.
// Ported from src/move.c lines 931-945.
func can_fly_at_outer_level(where int, c *command) int {
	outer := subloc(where)
	for loc_depth(outer) > LOC_region {
		if can_fly_here(outer, c) {
			return loc_depth(outer) - loc_depth(where)
		}
		outer = subloc(outer)
	}
	return 0
}

// v_fly is the FLY command handler.
// Handles flying movement to a destination or direction.
// Ported from src/move.c lines 948-1021.
func v_fly(c *command) int {
	where := subloc(c.who)
	checkOuter := true

	if numargs(c) < 1 {
		wout(c.who, "Specify direction or destination to FLY.")
		return FALSE
	}

	var v *exit_view
	for numargs(c) > 0 {
		v = parse_exit_dir(c, where, "fly")
		if v != nil {
			checkOuter = false
			if fly_check(c, v) {
				break
			}
			v = nil
		}
		cmd_shift(c)
	}

	if v == nil && checkOuter && can_fly_at_outer_level(where, c) != 0 {
		c.a = subloc(where)
		v = parse_exit_dir(c, where, "")
		if v != nil {
			if move_exit_fly(c, v, false) >= 0 {
				wout(c.who, "(assuming 'fly out' first)")
				prepend_order(player(c.who), c.who, cmd_to_string(c))
			}
		}
	}

	if v == nil {
		return FALSE
	}

	delay := move_exit_fly(c, v, true)
	if delay < 0 {
		return FALSE
	}

	v.distance = delay
	c.wait = delay
	save_v_array(c, v)
	leave_stack(c.who)

	if delay > 0 {
		vector_stack(c.who, true)
		wout(VECT, "Flying to %s will take %s day%s.",
			box_name(v.destination),
			nice_num(delay),
			add_s(delay))
	}

	suspend_stack_actions(c.who)
	clear_guard_flag(c.who)
	departure_message(c.who, v)

	return TRUE
}

// d_fly is the completion handler for FLY command.
// Ported from src/move.c lines 1024-1072.
func d_fly(c *command) int {
	var v exit_view
	restore_v_array(c, &v)

	if v.road != 0 {
		discover_road(c.who, subloc(c.who), &v)
	}

	vector_stack(c.who, true)
	wout(VECT, "Arrival at %s.", box_name(v.destination))

	if loc_depth(v.destination) == LOC_province &&
		viewloc(subloc(c.who)) != viewloc(v.destination) &&
		weather_here(v.destination, sub_fog) != 0 {
		wout(VECT, "The province is blanketed in fog.")
	}

	restore_stack_actions_impl(c.who)
	move_stack_impl(c.who, v.destination)
	arrival_message(c.who, &v)

	return TRUE
}

// v_exit is the EXIT command handler.
// Synonym for 'move out'.
// Ported from src/move.c lines 1079-1088.
func v_exit(c *command) int {
	c.a = 0 // clear destination, we'll parse "out"
	// Set up parse for "move out"
	oly_parse(c, "move out")
	return v_move(c)
}

// v_enter is the ENTER command handler.
// Enters a sublocation or 'move in'.
// Ported from src/move.c lines 1091-1107.
func v_enter(c *command) int {
	if numargs(c) < 1 {
		oly_parse(c, "move in")
		return v_move(c)
	}
	oly_parse(c, sout("move %s", get_parse_arg(c, 1)))
	return v_move(c)
}

// v_north is the NORTH command handler.
// Ported from src/move.c lines 1110-1119.
func v_north(c *command) int {
	oly_parse(c, "move north")
	return v_move(c)
}

// v_south is the SOUTH command handler.
// Ported from src/move.c lines 1122-1131.
func v_south(c *command) int {
	oly_parse(c, "move south")
	return v_move(c)
}

// v_east is the EAST command handler.
// Ported from src/move.c lines 1134-1143.
func v_east(c *command) int {
	oly_parse(c, "move east")
	return v_move(c)
}

// v_west is the WEST command handler.
// Ported from src/move.c lines 1146-1155.
func v_west(c *command) int {
	oly_parse(c, "move west")
	return v_move(c)
}

// check_captain_loses_sailors is defined in cmd_transfer.go

// move_exit_water calculates the delay for sailing movement.
// Ported from src/move.c lines 1230-1301.
func move_exit_water(c *command, v *exit_view, ship int, show bool) int {
	delay := v.distance
	handsShort := 0
	where := subloc(ship)

	switch subkind(ship) {
	case sub_roundship:
		n := has_item(c.who, item_sailor) + has_item(c.who, item_pirate)
		if n < 8 {
			handsShort = 8 - n
			s := "day"
			if handsShort > 1 {
				s = sout("%s days", nice_num(handsShort))
			}
			if show {
				nStr := "none"
				if n > 0 {
					nStr = nice_num(n)
				}
				wout(c.who, "The crew of a roundship is eight sailors, but you have %s. Travel will take an extra %s.", nStr, s)
			}
		}

	case sub_galley:
		n := has_item(c.who, item_sailor) + has_item(c.who, item_pirate)
		if n < 14 {
			handsShort = 14 - n
			s := "day"
			if handsShort > 1 {
				s = sout("%s days", nice_num(handsShort))
			}
			if show {
				nStr := "none"
				if n > 0 {
					nStr = nice_num(n)
				}
				wout(c.who, "The crew of a galley is fourteen slaves or sailors, but you have %s. Travel will take an extra %s.", nStr, s)
			}
		}

	default:
		return -1
	}

	windBonus := 0
	if subkind(ship) == sub_roundship && weather_here(where, sub_wind) != 0 && delay > 2 {
		windBonus = 1
		if show {
			wout(c.who, "Favorable winds speed our progress.")
		}
	}

	delay = delay + handsShort - windBonus
	return delay
}

// sail_depart_message outputs a message when a ship departs.
// Ported from src/move.c lines 1304-1320.
func sail_depart_message(ship int, v *exit_view) {
	desc := liner_desc(ship)

	to := ""
	if v.dest_hidden == 0 {
		to = sout(" for %s.", box_name(v.destination))
	}

	comma := ""
	if strings.Contains(desc, ",") {
		comma = ","
	}

	wout(v.orig, "%s%s departed%s", desc, comma, to)
}

// sail_arrive_message outputs a message when a ship arrives.
// Ported from src/move.c lines 1323-1350.
func sail_arrive_message(ship int, v *exit_view) {
	desc := liner_desc(ship)

	from := ""
	if v.orig_hidden == 0 {
		from = sout(" from %s", box_name(v.orig))
	}

	comma := ""
	if strings.Contains(desc, ",") {
		comma = ","
	}

	with := display_owner(ship)
	if with == "" {
		with = "."
	}

	show_to_garrison = true
	wout(v.destination, "%s%s arrived%s%s", desc, comma, from, with)
	show_owner_stack(v.destination, ship)
	show_to_garrison = false
}

// sail_check validates if sailing in a direction is possible.
// Ported from src/move.c lines 1353-1373.
func sail_check(c *command, v *exit_view, show bool) bool {
	if v.water == 0 {
		if show {
			wout(c.who, "There is no water route in that direction.")
		}
		return false
	}

	if v.impassable != 0 {
		if show {
			wout(c.who, "That route is impassable.")
		}
		return false
	}

	return true
}

// can_sail_here checks if a ship can sail in a direction from where.
// Ported from src/move.c lines 1376-1391.
func can_sail_here(where int, c *command, ship int) bool {
	v := parse_exit_dir(c, where, "")
	if v != nil && v.direction != DIR_IN && sail_check(c, v, false) && move_exit_water(c, v, ship, false) >= 0 {
		return true
	}
	return false
}

// can_sail_at_outer_level checks if sailing is possible from an outer location.
// Ported from src/move.c lines 1394-1432.
func can_sail_at_outer_level(ship, where int, c *command) int {
	if ship_cap(ship) > 0 {
		loaded := ship_weight(ship) * 100 / ship_cap(ship)
		if loaded > 100 {
			wout(c.who, "%s is too overloaded to sail.", box_name(ship))
			wout(c.who, "(ship capacity = %d, damage = %d, load = %d)",
				ship_cap_raw(ship), loc_damage(ship), ship_weight(ship))
			return 0
		}
	}

	outer := subloc(where)
	for loc_depth(outer) > LOC_region {
		if can_sail_here(outer, c, ship) {
			return loc_depth(outer) - loc_depth(where)
		}
		outer = subloc(outer)
	}

	return 0
}

// v_sail is the SAIL command handler.
// Handles ship sailing to a destination or direction.
// Ported from src/move.c lines 1435-1572.
func v_sail(c *command) int {
	ship := subloc(c.who)
	checkOuter := true

	if !is_ship(ship) {
		sk := subkind(ship)
		if sk == sub_galley_notdone || sk == sub_roundship_notdone {
			wout(c.who, "%s is not yet completed.", box_name(ship))
		} else {
			wout(c.who, "Must be on a sea-worthy ship to sail.")
		}
		return FALSE
	}

	if building_owner(ship) != c.who {
		wout(c.who, "Only the captain of a ship may sail.")
		return FALSE
	}

	if !has_skill_check(c.who, sk_pilot_ship) {
		wout(c.who, "Knowledge of %s is required to sail.", box_name(sk_pilot_ship))
		return FALSE
	}

	if numargs(c) < 1 {
		wout(c.who, "Specify direction or destination to sail.")
		return FALSE
	}

	outerLoc := subloc(ship)

	var v *exit_view
	for numargs(c) > 0 {
		v = parse_exit_dir(c, outerLoc, "sail")
		if v != nil {
			checkOuter = false
			if sail_check(c, v, true) {
				break
			}
			v = nil
		}
		cmd_shift(c)
	}

	if v == nil && checkOuter && can_sail_at_outer_level(ship, outerLoc, c) != 0 {
		c.a = subloc(outerLoc)
		v = parse_exit_dir(c, outerLoc, "")
		if v != nil {
			if move_exit_water(c, v, ship, false) >= 0 {
				wout(c.who, "(assuming 'sail out' first)")
				prepend_order(player(c.who), c.who, cmd_to_string(c))
			}
		}
	}

	if v == nil {
		return FALSE
	}

	if ship_cap(ship) > 0 {
		loaded := ship_weight(ship) * 100 / ship_cap(ship)
		if loaded > 100 {
			wout(c.who, "%s is too overloaded to sail.", box_name(ship))
			wout(c.who, "(ship capacity = %d, damage = %d, load = %d)",
				ship_cap_raw(ship), loc_damage(ship), ship_weight(ship))
			return FALSE
		}
	}

	if v.in_transit != 0 {
		wout(c.who, "Cannot sail while already in transit.")
		return FALSE
	}

	delay := move_exit_water(c, v, ship, true)
	if delay < 0 {
		return FALSE
	}

	c.wait = delay
	v.distance = delay
	save_v_array(c, v)

	if delay > 0 {
		vector_char_here(c.who)
		vector_add(c.who)
		wout(VECT, "Sailing to %s will take %s day%s.",
			box_name(v.destination),
			nice_num(delay),
			add_s(delay))
	}

	sail_depart_message(ship, v)

	p_subloc(ship).moving = teg.globals.sysclock.days_since_epoch

	if ferry_horn(ship) != 0 {
		p_magic(ship).ferry_flag = 0
	}

	return TRUE
}

// d_sail is the completion handler for SAIL command.
// Ported from src/move.c lines 1575-1615.
func d_sail(c *command) int {
	ship := subloc(c.who)

	if !is_ship(ship) {
		return FALSE
	}
	if building_owner(ship) != c.who {
		return FALSE
	}

	var v exit_view
	restore_v_array(c, &v)

	if v.road != 0 {
		discover_road(c.who, subloc(ship), &v)
	}

	vector_char_here(ship)
	wout(VECT, "Arrival at %s.", box_name(v.destination))

	if loc_depth(v.destination) == LOC_province &&
		viewloc(subloc(ship)) != viewloc(v.destination) &&
		weather_here(v.destination, sub_fog) != 0 {
		wout(VECT, "The province is blanketed in fog.")
	}

	p_subloc(ship).moving = 0
	set_where(ship, v.destination)
	mark_loc_stack_known(ship, v.destination)
	move_bound_storms(ship, v.destination)

	if ferry_horn(ship) != 0 {
		p_magic(ship).ferry_flag = 0
	}

	touch_loc_after_move(ship, v.destination)
	sail_arrive_message(ship, &v)

	if c.use_skill == 0 {
		add_skill_experience(c.who, sk_pilot_ship)
	}

	return TRUE
}

// i_sail is the interrupt handler for SAIL command.
// Called when sailing is interrupted to reset the ship's moving state.
// Ported from src/move.c lines 1623-1635.
func i_sail(c *command) int {
	ship := subloc(c.who)

	if !is_ship(ship) {
		return TRUE
	}

	p_subloc(ship).moving = 0

	if ferry_horn(ship) != 0 {
		p_magic(ship).ferry_flag = 0
	}

	return TRUE
}

// Helper functions

// count_stack_move_nobles counts the number of nobles in a stack for movement purposes.
// Ported from src/u.c lines 1076-1090.
func count_stack_move_nobles(who int) int {
	sum := 1
	pl := player(who)

	var chars []int
	loop_char_here(who, &chars)
	for _, i := range chars {
		if player(i) == pl || (!is_npc(i) && !is_prisoner(i)) {
			sum++
		}
	}
	return sum
}

// count_stack_figures counts the total figures in a stack.
// Ported from src/u.c lines 1093-1105.
func count_stack_figures(who int) int {
	sum := 0
	var stackMembers []int
	loop_stack(who, &stackMembers)
	for _, i := range stackMembers {
		sum += count_any(i)
	}
	return sum
}

// plural_man returns "man" or "men" based on count.
func plural_man(n int) string {
	if n == 1 {
		return "man"
	}
	return "men"
}

// add_s is defined in inventory.go

// cap_str is defined in code.go

// Stub functions for dependencies not yet implemented

// sout formats a string (simplified version for now).
func sout(format string, args ...any) string {
	// Use fmt.Sprintf when available
	return format // Placeholder - full implementation in later sprint
}

// viewloc is defined in loc.go

// garrison_here is defined in visibility.go

// garrison_notices returns true if a garrison notices a character.
// Stub for now.
func garrison_notices(garr, who int) bool {
	return false
}

// garrison_spot_check returns true if a garrison spot-checks a character.
// Stub for now.
func garrison_spot_check(garr, who int) bool {
	return false
}

// liner_desc is defined in lifecycle.go

// set_known is defined in knowledge.go

// mark_loc_stack_known is defined in loc.go

// update_weather_view_locs updates weather view for a stack.
// Stub for now.
func update_weather_view_locs(who, where int) {
}

// clear_contacts clears a character's contacts.
// Stub for now.
func clear_contacts(who int) {
}

// match_trades matches trades for a character at a city.
// Stub for now.
func match_trades(who int) {
}

// in_faery is defined in loc.go

// in_hades is defined in loc.go

// faery_attack_check checks for Faery attacks.
// Stub for now.
func faery_attack_check(who, where int) {
}

// hades_attack_check checks for Hades attacks.
// Stub for now.
func hades_attack_check(who, where int) {
}

// autocharge is defined in inventory.go

// prisoner_movement_escape_check is defined in stack.go

// cmd_shift is defined in cmd_meta.go

// prepend_order prepends an order to a character's order queue.
// Stub for now.
func prepend_order(pl, who int, line string) {
}

// cmd_to_string converts a command to its string representation.
// Stub for now.
func cmd_to_string(c *command) string {
	// c.line is *char (C-style), just return empty for now
	return ""
}

// oly_parse parses an order line into a command.
// Stub for now.
func oly_parse(c *command, line string) bool {
	return true
}

// find_command finds a command by name and returns its index.
// Stub for now.
func find_command(name string) int {
	return -1
}

// rp_command is defined in accessor.go

// move_bound_storms moves storms bound to a ship.
// Stub for now.
func move_bound_storms(ship, where int) {
}

// add_skill_experience adds experience for using a skill.
// Stub for now.
func add_skill_experience(who, skill int) {
}

// ferry_horn is defined in accessor.go (returns schar, use ferry_horn(x) != 0)

// display_owner returns a string describing the owner of a location.
// Stub for now.
func display_owner(who int) string {
	return ""
}

// show_owner_stack shows the owner's stack at a location.
// Stub for now.
func show_owner_stack(where, who int) {
}

// vector_clear clears the output vector.
// Stub for now.
func vector_clear() {
}

// vector_add is defined in destruction.go

// loop_stack collects all characters in a stack (including nested).
func loop_stack(who int, l *[]int) {
	*l = append(*l, who)
	p := rp_loc_info(who)
	if p == nil {
		return
	}
	for _, id := range p.here_list {
		if kind(id) == T_char {
			loop_stack(id, l)
		}
	}
}

// count_any counts a character plus their "men" items.
func count_any(who int) int {
	count := 1
	inv := teg.globals.inventories[who]
	for _, e := range inv {
		if is_man_item(e.item) {
			count += e.qty
		}
	}
	return count
}

// is_man_item returns true if an item represents "men" (soldiers, peasants, etc.)
func is_man_item(item int) bool {
	p := rp_item(item)
	if p != nil {
		return p.is_man_item != 0
	}
	return false
}

// is_npc is defined in lifecycle.go

// alive is defined in accessor.go

// indent is the current output indentation level.
var indent = 0
