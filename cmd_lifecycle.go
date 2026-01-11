// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

// cmd_lifecycle.go -- Command lifecycle and scheduling (Sprint 22)
//
// This file ports the command scheduling logic from input.c.
// Commands have four states:
//   STATE_DONE  - no command active, ready to load next
//   STATE_LOAD  - command parsed and ready to execute
//   STATE_RUN   - command is currently executing
//   STATE_ERROR - command has an error (will be reported to player)
//
// Commands have priorities (0 = highest, MAX_PRI-1 = lowest).
// Each day, the start_phase executes LOAD commands in priority order,
// and the evening_phase advances RUN commands.

const MAX_PRI = 5

// commandQueues holds the scheduling queues for command execution.
// These are runtime-only (not saved to DB).
type commandQueues struct {
	loadQ [MAX_PRI][]int // units with commands in STATE_LOAD, by priority
	runQ  []int          // units with commands in STATE_RUN
	curPri int           // current priority being processed
}

// initCommandQueues initializes the command scheduling queues.
func (e *Engine) initCommandQueues() {
	e.globals.cmdQueues = &commandQueues{}
	for i := 0; i < MAX_PRI; i++ {
		e.globals.cmdQueues.loadQ[i] = nil
	}
	e.globals.cmdQueues.runQ = nil
	e.globals.cmdQueues.curPri = 0
}

// set_state transitions a command to a new state, updating the scheduling queues.
// Port of C set_state() from input.c.
func (e *Engine) set_state(c *command, state int, newPri int) {
	if c == nil {
		return
	}

	queues := e.globals.cmdQueues
	if queues == nil {
		e.initCommandQueues()
		queues = e.globals.cmdQueues
	}

	// Remove from current queue based on old state
	switch c.state {
	case STATE_RUN:
		queues.runQ = ilistRemValue(queues.runQ, c.who)
	case STATE_LOAD:
		if c.pri >= 0 && int(c.pri) < MAX_PRI {
			queues.loadQ[c.pri] = ilistRemValue(queues.loadQ[c.pri], c.who)
		}
	}

	// Update state
	c.state = schar(state)

	// Add to new queue based on new state
	switch c.state {
	case STATE_RUN:
		queues.runQ = append(queues.runQ, c.who)
	case STATE_LOAD:
		if newPri >= 0 && newPri < MAX_PRI {
			queues.loadQ[newPri] = append(queues.loadQ[newPri], c.who)
			c.pri = schar(newPri) // update command priority to match queue
		}
	}
}

