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

// cmd_deprecated.go - Deprecated text-report command stubs from src/c1.c and src/c2.c
// Sprint 25.9: Stubs that panic if called, per TODO.md guidance.
//
// These commands were part of the legacy email/text-based report system
// which is replaced by the web frontend in the Go/DB version.

package taygete

// v_split is deprecated: legacy report splitting. Not used in Go/DB version.
// The C version already marked this as "no longer supported".
// Deprecated: use web frontend for report viewing.
func (e *Engine) v_split(c *command) int {
	panic("Deprecated: v_split (report splitting) is not supported in the Go/DB version")
}

// v_format is deprecated: legacy report format setting.
// Set player's report format preference for email reports.
// Deprecated: use web frontend for report viewing.
func (e *Engine) v_format(c *command) int {
	panic("Deprecated: v_format (report format) is not supported in the Go/DB version")
}

// v_notab is deprecated: legacy TAB character preference for email reports.
// Set whether TAB characters appear in turn reports.
// Deprecated: use web frontend for report viewing.
func (e *Engine) v_notab(c *command) int {
	panic("Deprecated: v_notab (no-TABs option) is not supported in the Go/DB version")
}

// v_times is deprecated: legacy Times newspaper subscription setting.
// Set whether player receives the Olympia Times via email.
// Deprecated: use web frontend for Times viewing.
func (e *Engine) v_times(c *command) int {
	panic("Deprecated: v_times (Times subscription) is not supported in the Go/DB version")
}

// open_times is deprecated: legacy Times file I/O.
// Opens rumor_fp and press_fp file handles for writing the Times.
// Deprecated: use web frontend for Times content.
func (e *Engine) open_times() {
	panic("Deprecated: open_times (Times file I/O) is not supported in the Go/DB version")
}

// times_masthead is deprecated: legacy Times masthead generation.
// Writes the Times header/masthead to file.
// Deprecated: use web frontend for Times content.
func (e *Engine) times_masthead() {
	panic("Deprecated: times_masthead (Times file I/O) is not supported in the Go/DB version")
}

// close_times is deprecated: legacy Times file I/O.
// Closes rumor_fp and press_fp file handles.
// Deprecated: use web frontend for Times content.
func (e *Engine) close_times() {
	panic("Deprecated: close_times (Times file I/O) is not supported in the Go/DB version")
}

// v_rumor is deprecated: legacy Times rumor submission.
// Submit anonymous rumors to the Olympia Times.
// Deprecated: use web frontend for Times submissions.
func (e *Engine) v_rumor(c *command) int {
	panic("Deprecated: v_rumor (Times submissions) is not supported in the Go/DB version")
}

// v_press is deprecated: legacy Times press submission.
// Submit attributed press releases to the Olympia Times.
// Deprecated: use web frontend for Times submissions.
func (e *Engine) v_press(c *command) int {
	panic("Deprecated: v_press (Times submissions) is not supported in the Go/DB version")
}

// text_list_free is deprecated: legacy text list memory management.
// Frees a plist of strings allocated during text parsing.
// Deprecated: Go uses garbage collection.
func (e *Engine) text_list_free(l []string) {
	panic("Deprecated: text_list_free (text list parsing) is not supported in the Go/DB version")
}

// line_length_check is deprecated: legacy text validation.
// Returns the maximum line length in a text list.
// Deprecated: use web frontend for text input validation.
func (e *Engine) line_length_check(l []string) int {
	panic("Deprecated: line_length_check (text list parsing) is not supported in the Go/DB version")
}

// parse_text_list is deprecated: legacy multi-line text parsing.
// Parses multi-line text from order stream (e.g., for POST, MESSAGE commands).
// Deprecated: use web frontend for text input.
func (e *Engine) parse_text_list(c *command) []string {
	panic("Deprecated: parse_text_list (text list parsing) is not supported in the Go/DB version")
}

// v_post is deprecated: legacy in-game posting.
// Creates a post entity at the current location.
// Deprecated: use web frontend for in-game messaging.
func (e *Engine) v_post(c *command) int {
	panic("Deprecated: v_post (in-game posting) is not supported in the Go/DB version")
}

// v_message is deprecated: legacy in-game messaging.
// Sends a multi-line message to a target entity.
// Deprecated: use web frontend for in-game messaging.
func (e *Engine) v_message(c *command) int {
	panic("Deprecated: v_message (in-game posting) is not supported in the Go/DB version")
}

// v_tell is deprecated: legacy knowledge sharing.
// The C version was already disabled with #if 0.
// Deprecated: TELL order was removed as of turn 50 in original game.
func (e *Engine) v_tell(c *command) int {
	panic("Deprecated: v_tell (TELL order) was removed as of turn 50")
}
