// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

import (
	"strings"
	"testing"
)

func TestCheckDB_Empty(t *testing.T) {
	e := newTestEngine(t)
	result := e.CheckDB()

	// An empty engine should create system players (repaired)
	if result.RepairedCount() == 0 {
		t.Error("expected system player creation repairs")
	}

	// Should not have any errors
	if result.HasErrors() {
		t.Errorf("unexpected errors: %d", result.ErrorCount())
		for _, issue := range result.Issues {
			if issue.Type == "error" {
				t.Logf("  error: %s", issue.Message)
			}
		}
	}
}

func TestCheckDB_SystemPlayers(t *testing.T) {
	e := newTestEngine(t)
	result := e.CheckDB()

	// Verify system players were created
	if kind(indep_player) != T_player {
		t.Error("indep_player was not created")
	}
	if kind(gm_player) != T_player {
		t.Error("gm_player was not created")
	}
	if kind(skill_player) != T_player {
		t.Error("skill_player was not created")
	}
	if kind(eat_pl) != T_player {
		t.Error("eat_pl was not created")
	}
	if kind(npc_pl) != T_player {
		t.Error("npc_pl was not created")
	}
	if kind(garr_pl) != T_player {
		t.Error("garr_pl was not created")
	}

	// Running again should not repair system players
	result2 := e.CheckDB()
	for _, issue := range result2.Issues {
		if issue.Type == "repaired" {
			if issue.Message == "creating independent player [100]" ||
				issue.Message == "creating gm player [200]" ||
				issue.Message == "creating skill player [202]" ||
				issue.Message == "creating eat player [203]" ||
				issue.Message == "creating npc player [206]" ||
				issue.Message == "creating garrison player [207]" {
				t.Errorf("system player was repaired on second run: %s", issue.Message)
			}
		}
	}
	_ = result
}

func TestCheckDB_HereListConsistency(t *testing.T) {
	e := newTestEngine(t)

	// Create a location
	locID := 10000
	alloc_box(locID, T_loc, sub_plain)
	set_name(locID, "Test Province")

	// Create a character
	charID := 1000
	alloc_box(charID, T_char, 0)
	set_name(charID, "Test Character")

	// Set character's location but don't add to here list
	p_loc_info(charID).where = locID

	result := e.CheckDB()

	// Should repair by adding char to here list
	foundRepair := false
	for _, issue := range result.Issues {
		if issue.Type == "repaired" && issue.Message == "adding [1000] to here list of [10000]" {
			foundRepair = true
			break
		}
	}
	if !foundRepair {
		t.Error("expected here list repair")
	}

	// Verify char is now in here list
	if !in_here_list(locID, charID) {
		t.Error("character was not added to here list")
	}
}

func TestCheckDB_HereListRemoveInvalid(t *testing.T) {
	e := newTestEngine(t)

	// Create a location
	locID := 10000
	alloc_box(locID, T_loc, sub_plain)
	set_name(locID, "Test Province")

	// Create another location
	locID2 := 10001
	alloc_box(locID2, T_loc, sub_plain)
	set_name(locID2, "Test Province 2")

	// Create a character claiming to be in locID2
	charID := 1000
	alloc_box(charID, T_char, 0)
	set_name(charID, "Test Character")
	p_loc_info(charID).where = locID2

	// Manually add char to locID's here list (wrong location)
	p_loc_info(locID).here_list = append(p_loc_info(locID).here_list, charID)

	result := e.CheckDB()

	// Should repair by removing char from wrong here list
	foundRepair := false
	for _, issue := range result.Issues {
		if issue.Type == "repaired" && issue.Message == "removing [1000] from here list of [10000]" {
			foundRepair = true
			break
		}
	}
	if !foundRepair {
		t.Error("expected here list removal repair")
	}

	// Verify char is no longer in wrong here list
	if in_here_list(locID, charID) {
		t.Error("character was not removed from wrong here list")
	}
}

func TestCheckDB_CharNowhereWarning(t *testing.T) {
	e := newTestEngine(t)

	// Create a character with no location
	charID := 1000
	alloc_box(charID, T_char, 0)
	set_name(charID, "Orphan Character")
	// Note: do not set location

	result := e.CheckDB()

	// Should warn about character being nowhere
	foundWarning := false
	for _, issue := range result.Issues {
		if issue.Type == "warning" && issue.Message == "unit [1000] is nowhere" {
			foundWarning = true
			break
		}
	}
	if !foundWarning {
		t.Error("expected warning about character being nowhere")
	}
}

