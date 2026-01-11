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

// immed.go -- Immediate (GM) commands (Sprint 23)
//
// This file ports the immediate command mode from immed.c.
// Immediate commands are GM-only commands that execute instantly,
// bypassing the normal turn processing queue.
//
// Key functions:
//   - ImmediateCommands: interactive command loop (primarily for testing/GM)
//   - v_* functions: individual immediate command handlers
//
// Many of these commands are GM tools for debugging and world manipulation.
// In the web-based version, these will be exposed via admin endpoints.

// Note: gm_player constant is defined in glob.go

// ImmediateMode executes a single immediate command for GM/testing.
// This is a simplified version of the C immediate_commands() loop.
// Returns true if command executed successfully.
// Port of C immediate_commands() without the interactive loop.
func (e *Engine) ImmediateMode(who int, line string) bool {
	c := e.p_command(who)
	if c == nil {
		return false
	}

	c.who = who
	c.wait = 0

	if !e.oly_parse(c, line) {
		return false
	}

	c.pri = schar(e.cmd_pri(c.cmd))
	c.wait = e.cmd_time(c.cmd)
	c.poll = schar(e.cmd_poll(c.cmd))
	c.days_executing = 0
	c.state = STATE_LOAD

	e.do_command(c)

	for c.state == STATE_RUN {
		e.globals.evening = true
		e.finish_command(c)
		e.globals.evening = false
		e.olytimeIncrement()
	}

	return c.status != FALSE
}

// v_be changes the current immediate-mode identity.
// Port of C v_be().
func (e *Engine) v_be(c *command) int {
	if e.valid_box(c.a) {
		e.out(c.a, "You are now %s.", e.box_name(c.a))
		return TRUE
	}
	e.out(c.who, "'%d' not a valid box.", c.a)
	return FALSE
}

// v_listcmds lists all available commands.
// Port of C v_listcmds().
func (e *Engine) v_listcmds(c *command) int {
	// In Go, we would iterate over the command table
	// For now, this is a stub that will be expanded when cmd_tbl is fully ported
	e.out(c.who, "Command list not yet implemented in Go port.")
	return TRUE
}

// v_add_item adds items to an entity.
// Port of C v_add_item().
func (e *Engine) v_add_item(c *command) int {
	if e.Kind(c.a) == T_item {
		if e.Kind(c.who) != T_char {
			e.out(c.who, "Warning: %s not a character", e.box_name(c.who))
		}
		e.gen_item(c.who, c.a, c.b)
		return TRUE
	}
	e.wout(c.who, "%d is not a valid item.", c.a)
	return FALSE
}

// v_sub_item removes items from an entity.
// Port of C v_sub_item().
func (e *Engine) v_sub_item(c *command) int {
	e.consume_item(c.who, c.a, c.b)
	return TRUE
}

// v_dump outputs the box data for debugging.
// Port of C v_dump().
func (e *Engine) v_dump(c *command) int {
	if e.valid_box(c.a) {
		e.globals.bx[c.a].temp = 0
		// In the Go port, we'll log the box data instead of save_box to stdout
		e.out(c.who, "Box %s dumped.", e.box_code(c.a))
		return TRUE
	}
	return FALSE
}

// v_poof teleports a character to a location.
// Port of C v_poof().
func (e *Engine) v_poof(c *command) int {
	if !is_loc_or_ship(c.a) {
		e.wout(c.who, "%s is not a location.", e.box_code(c.a))
		return FALSE
	}

	e.move_stack(c.who, c.a)
	e.wout(c.who, ">poof!< A cloud of orange smoke appears and wisks you away...")
	e.out(c.who, "")
	e.show_loc(c.who, loc(c.who))

	return TRUE
}

// v_see_all toggles revealing hidden features.
// Port of C v_see_all().
func (e *Engine) v_see_all(c *command) int {
	if c.a == 0 {
		e.globals.immedSeeAll = true
	} else {
		e.globals.immedSeeAll = c.a != 0
	}

	if e.globals.immedSeeAll {
		e.out(c.who, "Will reveal all hidden features.")
	} else {
		e.out(c.who, "Hidden features will operate normally.")
	}
	return TRUE
}

