// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

import (
	"encoding/json"
	"math/rand/v2"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveSeed(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "seed.json")

	globalPRNGState.pcg = rand.NewPCG(12345, 67890)
	globalPRNGState.r = rand.New(globalPRNGState.pcg)

	if err := save_seed(path); err != nil {
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

	if err := load_seed(path); err != nil {
		t.Fatalf("load_seed failed: %v", err)
	}

	if globalPRNGState.Seed1 != 111 {
		t.Errorf("expected Seed1=111, got %d", globalPRNGState.Seed1)
	}
	if globalPRNGState.Seed2 != 222 {
		t.Errorf("expected Seed2=222, got %d", globalPRNGState.Seed2)
	}
	if globalPRNGState.pcg == nil {
		t.Error("expected pcg to be initialized")
	}
	if globalPRNGState.r == nil {
		t.Error("expected r to be initialized")
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "seed.json")

	globalPRNGState.pcg = rand.NewPCG(99999, 88888)
	globalPRNGState.r = rand.New(globalPRNGState.pcg)

	// Generate some random numbers to advance state
	for i := 0; i < 100; i++ {
		rnd(1, 100)
	}

	// Save state
	if err := save_seed(path); err != nil {
		t.Fatalf("save_seed failed: %v", err)
	}

	// Generate more numbers and record them
	expected := make([]int, 10)
	for i := range expected {
		expected[i] = rnd(1, 1000)
	}

	// Restore state
	if err := load_seed(path); err != nil {
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
	err := load_seed("/nonexistent/path/seed.json")
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

	err := load_seed(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestRnd(t *testing.T) {
	globalPRNGState.pcg = rand.NewPCG(12345, 67890)
	globalPRNGState.r = rand.New(globalPRNGState.pcg)

	for i := 0; i < 1000; i++ {
		got := rnd(1, 10)
		if got < 1 || got > 10 {
			t.Errorf("rnd(1, 10) = %d, want value in [1, 10]", got)
		}
	}
}

func TestRndNegativeRange(t *testing.T) {
	globalPRNGState.pcg = rand.NewPCG(12345, 67890)
	globalPRNGState.r = rand.New(globalPRNGState.pcg)

	for i := 0; i < 1000; i++ {
		got := rnd(-10, -1)
		if got < -10 || got > -1 {
			t.Errorf("rnd(-10, -1) = %d, want value in [-10, -1]", got)
		}
	}
}

func TestRndDeterministic(t *testing.T) {
	globalPRNGState.pcg = rand.NewPCG(42, 42)
	globalPRNGState.r = rand.New(globalPRNGState.pcg)

	first := make([]int, 10)
	for i := range first {
		first[i] = rnd(1, 100)
	}

	globalPRNGState.pcg = rand.NewPCG(42, 42)
	globalPRNGState.r = rand.New(globalPRNGState.pcg)

	for i, want := range first {
		got := rnd(1, 100)
		if got != want {
			t.Errorf("rnd[%d]: got %d, want %d", i, got, want)
		}
	}
}
