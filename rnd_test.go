// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

import (
	"log/slog"
	"math/rand/v2"
	"testing"

	"github.com/mdhender/prng"
)

func TestSaveSeed(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	teg := Engine{
		db:     db,
		logger: slog.Default(),
		prng:   prng.New(rand.NewPCG(12345, 67890)),
	}

	if err := teg.savePrngState("test"); err != nil {
		t.Fatalf("savePrngState failed: %v", err)
	}

	var state []byte
	err = db.QueryRow("SELECT state FROM prng_state WHERE name = 'test'").Scan(&state)
	if err != nil {
		t.Fatalf("failed to read prng_state: %v", err)
	}

	if len(state) == 0 {
		t.Error("expected non-empty state blob")
	}
}

func TestLoadSeed(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	teg := Engine{
		db:     db,
		logger: slog.Default(),
		prng:   prng.New(rand.NewPCG(12_345, 67_890)),
	}

	// Save current state
	if err := teg.savePrngState("test"); err != nil {
		t.Fatalf("savePrngState failed: %v", err)
	}

	// Generate some values
	want1 := teg.prng.Uint64()
	want2 := teg.prng.Uint64()

	// Restore state
	if err := teg.restorePrngState("test"); err != nil {
		t.Fatalf("restorePrngState failed: %v", err)
	}

	// Should get same values
	got1 := teg.prng.Uint64()
	got2 := teg.prng.Uint64()

	if want1 != got1 {
		t.Errorf("expected value1=%d, got %d", want1, got1)
	}
	if want2 != got2 {
		t.Errorf("expected value2=%d, got %d", want2, got2)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	teg := Engine{
		db:     db,
		logger: slog.Default(),
		prng:   prng.New(rand.NewPCG(99_999, 88_888)),
	}

	// Generate some random numbers to advance state
	for i := 0; i < 100; i++ {
		teg.rnd(1, 100)
	}

	// Save state
	if err := teg.savePrngState("checkpoint"); err != nil {
		t.Fatalf("savePrngState failed: %v", err)
	}

	// Generate more numbers and record them
	expected := make([]int, 10)
	for i := range expected {
		expected[i] = teg.rnd(1, 1_000)
	}

	// Restore state
	if err := teg.restorePrngState("checkpoint"); err != nil {
		t.Fatalf("restorePrngState failed: %v", err)
	}

	// Generate numbers again - should match
	for i, want := range expected {
		got := teg.rnd(1, 1000)
		if got != want {
			t.Errorf("rnd[%d]: got %d, want %d", i, got, want)
		}
	}
}

func TestLoadSeedNotFound(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	teg := Engine{
		db:     db,
		logger: slog.Default(),
		prng:   prng.New(rand.NewPCG(99_999, 88_888)),
	}

	err = teg.restorePrngState("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent state")
	}
}

func TestRnd(t *testing.T) {
	teg := Engine{
		logger: slog.Default(),
		prng:   prng.New(rand.NewPCG(12_345, 67_890)),
	}

	for i := 0; i < 1000; i++ {
		got := teg.rnd(1, 10)
		if got < 1 || got > 10 {
			t.Errorf("rnd(1, 10) = %d, want value in [1, 10]", got)
		}
	}
}

func TestRndNegativeRange(t *testing.T) {
	teg := Engine{
		logger: slog.Default(),
		prng:   prng.New(rand.NewPCG(12_345, 67_890)),
	}

	for i := 0; i < 1000; i++ {
		got := teg.rnd(-10, -1)
		if got < -10 || got > -1 {
			t.Errorf("rnd(-10, -1) = %d, want value in [-10, -1]", got)
		}
	}
}

func TestRndDeterministic(t *testing.T) {
	teg := Engine{
		logger: slog.Default(),
		prng:   prng.New(rand.NewPCG(42, 42)),
	}

	first := make([]int, 10)
	for i := range first {
		first[i] = teg.rnd(1, 100)
	}

	teg = Engine{
		logger: slog.Default(),
		prng:   prng.New(rand.NewPCG(42, 42)),
	}

	for i, want := range first {
		got := teg.rnd(1, 100)
		if got != want {
			t.Errorf("rnd[%d]: got %d, want %d", i, got, want)
		}
	}
}
