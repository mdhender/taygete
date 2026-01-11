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

// setupStackTestWorld creates a test world with stacking hierarchy:
//
//	Region (1000) - sub_region
//	  └── Province (1001) - sub_forest
//	        ├── Character (2001) - stack leader
//	        │     ├── Character (2002) - stacked under 2001
//	        │     │     └── Character (2003) - stacked under 2002
//	        │     └── Character (2004) - stacked under 2001 (prisoner)
//	        ├── Character (2005) - independent
//	        └── Castle (1003) - sub_castle
//	              └── Character (2006)
func setupStackTestWorld(t *testing.T) func() {
	t.Helper()

	// Initialize boxes
	teg.globals.bx[1000] = &box{kind: T_loc, skind: sub_region}
	teg.globals.bx[1001] = &box{kind: T_loc, skind: sub_forest}
	teg.globals.bx[1003] = &box{kind: T_loc, skind: sub_castle}
	teg.globals.bx[2001] = &box{kind: T_char}
	teg.globals.bx[2002] = &box{kind: T_char}
	teg.globals.bx[2003] = &box{kind: T_char}
	teg.globals.bx[2004] = &box{kind: T_char}
	teg.globals.bx[2005] = &box{kind: T_char}
	teg.globals.bx[2006] = &box{kind: T_char}

	// Initialize char structs for prisoners
	teg.globals.bx[2004].x_char = &entity_char{prisoner: 1}

	// Set location hierarchy via x_loc_info.where
	teg.globals.bx[1000].x_loc_info.where = 0    // region has no parent
	teg.globals.bx[1001].x_loc_info.where = 1000 // province in region
	teg.globals.bx[1003].x_loc_info.where = 1001 // castle in province
	teg.globals.bx[2001].x_loc_info.where = 1001 // char in province (stack leader)
	teg.globals.bx[2002].x_loc_info.where = 2001 // char stacked under 2001
	teg.globals.bx[2003].x_loc_info.where = 2002 // char stacked under 2002
	teg.globals.bx[2004].x_loc_info.where = 2001 // prisoner stacked under 2001
	teg.globals.bx[2005].x_loc_info.where = 1001 // independent char in province
	teg.globals.bx[2006].x_loc_info.where = 1003 // char in castle

	// Set here_lists
	teg.globals.bx[1000].x_loc_info.here_list = []int{1001}
	teg.globals.bx[1001].x_loc_info.here_list = []int{2001, 2005, 1003}
	teg.globals.bx[1003].x_loc_info.here_list = []int{2006}
	teg.globals.bx[2001].x_loc_info.here_list = []int{2002, 2004}
	teg.globals.bx[2002].x_loc_info.here_list = []int{2003}
	teg.globals.bx[2003].x_loc_info.here_list = []int{}
	teg.globals.bx[2004].x_loc_info.here_list = []int{}
	teg.globals.bx[2005].x_loc_info.here_list = []int{}
	teg.globals.bx[2006].x_loc_info.here_list = []int{}

	return func() {
		teg.globals.bx[1000] = nil
		teg.globals.bx[1001] = nil
		teg.globals.bx[1003] = nil
		teg.globals.bx[2001] = nil
		teg.globals.bx[2002] = nil
		teg.globals.bx[2003] = nil
		teg.globals.bx[2004] = nil
		teg.globals.bx[2005] = nil
		teg.globals.bx[2006] = nil
	}
}

func TestHerePos(t *testing.T) {
	cleanup := setupStackTestWorld(t)
	defer cleanup()

	tests := []struct {
		name     string
		who      int
		expected int
	}{
		{"first in province (2001)", 2001, 0},
		{"second in province (2005)", 2005, 1},
		{"third in province (castle)", 1003, 2},
		{"first under 2001 (2002)", 2002, 0},
		{"second under 2001 (prisoner 2004)", 2004, 1},
		{"under 2002 (2003)", 2003, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := here_pos(tt.who)
			if got != tt.expected {
				t.Errorf("here_pos(%d) = %d, want %d", tt.who, got, tt.expected)
			}
		})
	}
}

