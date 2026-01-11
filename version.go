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

import "github.com/maloquacious/semver"

var (
	version = semver.Version{
		Major:      0,
		Minor:      23,
		Patch:      0,
		PreRelease: "alpha",
		Build:      semver.Commit(),
	}
)

func Version() semver.Version {
	return version
}

func (e *Engine) Version() semver.Version {
	return version
}
