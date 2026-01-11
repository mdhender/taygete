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

// stack.go - Stacking logic ported from src/stack.c

package taygete

// here_pos returns the position of who in its location's here_list.
// Panics if who is not found.
// Ported from src/stack.c lines 8-23.
func here_pos(who int) int {
	p := rp_loc_info(loc(who))
	if p == nil {
		panic("here_pos: nil loc_info")
	}

	ret := IListLookup(p.here_list, who)
	if ret < 0 {
		panic("here_pos: who not in here_list")
	}

	return ret
}

// here_precedes returns true if a comes before b in the here list.
// Returns false if they're in different locations.
// Ported from src/stack.c lines 30-51.
func here_precedes(a, b int) bool {
	if loc(a) != loc(b) {
		return false
	}

	p := rp_loc_info(loc(a))
	if p == nil {
		panic("here_precedes: nil loc_info")
	}

	for _, id := range p.here_list {
		if id == a {
			return true
		} else if id == b {
			return false
		}
	}

	panic("here_precedes: neither a nor b found in here_list")
}

// first_prisoner_pos returns the index of the first prisoner in where's here_list.
// Returns -1 if no prisoner is found.
// Ported from src/stack.c lines 54-73.
func first_prisoner_pos(where int) int {
	p := rp_loc_info(where)
	if p == nil {
		return -1
	}

	for i, id := range p.here_list {
		if kind(id) == T_char && is_prisoner(id) {
			return i
		}
	}

	return -1
}

// stack_parent returns the character that who is stacked under.
// Returns 0 if who is not stacked under a character.
// Ported from src/stack.c lines 76-87.
func stack_parent(who int) int {
	n := loc(who)
	if kind(n) == T_char {
		return n
	}
	return 0
}

// stack_leader returns the topmost character in a stack.
// Walks up the stack until it finds a character not stacked under another character.
// Ported from src/stack.c lines 90-108.
func stack_leader(who int) int {
	if kind(who) != T_char {
		panic("stack_leader: who is not a character")
	}

	count := 0
	n := who
	for kind(n) == T_char {
		who = n
		n = stack_parent(n)

		count++
		if count >= 1000 {
			panic("stack_leader: infinite loop detected")
		}
	}

	return who
}

// stacked_beneath returns true if b is stacked somewhere beneath a.
// Both a and b must be characters.
// Ported from src/stack.c lines 115-133.
func stacked_beneath(a, b int) bool {
	if kind(a) != T_char {
		panic("stacked_beneath: a is not a character")
	}
	if kind(b) != T_char {
		panic("stacked_beneath: b is not a character")
	}

	if a == b {
		return false
	}

	for b > 0 {
		b = stack_parent(b)
		if a == b {
			return true
		}
	}

	return false
}

// promote moves who to position new_pos in its location's here_list.
// who must already be at or after new_pos.
// Ported from src/stack.c lines 136-153.
func promote(who, new_pos int) {
	p := rp_loc_info(loc(who))
	if p == nil {
		panic("promote: nil loc_info")
	}

	who_pos := IListLookup(p.here_list, who)
	if who_pos < new_pos {
		panic("promote: who_pos < new_pos")
	}

	for i := who_pos; i > new_pos; i-- {
		p.here_list[i] = p.here_list[i-1]
	}
	p.here_list[new_pos] = who
}

// unstack removes who from its current stack and places it in the subloc.
// Ported from src/stack.c lines 156-187.
func unstack(who int) {
	leader := stack_leader(who)
	if !valid_box(leader) {
		panic("unstack: invalid leader")
	}

	if subloc(leader) != loc(leader) {
		panic("unstack: subloc(leader) != loc(leader)")
	}

	if release_swear(who) != 0 {
		p_magic(who).swear_on_release = 0
	}

	set_where(who, subloc(leader))
	promote(who, here_pos(leader)+1)

	restore_stack_actions(who)

	if loyal_kind(who) == LOY_summon {
		set_loyal(who, LOY_npc, 0)
	}
}

// leave_stack removes who from their current stack (if any) and outputs messages.
// Ported from src/stack.c lines 190-207.
func leave_stack(who int) {
	leader := stack_parent(who)
	if leader <= 0 {
		return
	}

	wout(leader, "%s unstacks from us.", box_name(who))
	wout(who, "%s unstacks from %s.", box_name(who), box_name(leader))

	vector_char_here(who)
	wout(VECT, "%s unstacks from %s.", box_name(who), box_name(leader))

	unstack(who)
}

