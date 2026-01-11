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

// lifecycle.go - Character death & revival ported from src/u.c (Sprint 26.1)

package taygete

// survive_fatal checks if a character survives a fatal wound via the sk_survive_fatal skill.
// If so, the skill is forgotten, health is restored to 100, and true is returned.
// Ported from src/u.c lines 50-79.
func survive_fatal(who int) bool {
	if !has_skill_check(who, sk_survive_fatal) {
		return false
	}

	if forget_skill(who, sk_survive_fatal) {
		wout(who, "%s would have died, but survived a fatal wound!", box_name(who))
		wout(who, "Forgot %s.", box_code(sk_survive_fatal))
		wout(who, "Health is now 100.")

		p_char(who).health = 100
		p_char(who).sick = FALSE

		return true
	}

	if is_magician(who) != 0 {
		p_magic(who).cur_aura = 0
	}

	return false
}

// Note: char_reclaim is defined in stack.go and updated to use kill_char.

// kill_stack_ocean kills all characters in a stack that sinks at sea.
// Characters with survive_fatal survive and swim to nearest land.
// Ported from src/u.c lines 14-47.
func kill_stack_ocean(who int) {
	var l []int

	for _, i := range loop_stack_list(who) {
		l = append(l, i)
	}

	for i := len(l) - 1; i >= 0; i-- {
		kill_char(l[i], 0)

		if kind(l[i]) == T_char { // not dead yet (survived)
			extract_stacked_unit(l[i])
			where := find_nearest_land(province(l[i]))

			out(l[i], "%s washed ashore at %s.", box_name(l[i]), box_name(where))
			log_write(LOG_SPECIAL, "kill_stack_ocean, swam ashore, who=%s", box_code_less(l[i]))
			move_stack(l[i], where)
		}
	}
}

// stackmate_inheritor finds who should inherit items from a dying character.
// First tries to find someone stacked below, then the stack parent.
// Ported from src/u.c lines 227-247.
func stackmate_inheritor(who int) int {
	target := 0

	p := rp_loc_info(who)
	if p != nil {
		for _, i := range p.here_list {
			if kind(i) == T_char && !is_prisoner(i) {
				target = i
				break
			}
		}
	}

	if target != 0 {
		return target
	}

	return stack_parent(who)
}

// take_unit_items transfers inventory from a dying unit.
// inherit specifies who receives items:
//
//	0 = discard items (no inheritor)
//	MATES = stackmate inheritor
//	MATES_SILENT = stackmate inheritor (no messages)
//	other = specific target
//
// how_many specifies what to transfer:
//
//	TAKE_ALL = all items
//	TAKE_SOME = random subset
//	TAKE_NI = add noble item first, then all
//
// Also handles prisoner transfer.
// Ported from src/u.c lines 250-366.
func take_unit_items(from, inherit, how_many int) {
	var to int
	var silent bool

	switch inherit {
	case 0:
		to = 0
		silent = true
	case MATES:
		to = stackmate_inheritor(from)
		silent = false
	case MATES_SILENT:
		to = stackmate_inheritor(from)
		silent = true
	default:
		to = inherit
		silent = false
	}

	if how_many == TAKE_NI {
		gen_item(from, int(noble_item(from)), 1)
	}

	first := true
	inv := teg.globals.inventories[from]
	for _, e := range inv {
		qty := e.qty

		if e.qty > 0 && how_many == TAKE_SOME && rnd(1, 2) == 1 {
			qty = rnd(0, e.qty)
		}

		if qty > 0 && !silent && valid_box(to) {
			if first {
				first = false
				wout(to, "Taken from %s:", box_name(from))
			}
			wout(to, "   %s", box_name_qty(e.item, qty))
		}

		move_item(from, to, e.item, qty)

		if e.item == item_gold && player(from) != player(to) {
			if player(from) == indep_player {
				gold_combat_indep += qty
			} else {
				gold_combat += qty
			}
		}

		if qty != e.qty {
			move_item(from, 0, e.item, e.qty-qty)
		}
	}

	p := rp_loc_info(from)
	if p != nil {
		for _, i := range p.here_list {
			if kind(i) == T_char && is_prisoner(i) {
				if to > 0 {
					if first && !silent {
						wout(to, "Taken from %s:", box_name(from))
						first = false
					}

					move_prisoner(from, to, i)

					if !silent {
						wout(to, "   %s", liner_desc(i))
					}

					if player(i) == player(to) {
						p_char(i).prisoner = FALSE
					}
				} else {
					p_magic(i).swear_on_release = FALSE
					drop_stack(from, i)
				}
			}
		}
	}
}

// add_char_damage adds damage to a character, potentially killing them.
// If health reaches 0, kill_char is called.
// If health doesn't reach 0 but falls below a random threshold, the character becomes sick.
// Ported from src/u.c lines 369-407.
func add_char_damage(who, amount, inherit int) {
	if amount <= 0 {
		return
	}

	p := p_char(who)

	if p.health == -1 {
		if amount >= 50 {
			kill_char(who, inherit)
		}
		return
	}

	if p.health > 0 {
		if amount > int(p.health) {
			amount = int(p.health)
		}

		p.health -= schar(amount)

		wout(who, "%s is wounded. Health is now %d.", box_name(who), p.health)
	}

	if p.health <= 0 {
		kill_char(who, inherit)
	} else if p.sick == 0 && rnd(1, 100) > int(p.health) {
		p.sick = TRUE
		wout(who, "%s has fallen ill.", box_name(who))
	}
}

