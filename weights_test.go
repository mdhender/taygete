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

func clearWeightsTest() {
	for i := range teg.globals.bx {
		teg.globals.bx[i] = nil
	}
	teg.globals.inventories = make(map[int][]item_ent)
}

func TestAddItemWeight(t *testing.T) {
	defer clearWeightsTest()
	clearWeightsTest()

	itemID := 100
	teg.globals.bx[itemID] = &box{
		kind:  T_item,
		skind: 0,
		x_item: &entity_item{
			weight:   10,
			land_cap: 0,
			ride_cap: 0,
			fly_cap:  0,
			animal:   0,
		},
	}

	horseID := 101
	teg.globals.bx[horseID] = &box{
		kind:  T_item,
		skind: 0,
		x_item: &entity_item{
			weight:   100,
			land_cap: 500,
			ride_cap: 300,
			fly_cap:  0,
			animal:   1,
		},
	}

	pegasusID := 102
	teg.globals.bx[pegasusID] = &box{
		kind:  T_item,
		skind: 0,
		x_item: &entity_item{
			weight:   80,
			land_cap: 400,
			ride_cap: 200,
			fly_cap:  500,
			animal:   1,
		},
	}

	t.Run("basic item weight", func(t *testing.T) {
		var w weights
		add_item_weight(itemID, 5, &w)

		if w.total_weight != 50 {
			t.Errorf("total_weight = %d, want 50", w.total_weight)
		}
		if w.land_weight != 50 {
			t.Errorf("land_weight = %d, want 50", w.land_weight)
		}
		if w.ride_weight != 50 {
			t.Errorf("ride_weight = %d, want 50", w.ride_weight)
		}
		if w.fly_weight != 50 {
			t.Errorf("fly_weight = %d, want 50", w.fly_weight)
		}
		if w.land_cap != 0 {
			t.Errorf("land_cap = %d, want 0", w.land_cap)
		}
		if w.animals != 0 {
			t.Errorf("animals = %d, want 0", w.animals)
		}
	})

	t.Run("horse with capacity", func(t *testing.T) {
		var w weights
		add_item_weight(horseID, 2, &w)

		if w.total_weight != 200 {
			t.Errorf("total_weight = %d, want 200", w.total_weight)
		}
		if w.land_cap != 1000 {
			t.Errorf("land_cap = %d, want 1000", w.land_cap)
		}
		if w.land_weight != 0 {
			t.Errorf("land_weight = %d, want 0", w.land_weight)
		}
		if w.ride_cap != 600 {
			t.Errorf("ride_cap = %d, want 600", w.ride_cap)
		}
		if w.ride_weight != 0 {
			t.Errorf("ride_weight = %d, want 0", w.ride_weight)
		}
		if w.fly_weight != 200 {
			t.Errorf("fly_weight = %d, want 200 (no fly cap)", w.fly_weight)
		}
		if w.animals != 2 {
			t.Errorf("animals = %d, want 2", w.animals)
		}
	})

	t.Run("pegasus with fly capacity", func(t *testing.T) {
		var w weights
		add_item_weight(pegasusID, 3, &w)

		if w.total_weight != 240 {
			t.Errorf("total_weight = %d, want 240", w.total_weight)
		}
		if w.land_cap != 1200 {
			t.Errorf("land_cap = %d, want 1200", w.land_cap)
		}
		if w.ride_cap != 600 {
			t.Errorf("ride_cap = %d, want 600", w.ride_cap)
		}
		if w.fly_cap != 1500 {
			t.Errorf("fly_cap = %d, want 1500", w.fly_cap)
		}
		if w.animals != 3 {
			t.Errorf("animals = %d, want 3", w.animals)
		}
	})

	t.Run("cumulative weights", func(t *testing.T) {
		var w weights
		add_item_weight(itemID, 10, &w)
		add_item_weight(horseID, 1, &w)

		if w.total_weight != 200 {
			t.Errorf("total_weight = %d, want 200", w.total_weight)
		}
		if w.land_cap != 500 {
			t.Errorf("land_cap = %d, want 500", w.land_cap)
		}
		if w.land_weight != 100 {
			t.Errorf("land_weight = %d, want 100", w.land_weight)
		}
		if w.animals != 1 {
			t.Errorf("animals = %d, want 1", w.animals)
		}
	})
}