// stack places who under target in the stacking hierarchy.
// who must not already be stacked under a character.
// Ported from src/stack.c lines 210-231.
func stack(who, target int) {
	if stack_parent(who) != 0 {
		panic("stack: who is already stacked")
	}

	set_where(who, target)
	p_char(who).moving = char_moving(target)

	if !is_prisoner(who) {
		pos := first_prisoner_pos(target)
		if pos >= 0 {
			promote(who, pos)
		}
	}
}

// join_stack makes who stack beneath target.
// Ported from src/stack.c lines 234-252.
func join_stack(who, target int) {
	if stacked_beneath(who, target) {
		panic("join_stack: who is already beneath target")
	}
	if is_prisoner(target) {
		panic("join_stack: target is a prisoner")
	}

	leave_stack(who)

	if subloc(target) != subloc(who) {
		panic("join_stack: subloc mismatch")
	}

	wout(who, "%s stacks beneath %s.", box_name(who), box_name(target))
	wout(target, "%s stacks beneath us.", box_name(who))

	vector_char_here(who)
	wout(VECT, "%s stacks beneath %s.", box_name(who), box_name(target))

	stack(who, target)
}

// check_prisoner_escape checks if a prisoner escapes.
// Returns true if the prisoner escapes, false otherwise.
// Ported from src/stack.c lines 255-286.
func check_prisoner_escape(who, chance int) bool {
	leader := stack_parent(who)
	hound := stack_has_item(who, item_hound)

	chance *= 10

	if hound > 0 {
		chance /= 2
	}

	n := rnd(1, 1000)

	if n > chance {
		if hound > 0 {
			vector_stack(leader, true)
			if hound == 1 {
				wout(VECT, "The hound is barking.")
			} else {
				wout(VECT, "The hounds are barking.")
			}
			return false
		}
		return false
	}

	prisoner_escapes(who)
	return true
}

// prisoner_escapes handles a prisoner escaping.
// Ported from src/stack.c lines 289-328.
func prisoner_escapes(who int) {
	leader := stack_parent(who)

	wout(leader, "Prisoner %s escaped!", box_name(who))
	p_char(who).prisoner = 0
	p_magic(who).swear_on_release = 0
	unstack(who)
	touch_loc(who)

	wout(who, "We escaped!")

	where := subloc(who)

	if loc_depth(where) <= LOC_province {
		return
	}

	out_one := loc(where)

	if is_ship(where) && subkind(out_one) == sub_ocean {
		out_one = find_nearest_land(out_one)

		wout(who, "After jumping over the side of the boat and "+
			"enduring a long, grueling, swim, we finally "+
			"washed ashore at %s.", box_name(out_one))

		wout(leader, "%s jumped overboard and presumably "+
			"drowned.", just_name(who))

		log_write(LOG_SPECIAL, "!! Someone swam ashore, who=%s",
			box_code_less(who))
	}

	move_stack(who, out_one)
}

// prisoner_movement_escape_check checks all prisoners under who for escape.
// Ported from src/stack.c lines 331-342.
func prisoner_movement_escape_check(who int) {
	var chars []int
	all_char_here(who, &chars)

	for _, i := range chars {
		if is_prisoner(i) {
			check_prisoner_escape(i, 2)
		}
	}
}

// weekly_prisoner_escape_check performs weekly escape checks for all prisoners.
// Ported from src/stack.c lines 345-371.
func weekly_prisoner_escape_check() {
	for who := kind_first(T_char); who > 0; who = kind_next(who) {
		if is_prisoner(who) {
			continue
		}

		if subkind(subloc(who)) == sub_ocean {
			continue
		}

		var hereList []int
		all_here(who, &hereList)

		for _, i := range hereList {
			if kind(i) == T_char && is_prisoner(i) && release_swear(i) == 0 {
				chance := 2
				if loc_depth(subloc(who)) >= LOC_build {
					chance = 1
				}
				check_prisoner_escape(i, chance)
			}
		}
	}
}

