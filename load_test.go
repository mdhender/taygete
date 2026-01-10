// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

import (
	"database/sql"
	"testing"
)

// insertTestWorld inserts a minimal test world into the database.
// Creates: 1 region, 2 provinces, 1 player, 1 character, 1 gate
func insertTestWorld(t *testing.T, db *sql.DB) {
	t.Helper()

	// Insert a region (id=58760, subkind=sub_region)
	_, err := db.Exec(`
		INSERT INTO entities (id, kind, subkind, name)
		VALUES (58760, ?, ?, 'Provinia')
	`, T_loc, sub_region)
	if err != nil {
		t.Fatalf("insert region entity: %v", err)
	}
	_, err = db.Exec(`
		INSERT INTO locations (id, terrain_subkind)
		VALUES (58760, ?)
	`, sub_region)
	if err != nil {
		t.Fatalf("insert region location: %v", err)
	}

	// Insert a province (id=10000, subkind=sub_plain)
	_, err = db.Exec(`
		INSERT INTO entities (id, kind, subkind, name, parent_loc_id)
		VALUES (10000, ?, ?, 'Greyfell', 58760)
	`, T_loc, sub_plain)
	if err != nil {
		t.Fatalf("insert province entity: %v", err)
	}
	_, err = db.Exec(`
		INSERT INTO locations (id, region_id, terrain_subkind, civ)
		VALUES (10000, 58760, ?, 5)
	`, sub_plain)
	if err != nil {
		t.Fatalf("insert province location: %v", err)
	}

	// Insert another province (id=10001, subkind=sub_forest)
	_, err = db.Exec(`
		INSERT INTO entities (id, kind, subkind, name, parent_loc_id)
		VALUES (10001, ?, ?, 'Darkwood', 58760)
	`, T_loc, sub_forest)
	if err != nil {
		t.Fatalf("insert province2 entity: %v", err)
	}
	_, err = db.Exec(`
		INSERT INTO locations (id, region_id, terrain_subkind, civ)
		VALUES (10001, 58760, ?, 3)
	`, sub_forest)
	if err != nil {
		t.Fatalf("insert province2 location: %v", err)
	}

	// Insert a player (id=50001)
	_, err = db.Exec(`
		INSERT INTO entities (id, kind, subkind, name)
		VALUES (50001, ?, ?, 'Test Faction')
	`, T_player, sub_pl_regular)
	if err != nil {
		t.Fatalf("insert player entity: %v", err)
	}
	_, err = db.Exec(`
		INSERT INTO players (id, code, name, subkind)
		VALUES (50001, 'aa1', 'Test Faction', ?)
	`, sub_pl_regular)
	if err != nil {
		t.Fatalf("insert player: %v", err)
	}

	// Insert a character (id=1001)
	_, err = db.Exec(`
		INSERT INTO entities (id, kind, subkind, name, parent_loc_id)
		VALUES (1001, ?, 0, 'Osswid', 10000)
	`, T_char)
	if err != nil {
		t.Fatalf("insert char entity: %v", err)
	}
	_, err = db.Exec(`
		INSERT INTO characters (id, player_id, loc_id, health, loy_kind, loy_rate)
		VALUES (1001, 50001, 10000, 100, ?, 100)
	`, LOY_oath)
	if err != nil {
		t.Fatalf("insert character: %v", err)
	}

	// Insert character magic data
	_, err = db.Exec(`
		INSERT INTO char_magic (char_id, cur_aura, max_aura, hide_mage)
		VALUES (1001, 5, 10, 1)
	`)
	if err != nil {
		t.Fatalf("insert char_magic: %v", err)
	}

	// Insert a gate (id=59001)
	_, err = db.Exec(`
		INSERT INTO entities (id, kind, subkind, name, parent_loc_id)
		VALUES (59001, ?, 0, 'Ancient Gate', 10000)
	`, T_gate)
	if err != nil {
		t.Fatalf("insert gate entity: %v", err)
	}
	_, err = db.Exec(`
		INSERT INTO gates (id, from_loc_id, to_loc_id, road_hidden)
		VALUES (59001, 10000, 10001, 0)
	`)
	if err != nil {
		t.Fatalf("insert gate: %v", err)
	}
}

