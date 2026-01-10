// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

import "testing"

// setupLocTestWorld creates a test world hierarchy:
//
//	Region (1000) - sub_region
//	  └── Province (1001) - sub_forest
//	        ├── City (1002) - sub_city (safe haven)
//	        │     └── Castle (1003) - sub_castle
//	        │           └── Character (2001)
//	        ├── Graveyard (1004) - sub_graveyard
//	        └── Character (2002)
//	              └── Character (2003) - stacked under 2002
func setupLocTestWorld(t *testing.T) func() {
	t.Helper()

	// Initialize boxes
	teg.globals.bx[1000] = &box{kind: T_loc, skind: sub_region}
	teg.globals.bx[1001] = &box{kind: T_loc, skind: sub_forest}
	teg.globals.bx[1002] = &box{kind: T_loc, skind: sub_city}
	teg.globals.bx[1003] = &box{kind: T_loc, skind: sub_castle}
	teg.globals.bx[1004] = &box{kind: T_loc, skind: sub_graveyard}
	teg.globals.bx[2001] = &box{kind: T_char}
	teg.globals.bx[2002] = &box{kind: T_char}
	teg.globals.bx[2003] = &box{kind: T_char}

	// Set location hierarchy via x_loc_info.where
	teg.globals.bx[1000].x_loc_info.where = 0    // region has no parent
	teg.globals.bx[1001].x_loc_info.where = 1000 // province in region
	teg.globals.bx[1002].x_loc_info.where = 1001 // city in province
	teg.globals.bx[1003].x_loc_info.where = 1002 // castle in city
	teg.globals.bx[1004].x_loc_info.where = 1001 // graveyard in province
	teg.globals.bx[2001].x_loc_info.where = 1003 // char in castle
	teg.globals.bx[2002].x_loc_info.where = 1001 // char in province
	teg.globals.bx[2003].x_loc_info.where = 2002 // char stacked under char

	// Mark city as safe haven
	teg.globals.bx[1002].x_subloc = &entity_subloc{safe: 1}

	return func() {
		teg.globals.bx[1000] = nil
		teg.globals.bx[1001] = nil
		teg.globals.bx[1002] = nil
		teg.globals.bx[1003] = nil
		teg.globals.bx[1004] = nil
		teg.globals.bx[2001] = nil
		teg.globals.bx[2002] = nil
		teg.globals.bx[2003] = nil
	}
}

func TestLocDepth(t *testing.T) {
	cleanup := setupLocTestWorld(t)
	defer cleanup()

	tests := []struct {
		name     string
		id       int
		expected int
	}{
		{"region", 1000, LOC_region},
		{"province (forest)", 1001, LOC_province},
		{"subloc (city)", 1002, LOC_subloc},
		{"building (castle)", 1003, LOC_build},
		{"subloc (graveyard)", 1004, LOC_subloc},
		{"character (not a location)", 2001, 0},
		{"invalid box", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := loc_depth(tt.id)
			if got != tt.expected {
				t.Errorf("loc_depth(%d) = %d, want %d", tt.id, got, tt.expected)
			}
		})
	}
}

func TestRegion(t *testing.T) {
	cleanup := setupLocTestWorld(t)
	defer cleanup()

	tests := []struct {
		name     string
		who      int
		expected int
	}{
		{"from region", 1000, 1000},
		{"from province", 1001, 1000},
		{"from city", 1002, 1000},
		{"from castle", 1003, 1000},
		{"from char in castle", 2001, 1000},
		{"from char in province", 2002, 1000},
		{"from stacked char", 2003, 1000},
		{"invalid", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := region(tt.who)
			if got != tt.expected {
				t.Errorf("region(%d) = %d, want %d", tt.who, got, tt.expected)
			}
		})
	}
}

func TestProvince(t *testing.T) {
	cleanup := setupLocTestWorld(t)
	defer cleanup()

	tests := []struct {
		name     string
		who      int
		expected int
	}{
		{"from region (no province)", 1000, 0},
		{"from province", 1001, 1001},
		{"from city", 1002, 1001},
		{"from castle", 1003, 1001},
		{"from char in castle", 2001, 1001},
		{"from char in province", 2002, 1001},
		{"from stacked char", 2003, 1001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := province(tt.who)
			if got != tt.expected {
				t.Errorf("province(%d) = %d, want %d", tt.who, got, tt.expected)
			}
		})
	}
}

