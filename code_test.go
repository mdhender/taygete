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

import (
	"testing"
)

func TestLetterVal(t *testing.T) {
	tests := []struct {
		c      byte
		let    string
		expect int
	}{
		{'a', letters, 0},
		{'b', letters, 1},
		{'z', letters, 25},
		{'a', letters2, 0},
		{'b', letters2, 1},
		{'c', letters2, 2},
		{'d', letters2, 3},
		{'f', letters2, 4}, // 'e' is skipped
		{'z', letters2, 19},
		{'x', letters2, 18}, // x is at index 18 in letters2
		{'e', letters2, 0},  // not in letters2, returns 0
	}

	for _, tc := range tests {
		got := letter_val(tc.c, tc.let)
		if got != tc.expect {
			t.Errorf("letter_val(%q, %q) = %d, want %d", tc.c, tc.let[:3]+"...", got, tc.expect)
		}
	}
}

func TestIntToCode(t *testing.T) {
	tests := []struct {
		input  int
		expect string
	}{
		// Simple numeric codes (items, chars, skills)
		{1, "1"},
		{100, "100"},
		{999, "999"},
		{1000, "1000"},
		{9999, "9999"},

		// CCNN format (provinces 10000-49999)
		// Format: letters2[b], letters2[a], digits n (2 digits)
		// offset = b*2000 + a*100 + n
		{10000, "aa00"}, // 0*2000 + 0*100 + 0 = 0
		{10001, "aa01"}, // 0*2000 + 0*100 + 1 = 1
		{10099, "aa99"}, // 0*2000 + 0*100 + 99 = 99
		{10100, "ab00"}, // 0*2000 + 1*100 + 0 = 100
		{10200, "ac00"}, // 0*2000 + 2*100 + 0 = 200
		{12000, "ba00"}, // 1*2000 + 0*100 + 0 = 2000
		{20000, "ga00"}, // 5*2000 + 0*100 + 0 = 10000
		{49999, "zz99"}, // 19*2000 + 19*100 + 99 = 39999

		// CCN format (players 50000-56759)
		{50000, "aa0"},
		{50001, "aa1"},
		{50010, "ab0"},
		{50260, "ba0"},
		{56759, "zz9"},

		// CNN format (lucky locs 56760-58759)
		// Format: letters2[b], digit a, digit n where offset = b*100 + a*10 + n
		{56760, "a00"}, // 0*100 + 0*10 + 0 = 0
		{56761, "a01"}, // 0*100 + 0*10 + 1 = 1
		{56770, "a10"}, // 0*100 + 1*10 + 0 = 10
		{56860, "b00"}, // 1*100 + 0*10 + 0 = 100
		{57560, "k00"}, // 8*100 + 0*10 + 0 = 800 (k is at index 8)
		{58759, "z99"}, // 19*100 + 9*10 + 9 = 1999 (z is at index 19)

		// Regions (58760-58999) - numeric
		{58760, "58760"},
		{58999, "58999"},

		// CNNN format (sublocs 59000-78999)
		// Format: letters2[a], digits n (3 digits) where offset = a*1000 + n
		{59000, "a000"}, // 0*1000 + 0 = 0
		{59001, "a001"}, // 0*1000 + 1 = 1
		{59999, "a999"}, // 0*1000 + 999 = 999
		{60000, "b000"}, // 1*1000 + 0 = 1000
		{67000, "k000"}, // 8*1000 + 0 = 8000 (k is at index 8)
		{78000, "z000"}, // 19*1000 + 0 = 19000 (z is at index 19)
		{78999, "z999"}, // 19*1000 + 999 = 19999

		// Storms (79000+) - numeric
		{79000, "79000"},
		{100000, "100000"},
	}

	for _, tc := range tests {
		got := int_to_code(tc.input)
		if got != tc.expect {
			t.Errorf("int_to_code(%d) = %q, want %q", tc.input, got, tc.expect)
		}
	}
}

func TestCodeToInt(t *testing.T) {
	tests := []struct {
		input  string
		expect int
	}{
		// Simple numeric codes
		{"1", 1},
		{"100", 100},
		{"999", 999},
		{"1000", 1000},
		{"9999", 9999},

		// CCNN format (provinces)
		{"aa00", 10000},
		{"aa01", 10001},
		{"aa99", 10099},
		{"ab00", 10100},
		{"ba00", 12000},
		{"zz99", 49999},

		// CCN format (players)
		{"aa0", 50000},
		{"aa1", 50001},
		{"ab0", 50010},
		{"ba0", 50260},
		{"zz9", 56759},

		// CNN format (lucky locs)
		{"a00", 56760},
		{"a01", 56761},
		{"a10", 56770},
		{"b00", 56860},
		{"k00", 57560},
		{"z99", 58759},

		// CNNN format (sublocs)
		{"a000", 59000},
		{"a001", 59001},
		{"a999", 59999},
		{"b000", 60000},
		{"k000", 67000},
		{"z000", 78000},
		{"z999", 78999},

		// Case insensitive
		{"AA00", 10000},
		{"Aa0", 50000},

		// Invalid codes
		{"", 0},
		{"@@@", 0},
		{"aaaaa", 0}, // too long for any format
	}

	for _, tc := range tests {
		got := code_to_int(tc.input)
		if got != tc.expect {
			t.Errorf("code_to_int(%q) = %d, want %d", tc.input, got, tc.expect)
		}
	}
}

