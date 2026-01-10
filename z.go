// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// IListLen returns the length of the integer list.
func IListLen(l []int) int {
	return len(l)
}

// IListAppend appends n to the end of the list.
func IListAppend(l *[]int, n int) {
	*l = append(*l, n)
}

// IListPrepend prepends n to the beginning of the list.
func IListPrepend(l *[]int, n int) {
	*l = append([]int{n}, *l...)
}

// IListDelete removes the element at index i.
func IListDelete(l *[]int, i int) {
	if i < 0 || i >= len(*l) {
		panic(fmt.Sprintf("ilist_delete: index %d out of bounds [0, %d)", i, len(*l)))
	}
	*l = append((*l)[:i], (*l)[i+1:]...)
}

// IListClear clears the list but retains the underlying array.
func IListClear(l *[]int) {
	*l = (*l)[:0]
}

// IListReclaim releases the list memory.
func IListReclaim(l *[]int) {
	*l = nil
}

// IListLookup returns the index of n in the list, or -1 if not found.
func IListLookup(l []int, n int) int {
	for i, v := range l {
		if v == n {
			return i
		}
	}
	return -1
}

// IListRemValue removes all occurrences of n from the list.
func IListRemValue(l *[]int, n int) {
	for i := len(*l) - 1; i >= 0; i-- {
		if (*l)[i] == n {
			*l = append((*l)[:i], (*l)[i+1:]...)
		}
	}
}

// IListRemValueUniq removes the first occurrence of n from the list.
func IListRemValueUniq(l *[]int, n int) {
	for i := len(*l) - 1; i >= 0; i-- {
		if (*l)[i] == n {
			*l = append((*l)[:i], (*l)[i+1:]...)
			break
		}
	}
}

// IListCopy returns a copy of the list.
func IListCopy(l []int) []int {
	if l == nil {
		return nil
	}
	cp := make([]int, len(l))
	copy(cp, l)
	return cp
}

// IListScramble performs a Fisher-Yates shuffle on the list.
func IListScramble(l []int) {
	n := len(l) - 1
	for i := 0; i < n; i++ {
		r := rnd(i, n)
		if r != i {
			l[i], l[r] = l[r], l[i]
		}
	}
}

// IListInsert inserts n at position pos in the list.
func IListInsert(l *[]int, pos, n int) {
	*l = append(*l, 0)
	copy((*l)[pos+1:], (*l)[pos:])
	(*l)[pos] = n
}

// IListSort sorts the list in ascending order.
func IListSort(l []int) {
	sort.Ints(l)
}

// PListLen returns the length of the pointer list.
func PListLen(l []any) int {
	return len(l)
}

// PListAppend appends n to the end of the list.
func PListAppend(l *[]any, n any) {
	*l = append(*l, n)
}

// PListPrepend prepends n to the beginning of the list.
func PListPrepend(l *[]any, n any) {
	*l = append([]any{n}, *l...)
}

// PListDelete removes the element at index i.
func PListDelete(l *[]any, i int) {
	if i < 0 || i >= len(*l) {
		panic(fmt.Sprintf("plist_delete: index %d out of bounds [0, %d)", i, len(*l)))
	}
	*l = append((*l)[:i], (*l)[i+1:]...)
}

// PListClear clears the list but retains the underlying array.
func PListClear(l *[]any) {
	*l = (*l)[:0]
}

// PListReclaim releases the list memory.
func PListReclaim(l *[]any) {
	*l = nil
}

// PListLookup returns the index of n in the list, or -1 if not found.
func PListLookup(l []any, n any) int {
	for i, v := range l {
		if v == n {
			return i
		}
	}
	return -1
}

// PListRemValue removes all occurrences of n from the list.
func PListRemValue(l *[]any, n any) {
	for i := len(*l) - 1; i >= 0; i-- {
		if (*l)[i] == n {
			*l = append((*l)[:i], (*l)[i+1:]...)
		}
	}
}

// PListRemValueUniq removes the first occurrence of n from the list.
func PListRemValueUniq(l *[]any, n any) {
	for i := len(*l) - 1; i >= 0; i-- {
		if (*l)[i] == n {
			*l = append((*l)[:i], (*l)[i+1:]...)
			break
		}
	}
}

// PListCopy returns a copy of the list.
func PListCopy(l []any) []any {
	if l == nil {
		return nil
	}
	cp := make([]any, len(l))
	copy(cp, l)
	return cp
}

// PListScramble performs a Fisher-Yates shuffle on the list.
func PListScramble(l []any) {
	n := len(l) - 1
	for i := 0; i < n; i++ {
		r := rnd(i, n)
		if r != i {
			l[i], l[r] = l[r], l[i]
		}
	}
}

// PListInsert inserts n at position pos in the list.
func PListInsert(l *[]any, pos int, n any) {
	*l = append(*l, nil)
	copy((*l)[pos+1:], (*l)[pos:])
	(*l)[pos] = n
}

// i_strcmp performs a case-insensitive string comparison.
// Returns 0 if equal, negative if s < t, positive if s > t.
func i_strcmp(s, t string) int {
	sl := strings.ToLower(s)
	tl := strings.ToLower(t)
	if sl < tl {
		return -1
	}
	if sl > tl {
		return 1
	}
	return 0
}

// i_strncmp performs a case-insensitive comparison of the first n characters.
// Returns 0 if equal, negative if s < t, positive if s > t.
func i_strncmp(s, t string, n int) int {
	if n <= 0 {
		return 0
	}
	if len(s) > n {
		s = s[:n]
	}
	if len(t) > n {
		t = t[:n]
	}
	return i_strcmp(s, t)
}

