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

// loc.go - Spatial model functions ported from src/loc.c and src/u.c

package taygete

// loc_depth returns the depth level of a location based on its subkind.
// Returns LOC_region, LOC_province, LOC_subloc, or LOC_build.
// Ported from src/u.c lines 150-220.
func loc_depth(n int) int {
	switch subkind(n) {
	case sub_region:
		return LOC_region

	case sub_ocean, sub_forest, sub_plain, sub_mountain, sub_desert,
		sub_swamp, sub_under, sub_cloud, sub_tunnel, sub_chamber:
		return LOC_province

	case sub_island, sub_stone_cir, sub_mallorn_grove, sub_bog, sub_cave,
		sub_city, sub_lair, sub_graveyard, sub_ruins, sub_battlefield,
		sub_ench_forest, sub_rocky_hill, sub_tree_circle, sub_pits,
		sub_pasture, sub_oasis, sub_yew_grove, sub_sand_pit,
		sub_sacred_grove, sub_poppy_field, sub_faery_hill, sub_hades_pit:
		return LOC_subloc

	case sub_temple, sub_galley, sub_roundship, sub_castle,
		sub_galley_notdone, sub_roundship_notdone, sub_ghost_ship,
		sub_temple_notdone, sub_inn, sub_inn_notdone, sub_castle_notdone,
		sub_mine, sub_mine_notdone, sub_mine_collapsed, sub_tower,
		sub_tower_notdone, sub_sewer:
		return LOC_build

	default:
		return 0
	}
}

// region returns the ultimate region containing who.
// Ported from src/loc.c lines 133-142.
func region(who int) int {
	for who > 0 && (kind(who) != T_loc || loc_depth(who) != LOC_region) {
		who = loc(who)
	}
	return who
}

// province returns the ultimate province containing who.
// Ported from src/loc.c lines 149-158.
func province(who int) int {
	for who > 0 && (kind(who) != T_loc || loc_depth(who) != LOC_province) {
		who = loc(who)
	}
	return who
}

// subloc returns the immediate location (T_loc or T_ship) containing who.
// Ignores stacked characters.
// Ported from src/loc.c lines 167-178.
func subloc(who int) int {
	for {
		who = loc(who)
		if who <= 0 || kind(who) == T_loc || kind(who) == T_ship {
			break
		}
	}
	return who
}

// viewloc returns the location to use for visibility calculations.
// Steps out from a location until reaching the appropriate viewing level.
// Provinces see into everything except cities, graveyards, sewers, and faery hills.
// Ported from src/loc.c lines 192-207.
func viewloc(who int) int {
	for who > 0 &&
		loc_depth(who) != LOC_province &&
		subkind(who) != sub_city &&
		subkind(who) != sub_graveyard &&
		subkind(who) != sub_sewer &&
		subkind(who) != sub_faery_hill {
		who = loc(who)
	}
	return who
}

// in_safe_now returns true if who is anywhere inside a safe haven.
// Ported from src/loc.c lines 210-224.
func in_safe_now(who int) bool {
	for {
		if safe_haven(who) != 0 {
			return true
		}
		who = loc(who)
		if who <= 0 {
			break
		}
	}
	return false
}

// somewhere_inside returns true if b is nested somewhere inside a.
// Ported from src/loc.c lines 11-26.
func somewhere_inside(a, b int) bool {
	if a == b {
		return false
	}
	for b > 0 {
		b = loc(b)
		if a == b {
			return true
		}
	}
	return false
}

// add_here recursively adds who and all entities in its here_list to l.
// Ported from src/loc.c lines 29-44.
func add_here(who int, l *[]int) {
	if !valid_box(who) {
		panic("add_here: invalid box")
	}

	*l = append(*l, who)

	p := rp_loc_info(who)
	if p == nil {
		panic("add_here: nil loc_info")
	}

	for _, id := range p.here_list {
		add_here(id, l)
	}
}

