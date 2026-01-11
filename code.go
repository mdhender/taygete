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
	"fmt"
	"strings"
	"unicode"
)

// code.go ports code.c - ID encoding/decoding and naming functions

// G2 Entity coding system:
//
//  range          extent    use
//  1-999             999    items
//  1000-8999        8000    chars
//  9000-9999        1000    skills
//
//  10,000-19,999   10,000   provinces       (CCNN: AA00-DV99)
//  20,000-49,999   30,000   more provinces  (CCNN: DW00-ZZ99)
//  50,000-56,759    6,760   player entities (CCN)
//  56,760-58,759    2,000   lucky locs      (CNN)
//  58,760-58,999      240   regions         (NNNNN)
//  59,000-78,999   20,000   sublocs, misc   (CNNN: A000-Z999)
//  79,000-100,000  21,000   storms          (NNNNN)
//
//  Note: restricted alphabet, no vowels (except a) or l:
//      "abcdfghjkmnpqrstvwxz"

// letters is the full 26-letter alphabet
const letters = "abcdefghijklmnopqrstuvwxyz"

// letters2 is the restricted alphabet (no vowels except 'a', no 'l')
const letters2 = "abcdfghjkmnpqrstvwxz"

// letter_val returns the index of character c in the letter string let.
// Returns 0 if not found (error case).
func letter_val(c byte, let string) int {
	for i := 0; i < len(let); i++ {
		if let[i] == c {
			return i
		}
	}
	return 0 // error
}

// int_to_code converts an integer entity ID to its alphanumeric code string.
func int_to_code(l int) string {
	if l < 10_000 {
		return fmt.Sprintf("%d", l)
	}

	if l < 50_000 { // CCNN
		l -= 10_000

		n := l % 100
		l /= 100

		a := l % 20
		b := l / 20

		return fmt.Sprintf("%c%c%02d", letters2[b], letters2[a], n)
	}

	if l < 56_760 { // CCN
		l -= 50_000

		n := l % 10
		l /= 10

		a := l % 26
		b := l / 26

		return fmt.Sprintf("%c%c%d", letters[b], letters[a], n)
	}

	if l < 58_760 { // CNN
		l -= 56_760

		n := l % 10
		l /= 10

		a := l % 10
		b := l / 10

		return fmt.Sprintf("%c%d%d", letters2[b], a, n)
	}

	if l < 59_000 {
		return fmt.Sprintf("%d", l)
	}

	if l < 79_000 { // CNNN
		l -= 59_000

		a := l / 1_000
		n := l % 1_000

		return fmt.Sprintf("%c%03d", letters2[a], n)
	}

	return fmt.Sprintf("%d", l)
}

// code_to_int converts an alphanumeric code string to its integer entity ID.
// Returns 0 if the string is not a valid code.
func code_to_int(s string) int {
	if len(s) == 0 {
		return 0
	}

	// If starts with digit, parse as integer
	if s[0] >= '0' && s[0] <= '9' {
		var n int
		_, err := fmt.Sscanf(s, "%d", &n)
		if err != nil {
			logger.Error("code_to_int", "sscanf", err)
		}
		return n
	}

	// Must start with a letter
	if !isAlpha(s[0]) {
		return 0
	}

	switch len(s) {
	case 3:
		if isAlpha(s[1]) && isDigit(s[2]) { // CCN
			a := int(toLower(s[0]) - 'a')
			b := int(toLower(s[1]) - 'a')
			c := int(s[2] - '0')

			return a*260 + b*10 + c + 50_000
		}

		if isDigit(s[1]) && isDigit(s[2]) { // CNN
			a := letter_val(toLower(s[0]), letters2)
			b := int(s[1] - '0')
			c := int(s[2] - '0')

			return a*100 + b*10 + c + 56_760
		}

		return 0

	case 4:
		if isAlpha(s[1]) && isDigit(s[2]) && isDigit(s[3]) { // CCNN
			a := letter_val(toLower(s[0]), letters2)
			b := letter_val(toLower(s[1]), letters2)
			c := int(s[2] - '0')
			d := int(s[3] - '0')

			return a*2_000 + b*100 + c*10 + d + 10_000
		}

		if isDigit(s[1]) && isDigit(s[2]) && isDigit(s[3]) { // CNNN
			a := letter_val(toLower(s[0]), letters2)
			b := int(s[1] - '0')
			c := int(s[2] - '0')
			d := int(s[3] - '0')

			return a*1_000 + b*100 + c*10 + d + 59_000
		}
		return 0

	default:
		return 0
	}
}

