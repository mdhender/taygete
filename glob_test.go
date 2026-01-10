// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

import (
	"testing"
)

func TestGlobInit(t *testing.T) {
	e := &Engine{}

	// Set some non-zero values first
	e.globals.box_head[T_char] = 100
	e.globals.box_head[T_player] = 200
	e.globals.sub_head[sub_city] = 300

	// GlobInit should reset everything to 0
	e.GlobInit()

	for i := 0; i < T_MAX; i++ {
		if e.globals.box_head[i] != 0 {
			t.Errorf("box_head[%d] = %d, want 0", i, e.globals.box_head[i])
		}
	}

	for i := 0; i < SUB_MAX; i++ {
		if e.globals.sub_head[i] != 0 {
			t.Errorf("sub_head[%d] = %d, want 0", i, e.globals.sub_head[i])
		}
	}
}

func TestSysclock(t *testing.T) {
	e := &Engine{}

	// Set and get sysclock
	testTime := olytime{day: 15, turn: 42, days_since_epoch: 1260}
	e.SetSysclock(testTime)

	got := e.Sysclock()
	if got.day != testTime.day {
		t.Errorf("Sysclock().day = %d, want %d", got.day, testTime.day)
	}
	if got.turn != testTime.turn {
		t.Errorf("Sysclock().turn = %d, want %d", got.turn, testTime.turn)
	}
	if got.days_since_epoch != testTime.days_since_epoch {
		t.Errorf("Sysclock().days_since_epoch = %d, want %d", got.days_since_epoch, testTime.days_since_epoch)
	}
}

func TestKindFirstNext(t *testing.T) {
	e := &Engine{}
	e.GlobInit()

	// Create a chain of characters: 10 -> 20 -> 30
	e.globals.bx[10] = &box{kind: T_char, x_next_kind: 20}
	e.globals.bx[20] = &box{kind: T_char, x_next_kind: 30}
	e.globals.bx[30] = &box{kind: T_char, x_next_kind: 0}
	e.globals.box_head[T_char] = 10

	// Test KindFirst
	first := e.KindFirst(T_char)
	if first != 10 {
		t.Errorf("KindFirst(T_char) = %d, want 10", first)
	}

	// Test KindNext chain
	next := e.KindNext(10)
	if next != 20 {
		t.Errorf("KindNext(10) = %d, want 20", next)
	}

	next = e.KindNext(20)
	if next != 30 {
		t.Errorf("KindNext(20) = %d, want 30", next)
	}

	next = e.KindNext(30)
	if next != 0 {
		t.Errorf("KindNext(30) = %d, want 0", next)
	}

	// Test invalid kind
	first = e.KindFirst(-1)
	if first != 0 {
		t.Errorf("KindFirst(-1) = %d, want 0", first)
	}

	first = e.KindFirst(T_MAX)
	if first != 0 {
		t.Errorf("KindFirst(T_MAX) = %d, want 0", first)
	}
}

func TestSubFirstNext(t *testing.T) {
	e := &Engine{}
	e.GlobInit()

	// Create a chain of cities: 100 -> 200 -> 300
	e.globals.bx[100] = &box{kind: T_loc, skind: sub_city, x_next_sub: 200}
	e.globals.bx[200] = &box{kind: T_loc, skind: sub_city, x_next_sub: 300}
	e.globals.bx[300] = &box{kind: T_loc, skind: sub_city, x_next_sub: 0}
	e.globals.sub_head[sub_city] = 100

	// Test SubFirst
	first := e.SubFirst(sub_city)
	if first != 100 {
		t.Errorf("SubFirst(sub_city) = %d, want 100", first)
	}

	// Test SubNext chain
	next := e.SubNext(100)
	if next != 200 {
		t.Errorf("SubNext(100) = %d, want 200", next)
	}

	next = e.SubNext(200)
	if next != 300 {
		t.Errorf("SubNext(200) = %d, want 300", next)
	}

	next = e.SubNext(300)
	if next != 0 {
		t.Errorf("SubNext(300) = %d, want 0", next)
	}
}

