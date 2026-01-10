// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

// read_pw loads the password file, searches for the key.
// returns the current password for the key, if found, or an error.
func read_pw(key string) (string, error) {
	return teg.readPassword(key)
}