// put_back_cookie returns an NPC's allocation cookie to its home location.
// Used when an NPC dies to allow respawn.
// Ported from src/u.c lines 410-421.
func put_back_cookie(who int) {
	p := rp_misc(who)

	if p == nil || p.npc_home == 0 {
		return
	}

	gen_item(p.npc_home, p.npc_cookie, 1)
}

// dead_char_body converts a character into a dead body item.
// NPCs and characters lost at sea don't leave bodies.
// Ported from src/u.c lines 456-495.
func dead_char_body(pl, who int) {
	grave := province(who)
	if subkind(grave) == sub_ocean {
		grave = find_nearest_land(grave)
	}

	set_where(who, 0)

	if char_melt_me(who) != 0 || is_npc(who) || grave == 0 {
		change_box_kind(who, T_deadchar)

		if subkind(who) != 0 {
			change_box_subkind(who, 0)
		}
		return
	}

	change_box_kind(who, T_item)
	change_box_subkind(who, sub_dead_body)

	// Store original name in savedNames map (replaces entity_misc.save_name)
	savedNames[who] = teg.getName(who)
	teg.setName(who, "dead body")

	pm := p_misc(who)
	pm.old_lord = pl

	pi := p_item(who)
	pi.weight = item_weight(item_peasant)
	teg.setPluralName(who, "dead bodies")

	hack_unique_item(who, grave)
}

// restore_dead_body revives a character from a dead body.
// Used by resurrection magic.
// Ported from src/u.c lines 521-591.
func restore_dead_body(owner, who int) {
	log_write(LOG_CODE, "dead body revived: who=%s, owner=%s, player=%s",
		box_code_less(who),
		box_code_less(owner),
		box_code_less(player(owner)))

	if !sub_item(owner, who, 1) {
		panic("restore_dead_body: sub_item failed")
	}

	p_item(who).who_has = 0

	change_box_kind(who, T_char)
	change_box_subkind(who, 0)

	pm := p_misc(who)
	pi := p_item(who)
	pc := p_char(who)

	pi.weight = 0
	teg.setPluralName(who, "")

	// Restore original name from savedNames map
	savedName := savedNames[who]
	if savedName != "" {
		teg.setName(who, savedName)
		delete(savedNames, who)
	}

	pc.health = 100
	pc.sick = FALSE

	set_where(who, subloc(owner))

	if kind(pm.old_lord) == T_player {
		wout(pm.old_lord, "%s has been brought back to life.", box_name(who))
		set_lord(who, pm.old_lord, LOY_UNCHANGED, 0)
	} else {
		set_lord(who, indep_player, LOY_UNCHANGED, 0)
	}

	pm.old_lord = 0

	def := char_defense(who)
	def -= 50
	if def < 0 {
		def = 0
	}
	p_char(who).defense = short(def)
}

// kill_char is the main death pipeline for characters.
// It handles:
// - survive_fatal check
// - Message output
// - Item transfer to inheritor
// - Stack extraction
// - Order flushing
// - Aura zeroing for mages
// - Transcend death (move to Hades)
// - Token handling
// - Unit desertion
// - Dead body creation
//
// Ported from src/u.c lines 594-693.
func kill_char(who, inherit int) {
	where := subloc(who)
	_ = where // used for logging in full implementation
	pl := player(who)

	if kind(who) != T_char {
		return // Not a character, nothing to do
	}

	if char_melt_me(who) == 0 && survive_fatal(who) {
		return
	}

	p_char(who).prisoner = FALSE
	msg := "died"
	if char_melt_me(who) != 0 {
		msg = "vanished"
	}
	wout(who, "*** %s has %s ***", just_name(who), msg)

	sp := stack_parent(who)
	if sp != 0 {
		wout(sp, "%s has %s.", box_name(who), msg)
	}

	p_char(who).prisoner = TRUE // suppress output during cleanup

	logMsg := "died"
	if char_melt_me(who) != 0 {
		logMsg = "melted"
	}
	log_write(LOG_DEATH, "%s %s in %s.", box_name(who), logMsg, char_rep_location(who))

	take_unit_items(who, inherit, TAKE_SOME)

	extract_stacked_unit(who)

	flush_unit_orders(player(who), who)
	interrupt_order(who)

	if is_magician(who) != 0 {
		p_magic(who).cur_aura = 0
	}

	if char_melt_me(who) == 0 && has_skill_check(who, sk_transcend_death) {
		hades_point := random_hades_loc()

		p_char(who).prisoner = FALSE

		log_write(LOG_SPECIAL, "%s transcends death", box_name(who))
		log_write(LOG_SPECIAL, "...%s moved to %s", box_name(who), box_name(hades_point))
		p_char(who).prisoner = FALSE
		p_char(who).sick = FALSE
		p_char(who).health = 100
		move_stack(who, hades_point)
		wout(who, "%s appears at %s.", box_name(who), box_name(hades_point))

		return
	}

	token_item := our_token(who)
	if token_item != 0 {
		who_has := item_unique(token_item)
		token_pl := p_player(token_item)

		if token_pl != nil {
			// Remove from player units - uses Engine helper since units is stored in globals
			removePlayerUnit(token_item, who)
		}

		if char_melt_me(who) == 0 {
			im := rp_item_magic(token_item)
			if im != nil {
				im.token_num--
			}
		}

		if item_token_num(token_item) <= 0 {
			if player(who_has) == sub_pl_regular {
				wout(who_has, "%s vanishes.", box_name(token_item))
			}
			destroy_unique_item(who_has, token_item)
		}
	}

	unit_deserts(who, 0, true, LOY_UNCHANGED, 0)

	put_back_cookie(who)
	p_char(who).death_time = teg.globals.sysclock
	p_char(who).prisoner = FALSE
	dead_char_body(pl, who)
}

