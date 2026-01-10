// Copyright (c) 2026 Michael D Henderson. All rights reserved.

// Package taygete implements Olympia in Go.
package taygete

import (
	"github.com/maloquacious/semver"
	"log/slog"
	"os"
)

var (
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	version = semver.Version{
		Major:      0,
		Minor:      1,
		Patch:      0,
		PreRelease: "alpha",
		Build:      semver.Commit(),
	}
)

func Version() semver.Version {
	return version
}
