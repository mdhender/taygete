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

// cmd_wait.go - WAIT/FLAG system ported from src/c1.c
// Sprint 26.7: WAIT/FLAG System

package taygete

import (
	"fmt"
	"unicode"
)

// flag_ent represents a flag signal raised by a character.
// Ported from src/c1.c lines 1083-1086.
type flag_ent struct {
	who  int    // who raised the flag
	flag string // the flag message
}

// flags is the global list of raised flags for the current turn.
// This list is cleared at the end of each turn.
// Ported from src/c1.c line 1088.
var flags []*flag_ent

// wait_list is the global list of units currently waiting.
// Used to check wait conditions each day.
// Ported from src/c1.c line 1579.
var wait_list []int

// waitArgExt extends wait_arg with string flag support for Go.
// Uses the existing wait_arg type from types.go but stores flag as string.
type waitArgExt struct {
	tag     int    // condition tag (index into wait_tags)
	a1      int    // first argument
	a2      int    // second argument (for item qty)
	flagStr string // flag string (for flag condition) - Go string version
}

// Wait condition tag constants.
// Ported from src/c1.c lines 1141-1168 (wait_tags array).
const (
	waitTagTime    = 0  // time n - wait n days
	waitTagDay     = 1  // day n - wait until day n
	waitTagUnit    = 2  // unit n - wait for unit to arrive
	waitTagGold    = 3  // gold n - wait until have n gold
	waitTagItem    = 4  // item n q - wait until have q of item n
	waitTagFlag    = 5  // flag f [n] - wait for flag signal
	waitTagLoc     = 6  // loc n - wait until at location n
	waitTagStack   = 7  // stack n - wait until stacked with n
	waitTagTop     = 8  // top - wait until stack leader
	waitTagFerry   = 9  // ferry n - wait for ferry signal
	waitTagShip    = 10 // ship n - wait for ship at location
	waitTagRain    = 11 // rain - wait for rain
	waitTagFog     = 12 // fog - wait for fog
	waitTagWind    = 13 // wind - wait for wind
	waitTagNot     = 14 // not - negate next condition
	waitTagOwner   = 15 // owner - wait until first character in loc
	waitTagRaining = 16 // raining -> rain (alias)
	waitTagFoggy   = 17 // foggy -> fog (alias)
	waitTagWindy   = 18 // windy -> wind (alias)
	waitTagClear   = 19 // clear - wait for clear weather
	waitTagShiploc = 20 // shiploc n - wait for ship location
	waitTagMonth   = 21 // month n - wait until turn n
	waitTagTurn    = 22 // turn n - wait until turn n
)

// wait_tags maps wait condition keywords to tag numbers.
// Ported from src/c1.c lines 1141-1168.
var wait_tags = []string{
	"time",    // 0
	"day",     // 1
	"unit",    // 2
	"gold",    // 3
	"item",    // 4
	"flag",    // 5
	"loc",     // 6
	"stack",   // 7
	"top",     // 8
	"ferry",   // 9
	"ship",    // 10
	"rain",    // 11
	"fog",     // 12
	"wind",    // 13
	"not",     // 14
	"owner",   // 15
	"raining", // 16 -> 11
	"foggy",   // 17 -> 12
	"windy",   // 18 -> 13
	"clear",   // 19
	"shiploc", // 20
	"month",   // 21
	"turn",    // 22
}

// flag_raised checks if a flag has been raised by the given who (or anyone if who=0).
// Returns the index into the flags slice, or -1 if not found.
// Ported from src/c1.c lines 1091-1108.
func flag_raised(who int, flag string) int {
	for i, f := range flags {
		if who != 0 && player(f.who) != who && f.who != who {
			continue
		}
		if i_strcmp(f.flag, flag) == 0 {
			return i
		}
	}
	return -1
}

// v_flag is the start routine for the FLAG command.
// Raises a flag signal that other units can wait for.
// Usage: FLAG <message>
// Ported from src/c1.c lines 1111-1138.
func v_flag(c *command) int {
	if numargs(c) < 1 {
		wout(c.who, "Must specify what message to signal.")
		return FALSE
	}

	flag := get_wait_parse_arg(c, 1)

	if flag_raised(c.who, flag) >= 0 {
		wout(c.who, "%s has already given that signal this month.", box_name(c.who))
		return FALSE
	}

	newFlag := &flag_ent{
		who:  c.who,
		flag: flag,
	}

	flags = append(flags, newFlag)

	return TRUE
}