func TestSubloc(t *testing.T) {
	cleanup := setupLocTestWorld(t)
	defer cleanup()

	tests := []struct {
		name     string
		who      int
		expected int
	}{
		{"from char in castle", 2001, 1003},
		{"from char in province", 2002, 1001},
		{"from stacked char", 2003, 1001},
		{"from castle", 1003, 1002},
		{"from city", 1002, 1001},
		{"from province", 1001, 1000},
		{"from region", 1000, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := subloc(tt.who)
			if got != tt.expected {
				t.Errorf("subloc(%d) = %d, want %d", tt.who, got, tt.expected)
			}
		})
	}
}

func TestViewloc(t *testing.T) {
	cleanup := setupLocTestWorld(t)
	defer cleanup()

	tests := []struct {
		name     string
		who      int
		expected int
	}{
		{"from province", 1001, 1001},
		{"from city (stops at city)", 1002, 1002},
		{"from castle (walks to city)", 1003, 1002},
		{"from graveyard (stops)", 1004, 1004},
		{"from char in castle (walks to city)", 2001, 1002},
		{"from char in province", 2002, 1001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := viewloc(tt.who)
			if got != tt.expected {
				t.Errorf("viewloc(%d) = %d, want %d", tt.who, got, tt.expected)
			}
		})
	}
}

func TestInSafeNow(t *testing.T) {
	cleanup := setupLocTestWorld(t)
	defer cleanup()

	tests := []struct {
		name     string
		who      int
		expected bool
	}{
		{"in city (safe haven)", 1002, true},
		{"in castle inside safe city", 1003, true},
		{"char in castle inside safe city", 2001, true},
		{"in province (not safe)", 1001, false},
		{"char in province (not safe)", 2002, false},
		{"in region (not safe)", 1000, false},
		{"invalid", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := in_safe_now(tt.who)
			if got != tt.expected {
				t.Errorf("in_safe_now(%d) = %v, want %v", tt.who, got, tt.expected)
			}
		})
	}
}

func TestSomewhereInside(t *testing.T) {
	cleanup := setupLocTestWorld(t)
	defer cleanup()

	tests := []struct {
		name     string
		a, b     int
		expected bool
	}{
		{"char in castle", 1003, 2001, true},
		{"char in city (through castle)", 1002, 2001, true},
		{"char in province (through city, castle)", 1001, 2001, true},
		{"char in region (through all)", 1000, 2001, true},
		{"castle in city", 1002, 1003, true},
		{"same entity", 1002, 1002, false},
		{"not inside", 1003, 2002, false},
		{"sibling locations", 1002, 1004, false},
		{"reverse order", 2001, 1003, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := somewhere_inside(tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("somewhere_inside(%d, %d) = %v, want %v", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestLocDepthSubkinds(t *testing.T) {
	tests := []struct {
		subkind  schar
		expected int
	}{
		// Province level
		{sub_ocean, LOC_province},
		{sub_forest, LOC_province},
		{sub_plain, LOC_province},
		{sub_mountain, LOC_province},
		{sub_desert, LOC_province},
		{sub_swamp, LOC_province},
		{sub_under, LOC_province},
		{sub_cloud, LOC_province},
		{sub_tunnel, LOC_province},
		{sub_chamber, LOC_province},
		// Subloc level
		{sub_island, LOC_subloc},
		{sub_city, LOC_subloc},
		{sub_cave, LOC_subloc},
		{sub_graveyard, LOC_subloc},
		{sub_faery_hill, LOC_subloc},
		{sub_hades_pit, LOC_subloc},
		// Build level
		{sub_castle, LOC_build},
		{sub_temple, LOC_build},
		{sub_tower, LOC_build},
		{sub_inn, LOC_build},
		{sub_mine, LOC_build},
		{sub_sewer, LOC_build},
		{sub_galley, LOC_build},
		{sub_roundship, LOC_build},
	}

	for _, tt := range tests {
		// Create a temporary box with this subkind
		teg.globals.bx[9999] = &box{kind: T_loc, skind: tt.subkind}
		got := loc_depth(9999)
		if got != tt.expected {
			t.Errorf("loc_depth for subkind %d = %d, want %d", tt.subkind, got, tt.expected)
		}
		teg.globals.bx[9999] = nil
	}
}
