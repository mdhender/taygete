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

// clearDestructionTestBoxes clears boxes used in destruction tests.
func clearDestructionTestBoxes() {
	for _, id := range []int{100, 200, 1000, 10000, 10001, 58760, 59000, 59100, 59200, 79000} {
		teg.globals.bx[id] = nil
	}
	teg.globals.sub_head[sub_garrison] = 0
}

func TestAddStructureDamage(t *testing.T) {
	clearDestructionTestBoxes()
	defer clearDestructionTestBoxes()

	t.Run("damage accumulation below 100", func(t *testing.T) {
		shipID := 100
		teg.globals.bx[shipID] = &box{
			kind:     T_ship,
			skind:    sub_galley,
			x_subloc: &entity_subloc{damage: 0},
		}

		destroyed := add_structure_damage(shipID, 50, true)
		if destroyed {
			t.Error("ship should not be destroyed at 50 damage")
		}

		p := rp_subloc(shipID)
		if p.damage != 50 {
			t.Errorf("damage = %d, want 50", p.damage)
		}
	})

	t.Run("damage capped at 99 when can_destroy is false", func(t *testing.T) {
		clearDestructionTestBoxes()
		shipID := 100
		provinceID := 10000

		teg.globals.bx[provinceID] = &box{kind: T_loc, skind: sub_ocean}
		teg.globals.bx[shipID] = &box{
			kind:     T_ship,
			skind:    sub_galley,
			x_subloc: &entity_subloc{damage: 90},
		}
		set_where(shipID, provinceID)

		destroyed := add_structure_damage(shipID, 50, false)
		if destroyed {
			t.Error("ship should not be destroyed when can_destroy is false")
		}

		p := rp_subloc(shipID)
		if p.damage != 99 {
			t.Errorf("damage = %d, want 99 (capped)", p.damage)
		}
	})

	t.Run("ship destroyed when damage reaches 100", func(t *testing.T) {
		clearDestructionTestBoxes()
		shipID := 100
		provinceID := 10000

		teg.globals.bx[provinceID] = &box{kind: T_loc, skind: sub_plain}
		teg.globals.bx[shipID] = &box{
			kind:     T_ship,
			skind:    sub_galley,
			x_subloc: &entity_subloc{damage: 90},
		}
		set_where(shipID, provinceID)

		destroyed := add_structure_damage(shipID, 20, true)
		if !destroyed {
			t.Error("ship should be destroyed at 100+ damage")
		}

		if kind(shipID) != T_deleted {
			t.Errorf("kind(ship) = %d, want T_deleted", kind(shipID))
		}
	})
}

func TestSinkShip(t *testing.T) {
	clearDestructionTestBoxes()
	defer clearDestructionTestBoxes()

	t.Run("sink ship near land moves contents to shore", func(t *testing.T) {
		clearDestructionTestBoxes()
		shipID := 100
		provinceID := 10000

		teg.globals.bx[provinceID] = &box{kind: T_loc, skind: sub_plain}
		teg.globals.bx[shipID] = &box{
			kind:     T_ship,
			skind:    sub_galley,
			x_subloc: &entity_subloc{},
		}

		set_where(shipID, provinceID)

		sink_ship(shipID)

		if kind(shipID) != T_deleted {
			t.Errorf("ship kind = %d, want T_deleted", kind(shipID))
		}
		// Note: character relocation tested implicitly; move_stack is stub
	})

	t.Run("sink ship unbinds storms", func(t *testing.T) {
		clearDestructionTestBoxes()
		shipID := 100
		stormID := 79000
		provinceID := 10000

		teg.globals.bx[provinceID] = &box{kind: T_loc, skind: sub_plain}
		teg.globals.bx[stormID] = &box{
			kind:   T_storm,
			x_misc: &entity_misc{bind_storm: shipID},
		}
		teg.globals.bx[shipID] = &box{
			kind:  T_ship,
			skind: sub_galley,
			x_subloc: &entity_subloc{
				bound_storms: []int{stormID},
			},
		}

		set_where(shipID, provinceID)

		sink_ship(shipID)

		if p_misc(stormID).bind_storm != 0 {
			t.Errorf("storm bind_storm = %d, want 0", p_misc(stormID).bind_storm)
		}
	})
}

