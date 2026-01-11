// taygete - a game engine for a game.
// Copyright (c) 2026 Michael D Henderson.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package taygete

import (
	"testing"
)

// TestCommandStates tests the state constants.
func TestCommandStates(t *testing.T) {
	// Verify state constants match C definitions
	if STATE_DONE != 0 {
		t.Errorf("STATE_DONE = %d, want 0", STATE_DONE)
	}
	if STATE_LOAD != 1 {
		t.Errorf("STATE_LOAD = %d, want 1", STATE_LOAD)
	}
	if STATE_RUN != 2 {
		t.Errorf("STATE_RUN = %d, want 2", STATE_RUN)
	}
	if STATE_ERROR != 3 {
		t.Errorf("STATE_ERROR = %d, want 3", STATE_ERROR)
	}
}

// TestMaxPriority tests the priority constant.
func TestMaxPriority(t *testing.T) {
	// MAX_PRI should be 5 to match C
	if MAX_PRI != 5 {
		t.Errorf("MAX_PRI = %d, want 5", MAX_PRI)
	}
}

// TestInitCommandQueues tests queue initialization.
func TestInitCommandQueues(t *testing.T) {
	e := &Engine{}

	e.initCommandQueues()

	if e.globals.cmdQueues == nil {
		t.Fatal("cmdQueues is nil after init")
	}

	// All load queues should be empty
	for pri := 0; pri < MAX_PRI; pri++ {
		if len(e.globals.cmdQueues.loadQ[pri]) != 0 {
			t.Errorf("loadQ[%d] not empty after init", pri)
		}
	}

	// Run queue should be empty
	if len(e.globals.cmdQueues.runQ) != 0 {
		t.Error("runQ not empty after init")
	}
}

// TestSetState tests state transitions.
func TestSetState(t *testing.T) {
	e := &Engine{}
	e.initCommandQueues()

	// Create a test box and command
	who := 1001
	e.globals.bx[who] = &box{kind: T_char}
	c := e.p_command(who)
	c.who = who
	c.state = STATE_DONE

	// Transition to STATE_LOAD at priority 2
	e.set_state(c, STATE_LOAD, 2)

	if c.state != STATE_LOAD {
		t.Errorf("state = %d, want %d", c.state, STATE_LOAD)
	}
	if !contains(e.globals.cmdQueues.loadQ[2], who) {
		t.Error("unit not in loadQ[2]")
	}

	// Transition to STATE_RUN
	e.set_state(c, STATE_RUN, 0)

	if c.state != STATE_RUN {
		t.Errorf("state = %d, want %d", c.state, STATE_RUN)
	}
	if contains(e.globals.cmdQueues.loadQ[2], who) {
		t.Error("unit still in loadQ[2] after transition to RUN")
	}
	if !contains(e.globals.cmdQueues.runQ, who) {
		t.Error("unit not in runQ")
	}

	// Transition to STATE_DONE
	e.set_state(c, STATE_DONE, 0)

	if c.state != STATE_DONE {
		t.Errorf("state = %d, want %d", c.state, STATE_DONE)
	}
	if contains(e.globals.cmdQueues.runQ, who) {
		t.Error("unit still in runQ after transition to DONE")
	}
}

// TestMinPriReady tests priority scheduling.
func TestMinPriReady(t *testing.T) {
	e := &Engine{}
	e.initCommandQueues()

	// No commands = priority 99
	if pri := e.min_pri_ready(); pri != 99 {
		t.Errorf("min_pri_ready() = %d, want 99 (no commands)", pri)
	}

	// Create two characters with commands at different priorities
	who1, who2 := 1001, 1002
	e.globals.bx[who1] = &box{kind: T_char, x_char: &entity_char{}}
	e.globals.bx[who2] = &box{kind: T_char, x_char: &entity_char{}}

	c1 := e.p_command(who1)
	c1.who = who1
	c1.pri = 3
	e.set_state(c1, STATE_LOAD, 3)

	c2 := e.p_command(who2)
	c2.who = who2
	c2.pri = 1
	e.set_state(c2, STATE_LOAD, 1)

	// Should return the lower priority (higher precedence)
	if pri := e.min_pri_ready(); pri != 1 {
		t.Errorf("min_pri_ready() = %d, want 1", pri)
	}

	// Remove priority 1 command
	e.set_state(c2, STATE_DONE, 0)

	// Now should return priority 3
	if pri := e.min_pri_ready(); pri != 3 {
		t.Errorf("min_pri_ready() = %d, want 3", pri)
	}
}

// TestIlistRemValue tests the helper function.
func TestIlistRemValue(t *testing.T) {
	l := []int{1, 2, 3, 4, 5}

	// Remove middle value
	l = ilistRemValue(l, 3)
	if contains(l, 3) {
		t.Error("3 still in list after removal")
	}
	if len(l) != 4 {
		t.Errorf("len = %d, want 4", len(l))
	}

	// Remove first value
	l = ilistRemValue(l, 1)
	if contains(l, 1) {
		t.Error("1 still in list after removal")
	}

	// Remove last value
	l = ilistRemValue(l, 5)
	if contains(l, 5) {
		t.Error("5 still in list after removal")
	}

	// Remove non-existent value (should not panic)
	l = ilistRemValue(l, 99)
	if len(l) != 2 {
		t.Errorf("len = %d, want 2", len(l))
	}
}