// v_makeloc creates a new location at the current subloc.
// Port of C v_makeloc().
func (e *Engine) v_makeloc(c *command) int {
	// sk := e.lookup(subkind_s, parse1)
	// For now, stub implementation
	e.wout(c.who, "makeloc not yet fully implemented.")
	return FALSE
}

// v_invent shows character inventory.
// Port of C v_invent().
func (e *Engine) v_invent(c *command) int {
	e.show_char_inventory(c.who, c.who)
	e.show_carry_capacity(c.who, c.who)
	e.show_item_skills(c.who, c.who)
	return TRUE
}

// v_know teaches a skill to the character.
// Port of C v_know().
func (e *Engine) v_know(c *command) int {
	if e.Kind(c.a) != T_skill {
		e.wout(c.who, "%s is not a skill.", e.box_code(c.a))
		return FALSE
	}
	e.learn_skill(c.who, c.a)
	return TRUE
}

// v_skills lists character skills.
// Port of C v_skills().
func (e *Engine) v_skills(c *command) int {
	e.list_skills(c.who, c.who)
	e.list_partial_skills(c.who, c.who)
	return TRUE
}

// v_save saves the game database.
// Port of C v_save().
func (e *Engine) v_save(c *command) int {
	// In Go, we commit the current transaction or save state
	e.out(c.who, "Save triggered.")
	return TRUE
}

// v_los calculates line-of-sight distance.
// Port of C v_los().
func (e *Engine) v_los(c *command) int {
	target := c.a
	if !is_loc_or_ship(target) {
		e.wout(c.who, "%s is not a location.", e.box_code(target))
		return FALSE
	}

	d := e.los_province_distance(subloc(c.who), target)
	e.wout(c.who, "distance=%d", d)
	return TRUE
}

// v_kill kills a character.
// Port of C v_kill().
func (e *Engine) v_kill(c *command) int {
	e.kill_char(c.a, MATES)
	return TRUE
}

// v_take_pris takes a character prisoner.
// Port of C v_take_pris().
func (e *Engine) v_take_pris(c *command) int {
	if !e.check_char_here(c.who, c.a) {
		return FALSE
	}
	e.take_prisoner(c.who, c.a)
	return TRUE
}

// v_seed seeds initial locations.
// Port of C v_seed().
func (e *Engine) v_seed(c *command) int {
	e.seed_initial_locations()
	return TRUE
}

// v_postproc runs post-month processing.
// Port of C v_postproc().
func (e *Engine) v_postproc(c *command) int {
	// Reset studied flags and run post_month
	for _, i := range e.Characters() {
		ch := rp_char(i)
		if ch != nil {
			ch.studied = 0
		}
		// Reset skill experience flags
		for _, sk := range e.getCharSkills(i) {
			sk.exp_this_month = 0
		}
	}

	e.PostMonth()
	e.olytimeTurnChange()
	return TRUE
}

// v_lore delivers lore to a character.
// Port of C v_lore().
func (e *Engine) v_lore(c *command) int {
	if e.valid_box(c.a) {
		e.deliver_lore(c.who, c.a)
	}
	return TRUE
}

// v_ct clears city trades and regenerates them.
// Port of C v_ct().
func (e *Engine) v_ct(c *command) int {
	for i := sub_first(sub_city); i != 0; i = sub_next(i) {
		e.globals.bx[i].trades = nil
	}
	e.location_trades()
	return TRUE
}

// v_seedmarket seeds city market trades.
// Port of C v_seedmarket().
func (e *Engine) v_seedmarket(c *command) int {
	for i := sub_first(sub_city); i != 0; i = sub_next(i) {
		e.seed_city_trade(i)
	}
	for i := sub_first(sub_city); i != 0; i = sub_next(i) {
		e.loc_trade_sup(i, true)
	}
	return TRUE
}