// clear_flags clears all raised flags (called at end of turn).
func clear_flags() {
	flags = nil
}

// clear_wait_parse clears the parsed wait arguments from a command.
// Ported from src/c1.c lines 1172-1184.
func clear_wait_parse(c *command) {
	c.wait_parse = nil
}

// parse_wait_args parses the WAIT command arguments into wait_arg structures.
// Returns an error message string if parsing fails, or empty string on success.
// Ported from src/c1.c lines 1187-1283.
func parse_wait_args(c *command) string {
	// Clear any existing parsed args
	clearWaitParseList(c)

	args := getCommandParseArgs(c)
	i := 1 // skip command name at index 0

	for i < len(args) {
		tagStr := args[i]
		tag := lookup(wait_tags, tagStr)

		// Map aliases
		switch tag {
		case 16: // raining -> rain
			tag = 11
		case 17: // foggy -> fog
			tag = 12
		case 18: // windy -> wind
			tag = 13
		}

		if tag < 0 {
			return fmt.Sprintf("Unknown condition '%s'.", tagStr)
		}
		i++

		newArg := &waitArgExt{
			tag: tag,
		}

		switch tag {
		case waitTagTime, waitTagDay, waitTagMonth, waitTagTurn,
			waitTagUnit, waitTagGold, waitTagLoc, waitTagShiploc,
			waitTagStack, waitTagFerry, waitTagShip:
			// These require one argument
			if i < len(args) {
				newArg.a1 = parse_wait_arg_value(c.who, args[i])
				i++
			} else {
				return fmt.Sprintf("Argument missing for '%s'.", tagStr)
			}

		case waitTagItem:
			// Requires item and optional quantity
			if i < len(args) {
				newArg.a1 = parse_wait_arg_value(c.who, args[i])
				i++
			} else {
				return fmt.Sprintf("Argument missing for '%s'.", tagStr)
			}

			if i < len(args) {
				newArg.a2 = parse_wait_arg_value(c.who, args[i])
				i++
			} else {
				newArg.a2 = 1 // default quantity
			}

		case waitTagFlag:
			// Requires flag string and optional who
			if i < len(args) {
				newArg.flagStr = args[i]
				i++
			} else {
				return "Flag missing."
			}

			newArg.a1 = player(c.who) // default to own faction

			if i < len(args) {
				// Check if next arg is numeric or a valid ID
				nextArg := args[i]
				if len(nextArg) > 0 && (unicode.IsDigit(rune(nextArg[0])) || parse_wait_arg_value(c.who, nextArg) != 0) {
					newArg.a1 = parse_wait_arg_value(c.who, nextArg)
					i++
				}
			}

		case waitTagTop, waitTagRain, waitTagFog, waitTagWind, waitTagNot, waitTagOwner, waitTagClear:
			// No arguments needed
		}

		appendWaitParse(c, newArg)
	}

	return ""
}