func TestLoadWorldEntities(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	// Insert test world
	insertTestWorld(t, db)

	// Create engine and load world
	e := &Engine{db: db}
	e.globals.names = make(map[int]string)
	e.globals.banners = make(map[int]string)
	e.globals.pluralNames = make(map[int]string)

	err = e.LoadWorld()
	if err != nil {
		t.Fatalf("LoadWorld: %v", err)
	}

	// Verify region loaded
	if e.globals.bx[58760] == nil {
		t.Error("region 58760 not loaded")
	} else {
		if e.globals.bx[58760].kind != T_loc {
			t.Errorf("region kind = %d, want %d", e.globals.bx[58760].kind, T_loc)
		}
		if e.globals.bx[58760].skind != sub_region {
			t.Errorf("region subkind = %d, want %d", e.globals.bx[58760].skind, sub_region)
		}
		if e.globals.names[58760] != "Provinia" {
			t.Errorf("region name = %q, want 'Provinia'", e.globals.names[58760])
		}
	}

	// Verify province loaded
	if e.globals.bx[10000] == nil {
		t.Error("province 10000 not loaded")
	} else {
		if e.globals.bx[10000].kind != T_loc {
			t.Errorf("province kind = %d, want %d", e.globals.bx[10000].kind, T_loc)
		}
		if e.globals.bx[10000].skind != sub_plain {
			t.Errorf("province subkind = %d, want %d", e.globals.bx[10000].skind, sub_plain)
		}
		if e.globals.names[10000] != "Greyfell" {
			t.Errorf("province name = %q, want 'Greyfell'", e.globals.names[10000])
		}
		if e.globals.bx[10000].x_loc_info.where != 58760 {
			t.Errorf("province parent = %d, want 58760", e.globals.bx[10000].x_loc_info.where)
		}
	}

	// Verify location details loaded
	if e.globals.bx[10000].x_loc == nil {
		t.Error("province x_loc not allocated")
	} else {
		if e.globals.bx[10000].x_loc.civ != 5 {
			t.Errorf("province civ = %d, want 5", e.globals.bx[10000].x_loc.civ)
		}
	}
}

func TestLoadWorldCharacters(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	insertTestWorld(t, db)

	e := &Engine{db: db}
	e.globals.names = make(map[int]string)
	e.globals.banners = make(map[int]string)
	e.globals.pluralNames = make(map[int]string)

	err = e.LoadWorld()
	if err != nil {
		t.Fatalf("LoadWorld: %v", err)
	}

	// Verify character loaded
	if e.globals.bx[1001] == nil {
		t.Fatal("character 1001 not loaded")
	}

	if e.globals.bx[1001].kind != T_char {
		t.Errorf("char kind = %d, want %d", e.globals.bx[1001].kind, T_char)
	}
	if e.globals.names[1001] != "Osswid" {
		t.Errorf("char name = %q, want 'Osswid'", e.globals.names[1001])
	}

	// Verify character details
	ch := e.globals.bx[1001].x_char
	if ch == nil {
		t.Fatal("character x_char not allocated")
	}
	if ch.health != 100 {
		t.Errorf("char health = %d, want 100", ch.health)
	}
	if ch.loy_kind != LOY_oath {
		t.Errorf("char loy_kind = %d, want %d", ch.loy_kind, LOY_oath)
	}
	if ch.loy_rate != 100 {
		t.Errorf("char loy_rate = %d, want 100", ch.loy_rate)
	}
	if ch.unit_lord != 50001 {
		t.Errorf("char unit_lord = %d, want 50001", ch.unit_lord)
	}

	// Verify character location
	if e.globals.bx[1001].x_loc_info.where != 10000 {
		t.Errorf("char location = %d, want 10000", e.globals.bx[1001].x_loc_info.where)
	}

	// Verify character magic
	m := ch.x_char_magic
	if m == nil {
		t.Fatal("character x_char_magic not allocated")
	}
	if m.cur_aura != 5 {
		t.Errorf("char cur_aura = %d, want 5", m.cur_aura)
	}
	if m.max_aura != 10 {
		t.Errorf("char max_aura = %d, want 10", m.max_aura)
	}
	if m.hide_mage != 1 {
		t.Errorf("char hide_mage = %d, want 1", m.hide_mage)
	}
}

