// Copyright (c) 2026 Michael D Henderson. All rights reserved.

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
		names        map[int]string // entity names by ID
		banners      map[int]string // display banners by ID
		pluralNames  map[int]string // item plural names by ID
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
