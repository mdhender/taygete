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

import "testing"

func TestGatesHere(t *testing.T) {
	defer clearBx()
	clearBx()

	// Create a province location
	teg.globals.bx[100] = &box{
		kind:       T_loc,
		skind:      sub_plain,
		x_loc_info: loc_info{here_list: []int{200, 300, 400}},
	}

	// Create entities in the province
	teg.globals.bx[200] = &box{kind: T_gate} // gate
	teg.globals.bx[300] = &box{kind: T_char} // character
	teg.globals.bx[400] = &box{kind: T_gate} // another gate

	gates := gates_here(100)

	if len(gates) != 2 {
		t.Errorf("gates_here(100) returned %d gates, want 2", len(gates))
	}

	if len(gates) >= 1 && gates[0] != 200 {
		t.Errorf("gates_here(100)[0] = %d, want 200", gates[0])
	}
	if len(gates) >= 2 && gates[1] != 400 {
		t.Errorf("gates_here(100)[1] = %d, want 400", gates[1])
	}
}

func TestGatesHereEmpty(t *testing.T) {
	defer clearBx()
	clearBx()

	// Create a province with no gates
	teg.globals.bx[100] = &box{
		kind:       T_loc,
		skind:      sub_plain,
		x_loc_info: loc_info{here_list: []int{200, 300}},
	}

	teg.globals.bx[200] = &box{kind: T_char}
	teg.globals.bx[300] = &box{kind: T_loc}

	gates := gates_here(100)

	if len(gates) != 0 {
		t.Errorf("gates_here(100) returned %d gates, want 0", len(gates))
	}
}

func TestGatesHereInvalidBox(t *testing.T) {
	defer clearBx()
	clearBx()

	gates := gates_here(999)

	if gates != nil {
		t.Error("gates_here(999) should return nil for invalid box")
	}
}

func TestProvinceGateHere(t *testing.T) {
	defer clearBx()
	clearBx()

	// Create a province with a gate
	teg.globals.bx[100] = &box{
		kind:       T_loc,
		skind:      sub_plain,
		x_loc_info: loc_info{here_list: []int{200, 300, 400}},
	}

	teg.globals.bx[200] = &box{kind: T_char}
	teg.globals.bx[300] = &box{kind: T_gate}
	teg.globals.bx[400] = &box{kind: T_gate}

	gate := province_gate_here(100)

	if gate != 300 {
		t.Errorf("province_gate_here(100) = %d, want 300 (first gate)", gate)
	}
}

func TestProvinceGateHereNoGate(t *testing.T) {
	defer clearBx()
	clearBx()

	teg.globals.bx[100] = &box{
		kind:       T_loc,
		skind:      sub_plain,
		x_loc_info: loc_info{here_list: []int{200}},
	}

	teg.globals.bx[200] = &box{kind: T_char}

	gate := province_gate_here(100)

	if gate != 0 {
		t.Errorf("province_gate_here(100) = %d, want 0 (no gate)", gate)
	}
}

func TestIsGate(t *testing.T) {
	defer clearBx()
	clearBx()

	teg.globals.bx[100] = &box{kind: T_gate}
	teg.globals.bx[200] = &box{kind: T_road}
	teg.globals.bx[300] = &box{kind: T_loc}

	if !is_gate(100) {
		t.Error("is_gate(100) should be true")
	}
	if is_gate(200) {
		t.Error("is_gate(200) should be false for road")
	}
	if is_gate(300) {
		t.Error("is_gate(300) should be false for loc")
	}
}

func TestIsRoad(t *testing.T) {
	defer clearBx()
	clearBx()

	teg.globals.bx[100] = &box{kind: T_gate}
	teg.globals.bx[200] = &box{kind: T_road}
	teg.globals.bx[300] = &box{kind: T_loc}

	if is_road(100) {
		t.Error("is_road(100) should be false for gate")
	}
	if !is_road(200) {
		t.Error("is_road(200) should be true")
	}
	if is_road(300) {
		t.Error("is_road(300) should be false for loc")
	}
}

func TestRoadDest(t *testing.T) {
	defer clearBx()
	clearBx()

	teg.globals.bx[100] = &box{
		kind: T_road,
		x_gate: &entity_gate{
			to_loc: 500,
		},
	}

	dest := road_dest(100)
	if dest != 500 {
		t.Errorf("road_dest(100) = %d, want 500", dest)
	}
}

func TestRoadDestNil(t *testing.T) {
	defer clearBx()
	clearBx()

	teg.globals.bx[100] = &box{kind: T_road}

	dest := road_dest(100)
	if dest != 0 {
		t.Errorf("road_dest(100) = %d, want 0 for nil gate struct", dest)
	}
}

