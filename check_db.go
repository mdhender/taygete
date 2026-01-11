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

// check_db.go - Database consistency checks ported from src/check.c

package taygete

import (
	"fmt"
	"log/slog"
)

// CheckIssue represents a database integrity issue found during checking.
type CheckIssue struct {
	Type    string // category of issue: "error", "warning", "repaired"
	Message string // description of the issue
}

// CheckResult holds the results of a database consistency check.
type CheckResult struct {
	Issues []CheckIssue
}

// AddError adds an error issue to the result.
func (r *CheckResult) AddError(format string, args ...any) {
	r.Issues = append(r.Issues, CheckIssue{
		Type:    "error",
		Message: fmt.Sprintf(format, args...),
	})
}

// AddWarning adds a warning issue to the result.
func (r *CheckResult) AddWarning(format string, args ...any) {
	r.Issues = append(r.Issues, CheckIssue{
		Type:    "warning",
		Message: fmt.Sprintf(format, args...),
	})
}

// AddRepaired adds a repaired issue to the result.
func (r *CheckResult) AddRepaired(format string, args ...any) {
	r.Issues = append(r.Issues, CheckIssue{
		Type:    "repaired",
		Message: fmt.Sprintf(format, args...),
	})
}

// HasErrors returns true if any error issues were found.
func (r *CheckResult) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Type == "error" {
			return true
		}
	}
	return false
}

// ErrorCount returns the number of error issues.
func (r *CheckResult) ErrorCount() int {
	count := 0
	for _, issue := range r.Issues {
		if issue.Type == "error" {
			count++
		}
	}
	return count
}

// WarningCount returns the number of warning issues.
func (r *CheckResult) WarningCount() int {
	count := 0
	for _, issue := range r.Issues {
		if issue.Type == "warning" {
			count++
		}
	}
	return count
}

// RepairedCount returns the number of repaired issues.
func (r *CheckResult) RepairedCount() int {
	count := 0
	for _, issue := range r.Issues {
		if issue.Type == "repaired" {
			count++
		}
	}
	return count
}

// CheckDB checks database integrity and effects minor repairs.
// Fixes minor problems in backlinks and lists.
// Always notes database corrections in the returned CheckResult.
// Ported from src/check.c check_db().
func (e *Engine) CheckDB() *CheckResult {
	result := &CheckResult{}

	e.checkGlob(result)
	e.checkHere(result)
	e.checkSwear(result)
	e.checkIndep(result)
	e.checkGM(result)
	e.checkSkillPlayer(result)
	e.checkEatPlayer(result)
	e.checkNPCPlayer(result)
	e.checkGarrPlayer(result)
	e.checkNowhere(result)
	e.checkSkills(result)
	e.checkItemCounts(result)
	e.checkLocNameLengths(result)
	e.checkMoving(result)
	e.checkPrisoner(result)

	if e.globals.bx[e.globals.garrison_magic] != nil {
		result.AddWarning("%s should not be allocated, reserved for garrison_magic",
			box_code(e.globals.garrison_magic))
	}

	return result
}

// checkGlob verifies that T_MAX and SUB_MAX match the string tables.
// Ported from src/check.c check_glob().
func (e *Engine) checkGlob(result *CheckResult) {
	if len(kind_s) != T_MAX {
		result.AddError("kind_s length %d != T_MAX %d", len(kind_s), T_MAX)
	}

	if len(subkind_s) != SUB_MAX {
		result.AddError("subkind_s length %d != SUB_MAX %d", len(subkind_s), SUB_MAX)
	}
}

// checkHere verifies and repairs here_list consistency.
// 1. If box claims it's in a location but not in here_list, add it.
// 2. If box is in here_list but claims different location, remove it.
// Ported from src/check.c check_here().
func (e *Engine) checkHere(result *CheckResult) {
	for i := 1; i < MAX_BOXES; i++ {
		if e.globals.bx[i] == nil {
			continue
		}

		where := loc(i)
		if where > 0 && !in_here_list(where, i) {
			result.AddRepaired("adding [%d] to here list of [%d]", i, where)
			add_to_here_list(where, i)
		}
	}

	for i := 1; i < MAX_BOXES; i++ {
		if e.globals.bx[i] == nil {
			continue
		}

		li := rp_loc_info(i)
		if li == nil {
			continue
		}

		toRemove := []int{}
		for _, j := range li.here_list {
			where := loc(j)
			if where != i {
				result.AddRepaired("removing [%d] from here list of [%d]", j, i)
				toRemove = append(toRemove, j)
			}
		}

		for _, j := range toRemove {
			IListRemValue(&li.here_list, j)
		}
	}
}

