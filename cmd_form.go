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

// cmd_form.go - Noble formation commands ported from src/c1.c
// Sprint 26.6: Noble Formation

package taygete

// noble_cost returns the NP cost to form a new noble for the given player.
// Currently fixed at 1 NP per noble.
// Ported from src/c1.c line 746: #define noble_cost(pl) (1)
func noble_cost(pl int) int {
	return 1
}

// next_np_turn calculates the number of turns until the next NP grant.
// All factions gain noble points on turns 8, 16, 24, etc (every NUM_MONTHS turns).
// Returns 0 if NP is granted this turn.
// Ported from src/c1.c lines 750-769.
func next_np_turn(pl int) int {
	// ct = (7 - (sysclock.turn + 1) % NUM_MONTHS)
	ct := (7 - int(teg.globals.sysclock.turn+1)%NUM_MONTHS)

	// by Cappinator: All factions gain noble points on turns 8, 16, 24, etc...
	// ft = p->first_turn % NUM_MONTHS -> simplified to 2
	ft := 2

	n := (ft + ct) % NUM_MONTHS

	return n
}

// print_hiring_status displays when the player will receive their next NP.
// Ported from src/c1.c lines 773-789.
func print_hiring_status(pl int) {
	if kind(pl) != T_player {
		return
	}

	if subkind(pl) != sub_pl_regular {
		return
	}

	n := next_np_turn(pl)

	if n == 0 {
		n += NUM_MONTHS
	}

	wout(pl, "The next NP will be received at the end of turn %d.",
		int(teg.globals.sysclock.turn)+n)
}

// print_unformed displays the list of unformed noble IDs for a player.
// Ported from src/c1.c lines 793-812.
func print_unformed(pl int) {
	p := rp_player(pl)
	if p == nil {
		return
	}

	unformed := getPlayerUnformed(pl)
	n := len(unformed)
	if n < 1 {
		return
	}

	// Show at most 5 unformed nobles
	var codes []string
	for i := 0; i < n && i < 5; i++ {
		codes = append(codes, box_code_less(unformed[i]))
	}

	out(pl, "")
	wout(pl, "The next %s nobles formed will be: %s", nice_num(n), joinStrings(codes, " "))
	out(pl, "")
}

// form_new_noble creates a new noble from an unformed entity.
// Sets up the noble's stats, location, loyalty, and joins it to the forming character's stack.
// Ported from src/c1.c lines 903-933.
func form_new_noble(who int, name string, new int) {
	if kind(new) != T_unform {
		return
	}

	change_box_kind(new, T_char)

	p := p_char(new)
	op := p_char(who)

	p.behind = op.behind
	p.fresh_hire = TRUE
	p.health = 100

	p.attack = 80
	p.defense = 80
	p.break_point = 50

	set_name(new, name)

	set_where(new, subloc(who))
	set_lord(new, player(who), LOY_contract, 500)

	join_stack(new, who)
}

// v_form is the start routine for the FORM command.
// Validates that the character is in a city and has enough NP.
// Usage: FORM [entity_id] [name]
// Ported from src/c1.c lines 937-959.
func v_form(c *command) int {
	if subkind(subloc(c.who)) != sub_city {
		wout(c.who, "Nobles may only be formed in cities.")
		return FALSE
	}

	pl := player(c.who)
	cost := noble_cost(pl)

	if int(player_np(pl)) < cost {
		wout(c.who, "To form another noble requires %d Noble Point%s.", cost, add_s(cost))
		return FALSE
	}

	return TRUE
}

// d_form is the finish routine for the FORM command.
// Deducts NP and creates the new noble.
// Ported from src/c1.c lines 963-1026.
func d_form(c *command) int {
	pl := player(c.who)
	cost := noble_cost(pl)

	if int(player_np(pl)) < cost {
		wout(c.who, "To form another noble requires %d Noble Point%s.",
			cost, add_s(cost))
		return FALSE
	}

	unformed := getPlayerUnformed(pl)
	new := c.a

	if new != 0 {
		if kind(new) != T_unform || !containsInt(unformed, new) {
			wout(c.who, "%s is not a valid unformed noble entity.", box_code(new))
			new = 0
		}
	}

	if new == 0 && len(unformed) > 0 {
		new = unformed[0]
	}

	if new == 0 {
		wout(c.who, "No further nobles may be formed this turn.")
		return FALSE
	}

	new_name := "New noble"
	if numargs(c) >= 2 {
		parsed := get_parse_arg(c, 2)
		if parsed != "" {
			new_name = parsed
		}
	}

	form_new_noble(c.who, new_name, new)

	removePlayerUnformed(pl, new)
	deduct_np(pl, cost)

	return TRUE
}

// getPlayerUnformed returns the unformed noble IDs for a player.
// Uses a workaround since ilist is a C-style pointer type.
func getPlayerUnformed(pl int) []int {
	if kind(pl) != T_player {
		return nil
	}
	if teg.globals.playerUnits == nil {
		return nil
	}
	return teg.globals.playerUnits[pl+100_000]
}

// setPlayerUnformed sets the unformed noble IDs for a player.
func setPlayerUnformed(pl int, unformed []int) {
	if teg.globals.playerUnits == nil {
		teg.globals.playerUnits = make(map[int][]int)
	}
	teg.globals.playerUnits[pl+100_000] = unformed
}

// addPlayerUnformed adds an unformed noble ID to a player's list.
func addPlayerUnformed(pl, id int) {
	unformed := getPlayerUnformed(pl)
	unformed = append(unformed, id)
	setPlayerUnformed(pl, unformed)
}

// removePlayerUnformed removes an unformed noble ID from a player's list.
func removePlayerUnformed(pl, id int) {
	unformed := getPlayerUnformed(pl)
	for i, v := range unformed {
		if v == id {
			unformed = append(unformed[:i], unformed[i+1:]...)
			break
		}
	}
	setPlayerUnformed(pl, unformed)
}

// containsInt checks if a slice contains a value.
func containsInt(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// joinStrings joins strings with a separator.
func joinStrings(s []string, sep string) string {
	if len(s) == 0 {
		return ""
	}
	result := s[0]
	for i := 1; i < len(s); i++ {
		result += sep + s[i]
	}
	return result
}


