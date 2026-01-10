// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

import (
	"encoding/binary"
	"encoding/json"
	"math/rand/v2"
	"os"
)

// rnd returns a number in the range [low, high].
func rnd(low, high int) int {
	return globalPRNGState.r.IntN(high-low) + low
}

var (
	globalPRNGState struct {
		Seed1 uint64 `json:"seed1,omitempty"`
		Seed2 uint64 `json:"seed2,omitempty"`
		pcg   *rand.PCG
		r     *rand.Rand
	}
)

// load_seed restores our global prng state from a file.
func load_seed(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		logger.Error("prng state file could not be opened", "path", path, "err", err)
		return err
	}
	if err = json.Unmarshal(data, &globalPRNGState); err != nil {
		logger.Error("prng state file is invalid", "path", path, "err", err)
		return err
	}
	globalPRNGState.pcg = rand.NewPCG(globalPRNGState.Seed1, globalPRNGState.Seed2)
	binData := make([]byte, 16)
	binary.LittleEndian.PutUint64(binData[0:8], globalPRNGState.Seed1)
	binary.LittleEndian.PutUint64(binData[8:16], globalPRNGState.Seed2)
	if err = globalPRNGState.pcg.UnmarshalBinary(binData); err != nil {
		logger.Error("prng state could not be restored", "path", path, "err", err)
		return err
	}
	globalPRNGState.r = rand.New(globalPRNGState.pcg)
	return nil
}

// save_seed writes our global prng state to a file.
func save_seed(path string) error {
	binData, err := globalPRNGState.pcg.MarshalBinary()
	if err != nil {
		logger.Error("prng state could not be marshaled", "path", path, "err", err)
		return err
	}
	globalPRNGState.Seed1 = binary.LittleEndian.Uint64(binData[0:8])
	globalPRNGState.Seed2 = binary.LittleEndian.Uint64(binData[8:16])

	data, err := json.MarshalIndent(globalPRNGState, "", "  ")
	if err != nil {
		logger.Error("prng state could not be marshaled", "path", path, "err", err)
		return err
	}
	err = os.WriteFile(path, data, 0o644)
	if err != nil {
		logger.Error("prng state save failed", "path", path, "err", err)
		return err
	}
	return nil
}