// ilistRemValue removes a value from an int slice (helper function).
func ilistRemValue(l []int, val int) []int {
	for i, v := range l {
		if v == val {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

// p_command returns the command structure for an entity, creating one if needed.
// Port of C p_command().
func (e *Engine) p_command(who int) *command {
	if who <= 0 || who >= MAX_BOXES {
		return nil
	}
	if e.globals.bx[who] == nil {
		return nil
	}
	if e.globals.bx[who].cmd == nil {
		e.globals.bx[who].cmd = &command{who: who, state: STATE_DONE}
	}
	return e.globals.bx[who].cmd
}

// get_command retrieves the next order from the unit's queue.
// Returns the order text and true if an order exists, or empty string and false.
// Port of C get_command().
func (e *Engine) get_command(who int) (string, bool) {
	pl := e.player(who)
	if pl == 0 {
		return "", false
	}

	s := e.top_order(pl, who)
	if s == "" {
		return "", false
	}

	e.pop_order(pl, who)
	return s, true
}

// load_command loads the next command from the order queue into the command structure.
// Returns true if a command was loaded, false if no commands remain.
// Port of C load_command() from input.c.
//
// Sets c.state to:
//   STATE_DONE  - no more commands remain in the queue
//   STATE_LOAD  - command loaded and ready to run
//   STATE_ERROR - player command has an error
func (e *Engine) load_command(c *command) bool {
	if c == nil {
		return false
	}

	line, ok := e.get_command(c.who)
	if !ok {
		e.set_state(c, STATE_DONE, 0)
		return false
	}

	// Parse the command line
	if e.oly_parse(c, line) {
		// Get priority from command table
		pri := e.cmd_pri(c.cmd)
		e.set_state(c, STATE_LOAD, pri)

		c.pri = schar(pri)
		c.wait = e.cmd_time(c.cmd)
		c.poll = schar(e.cmd_poll(c.cmd))
		c.days_executing = 0
	} else {
		e.set_state(c, STATE_ERROR, 0)
	}

	return true
}

// oly_parse parses a command line and populates the command structure.
// This is a stub that will be fully implemented when command parsing is ported.
// For now, it does basic parsing to extract the command name.
// Port of C oly_parse().
func (e *Engine) oly_parse(c *command, line string) bool {
	if c == nil || line == "" {
		return false
	}

	c.line = strToCharPtr(line)

	// Parse the line into arguments
	args := parse_line(line)
	if len(args) == 0 {
		return false
	}

	// Look up the command
	cmdIndex := e.find_command(args[0])
	if cmdIndex < 0 {
		return false
	}

	c.cmd = cmdIndex
	// Store parsed arguments - will be expanded in later sprints
	return true
}

// strToCharPtr converts a Go string to a *char (for compatibility with C port).
// In Go, we just return a pointer to the first byte as schar.
func strToCharPtr(s string) *char {
	if s == "" {
		return nil
	}
	b := []byte(s)
	c := char(b[0])
	return &c
}

// parse_line splits a command line into arguments, respecting quotes.
// Port of C parse_line() from input.c.
func parse_line(line string) []string {
	var result []string
	i := 0
	n := len(line)

	for i < n {
		// Skip whitespace
		for i < n && iswhite(line[i]) {
			i++
		}
		if i >= n {
			break
		}

		var arg string
		if line[i] == '"' {
			// Double-quoted string
			i++
			start := i
			for i < n && line[i] != '"' {
				i++
			}
			arg = line[start:i]
			if i < n {
				i++ // skip closing quote
			}
		} else if line[i] == '\'' {
			// Single-quoted string
			i++
			start := i
			for i < n && line[i] != '\'' {
				i++
			}
			arg = line[start:i]
			if i < n {
				i++ // skip closing quote
			}
		} else {
			// Unquoted argument
			start := i
			for i < n && !iswhite(line[i]) {
				i++
			}
			arg = line[start:i]
		}

		// Trim whitespace from argument
		arg = trimWhitespace(arg)
		if arg != "" {
			result = append(result, arg)
		}
	}

	return result
}

// trimWhitespace removes leading and trailing whitespace from a string.
func trimWhitespace(s string) string {
	start := 0
	end := len(s)
	for start < end && iswhite(s[start]) {
		start++
	}
	for end > start && iswhite(s[end-1]) {
		end--
	}
	return s[start:end]
}

// find_command looks up a command by name and returns its index.
// Returns -1 if not found.
// This is a stub that will be fully implemented when the command table is ported.
func (e *Engine) find_command(name string) int {
	// For now, just return 0 for any non-empty command name
	// The actual command table will be implemented in later sprints
	if name == "" {
		return -1
	}
	// Check for "stop" command which is special
	if i_strcmp(name, "stop") == 0 {
		return 0 // cmd_stop = 0
	}
	// Return a placeholder command index
	return 1
}

// cmd_pri returns the priority for a command from the command table.
// Priority 0 = highest, 4 = lowest.
// This is a stub that will be fully implemented when the command table is ported.
func (e *Engine) cmd_pri(cmdIndex int) int {
	// Default priority is 2 (middle)
	return 2
}

// cmd_time returns the execution time for a command from the command table.
// This is a stub that will be fully implemented when the command table is ported.
func (e *Engine) cmd_time(cmdIndex int) int {
	// Default execution time is 0 (immediate)
	return 0
}

// cmd_poll returns whether a command should be polled each day.
// This is a stub that will be fully implemented when the command table is ported.
func (e *Engine) cmd_poll(cmdIndex int) int {
	return 0
}

// command_done marks a command as done and loads the next one.
// Port of C command_done() from input.c.
func (e *Engine) commandDone(c *command) {
	if c == nil {
		return
	}

	// In immediate mode, just set state to done
	if e.globals.immediate {
		e.set_state(c, STATE_DONE, 0)
		return
	}

	// Load the next command
	if e.load_command(c) {
		// If new command has higher priority, update cur_pri
		if e.globals.cmdQueues != nil && int(c.pri) < e.globals.cmdQueues.curPri {
			e.globals.cmdQueues.curPri = int(c.pri)
		}
	}
}

// init_load_sup initializes command loading for a single entity.
// Port of C init_load_sup() from input.c.
func (e *Engine) init_load_sup(who int) {
	c := e.rp_command(who)

	// All characters should have a command structure.
	// Create one if they don't.
	if c == nil {
		c = e.p_command(who)
		if c == nil {
			return
		}
		c.who = who
		c.state = STATE_DONE
	}

	queues := e.globals.cmdQueues
	if queues == nil {
		e.initCommandQueues()
		queues = e.globals.cmdQueues
	}

	switch c.state {
	case STATE_LOAD:
		if c.pri >= 0 && int(c.pri) < MAX_PRI {
			queues.loadQ[c.pri] = append(queues.loadQ[c.pri], c.who)
		}
	case STATE_RUN:
		queues.runQ = append(queues.runQ, c.who)
	case STATE_DONE:
		e.load_command(c)
	}
}

// initialCommandLoad loads initial commands for all characters and players.
// Port of C initial_command_load() from input.c.
// This replaces the stub in day.go.
func (e *Engine) initialCommandLoadImpl() {
	e.initCommandQueues()

	// Load commands for all characters
	for _, who := range e.Characters() {
		e.init_load_sup(who)
	}

	// Load commands for all players (player entities can have commands too)
	for _, who := range e.Players() {
		e.init_load_sup(who)
	}
}

// min_pri_ready returns the minimum priority that has a ready command.
// Returns 99 if no commands are ready.
// Port of C min_pri_ready() from input.c.
func (e *Engine) min_pri_ready() int {
	queues := e.globals.cmdQueues
	if queues == nil {
		return 99
	}

	for pri := 0; pri < MAX_PRI; pri++ {
		for _, who := range queues.loadQ[pri] {
			c := e.rp_command(who)
			if c == nil {
				continue
			}
			if c.state != STATE_LOAD {
				continue
			}
			if c.pri != schar(pri) {
				continue
			}
			// Check if unit can execute
			if e.is_prisoner(who) {
				continue
			}
			if e.char_moving(who) {
				continue
			}
			if c.second_wait != 0 {
				continue
			}
			return pri
		}
	}

	return 99
}

// is_prisoner returns true if the character is a prisoner.
// This is a stub that will be fully implemented when combat/capture is ported.
func (e *Engine) is_prisoner(who int) bool {
	if e.globals.bx[who] == nil {
		return false
	}
	ch := e.globals.bx[who].x_char
	if ch == nil {
		return false
	}
	return ch.prisoner != 0
}

// start_phase executes loaded commands in priority order.
// Port of C start_phase() from input.c.
func (e *Engine) start_phase() {
	queues := e.globals.cmdQueues
	if queues == nil {
		return
	}

	for {
		pri := e.min_pri_ready()
		queues.curPri = pri

		// Check for auto-attacks at priority 3
		if e.globals.autoAttackFlag && pri >= 3 {
			e.checkAllAutoAttacks()
			e.globals.autoAttackFlag = false
		}

		if pri == 99 {
			return
		}

		// Make a copy of the load queue at this priority
		l := make([]int, len(queues.loadQ[pri]))
		copy(l, queues.loadQ[pri])

		for _, who := range l {
			c := e.rp_command(who)
			if c == nil {
				continue
			}

			if c.state == STATE_LOAD &&
				int(c.pri) == pri &&
				!e.is_prisoner(who) &&
				!e.char_moving(who) &&
				c.second_wait == 0 {
				e.do_command(c)
				e.checkAllWaits()

				if pri != queues.curPri {
					break
				}
			}
		}
	}
}

// do_command executes a single command.
// Port of C do_command() from input.c.
// This is a stub that will be fully implemented when command execution is ported.
func (e *Engine) do_command(c *command) {
	if c == nil {
		return
	}

	if c.state == STATE_ERROR {
		// Report error to player
		c.status = FALSE
		e.commandDone(c)
		return
	}

	// Check if command is allowed - stub for now
	// Transition to running state
	e.set_state(c, STATE_RUN, 0)
	c.days_executing = 0
	c.inhibit_finish = 0
	c.status = TRUE // Assume success for now

	// If command has no wait time, finish immediately
	if c.wait == 0 && c.state == STATE_RUN {
		e.finish_command(c)
	}
}

// finish_command completes a running command.
// Port of C finish_command() from input.c.
func (e *Engine) finish_command(c *command) bool {
	if c == nil {
		return false
	}

	if e.Kind(c.who) == T_deadchar {
		e.commandDone(c)
		return false
	}

	// Characters stacked under moving units have commands suspended
	if e.char_gone(c.who) && e.stack_leader(c.who) != c.who {
		// Suspend unless it's a wait command
		return true
	}

	if c.wait > 0 {
		c.wait--
	}

	// Call finish routine when done or if polling
	if c.wait <= 0 || c.poll != 0 {
		// Finish routine would be called here
		c.status = TRUE
	}

	if c.state == STATE_RUN && (c.status == FALSE || c.wait == 0) {
		e.commandDone(c)
	}

	return c.status != 0
}

// evening_phase processes running commands at end of day.
// Port of C evening_phase() from input.c.
func (e *Engine) evening_phase() {
	e.globals.evening = true

	queues := e.globals.cmdQueues
	if queues == nil {
		return
	}

	// Make a copy of the run queue
	l := make([]int, len(queues.runQ))
	copy(l, queues.runQ)

	for _, who := range l {
		c := e.rp_command(who)
		if c == nil {
			continue
		}

		if c.state != STATE_RUN {
			continue
		}

		if c.second_wait != 0 {
			continue
		}

		c.days_executing++
		e.finish_command(c)
	}

	e.globals.evening = false
}

// checkAllWaits checks all waiting commands.
// Port of C check_all_waits() from input.c.
// This is a stub that will be fully implemented when wait commands are ported.
func (e *Engine) checkAllWaits() {
	// Will be implemented with wait command handling
}

// checkAllAutoAttacks checks for automatic attacks.
// This is a stub that will be fully implemented when combat is ported.
func (e *Engine) checkAllAutoAttacks() {
	// Will be implemented with combat
}

// stack_leader returns the stack leader for a character.
// This is a stub - the full implementation exists in stack.go.
func (e *Engine) stack_leader(who int) int {
	// Return self for now - full implementation in stack.go
	return who
}

// char_moving returns the movement timestamp for a character (Engine method wrapper).
func (e *Engine) char_moving(who int) bool {
	return char_moving(who) != 0
}

// char_gone returns true if the character is currently moving (Engine method wrapper).
func (e *Engine) char_gone(who int) bool {
	return char_gone(who) != 0
}
