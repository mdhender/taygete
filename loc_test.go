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

// setupHereListTestWorld creates a test world with here_lists populated:
//
//	Province (1001) - sub_forest
//	  ├── City (1002) - sub_city
//	  │     ├── Castle (1003) - sub_castle
//	  │     │     └── Character (2001)
//	  │     └── Inn (1005) - sub_inn
//	  ├── Graveyard (1004) - sub_graveyard
//	  ├── Character (2002)
//	  │     └── Character (2003) - stacked under 2002
//	  └── Item (3001) - some item at province level
func setupHereListTestWorld(t *testing.T) func() {
	t.Helper()

	// Initialize boxes
	teg.globals.bx[1001] = &box{kind: T_loc, skind: sub_forest}
	teg.globals.bx[1002] = &box{kind: T_loc, skind: sub_city}
	teg.globals.bx[1003] = &box{kind: T_loc, skind: sub_castle}
	teg.globals.bx[1004] = &box{kind: T_loc, skind: sub_graveyard}
	teg.globals.bx[1005] = &box{kind: T_loc, skind: sub_inn}
	teg.globals.bx[2001] = &box{kind: T_char}
	teg.globals.bx[2002] = &box{kind: T_char}
	teg.globals.bx[2003] = &box{kind: T_char}
	teg.globals.bx[3001] = &box{kind: T_item}

	// Set location hierarchy via x_loc_info.where
	teg.globals.bx[1001].x_loc_info.where = 0    // province has no parent (for this test)
	teg.globals.bx[1002].x_loc_info.where = 1001 // city in province
	teg.globals.bx[1003].x_loc_info.where = 1002 // castle in city
	teg.globals.bx[1004].x_loc_info.where = 1001 // graveyard in province
	teg.globals.bx[1005].x_loc_info.where = 1002 // inn in city
	teg.globals.bx[2001].x_loc_info.where = 1003 // char in castle
	teg.globals.bx[2002].x_loc_info.where = 1001 // char in province
	teg.globals.bx[2003].x_loc_info.where = 2002 // char stacked under char
	teg.globals.bx[3001].x_loc_info.where = 1001 // item in province

	// Set here_lists
	teg.globals.bx[1001].x_loc_info.here_list = []int{1002, 1004, 2002, 3001}
	teg.globals.bx[1002].x_loc_info.here_list = []int{1003, 1005}
	teg.globals.bx[1003].x_loc_info.here_list = []int{2001}
	teg.globals.bx[1004].x_loc_info.here_list = []int{}
	teg.globals.bx[1005].x_loc_info.here_list = []int{}
	teg.globals.bx[2001].x_loc_info.here_list = []int{}
	teg.globals.bx[2002].x_loc_info.here_list = []int{2003}
	teg.globals.bx[2003].x_loc_info.here_list = []int{}
	teg.globals.bx[3001].x_loc_info.here_list = []int{}

	return func() {
		teg.globals.bx[1001] = nil
		teg.globals.bx[1002] = nil
		teg.globals.bx[1003] = nil
		teg.globals.bx[1004] = nil
		teg.globals.bx[1005] = nil
		teg.globals.bx[2001] = nil
		teg.globals.bx[2002] = nil
		teg.globals.bx[2003] = nil
		teg.globals.bx[3001] = nil
	}
}

func TestInHereList(t *testing.T) {
	cleanup := setupHereListTestWorld(t)
	defer cleanup()

	tests := []struct {
		name     string
		loc, who int
		expected bool
	}{
		{"city in province", 1001, 1002, true},
		{"graveyard in province", 1001, 1004, true},
		{"char in province", 1001, 2002, true},
		{"castle in city", 1002, 1003, true},
		{"char in castle", 1003, 2001, true},
		{"stacked char under char", 2002, 2003, true},
		{"castle not in province", 1001, 1003, false},
		{"char not in province", 1001, 2001, false},
		{"char not in castle", 1003, 2002, false},
		{"invalid loc", 9999, 2001, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := in_here_list(tt.loc, tt.who)
			if got != tt.expected {
				t.Errorf("in_here_list(%d, %d) = %v, want %v", tt.loc, tt.who, got, tt.expected)
			}
		})
	}
}

