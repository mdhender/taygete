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
	"strings"
	"testing"
)

// TestDeprecatedCommandsPanic verifies that deprecated command stubs panic when called.
// Sprint 25.9: Unit tests verify stubs panic.
func TestDeprecatedCommandsPanic(t *testing.T) {
	e := &Engine{}
	c := &command{}

	tests := []struct {
		name     string
		fn       func()
		contains string
	}{
		{
			name:     "v_split",
			fn:       func() { e.v_split(c) },
			contains: "v_split",
		},
		{
			name:     "v_format",
			fn:       func() { e.v_format(c) },
			contains: "v_format",
		},
		{
			name:     "v_notab",
			fn:       func() { e.v_notab(c) },
			contains: "v_notab",
		},
		{
			name:     "v_times",
			fn:       func() { e.v_times(c) },
			contains: "v_times",
		},
		{
			name:     "open_times",
			fn:       func() { e.open_times() },
			contains: "open_times",
		},
		{
			name:     "times_masthead",
			fn:       func() { e.times_masthead() },
			contains: "times_masthead",
		},
		{
			name:     "close_times",
			fn:       func() { e.close_times() },
			contains: "close_times",
		},
		{
			name:     "v_rumor",
			fn:       func() { e.v_rumor(c) },
			contains: "v_rumor",
		},
		{
			name:     "v_press",
			fn:       func() { e.v_press(c) },
			contains: "v_press",
		},
		{
			name:     "text_list_free",
			fn:       func() { e.text_list_free(nil) },
			contains: "text_list_free",
		},
		{
			name:     "line_length_check",
			fn:       func() { e.line_length_check(nil) },
			contains: "line_length_check",
		},
		{
			name:     "parse_text_list",
			fn:       func() { e.parse_text_list(c) },
			contains: "parse_text_list",
		},
		{
			name:     "v_post",
			fn:       func() { e.v_post(c) },
			contains: "v_post",
		},
		{
			name:     "v_message",
			fn:       func() { e.v_message(c) },
			contains: "v_message",
		},
		{
			name:     "v_tell",
			fn:       func() { e.v_tell(c) },
			contains: "v_tell",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Errorf("%s did not panic", tt.name)
					return
				}
				msg, ok := r.(string)
				if !ok {
					t.Errorf("%s panicked with non-string: %v", tt.name, r)
					return
				}
				if !strings.Contains(msg, "Deprecated") {
					t.Errorf("%s panic message missing 'Deprecated': %s", tt.name, msg)
				}
				if !strings.Contains(msg, tt.contains) {
					t.Errorf("%s panic message missing '%s': %s", tt.name, tt.contains, msg)
				}
			}()
			tt.fn()
		})
	}
}