// checkSwear verifies and repairs player unit list consistency.
// 1. If char claims to be in a faction, ensure it's in faction's unit list.
// 2. If char is in faction's unit list but claims different faction, remove it.
// Ported from src/check.c check_swear().
func (e *Engine) checkSwear(result *CheckResult) {
	for _, i := range e.Characters() {
		over := player(i)
		if over > 0 && !e.isUnit(over, i) {
			result.AddRepaired("adding [%d] to player [%d]", i, over)
			e.addUnit(over, i)
		}
	}

	for _, i := range e.Players() {
		p := rp_player(i)
		if p == nil {
			continue
		}

		toRemove := []int{}
		units := e.getPlayerUnits(i)
		for _, j := range units {
			over := player(j)
			if over != i {
				result.AddRepaired("removing [%d] from player list of [%d]", j, i)
				toRemove = append(toRemove, j)
			}
		}

		for _, j := range toRemove {
			e.removeUnit(i, j)
		}
	}
}

// checkIndep ensures the independent player exists and all orphan chars are assigned.
// Ported from src/check.c check_indep().
func (e *Engine) checkIndep(result *CheckResult) {
	if e.globals.bx[indep_player] == nil {
		result.AddRepaired("creating independent player [%d]", indep_player)
		alloc_box(indep_player, T_player, sub_pl_npc)
	}

	if kind(indep_player) != T_player {
		result.AddError("indep_player [%d] is not T_player", indep_player)
		return
	}

	if name(indep_player) == "" {
		set_name(indep_player, "Independent player")
	}

	for _, i := range e.Characters() {
		if player(i) == 0 {
			result.AddRepaired("swearing unit [%d] to %s", i, box_name(indep_player))
			set_lord(i, indep_player, LOY_unsworn, 0)
		}
	}
}

// checkGM ensures the gamemaster player exists.
// Ported from src/check.c check_gm().
func (e *Engine) checkGM(result *CheckResult) {
	if e.globals.bx[gm_player] == nil {
		result.AddRepaired("creating gm player [%d]", gm_player)
		alloc_box(gm_player, T_player, sub_pl_system)
	}

	if kind(gm_player) != T_player {
		result.AddError("gm_player [%d] is not T_player", gm_player)
		return
	}

	if name(gm_player) == "" {
		set_name(gm_player, "Gamemaster")
	}
}

// checkSkillPlayer ensures the skill player exists.
// Ported from src/check.c check_skill_player().
func (e *Engine) checkSkillPlayer(result *CheckResult) {
	if e.globals.bx[skill_player] == nil {
		result.AddRepaired("creating skill player [%d]", skill_player)
		alloc_box(skill_player, T_player, sub_pl_system)
	}

	if kind(skill_player) != T_player {
		result.AddError("skill_player [%d] is not T_player", skill_player)
		return
	}

	if name(skill_player) == "" {
		set_name(skill_player, "Skill list")
	}
}

// checkEatPlayer ensures the order eater player exists.
// Ported from src/check.c check_eat_player().
func (e *Engine) checkEatPlayer(result *CheckResult) {
	if e.globals.bx[eat_pl] == nil {
		result.AddRepaired("creating eat player [%d]", eat_pl)
		alloc_box(eat_pl, T_player, sub_pl_system)
	}

	if kind(eat_pl) != T_player {
		result.AddError("eat_pl [%d] is not T_player", eat_pl)
		return
	}

	if name(eat_pl) == "" {
		set_name(eat_pl, "Order eater")
	}
}

// checkNPCPlayer ensures the NPC control player exists.
// Ported from src/check.c check_npc_player().
func (e *Engine) checkNPCPlayer(result *CheckResult) {
	if e.globals.bx[npc_pl] == nil {
		result.AddRepaired("creating npc player [%d]", npc_pl)
		alloc_box(npc_pl, T_player, sub_pl_silent)
	}

	if kind(npc_pl) != T_player {
		result.AddError("npc_pl [%d] is not T_player", npc_pl)
		return
	}

	if name(npc_pl) == "" {
		set_name(npc_pl, "NPC control")
	}
}

// checkGarrPlayer ensures the garrison player exists.
// Ported from src/check.c check_garr_player().
func (e *Engine) checkGarrPlayer(result *CheckResult) {
	if e.globals.bx[garr_pl] == nil {
		result.AddRepaired("creating garrison player [%d]", garr_pl)
		alloc_box(garr_pl, T_player, sub_pl_silent)
	}

	if kind(garr_pl) != T_player {
		result.AddError("garr_pl [%d] is not T_player", garr_pl)
		return
	}

	if name(garr_pl) == "" {
		set_name(garr_pl, "Garrison units")
	}
}

// checkNowhere warns about units and locations that have no parent location.
// Ported from src/check.c check_nowhere().
func (e *Engine) checkNowhere(result *CheckResult) {
	for _, i := range e.Characters() {
		if loc(i) == 0 {
			result.AddWarning("unit %s is nowhere", box_code(i))
		}
	}

	for _, i := range e.LocsAndShips() {
		if loc_depth(i) > LOC_region && loc(i) == 0 {
			result.AddWarning("loc %s is nowhere", box_code(i))
		}
	}
}

