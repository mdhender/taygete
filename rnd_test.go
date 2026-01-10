// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

import (
	"encoding/json"
	"math/rand/v2"
	"os"
	"path/filepath"
	"testing"

	"github.com/mdhender/prng"
)

func TestSaveSeed(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "seed.json")

	teg := Engine{
		prng: prng.New(rand.NewPCG(12345, 67890)),
	}

	if err := teg.savePrngState(path); err != nil {
		t.Fatalf("save_seed failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read seed file: %v", err)
	}

	var state struct {
		Seed1 uint64 `json:"seed1,omitempty"`
		Seed2 uint64 `json:"seed2,omitempty"`
	}
	if err := json.Unmarshal(data, &state); err != nil {
		t.Fatalf("seed file is not valid JSON: %v", err)
	}

	if state.Seed1 == 0 && state.Seed2 == 0 {
		t.Error("expected non-zero seeds in saved state")
	}
}

func TestLoadSeed(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "seed.json")

	seedJSON := `{"seed1": 111, "seed2": 222}`
	if err := os.WriteFile(path, []byte(seedJSON), 0o644); err != nil {
		t.Fatalf("failed to write seed file: %v", err)
	}

	teg := Engine{
		prng: prng.New(rand.NewPCG(12_345, 67_890)),
	}
	type state struct {
		seed1, seed2 uint64
	}

	want := state{
		seed1: teg.prng.Uint64(),
		seed2: teg.prng.Uint64(),
	}

	if err := teg.restorePrngState(path); err != nil {
		t.Fatalf("load_seed failed: %v", err)
	}

	if teg.prng == nil {
		t.Fatalf("expected prng to be initialized")
	}

	got := state{
		seed1: teg.prng.Uint64(),
		seed2: teg.prng.Uint64(),
	}

	if want.seed1 != got.seed1 {
		t.Errorf("expected seed1=%d, got %d", want.seed1, got.seed1)
	}
	if want.seed2 != got.seed2 {
		t.Errorf("expected seed2=%d, got %d", want.seed2, got.seed2)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "seed.json")

	teg := Engine{
		prng: prng.New(rand.NewPCG(99_999, 88_888)),
	}
	type state struct {
		seed1, seed2 uint64
	}

	// Generate some random numbers to advance state
	for i := 0; i < 100; i++ {
		rnd(1, 100)
	}

	// Save state
	if err := teg.savePrngState(path); err != nil {
		t.Fatalf("save_seed failed: %v", err)
	}

	// Generate more numbers and record them
	expected := make([]int, 10)
	for i := range expected {
		expected[i] = rnd(1, 1_000)
	}

	// Restore state
	if err := teg.restorePrngState(path); err != nil {
		t.Fatalf("load_seed failed: %v", err)
	}

	// Generate numbers again - should match
	for i, want := range expected {
		got := rnd(1, 1000)
		if got != want {
			t.Errorf("rnd[%d]: got %d, want %d", i, got, want)
		}
	}
}

func TestLoadSeedFileNotFound(t *testing.T) {
	teg := Engine{
		prng: prng.New(rand.NewPCG(99_999, 88_888)),
	}
	err := teg.restorePrngState("/nonexistent/path/seed.json")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestLoadSeedInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "seed.json")

	if err := os.WriteFile(path, []byte("not valid json"), 0o644); err != nil {
		t.Fatalf("failed to write seed file: %v", err)
	}

	teg := Engine{
		prng: prng.New(rand.NewPCG(99_999, 88_888)),
	}

	err := teg.restorePrngState(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestRnd(t *testing.T) {
	teg := Engine{
		prng: prng.New(rand.NewPCG(12_345, 67_890)),
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
		prng: prng.New(rand.NewPCG(12_345, 67_890)),
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
		prng: prng.New(rand.NewPCG(42, 42)),
	}

	first := make([]int, 10)
	for i := range first {
		first[i] = teg.rnd(1, 100)
	}

	teg = Engine{
		prng: prng.New(rand.NewPCG(42, 42)),
	}

	for i, want := range first {
		got := teg.rnd(1, 100)
		if got != want {
			t.Errorf("rnd[%d]: got %d, want %d", i, got, want)
		}
	}
}