func TestCharactersIterator(t *testing.T) {
	e := &Engine{}
	e.GlobInit()

	// Create a chain of characters: 10 -> 20 -> 30
	e.globals.bx[10] = &box{kind: T_char, x_next_kind: 20}
	e.globals.bx[20] = &box{kind: T_char, x_next_kind: 30}
	e.globals.bx[30] = &box{kind: T_char, x_next_kind: 0}
	e.globals.box_head[T_char] = 10

	chars := e.Characters()

	if len(chars) != 3 {
		t.Errorf("len(Characters()) = %d, want 3", len(chars))
	}

	expected := []int{10, 20, 30}
	for i, id := range chars {
		if id != expected[i] {
			t.Errorf("Characters()[%d] = %d, want %d", i, id, expected[i])
		}
	}
}

func TestPlayersIterator(t *testing.T) {
	e := &Engine{}
	e.GlobInit()

	// Create players: 100 -> 200
	e.globals.bx[100] = &box{kind: T_player, x_next_kind: 200}
	e.globals.bx[200] = &box{kind: T_player, x_next_kind: 0}
	e.globals.box_head[T_player] = 100

	players := e.Players()

	if len(players) != 2 {
		t.Errorf("len(Players()) = %d, want 2", len(players))
	}

	expected := []int{100, 200}
	for i, id := range players {
		if id != expected[i] {
			t.Errorf("Players()[%d] = %d, want %d", i, id, expected[i])
		}
	}
}

func TestCitiesIterator(t *testing.T) {
	e := &Engine{}
	e.GlobInit()

	// Create cities via subkind chain: 1000 -> 2000
	e.globals.bx[1000] = &box{kind: T_loc, skind: sub_city, x_next_sub: 2000}
	e.globals.bx[2000] = &box{kind: T_loc, skind: sub_city, x_next_sub: 0}
	e.globals.sub_head[sub_city] = 1000

	cities := e.Cities()

	if len(cities) != 2 {
		t.Errorf("len(Cities()) = %d, want 2", len(cities))
	}

	expected := []int{1000, 2000}
	for i, id := range cities {
		if id != expected[i] {
			t.Errorf("Cities()[%d] = %d, want %d", i, id, expected[i])
		}
	}
}

func TestEmptyIterators(t *testing.T) {
	e := &Engine{}
	e.GlobInit()

	// All iterators should return empty slices when no entities exist
	if len(e.Characters()) != 0 {
		t.Error("Characters() should be empty")
	}
	if len(e.Players()) != 0 {
		t.Error("Players() should be empty")
	}
	if len(e.Locations()) != 0 {
		t.Error("Locations() should be empty")
	}
	if len(e.Cities()) != 0 {
		t.Error("Cities() should be empty")
	}
	if len(e.Ships()) != 0 {
		t.Error("Ships() should be empty")
	}
}

func TestStringTables(t *testing.T) {
	// Verify kind_s has correct entries
	if kind_s[T_deleted] != "deleted" {
		t.Errorf("kind_s[T_deleted] = %q, want %q", kind_s[T_deleted], "deleted")
	}
	if kind_s[T_char] != "char" {
		t.Errorf("kind_s[T_char] = %q, want %q", kind_s[T_char], "char")
	}
	if kind_s[T_loc] != "loc" {
		t.Errorf("kind_s[T_loc] = %q, want %q", kind_s[T_loc], "loc")
	}

	// Verify subkind_s
	if subkind_s[sub_city] != "city" {
		t.Errorf("subkind_s[sub_city] = %q, want %q", subkind_s[sub_city], "city")
	}
	if subkind_s[sub_castle] != "castle" {
		t.Errorf("subkind_s[sub_castle] = %q, want %q", subkind_s[sub_castle], "castle")
	}

	// Verify direction strings
	if full_dir_s[DIR_N] != "north" {
		t.Errorf("full_dir_s[DIR_N] = %q, want %q", full_dir_s[DIR_N], "north")
	}
	if short_dir_s[DIR_N] != "n" {
		t.Errorf("short_dir_s[DIR_N] = %q, want %q", short_dir_s[DIR_N], "n")
	}

	// Verify exit_opposite
	if exit_opposite[DIR_N] != DIR_S {
		t.Errorf("exit_opposite[DIR_N] = %d, want %d", exit_opposite[DIR_N], DIR_S)
	}
	if exit_opposite[DIR_E] != DIR_W {
		t.Errorf("exit_opposite[DIR_E] = %d, want %d", exit_opposite[DIR_E], DIR_W)
	}

	// Verify month names
	if len(month_names) != 8 {
		t.Errorf("len(month_names) = %d, want 8", len(month_names))
	}
	if month_names[0] != "Fierce winds" {
		t.Errorf("month_names[0] = %q, want %q", month_names[0], "Fierce winds")
	}
}
