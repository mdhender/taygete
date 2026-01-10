// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

import (
	"testing"
)

func TestSaveWorldEntities(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	e := &Engine{db: db}
	e.globals.names = make(map[int]string)
	e.globals.banners = make(map[int]string)
	e.globals.pluralNames = make(map[int]string)

	// Create a location in memory
	e.globals.bx[10000] = &box{
		kind:  T_loc,
		skind: sub_plain,
	}
	e.globals.bx[10000].x_loc = &entity_loc{
		civ: 5,
	}
	e.globals.names[10000] = "Test Province"
	e.addToKindChain(10000)
	e.addToSubkindChain(10000)

	// Save to database
	err = e.SaveWorld()
	if err != nil {
		t.Fatalf("SaveWorld: %v", err)
	}

	// Verify entity was saved
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM entities WHERE id = 10000").Scan(&count)
	if err != nil {
		t.Fatalf("query entities: %v", err)
	}
	if count != 1 {
		t.Errorf("entities count = %d, want 1", count)
	}

	// Verify location was saved
	err = db.QueryRow("SELECT COUNT(*) FROM locations WHERE id = 10000").Scan(&count)
	if err != nil {
		t.Fatalf("query locations: %v", err)
	}
	if count != 1 {
		t.Errorf("locations count = %d, want 1", count)
	}
}