// drop_stack drops to_drop from who's stack.
// Ported from src/stack.c lines 374-440.
func drop_stack(who, to_drop int) {
	if stack_parent(to_drop) != who {
		panic("drop_stack: to_drop is not stacked under who")
	}

	release_swear_flag := false

	if is_prisoner(to_drop) {
		p_char(to_drop).prisoner = 0
		touch_loc(to_drop)
		wout(who, "Freed prisoner %s.", box_name(to_drop))
		wout(to_drop, "%s set us free.", box_name(who))

		if release_swear(to_drop) != 0 {
			release_swear_flag = true
		}
	} else {
		wout(who, "Dropped %s from stack.", box_name(to_drop))
		wout(to_drop, "%s dropped us from the stack.", box_name(who))

		vector_char_here(to_drop)
		wout(VECT, "%s dropped %s from the stack.", box_name(who), box_name(to_drop))
	}

	unstack(to_drop)

	if release_swear_flag {
		log_write(LOG_SPECIAL, "%s frees a swear_on_release prisoner", box_name(who))

		if rnd(1, 5) < 5 {
			wout(who, "%s is grateful for your gallantry.", box_name(to_drop))
			wout(who, "%s pledges fealty to us.", box_name(to_drop))

			set_lord(to_drop, player(who), LOY_oath, 1)
		} else {
			switch rnd(1, 3) {
			case 1:
				wout(who, "%s spits on you, and vanishes in a cloud of orange smoke.", box_name(to_drop))
			case 2:
				wout(who, "%s cackles wildly and vanishes.", box_name(to_drop))
			case 3:
				wout(who, "%s smiles briefly at you, then vanishes.", box_name(to_drop))
			}

			unit_deserts(to_drop, 0, true, LOY_unsworn, 0)
			put_back_cookie(to_drop)
			set_where(to_drop, 0)
			change_box_kind(to_drop, T_deadchar)
		}
	}
}

// free_all_prisoners frees all prisoners stacked under who.
// Ported from src/stack.c lines 443-454.
func free_all_prisoners(who int) {
	var hereList []int
	all_here(who, &hereList)

	for _, i := range hereList {
		if kind(i) == T_char && is_prisoner(i) {
			drop_stack(who, i)
		}
	}
}

// extract_stacked_unit removes who from a stack, leaving those above and below behind.
// Ported from src/stack.c lines 464-518.
func extract_stacked_unit(who int) {
	var first int

	var hereList []int
	all_here(who, &hereList)

	for _, i := range hereList {
		if kind(i) == T_char && !is_prisoner(i) {
			first = i
			break
		}
	}

	if first != 0 {
		all_here(who, &hereList)
		for _, i := range hereList {
			if i != first && kind(i) == T_char {
				set_where(i, first)
			}
		}

		set_where(first, loc(who))
		promote(first, here_pos(who)+1)
	}

	all_here(who, &hereList)
	for _, i := range hereList {
		if kind(i) == T_char && is_prisoner(i) {
			prisoner_escapes(i)
		}
	}

	leave_stack(who)
}

// promote_stack promotes lower to be before higher in the location order.
// Ported from src/stack.c lines 525-556.
func promote_stack(lower, higher int) {
	if stacked_beneath(lower, higher) {
		panic("promote_stack: lower is beneath higher")
	}

	set_where(lower, loc(higher))

	p := rp_loc_info(loc(higher))
	if p == nil {
		panic("promote_stack: nil loc_info")
	}

	if p.here_list[len(p.here_list)-1] != lower {
		panic("promote_stack: lower is not last in here_list")
	}

	pos := IListLookup(p.here_list, higher)
	promote(lower, pos)

	wout(higher, "Promoted %s.", box_name(lower))
	wout(lower, "%s promoted us.", box_name(higher))
}

// take_prisoner makes who take target prisoner.
// Ported from src/stack.c lines 559-610.
func take_prisoner(who, target int) {
	if who == target {
		panic("take_prisoner: who == target")
	}
	if kind(who) != T_char {
		panic("take_prisoner: who is not a character")
	}
	if kind(target) != T_char {
		panic("take_prisoner: target is not a character")
	}

	ni := false
	if subkind(target) == sub_ni && beast_capturable(target) {
		ni = true
	}

	vector_stack(stack_leader(target), true)
	vector_add(who)

	if ni {
		wout(VECT, "%s disbands.", box_name(target))
	} else {
		wout(VECT, "%s is taken prisoner by %s.", box_name(target), box_name(who))
	}

	p_char(target).prisoner = 1

	if ni {
		take_unit_items(target, who, TAKE_NI)
	} else {
		take_unit_items(target, who, TAKE_ALL)
	}

	extract_stacked_unit(target)
	interrupt_order(target)

	if ni {
		unit_deserts(target, 0, true, LOY_unsworn, 0)
		put_back_cookie(target)
		set_where(target, 0)
		change_box_kind(target, T_deadchar)
	} else {
		stack(target, who)
	}
}

// has_prisoner returns true if who has pris as a prisoner.
// Ported from src/stack.c lines 613-630.
func has_prisoner(who, pris int) bool {
	var hereList []int
	all_here(who, &hereList)

	for _, i := range hereList {
		if i == pris && is_prisoner(i) {
			return true
		}
	}

	return false
}

