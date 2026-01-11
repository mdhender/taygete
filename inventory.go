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

// inventory.go - Inventory & gold helpers ported from src/u.c (Sprint 25.1)

package taygete

// has_item returns the quantity of item held by who.
// Returns 0 if who doesn't have any of the item.
// Ported from src/u.c lines 1677-1698.
func has_item(who, item int) int {
	if !valid_box(who) {
		return 0
	}
	if !valid_box(item) {
		return 0
	}

	inv := teg.globals.inventories[who]
	for _, e := range inv {
		if e.item == item {
			return e.qty
		}
	}
	return 0
}

// add_item adds qty of item to who's inventory.
// Also handles lore delivery for new items.
// Ported from src/u.c lines 1702-1742.
func add_item(who, item, qty int) {
	if !valid_box(who) {
		panic("add_item: invalid who")
	}
	if !valid_box(item) {
		panic("add_item: invalid item")
	}
	if qty < 0 {
		panic("add_item: qty < 0")
	}

	lore := item_lore(item)
	if lore != 0 && kind(who) == T_char && !test_known(who, item) {
		queue_lore(who, item, false)
	}

	if teg.globals.inventories == nil {
		teg.globals.inventories = make(map[int][]item_ent)
	}

	inv := teg.globals.inventories[who]
	for i := range inv {
		if inv[i].item == item {
			inv[i].qty += qty
			investigate_possible_trade(who, item, inv[i].qty-qty)
			return
		}
	}

	teg.globals.inventories[who] = append(teg.globals.inventories[who], item_ent{item: item, qty: qty})
	investigate_possible_trade(who, item, 0)
}

// sub_item subtracts qty of item from who's inventory.
// Returns true if successful, false if who doesn't have enough.
// Ported from src/u.c lines 498-518.
func sub_item(who, item, qty int) bool {
	if !valid_box(who) {
		return false
	}
	if !valid_box(item) {
		return false
	}
	if qty < 0 {
		return false
	}

	inv := teg.globals.inventories[who]
	for i := range inv {
		if inv[i].item == item {
			if inv[i].qty < qty {
				return false
			}
			inv[i].qty -= qty
			return true
		}
	}
	return false
}

// gen_item generates qty of a non-unique item for who.
// Panics if item is unique.
// Ported from src/u.c lines 1745-1751.
func gen_item(who, item, qty int) {
	if item_unique(item) != 0 {
		panic("gen_item: item is unique")
	}
	add_item(who, item, qty)
}

// consume_item consumes qty of a non-unique item from who.
// Panics if item is unique.
// Returns true if successful, false if who doesn't have enough.
// Ported from src/u.c lines 1754-1760.
func consume_item(who, item, qty int) bool {
	if item_unique(item) != 0 {
		panic("consume_item: item is unique")
	}
	return sub_item(who, item, qty)
}

// move_item moves qty of item from one unit to another.
// If to is 0, the items are discarded (via drop_item).
// Returns true if successful, false if from doesn't have enough.
// Ported from src/u.c lines 1768-1795.
func move_item(from, to, item, qty int) bool {
	if qty <= 0 {
		return true
	}

	if to == 0 {
		return drop_item(from, item, qty)
	}

	if !sub_item(from, item, qty) {
		return false
	}

	add_item(to, item, qty)

	if item_unique(item) != 0 {
		if qty != 1 {
			panic("move_item: unique item qty != 1")
		}
		p_item(item).who_has = to

		if subkind(item) == sub_npc_token {
			move_token(item, from, to)
		}
	}

	return true
}

// hack_unique_item sets the owner of a unique item without transfer checks.
// Used for initial setup and special cases.
// Ported from src/u.c lines 1798-1804.
func hack_unique_item(item, owner int) {
	p_item(item).who_has = owner
	add_item(owner, item, 1)
}

// create_unique_item creates a new unique item and gives it to who.
// Returns the new item ID, or -1 on failure.
// Ported from src/u.c lines 1807-1821.
func create_unique_item(who int, sk schar) int {
	newItem := new_ent(T_item, sk)
	if newItem < 0 {
		return -1
	}

	p_item(newItem).who_has = who
	add_item(who, newItem, 1)

	return newItem
}

