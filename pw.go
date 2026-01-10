// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

import (
	"encoding/json"
	"fmt"
	"os"
)

const plaintextPasswordFile = "PWS"

type plaintextPasswordItem struct {
	Key      string `json:"key,omitempty"`
	Password string `json:"password,omitempty"`
}

// read_pw loads the password file, searches for the key.
// returns the current password for the key, if found, or an error.
func read_pw(key string) (string, error) {
	data, err := os.ReadFile(plaintextPasswordFile)
	if err != nil {
		return "", fmt.Errorf("can't read password file: %w", err)
	}
	var items []plaintextPasswordItem
	if err := json.Unmarshal(data, &items); err != nil {
		return "", fmt.Errorf("%s: unmarshal: %w", plaintextPasswordFile, err)
	}
	for _, item := range items {
		if item.Key == key {
			return item.Password, nil
		}
	}
	return "", fmt.Errorf("%s: not found", key)
}