// checkSkills checks skill tree consistency.
// Verifies that skills have valid learn times and valid skill tree references.
// Ported from src/check.c check_skills().
func (e *Engine) checkSkills(result *CheckResult) {
	for _, sk := range e.Skills() {
		if learn_time(sk) == 0 {
			result.AddWarning("learn time of %s is 0", box_name(sk))
		}

		s := rp_skill(sk)
		if s == nil {
			continue
		}

		for _, child := range s.offered {
			if kind(child) != T_skill {
				result.AddError("skill %s offered list contains non-skill %s", box_name(sk), box_code(child))
			}
		}

		for _, child := range s.research {
			if kind(child) != T_skill {
				result.AddError("skill %s research list contains non-skill %s", box_name(sk), box_code(child))
			}
		}
	}
}

// checkItemCounts verifies unique item counts and ownership.
// Ported from src/check.c check_item_counts().
func (e *Engine) checkItemCounts(result *CheckResult) {
	for _, i := range e.Items() {
		e.globals.bx[i].temp = 0
	}

	for i := 1; i < MAX_BOXES; i++ {
		if e.globals.bx[i] == nil {
			continue
		}

		inv := e.getInventory(i)
		for _, ent := range inv {
			if kind(ent.item) != T_item {
				result.AddError("%s has non-item %s", box_name(i), box_name(ent.item))
				continue
			}

			if item_unique(ent.item) == 0 {
				continue
			}

			if item_unique(ent.item) != i {
				result.AddRepaired("unique item %s: whohas=%s, actual=%s",
					box_name(ent.item), box_name(item_unique(ent.item)), box_name(i))
				p_item(ent.item).who_has = i
			}

			if ent.qty != 1 {
				result.AddError("%s has qty %d of unique item %s",
					box_name(i), ent.qty, box_name(ent.item))
			}

			if e.globals.bx[ent.item] != nil {
				e.globals.bx[ent.item].temp += ent.qty
			}
		}
	}

	for _, i := range e.Items() {
		if item_unique(i) != 0 {
			if e.globals.bx[i].temp != 1 {
				result.AddError("unique item %s count %d", box_name(i), e.globals.bx[i].temp)
			}
		}
	}
}

// checkLocNameLengths warns about location names that are too long.
// Ported from src/check.c check_loc_name_lengths().
func (e *Engine) checkLocNameLengths(result *CheckResult) {
	for _, i := range e.Locations() {
		n := name(i)
		if len(n) > 25 {
			result.AddWarning("%s name too long (%d chars)", box_name(i), len(n))
		}
	}
}

// checkMoving verifies that moving characters have valid commands.
// Ported from src/check.c check_moving().
func (e *Engine) checkMoving(result *CheckResult) {
	for _, i := range e.Characters() {
		if stack_leader(i) != i || char_moving(i) == 0 {
			continue
		}

		c := rp_command(i)

		if c == nil || c.state != STATE_RUN {
			result.AddRepaired("%s moving but no command", box_name(i))
			restore_stack_actions(i)
		}
	}

	for _, i := range e.Characters() {
		leader := stack_leader(i)

		if leader == i || char_moving(i) == char_moving(leader) {
			continue
		}

		result.AddRepaired("%s moving disagrees with leader", box_name(i))
		p_char(i).moving = char_moving(leader)
	}
}

// checkPrisoner verifies that prisoners are properly stacked.
// Ported from src/check.c check_prisoner().
func (e *Engine) checkPrisoner(result *CheckResult) {
	for _, who := range e.Characters() {
		if !is_prisoner(who) {
			continue
		}

		if stack_parent(who) == 0 {
			result.AddRepaired("%s prisoner but unstacked", box_name(who))
			p_char(who).prisoner = FALSE
		}
	}
}

// isUnit returns true if char is in player's unit list.
func (e *Engine) isUnit(pl, char int) bool {
	units := e.getPlayerUnits(pl)
	return IListLookup(units, char) != -1
}

// addUnit adds char to player's unit list.
func (e *Engine) addUnit(pl, char int) {
	if e.globals.playerUnits == nil {
		e.globals.playerUnits = make(map[int][]int)
	}
	units := e.globals.playerUnits[pl]
	IListAppend(&units, char)
	IListSort(units)
	e.globals.playerUnits[pl] = units
}

// removeUnit removes char from player's unit list.
func (e *Engine) removeUnit(pl, char int) {
	if e.globals.playerUnits == nil {
		return
	}
	units := e.globals.playerUnits[pl]
	IListRemValue(&units, char)
	e.globals.playerUnits[pl] = units
}

// getPlayerUnits returns the units list for a player.
func (e *Engine) getPlayerUnits(pl int) []int {
	if e.globals.playerUnits == nil {
		return nil
	}
	return e.globals.playerUnits[pl]
}

// getInventory returns the inventory for an entity.
// Returns empty slice if entity has no inventory.
func (e *Engine) getInventory(n int) []item_ent {
	if e.globals.inventories == nil {
		return nil
	}
	return e.globals.inventories[n]
}

// LogCheckResult logs all issues in a CheckResult using slog.
func LogCheckResult(logger *slog.Logger, result *CheckResult) {
	for _, issue := range result.Issues {
		switch issue.Type {
		case "error":
			logger.Error("check_db", "issue", issue.Message)
		case "warning":
			logger.Warn("check_db", "issue", issue.Message)
		case "repaired":
			logger.Info("check_db", "repaired", issue.Message)
		}
	}
}