func TestLoadWorldPlayers(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	insertTestWorld(t, db)

	e := &Engine{db: db}
	e.globals.names = make(map[int]string)
	e.globals.banners = make(map[int]string)
	e.globals.pluralNames = make(map[int]string)

	err = e.LoadWorld()
	if err != nil {
		t.Fatalf("LoadWorld: %v", err)
	}

	// Verify player loaded
	if e.globals.bx[50001] == nil {
		t.Fatal("player 50001 not loaded")
	}

	if e.globals.bx[50001].kind != T_player {
		t.Errorf("player kind = %d, want %d", e.globals.bx[50001].kind, T_player)
	}
	if e.globals.bx[50001].skind != sub_pl_regular {
		t.Errorf("player subkind = %d, want %d", e.globals.bx[50001].skind, sub_pl_regular)
	}
	if e.globals.names[50001] != "Test Faction" {
		t.Errorf("player name = %q, want 'Test Faction'", e.globals.names[50001])
	}
}

func TestLoadWorldGates(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	insertTestWorld(t, db)

	e := &Engine{db: db}
	e.globals.names = make(map[int]string)
	e.globals.banners = make(map[int]string)
	e.globals.pluralNames = make(map[int]string)

	err = e.LoadWorld()
	if err != nil {
		t.Fatalf("LoadWorld: %v", err)
	}

	// Verify gate loaded
	if e.globals.bx[59001] == nil {
		t.Fatal("gate 59001 not loaded")
	}

	if e.globals.bx[59001].kind != T_gate {
		t.Errorf("gate kind = %d, want %d", e.globals.bx[59001].kind, T_gate)
	}
	if e.globals.names[59001] != "Ancient Gate" {
		t.Errorf("gate name = %q, want 'Ancient Gate'", e.globals.names[59001])
	}

	// Verify gate details
	g := e.globals.bx[59001].x_gate
	if g == nil {
		t.Fatal("gate x_gate not allocated")
	}
	if g.to_loc != 10001 {
		t.Errorf("gate to_loc = %d, want 10001", g.to_loc)
	}
	if g.road_hidden != 0 {
		t.Errorf("gate road_hidden = %d, want 0", g.road_hidden)
	}

	// Verify gate location (from_loc)
	if e.globals.bx[59001].x_loc_info.where != 10000 {
		t.Errorf("gate location = %d, want 10000", e.globals.bx[59001].x_loc_info.where)
	}
}

func TestLoadWorldKindChains(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	insertTestWorld(t, db)

	e := &Engine{db: db}
	e.globals.names = make(map[int]string)
	e.globals.banners = make(map[int]string)
	e.globals.pluralNames = make(map[int]string)

	err = e.LoadWorld()
	if err != nil {
		t.Fatalf("LoadWorld: %v", err)
	}

	// Verify location kind chain
	locCount := 0
	for id := e.KindFirst(T_loc); id > 0; id = e.KindNext(id) {
		locCount++
		if locCount > 100 {
			t.Fatal("infinite loop in location kind chain")
		}
	}
	if locCount != 3 { // region + 2 provinces
		t.Errorf("location count = %d, want 3", locCount)
	}

	// Verify character kind chain
	charCount := 0
	for id := e.KindFirst(T_char); id > 0; id = e.KindNext(id) {
		charCount++
	}
	if charCount != 1 {
		t.Errorf("character count = %d, want 1", charCount)
	}

	// Verify player kind chain
	playerCount := 0
	for id := e.KindFirst(T_player); id > 0; id = e.KindNext(id) {
		playerCount++
	}
	if playerCount != 1 {
		t.Errorf("player count = %d, want 1", playerCount)
	}
}