// move_prisoner moves pris from who to target.
// Ported from src/stack.c lines 633-647.
func move_prisoner(who, target, pris int) int {
	rs := release_swear(pris)

	unstack(pris)
	stack(pris, target)

	if rs != 0 {
		p_magic(pris).swear_on_release = 1
	}

	if rp_char(pris).prisoner == 0 {
		panic("move_prisoner: pris is not a prisoner")
	}

	return 0
}

// give_prisoner transfers pris from who to target.
// Returns true on success.
// Ported from src/stack.c lines 650-666.
func give_prisoner(who, target, pris int) bool {
	if check_prisoner_escape(pris, 2) {
		return false
	}

	move_prisoner(who, target, pris)

	wout(who, "Transferred prisoner %s to %s.", box_name(pris), box_name(target))
	wout(target, "%s transferred the prisoner %s to us.", box_name(who), box_name(pris))

	return true
}

// v_stack is the STACK command handler.
// Ported from src/stack.c lines 669-707.
func v_stack(c *command) int {
	target := c.a

	if !check_char_gone(c.who, target) {
		return FALSE
	}

	if target == c.who {
		wout(c.who, "Can't stack beneath oneself.")
		return FALSE
	}

	if stacked_beneath(c.who, target) {
		wout(c.who, "Cannot stack beneath %s since %s is stacked under you.",
			box_name(target), just_name(target))
		return FALSE
	}

	if is_prisoner(target) {
		wout(c.who, "Can't stack beneath prisoners.")
		return FALSE
	}

	if !will_admit(target, c.who, target) {
		wout(c.who, "%s refuses to let us stack.", box_name(target))
		wout(target, "Refused to let %s stack with us.", box_name(c.who))
		return FALSE
	}

	join_stack(c.who, target)
	return TRUE
}

// v_unstack is the UNSTACK command handler.
// Ported from src/stack.c lines 710-753.
func v_unstack(c *command) int {
	target := c.a

	if numargs(c) < 1 {
		if stack_parent(c.who) <= 0 {
			wout(c.who, "Not stacked under anyone.")
			return FALSE
		}

		leave_stack(c.who)
		return TRUE
	}

	if c.who == target {
		extract_stacked_unit(c.who)
		return TRUE
	}

	if !valid_box(target) || stack_parent(target) != c.who {
		wout(c.who, "%s is not stacked beneath us.", get_parse_arg(c, 1))
		return FALSE
	}

	drop_stack(c.who, target)
	return TRUE
}

// v_surrender is the SURRENDER command handler.
// Ported from src/stack.c lines 756-787.
func v_surrender(c *command) int {
	target := c.a

	if !check_char_gone(c.who, target) {
		return FALSE
	}

	if player(target) == player(c.who) {
		wout(c.who, "Can't surrender to oneself.")
		return FALSE
	}

	if is_prisoner(target) {
		wout(c.who, "Can't surrender to a prisoner.")
		return FALSE
	}

	log_write(LOG_SPECIAL, "Player %s surrenders %s",
		box_code_less(player(c.who)), box_name(c.who))

	vector_stack(stack_leader(c.who), true)
	vector_stack(stack_leader(target), false)

	wout(VECT, "%s surrenders to %s.", box_name(c.who), box_name(target))

	take_prisoner(target, c.who)
	return TRUE
}

// promote_after returns true if b appears later in location order than a.
// Ported from src/stack.c lines 794-827.
func promote_after(a, b int) bool {
	where := subloc(a)
	if subloc(b) != where {
		panic("promote_after: different sublocs")
	}

	var chars []int
	all_char_here(where, &chars)

	for _, i := range chars {
		if i == a {
			return true
		} else if i == b {
			return false
		}
	}

	panic("promote_after: neither a nor b found")
}

// v_promote is the PROMOTE command handler.
// Ported from src/stack.c lines 830-891.
func v_promote(c *command) int {
	target := c.a

	if numargs(c) < 1 {
		wout(c.who, "Must specify which character to promote.")
		return FALSE
	}

	if kind(target) != T_char {
		wout(c.who, "%s is not a character.", get_parse_arg(c, 1))
		return FALSE
	}

	if is_prisoner(target) {
		wout(c.who, "Can't promote prisoners.")
		return FALSE
	}

	targ_par := stack_parent(target)

	if !check_char_here(c.who, target) {
		return FALSE
	}

	if target == c.who {
		wout(c.who, "Can't promote oneself.")
		return FALSE
	}

	if player(c.who) == player(target) {
		if !promote_after(c.who, target) {
			wout(c.who, "%s already comes before us in location order.", box_name(target))
			return FALSE
		}
	} else if targ_par != c.who && !here_precedes(c.who, target) {
		wout(c.who, "Characters to be promoted must be stacked "+
			"immediately beneath the promoter, or be listed after the "+
			"promoter at the same level.")
		return FALSE
	}

	promote_stack(target, c.who)
	return TRUE
}