// create_unique_item_alloc creates a unique item with a specific pre-allocated ID.
// Returns the item ID.
// Ported from src/u.c lines 1824-1834.
func create_unique_item_alloc(newItem, who int, sk schar) int {
	alloc_box(newItem, T_item, sk)

	p_item(newItem).who_has = who
	add_item(who, newItem, 1)

	return newItem
}

// destroy_unique_item destroys a unique item.
// If the item is a dead body, grants NP to the original lord.
// Ported from src/u.c lines 1837-1867.
func destroy_unique_item(who, item int) {
	if kind(item) != T_item {
		panic("destroy_unique_item: item is not T_item")
	}
	if item_unique(item) == 0 {
		panic("destroy_unique_item: item is not unique")
	}

	if subkind(item) == sub_dead_body {
		pl := body_old_lord(item)

		if pl == indep_player {
			pl = char_prev_lord(item)
		}

		if kind(pl) == T_player {
			nps := char_np_total(item)
			out(pl, "%s~%s has passed on.  Gained %d NP%s.",
				save_name(item), box_code(item),
				nps, add_s(nps))
			add_np(pl, nps)
		}
	}

	if !sub_item(who, item, 1) {
		panic("destroy_unique_item: sub_item failed")
	}

	delete_box(item)
}

// drop_item drops an item from who's inventory.
// Non-unique items are simply consumed.
// Unique items are moved to the province (or nearest land if at sea).
// Ported from src/u.c lines 1962-1992.
func drop_item(who, item, qty int) bool {
	if item_unique(item) == 0 {
		return consume_item(who, item, qty)
	}

	whoGets := province(who)

	if subkind(item) == sub_dead_body {
		if whoGets == 0 {
			destroy_unique_item(who, item)
			return true
		}
	}

	if subkind(whoGets) == sub_ocean {
		whoGets = find_nearest_land(whoGets)
	}

	if whoGets == 0 {
		whoGets = province(who)
	}

	log_write(LOG_CODE, "drop_item: %s from %s to %s",
		box_name(item), box_name(subloc(who)), box_name(whoGets))

	return move_item(who, whoGets, item, qty)
}

// can_pay returns true if who has at least amount gold.
// Ported from src/u.c lines 1995-2000.
func can_pay(who, amount int) bool {
	return has_item(who, item_gold) >= amount
}

// charge deducts amount gold from who.
// Returns true if successful, false if who doesn't have enough.
// Ported from src/u.c lines 2003-2007.
func charge(who, amount int) bool {
	return sub_item(who, item_gold, amount)
}

// stack_has_item returns the total quantity of item held by who's stack.
// Only counts units owned by the same player as who.
// Ported from src/u.c lines 2010-2027.
// Note: This replaces the stub in stack.go.
func stack_has_item(who, item int) int {
	head := stack_leader(who)
	sum := 0

	for _, i := range stackMembers(head) {
		if player(i) != player(who) {
			continue
		}
		sum += has_item(i, item)
	}

	return sum
}

// has_use_key returns the first item in who's inventory with the given use key.
// Returns 0 if no such item is found.
// Ported from src/u.c lines 2030-2050.
func has_use_key(who, key int) int {
	inv := teg.globals.inventories[who]
	for _, e := range inv {
		p := rp_item_magic(e.item)
		if p != nil && int(p.use_key) == key {
			return e.item
		}
	}
	return 0
}

// stack_has_use_key returns the first item with the given use key in who's stack.
// Only checks units owned by the same player as who.
// Returns 0 if no such item is found.
// Ported from src/u.c lines 2053-2072.
func stack_has_use_key(who, key int) int {
	head := stack_leader(who)

	for _, i := range stackMembers(head) {
		if player(i) != player(who) {
			continue
		}
		ret := has_use_key(i, key)
		if ret != 0 {
			return ret
		}
	}

	return 0
}

// stack_sub_item subtracts qty of item from who's stack.
// First tries who, then borrows from friendly stackmates.
// Returns true if successful, false if the stack doesn't have enough.
// Ported from src/u.c lines 2083-2146.
func stack_sub_item(who, item, qty int) bool {
	if stack_has_item(who, item) < qty {
		return false
	}

	n := min(has_item(who, item), qty)
	if n > 0 {
		qty -= n
		sub_item(who, item, n)
	}

	if qty == 0 {
		return true
	}

	head := stack_leader(who)

	for _, i := range stackMembers(head) {
		if qty <= 0 {
			break
		}

		if player(i) != player(who) {
			continue
		}

		n = min(has_item(i, item), qty)
		if n > 0 {
			qty -= n
			sub_item(i, item, n)
		}
	}

	if qty != 0 {
		panic("stack_sub_item: qty != 0 after borrowing")
	}

	return true
}

