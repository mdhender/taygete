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

// cmd_explore.go - Exploration commands ported from src/c1.c
// Sprint 26.4: EXPLORE command

package taygete

// find_lost_items attempts to find unique items in a location.
// Looks for unique items in the location's inventory and gives one to who
// with a probability based on location depth (100% in sublocs, 40% in provinces).
// Returns true if an item was found.
// Ported from src/c1.c lines 33-78.
func find_lost_items(who, where int) bool {
	var item int

	inv := teg.globals.inventories[where]
	for _, e := range inv {
		if item_unique(e.item) == 0 {
			continue
		}

		// Don't take dead bodies out of graveyards; that's what EXHUME is for
		if subkind(where) == sub_graveyard && subkind(e.item) == sub_dead_body {
			continue
		}

		// Don't take magic rings from the market
		if subkind(where) == sub_city && subkind(e.item) == sub_suffuse_ring {
			continue
		}

		item = e.item
		break
	}

	var chance int
	if loc_depth(where) >= LOC_subloc {
		chance = 100
	} else {
		chance = 40
	}

	if item == 0 || rnd(1, 100) > chance {
		return false
	}

	move_item(where, who, item, 1)
	wout(who, "%s found one %s.", box_name(who), box_name(item))

	log_write(LOG_MISC, "%s found %s in %s.",
		box_name(who), box_name(item),
		char_rep_location(where))

	return true
}

// v_explore is the start function for the EXPLORE command.
// Always returns TRUE to begin exploration.
// Ported from src/c1.c lines 26-29.
func v_explore(c *command) int {
	return TRUE
}

// d_explore is the finish function for the EXPLORE command.
// Probability breakdown:
//   - 50% fail
//   - 33% success (find hidden exit)
//   - 17% fail, but message indicating if there is something to find
//
// Also attempts to find lost items first.
// Ported from src/c1.c lines 91-178.
func d_explore(c *command) int {
	where := subloc(c.who)

	if find_lost_items(c.who, where) {
		return TRUE
	}

	// Explore in a ship should explore the surrounding ocean region
	if is_ship(where) && subkind(loc(where)) == sub_ocean {
		where = loc(where)
		find_lost_items(c.who, where)
	}

	r := rnd(1, 100)

	if r <= 50 {
		wout(c.who, "Exploration of %s uncovers no new features.",
			box_code(where))
		return FALSE
	}

	l := exits_from_loc(c.who, where)

	hiddenExits := count_hidden_exits(l)

	// Nothing to find
	if hiddenExits <= 0 {
		wout(c.who, "Exploration of %s uncovers no new features.",
			box_code(where))
		return FALSE
	}

	// Something to find, but a bad roll
	if r <= 67 {
		switch rnd(1, 4) {
		case 1:
			wout(c.who, "Rumors speak of hidden features here, "+
				"but none were found.")
		case 2:
			wout(c.who, "We suspect something is hidden here, "+
				"but did not find anything.")
		case 3:
			wout(c.who, "Something may be hidden here.  "+
				"Further exploration is needed.")
		case 4:
			wout(c.who, "Nothing was found, but further "+
				"exploration looks promising.")
		}
		return FALSE
	}

	// Choose what we found randomly
	i := rnd(1, hiddenExits)

	find_hidden_exit(c.who, l, hidden_count_to_index(i, l))

	return TRUE
}

// exits_from_loc returns a list of exit views from a location.
// Stub: will be implemented with movement system.
func exits_from_loc(who, where int) []*exit_view {
	return nil
}

// count_hidden_exits counts the number of hidden exits in an exit list.
// Stub: will be implemented with movement system.
func count_hidden_exits(l []*exit_view) int {
	if l == nil {
		return 0
	}
	count := 0
	for _, e := range l {
		if e.hidden != 0 {
			count++
		}
	}
	return count
}

// hidden_count_to_index converts a 1-based count to an exit list index.
// Returns the index of the nth hidden exit.
// Stub: will be implemented with movement system.
func hidden_count_to_index(which int, l []*exit_view) int {
	if l == nil {
		return 0
	}
	count := 0
	for i, e := range l {
		if e.hidden != 0 {
			count++
			if count == which {
				return i
			}
		}
	}
	return 0
}

// find_hidden_exit reveals a hidden exit to a character.
// Stub: will be implemented with movement system.
func find_hidden_exit(who int, l []*exit_view, which int) {
	if l == nil || which < 0 || which >= len(l) {
		return
	}
	e := l[which]
	if e == nil {
		return
	}

	// Mark as known to the player
	set_known(who, e.destination)

	wout(who, "Found %s.", box_name(e.destination))
}


