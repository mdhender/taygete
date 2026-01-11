// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

import (
	"testing"
)

func TestProcessOrdersNoOp(t *testing.T) {
	e := teg

	// Save initial state
	initialTurn := e.globals.sysclock.turn
	initialDay := e.globals.sysclock.day

	// Reset state for clean test
	e.globals.sysclock.turn = 0
	e.globals.sysclock.day = 0
	e.globals.monthDone = false

	// Run process_orders with no entities
	err := e.ProcessOrders()
	if err != nil {
		t.Fatalf("ProcessOrders returned error: %v", err)
	}

	// Verify turn was incremented
	if e.globals.sysclock.turn != 1 {
		t.Errorf("sysclock.turn: got %d, want 1", e.globals.sysclock.turn)
	}

	// Verify day advanced to MONTH_DAYS
	if e.globals.sysclock.day != MONTH_DAYS {
		t.Errorf("sysclock.day: got %d, want %d", e.globals.sysclock.day, MONTH_DAYS)
	}

	// Verify monthDone flag is set
	if !e.globals.monthDone {
		t.Error("monthDone should be true after ProcessOrders")
	}

	// Restore initial state
	e.globals.sysclock.turn = initialTurn
	e.globals.sysclock.day = initialDay
	e.globals.monthDone = false
}

func TestPostMonthNoOp(t *testing.T) {
	e := teg

	// Save initial state
	initialPostRun := e.globals.post_has_been_run

	// Reset state
	e.globals.post_has_been_run = false
	e.globals.sysclock.turn = 1 // month 1

	// Run post_month with no entities
	err := e.PostMonth()
	if err != nil {
		t.Fatalf("PostMonth returned error: %v", err)
	}

	// Verify post_has_been_run is set
	if !e.globals.post_has_been_run {
		t.Error("post_has_been_run should be true after PostMonth")
	}

	// Restore initial state
	e.globals.post_has_been_run = initialPostRun
}

func TestRunTurnNoOp(t *testing.T) {
	e := teg

	// Save initial state
	initialTurn := e.globals.sysclock.turn
	initialDay := e.globals.sysclock.day
	initialPostRun := e.globals.post_has_been_run
	initialMonthDone := e.globals.monthDone

	// Reset state
	e.globals.sysclock.turn = 0
	e.globals.sysclock.day = 0
	e.globals.post_has_been_run = false
	e.globals.monthDone = false

	// Run complete turn
	err := e.RunTurn()
	if err != nil {
		t.Fatalf("RunTurn returned error: %v", err)
	}

	// Verify both phases completed
	if !e.globals.monthDone {
		t.Error("monthDone should be true after RunTurn")
	}
	if !e.globals.post_has_been_run {
		t.Error("post_has_been_run should be true after RunTurn")
	}

	// Restore initial state
	e.globals.sysclock.turn = initialTurn
	e.globals.sysclock.day = initialDay
	e.globals.post_has_been_run = initialPostRun
	e.globals.monthDone = initialMonthDone
}

func TestOlyMonth(t *testing.T) {
	e := teg

	tests := []struct {
		turn int
		want int
	}{
		{1, 1},  // turn 1 = month 1
		{2, 2},  // turn 2 = month 2
		{8, 8},  // turn 8 = month 8
		{9, 1},  // turn 9 = month 1 (wraps)
		{10, 2}, // turn 10 = month 2
		{16, 8}, // turn 16 = month 8
		{17, 1}, // turn 17 = month 1
	}

	for _, tc := range tests {
		e.globals.sysclock.turn = short(tc.turn)
		got := e.olyMonth()
		if got != tc.want {
			t.Errorf("olyMonth() with turn=%d: got %d, want %d", tc.turn, got, tc.want)
		}
	}
}

func TestOlytimeIncrement(t *testing.T) {
	e := teg

	// Save initial state
	initialDay := e.globals.sysclock.day
	initialEpoch := e.globals.sysclock.days_since_epoch

	// Reset
	e.globals.sysclock.day = 5
	e.globals.sysclock.days_since_epoch = 100

	e.olytimeIncrement()

	if e.globals.sysclock.day != 6 {
		t.Errorf("day after increment: got %d, want 6", e.globals.sysclock.day)
	}
	if e.globals.sysclock.days_since_epoch != 101 {
		t.Errorf("days_since_epoch after increment: got %d, want 101", e.globals.sysclock.days_since_epoch)
	}

	// Restore
	e.globals.sysclock.day = initialDay
	e.globals.sysclock.days_since_epoch = initialEpoch
}

func TestOlytimeTurnChange(t *testing.T) {
	e := teg

	// Save initial state
	initialTurn := e.globals.sysclock.turn
	initialDay := e.globals.sysclock.day

	// Reset
	e.globals.sysclock.turn = 5
	e.globals.sysclock.day = 15

	e.olytimeTurnChange()

	if e.globals.sysclock.turn != 6 {
		t.Errorf("turn after change: got %d, want 6", e.globals.sysclock.turn)
	}
	if e.globals.sysclock.day != 0 {
		t.Errorf("day after turn change: got %d, want 0", e.globals.sysclock.day)
	}

	// Restore
	e.globals.sysclock.turn = initialTurn
	e.globals.sysclock.day = initialDay
}

func TestMonthDays(t *testing.T) {
	if MONTH_DAYS != 30 {
		t.Errorf("MONTH_DAYS: got %d, want 30", MONTH_DAYS)
	}
}
