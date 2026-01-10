// Copyright (c) 2026 Michael D Henderson. All rights reserved.

// Package taygete implements Olympia in Go.
package taygete

import (
	"log/slog"
	"os"
)

var (
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
)