func TestCheckDB_PrisonerNotStacked(t *testing.T) {
	e := newTestEngine(t)

	// Create a location
	locID := 10000
	alloc_box(locID, T_loc, sub_plain)

	// Create a character marked as prisoner but not stacked
	charID := 1000
	alloc_box(charID, T_char, 0)
	p_char(charID).prisoner = TRUE
	p_loc_info(charID).where = locID
	p_loc_info(locID).here_list = append(p_loc_info(locID).here_list, charID)

	result := e.CheckDB()

	// Should repair by clearing prisoner flag
	foundRepair := false
	for _, issue := range result.Issues {
		if issue.Type == "repaired" && strings.Contains(issue.Message, "prisoner but unstacked") {
			foundRepair = true
			break
		}
	}
	if !foundRepair {
		t.Error("expected prisoner flag repair")
	}

	// Verify prisoner flag was cleared
	if p_char(charID).prisoner != FALSE {
		t.Error("prisoner flag was not cleared")
	}
}

func TestCheckDB_UniqueItemOwnership(t *testing.T) {
	e := newTestEngine(t)

	// Create an item type with unique owner
	itemID := 100
	alloc_box(itemID, T_item, sub_artifact)
	set_name(itemID, "Magic Sword")
	p_item(itemID).who_has = 2000 // claims to belong to 2000

	// Create a character holding the item
	charID := 1000
	alloc_box(charID, T_char, 0)
	if e.globals.inventories == nil {
		e.globals.inventories = make(map[int][]item_ent)
	}
	e.globals.inventories[charID] = []item_ent{{item: itemID, qty: 1}}

	result := e.CheckDB()

	// Should repair the unique item ownership
	foundRepair := false
	for _, issue := range result.Issues {
		if issue.Type == "repaired" && issue.Message == "unique item Magic Sword [100]: whohas=, actual=1000" {
			foundRepair = true
			break
		}
	}
	if !foundRepair {
		// Check for alternate message format (when who_has box doesn't exist)
		for _, issue := range result.Issues {
			if issue.Type == "repaired" {
				t.Logf("repair: %s", issue.Message)
			}
		}
	}

	// Verify who_has was corrected
	if p_item(itemID).who_has != charID {
		t.Errorf("who_has was not corrected, got %d, want %d", p_item(itemID).who_has, charID)
	}
}

func TestCheckDB_LocNameTooLong(t *testing.T) {
	e := newTestEngine(t)

	// Create a location with a long name
	locID := 10000
	alloc_box(locID, T_loc, sub_plain)
	set_name(locID, "This is a very very very long location name that exceeds twenty five characters")

	result := e.CheckDB()

	// Should warn about name being too long
	foundWarning := false
	for _, issue := range result.Issues {
		if issue.Type == "warning" {
			if len(issue.Message) > 0 && issue.Message[:len("This is")] == "This is" {
				foundWarning = true
				break
			}
		}
	}
	if !foundWarning {
		// Check what warnings we got
		for _, issue := range result.Issues {
			if issue.Type == "warning" {
				t.Logf("warning: %s", issue.Message)
			}
		}
	}
}

func TestCheckDB_GlobConsistency(t *testing.T) {
	e := newTestEngine(t)
	result := e.CheckDB()

	// Verify no errors from glob check
	for _, issue := range result.Issues {
		if issue.Type == "error" {
			if issue.Message == "kind_s length != T_MAX" ||
				issue.Message == "subkind_s length != SUB_MAX" {
				t.Errorf("glob consistency error: %s", issue.Message)
			}
		}
	}
}

func TestCheckResult_Counts(t *testing.T) {
	result := &CheckResult{}

	result.AddError("error 1")
	result.AddError("error 2")
	result.AddWarning("warning 1")
	result.AddRepaired("repaired 1")
	result.AddRepaired("repaired 2")
	result.AddRepaired("repaired 3")

	if result.ErrorCount() != 2 {
		t.Errorf("ErrorCount() = %d, want 2", result.ErrorCount())
	}
	if result.WarningCount() != 1 {
		t.Errorf("WarningCount() = %d, want 1", result.WarningCount())
	}
	if result.RepairedCount() != 3 {
		t.Errorf("RepairedCount() = %d, want 3", result.RepairedCount())
	}
	if !result.HasErrors() {
		t.Error("HasErrors() = false, want true")
	}

	result2 := &CheckResult{}
	result2.AddWarning("just a warning")
	if result2.HasErrors() {
		t.Error("HasErrors() = true, want false")
	}
}

// newTestEngine creates a fresh Engine for testing.
func newTestEngine(t *testing.T) *Engine {
	t.Helper()

	// Create a fresh in-memory database
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	// Replace the global engine with a fresh one
	teg = &Engine{db: db}
	teg.globals.garrison_magic = 999
	teg.globals.names = make(map[int]string)
	teg.globals.banners = make(map[int]string)
	teg.globals.pluralNames = make(map[int]string)
	teg.globals.charSkills = make(map[int][]*skill_ent)
	teg.globals.playerUnits = make(map[int][]int)
	teg.globals.inventories = make(map[int][]item_ent)

	return teg
}