// scode parses a code string, stripping leading and trailing brackets if present.
func scode(s string) int {
	if len(s) == 0 {
		return 0
	}
	// Strip leading bracket
	if s[0] == '[' || s[0] == '(' {
		s = s[1:]
	}
	// Strip trailing bracket
	if len(s) > 0 && (s[len(s)-1] == ']' || s[len(s)-1] == ')') {
		s = s[:len(s)-1]
	}
	return code_to_int(s)
}

// name returns the name of entity n.
// Returns empty string if entity doesn't exist or has no name.
func name(n int) string {
	if !valid_box(n) {
		return ""
	}
	return teg.getName(n)
}

// set_name sets the name of entity n.
// Replaces brackets with braces in the name.
func set_name(n int, s string) {
	if !valid_box(n) {
		return
	}
	// Replace brackets with braces (as in C code)
	s = strings.ReplaceAll(s, "[", "{")
	s = strings.ReplaceAll(s, "]", "}")
	teg.setName(n, s)
}

// set_banner sets the display banner for entity n.
// Truncates to 50 characters max.
func set_banner(n int, s string) {
	if s != "" && len(s) > 50 {
		s = s[:50]
	}
	teg.setBanner(n, s)
}

// banner returns the display banner for entity n.
func get_banner(n int) string {
	return teg.getBanner(n)
}

// display_name returns the display name for entity n.
// Falls back to kind-based names if entity has no name.
func display_name(n int) string {
	if !valid_box(n) {
		return ""
	}

	s := name(n)
	if s != "" {
		return s
	}

	switch kind(n) {
	case T_player:
		return "Player"
	case T_gate:
		return "Gate"
	case T_post:
		return "Sign"
	}

	if i := noble_item(n); i != 0 {
		return cap_str(plural_item_name(int(i), 1))
	}

	sk := subkind(n)
	if int(sk) < len(subkind_s) {
		return cap_str(subkind_s[sk])
	}
	return ""
}

// display_kind returns the kind description for entity n.
func display_kind(n int) string {
	sk := subkind(n)

	switch sk {
	case sub_city:
		if is_port_city(n) {
			return "port city"
		}
		return "city"

	case sub_fog, sub_rain, sub_wind:
		return "storm"

	default:
		if int(sk) < len(subkind_s) {
			return subkind_s[sk]
		}
		return ""
	}
}

// box_code_less returns the code for entity n without brackets.
func box_code_less(n int) string {
	return int_to_code(n)
}

// box_code returns the code for entity n with brackets.
func box_code(n int) string {
	if n == teg.globals.garrison_magic {
		return "Garrison"
	}
	return fmt.Sprintf("[%s]", int_to_code(n))
}

// box_name returns the display name and code for entity n.
func box_name(n int) string {
	if n == teg.globals.garrison_magic {
		return "Garrison"
	}

	if valid_box(n) {
		s := display_name(n)
		if s != "" {
			return fmt.Sprintf("%s~%s", s, box_code(n))
		}
	}

	return box_code(n)
}

// just_name returns just the display name, or the code if no name.
func just_name(n int) string {
	if n == teg.globals.garrison_magic {
		return "Garrison"
	}

	if valid_box(n) {
		s := display_name(n)
		if s != "" {
			return s
		}
	}

	return box_code(n)
}

// plural_item_name returns the plural name for an item.
func plural_item_name(item, qty int) string {
	if qty == 1 {
		return display_name(item)
	}

	s := teg.getPluralName(item)
	if s != "" {
		return s
	}

	// Fall back to singular name
	return display_name(item)
}

// plural_item_box returns the plural name with box code.
func plural_item_box(item, qty int) string {
	if qty == 1 {
		return box_name(item)
	}

	s := plural_item_name(item, qty)
	return fmt.Sprintf("%s~%s", s, box_code(item))
}

// just_name_qty returns quantity and plural item name.
func just_name_qty(item, qty int) string {
	return fmt.Sprintf("%s~%s", nice_num(qty), plural_item_name(item, qty))
}