func TestRoadHidden(t *testing.T) {
	defer clearBx()
	clearBx()

	teg.globals.bx[100] = &box{
		kind: T_road,
		x_gate: &entity_gate{
			road_hidden: 1,
		},
	}

	hidden := road_hidden(100)
	if hidden != 1 {
		t.Errorf("road_hidden(100) = %d, want 1", hidden)
	}
}

func TestRoadHiddenFalse(t *testing.T) {
	defer clearBx()
	clearBx()

	teg.globals.bx[100] = &box{
		kind: T_road,
		x_gate: &entity_gate{
			road_hidden: 0,
		},
	}

	hidden := road_hidden(100)
	if hidden != 0 {
		t.Errorf("road_hidden(100) = %d, want 0", hidden)
	}
}

func TestGateNotifyJumps(t *testing.T) {
	defer clearBx()
	clearBx()

	teg.globals.bx[100] = &box{
		kind: T_gate,
		x_gate: &entity_gate{
			notify_jumps: 999,
		},
	}

	notify := gate_notify_jumps(100)
	if notify != 999 {
		t.Errorf("gate_notify_jumps(100) = %d, want 999", notify)
	}
}

func TestGateNotifyUnseal(t *testing.T) {
	defer clearBx()
	clearBx()

	teg.globals.bx[100] = &box{
		kind: T_gate,
		x_gate: &entity_gate{
			notify_unseal: 888,
		},
	}

	notify := gate_notify_unseal(100)
	if notify != 888 {
		t.Errorf("gate_notify_unseal(100) = %d, want 888", notify)
	}
}

func TestCheckGateHere(t *testing.T) {
	defer clearBx()
	clearBx()

	// Create a province
	teg.globals.bx[100] = &box{
		kind:       T_loc,
		skind:      sub_plain,
		x_loc_info: loc_info{here_list: []int{200, 300}},
	}

	// Character at province
	teg.globals.bx[200] = &box{
		kind:       T_char,
		x_loc_info: loc_info{where: 100},
	}

	// Gate at same province
	teg.globals.bx[300] = &box{
		kind:       T_gate,
		x_loc_info: loc_info{where: 100},
	}

	// Gate at different location
	teg.globals.bx[400] = &box{
		kind:       T_gate,
		x_loc_info: loc_info{where: 500},
	}

	teg.globals.bx[500] = &box{
		kind:  T_loc,
		skind: sub_forest,
	}

	if !check_gate_here(200, 300) {
		t.Error("check_gate_here(200, 300) should be true (same subloc)")
	}

	if check_gate_here(200, 400) {
		t.Error("check_gate_here(200, 400) should be false (different subloc)")
	}

	// Non-gate should fail
	if check_gate_here(200, 200) {
		t.Error("check_gate_here(200, 200) should be false (not a gate)")
	}
}

func TestDiffRegion(t *testing.T) {
	defer clearBx()
	clearBx()

	// Save and set special region IDs
	oldFaery := teg.globals.faeryRegion
	teg.globals.faeryRegion = 3
	defer func() { teg.globals.faeryRegion = oldFaery }()

	// Create normal regions (1 and 2) and Faery region (3)
	teg.globals.bx[1] = &box{
		kind:  T_loc,
		skind: sub_region,
	}
	teg.globals.bx[2] = &box{
		kind:  T_loc,
		skind: sub_region,
	}
	teg.globals.bx[3] = &box{
		kind:  T_loc,
		skind: sub_region,
	}

	// Provinces in different normal regions
	teg.globals.bx[100] = &box{
		kind:       T_loc,
		skind:      sub_plain,
		x_loc_info: loc_info{where: 1},
	}
	teg.globals.bx[200] = &box{
		kind:       T_loc,
		skind:      sub_plain,
		x_loc_info: loc_info{where: 2},
	}

	// Same region as 100
	teg.globals.bx[300] = &box{
		kind:       T_loc,
		skind:      sub_plain,
		x_loc_info: loc_info{where: 1},
	}

	// Province in Faery
	teg.globals.bx[400] = &box{
		kind:       T_loc,
		skind:      sub_plain,
		x_loc_info: loc_info{where: 3},
	}

	// diff_region checks greater_region, which returns 0 for normal world regions.
	// Two normal regions both return greater_region=0, so diff_region is false.
	if diff_region(100, 200) {
		t.Error("diff_region(100, 200) should be false (both in normal world)")
	}

	if diff_region(100, 300) {
		t.Error("diff_region(100, 300) should be false (same region)")
	}

	// Normal vs Faery should be different
	if !diff_region(100, 400) {
		t.Error("diff_region(100, 400) should be true (normal vs faery)")
	}
}
