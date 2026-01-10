// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// GateFile looks something like
// Gate*
type GateFile struct {
	Path  string `json:"path,omitempty"`
	Gates []Gate `json:"gates,omitempty"`
}

// Gate looks something like
// ID Kind Tag
// Sections*
// blankLine?
type Gate struct {
	ID       int       `json:"id,omitempty"`
	Kind     string    `json:"kind,omitempty"`
	Tag      int       `json:"tag,omitempty"`
	Sections []Section `json:"sections,omitempty"`
}

// Section looks something like
// ID
// <indent>Item
type Section struct {
	ID    string `json:"id,omitempty"`
	Items []Item `json:"items,omitempty"`
}

// Item looks something like
// <indent>ID Number
type Item struct {
	ID       string `json:"id,omitempty"`
	Quantity int    `json:"quantity,omitempty"`
}

func loadGateFile(input string) (GateFile, error) {
	gf := GateFile{Path: input}

	data, err := os.ReadFile(input)
	if err != nil {
		return gf, err
	}

	lines := strings.Split(string(data), "\n")
	i := 0

	for i < len(lines) {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			i++
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 || fields[1] != "gate" {
			return gf, fmt.Errorf("line %d: expected gate header, got %q", i+1, line)
		}

		id, err := strconv.Atoi(fields[0])
		if err != nil {
			return gf, fmt.Errorf("line %d: invalid gate id %q: %w", i+1, fields[0], err)
		}
		tag, err := strconv.Atoi(fields[2])
		if err != nil {
			return gf, fmt.Errorf("line %d: invalid gate tag %q: %w", i+1, fields[2], err)
		}

		gate := Gate{ID: id, Kind: fields[1], Tag: tag}
		i++

		for i < len(lines) {
			line = lines[i]
			if strings.TrimSpace(line) == "" {
				i++
				break
			}

			if !strings.HasPrefix(line, " ") {
				section := Section{ID: strings.TrimSpace(line)}
				i++

				for i < len(lines) {
					line = lines[i]
					if strings.TrimSpace(line) == "" || !strings.HasPrefix(line, " ") {
						break
					}
					itemFields := strings.Fields(line)
					if len(itemFields) >= 2 {
						qty, err := strconv.Atoi(itemFields[1])
						if err != nil {
							return gf, fmt.Errorf("line %d: invalid item quantity %q: %w", i+1, itemFields[1], err)
						}
						section.Items = append(section.Items, Item{ID: itemFields[0], Quantity: qty})
					}
					i++
				}

				gate.Sections = append(gate.Sections, section)
			} else {
				i++
			}
		}

		gf.Gates = append(gf.Gates, gate)
	}

	return gf, nil
}