// TestParseLine tests command line parsing.
func TestParseLine(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"move north", []string{"move", "north"}},
		{"attack 1234", []string{"attack", "1234"}},
		{`say "hello world"`, []string{"say", "hello world"}},
		{`name 'Test Unit'`, []string{"name", "Test Unit"}},
		{"  wait  time  7 ", []string{"wait", "time", "7"}},
		{"", nil},
		{"   ", nil},
		{"study 600 14", []string{"study", "600", "14"}},
	}

	for _, tc := range tests {
		got := parse_line(tc.input)
		if len(got) != len(tc.want) {
			t.Errorf("parse_line(%q) = %v, want %v", tc.input, got, tc.want)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("parse_line(%q)[%d] = %q, want %q", tc.input, i, got[i], tc.want[i])
			}
		}
	}
}

// TestLoadCommand tests command loading from order queue.
func TestLoadCommand(t *testing.T) {
	e := &Engine{}
	e.initCommandQueues()

	// Create a player and character
	playerID := 100
	who := 1001
	e.globals.bx[playerID] = &box{kind: T_player, x_player: &entity_player{}}
	e.globals.bx[who] = &box{kind: T_char, x_char: &entity_char{unit_lord: playerID}}

	c := e.p_command(who)
	c.who = who

	// No orders = should return false and set state to DONE
	if e.load_command(c) {
		t.Error("load_command returned true with empty queue")
	}
	if c.state != STATE_DONE {
		t.Errorf("state = %d, want %d", c.state, STATE_DONE)
	}

	// Add an order
	e.queue_order(playerID, who, "wait time 7")

	// Now should load successfully
	if !e.load_command(c) {
		t.Error("load_command returned false with order in queue")
	}
	if c.state != STATE_LOAD {
		t.Errorf("state = %d, want %d", c.state, STATE_LOAD)
	}
}

// TestInitLoadSup tests initialization of command loading for a unit.
func TestInitLoadSup(t *testing.T) {
	e := &Engine{}
	e.initCommandQueues()

	playerID := 100
	who := 1001
	e.globals.bx[playerID] = &box{kind: T_player, x_player: &entity_player{}}
	e.globals.bx[who] = &box{kind: T_char, x_char: &entity_char{unit_lord: playerID}}

	// Set up player in kind list for Players() to work
	e.globals.box_head[T_player] = playerID

	// Add an order
	e.queue_order(playerID, who, "move north")

	// Initialize command loading
	e.init_load_sup(who)

	// Command should be loaded
	c := e.rp_command(who)
	if c == nil {
		t.Fatal("rp_command returned nil")
	}
	if c.state != STATE_LOAD {
		t.Errorf("state = %d, want %d", c.state, STATE_LOAD)
	}
}

// TestCommandDone tests command completion.
func TestCommandDone(t *testing.T) {
	e := &Engine{}
	e.initCommandQueues()

	playerID := 100
	who := 1001
	e.globals.bx[playerID] = &box{kind: T_player, x_player: &entity_player{}}
	e.globals.bx[who] = &box{kind: T_char, x_char: &entity_char{unit_lord: playerID}}

	c := e.p_command(who)
	c.who = who
	e.set_state(c, STATE_RUN, 0)

	// No more orders = should transition to DONE
	e.commandDone(c)

	if c.state != STATE_DONE {
		t.Errorf("state = %d, want %d", c.state, STATE_DONE)
	}

	// Add orders and reset to RUN
	e.queue_order(playerID, who, "wait time 1")
	e.queue_order(playerID, who, "move north")
	e.set_state(c, STATE_RUN, 0)

	// Should load next command
	e.commandDone(c)

	if c.state != STATE_LOAD {
		t.Errorf("state = %d, want %d after loading next", c.state, STATE_LOAD)
	}
}

// TestEveningPhase tests the evening phase processing.
func TestEveningPhase(t *testing.T) {
	e := &Engine{}
	e.initCommandQueues()

	who := 1001
	e.globals.bx[who] = &box{kind: T_char, x_char: &entity_char{}}

	c := e.p_command(who)
	c.who = who
	c.wait = 3      // 3 days to complete
	c.status = TRUE // command is in progress (success so far)
	c.days_executing = 0
	e.set_state(c, STATE_RUN, 0)

	// Run evening phase
	e.evening_phase()

	// Days executing should increment
	if c.days_executing != 1 {
		t.Errorf("days_executing = %d, want 1", c.days_executing)
	}

	// Wait should decrement
	if c.wait != 2 {
		t.Errorf("wait = %d, want 2", c.wait)
	}

	// Should still be running
	if c.state != STATE_RUN {
		t.Errorf("state = %d, want %d", c.state, STATE_RUN)
	}
}

// TestPrisonerCannotExecute tests that prisoners cannot execute commands.
func TestPrisonerCannotExecute(t *testing.T) {
	e := &Engine{}
	e.initCommandQueues()

	who := 1001
	e.globals.bx[who] = &box{kind: T_char, x_char: &entity_char{prisoner: 1}}

	c := e.p_command(who)
	c.who = who
	c.pri = 2
	e.set_state(c, STATE_LOAD, 2)

	// min_pri_ready should skip prisoners
	if pri := e.min_pri_ready(); pri != 99 {
		t.Errorf("min_pri_ready() = %d, want 99 (prisoner should be skipped)", pri)
	}
}

// contains checks if a slice contains a value.
func contains(l []int, val int) bool {
	for _, v := range l {
		if v == val {
			return true
		}
	}
	return false
}