func TestLoadWorldSubkindChains(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	insertTestWorld(t, db)

	e := &Engine{db: db}
	e.globals.names = make(map[int]string)
	e.globals.banners = make(map[int]string)
	e.globals.pluralNames = make(map[int]string)

	err = e.LoadWorld()
	if err != nil {
		t.Fatalf("LoadWorld: %v", err)
	}

	// Verify plain subkind chain
	plainCount := 0
	for id := e.SubFirst(sub_plain); id > 0; id = e.SubNext(id) {
		plainCount++
		if plainCount > 100 {
			t.Fatal("infinite loop in plain subkind chain")
		}
	}
	if plainCount != 1 {
		t.Errorf("plain count = %d, want 1", plainCount)
	}

	// Verify forest subkind chain
	forestCount := 0
	for id := e.SubFirst(sub_forest); id > 0; id = e.SubNext(id) {
		forestCount++
	}
	if forestCount != 1 {
		t.Errorf("forest count = %d, want 1", forestCount)
	}

	// Verify region subkind chain
	regionCount := 0
	for id := e.SubFirst(sub_region); id > 0; id = e.SubNext(id) {
		regionCount++
	}
	if regionCount != 1 {
		t.Errorf("region count = %d, want 1", regionCount)
	}
}

func TestLoadWorldClearsExisting(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	e := &Engine{db: db}
	e.globals.names = make(map[int]string)
	e.globals.banners = make(map[int]string)
	e.globals.pluralNames = make(map[int]string)

	// Manually create a box
	e.globals.bx[999] = &box{kind: T_item}
	e.globals.names[999] = "Old Item"

	// Insert minimal world
	insertTestWorld(t, db)

	// Load world - should clear old data
	err = e.LoadWorld()
	if err != nil {
		t.Fatalf("LoadWorld: %v", err)
	}

	// Verify old box was cleared
	if e.globals.bx[999] != nil {
		t.Error("old box 999 not cleared")
	}
	if e.globals.names[999] != "" {
		t.Error("old name 999 not cleared")
	}
}

func TestLoadWorldItemTypes(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	// Insert test item type
	_, err = db.Exec(`
		INSERT INTO item_types (id, subkind, name, weight, is_animal, prominent)
		VALUES (1, 0, 'gold', 0, 0, 0)
	`)
	if err != nil {
		t.Fatalf("insert item_type: %v", err)
	}
	_, err = db.Exec(`
		INSERT INTO item_types (id, subkind, name, weight, is_animal, prominent)
		VALUES (10, 0, 'peasant', 100, 1, 0)
	`)
	if err != nil {
		t.Fatalf("insert item_type 2: %v", err)
	}

	e := &Engine{db: db}
	e.globals.names = make(map[int]string)
	e.globals.banners = make(map[int]string)
	e.globals.pluralNames = make(map[int]string)
	e.globals.charSkills = make(map[int][]*skill_ent)

	err = e.LoadWorld()
	if err != nil {
		t.Fatalf("LoadWorld: %v", err)
	}

	// Verify gold loaded
	if e.globals.bx[1] == nil {
		t.Fatal("item 1 (gold) not loaded")
	}
	if e.globals.bx[1].kind != T_item {
		t.Errorf("item 1 kind = %d, want %d", e.globals.bx[1].kind, T_item)
	}
	if e.globals.names[1] != "gold" {
		t.Errorf("item 1 name = %q, want 'gold'", e.globals.names[1])
	}

	// Verify peasant loaded with weight
	if e.globals.bx[10] == nil {
		t.Fatal("item 10 (peasant) not loaded")
	}
	if e.globals.bx[10].x_item == nil {
		t.Fatal("item 10 x_item not allocated")
	}
	if e.globals.bx[10].x_item.weight != 100 {
		t.Errorf("item 10 weight = %d, want 100", e.globals.bx[10].x_item.weight)
	}
}