// autocharge charges amount gold from who's stack.
// Returns true if successful, false if the stack doesn't have enough.
// Ported from src/u.c lines 2149-2154.
func autocharge(who, amount int) bool {
	return stack_sub_item(who, item_gold, amount)
}

// stackMembers returns all characters stacked under leader (including leader).
// This is a helper for iterating over stack members in Go style.
func stackMembers(leader int) []int {
	var result []int
	result = append(result, leader)

	p := rp_loc_info(leader)
	if p == nil {
		return result
	}

	for _, id := range p.here_list {
		if kind(id) == T_char && !is_prisoner(id) {
			result = append(result, stackMembers(id)...)
		}
	}

	return result
}

// add_np adds noble points to a player.
// Ported from src/u.c lines 1619-1629.
func add_np(pl, num int) {
	if kind(pl) != T_player {
		panic("add_np: not a player")
	}
	p := p_player(pl)
	p.noble_points += short(num)
	p.np_gained += short(num)
}

// deduct_np deducts noble points from a player.
// Returns true if successful, false if the player doesn't have enough.
// Ported from src/u.c lines 1601-1616.
func deduct_np(pl, num int) bool {
	if kind(pl) != T_player {
		panic("deduct_np: not a player")
	}
	p := p_player(pl)
	if int(p.noble_points) < num {
		return false
	}
	p.noble_points -= short(num)
	p.np_spent += short(num)
	return true
}

// deduct_aura deducts aura from a character.
// Returns true if successful, false if the character doesn't have enough aura.
// Ported from src/u.c lines 1632-1644.
func deduct_aura(who, amount int) bool {
	p := rp_magic(who)
	if p == nil || p.cur_aura < amount {
		return false
	}
	p.cur_aura -= amount
	return true
}

// charge_aura deducts aura from a character and reports failure.
// Returns true if successful, false if insufficient aura (with message).
// Ported from src/u.c lines 1647-1659.
func charge_aura(who, amount int) bool {
	if !deduct_aura(who, amount) {
		wout(who, "%s aura required, current level is %s.",
			cap(nice_num(amount)), nice_num(char_cur_aura(who)))
		return false
	}
	return true
}

// check_aura checks if a character has enough aura (without deducting).
// Returns true if sufficient, false otherwise (with message).
// Ported from src/u.c lines 1662-1674.
func check_aura(who, amount int) bool {
	if char_cur_aura(who) < amount {
		wout(who, "%s aura required, current level is %s.",
			cap(nice_num(amount)), nice_num(char_cur_aura(who)))
		return false
	}
	return true
}

// Stub functions for dependencies not yet implemented.
// These will be implemented in later sprints.

// queue_lore queues lore delivery for an item.
// TODO: Implement in later sprint.
func queue_lore(who, item int, anyway bool) {
}

// investigate_possible_trade checks if a trade can be made.
// TODO: Implement in later sprint (trade system).
func investigate_possible_trade(who, item, oldQty int) {
}

// move_token handles NPC token movement.
// TODO: Implement in later sprint (NPC system).
func move_token(token, from, to int) {
}

// char_prev_lord returns the previous lord of a character.
func char_prev_lord(n int) int {
	c := rp_char(n)
	if c == nil {
		return 0
	}
	return c.prev_lord
}

// save_name returns the saved name for a dead body.
// Note: save_name is stored as *char (legacy C type), currently returns empty.
// TODO: Implement proper string storage for save_name when needed.
func save_name(n int) string {
	return ""
}

// char_np_total returns the total NP cost for a character.
// TODO: Implement properly (uses skill system).
func char_np_total(who int) int {
	return 1
}

// add_s returns "s" for plural or "" for singular.
func add_s(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// out outputs a message to a player or character.
// TODO: Implement in later sprint (output system).
func out(who int, format string, args ...any) {
}

// box_code and box_name are defined in code.go

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the larger of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