// Stub functions for dependencies not yet implemented.
// These will be implemented in later sprints.

// Deprecated: restore_stack_actions not yet implemented.
func restore_stack_actions(who int) {
	// TODO: Implement in later sprint (order processing)
}

// Deprecated: wout not yet implemented.
func wout(who int, format string, args ...any) {
	// TODO: Implement in later sprint (output/reporting)
}

// Deprecated: log_write not yet implemented.
func log_write(k int, format string, args ...any) {
	// TODO: Implement in later sprint (logging)
}

// Deprecated: vector_char_here not yet implemented.
func vector_char_here(where int) {
	// TODO: Implement in later sprint (output/reporting)
}

// Deprecated: vector_stack not yet implemented.
func vector_stack(who int, clear bool) {
	// TODO: Implement in later sprint (output/reporting)
}

// Note: vector_add implemented in destruction.go

// Deprecated: touch_loc not yet implemented.
func touch_loc(who int) {
	// TODO: Implement in later sprint (day.c)
}

// Note: find_nearest_land implemented in destruction.go

// Deprecated: move_stack not yet implemented.
func move_stack(who, dest int) {
	// TODO: Implement in later sprint (movement)
}

// set_lord sets the lord (owner) of a character.
// Ported from src/swear.c (simplified implementation for Sprint 25.8).
func set_lord(who, new_lord, k, lev int) {
	c := rp_char(who)
	if c == nil {
		return
	}

	oldLord := c.unit_lord
	c.unit_lord = new_lord
	c.loy_kind = schar(k)
	c.loy_rate = lev

	if oldLord != new_lord && oldLord != 0 {
		c.prev_lord = oldLord
	}
}

// Deprecated: set_loyal not yet implemented.
func set_loyal(who, k, lev int) {
	// TODO: Implement in later sprint (loyalty)
}

// unit_deserts handles a unit deserting to a new player.
// Ported from src/swear.c lines 304-343.
func unit_deserts(who, to_who int, loy_check bool, k, lev int) {
	sp := player(who)

	if to_who != 0 && sp != 0 {
		wout(sp, "%s renounces loyalty to us.", box_name(who))
		wout(who, "%s renounces loyalty.", box_name(who))
	}

	if to_who != 0 && is_prisoner(who) && player(to_who) == player(stack_parent(who)) {
		p_char(who).prisoner = FALSE
	} else if !is_prisoner(who) {
		extract_stacked_unit(who)
	}

	set_lord(who, to_who, k, lev)

	if to_who != 0 {
		wout(who, "%s pledges fealty to us.", box_name(who))
		wout(to_who, "%s pledges fealty to us.", box_name(who))
		p_char(who).new_lord = 1
	}
}

// char_reclaim marks a character for melting and triggers death.
// Used by QUIT/RECLAIM commands.
// Ported from src/u.c lines 82-92.
func char_reclaim(who int) {
	p_char(who).melt_me = TRUE
	kill_char(who, 0) // QUIT shouldn't give items to stackmates
}

// Note: put_back_cookie is implemented in lifecycle.go
// Note: take_unit_items is implemented in lifecycle.go
// Note: interrupt_order is implemented in lifecycle.go

// Note: check_char_gone is implemented in visibility.go
// Note: check_char_here is implemented in visibility.go

// Deprecated: will_admit not yet implemented.
func will_admit(who, target, where int) bool {
	// TODO: Implement in later sprint (permissions)
	return true
}

// numargs returns the number of arguments in a parsed command.
// The first element (index 0) is the command name, so numargs = len - 1.
func numargs(c *command) int {
	// TODO: parse is **char in C; for now return based on command args
	// This is a simplified implementation until order parsing is done.
	if c.a != 0 {
		return 1
	}
	return 0
}

// get_parse_arg returns the parsed argument at index i as a string.
func get_parse_arg(c *command, i int) string {
	// TODO: Implement properly when order parsing is done
	return ""
}

// stack_has_item is implemented in inventory.go

// Deprecated: beast_capturable not yet implemented.
func beast_capturable(who int) bool {
	// TODO: Implement in later sprint (beasts)
	return false
}
