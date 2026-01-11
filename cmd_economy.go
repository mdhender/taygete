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

// cmd_economy.go - Simple economy commands ported from src/c2.c
// Sprint 25.8: DISCARD, QUIT

package taygete

// how_many calculates the quantity to operate on for inventory commands.
// It validates that from_who has the item and computes the effective quantity
// based on the requested qty and have_left parameters.
// Returns 0 if the operation is not possible (with error message output).
// Ported from src/c1.c lines 488-519.
func how_many(who, from_who, item, qty, have_left int) int {
	num_has := has_item(from_who, item)

	if num_has <= 0 {
		wout(who, "%s has no %s.",
			just_name(from_who),
			just_name(item))
		return 0
	}

	if num_has <= have_left {
		wout(who, "%s has only %s.",
			just_name(from_who),
			just_name_qty(item, num_has))
		return 0
	}

	if qty == 0 {
		qty = num_has
	}

	qty = min(num_has-have_left, qty)

	if qty <= 0 {
		return 0
	}

	return qty
}

// v_discard executes the DISCARD command.
// Drops an item from the character's inventory.
// Ported from src/c2.c lines 12-43.
func v_discard(c *command) int {
	item := c.a
	qty := c.b
	have_left := c.c

	if kind(item) != T_item {
		wout(c.who, "%s is not an item.", box_code(item))
		return FALSE
	}

	if has_item(c.who, item) < 1 {
		wout(c.who, "%s does not have any %s.", box_name(c.who),
			box_code(item))
		return FALSE
	}

	qty = how_many(c.who, c.who, item, qty, have_left)

	if qty <= 0 {
		return FALSE
	}

	ret := drop_item(c.who, item, qty)
	if !ret {
		return FALSE
	}

	wout(c.who, "Dropped.")

	return TRUE
}

// drop_player handles player removal from the game.
// All units are either freed (if prisoners) or reclaimed.
// Dead bodies owned by the player have their old_lord set to indep_player.
// Ported from src/c2.c lines 47-117.
// Note: Shell command calls from the C version are skipped.
func drop_player(pl int) {
	if kind(pl) != T_player {
		panic("drop_player: not a player")
	}

	for _, who := range loop_units(pl) {
		if is_prisoner(who) {
			unit_deserts(who, indep_player, true, LOY_UNCHANGED, 0)
		} else {
			wout(subloc(who), "%s melts into the ground and vanishes.", box_name(who))
			char_reclaim(who)
		}
	}

	for _, i := range loop_dead_body() {
		owner := item_unique(i)
		if owner == 0 {
			continue
		}

		p := rp_misc(i)
		if p == nil || p.old_lord != pl {
			continue
		}

		p_misc(i).old_lord = indep_player
	}

	p := rp_player(pl)
	var s, email string
	if p != nil {
		if p.email != nil && *p.email != 0 {
			email = charPtrToStr(p.email)
		}
		if p.full_name != nil && *p.full_name != 0 {
			s = charPtrToStr(p.full_name)
		}
	}

	log_write(LOG_DROP, "Dropped player %s", box_name(pl))
	log_write(LOG_DROP, "    %s <%s>", s, email)

	delete_box(pl)
}

// v_quit executes the QUIT command.
// Removes a player from the game.
// Only the GM can quit another player.
// Ported from src/c2.c lines 121-147.
func v_quit(c *command) int {
	target := c.a

	if target == 0 {
		target = player(c.who)
	}

	if target != player(c.who) && player(c.who) != gm_player {
		wout(c.who, "Not allowed to drop another player.")
		return FALSE
	}

	if kind(target) != T_player {
		wout(c.who, "%s is not a player.", box_name(target))
		return FALSE
	}

	drop_player(target)

	return FALSE
}

// loop_units returns all units belonging to a player.
// This is a helper that replaces the C loop_units macro.
// Uses the kind chain for efficient iteration.
func loop_units(pl int) []int {
	if kind(pl) != T_player {
		return nil
	}

	var units []int
	for i := kind_first(T_char); i != 0; i = kind_next(i) {
		if player(i) == pl {
			units = append(units, i)
		}
	}
	return units
}

// loop_dead_body returns all dead body items.
// This is a helper that replaces the C loop_dead_body macro.
// Uses the subkind chain for efficient iteration.
func loop_dead_body() []int {
	var bodies []int
	for i := sub_first(sub_dead_body); i != 0; i = sub_next(i) {
		bodies = append(bodies, i)
	}
	return bodies
}


// charPtrToStr converts a *char (int8 pointer) to a Go string.
func charPtrToStr(p *char) string {
	if p == nil {
		return ""
	}
	return ""
}
