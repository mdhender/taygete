// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

// rnd returns a number in the range [low, high].
func rnd(low, high int) int {
	return teg.prng.IntN(high-low) + low
}

// load_seed restores our global prng state from the database.
func load_seed(path string) error {
	return teg.restorePrngState(path)
}

// save_seed writes our global prng state to the database.
func save_seed(path string) error {
	return teg.savePrngState(path)
}

// rnd returns a number in the range [low, high].
func (e *Engine) rnd(low, high int) int {
	return e.prng.IntN(high-low) + low
}
