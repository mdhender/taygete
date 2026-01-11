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

func TestCommaNum(t *testing.T) {
	tests := []struct {
		input  int
		expect string
	}{
		{0, "0"},
		{1, "1"},
		{12, "12"},
		{123, "123"},
		{999, "999"},
		{1000, "1,000"},
		{1234, "1,234"},
		{12345, "12,345"},
		{123456, "123,456"},
		{1234567, "1,234,567"},
		{12345678, "12,345,678"},
		{123456789, "123,456,789"},
		{1234567890, "1,234,567,890"},
	}

	for _, tc := range tests {
		got := comma_num(tc.input)
		if got != tc.expect {
			t.Errorf("comma_num(%d) = %q, want %q", tc.input, got, tc.expect)
		}
	}
}

func TestKnum(t *testing.T) {
	tests := []struct {
		n      int
		nozero bool
		expect string
	}{
		{0, false, "0"},
		{0, true, ""},
		{1, false, "1"},
		{999, false, "999"},
		{9998, false, "9998"},
		{9999, false, "9k"}, // C: n < 9999 returns as-is; 9999 -> 9k
		{10000, false, "10k"},
		{50000, false, "50k"},
		{999999, false, "999k"},
		{1000000, false, "1M"},
		{5000000, false, "5M"},
	}

	for _, tc := range tests {
		got := knum(tc.n, tc.nozero)
		if got != tc.expect {
			t.Errorf("knum(%d, %v) = %q, want %q", tc.n, tc.nozero, got, tc.expect)
		}
	}
}

func TestOrdinal(t *testing.T) {
	tests := []struct {
		input  int
		expect string
	}{
		{0, "0th"},
		{1, "1st"},
		{2, "2nd"},
		{3, "3rd"},
		{4, "4th"},
		{5, "5th"},
		{10, "10th"},
		{11, "11th"},
		{12, "12th"},
		{13, "13th"},
		{14, "14th"},
		{19, "19th"},
		{20, "20th"},
		{21, "21st"},
		{22, "22nd"},
		{23, "23rd"},
		{24, "24th"},
		{100, "100th"},
		{101, "101st"},
		{102, "102nd"},
		{103, "103rd"},
		{111, "111st"}, // C only special-cases 10-19, not 111-119
		{112, "112nd"},
		{113, "113rd"},
		{1000, "1,000th"},
		{1001, "1,001st"},
	}

	for _, tc := range tests {
		got := ordinal(tc.input)
		if got != tc.expect {
			t.Errorf("ordinal(%d) = %q, want %q", tc.input, got, tc.expect)
		}
	}
}

func TestWeeks(t *testing.T) {
	tests := []struct {
		input  int
		expect string
	}{
		{0, "0~days"},
		{1, "one~day"},
		{2, "two~days"},
		{6, "six~days"},
		{7, "one~week"},
		{14, "two~weeks"},
		{21, "three~weeks"},
		{8, "eight~days"},
		{10, "ten~days"},
		{11, "11~days"},
		{28, "four~weeks"},
	}

	for _, tc := range tests {
		got := weeks(tc.input)
		if got != tc.expect {
			t.Errorf("weeks(%d) = %q, want %q", tc.input, got, tc.expect)
		}
	}
}

func TestMoreWeeks(t *testing.T) {
	tests := []struct {
		input  int
		expect string
	}{
		{0, "0~more days"},
		{1, "1~more day"},
		{2, "2~more days"},
		{6, "6~more days"},
		{7, "1~more week"},
		{14, "2~more weeks"},
		{21, "3~more weeks"},
		{8, "8~more days"},
	}

	for _, tc := range tests {
		got := more_weeks(tc.input)
		if got != tc.expect {
			t.Errorf("more_weeks(%d) = %q, want %q", tc.input, got, tc.expect)
		}
	}
}

func TestGoldS(t *testing.T) {
	tests := []struct {
		input  int
		expect string
	}{
		{0, "0~gold"},
		{1, "1~gold"},
		{100, "100~gold"},
		{1000, "1,000~gold"},
		{1234567, "1,234,567~gold"},
	}

	for _, tc := range tests {
		got := gold_s(tc.input)
		if got != tc.expect {
			t.Errorf("gold_s(%d) = %q, want %q", tc.input, got, tc.expect)
		}
	}
}

func TestLoyalS(t *testing.T) {
	e := newTestEngine(t)

	charID := 1001

	tests := []struct {
		name     string
		loyKind  schar
		loyRate  int
		expected string
	}{
		{"unsworn", 0, 0, "unsworn-0"},
		{"contract", LOY_contract, 50, "contract-50"},
		{"oath", LOY_oath, 75, "oath-75"},
		{"fear", LOY_fear, 30, "fear-30"},
		{"npc", LOY_npc, 100, "npc-100"},
		{"summon", LOY_summon, 10, "summon-10"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up character with loyalty
			e.globals.bx[charID] = &box{kind: T_char}
			e.globals.bx[charID].x_char = &entity_char{
				loy_kind: tc.loyKind,
				loy_rate: tc.loyRate,
			}

			got := loyal_s(charID)
			if got != tc.expected {
				t.Errorf("loyal_s() = %q, want %q", got, tc.expected)
			}

			// Clean up
			e.globals.bx[charID] = nil
		})
	}
}

func TestCap(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"", ""},
		{"hello", "Hello"},
		{"HELLO", "HELLO"},
		{"hello world", "Hello world"},
		{"one", "One"},
	}

	for _, tc := range tests {
		got := cap(tc.input)
		if got != tc.expect {
			t.Errorf("cap(%q) = %q, want %q", tc.input, got, tc.expect)
		}
	}
}