// all_here returns all entities at and below who.
// Ported from src/loc.c lines 47-64.
func all_here(who int, l *[]int) {
	if !valid_box(who) {
		panic("all_here: invalid box")
	}

	*l = (*l)[:0] // clear the list

	p := rp_loc_info(who)
	if p == nil {
		return
	}

	for _, id := range p.here_list {
		add_here(id, l)
	}
}

// add_char_here recursively adds who and all characters in its here_list to l.
// Ported from src/loc.c lines 67-83.
func add_char_here(who int, l *[]int) {
	if !valid_box(who) {
		panic("add_char_here: invalid box")
	}

	*l = append(*l, who)

	p := rp_loc_info(who)
	if p == nil {
		panic("add_char_here: nil loc_info")
	}

	for _, id := range p.here_list {
		if kind(id) == T_char {
			add_char_here(id, l)
		}
	}
}

// all_char_here returns all characters at and below who.
// Ported from src/loc.c lines 86-104.
func all_char_here(who int, l *[]int) {
	if !valid_box(who) {
		panic("all_char_here: invalid box")
	}

	*l = (*l)[:0] // clear the list

	p := rp_loc_info(who)
	if p == nil {
		return
	}

	for _, id := range p.here_list {
		if kind(id) == T_char {
			add_char_here(id, l)
		}
	}
}

// all_stack returns who plus all characters stacked under who.
// Ported from src/loc.c lines 107-126.
func all_stack(who int, l *[]int) {
	if !valid_box(who) {
		panic("all_stack: invalid box")
	}

	*l = (*l)[:0] // clear the list
	*l = append(*l, who)

	p := rp_loc_info(who)
	if p == nil {
		return
	}

	for _, id := range p.here_list {
		if kind(id) == T_char {
			add_char_here(id, l)
		}
	}
}

// in_here_list returns true if who is in the here_list of loc.
// Ported from src/loc.c lines 298-309.
func in_here_list(loc, who int) bool {
	p := rp_loc_info(loc)
	if p == nil {
		return false
	}
	return IListLookup(p.here_list, who) != -1
}

// add_to_here_list adds who to the here_list of loc.
// Ported from src/loc.c lines 227-234.
func add_to_here_list(loc, who int) {
	if in_here_list(loc, who) {
		panic("add_to_here_list: already in here_list")
	}
	p := p_loc_info(loc)
	p.here_list = append(p.here_list, who)
	if !in_here_list(loc, who) {
		panic("add_to_here_list: failed to add")
	}
}

// remove_from_here_list removes who from the here_list of loc.
// Ported from src/loc.c lines 237-244.
func remove_from_here_list(loc, who int) {
	if !in_here_list(loc, who) {
		panic("remove_from_here_list: not in here_list")
	}
	p := rp_loc_info(loc)
	IListRemValue(&p.here_list, who)
	if in_here_list(loc, who) {
		panic("remove_from_here_list: failed to remove")
	}
}

// set_where moves who from its current location to new_loc.
// Ported from src/loc.c lines 248-274.
func set_where(who, new_loc int) {
	if who == new_loc {
		panic("set_where: who == new_loc")
	}

	old_loc := loc(who)

	if old_loc > 0 {
		remove_from_here_list(old_loc, who)
	}

	if new_loc > 0 {
		add_to_here_list(new_loc, who)
	}

	p_loc_info(who).where = new_loc
}

// first_character returns the first character in the here_list of where.
// Ported from src/loc.c lines 312-329.
func first_character(where int) int {
	p := rp_loc_info(where)
	if p == nil {
		return 0
	}

	for _, id := range p.here_list {
		if kind(id) == T_char {
			return id
		}
	}
	return 0
}

// subloc_here returns the first sublocation of subkind sk in the here_list of where.
// Ported from src/loc.c lines 332-349.
func subloc_here(where int, sk schar) int {
	p := rp_loc_info(where)
	if p == nil {
		return 0
	}

	for _, id := range p.here_list {
		if kind(id) == T_loc && subkind(id) == sk {
			return id
		}
	}
	return 0
}

