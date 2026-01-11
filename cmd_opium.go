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

// cmd_opium.go - Opium & Misc Commands ported from src/c2.c
// Sprint 26.10: Opium & Misc Commands

package taygete

// v_improve_opium starts opium improvement.
// Must be in a poppy field. Training takes 7 days (see use.c skill table).
// Ported from src/c2.c lines 565-576.
func v_improve_opium(c *command) int {
	where := subloc(c.who)

	if subkind(where) != sub_poppy_field {
		wout(c.who, "Opium is produced only in poppy fields.")
		return FALSE
	}

	return TRUE
}

// d_improve_opium executes opium improvement, doubling production for one turn.
// Sets the opium_double flag on the poppy field.
// Ported from src/c2.c lines 580-593.
func d_improve_opium(c *command) int {
	where := subloc(c.who)

	if subkind(where) != sub_poppy_field {
		wout(c.who, "Not in a poppy field anymore.")
		return FALSE
	}

	p_misc(where).opium_double = TRUE

	return TRUE
}

// v_die is the suicide command. Calls kill_char to handle death.
// Ported from src/c2.c lines 597-602.
func v_die(c *command) int {
	kill_char(c.who, MATES)
	return TRUE
}
