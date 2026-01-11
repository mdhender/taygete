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

// gate.go - Gate and road functions ported from src/gate.c

package taygete

// gates_here returns all gates in the here_list at the given location.
// Ported from loop_gates_here macro in src/loop.h lines 256-267.
func gates_here(where int) []int {
	if !valid_box(where) {
		return nil
	}
	p := rp_loc_info(where)
	if p == nil {
		return nil
	}

	var gates []int
	for _, id := range p.here_list {
		if kind(id) == T_gate {
			gates = append(gates, id)
		}
	}
	return gates
}

// province_gate_here finds a gate in the here_list of a location.
// Returns 0 if no gate found.
// Ported from src/gate.c lines 51-68.
func province_gate_here(where int) int {
	p := rp_loc_info(where)
	if p == nil {
		return 0
	}

	for _, id := range p.here_list {
		if kind(id) == T_gate {
			return id
		}
	}
	return 0
}

// is_gate returns true if the entity is a gate.
func is_gate(n int) bool {
	return kind(n) == T_gate
}

// is_road returns true if the entity is a road.
func is_road(n int) bool {
	return kind(n) == T_road
}

// gate_notify_jumps returns the character to notify when someone jumps through the gate.
// Ported from src/oly.h entity_gate structure.
func gate_notify_jumps(n int) int {
	g := rp_gate(n)
	if g == nil {
		return 0
	}
	return g.notify_jumps
}

// gate_notify_unseal returns the character to notify when someone unseals the gate.
// Ported from src/oly.h entity_gate structure.
func gate_notify_unseal(n int) int {
	g := rp_gate(n)
	if g == nil {
		return 0
	}
	return g.notify_unseal
}

// check_gate_here validates that a gate exists and is at the same sublocation as who.
// Returns true if the gate is valid and at the same location.
// Ported from src/gate.c lines 98-109.
func check_gate_here(who, gate int) bool {
	if kind(gate) != T_gate || subloc(gate) != subloc(who) {
		return false
	}
	return true
}