// count_loc_structures counts locations with subkind a or b in the here_list of where.
// Ported from src/loc.c lines 352-366.
func count_loc_structures(where int, a, b schar) int {
	p := rp_loc_info(where)
	if p == nil {
		return 0
	}

	sum := 0
	for _, id := range p.here_list {
		if kind(id) == T_loc && (subkind(id) == a || subkind(id) == b) {
			sum++
		}
	}
	return sum
}

// building_owner returns the first character in a building.
// Ported from src/loc.c lines 391-397.
func building_owner(where int) int {
	if loc_depth(where) != LOC_build {
		panic("building_owner: not a building")
	}
	return first_character(where)
}

// city_here returns the city sublocation in the here_list of where.
// Equivalent to C macro: #define city_here(a) subloc_here((a), sub_city)
func city_here(where int) int {
	return subloc_here(where, sub_city)
}

// mark_loc_stack_known marks that each member of a stack (or ship) has visited a location.
// Ported from src/loc.c lines 281-295.
func mark_loc_stack_known(stack, where int) {
	if kind(stack) == T_char {
		set_known(stack, where)
	}

	var chars []int
	all_char_here(stack, &chars)
	for _, id := range chars {
		if !is_prisoner(id) {
			set_known(id, where)
		}
	}
}

// loc_hidden returns true if the location is hidden.
// Ported from src/u.c lines 2284-2293.
func loc_hidden(n int) bool {
	l := rp_loc(n)
	if l == nil {
		return false
	}
	return l.hidden != 0
}

// nprovinces returns the total count of provinces in the world.
// The result is cached after first computation.
// Ported from src/u.c lines 2317-2332.
func nprovinces() int {
	if teg.globals.nprov > 0 {
		return teg.globals.nprov
	}

	for id := teg.KindFirst(T_loc); id > 0; id = teg.KindNext(id) {
		if loc_depth(id) == LOC_province {
			teg.globals.nprov++
		}
	}
	return teg.globals.nprov
}

// greater_region returns the "greater region" containing who.
// Returns 0 for normal world regions (not Faery, Hades, Nowhere, Cloud, Tunnel, or Under).
// Returns the special region ID for entities in those realms.
// Ported from src/u.c lines 2415-2430.
func greater_region(who int) int {
	reg := region(who)

	if reg != teg.globals.faeryRegion &&
		reg != teg.globals.hadesRegion &&
		reg != teg.globals.nowhereRegion &&
		reg != teg.globals.cloudRegion &&
		reg != teg.globals.tunnelRegion &&
		reg != teg.globals.underRegion {
		return 0
	}
	return reg
}

// diff_region returns true if two entities are in different greater regions.
// Used to check if teleportation is possible between locations.
// Ported from src/u.c lines 2433-2436.
func diff_region(a, b int) bool {
	return greater_region(a) != greater_region(b)
}

// clear_temps clears the temp field on all entities of the given kind.
// Ported from src/u.c lines 1213-1222.
func clear_temps(k schar) {
	for id := teg.KindFirst(int(k)); id > 0; id = teg.KindNext(id) {
		if b := teg.globals.bx[id]; b != nil {
			b.temp = 0
		}
	}
}

// lookup searches a string table for a case-insensitive match.
// Returns the index of the match, or -1 if not found.
// Ported from src/u.c lines 1415-1426.
func lookup(table []string, s string) int {
	for i, entry := range table {
		if strcasecmp(entry, s) == 0 {
			return i
		}
	}
	return -1
}

// strcasecmp performs a case-insensitive string comparison.
// Returns 0 if strings are equal (ignoring case).
func strcasecmp(a, b string) int {
	if len(a) != len(b) {
		if len(a) < len(b) {
			return -1
		}
		return 1
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			if ca < cb {
				return -1
			}
			return 1
		}
	}
	return 0
}

// in_faery returns true if n is in the Faery region.
func in_faery(n int) bool {
	return region(n) == teg.globals.faeryRegion
}

// in_hades returns true if n is in the Hades region.
func in_hades(n int) bool {
	return region(n) == teg.globals.hadesRegion
}

// in_clouds returns true if n is in the Cloud region.
func in_clouds(n int) bool {
	return region(n) == teg.globals.cloudRegion
}


