// Copyright (c) 2026 Michael D Henderson. All rights reserved.

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
