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

// cmd_transfer.go - Inventory transfer commands ported from src/c1.c and src/c2.c
// Sprint 26.5: Inventory Transfer Commands

package taygete

// v_accept sets accept rules for receiving items from other factions.
// Usage: ACCEPT <from_who> <item> <qty>
// All parameters can be 0 for wildcards.
// Ported from src/c1.c lines 397-417.
func v_accept(c *command) int {
	from_who := c.a
	item := c.b
	qty := c.c

	p := p_char(c.who)

	newEnt := &accept_ent{
		item:     item,
		from_who: from_who,
		qty:      qty,
	}

	p.accept = append(p.accept, newEnt)

	return TRUE
}

// will_accept_sup checks accept rules for a specific who/player.
// Returns true if the accept list contains a matching rule.
// Ported from src/c1.c lines 420-450.
func will_accept_sup(who, item, from, qty int) bool {
	p := rp_char(who)
	if p == nil {
		return false
	}

	for _, ae := range p.accept {
		item_match := ae.item == item || ae.item == 0
		from_match := ae.from_who == from || ae.from_who == 0
		qty_match := ae.qty >= qty || ae.qty == 0

		if item_match && from_match && qty_match {
			if ae.qty != 0 {
				ae.qty -= qty
			}
			return true
		}
	}

	return false
}

// will_accept checks if who will accept item from from.
// Gold is always accepted.
// Same faction is always accepted.
// Garrisons accept from rulers.
// Otherwise, checks ACCEPT rules on character and player.
// Ported from src/c1.c lines 453-484.
func will_accept(who, item, from, qty int) bool {
	if item == item_gold {
		return true
	}

	if player(who) == player(from) {
		return true
	}

	if subkind(who) == sub_garrison {
		if may_rule_here(from, who) {
			return true
		}
		wout(from, "%s is not under your control.", box_name(who))
		return false
	}

	if will_accept_sup(who, item, from, qty) ||
		will_accept_sup(player(who), item, from, qty) ||
		will_accept_sup(who, item, player(from), qty) ||
		will_accept_sup(player(who), item, player(from), qty) {
		return true
	}

	wout(who, "Refused %s from %s.", just_name_qty(item, qty), box_name(from))
	wout(from, "Refused by %s.", just_name(who))

	return false
}

// Note: how_many is defined in cmd_economy.go

// v_give gives items to a target character.
// Usage: GIVE <who> <what> [qty] [have-left]
// Handles argument correction if target and item are swapped.
// Ported from src/c1.c lines 526-628.
func v_give(c *command) int {
	target := c.a
	item := c.b
	qty := c.c
	have_left := c.d

	if numargs(c) >= 2 &&
		(kind(target) == T_item || has_prisoner(c.who, target)) &&
		(kind(item) == T_char && subloc(c.who) == subloc(item) && !has_prisoner(c.who, item)) {
		c.a, c.b = c.b, c.a
		target = c.a
		item = c.b

		switch numargs(c) {
		case 2:
			wout(c.who, "(assuming you meant 'give %d %d')", target, item)
		case 3:
			wout(c.who, "(assuming you meant 'give %d %d %d')", target, item, qty)
		default:
			wout(c.who, "(assuming you meant 'give %d %d %d %d')", target, item, qty, have_left)
		}
	}

	if !check_char_here(c.who, target) {
		return FALSE
	}
	if !check_char_gone(c.who, target) {
		return FALSE
	}

	if is_prisoner(target) {
		wout(c.who, "Prisoners may not be given anything.")
		return FALSE
	}

	if loyal_kind(target) == LOY_summon {
		wout(c.who, "Summoned entities may not be given anything.")
		return FALSE
	}

	if has_prisoner(c.who, item) {
		if !will_accept(target, item, c.who, 1) {
			return FALSE
		}
		if give_prisoner(c.who, target, item) {
			return TRUE
		}
		return FALSE
	}

	if kind(item) != T_item {
		wout(c.who, "%s is not an item or a prisoner.", box_code(item))
		return FALSE
	}

	if has_item(c.who, item) < 1 {
		wout(c.who, "%s does not have any %s.", box_name(c.who), box_code(item))
		return FALSE
	}

	qty = how_many(c.who, c.who, item, qty, have_left)
	if qty <= 0 {
		return FALSE
	}

	if !will_accept(target, item, c.who, qty) {
		return FALSE
	}

	ret := move_item(c.who, target, item, qty)
	if !ret {
		panic("v_give: move_item failed")
	}

	wout(c.who, "Gave %s to %s.", just_name_qty(item, qty), box_name(target))
	wout(target, "Received %s from %s.", box_name_qty(item, qty), box_name(c.who))

	return TRUE
}

