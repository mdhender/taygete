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
	"database/sql"
	"log/slog"
	"math/rand/v2"

	"github.com/mdhender/prng"
)

var (
	// use a global Engine while porting
	teg *Engine
)

type Engine struct {
	db     *sql.DB
	logger *slog.Logger
	prng   *prng.Rand
	// use this globals struct for C globals while porting.
	// as we refactor, these will become state in Engine.
	globals struct {
		bx       [MAX_BOXES]*box
		box_head [T_MAX]int
		sub_head [SUB_MAX]int

		// game state flags (from glob.c)
		show_day          bool
		post_has_been_run bool
		garrison_magic    int
		mount_olympus     int
		combat_pl         int // Combat log player
		sysclock          olytime
		evening           bool // are we in the evening phase?

		// String storage - Go uses strings, not char pointers
		names       map[int]string // entity names by ID
		banners     map[int]string // display banners by ID
		pluralNames map[int]string // item plural names by ID

		// Skill storage - workaround for C-style **skill_ent plist
		charSkills map[int][]*skill_ent // character skills by char ID

		// Player units - workaround for C-style ilist in entity_player
		playerUnits map[int][]int

		// Inventory storage - workaround for C-style **item_ent in box
		inventories map[int][]item_ent

		// Order queues - maps player ID -> (unit ID -> OrderQueue)
		// Replaces C entity_player.orders plist
		orderQueues map[int]map[int]*OrderQueue

		// Turn processing state (from day.c / input.c)
		monthDone     bool // set when monthly processing completes
		autoQuitTurns int  // turns without orders before auto-quit (0 = disabled)

		// Command scheduling (from input.c - Sprint 22)
		cmdQueues      *commandQueues // scheduling queues for command execution
		immediate      bool           // true during immediate command execution
		autoAttackFlag bool           // check for auto-attacks once per day

		// Immediate mode state (from immed.c - Sprint 23)
		immedSeeAll bool // reveal all hidden features in immediate mode

		// Player knowledge sets (Sprint 25.2)
		// Maps player ID -> set of known entity IDs
		playerKnowledge map[int]map[int]bool

		// Region helpers (Sprint 25.6)
		nprov          int // cached province count (0 = not computed)
		faeryRegion    int // Faery realm region ID
		hadesRegion    int // Hades realm region ID
		nowhereRegion  int // Nowhere region ID
		cloudRegion    int // Cloud realm region ID
		tunnelRegion   int // Tunnel realm region ID
		underRegion    int // Underground realm region ID
	}
}

// getName returns the name for entity n.
func (e *Engine) getName(n int) string {
	if e.globals.names == nil {
		return ""
	}
	return e.globals.names[n]
}

// setName sets the name for entity n.
func (e *Engine) setName(n int, s string) {
	if e.globals.names == nil {
		e.globals.names = make(map[int]string)
	}
	if s == "" {
		delete(e.globals.names, n)
	} else {
		e.globals.names[n] = s
	}
}

// getBanner returns the display banner for entity n.
func (e *Engine) getBanner(n int) string {
	if e.globals.banners == nil {
		return ""
	}
	return e.globals.banners[n]
}

// setBanner sets the display banner for entity n.
func (e *Engine) setBanner(n int, s string) {
	if e.globals.banners == nil {
		e.globals.banners = make(map[int]string)
	}
	if s == "" {
		delete(e.globals.banners, n)
	} else {
		e.globals.banners[n] = s
	}
}

// getPluralName returns the plural name for item n.
func (e *Engine) getPluralName(n int) string {
	if e.globals.pluralNames == nil {
		return ""
	}
	return e.globals.pluralNames[n]
}

// setPluralName sets the plural name for item n.
func (e *Engine) setPluralName(n int, s string) {
	if e.globals.pluralNames == nil {
		e.globals.pluralNames = make(map[int]string)
	}
	if s == "" {
		delete(e.globals.pluralNames, n)
	} else {
		e.globals.pluralNames[n] = s
	}
}

// getPlayerKnowledge returns the knowledge set for player pl.
func (e *Engine) getPlayerKnowledge(pl int) map[int]bool {
	if e.globals.playerKnowledge == nil {
		return nil
	}
	return e.globals.playerKnowledge[pl]
}

// setPlayerKnowledge marks entity i as known to player pl.
func (e *Engine) setPlayerKnowledge(pl, i int) {
	if e.globals.playerKnowledge == nil {
		e.globals.playerKnowledge = make(map[int]map[int]bool)
	}
	if e.globals.playerKnowledge[pl] == nil {
		e.globals.playerKnowledge[pl] = make(map[int]bool)
	}
	e.globals.playerKnowledge[pl][i] = true
}

// clearPlayerKnowledge clears all knowledge for player pl.
func (e *Engine) clearPlayerKnowledge(pl int) {
	if e.globals.playerKnowledge == nil {
		return
	}
	if known := e.globals.playerKnowledge[pl]; known != nil {
		clear_know_rec(known)
	}
}

func NewEngine(db *sql.DB, p *prng.Rand) (*Engine, error) {
	if p == nil {
		p = prng.New(rand.NewPCG(0xC0FFEECAFE, 0xBEEFF00D))
	}
	e := &Engine{db: db, prng: p}
	err := e.restorePrngState(".")
	if err != nil {
		e.logger.Error("new engine", "err", err)
		return nil, err
	}
	return e, nil
}

func init() {
	db, err := OpenTestDB()
	if err != nil {
		panic(err)
	}
	teg = &Engine{
		db:   db,
		prng: prng.New(rand.NewPCG(0xC0FFEECAFE, 0xBEEFF00D)),
	}
	// initialize game state flags
	teg.globals.garrison_magic = 999
}