func TestHerePrecedes(t *testing.T) {
	cleanup := setupStackTestWorld(t)
	defer cleanup()

	tests := []struct {
		name     string
		a, b     int
		expected bool
	}{
		{"2001 precedes 2005", 2001, 2005, true},
		{"2005 does not precede 2001", 2005, 2001, false},
		{"2002 precedes 2004 (under same parent)", 2002, 2004, true},
		{"different locations (2001 vs 2006)", 2001, 2006, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := here_precedes(tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("here_precedes(%d, %d) = %v, want %v", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestFirstPrisonerPos(t *testing.T) {
	cleanup := setupStackTestWorld(t)
	defer cleanup()

	tests := []struct {
		name     string
		where    int
		expected int
	}{
		{"2001 has prisoner at pos 1", 2001, 1},
		{"2002 has no prisoner", 2002, -1},
		{"province has no direct prisoner", 1001, -1},
		{"castle has no prisoner", 1003, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := first_prisoner_pos(tt.where)
			if got != tt.expected {
				t.Errorf("first_prisoner_pos(%d) = %d, want %d", tt.where, got, tt.expected)
			}
		})
	}
}

func TestStackParent(t *testing.T) {
	cleanup := setupStackTestWorld(t)
	defer cleanup()

	tests := []struct {
		name     string
		who      int
		expected int
	}{
		{"2001 is not stacked (in province)", 2001, 0},
		{"2002 is stacked under 2001", 2002, 2001},
		{"2003 is stacked under 2002", 2003, 2002},
		{"2004 is stacked under 2001", 2004, 2001},
		{"2005 is not stacked", 2005, 0},
		{"2006 is in castle (not stacked)", 2006, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stack_parent(tt.who)
			if got != tt.expected {
				t.Errorf("stack_parent(%d) = %d, want %d", tt.who, got, tt.expected)
			}
		})
	}
}

func TestStackLeader(t *testing.T) {
	cleanup := setupStackTestWorld(t)
	defer cleanup()

	tests := []struct {
		name     string
		who      int
		expected int
	}{
		{"2001 is its own leader", 2001, 2001},
		{"2002's leader is 2001", 2002, 2001},
		{"2003's leader is 2001 (through 2002)", 2003, 2001},
		{"2004's leader is 2001", 2004, 2001},
		{"2005 is its own leader", 2005, 2005},
		{"2006 is its own leader", 2006, 2006},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stack_leader(tt.who)
			if got != tt.expected {
				t.Errorf("stack_leader(%d) = %d, want %d", tt.who, got, tt.expected)
			}
		})
	}
}

func TestStackedBeneath(t *testing.T) {
	cleanup := setupStackTestWorld(t)
	defer cleanup()

	tests := []struct {
		name     string
		a, b     int
		expected bool
	}{
		{"2002 is beneath 2001", 2001, 2002, true},
		{"2003 is beneath 2001 (through 2002)", 2001, 2003, true},
		{"2003 is beneath 2002", 2002, 2003, true},
		{"2004 is beneath 2001", 2001, 2004, true},
		{"2001 is not beneath 2002", 2002, 2001, false},
		{"same entity", 2001, 2001, false},
		{"2005 is not beneath 2001", 2001, 2005, false},
		{"2006 is not beneath anyone", 2001, 2006, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stacked_beneath(tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("stacked_beneath(%d, %d) = %v, want %v", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestPromote(t *testing.T) {
	cleanup := setupStackTestWorld(t)
	defer cleanup()

	// Province here_list is [2001, 2005, 1003]
	// Promote 1003 to position 1 (before 2005)
	promote(1003, 1)

	p := rp_loc_info(1001)
	expected := []int{2001, 1003, 2005}

	if len(p.here_list) != len(expected) {
		t.Fatalf("promote: here_list length = %d, want %d", len(p.here_list), len(expected))
	}

	for i, v := range expected {
		if p.here_list[i] != v {
			t.Errorf("promote: here_list[%d] = %d, want %d", i, p.here_list[i], v)
		}
	}
}

func TestHasPrisoner(t *testing.T) {
	cleanup := setupStackTestWorld(t)
	defer cleanup()

	tests := []struct {
		name     string
		who      int
		pris     int
		expected bool
	}{
		{"2001 has 2004 as prisoner", 2001, 2004, true},
		{"2001 has 2002 but not prisoner", 2001, 2002, false},
		{"2005 has no prisoners", 2005, 2004, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := has_prisoner(tt.who, tt.pris)
			if got != tt.expected {
				t.Errorf("has_prisoner(%d, %d) = %v, want %v", tt.who, tt.pris, got, tt.expected)
			}
		})
	}
}

func TestPlayer(t *testing.T) {
	cleanup := setupStackTestWorld(t)
	defer cleanup()

	// Create a player and link characters to it
	teg.globals.bx[3001] = &box{kind: T_player}
	teg.globals.bx[2001].x_char = &entity_char{unit_lord: 3001}
	teg.globals.bx[2002].x_char = &entity_char{unit_lord: 2001}
	teg.globals.bx[2005].x_char = &entity_char{unit_lord: 3001}

	defer func() {
		teg.globals.bx[3001] = nil
	}()

	tests := []struct {
		name     string
		who      int
		expected int
	}{
		{"player directly owns 2001", 2001, 3001},
		{"player owns 2002 through 2001", 2002, 3001},
		{"player directly owns 2005", 2005, 3001},
		{"player entity returns itself", 3001, 3001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := player(tt.who)
			if got != tt.expected {
				t.Errorf("player(%d) = %d, want %d", tt.who, got, tt.expected)
			}
		})
	}
}