func TestLoadWorldSkills(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	// Insert test skill
	_, err = db.Exec(`
		INSERT INTO skills (id, name, category, is_magic)
		VALUES (600, 'Shipcraft', 'craft', 0)
	`)
	if err != nil {
		t.Fatalf("insert skill: %v", err)
	}
	_, err = db.Exec(`
		INSERT INTO skills (id, name, category, is_magic)
		VALUES (800, 'Basic Magic', 'magic', 1)
	`)
	if err != nil {
		t.Fatalf("insert skill 2: %v", err)
	}

	e := &Engine{db: db}
	e.globals.names = make(map[int]string)
	e.globals.banners = make(map[int]string)
	e.globals.pluralNames = make(map[int]string)
	e.globals.charSkills = make(map[int][]*skill_ent)

	err = e.LoadWorld()
	if err != nil {
		t.Fatalf("LoadWorld: %v", err)
	}

	// Verify Shipcraft loaded
	if e.globals.bx[600] == nil {
		t.Fatal("skill 600 (Shipcraft) not loaded")
	}
	if e.globals.bx[600].kind != T_skill {
		t.Errorf("skill 600 kind = %d, want %d", e.globals.bx[600].kind, T_skill)
	}
	if e.globals.names[600] != "Shipcraft" {
		t.Errorf("skill 600 name = %q, want 'Shipcraft'", e.globals.names[600])
	}

	// Verify Basic Magic loaded with magic subkind
	if e.globals.bx[800] == nil {
		t.Fatal("skill 800 (Basic Magic) not loaded")
	}
	if e.globals.bx[800].skind != sub_magic {
		t.Errorf("skill 800 skind = %d, want %d", e.globals.bx[800].skind, sub_magic)
	}
}

func TestLoadWorldCharSkills(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	// Insert test world with a character
	insertTestWorld(t, db)

	// Insert skill and char_skill
	_, err = db.Exec(`
		INSERT INTO skills (id, name, category, is_magic)
		VALUES (600, 'Shipcraft', 'craft', 0)
	`)
	if err != nil {
		t.Fatalf("insert skill: %v", err)
	}
	_, err = db.Exec(`
		INSERT INTO char_skills (char_id, skill_id, level, experience)
		VALUES (1001, 600, 14, 50)
	`)
	if err != nil {
		t.Fatalf("insert char_skill: %v", err)
	}

	e := &Engine{db: db}
	e.globals.names = make(map[int]string)
	e.globals.banners = make(map[int]string)
	e.globals.pluralNames = make(map[int]string)
	e.globals.charSkills = make(map[int][]*skill_ent)

	err = e.LoadWorld()
	if err != nil {
		t.Fatalf("LoadWorld: %v", err)
	}

	// Verify character skill loaded
	skills := e.getCharSkills(1001)
	if len(skills) != 1 {
		t.Fatalf("char 1001 skills count = %d, want 1", len(skills))
	}
	if skills[0].skill != 600 {
		t.Errorf("char 1001 skill[0].skill = %d, want 600", skills[0].skill)
	}
	if skills[0].days_studied != 14 {
		t.Errorf("char 1001 skill[0].days_studied = %d, want 14", skills[0].days_studied)
	}
	if skills[0].experience != 50 {
		t.Errorf("char 1001 skill[0].experience = %d, want 50", skills[0].experience)
	}
}