func TestCodeRoundTrip(t *testing.T) {
	// Test that encoding then decoding gives the original value
	testValues := []int{
		// Items/chars/skills (numeric)
		1, 100, 999, 5000, 9999,
		// Provinces (CCNN)
		10000, 10001, 20000, 35000, 49999,
		// Players (CCN)
		50000, 50001, 53000, 56759,
		// Lucky locs (CNN)
		56760, 57000, 58000, 58759,
		// Sublocs (CNNN)
		59000, 59001, 65000, 78999,
	}

	for _, n := range testValues {
		code := int_to_code(n)
		back := code_to_int(code)
		if back != n {
			t.Errorf("round trip failed: %d -> %q -> %d", n, code, back)
		}
	}
}

func TestScode(t *testing.T) {
	tests := []struct {
		input  string
		expect int
	}{
		{"aa00", 10000},
		{"[aa00]", 10000},
		{"(aa00)", 10000},
		{"1000", 1000},
		{"[1000]", 1000},
		{"", 0},
	}

	for _, tc := range tests {
		got := scode(tc.input)
		if got != tc.expect {
			t.Errorf("scode(%q) = %d, want %d", tc.input, got, tc.expect)
		}
	}
}

func TestBoxCode(t *testing.T) {
	got := box_code(10000)
	want := "[aa00]"
	if got != want {
		t.Errorf("box_code(10000) = %q, want %q", got, want)
	}

	// Test garrison_magic special case
	got = box_code(teg.globals.garrison_magic)
	want = "Garrison"
	if got != want {
		t.Errorf("box_code(garrison_magic) = %q, want %q", got, want)
	}
}

func TestBoxCodeLess(t *testing.T) {
	got := box_code_less(10000)
	want := "aa00"
	if got != want {
		t.Errorf("box_code_less(10000) = %q, want %q", got, want)
	}
}

func TestNameSetName(t *testing.T) {
	// Allocate a test entity
	testID := 1000
	teg.globals.bx[testID] = &box{kind: T_char, skind: 0}

	// Initially no name
	if n := name(testID); n != "" {
		t.Errorf("name(%d) = %q, want empty string", testID, n)
	}

	// Set a name
	set_name(testID, "Test Character")
	if n := name(testID); n != "Test Character" {
		t.Errorf("name(%d) = %q, want %q", testID, n, "Test Character")
	}

	// Brackets should be replaced with braces
	set_name(testID, "Char [with] brackets")
	if n := name(testID); n != "Char {with} brackets" {
		t.Errorf("name(%d) = %q, want %q", testID, n, "Char {with} brackets")
	}

	// Clean up
	teg.globals.bx[testID] = nil
	delete(teg.globals.names, testID)
}

func TestNiceNum(t *testing.T) {
	tests := []struct {
		input  int
		expect string
	}{
		{0, "0"},
		{1, "1"},
		{999, "999"},
		{1000, "1,000"},
		{1234, "1,234"},
		{12345, "12,345"},
		{123456, "123,456"},
		{1234567, "1,234,567"},
	}

	for _, tc := range tests {
		got := nice_num(tc.input)
		if got != tc.expect {
			t.Errorf("nice_num(%d) = %q, want %q", tc.input, got, tc.expect)
		}
	}
}

func TestCapStr(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"", ""},
		{"hello", "Hello"},
		{"HELLO", "HELLO"},
		{"hELLO", "HELLO"},
		{"castle", "Castle"},
	}

	for _, tc := range tests {
		got := cap_str(tc.input)
		if got != tc.expect {
			t.Errorf("cap_str(%q) = %q, want %q", tc.input, got, tc.expect)
		}
	}
}

func TestAllocBoxDeleteBox(t *testing.T) {
	testID := 1001

	// Ensure clean state
	teg.globals.bx[testID] = nil

	// Allocate
	alloc_box(testID, T_char, 0)
	if teg.globals.bx[testID] == nil {
		t.Fatalf("alloc_box(%d) did not create box", testID)
	}
	if kind(testID) != T_char {
		t.Errorf("kind(%d) = %d, want %d", testID, kind(testID), T_char)
	}

	// Delete
	delete_box(testID)
	if kind(testID) != T_deleted {
		t.Errorf("after delete, kind(%d) = %d, want %d", testID, kind(testID), T_deleted)
	}

	// Clean up
	teg.globals.bx[testID] = nil
}

func TestChangeBoxKind(t *testing.T) {
	testID := 1002

	// Allocate as char
	teg.globals.bx[testID] = nil
	alloc_box(testID, T_char, 0)

	// Change to player
	change_box_kind(testID, T_player)
	if kind(testID) != T_player {
		t.Errorf("after change, kind(%d) = %d, want %d", testID, kind(testID), T_player)
	}

	// Clean up
	teg.globals.bx[testID] = nil
}

func TestChangeBoxSubkind(t *testing.T) {
	testID := 1003

	// Allocate
	teg.globals.bx[testID] = nil
	alloc_box(testID, T_loc, sub_forest)

	// Change subkind
	change_box_subkind(testID, sub_mountain)
	if subkind(testID) != sub_mountain {
		t.Errorf("after change, subkind(%d) = %d, want %d", testID, subkind(testID), sub_mountain)
	}

	// No change if same
	change_box_subkind(testID, sub_mountain)
	if subkind(testID) != sub_mountain {
		t.Errorf("after no-op change, subkind(%d) = %d, want %d", testID, subkind(testID), sub_mountain)
	}

	// Clean up
	teg.globals.bx[testID] = nil
}