func TestDetermineUnitWeights(t *testing.T) {
	defer clearWeightsTest()
	clearWeightsTest()

	peasantID := item_peasant
	teg.globals.bx[peasantID] = &box{
		kind:  T_item,
		skind: 0,
		x_item: &entity_item{
			weight:   10,
			land_cap: 0,
			ride_cap: 0,
			fly_cap:  0,
			animal:   0,
		},
	}

	goldID := item_gold
	teg.globals.bx[goldID] = &box{
		kind:  T_item,
		skind: 0,
		x_item: &entity_item{
			weight:   1,
			land_cap: 0,
			ride_cap: 0,
			fly_cap:  0,
			animal:   0,
		},
	}

	charID := 1001
	teg.globals.bx[charID] = &box{
		kind:  T_char,
		skind: 0,
		x_char: &entity_char{
			unit_item: 0,
		},
	}
	teg.globals.inventories[charID] = []item_ent{
		{item: goldID, qty: 100},
	}

	t.Run("unit with default noble_item", func(t *testing.T) {
		var w weights
		determine_unit_weights(charID, &w)

		if w.total_weight != 110 {
			t.Errorf("total_weight = %d, want 110 (10 peasant + 100 gold)", w.total_weight)
		}
	})

	warriorID := 50
	teg.globals.bx[warriorID] = &box{
		kind:  T_item,
		skind: 0,
		x_item: &entity_item{
			weight:   20,
			land_cap: 0,
			ride_cap: 0,
			fly_cap:  0,
			animal:   0,
		},
	}

	charID2 := 1002
	teg.globals.bx[charID2] = &box{
		kind:  T_char,
		skind: 0,
		x_char: &entity_char{
			unit_item: schar(warriorID),
		},
	}
	teg.globals.inventories[charID2] = []item_ent{}

	t.Run("unit with custom noble_item", func(t *testing.T) {
		var w weights
		determine_unit_weights(charID2, &w)

		if w.total_weight != 20 {
			t.Errorf("total_weight = %d, want 20 (warrior)", w.total_weight)
		}
	})
}

func TestDetermineStackWeights(t *testing.T) {
	defer clearWeightsTest()
	clearWeightsTest()

	peasantID := item_peasant
	teg.globals.bx[peasantID] = &box{
		kind:  T_item,
		skind: 0,
		x_item: &entity_item{
			weight:   10,
			land_cap: 0,
			ride_cap: 0,
			fly_cap:  0,
			animal:   0,
		},
	}

	goldID := item_gold
	teg.globals.bx[goldID] = &box{
		kind:  T_item,
		skind: 0,
		x_item: &entity_item{
			weight:   1,
			land_cap: 0,
			ride_cap: 0,
			fly_cap:  0,
			animal:   0,
		},
	}

	leaderID := 1001
	teg.globals.bx[leaderID] = &box{
		kind:  T_char,
		skind: 0,
		x_char: &entity_char{
			unit_item: 0,
		},
	}
	teg.globals.inventories[leaderID] = []item_ent{
		{item: goldID, qty: 50},
	}

	followerID := 1002
	teg.globals.bx[followerID] = &box{
		kind:  T_char,
		skind: 0,
		x_char: &entity_char{
			unit_item: 0,
		},
	}
	teg.globals.inventories[followerID] = []item_ent{
		{item: goldID, qty: 30},
	}

	teg.globals.bx[leaderID].x_loc_info.here_list = []int{followerID}

	t.Run("stack with leader and follower", func(t *testing.T) {
		var w weights
		determine_stack_weights(leaderID, &w)

		expectedLeader := 10 + 50
		expectedFollower := 10 + 30
		expected := expectedLeader + expectedFollower
		if w.total_weight != expected {
			t.Errorf("total_weight = %d, want %d", w.total_weight, expected)
		}
	})
}