// fuzzy_transpose checks if one and two differ only by adjacent transposed characters.
func fuzzy_transpose(one, two string, l1, l2 int) bool {
	if l1 != l2 {
		return false
	}
	buf := []byte(two)
	for i := 0; i < l2-1; i++ {
		buf[i], buf[i+1] = buf[i+1], buf[i]
		if i_strcmp(one, string(buf)) == 0 {
			return true
		}
		buf[i], buf[i+1] = buf[i+1], buf[i]
	}
	return false
}

// fuzzy_one_less checks if one has exactly one extra character compared to two.
// one is longer than two by 1 character.
func fuzzy_one_less(one, two string, l1, l2 int) bool {
	if l1 != l2+1 {
		return false
	}
	count := 0
	i, j := 0, 0
	for j < l2 {
		if toLowerByte(one[i]) != toLowerByte(two[j]) {
			count++
			if count > 1 {
				return false
			}
			i++
		} else {
			i++
			j++
		}
	}
	return true
}

// fuzzy_one_extra checks if one has exactly one fewer character compared to two.
// one is shorter than two by 1 character.
func fuzzy_one_extra(one, two string, l1, l2 int) bool {
	if l1 != l2-1 {
		return false
	}
	count := 0
	i, j := 0, 0
	for i < l1 {
		if toLowerByte(one[i]) != toLowerByte(two[j]) {
			count++
			if count > 1 {
				return false
			}
			j++
		} else {
			i++
			j++
		}
	}
	return true
}

// fuzzy_one_bad checks if one and two differ by exactly one character.
func fuzzy_one_bad(one, two string, l1, l2 int) bool {
	if l1 != l2 {
		return false
	}
	count := 0
	for i := 0; i < l2; i++ {
		if toLowerByte(one[i]) != toLowerByte(two[i]) {
			count++
			if count > 1 {
				return false
			}
		}
	}
	return true
}

// fuzzy_strcmp returns true if the strings match with minor typos.
// Checks for: transposed chars, one missing char, one extra char, one wrong char.
func fuzzy_strcmp(one, two string) bool {
	l1 := len(one)
	l2 := len(two)

	if l2 >= 4 && fuzzy_transpose(one, two, l1, l2) {
		return true
	}
	if l2 >= 5 && fuzzy_one_less(one, two, l1, l2) {
		return true
	}
	if l2 >= 5 && fuzzy_one_extra(one, two, l1, l2) {
		return true
	}
	if l2 >= 5 && fuzzy_one_bad(one, two, l1, l2) {
		return true
	}
	return false
}

// lower_array is a lookup table for fast lowercase conversion.
// Must be initialized by calling init_lower() before use.
var lower_array [256]byte

// init_lower initializes the lower_array lookup table for fast case conversion.
func init_lower() {
	for i := 0; i < 256; i++ {
		lower_array[i] = byte(i)
	}
	for i := byte('A'); i <= byte('Z'); i++ {
		lower_array[i] = i - 'A' + 'a'
	}
}

// toLowerByte converts a single byte to lowercase using the lookup table.
func toLowerByte(c byte) byte {
	return lower_array[c]
}

// lcase converts a string to lowercase.
func lcase(s string) string {
	return strings.ToLower(s)
}

// isalpha returns true if c is an alphabetic character (a-z or A-Z).
func isalpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// isdigit returns true if c is a digit (0-9).
func isdigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// toupper converts a single byte to uppercase.
func toupper(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - 'a' + 'A'
	}
	return c
}

// asfail panics with a formatted assertion failure message.
func asfail(file string, line int, cond string) {
	panic(fmt.Sprintf("assertion failure: %s (%d): %s", file, line, cond))
}

// iswhite returns true if c is a space or tab.
func iswhite(c byte) bool {
	return c == ' ' || c == '\t'
}

// int_comp compares two integers for sorting.
// Returns negative if a < b, 0 if equal, positive if a > b.
func int_comp(a, b int) int {
	return a - b
}

// Deprecated: readfile is not needed; game state is stored in SQLite.
// This function will panic if called.
func readfile(path string) int {
	panic("readfile: deprecated - use SQLite for game state persistence")
}

// Deprecated: readlin is not needed; game state is stored in SQLite.
// This function will panic if called.
func readlin() string {
	panic("readlin: deprecated - use SQLite for game state persistence")
}

// Deprecated: readlin_ew is not needed; game state is stored in SQLite.
// This function will panic if called.
func readlin_ew() string {
	panic("readlin_ew: deprecated - use SQLite for game state persistence")
}

// Deprecated: closefile is not needed; game state is stored in SQLite.
// This function will panic if called.
func closefile(path string) {
	panic("closefile: deprecated - use SQLite for game state persistence")
}

// Deprecated: getlin is not needed; game state is stored in SQLite.
// This function will panic if called.
func getlin(fp *os.File) string {
	panic("getlin: deprecated - use SQLite for game state persistence")
}

// Deprecated: getlin_ew is not needed; game state is stored in SQLite.
// This function will panic if called.
func getlin_ew(fp *os.File) string {
	panic("getlin_ew: deprecated - use SQLite for game state persistence")
}

// Deprecated: copy_fp is not needed; game state is stored in SQLite.
// This function will panic if called.
func copy_fp(a, b *os.File) {
	panic("copy_fp: deprecated - use SQLite for game state persistence")
}

// Deprecated: eat_leading_trailing_whitespace is not needed; game state is stored in SQLite.
// This function will panic if called.
func eat_leading_trailing_whitespace(s string) string {
	panic("eat_leading_trailing_whitespace: deprecated - use SQLite for game state persistence")
}