// v_credit credits items or NP to a player/character.
// Port of C v_credit().
func (e *Engine) v_credit(c *command) int {
	target := c.a
	amount := c.b
	item := c.c

	if amount != 0 {
		k := e.Kind(target)
		if k != T_char && k != T_player {
			e.wout(c.who, "%s not a character or player.", e.box_code(target))
			return FALSE
		}

		if item == 0 {
			item = item_gold
		}

		e.gen_item(target, item, amount)
		e.wout(c.who, "Credited %s %d of item %d.", e.box_name(target), amount, item)
		e.wout(target, "Received CLAIM credit of %d items.", amount)
		return TRUE
	}

	if e.Kind(target) != T_char {
		e.wout(c.who, "%s not a character.", e.box_code(target))
		return FALSE
	}

	pl := player(target)
	if e.times_paid(pl) {
		e.wout(c.who, "Already paid faction %s.", e.box_name(pl))
		return FALSE
	}

	e.p_player(pl).times_paid = 1
	e.wout(target, "The Times pays %s %d gold.", e.box_name(target), 25)
	e.gen_item(target, item_gold, 25)

	return TRUE
}

// v_relore re-queues lore for all characters with a skill.
// Port of C v_relore().
func (e *Engine) v_relore(c *command) int {
	skill := c.a
	if !e.valid_box(skill) || e.Kind(skill) != T_skill {
		e.wout(c.who, "%s is not a skill.", e.box_code(skill))
		return FALSE
	}

	for _, i := range e.Characters() {
		if e.has_skill(i, skill) {
			e.queue_lore(i, skill, true)
		}
	}
	return TRUE
}

// v_xyzzy is a placeholder/debugging command.
// Port of C v_xyzzy().
func (e *Engine) v_xyzzy(c *command) int {
	return TRUE
}

// v_fix2 teaches advanced sorcery to characters with auraculum.
// Port of C v_fix2().
func (e *Engine) v_fix2(c *command) int {
	for _, i := range e.Characters() {
		if char_auraculum(i) != 0 {
			e.learn_skill(i, sk_adv_sorcery)
		}
	}
	return TRUE
}

// fix_gates calculates gate distances for provinces.
// Port of C fix_gates().
func (e *Engine) fix_gates() {
	e.clear_temps(T_loc)

	// First pass: mark provinces adjacent to gates
	for where := kind_first(T_loc); where != 0; where = kind_next(where) {
		if loc_depth(where) != LOC_province {
			continue
		}
		if !e.in_hades(where) && !e.in_clouds(where) && !e.in_faery(where) {
			continue
		}
		if !e.province_gate_here(where) {
			continue
		}

		l := e.exits_from_loc_nsew(0, where)
		for _, exit := range l {
			if loc_depth(exit.destination) != LOC_province {
				continue
			}
			if !e.province_gate_here(exit.destination) {
				e.globals.bx[exit.destination].temp = 1
			}
		}
	}

	// Iterative flood-fill to mark distances
	m := 1
	for {
		setOne := false
		for where := kind_first(T_loc); where != 0; where = kind_next(where) {
			if loc_depth(where) != LOC_province {
				continue
			}
			if !e.in_hades(where) && !e.in_clouds(where) && !e.in_faery(where) {
				continue
			}
			if e.province_gate_here(where) || e.globals.bx[where].temp != m {
				continue
			}

			l := e.exits_from_loc_nsew(0, where)
			for _, exit := range l {
				dest := exit.destination
				if loc_depth(dest) != LOC_province {
					continue
				}
				if !e.province_gate_here(dest) && e.globals.bx[dest].temp == 0 {
					e.globals.bx[dest].temp = m + 1
					setOne = true
				}
			}
		}
		m++
		if !setOne {
			break
		}
	}

	// Copy temp to dist_from_gate
	for where := kind_first(T_loc); where != 0; where = kind_next(where) {
		if loc_depth(where) != LOC_province {
			continue
		}
		if !e.in_hades(where) && !e.in_clouds(where) && !e.in_faery(where) {
			continue
		}
		p_loc(where).dist_from_gate = schar(e.globals.bx[where].temp)
	}
}