// box_name_qty returns quantity and plural item box name.
func box_name_qty(item, qty int) string {
	return fmt.Sprintf("%s~%s", nice_num(qty), plural_item_box(item, qty))
}

// box_name_kind returns box name with kind description.
func box_name_kind(n int) string {
	return fmt.Sprintf("%s, %s", box_name(n), display_kind(n))
}

// Entity allocation and chain management functions

// add_next_chain adds entity n to the kind chain.
func add_next_chain(n int) {
	if teg.globals.bx[n] == nil {
		return
	}
	k := int(teg.globals.bx[n].kind)
	if k == 0 {
		return
	}

	// Find insertion point (keep sorted by ID)
	if teg.globals.box_head[k] == 0 {
		teg.globals.box_head[k] = n
		teg.globals.bx[n].x_next_kind = 0
		return
	}

	if n < teg.globals.box_head[k] {
		teg.globals.bx[n].x_next_kind = teg.globals.box_head[k]
		teg.globals.box_head[k] = n
		return
	}

	i := teg.globals.box_head[k]
	for teg.globals.bx[i].x_next_kind > 0 && teg.globals.bx[i].x_next_kind < n {
		i = teg.globals.bx[i].x_next_kind
	}

	teg.globals.bx[n].x_next_kind = teg.globals.bx[i].x_next_kind
	teg.globals.bx[i].x_next_kind = n
}

// remove_next_chain removes entity n from the kind chain.
func remove_next_chain(n int) {
	if teg.globals.bx[n] == nil {
		return
	}

	k := int(teg.globals.bx[n].kind)
	i := teg.globals.box_head[k]

	if i == n {
		teg.globals.box_head[k] = teg.globals.bx[n].x_next_kind
	} else {
		for i > 0 && teg.globals.bx[i].x_next_kind != n {
			i = teg.globals.bx[i].x_next_kind
		}
		if i > 0 {
			teg.globals.bx[i].x_next_kind = teg.globals.bx[n].x_next_kind
		}
	}

	teg.globals.bx[n].x_next_kind = 0
}

// add_sub_chain adds entity n to the subkind chain.
func add_sub_chain(n int) {
	if teg.globals.bx[n] == nil {
		return
	}
	sk := int(teg.globals.bx[n].skind)

	// Find insertion point (keep sorted by ID)
	if teg.globals.sub_head[sk] == 0 {
		teg.globals.sub_head[sk] = n
		teg.globals.bx[n].x_next_sub = 0
		return
	}

	if n < teg.globals.sub_head[sk] {
		teg.globals.bx[n].x_next_sub = teg.globals.sub_head[sk]
		teg.globals.sub_head[sk] = n
		return
	}

	i := teg.globals.sub_head[sk]
	for teg.globals.bx[i].x_next_sub > 0 && teg.globals.bx[i].x_next_sub < n {
		i = teg.globals.bx[i].x_next_sub
	}

	teg.globals.bx[n].x_next_sub = teg.globals.bx[i].x_next_sub
	teg.globals.bx[i].x_next_sub = n
}

// remove_sub_chain removes entity n from the subkind chain.
func remove_sub_chain(n int) {
	if teg.globals.bx[n] == nil {
		return
	}

	sk := int(teg.globals.bx[n].skind)
	i := teg.globals.sub_head[sk]

	if i == n {
		teg.globals.sub_head[sk] = teg.globals.bx[n].x_next_sub
	} else {
		for i > 0 && teg.globals.bx[i].x_next_sub != n {
			i = teg.globals.bx[i].x_next_sub
		}
		if i > 0 {
			teg.globals.bx[i].x_next_sub = teg.globals.bx[n].x_next_sub
		}
	}

	teg.globals.bx[n].x_next_sub = 0
}

// delete_box marks an entity as deleted.
func delete_box(n int) {
	remove_next_chain(n)
	remove_sub_chain(n)
	teg.globals.bx[n].kind = T_deleted
}

// change_box_kind changes the kind of entity n.
func change_box_kind(n int, k schar) {
	remove_next_chain(n)
	teg.globals.bx[n].kind = k
	add_next_chain(n)
}

// change_box_subkind changes the subkind of entity n.
func change_box_subkind(n int, sk schar) {
	if subkind(n) == sk {
		return
	}
	remove_sub_chain(n)
	teg.globals.bx[n].skind = sk
	add_sub_chain(n)
}

