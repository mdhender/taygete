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

// rnd returns a number in the range [low, high].
func rnd(low, high int) int {
	return teg.prng.IntN(high-low) + low
}

// load_seed restores our global prng state from the database.
func load_seed(path string) error {
	return teg.restorePrngState(path)
}

// save_seed writes our global prng state to the database.
func save_seed(path string) error {
	return teg.savePrngState(path)
}

// rnd returns a number in the range [low, high].
func (e *Engine) rnd(low, high int) int {
	return e.prng.IntN(high-low) + low
}