// check_wait_conditions evaluates all WAIT conditions for a command.
// Returns a status message if a condition is met, or empty string if still waiting.
// Ported from src/c1.c lines 1286-1576.
func check_wait_conditions(c *command) string {
	whereShip := subloc(c.who)
	if is_ship_either(whereShip) {
		whereShip = subloc(whereShip)
	}

	waitArgs := getWaitParse(c)
	if len(waitArgs) < 1 {
		ret := parse_wait_args(c)
		if ret != "" {
			return ret
		}
		waitArgs = getWaitParse(c)
		if len(waitArgs) == 0 {
			return "No wait conditions specified."
		}
	}

	not := false
	setNot := false

	for _, p := range waitArgs {
		if setNot {
			setNot = false
		} else if not {
			not = false
		}

		var cond bool

		switch p.tag {
		case waitTagTime: // time n - wait n days
			cond = (command_days(c) >= p.a1)
			cond = evalNotCond(not, cond)
			if cond {
				if not {
					return fmt.Sprintf("%s day%s have not passed.", nice_num(p.a1), pluralS(p.a1))
				}
				return fmt.Sprintf("%s day%s passed.", nice_num(p.a1), pluralHaveHas(p.a1))
			}

		case waitTagDay: // day n - wait until day n
			cond = (int(teg.globals.sysclock.day) >= p.a1)
			cond = evalNotCond(not, cond)
			if cond {
				if not {
					return fmt.Sprintf("today is not day %d.", p.a1)
				}
				return fmt.Sprintf("today is day %d.", int(teg.globals.sysclock.day))
			}

		case waitTagMonth, waitTagTurn: // month/turn n - wait until turn n
			cond = (int(teg.globals.sysclock.turn) >= p.a1)
			cond = evalNotCond(not, cond)
			if cond {
				if not {
					return fmt.Sprintf("it is not turn %d.", p.a1)
				}
				return fmt.Sprintf("it is turn %d.", int(teg.globals.sysclock.turn))
			}

		case waitTagUnit: // unit n - wait for unit to arrive
			if !valid_box(p.a1) {
				return fmt.Sprintf("%s does not exist.", box_code(p.a1))
			}
			cond = char_here(c.who, p.a1)
			cond = evalNotCond(not, cond)
			if cond {
				return fmt.Sprintf("%s is%s here.", box_code(p.a1), notStr(not))
			}

		case waitTagGold: // gold n - wait until have n gold
			cond = (has_item(c.who, item_gold) >= p.a1)
			cond = evalNotCond(not, cond)
			if cond {
				if not {
					return fmt.Sprintf("%s doesn't have %s.", just_name(c.who), gold_s(p.a1))
				}
				return fmt.Sprintf("%s has %s.", just_name(c.who), gold_s(has_item(c.who, item_gold)))
			}

		case waitTagItem: // item n q - wait until have q of item n
			cond = (kind(p.a1) == T_item && has_item(c.who, p.a1) >= p.a2)
			cond = evalNotCond(not, cond)
			if cond {
				if not {
					return fmt.Sprintf("%s doesn't have %s.", just_name(c.who), just_name_qty(p.a1, p.a2))
				}
				return fmt.Sprintf("%s has %s.", just_name(c.who), just_name_qty(p.a1, has_item(c.who, p.a1)))
			}

		case waitTagFlag: // flag - wait for flag signal
			if p.a1 != 0 && !valid_box(p.a1) {
				return fmt.Sprintf("%s does not exist.", box_code(p.a1))
			}

			j := flag_raised(p.a1, p.flagStr)

			if not {
				if j < 0 {
					return "received no signal"
				}
			} else {
				if j >= 0 {
					return fmt.Sprintf("%s signaled '%s'", box_name(flags[j].who), flags[j].flag)
				}
			}

		case waitTagLoc: // loc n - wait until at location n
			if !is_loc_or_ship(p.a1) {
				return fmt.Sprintf("%s is not a location or ship.", box_code(p.a1))
			}
			cond = (subloc(c.who) == p.a1)
			cond = evalNotCond(not, cond)
			if cond {
				return fmt.Sprintf("%sat %s.", notStr2(not), box_name(p.a1))
			}

		case waitTagShiploc: // shiploc n - wait for ship location
			ship := subloc(c.who)
			if !is_ship(ship) && !is_ship_notdone(ship) {
				return fmt.Sprintf("%s is not on a ship.", box_name(c.who))
			}
			if !is_loc_or_ship(p.a1) {
				return fmt.Sprintf("%s is not a location or ship.", box_code(p.a1))
			}
			where := subloc(ship)
			cond = (where == p.a1)
			cond = evalNotCond(not, cond)
			if cond {
				return fmt.Sprintf("%sat %s.", notStr2(not), box_name(p.a1))
			}

		case waitTagStack: // stack n - wait until stacked with n
			if kind(p.a1) != T_char {
				break // just hang (don't error)
			}
			cond = (stack_leader(c.who) == stack_leader(p.a1))
			cond = evalNotCond(not, cond)
			if cond {
				return fmt.Sprintf("%s is%s stacked with us.", box_name(p.a1), notStr(not))
			}

		case waitTagTop: // top - wait until stack leader
			cond = (stack_leader(c.who) == c.who)
			cond = evalNotCond(not, cond)
			if cond {
				return fmt.Sprintf("we are%s the stack leader", notStr(not))
			}

		case waitTagFerry: // ferry n - wait for ferry signal
			if !is_ship(p.a1) {
				return fmt.Sprintf("%s is not a ship", box_code(p.a1))
			}
			cond = (subloc(p.a1) == subloc(c.who) && ferry_horn(p.a1) != 0)
			cond = evalNotCond(not, cond)
			if cond {
				return fmt.Sprintf("the ferry has%s signaled.", notStr(not))
			}

		case waitTagShip: // ship n - wait for ship at location
			if kind(p.a1) != T_ship {
				return fmt.Sprintf("%s is not a ship.", box_code(p.a1))
			}
			cond = (whereShip == subloc(p.a1))
			cond = evalNotCond(not, cond)
			if cond {
				return fmt.Sprintf("%s is%s here.", box_code(p.a1), notStr(not))
			}

		case waitTagRain: // rain - wait for rain
			cond = (weather_here(province(c.who), sub_rain) != 0)
			cond = evalNotCond(not, cond)
			if cond {
				return fmt.Sprintf("it is%s raining.", notStr(not))
			}

		case waitTagFog: // fog - wait for fog
			cond = (weather_here(province(c.who), sub_fog) != 0)
			cond = evalNotCond(not, cond)
			if cond {
				return fmt.Sprintf("it is%s foggy.", notStr(not))
			}

		case waitTagWind: // wind - wait for wind
			cond = (weather_here(province(c.who), sub_wind) != 0)
			cond = evalNotCond(not, cond)
			if cond {
				return fmt.Sprintf("it is%s windy.", notStr(not))
			}

		case waitTagNot: // not - negate next condition
			not = true
			setNot = true

		case waitTagOwner: // owner - wait until first character in loc
			cond = (first_character(subloc(c.who)) == c.who)
			cond = evalNotCond(not, cond)
			if cond {
				return fmt.Sprintf("we are%s the first character here", notStr(not))
			}

		case waitTagClear: // clear - wait for clear weather
			cond = (weather_here(subloc(c.who), sub_fog) == 0 &&
				weather_here(subloc(c.who), sub_rain) == 0 &&
				weather_here(subloc(c.who), sub_wind) == 0)
			cond = evalNotCond(not, cond)
			if cond {
				return fmt.Sprintf("it is%s clear.", notStr(not))
			}
		}
	}

	return ""
}

