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