func TestSaveWorldRoundTrip(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	// Insert test world via SQL
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

	// Capture original state
	origRegionKind := e.globals.bx[58760].kind
	origRegionSubkind := e.globals.bx[58760].skind
	origRegionName := e.globals.names[58760]

	origProvinceKind := e.globals.bx[10000].kind
	origProvinceSubkind := e.globals.bx[10000].skind
	origProvinceName := e.globals.names[10000]
	origProvinceCiv := e.globals.bx[10000].x_loc.civ

	origCharKind := e.globals.bx[1001].kind
	origCharName := e.globals.names[1001]
	origCharHealth := e.globals.bx[1001].x_char.health
	origCharLoyKind := e.globals.bx[1001].x_char.loy_kind
	origCharCurAura := e.globals.bx[1001].x_char.x_char_magic.cur_aura
	origCharMaxAura := e.globals.bx[1001].x_char.x_char_magic.max_aura

	origPlayerKind := e.globals.bx[50001].kind
	origPlayerSubkind := e.globals.bx[50001].skind
	origPlayerName := e.globals.names[50001]

	origGateKind := e.globals.bx[59001].kind
	origGateName := e.globals.names[59001]
	origGateToLoc := e.globals.bx[59001].x_gate.to_loc
	origGateLoc := e.globals.bx[59001].x_loc_info.where

	// Save to database
	err = e.SaveWorld()
	if err != nil {
		t.Fatalf("SaveWorld: %v", err)
	}

	// Clear in-memory state
	e.clearWorld()

	// Reload from database
	err = e.LoadWorld()
	if err != nil {
		t.Fatalf("LoadWorld (after save): %v", err)
	}

	// Verify region
	if e.globals.bx[58760] == nil {
		t.Fatal("region 58760 not reloaded")
	}
	if e.globals.bx[58760].kind != origRegionKind {
		t.Errorf("region kind = %d, want %d", e.globals.bx[58760].kind, origRegionKind)
	}
	if e.globals.bx[58760].skind != origRegionSubkind {
		t.Errorf("region subkind = %d, want %d", e.globals.bx[58760].skind, origRegionSubkind)
	}
	if e.globals.names[58760] != origRegionName {
		t.Errorf("region name = %q, want %q", e.globals.names[58760], origRegionName)
	}

	// Verify province
	if e.globals.bx[10000] == nil {
		t.Fatal("province 10000 not reloaded")
	}
	if e.globals.bx[10000].kind != origProvinceKind {
		t.Errorf("province kind = %d, want %d", e.globals.bx[10000].kind, origProvinceKind)
	}
	if e.globals.bx[10000].skind != origProvinceSubkind {
		t.Errorf("province subkind = %d, want %d", e.globals.bx[10000].skind, origProvinceSubkind)
	}
	if e.globals.names[10000] != origProvinceName {
		t.Errorf("province name = %q, want %q", e.globals.names[10000], origProvinceName)
	}
	if e.globals.bx[10000].x_loc == nil {
		t.Fatal("province x_loc not reloaded")
	}
	if e.globals.bx[10000].x_loc.civ != origProvinceCiv {
		t.Errorf("province civ = %d, want %d", e.globals.bx[10000].x_loc.civ, origProvinceCiv)
	}

	// Verify character
	if e.globals.bx[1001] == nil {
		t.Fatal("character 1001 not reloaded")
	}
	if e.globals.bx[1001].kind != origCharKind {
		t.Errorf("char kind = %d, want %d", e.globals.bx[1001].kind, origCharKind)
	}
	if e.globals.names[1001] != origCharName {
		t.Errorf("char name = %q, want %q", e.globals.names[1001], origCharName)
	}
	if e.globals.bx[1001].x_char == nil {
		t.Fatal("char x_char not reloaded")
	}
	if e.globals.bx[1001].x_char.health != origCharHealth {
		t.Errorf("char health = %d, want %d", e.globals.bx[1001].x_char.health, origCharHealth)
	}
	if e.globals.bx[1001].x_char.loy_kind != origCharLoyKind {
		t.Errorf("char loy_kind = %d, want %d", e.globals.bx[1001].x_char.loy_kind, origCharLoyKind)
	}
	if e.globals.bx[1001].x_char.x_char_magic == nil {
		t.Fatal("char x_char_magic not reloaded")
	}
	if e.globals.bx[1001].x_char.x_char_magic.cur_aura != origCharCurAura {
		t.Errorf("char cur_aura = %d, want %d", e.globals.bx[1001].x_char.x_char_magic.cur_aura, origCharCurAura)
	}
	if e.globals.bx[1001].x_char.x_char_magic.max_aura != origCharMaxAura {
		t.Errorf("char max_aura = %d, want %d", e.globals.bx[1001].x_char.x_char_magic.max_aura, origCharMaxAura)
	}

	// Verify player
	if e.globals.bx[50001] == nil {
		t.Fatal("player 50001 not reloaded")
	}
	if e.globals.bx[50001].kind != origPlayerKind {
		t.Errorf("player kind = %d, want %d", e.globals.bx[50001].kind, origPlayerKind)
	}
	if e.globals.bx[50001].skind != origPlayerSubkind {
		t.Errorf("player subkind = %d, want %d", e.globals.bx[50001].skind, origPlayerSubkind)
	}
	if e.globals.names[50001] != origPlayerName {
		t.Errorf("player name = %q, want %q", e.globals.names[50001], origPlayerName)
	}

	// Verify gate
	if e.globals.bx[59001] == nil {
		t.Fatal("gate 59001 not reloaded")
	}
	if e.globals.bx[59001].kind != origGateKind {
		t.Errorf("gate kind = %d, want %d", e.globals.bx[59001].kind, origGateKind)
	}
	if e.globals.names[59001] != origGateName {
		t.Errorf("gate name = %q, want %q", e.globals.names[59001], origGateName)
	}
	if e.globals.bx[59001].x_gate == nil {
		t.Fatal("gate x_gate not reloaded")
	}
	if e.globals.bx[59001].x_gate.to_loc != origGateToLoc {
		t.Errorf("gate to_loc = %d, want %d", e.globals.bx[59001].x_gate.to_loc, origGateToLoc)
	}
	if e.globals.bx[59001].x_loc_info.where != origGateLoc {
		t.Errorf("gate location = %d, want %d", e.globals.bx[59001].x_loc_info.where, origGateLoc)
	}
}

func TestSaveWorldRoundTripKindChains(t *testing.T) {
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

	// Count entities by kind before save
	origLocCount := 0
	for id := e.KindFirst(T_loc); id > 0; id = e.KindNext(id) {
		origLocCount++
	}
	origCharCount := 0
	for id := e.KindFirst(T_char); id > 0; id = e.KindNext(id) {
		origCharCount++
	}
	origPlayerCount := 0
	for id := e.KindFirst(T_player); id > 0; id = e.KindNext(id) {
		origPlayerCount++
	}
	origGateCount := 0
	for id := e.KindFirst(T_gate); id > 0; id = e.KindNext(id) {
		origGateCount++
	}

	// Save and reload
	err = e.SaveWorld()
	if err != nil {
		t.Fatalf("SaveWorld: %v", err)
	}
	e.clearWorld()
	err = e.LoadWorld()
	if err != nil {
		t.Fatalf("LoadWorld (after save): %v", err)
	}

	// Count entities by kind after reload
	locCount := 0
	for id := e.KindFirst(T_loc); id > 0; id = e.KindNext(id) {
		locCount++
		if locCount > 100 {
			t.Fatal("infinite loop in location kind chain")
		}
	}
	if locCount != origLocCount {
		t.Errorf("location count = %d, want %d", locCount, origLocCount)
	}

	charCount := 0
	for id := e.KindFirst(T_char); id > 0; id = e.KindNext(id) {
		charCount++
	}
	if charCount != origCharCount {
		t.Errorf("character count = %d, want %d", charCount, origCharCount)
	}

	playerCount := 0
	for id := e.KindFirst(T_player); id > 0; id = e.KindNext(id) {
		playerCount++
	}
	if playerCount != origPlayerCount {
		t.Errorf("player count = %d, want %d", playerCount, origPlayerCount)
	}

	gateCount := 0
	for id := e.KindFirst(T_gate); id > 0; id = e.KindNext(id) {
		gateCount++
	}
	if gateCount != origGateCount {
		t.Errorf("gate count = %d, want %d", gateCount, origGateCount)
	}
}