// v_pay pays gold to a target character.
// Usage: PAY <who> <qty> [have-left]
// This is a wrapper around v_give that rewrites the command.
// Ported from src/c1.c lines 631-643.
func v_pay(c *command) int {
	target := c.a
	qty := c.b
	have_left := c.c

	c.a = target
	c.b = item_gold
	c.c = qty
	c.d = have_left

	return v_give(c)
}

// may_take checks if who may take items from target.
// Returns true if target is a controlled garrison, a prisoner, or same faction.
// Ported from src/c1.c lines 646-674.
func may_take(who, target int) bool {
	if !check_char_here(who, target) {
		return false
	}
	if !check_char_gone(who, target) {
		return false
	}

	if subkind(target) == sub_garrison {
		if may_rule_here(who, target) {
			return true
		}
		wout(who, "%s is not under your control.", box_name(target))
		return false
	}

	if !my_prisoner(who, target) && player(target) != player(who) {
		wout(who, "May only take items from other units in your faction.")
		return false
	}

	return true
}

// v_get takes items from a target character.
// Usage: GET <who> <what> [qty] [have-left]
// Ported from src/c1.c lines 681-742.
func v_get(c *command) int {
	target := c.a
	item := c.b
	qty := c.c
	have_left := c.d

	if !may_take(c.who, target) {
		return FALSE
	}

	if has_prisoner(target, item) {
		if give_prisoner(target, c.who, item) {
			return TRUE
		}
		return FALSE
	}

	if kind(item) != T_item {
		wout(c.who, "%s is not an item or a prisoner.", box_code(item))
		return FALSE
	}

	if has_item(target, item) < 1 {
		wout(c.who, "%s does not have any %s.", box_name(target), box_code(item))
		return FALSE
	}

	qty = how_many(c.who, target, item, qty, have_left)
	if qty <= 0 {
		return FALSE
	}

	if subkind(target) == sub_garrison && man_item(item) != 0 {
		garr_men := count_man_items(target)
		garr_men -= qty
		if garr_men < 10 {
			wout(c.who, "Garrisons must be left with a minimum of ten men.")
			return FALSE
		}
	}

	ret := move_item(target, c.who, item, qty)
	if !ret {
		panic("v_get: move_item failed")
	}

	wout(c.who, "Took %s from %s.", just_name_qty(item, qty), box_name(target))
	wout(target, "%s took %s from us.", box_name(c.who), box_name_qty(item, qty))

	if item == item_sailor || item == item_pirate {
		check_captain_loses_sailors(qty, target, c.who)
	}

	return TRUE
}

// v_claim claims items from the faction treasury.
// Usage: CLAIM <item> <qty> [have-left]
// Common mistake: CLAIM 500 is corrected to CLAIM 1 500 (claim gold).
// Cannot be used in Cloudlands, Hades, or Faery.
// Ported from src/c2.c lines 739-797.
func v_claim(c *command) int {
	item := c.a
	qty := c.b
	have_left := c.c
	pl := player(c.who)

	if region(c.who) == cloud_region ||
		region(c.who) == hades_region ||
		region(c.who) == faery_region {
		wout(c.who, "CLAIM may not be used in the Cloudlands, Hades or Faery.")
		return FALSE
	}

	if numargs(c) < 2 &&
		(kind(item) != T_item || has_item(pl, item) < 1) &&
		qty == 0 {
		wout(c.who, "(assuming you meant CLAIM %d %d)", item_gold, item)
		qty = item
		item = item_gold
	}

	if kind(item) != T_item {
		wout(c.who, "%s is not an item.", box_code(item))
		return FALSE
	}

	if has_item(pl, item) < 1 {
		wout(c.who, "No %s for you to claim.", box_code(item))
		return FALSE
	}

	qty = how_many(c.who, pl, item, qty, have_left)
	if qty <= 0 {
		return FALSE
	}

	ret := move_item(pl, c.who, item, qty)
	if !ret {
		panic("v_claim: move_item failed")
	}

	wout(c.who, "Claimed %s.", just_name_qty(item, qty))
	return TRUE
}

// count_man_items counts the number of man-type items held by who.
// Ported from src/u.c lines 1039-1056.
func count_man_items(who int) int {
	count := 0
	inv := teg.globals.inventories[who]
	for _, e := range inv {
		if man_item(e.item) != 0 {
			count += e.qty
		}
	}
	return count
}

// my_prisoner checks if pris is a prisoner of who.
// Simply checks if pris is a prisoner and their location is who.
// Ported from src/u.c lines 2336-2349.
func my_prisoner(who, pris int) bool {
	if kind(pris) != T_char {
		return false
	}

	if !is_prisoner(pris) {
		return false
	}

	if loc(pris) != who {
		return false
	}

	return true
}

// check_captain_loses_sailors handles sailor/pirate transfer affecting captaincy.
// Ported from src/move.c lines 1159-1185.
func check_captain_loses_sailors(qty, target, inform int) {
	// Stub: will be fully implemented in movement sprint
	// For now, just a placeholder
}
