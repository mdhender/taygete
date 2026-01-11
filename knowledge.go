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

// Bitset & knowledge helpers ported from src/u.c (Sprint 25.2).
//
// In C, sparse is typedef'd as int* and used as an ilist (dynamic array).
// In Go, we use map[int]bool for efficient set operations.

// test_bit checks if a value exists in a sparse set.
// Ported from src/u.c test_bit().
func test_bit(kr map[int]bool, i int) bool {
	if kr == nil {
		return false
	}
	return kr[i]
}

// set_bit adds a value to a sparse set (idempotent).
// Ported from src/u.c set_bit().
func set_bit(kr map[int]bool, i int) map[int]bool {
	if kr == nil {
		kr = make(map[int]bool)
	}
	kr[i] = true
	return kr
}

// clear_know_rec clears all entries from a sparse set.
// Ported from src/u.c clear_know_rec().
func clear_know_rec(kr map[int]bool) {
	for k := range kr {
		delete(kr, k)
	}
}

// test_known checks if entity i is known to the player owning who.
// Ported from src/u.c test_known().
func test_known(who, i int) bool {
	if who == 0 {
		return false
	}
	if !valid_box(who) || !valid_box(i) {
		return false
	}

	pl := player(who)
	if pl == 0 {
		return false
	}

	known := teg.getPlayerKnowledge(pl)
	return test_bit(known, i)
}

// set_known marks entity i as known to the player owning who.
// Ported from src/u.c set_known().
func set_known(who, i int) {
	if !valid_box(who) || !valid_box(i) {
		return
	}

	pl := player(who)
	if !valid_box(pl) {
		return
	}

	teg.setPlayerKnowledge(pl, i)
}