func TestSaveWorldClearsOldData(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	// Insert initial world
	insertTestWorld(t, db)

	e := &Engine{db: db}
	e.globals.names = make(map[int]string)
	e.globals.banners = make(map[int]string)
	e.globals.pluralNames = make(map[int]string)

	err = e.LoadWorld()
	if err != nil {
		t.Fatalf("LoadWorld: %v", err)
	}

	// Remove the gate (which references location 10001) and the location
	e.globals.bx[59001] = nil
	e.globals.bx[10001] = nil

	// Save
	err = e.SaveWorld()
	if err != nil {
		t.Fatalf("SaveWorld: %v", err)
	}

	// Verify removed location is gone from DB
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM entities WHERE id = 10001").Scan(&count)
	if err != nil {
		t.Fatalf("query entities: %v", err)
	}
	if count != 0 {
		t.Errorf("entity 10001 should be deleted, but count = %d", count)
	}

	// Verify gate is also gone
	err = db.QueryRow("SELECT COUNT(*) FROM entities WHERE id = 59001").Scan(&count)
	if err != nil {
		t.Fatalf("query gate: %v", err)
	}
	if count != 0 {
		t.Errorf("gate 59001 should be deleted, but count = %d", count)
	}
}

func TestSaveWorldCharacterMagic(t *testing.T) {
	db, err := OpenTestDB()
	if err != nil {
		t.Fatalf("OpenTestDB: %v", err)
	}
	defer db.Close()

	e := &Engine{db: db}
	e.globals.names = make(map[int]string)
	e.globals.banners = make(map[int]string)
	e.globals.pluralNames = make(map[int]string)

	// Create a character with magic in memory
	e.globals.bx[2001] = &box{
		kind:  T_char,
		skind: 0,
	}
	e.globals.bx[2001].x_char = &entity_char{
		health:   100,
		loy_kind: LOY_oath,
		loy_rate: 50,
	}
	e.globals.bx[2001].x_char.x_char_magic = &char_magic{
		cur_aura:   7,
		max_aura:   15,
		hide_mage:  1,
		hide_self:  1,
		vis_protect: 2,
	}
	e.globals.names[2001] = "Mage Test"
	e.addToKindChain(2001)
	e.addToSubkindChain(2001)

	// Save
	err = e.SaveWorld()
	if err != nil {
		t.Fatalf("SaveWorld: %v", err)
	}

	// Clear and reload
	e.clearWorld()
	err = e.LoadWorld()
	if err != nil {
		t.Fatalf("LoadWorld: %v", err)
	}

	// Verify magic data
	if e.globals.bx[2001] == nil {
		t.Fatal("character 2001 not reloaded")
	}
	if e.globals.bx[2001].x_char == nil {
		t.Fatal("x_char not reloaded")
	}
	if e.globals.bx[2001].x_char.x_char_magic == nil {
		t.Fatal("x_char_magic not reloaded")
	}
	m := e.globals.bx[2001].x_char.x_char_magic
	if m.cur_aura != 7 {
		t.Errorf("cur_aura = %d, want 7", m.cur_aura)
	}
	if m.max_aura != 15 {
		t.Errorf("max_aura = %d, want 15", m.max_aura)
	}
	if m.hide_mage != 1 {
		t.Errorf("hide_mage = %d, want 1", m.hide_mage)
	}
	if m.hide_self != 1 {
		t.Errorf("hide_self = %d, want 1", m.hide_self)
	}
	if m.vis_protect != 2 {
		t.Errorf("vis_protect = %d, want 2", m.vis_protect)
	}
}