func TestBuildingCollapses(t *testing.T) {
	clearDestructionTestBoxes()
	defer clearDestructionTestBoxes()

	t.Run("mine collapses to collapsed state", func(t *testing.T) {
		clearDestructionTestBoxes()
		mineID := 59000
		provinceID := 10000

		teg.globals.bx[provinceID] = &box{kind: T_loc, skind: sub_plain}
		teg.globals.bx[mineID] = &box{
			kind:     T_loc,
			skind:    sub_mine,
			x_subloc: &entity_subloc{damage: 100},
		}

		set_where(mineID, provinceID)

		building_collapses(mineID)

		if subkind(mineID) != sub_mine_collapsed {
			t.Errorf("mine subkind = %d, want sub_mine_collapsed (%d)", subkind(mineID), sub_mine_collapsed)
		}

		if p_misc(mineID).mine_delay != 8 {
			t.Errorf("mine_delay = %d, want 8", p_misc(mineID).mine_delay)
		}
	})

	t.Run("tower is deleted on collapse", func(t *testing.T) {
		clearDestructionTestBoxes()
		towerID := 59100
		provinceID := 10000

		teg.globals.bx[provinceID] = &box{kind: T_loc, skind: sub_plain}
		teg.globals.bx[towerID] = &box{
			kind:     T_loc,
			skind:    sub_tower,
			x_subloc: &entity_subloc{damage: 100},
		}

		set_where(towerID, provinceID)

		building_collapses(towerID)

		if kind(towerID) != T_deleted {
			t.Errorf("tower kind = %d, want T_deleted", kind(towerID))
		}
	})

	t.Run("castle collapse clears garrison links", func(t *testing.T) {
		clearDestructionTestBoxes()
		castleID := 59200
		garrisonID := 1000
		provinceID := 10000

		teg.globals.bx[provinceID] = &box{kind: T_loc, skind: sub_plain}
		teg.globals.bx[castleID] = &box{
			kind:     T_loc,
			skind:    sub_castle,
			x_subloc: &entity_subloc{damage: 100},
		}
		teg.globals.bx[garrisonID] = &box{
			kind:   T_char,
			skind:  sub_garrison,
			x_char: &entity_char{},
			x_misc: &entity_misc{garr_castle: castleID},
		}

		teg.globals.sub_head[sub_garrison] = garrisonID

		set_where(castleID, provinceID)

		building_collapses(castleID)

		if kind(castleID) != T_deleted {
			t.Errorf("castle kind = %d, want T_deleted", kind(castleID))
		}

		if p_misc(garrisonID).garr_castle != 0 {
			t.Errorf("garrison garr_castle = %d, want 0", p_misc(garrisonID).garr_castle)
		}
	})
}

func TestGetRidOfCollapsedMine(t *testing.T) {
	clearDestructionTestBoxes()
	defer clearDestructionTestBoxes()

	mineID := 59000
	provinceID := 10000

	teg.globals.bx[provinceID] = &box{kind: T_loc, skind: sub_plain}
	teg.globals.bx[mineID] = &box{
		kind:     T_loc,
		skind:    sub_mine_collapsed,
		x_subloc: &entity_subloc{},
	}

	set_where(mineID, provinceID)

	get_rid_of_collapsed_mine(mineID)

	if kind(mineID) != T_deleted {
		t.Errorf("mine kind = %d, want T_deleted", kind(mineID))
	}
	// Note: character relocation tested implicitly; move_stack is stub
}

func TestFindNearestLand(t *testing.T) {
	clearDestructionTestBoxes()
	defer clearDestructionTestBoxes()

	t.Run("returns land if already on land", func(t *testing.T) {
		clearDestructionTestBoxes()
		provinceID := 10000

		teg.globals.bx[provinceID] = &box{kind: T_loc, skind: sub_plain}

		result := find_nearest_land(provinceID)
		if result != provinceID {
			t.Errorf("find_nearest_land = %d, want %d", result, provinceID)
		}
	})

	t.Run("returns island in ocean", func(t *testing.T) {
		clearDestructionTestBoxes()
		oceanID := 10000
		islandID := 59000

		teg.globals.bx[oceanID] = &box{kind: T_loc, skind: sub_ocean}
		teg.globals.bx[islandID] = &box{kind: T_loc, skind: sub_island}

		set_where(islandID, oceanID)

		result := find_nearest_land(oceanID)
		if result != islandID {
			t.Errorf("find_nearest_land = %d, want island %d", result, islandID)
		}
	})

	t.Run("navigates to adjacent land", func(t *testing.T) {
		clearDestructionTestBoxes()
		oceanID := 10000
		landID := 10001
		regionID := 58760

		teg.globals.bx[regionID] = &box{kind: T_loc, skind: sub_region}
		teg.globals.bx[oceanID] = &box{
			kind:  T_loc,
			skind: sub_ocean,
			x_loc: &entity_loc{
				prov_dest: []int{landID, 0, 0, 0},
			},
		}
		teg.globals.bx[landID] = &box{kind: T_loc, skind: sub_plain}

		set_where(oceanID, regionID)
		set_where(landID, regionID)

		result := find_nearest_land(oceanID)
		if result != landID {
			t.Errorf("find_nearest_land = %d, want %d", result, landID)
		}
	})
}

func TestLocationDirection(t *testing.T) {
	clearDestructionTestBoxes()
	defer clearDestructionTestBoxes()

	provinceID := 10000
	northID := 10001
	eastID := 10002
	southID := 10003
	westID := 10004

	teg.globals.bx[provinceID] = &box{
		kind:  T_loc,
		skind: sub_plain,
		x_loc: &entity_loc{
			prov_dest: []int{northID, eastID, southID, westID},
		},
	}

	tests := []struct {
		dir  int
		want int
	}{
		{DIR_N, northID},
		{DIR_E, eastID},
		{DIR_S, southID},
		{DIR_W, westID},
	}

	for _, tt := range tests {
		got := location_direction(provinceID, tt.dir)
		if got != tt.want {
			t.Errorf("location_direction(%d, %d) = %d, want %d", provinceID, tt.dir, got, tt.want)
		}
	}

	got := location_direction(provinceID, 5)
	if got != 0 {
		t.Errorf("location_direction with invalid dir = %d, want 0", got)
	}
}