// alloc_box allocates a new box at position n with given kind and subkind.
func alloc_box(n int, k, sk schar) {
	if n <= 0 || n >= MAX_BOXES {
		logger.Error("alloc_box", "invalid box id", n)
		panic(fmt.Sprintf("alloc_box: invalid box ID %d", n))
	}
	if teg.globals.bx[n] != nil {
		logger.Error("alloc_box", "duplicate box id", n)
		panic(fmt.Sprintf("alloc_box: DUP box %d", n))
	}

	teg.globals.bx[n] = &box{
		kind:  k,
		skind: sk,
	}
	add_next_chain(n)
	add_sub_chain(n)
}

// rnd_alloc_num finds a random unallocated box ID in range [low, high].
// Returns -1 if no free slot is found.
func rnd_alloc_num(low, high int) int {
	n := rnd(low, high)

	// Search from n to high
	for i := n; i <= high; i++ {
		if teg.globals.bx[i] == nil {
			return i
		}
	}

	// Search from low to n-1
	for i := low; i < n; i++ {
		if teg.globals.bx[i] == nil {
			return i
		}
	}

	return -1
}

// new_ent allocates a new entity of the given kind and subkind.
// Returns the entity ID, or panics if no space available.
func new_ent(k, sk schar) int {
	var n int = -1

	switch k {
	case T_player:
		n = rnd_alloc_num(50_000, 56_759)

	case T_char, T_unform:
		if sk == sub_ni {
			n = rnd_alloc_num(79_000, MAX_BOXES-1)
		} else {
			n = rnd_alloc_num(1_000, 9_999)
		}
		if n < 0 {
			n = rnd_alloc_num(59_000, 78_999)
		}
		if n < 0 {
			n = rnd_alloc_num(79_000, MAX_BOXES-1)
		}

	case T_skill:
		panic("new_ent: cannot allocate skill entities this way")

	case T_loc:
		switch sk {
		case sub_city:
			n = rnd_alloc_num(56_760, 58_759)
		case sub_region:
			n = rnd_alloc_num(58_760, 58_999)
		case sub_under, sub_forest, sub_ocean, sub_cloud, sub_tunnel:
			n = rnd_alloc_num(20_000, 49_999)
		default:
			n = rnd_alloc_num(59_000, 78_999)
		}

	case T_storm:
		n = rnd_alloc_num(79_000, MAX_BOXES-1)

	default:
		n = rnd_alloc_num(59_000, 78_999)
	}

	if n < 0 {
		n = rnd_alloc_num(79_000, MAX_BOXES-1)
	}

	if n < 0 {
		panic("new_ent: out of entity space")
	}

	alloc_box(n, k, sk)
	return n
}

// print_box_usage prints entity space usage statistics.
func print_box_usage() {
	fmt.Println("entity space usage:")
	print_box_usage_sup(1_000, 9_999, "char")
	print_box_usage_sup(20_000, 59_999, "locs")
	print_box_usage_sup(50_000, 56_759, "play")
	print_box_usage_sup(56_760, 58_759, "lloc")
	print_box_usage_sup(59_000, 78_999, "subl")
	print_box_usage_sup(79_000, MAX_BOXES-1, "rest")
}

// print_box_usage_sup prints usage for a specific range.
func print_box_usage_sup(low, high int, label string) {
	used := 0
	for i := low; i <= high; i++ {
		if teg.globals.bx[i] != nil {
			used++
		}
	}
	total := high - low
	pct := 0
	if total > 0 {
		pct = used * 100 / total
	}
	fmt.Printf("\t%5d - %5d  %4s  %d/%d used (%d%%)\n",
		low, high, label, used, total, pct)
}

// Helper functions

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 32
	}
	return c
}

// cap_str capitalizes the first letter of a string.
func cap_str(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// nice_num formats a number with commas.
func nice_num(n int) string {
	if n < 1_000 {
		return fmt.Sprintf("%d", n)
	}
	// Simple comma formatting
	s := fmt.Sprintf("%d", n)
	result := ""
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += ","
		}
		result += string(c)
	}
	return result
}

// is_port_city returns true if the location is a port city.
// Stub implementation - will be completed when loc.c is ported.
func is_port_city(n int) bool {
	// TODO: implement when loc.c is ported
	logger.Error("is_port_city",
		"todo", "implement when `loc.c` is ported",
	)
	return false
}
