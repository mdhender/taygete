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

// format.go -- Time & numeric helpers from u.c
//
// This file ports formatting functions from u.c (Sprint 25.3):
//   - comma_num: format number with commas (e.g., "1,234,567")
//   - nice_num: convert 0-10 to words, else comma_num
//   - knum: compact notation (e.g., "1k", "1M")
//   - ordinal: ordinal suffix (e.g., "1st", "2nd", "3rd")
//   - weeks: format days as weeks/days
//   - more_weeks: format as "N more weeks/days"
//   - gold_s: format gold amount
//   - loyal_s: format loyalty type and level

import "fmt"

// num_s maps 0-10 to English words.
var num_s = []string{
	"zero", "one", "two", "three", "four", "five",
	"six", "seven", "eight", "nine", "ten",
}

// comma_num formats a number with commas.
// Port of C comma_num() from u.c.
func comma_num(n int) string {
	further := n / 1_000_000_000
	n = n % 1_000_000_000

	millions := n / 1_000_000
	n = n % 1_000_000

	thousands := n / 1_000
	ones := n % 1_000

	if further == 0 && millions == 0 && thousands == 0 {
		return fmt.Sprintf("%d", ones)
	} else if further == 0 && millions == 0 {
		return fmt.Sprintf("%d,%03d", thousands, ones)
	} else if further == 0 {
		return fmt.Sprintf("%d,%03d,%03d", millions, thousands, ones)
	}
	return fmt.Sprintf("%d,%03d,%03d,%03d", further, millions, thousands, ones)
}

// nice_num converts 0-10 to English words, otherwise uses comma_num.
// Port of C nice_num() from u.c.
// Note: This replaces the previous nice_num in code.go which only did comma formatting.
func nice_num_words(n int) string {
	if n >= 0 && n <= 10 {
		return num_s[n]
	}
	return comma_num(n)
}

// knum returns a compact representation of a number.
// If n is 0 and nozero is true, returns empty string.
// For n < 9999, returns the number as-is.
// For n < 1000000, returns "Nk".
// Otherwise returns "NM".
// Port of C knum() from u.c.
func knum(n int, nozero bool) string {
	if n == 0 && nozero {
		return ""
	}

	if n < 9999 {
		return fmt.Sprintf("%d", n)
	}

	if n < 1_000_000 {
		return fmt.Sprintf("%dk", n/1000)
	}

	return fmt.Sprintf("%dM", n/1_000_000)
}

// ordinal returns the ordinal representation of a number (e.g., "1st", "2nd", "3rd").
// Port of C ordinal() from u.c.
func ordinal(n int) string {
	// Special case: 10-19 all end in "th"
	if n >= 10 && n <= 19 {
		return fmt.Sprintf("%sth", comma_num(n))
	}

	switch n % 10 {
	case 1:
		return fmt.Sprintf("%sst", comma_num(n))
	case 2:
		return fmt.Sprintf("%snd", comma_num(n))
	case 3:
		return fmt.Sprintf("%srd", comma_num(n))
	default:
		return fmt.Sprintf("%sth", comma_num(n))
	}
}

// weeks formats a number of days as weeks or days.
// If n is divisible by 7, returns "N week(s)".
// Otherwise returns "N day(s)".
// Port of C weeks() from u.c.
func weeks(n int) string {
	if n == 0 {
		return "0~days"
	}

	if n%7 == 0 {
		w := n / 7
		return fmt.Sprintf("%s~week%s", nice_num_words(w), add_s(w))
	}

	return fmt.Sprintf("%s~day%s", nice_num_words(n), add_s(n))
}

// more_weeks formats a number of days as "N more weeks/days".
// Port of C more_weeks() from u.c.
func more_weeks(n int) string {
	if n == 0 {
		return "0~more days"
	}

	if n%7 == 0 {
		w := n / 7
		return fmt.Sprintf("%d~more week%s", w, add_s(w))
	}

	return fmt.Sprintf("%d~more day%s", n, add_s(n))
}

// gold_s formats a gold amount (e.g., "1,234~gold").
// Port of C gold_s() from u.c.
func gold_s(n int) string {
	return fmt.Sprintf("%s~gold", comma_num(n))
}

// loyal_s returns a string describing the loyalty type and level.
// Port of C loyal_s() from u.c.
func loyal_s(who int) string {
	var s string

	switch loyal_kind(who) {
	case 0:
		s = "unsworn"
	case LOY_contract:
		s = "contract"
	case LOY_oath:
		s = "oath"
	case LOY_fear:
		s = "fear"
	case LOY_npc:
		s = "npc"
	case LOY_summon:
		s = "summon"
	default:
		s = "unknown"
	}

	return fmt.Sprintf("%s-%d", s, loyal_rate(who))
}

// cap capitalizes the first character of a string.
// Port of C cap() from u.c.
// Note: This is an alias for cap_str which already exists in code.go.
func cap(s string) string {
	return cap_str(s)
}