// v_wait is the start routine for the WAIT command.
// Waits for one or more conditions to be met.
// Usage: WAIT <condition> [args...]
// Ported from src/c1.c lines 1582-1606.
func v_wait(c *command) int {
	if numargs(c) < 1 {
		wout(c.who, "Must say what condition to wait for.")
		return FALSE
	}

	clear_wait_parse(c)

	if s := check_wait_conditions(c); s != "" {
		wout(c.who, "Wait finished: %s", s)

		c.wait = 0
		c.inhibit_finish = TRUE // don't call d_wait
		return TRUE
	}

	IListAppend(&wait_list, c.who)
	return TRUE
}

// d_wait is the finish routine for the WAIT command.
// Called each day to check if wait conditions are met.
// Ported from src/c1.c lines 1609-1625.
func d_wait(c *command) int {
	if s := check_wait_conditions(c); s != "" {
		wout(c.who, "Wait finished: %s", s)
		IListRemValue(&wait_list, c.who)

		c.wait = 0
		c.inhibit_finish = TRUE
		return TRUE
	}

	return TRUE
}

// i_wait is the interrupt routine for the WAIT command.
// Called when the WAIT command is interrupted.
// Ported from src/c1.c lines 1628-1634.
func i_wait(c *command) int {
	IListRemValue(&wait_list, c.who)
	return TRUE
}

// Helper functions for WAIT condition evaluation