func TestShipWeight(t *testing.T) {
	defer clearWeightsTest()
	clearWeightsTest()

	peasantID := item_peasant
	teg.globals.bx[peasantID] = &box{
		kind:  T_item,
		skind: 0,
		x_item: &entity_item{
			weight:   10,
			land_cap: 0,
			ride_cap: 0,
			fly_cap:  0,
			animal:   0,
		},
	}

	goldID := item_gold
	teg.globals.bx[goldID] = &box{
		kind:  T_item,
		skind: 0,
		x_item: &entity_item{
			weight:   1,
			land_cap: 0,
			ride_cap: 0,
			fly_cap:  0,
			animal:   0,
		},
	}

	shipID := 500
	teg.globals.bx[shipID] = &box{
		kind:  T_ship,
		skind: 0,
		x_subloc: &entity_subloc{
			capacity: 10000,
			damage:   0,
		},
	}

	charID := 1001
	teg.globals.bx[charID] = &box{
		kind:  T_char,
		skind: 0,
		x_char: &entity_char{
			unit_item: 0,
		},
	}
	teg.globals.inventories[charID] = []item_ent{
		{item: goldID, qty: 100},
	}

	teg.globals.bx[shipID].x_loc_info.here_list = []int{charID}

	t.Run("ship with one passenger", func(t *testing.T) {
		w := ship_weight(shipID)

		expected := 10 + 100
		if w != expected {
			t.Errorf("ship_weight = %d, want %d", w, expected)
		}
	})
}

func TestShipCap(t *testing.T) {
	defer clearWeightsTest()
	clearWeightsTest()

	shipID := 500
	teg.globals.bx[shipID] = &box{
		kind:  T_ship,
		skind: 0,
		x_subloc: &entity_subloc{
			capacity: 10000,
			damage:   0,
		},
	}

	t.Run("undamaged ship", func(t *testing.T) {
		cap := ship_cap(shipID)
		if cap != 10000 {
			t.Errorf("ship_cap = %d, want 10000", cap)
		}
	})

	t.Run("50% damaged ship", func(t *testing.T) {
		teg.globals.bx[shipID].x_subloc.damage = 50
		cap := ship_cap(shipID)
		if cap != 5000 {
			t.Errorf("ship_cap = %d, want 5000", cap)
		}
	})

	t.Run("25% damaged ship", func(t *testing.T) {
		teg.globals.bx[shipID].x_subloc.damage = 25
		cap := ship_cap(shipID)
		if cap != 7500 {
			t.Errorf("ship_cap = %d, want 7500", cap)
		}
	})

	t.Run("100% damaged ship", func(t *testing.T) {
		teg.globals.bx[shipID].x_subloc.damage = 100
		cap := ship_cap(shipID)
		if cap != 0 {
			t.Errorf("ship_cap = %d, want 0", cap)
		}
	})
}

func TestShipCapRaw(t *testing.T) {
	defer clearWeightsTest()
	clearWeightsTest()

	t.Run("ship with capacity", func(t *testing.T) {
		shipID := 500
		teg.globals.bx[shipID] = &box{
			kind:  T_ship,
			skind: 0,
			x_subloc: &entity_subloc{
				capacity: 25000,
			},
		}

		cap := ship_cap_raw(shipID)
		if cap != 25000 {
			t.Errorf("ship_cap_raw = %d, want 25000", cap)
		}
	})

	t.Run("nil subloc", func(t *testing.T) {
		emptyID := 501
		teg.globals.bx[emptyID] = &box{
			kind:  T_ship,
			skind: 0,
		}

		cap := ship_cap_raw(emptyID)
		if cap != 0 {
			t.Errorf("ship_cap_raw = %d, want 0", cap)
		}
	})
}



func TestLoopCharHere(t *testing.T) {
	defer clearWeightsTest()
	clearWeightsTest()

	locID := 1000
	char1ID := 2001
	char2ID := 2002
	itemID := 3001

	teg.globals.bx[locID] = &box{
		kind:  T_loc,
		skind: 0,
	}
	teg.globals.bx[char1ID] = &box{
		kind:  T_char,
		skind: 0,
	}
	teg.globals.bx[char2ID] = &box{
		kind:  T_char,
		skind: 0,
	}
	teg.globals.bx[itemID] = &box{
		kind:  T_item,
		skind: 0,
	}

	teg.globals.bx[locID].x_loc_info.here_list = []int{char1ID, itemID, char2ID}

	t.Run("filters characters only", func(t *testing.T) {
		var chars []int
		loop_char_here(locID, &chars)

		if len(chars) != 2 {
			t.Errorf("loop_char_here returned %d chars, want 2", len(chars))
		}
		if chars[0] != char1ID || chars[1] != char2ID {
			t.Errorf("loop_char_here = %v, want [%d, %d]", chars, char1ID, char2ID)
		}
	})
}