// Global combat gold tracking (from u.c)
var gold_combat int       // gold taken from players in combat
var gold_combat_indep int // gold taken from independents in combat

// savedNames stores the original names for dead bodies (replaces entity_misc.save_name)
// Maps entity ID -> saved name
var savedNames = make(map[int]string)

// Helper functions and stubs for dependencies

// has_skill_check checks if a character has the given skill.
// This wraps the Engine method for package-level use.
func has_skill_check(who, skill int) bool {
	skills := teg.getCharSkills(who)
	for _, s := range skills {
		if s.skill == skill && s.know == SKILL_know {
			return true
		}
	}
	return false
}

// forget_skill removes a skill from a character.
// Returns true if the skill was known and removed.
// Ported from src/use.c lines 1008-1027.
func forget_skill(who, skill int) bool {
	skills := teg.getCharSkills(who)
	for i, s := range skills {
		if s.skill == skill {
			if s.know == SKILL_know {
				skills[i].know = SKILL_dont
				return true
			}
			return false
		}
	}
	return false
}

// char_melt_me returns 1 if the character is marked for melting.
func char_melt_me(n int) schar {
	p := rp_char(n)
	if p == nil {
		return 0
	}
	return p.melt_me
}

// is_npc returns true if the entity is an NPC.
// #define is_npc(n) (subkind(n) || loyal_kind(n) == LOY_npc || loyal_kind(n) == LOY_summon)
func is_npc(n int) bool {
	if subkind(n) != 0 {
		return true
	}
	lk := loyal_kind(n)
	return lk == LOY_npc || lk == LOY_summon
}

// Note: our_token is defined in accessor.go
// Note: item_token_num is defined in accessor.go

// loop_stack_list returns all members of a stack (including nested stacks).
func loop_stack_list(who int) []int {
	var result []int
	result = append(result, who)
	p := rp_loc_info(who)
	if p != nil {
		for _, i := range p.here_list {
			if kind(i) == T_char {
				result = append(result, loop_stack_list(i)...)
			}
		}
	}
	return result
}

// removePlayerUnit removes a unit from a player's unit list.
// Uses the Engine's playerUnits map instead of entity_player.units.
func removePlayerUnit(pl, who int) {
	if teg.globals.playerUnits == nil {
		return
	}
	units := teg.globals.playerUnits[pl]
	for i, u := range units {
		if u == who {
			teg.globals.playerUnits[pl] = append(units[:i], units[i+1:]...)
			return
		}
	}
}

// Note: move_prisoner is defined in stack.go
// Note: drop_stack is defined in stack.go

// liner_desc returns a one-line description of an entity.
// Stub: returns box_name for now.
func liner_desc(n int) string {
	return box_name(n)
}

// char_rep_location returns a printable location for a character.
// Stub: returns box_name of subloc.
func char_rep_location(who int) string {
	return box_name(subloc(who))
}

// flush_unit_orders removes all pending orders for a unit.
// Stub: will be implemented with order system.
func flush_unit_orders(pl, who int) {
	if teg.globals.orderQueues == nil {
		return
	}
	if teg.globals.orderQueues[pl] == nil {
		return
	}
	delete(teg.globals.orderQueues[pl], who)
}

// interrupt_order interrupts the current order for a unit.
// Stub: will be implemented with command execution system.
func interrupt_order(who int) {
	if teg.globals.bx[who] == nil {
		return
	}
	if teg.globals.bx[who].cmd != nil {
		teg.globals.bx[who].cmd.state = STATE_DONE
		teg.globals.bx[who].cmd.status = FALSE
	}
}

// random_hades_loc returns a random location in Hades.
// Stub: returns 0 for now; will be implemented with Hades system.
func random_hades_loc() int {
	return 0
}

// Note: unit_deserts is defined in stack.go and updated to call kill_char
// Note: change_box_kind is defined in code.go
// Note: change_box_subkind is defined in code.go
// Note: char_defense is defined in accessor.go
// Note: set_lord is defined in stack.go

// Note: IListRemValue is defined in z.go
// Note: charPtrToStr is defined in cmd_economy.go
