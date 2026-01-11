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
	"cmp"
	"fmt"
	"slices"
)

// List is a generic list type that replaces the legacy C-style ilist.
// It wraps a slice and provides encapsulated operations.
type List[T comparable] struct{ values []T }

// IList is a backward-compatible alias for List[int].
// Use this for existing code that expects integer lists.
type IList = List[int]

// NewList creates a new List with the given initial values.
// The returned list does not share memory with the input.
func NewList[T comparable](vals ...T) List[T] {
	if len(vals) == 0 {
		return List[T]{}
	}
	return List[T]{values: append([]T(nil), vals...)}
}

// Values returns the underlying slice.
// The returned slice shares memory with the list.
func (l *List[T]) Values() []T {
	if l == nil {
		return nil
	}
	return l.values
}

// Append appends value to the end of the list.
func (l *List[T]) Append(value T) {
	l.values = append(l.values, value)
}

// Clear clears the list and retains the underlying array.
func (l *List[T]) Clear() {
	l.values = l.values[:0]
}

// Copy returns a copy of the list.
// The copy does NOT share memory with the source list.
func (l *List[T]) Copy() *List[T] {
	if l == nil {
		return nil
	}
	if l.values == nil {
		return &List[T]{}
	}
	cp := &List[T]{values: make([]T, len(l.values))}
	copy(cp.values, l.values)
	return cp
}

// Delete removes the element at index i.
// Panics if i is out of bounds.
func (l *List[T]) Delete(i int) {
	if i < 0 || i >= len(l.values) {
		panic(fmt.Sprintf("List.Delete: index %d out of bounds [0, %d)", i, len(l.values)))
	}
	l.values = append(l.values[:i], l.values[i+1:]...)
}

// IndexOf returns the index of the first occurrence of value in the list,
// or -1, false if not found.
func (l *List[T]) IndexOf(value T) (int, bool) {
	for i, v := range l.values {
		if v == value {
			return i, true
		}
	}
	return -1, false
}

// Insert inserts value at position pos in the list.
// Panics if pos is out of bounds (must be 0 <= pos <= Len()).
func (l *List[T]) Insert(pos int, value T) {
	if pos < 0 || pos > len(l.values) {
		panic(fmt.Sprintf("List.Insert: position %d out of bounds [0, %d]", pos, len(l.values)))
	}
	var zero T
	l.values = append(l.values, zero)
	copy(l.values[pos+1:], l.values[pos:])
	l.values[pos] = value
}

// Len returns the length of the list.
func (l *List[T]) Len() int {
	if l == nil {
		return 0
	}
	return len(l.values)
}

// Lookup returns the index of the first occurrence of value in the list,
// or -1 if not found.
func (l *List[T]) Lookup(value T) int {
	idx, _ := l.IndexOf(value)
	return idx
}

// Prepend prepends value to the beginning of the list.
func (l *List[T]) Prepend(value T) {
	l.values = append([]T{value}, l.values...)
}

// Reclaim releases the list memory.
func (l *List[T]) Reclaim() {
	l.values = nil
}

// RemValue removes all occurrences of value from the list.
func (l *List[T]) RemValue(value T) {
	j := 0
	for _, v := range l.values {
		if v != value {
			l.values[j] = v
			j++
		}
	}
	l.values = l.values[:j]
}

// RemValueUniq removes the last occurrence of value from the list.
// Searches from the back (largest index) to match original C behavior.
func (l *List[T]) RemValueUniq(value T) {
	for i := len(l.values) - 1; i >= 0; i-- {
		if l.values[i] == value {
			l.values = append(l.values[:i], l.values[i+1:]...)
			return
		}
	}
}

// Scramble performs a Fisher-Yates shuffle on the list.
func (l *List[T]) Scramble() {
	n := len(l.values) - 1
	for i := 0; i < n; i++ {
		r := rnd(i, n)
		if r != i {
			l.values[i], l.values[r] = l.values[r], l.values[i]
		}
	}
}

// SortList sorts a list of ordered values in ascending order.
// This is a standalone function because Sort requires the cmp.Ordered constraint,
// which is more restrictive than the comparable constraint used by List[T].
func SortList[T cmp.Ordered](l *List[T]) {
	if l == nil {
		return
	}
	slices.Sort(l.values)
}

// SortListFunc sorts a list using a custom comparison function.
// The cmp function should return a negative number when a < b,
// zero when a == b, and a positive number when a > b.
func SortListFunc[T comparable](l *List[T], cmpFunc func(a, b T) int) {
	if l == nil {
		return
	}
	slices.SortFunc(l.values, cmpFunc)
}
