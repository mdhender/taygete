// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

import (
	"database/sql"
	"log/slog"

	"github.com/mdhender/prng"
)

var (
	teg *Engine
)

type Engine struct {
	db     *sql.DB
	logger *slog.Logger
	prng   *prng.Rand
}

func NewEngine(db *sql.DB) (*Engine, error) {
	e := &Engine{db: db}
	err := e.restorePrngState(".")
	if err != nil {
		e.logger.Error("new engine", "err", err)
		return nil, err
	}
	return e, nil
}
