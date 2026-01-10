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
