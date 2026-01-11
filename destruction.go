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

// destruction.go - Ship/building destruction ported from src/u.c (Sprint 26.2)

package taygete

// sink_ship handles ship destruction, unbinding any storms and relocating contents.
// If in ocean, kills everyone aboard via kill_stack_ocean.
// If near land, moves everyone to the shore.
// Ported from src/u.c lines 874-928.
func sink_ship(ship int) {
	where := subloc(ship)

	log_write(LOG_SPECIAL, "%s has sunk in %s.", box_name(ship), box_name(subloc(ship)))

	wout(ship, "%s has sunk!", box_name(ship))
	wout(subloc(ship), "%s has sunk!", box_name(ship))

	if subkind(where) == sub_ocean {
		p := rp_loc_info(ship)
		if p != nil {
			for _, who := range p.here_list {
				if kind(who) == T_char {
					kill_stack_ocean(who)
				}
			}
		}
	} else {
		p := rp_loc_info(ship)
		if p != nil {
			for _, who := range p.here_list {
				if kind(who) == T_char {
					move_stack(who, where)
				} else {
					set_where(who, where)
				}
			}
		}
	}

	p := rp_subloc(ship)
	if p != nil && len(p.bound_storms) > 0 {
		for _, storm := range p.bound_storms {
			if kind(storm) == T_storm {
				m := p_misc(storm)
				if m != nil {
					m.bind_storm = 0
				}
			}
		}
		p.bound_storms = nil
	}

	set_where(ship, 0)
	delete_box(ship)
}

// get_rid_of_collapsed_mine removes a collapsed mine from the game.
// Moves any contents out first.
// Ported from src/u.c lines 932-954.
func get_rid_of_collapsed_mine(fort int) {
	if subkind(fort) != sub_mine_collapsed {
		panic("get_rid_of_collapsed_mine: not a collapsed mine")
	}

	where := subloc(fort)

	p := rp_loc_info(fort)
	if p != nil {
		for _, who := range p.here_list {
			if kind(who) == T_char {
				move_stack(who, where)
			} else {
				set_where(who, where)
			}
		}
	}

	set_where(fort, 0)
	delete_box(fort)
}

// building_collapses handles building destruction.
// Mines collapse to sub_mine_collapsed instead of being deleted.
// Castles clear any garrison links.
// Ported from src/u.c lines 958-999.
func building_collapses(fort int) {
	where := subloc(fort)

	log_write(LOG_SPECIAL, "%s collapsed in %s.", box_name(fort), box_name(where))

	vector_char_here(fort)
	vector_add(where)
	wout(VECT, "%s collapses!", box_name(fort))

	p := rp_loc_info(fort)
	if p != nil {
		for _, who := range p.here_list {
			if kind(who) == T_char {
				move_stack(who, where)
			} else {
				set_where(who, where)
			}
		}
	}

	if subkind(fort) == sub_mine {
		change_box_subkind(fort, sub_mine_collapsed)
		p_misc(fort).mine_delay = 8
		return
	}

	if subkind(fort) == sub_castle {
		for i := sub_first(sub_garrison); i != 0; i = sub_next(i) {
			if garrison_castle(i) == fort {
				p_misc(i).garr_castle = 0
			}
		}
	}

	set_where(fort, 0)
	delete_box(fort)
}

// add_structure_damage adds damage to a structure (ship or building).
// Returns true if the structure was destroyed.
// If can_destroy is false, damage is capped at 99.
// Ported from src/u.c lines 1004-1035.
func add_structure_damage(fort, damage int, can_destroy bool) bool {
	if damage < 0 {
		panic("add_structure_damage: negative damage")
	}

	p := p_subloc(fort)

	if int(p.damage)+damage > 100 {
		p.damage = 100
	} else {
		p.damage += uchar(damage)
	}

	if p.damage < 100 {
		return false
	}

	if !can_destroy {
		p.damage = 99
		return false
	}

	if is_ship(fort) {
		sink_ship(fort)
	} else {
		building_collapses(fort)
	}

	return true
}

// find_nearest_land finds the nearest land province from an ocean location.
// Uses random direction walking with fallback to region scan.
// Ported from src/u.c lines 1871-1951.
func find_nearest_land(where int) int {
	origWhere := where
	ret := 0

	tryTwo := 100
	for tryTwo > 0 {
		tryTwo--
		dir := rnd(1, 4)

		tryOne := 1000
		for tryOne > 0 {
			tryOne--

			if subkind(where) != sub_ocean {
				return where
			}

			p := rp_loc_info(where)
			if p != nil {
				for _, i := range p.here_list {
					if subkind(i) == sub_island {
						if kind(i) != T_loc {
							panic("find_nearest_land: island not T_loc")
						}
						ret = i
						break
					}
				}
			}

			if ret != 0 {
				return ret
			}

			where = location_direction(where, dir)

			check := 0
			for where == 0 {
				where = origWhere
				dir = (dir % 4) + 1
				check++
				if check > 4 {
					// No valid direction found, break out to try fallback
					break
				}
				where = location_direction(where, dir)
			}
			if check > 4 {
				break
			}
		}

		if tryTwo == 99 {
			log_write(LOG_CODE, "find_nearest_land: Plan B")
		}
	}

	log_write(LOG_CODE, "find_nearest_land: Plan C")

	var l []int
	for i := kind_first(T_loc); i != 0; i = kind_next(i) {
		if region(i) != region(origWhere) {
			continue
		}
		if loc_depth(i) != LOC_province {
			continue
		}
		if subkind(i) == sub_ocean {
			continue
		}
		l = append(l, i)
	}

	if len(l) < 1 {
		return 0
	}

	ret = l[rnd(0, len(l)-1)]
	return ret
}

// location_direction returns the province in the given direction from where.
// dir is 1-4 for N/E/S/W.
// Ported from src/dir.c lines 145-156.
func location_direction(where, dir int) int {
	dir--

	p := rp_loc(where)
	if p == nil {
		return 0
	}

	provDest := p.prov_dest
	if provDest == nil || dir >= len(provDest) {
		return 0
	}

	return provDest[dir]
}

// Note: garrison_castle is defined in accessor.go
// Note: sub_garrison is defined in glob.go as 64

// vector_add adds to the current vector list.
// Stub implementation.
func vector_add(who int) {
	// TODO: Implement vector system for output
}