func TestAllHere(t *testing.T) {
	cleanup := setupHereListTestWorld(t)
	defer cleanup()

	var result []int

	// Test all_here from province - should get everything nested below
	all_here(1001, &result)

	// Expected: 1002, 1003, 2001, 1005, 1004, 2002, 2003, 3001
	// Order depends on recursion order
	expected := map[int]bool{
		1002: true, 1003: true, 1005: true, 1004: true,
		2001: true, 2002: true, 2003: true, 3001: true,
	}

	if len(result) != len(expected) {
		t.Errorf("all_here(1001) returned %d items, want %d", len(result), len(expected))
	}

	for _, id := range result {
		if !expected[id] {
			t.Errorf("all_here(1001) unexpectedly included %d", id)
		}
	}

	// Test all_here from city - should get castle, inn, and char in castle
	all_here(1002, &result)
	expectedCity := map[int]bool{1003: true, 1005: true, 2001: true}

	if len(result) != len(expectedCity) {
		t.Errorf("all_here(1002) returned %d items, want %d", len(result), len(expectedCity))
	}

	for _, id := range result {
		if !expectedCity[id] {
			t.Errorf("all_here(1002) unexpectedly included %d", id)
		}
	}

	// Test all_here from empty location (graveyard)
	all_here(1004, &result)
	if len(result) != 0 {
		t.Errorf("all_here(1004) returned %d items, want 0", len(result))
	}
}

func TestAllCharHere(t *testing.T) {
	cleanup := setupHereListTestWorld(t)
	defer cleanup()

	var result []int

	// Test all_char_here from province - should get only characters
	all_char_here(1001, &result)

	// Expected: 2002, 2003 (not 2001 because it's not directly in province's here_list)
	expected := map[int]bool{2002: true, 2003: true}

	if len(result) != len(expected) {
		t.Errorf("all_char_here(1001) returned %d items, want %d: got %v", len(result), len(expected), result)
	}

	for _, id := range result {
		if !expected[id] {
			t.Errorf("all_char_here(1001) unexpectedly included %d", id)
		}
	}

	// Test all_char_here from castle - should get char 2001
	all_char_here(1003, &result)
	if len(result) != 1 || result[0] != 2001 {
		t.Errorf("all_char_here(1003) = %v, want [2001]", result)
	}

	// Test all_char_here from empty graveyard
	all_char_here(1004, &result)
	if len(result) != 0 {
		t.Errorf("all_char_here(1004) returned %d items, want 0", len(result))
	}
}

func TestAllStack(t *testing.T) {
	cleanup := setupHereListTestWorld(t)
	defer cleanup()

	var result []int

	// Test all_stack from char 2002 - should get 2002 and 2003
	all_stack(2002, &result)

	if len(result) != 2 {
		t.Errorf("all_stack(2002) returned %d items, want 2: got %v", len(result), result)
	}
	if result[0] != 2002 {
		t.Errorf("all_stack(2002)[0] = %d, want 2002", result[0])
	}

	// Test all_stack from char with no stack
	all_stack(2001, &result)
	if len(result) != 1 || result[0] != 2001 {
		t.Errorf("all_stack(2001) = %v, want [2001]", result)
	}
}

func TestFirstCharacter(t *testing.T) {
	cleanup := setupHereListTestWorld(t)
	defer cleanup()

	tests := []struct {
		name     string
		where    int
		expected int
	}{
		{"castle with char", 1003, 2001},
		{"province (first char is 2002)", 1001, 2002},
		{"empty graveyard", 1004, 0},
		{"city (no direct chars)", 1002, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := first_character(tt.where)
			if got != tt.expected {
				t.Errorf("first_character(%d) = %d, want %d", tt.where, got, tt.expected)
			}
		})
	}
}

