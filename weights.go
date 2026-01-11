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

// weights.go - Weight & capacity helpers ported from src/u.c (Sprint 25.5)

package taygete

// add_item_weight adds the weight and capacity of an item to a weights struct.
// Ported from src/u.c lines 1303-1331.
func add_item_weight(item, qty int, w *weights) {
	wt := int(item_weight(item)) * qty
	lc := int(item_land_cap(item))
	rc := int(item_ride_cap(item))
	fc := int(item_fly_cap(item))

	if lc != 0 {
		w.land_cap += max(lc, 0) * qty
	} else {
		w.land_weight += wt
	}

	if rc != 0 {
		w.ride_cap += max(rc, 0) * qty
	} else {
		w.ride_weight += wt
	}

	if fc != 0 {
		w.fly_cap += max(fc, 0) * qty
	} else {
		w.fly_weight += wt
	}

	w.total_weight += wt

	if item_animal(item) != 0 {
		w.animals += qty
	}
}

// determine_unit_weights calculates the weight and capacity of a single unit.
// Ported from src/u.c lines 1335-1355.
func determine_unit_weights(who int, w *weights) {
	if kind(who) != T_char {
		panic("determine_unit_weights: who is not a character")
	}

	*w = weights{}

	unitBase := int(noble_item(who))
	if unitBase == 0 {
		unitBase = item_peasant
	}

	add_item_weight(unitBase, 1, w)

	inv := teg.globals.inventories[who]
	for _, e := range inv {
		add_item_weight(e.item, e.qty, w)
	}
}

// determine_stack_weights calculates the weight and capacity of a stack.
// Ported from src/u.c lines 1359-1379.
func determine_stack_weights(who int, w *weights) {
	determine_unit_weights(who, w)

	p := rp_loc_info(who)
	if p == nil {
		return
	}

	for _, i := range p.here_list {
		if kind(i) != T_char {
			continue
		}
		var v weights
		determine_unit_weights(i, &v)
		w.total_weight += v.total_weight
		w.land_cap += v.land_cap
		w.ride_cap += v.ride_cap
		w.fly_cap += v.fly_cap
		w.land_weight += v.land_weight
		w.ride_weight += v.ride_weight
		w.fly_weight += v.fly_weight
		w.animals += v.animals
	}
}

// ship_weight calculates the total weight of all cargo on a ship.
// Ported from src/u.c lines 1385-1411.
func ship_weight(ship int) int {
	if kind(ship) != T_ship {
		panic("ship_weight: ship is not a ship")
	}

	sum := 0

	var chars []int
	loop_char_here(ship, &chars)

	for _, i := range chars {
		var w weights
		determine_unit_weights(i, &w)
		sum += w.total_weight
	}

	return sum
}

// ship_cap returns the effective capacity of a ship, reduced by damage.
// Ported from src/u.c lines 2404-2412.
func ship_cap(ship int) int {
	sc := ship_cap_raw(ship)
	dam := int(loc_damage(ship))

	sc -= sc * dam / 100

	return sc
}

// loop_char_here collects all characters directly at a location.
// This is a helper that matches the C loop_char_here macro.
func loop_char_here(where int, l *[]int) {
	*l = nil
	p := rp_loc_info(where)
	if p == nil {
		return
	}
	for _, id := range p.here_list {
		if kind(id) == T_char {
			*l = append(*l, id)
		}
	}
}