// evalNotCond applies the "not" modifier to a condition result.
// Mimics the C logic: cond = not - cond
func evalNotCond(not, cond bool) bool {
	if not {
		return !cond
	}
	return cond
}

// notStr returns " not" if not is true, empty string otherwise.
func notStr(not bool) string {
	if not {
		return " not"
	}
	return ""
}

// notStr2 returns "not " if not is true, empty string otherwise.
func notStr2(not bool) string {
	if not {
		return "not "
	}
	return ""
}

// pluralS returns "s" for plural, empty for singular.
func pluralS(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// pluralHaveHas returns " has" for singular, "s have" for plural.
func pluralHaveHas(n int) string {
	if n == 1 {
		return " has"
	}
	return "s have"
}

// command_days returns the number of days a command has been executing.
func command_days(c *command) int {
	return c.days_executing
}

// parse_wait_arg_value parses a wait argument value (number or entity code).
func parse_wait_arg_value(who int, s string) int {
	// Try parsing as a number first
	var n int
	if _, err := fmt.Sscanf(s, "%d", &n); err == nil {
		return n
	}
	// Try parsing as entity code
	return code_to_int(s)
}

// get_wait_parse_arg returns the parsed argument at index i as a string.
// This is a workaround for accessing parsed command arguments.
func get_wait_parse_arg(c *command, i int) string {
	args := getCommandParseArgs(c)
	if i < len(args) {
		return args[i]
	}
	return ""
}

// getCommandParseArgs returns the parsed arguments for a command.
// Workaround since c.parse is **char in C.
func getCommandParseArgs(c *command) []string {
	// For now, construct from command args field.
	// This will be replaced when full order parsing is implemented.
	var args []string
	args = append(args, "wait") // command name placeholder

	// Build args from command.a, .b, etc based on actual use
	if c.a != 0 {
		args = append(args, fmt.Sprintf("%d", c.a))
	}
	if c.b != 0 {
		args = append(args, fmt.Sprintf("%d", c.b))
	}
	if c.c != 0 {
		args = append(args, fmt.Sprintf("%d", c.c))
	}
	if c.d != 0 {
		args = append(args, fmt.Sprintf("%d", c.d))
	}

	return args
}

// appendWaitParse appends a waitArgExt to the command's wait_parse list.
func appendWaitParse(c *command, arg *waitArgExt) {
	// Workaround: wait_parse is **wait_arg in C
	// We store in engine globals keyed by command pointer
	waitParseLists[c] = append(waitParseLists[c], arg)
}

// getWaitParse returns the wait_parse list for a command.
func getWaitParse(c *command) []*waitArgExt {
	return waitParseLists[c]
}

// waitParseLists stores wait_parse lists keyed by command pointer.
// Workaround since command.wait_parse is **wait_arg (C plist).
var waitParseLists = make(map[*command][]*waitArgExt)

// clearWaitParseLists clears all stored wait parse lists.
// Called at end of turn processing.
func clearWaitParseLists() {
	waitParseLists = make(map[*command][]*waitArgExt)
}

// clearWaitParseList clears the wait parse list for a specific command.
func clearWaitParseList(c *command) {
	delete(waitParseLists, c)
}

// init_wait_list initializes the wait_list at the start of a turn.
// Called from input.c init_wait_list.
func init_wait_list() {
	wait_list = nil

	// Scan all characters for those with running WAIT commands
	for id := teg.KindFirst(T_char); id > 0; id = teg.KindNext(id) {
		c := rp_command(id)
		if c != nil && c.state == STATE_RUN && is_wait_command(c.cmd) {
			IListAppend(&wait_list, id)
		}
	}
}

// check_all_wait_conditions checks wait conditions for all waiting units.
// Called each day during turn processing.
func check_all_wait_conditions() {
	for i := 0; i < len(wait_list); i++ {
		c := rp_command(wait_list[i])
		if c != nil {
			check_wait_conditions(c)
		}
	}
}

// is_wait_command checks if cmd is the WAIT command.
// TODO: Replace with proper command lookup when cmd_tbl is implemented.
func is_wait_command(cmd int) bool {
	// Placeholder - will be implemented when command table is available
	return false
}