func TestSublocHere(t *testing.T) {
	cleanup := setupHereListTestWorld(t)
	defer cleanup()

	tests := []struct {
		name     string
		where    int
		sk       schar
		expected int
	}{
		{"city in province", 1001, sub_city, 1002},
		{"graveyard in province", 1001, sub_graveyard, 1004},
		{"castle in city", 1002, sub_castle, 1003},
		{"inn in city", 1002, sub_inn, 1005},
		{"no temple in province", 1001, sub_temple, 0},
		{"no castle in province", 1001, sub_castle, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := subloc_here(tt.where, tt.sk)
			if got != tt.expected {
				t.Errorf("subloc_here(%d, %d) = %d, want %d", tt.where, tt.sk, got, tt.expected)
			}
		})
	}
}

func TestCountLocStructures(t *testing.T) {
	cleanup := setupHereListTestWorld(t)
	defer cleanup()

	tests := []struct {
		name     string
		where    int
		a, b     schar
		expected int
	}{
		{"city and graveyard in province", 1001, sub_city, sub_graveyard, 2},
		{"only city in province", 1001, sub_city, sub_temple, 1},
		{"castle and inn in city", 1002, sub_castle, sub_inn, 2},
		{"no matching in graveyard", 1004, sub_castle, sub_inn, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := count_loc_structures(tt.where, tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("count_loc_structures(%d, %d, %d) = %d, want %d", tt.where, tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestCityHere(t *testing.T) {
	cleanup := setupHereListTestWorld(t)
	defer cleanup()

	got := city_here(1001)
	if got != 1002 {
		t.Errorf("city_here(1001) = %d, want 1002", got)
	}

	got = city_here(1002)
	if got != 0 {
		t.Errorf("city_here(1002) = %d, want 0", got)
	}
}

func TestBuildingOwner(t *testing.T) {
	cleanup := setupHereListTestWorld(t)
	defer cleanup()

	// Castle 1003 has char 2001 in it
	got := building_owner(1003)
	if got != 2001 {
		t.Errorf("building_owner(1003) = %d, want 2001", got)
	}

	// Inn 1005 is empty
	got = building_owner(1005)
	if got != 0 {
		t.Errorf("building_owner(1005) = %d, want 0", got)
	}
}

func TestAddRemoveHereList(t *testing.T) {
	cleanup := setupHereListTestWorld(t)
	defer cleanup()

	// Create a new character
	teg.globals.bx[2004] = &box{kind: T_char}
	defer func() { teg.globals.bx[2004] = nil }()

	// Add to graveyard's here_list
	add_to_here_list(1004, 2004)

	if !in_here_list(1004, 2004) {
		t.Error("add_to_here_list failed: 2004 not in 1004's here_list")
	}

	// Remove from here_list
	remove_from_here_list(1004, 2004)

	if in_here_list(1004, 2004) {
		t.Error("remove_from_here_list failed: 2004 still in 1004's here_list")
	}
}

func TestSetWhere(t *testing.T) {
	cleanup := setupHereListTestWorld(t)
	defer cleanup()

	// Create a new character in graveyard
	teg.globals.bx[2004] = &box{kind: T_char}
	teg.globals.bx[2004].x_loc_info.where = 1004
	teg.globals.bx[1004].x_loc_info.here_list = []int{2004}
	defer func() { teg.globals.bx[2004] = nil }()

	// Verify initial state
	if !in_here_list(1004, 2004) {
		t.Fatal("setup failed: 2004 not in graveyard's here_list")
	}
	if loc(2004) != 1004 {
		t.Fatalf("setup failed: loc(2004) = %d, want 1004", loc(2004))
	}

	// Move character from graveyard to castle
	set_where(2004, 1003)

	// Check old location
	if in_here_list(1004, 2004) {
		t.Error("set_where failed: 2004 still in graveyard's here_list")
	}

	// Check new location
	if !in_here_list(1003, 2004) {
		t.Error("set_where failed: 2004 not in castle's here_list")
	}

	// Check where field
	if loc(2004) != 1003 {
		t.Errorf("set_where failed: loc(2004) = %d, want 1003", loc(2004))
	}
}
