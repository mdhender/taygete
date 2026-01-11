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

// storm.go ports storm.c and cloud.c - movement timing functions
//
// The movement timing functions track when entities (ships or characters)
// started moving and how long they've been moving. This is used to determine
// if an entity is "gone" (unavailable for interaction because it's in transit).
//
// The C code uses a global `evening` variable to track if we're in the
// evening phase of turn processing. For now, we use a simplified version
// that just checks if movement has started.
//
// Note: ship_moving and char_moving are defined in accessor.go

// ship_gone returns how long the ship has been moving (days), or 0 if not moving.
// When a ship is "gone", it cannot be interacted with (boarded, attacked, etc.)
// Ports: #define ship_gone(n) (ship_moving(n) ? sysclock.days_since_epoch - ship_moving(n) + evening : 0)
func ship_gone(n int) int {
	moving := ship_moving(n)
	if moving == 0 {
		return 0
	}
	return teg.globals.sysclock.days_since_epoch - moving + boolToInt(teg.globals.evening)
}

// char_gone returns how long the character has been moving (days), or 0 if not moving.
// When a character is "gone", it cannot be targeted for certain actions.
//
// Note: The C code has two versions controlled by #if 0:
//   - Full version: sysclock.days_since_epoch - char_moving(n) + evening
//   - Simplified: just 1 if moving, 0 if not
//
// We implement the simplified version to match the active C code.
// Ports: #define char_gone(n) (char_moving(n) ? 1 : 0)
func char_gone(n int) int {
	if char_moving(n) != 0 {
		return 1
	}
	return 0
}

// char_gone_full returns the actual days elapsed since movement started.
// This is the "full" version that was #if 0'd out in the C code.
// Provided for completeness and potential future use.
func char_gone_full(n int) int {
	moving := char_moving(n)
	if moving == 0 {
		return 0
	}
	return teg.globals.sysclock.days_since_epoch - moving + boolToInt(teg.globals.evening)
}

// boolToInt converts a boolean to an integer (0 or 1).
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