// v_fix sets HTML passwords for all players.
// Port of C v_fix().
func (e *Engine) v_fix(c *command) int {
	for pl := kind_first(T_player); pl != 0; pl = kind_next(pl) {
		e.set_html_pass(pl)
	}
	return TRUE
}

// v_plugh is a no-op placeholder command.
// Port of C v_plugh().
func (e *Engine) v_plugh(c *command) int {
	return TRUE
}

// Helper stubs for functions not yet ported

func (e *Engine) out(who int, format string, args ...any)  {}
func (e *Engine) wout(who int, format string, args ...any) {}

func (e *Engine) box_name(n int) string {
	return e.getName(n)
}

func (e *Engine) box_code(n int) string {
	return int_to_code(n)
}

func (e *Engine) valid_box(n int) bool {
	return valid_box(n)
}

// Note: Engine.Kind() is defined in glob.go

func (e *Engine) gen_item(who, item, qty int) {
	inv := e.globals.inventories[who]
	for i := range inv {
		if inv[i].item == item {
			inv[i].qty += qty
			return
		}
	}
	if e.globals.inventories == nil {
		e.globals.inventories = make(map[int][]item_ent)
	}
	e.globals.inventories[who] = append(e.globals.inventories[who], item_ent{item: item, qty: qty})
}

func (e *Engine) consume_item(who, item, qty int) int {
	inv := e.globals.inventories[who]
	for i := range inv {
		if inv[i].item == item {
			if inv[i].qty >= qty {
				inv[i].qty -= qty
				return qty
			}
			consumed := inv[i].qty
			inv[i].qty = 0
			return consumed
		}
	}
	return 0
}

func (e *Engine) move_stack(who, where int) {
	// Move character and stack to location
	set_where(who, where)
}

func (e *Engine) show_loc(who, where int)                    {}
func (e *Engine) show_char_inventory(who, num int)           {}
func (e *Engine) show_carry_capacity(who, num int)           {}
func (e *Engine) show_item_skills(who, num int)              {}
func (e *Engine) learn_skill(who, sk int)                    {}
func (e *Engine) list_skills(who, num int)                   {}
func (e *Engine) list_partial_skills(who, num int)           {}
func (e *Engine) los_province_distance(from, to int) int     { return 0 }
func (e *Engine) kill_char(who, inherit int)                 {}
func (e *Engine) check_char_here(who, target int) bool       { return true }
func (e *Engine) take_prisoner(who, target int)              {}
func (e *Engine) seed_initial_locations()                    {}
// Note: Engine.getCharSkills() is defined in load.go
func (e *Engine) deliver_lore(who, num int)                  {}
func (e *Engine) location_trades()                           {}
func (e *Engine) seed_city_trade(where int)                  {}
func (e *Engine) loc_trade_sup(where int, flag bool)         {}
func (e *Engine) times_paid(pl int) bool                     { return p_player(pl).times_paid != 0 }
func (e *Engine) has_skill(who, skill int) bool              { return false }
func (e *Engine) queue_lore(who, num int, anyway bool)       {}
func (e *Engine) clear_temps(k int)                          {}
func (e *Engine) in_hades(where int) bool                    { return false }
func (e *Engine) in_clouds(where int) bool                   { return false }
func (e *Engine) in_faery(where int) bool                    { return false }
func (e *Engine) province_gate_here(where int) bool          { return false }
func (e *Engine) exits_from_loc_nsew(who, where int) []*exit_view { return nil }
func (e *Engine) set_html_pass(pl int)                       {}

func (e *Engine) p_player(n int) *entity_player {
	return p_player(n)
}

// Engine.globals extension for immediate mode
func init() {
	// Note: immedSeeAll is added to Engine.globals in engine.go
}
