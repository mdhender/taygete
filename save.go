// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

import (
	"database/sql"
	"fmt"
)

// SaveWorld saves the in-memory world state to the database.
// It clears existing data and writes all entities from the bx array.
func (e *Engine) SaveWorld() error {
	tx, err := e.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Clear existing data (in reverse order of foreign key dependencies)
	if err := e.clearDBTables(tx); err != nil {
		return fmt.Errorf("clear tables: %w", err)
	}

	// Save entities
	if err := e.saveEntities(tx); err != nil {
		return fmt.Errorf("save entities: %w", err)
	}

	// Save locations
	if err := e.saveLocations(tx); err != nil {
		return fmt.Errorf("save locations: %w", err)
	}

	// Save players (before characters due to FK)
	if err := e.savePlayers(tx); err != nil {
		return fmt.Errorf("save players: %w", err)
	}

	// Save characters
	if err := e.saveCharacters(tx); err != nil {
		return fmt.Errorf("save characters: %w", err)
	}

	// Save character magic
	if err := e.saveCharMagic(tx); err != nil {
		return fmt.Errorf("save char_magic: %w", err)
	}

	// Save character skills
	if err := e.saveCharSkills(tx); err != nil {
		return fmt.Errorf("save char_skills: %w", err)
	}

	// Save item types
	if err := e.saveItemTypes(tx); err != nil {
		return fmt.Errorf("save item_types: %w", err)
	}

	// Save skills
	if err := e.saveSkills(tx); err != nil {
		return fmt.Errorf("save skills: %w", err)
	}

	// Save gates
	if err := e.saveGates(tx); err != nil {
		return fmt.Errorf("save gates: %w", err)
	}

	// Save storms
	if err := e.saveStorms(tx); err != nil {
		return fmt.Errorf("save storms: %w", err)
	}

	// Save ships
	if err := e.saveShips(tx); err != nil {
		return fmt.Errorf("save ships: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// clearDBTables clears all entity-related tables in reverse FK order.
func (e *Engine) clearDBTables(tx *sql.Tx) error {
	tables := []string{
		"char_skills",
		"char_magic",
		"ships",
		"storms",
		"gates",
		"characters",
		"players",
		"locations",
		"entities",
		"item_types",
		"skills",
	}

	for _, table := range tables {
		if _, err := tx.Exec("DELETE FROM " + table); err != nil {
			return fmt.Errorf("delete from %s: %w", table, err)
		}
	}

	return nil
}

// saveEntities saves all boxes to the entities table.
func (e *Engine) saveEntities(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`
		INSERT INTO entities (id, kind, subkind, name, parent_loc_id)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id := 1; id < MAX_BOXES; id++ {
		b := e.globals.bx[id]
		if b == nil {
			continue
		}

		var name sql.NullString
		if n := e.globals.names[id]; n != "" {
			name = sql.NullString{String: n, Valid: true}
		}

		var parentLocID sql.NullInt64
		if b.x_loc_info.where > 0 {
			parentLocID = sql.NullInt64{Int64: int64(b.x_loc_info.where), Valid: true}
		}

		if _, err := stmt.Exec(id, int(b.kind), int(b.skind), name, parentLocID); err != nil {
			return fmt.Errorf("insert entity %d: %w", id, err)
		}
	}

	return nil
}

// saveLocations saves location data to the locations table.
func (e *Engine) saveLocations(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`
		INSERT INTO locations (id, region_id, province_id, parent_loc_id, terrain_subkind,
		                       barrier, shroud, civ, sea_lane, is_safe_haven)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id := 1; id < MAX_BOXES; id++ {
		b := e.globals.bx[id]
		if b == nil || b.kind != T_loc {
			continue
		}

		var regionID, provinceID, parentLocID sql.NullInt64
		if b.x_loc_info.where > 0 {
			parentLocID = sql.NullInt64{Int64: int64(b.x_loc_info.where), Valid: true}
		}

		barrier, shroud, civ, seaLane := 0, 0, 0, 0
		safeHaven := 0
		if b.x_loc != nil {
			barrier = int(b.x_loc.barrier)
			shroud = int(b.x_loc.shroud)
			civ = int(b.x_loc.civ)
			seaLane = int(b.x_loc.sea_lane)
		}
		if b.x_subloc != nil && b.x_subloc.safe != 0 {
			safeHaven = 1
		}

		if _, err := stmt.Exec(id, regionID, provinceID, parentLocID, int(b.skind),
			barrier, shroud, civ, seaLane, safeHaven); err != nil {
			return fmt.Errorf("insert location %d: %w", id, err)
		}
	}

	return nil
}

// saveCharacters saves character data to the characters table.
func (e *Engine) saveCharacters(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`
		INSERT INTO characters (id, player_id, loc_id, health, sick, loy_kind, loy_rate,
		                        unit_item, guard, npc_prog, moving_since, gone_flag,
		                        is_npc, is_dead)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id := 1; id < MAX_BOXES; id++ {
		b := e.globals.bx[id]
		if b == nil || b.kind != T_char {
			continue
		}

		var playerID, locID sql.NullInt64
		var health, sick int
		var loyKind, loyRate, unitItem, guard sql.NullInt64
		var npcProg, movingSince, goneFlag sql.NullInt64
		var isNPC, isDead int

		if b.x_loc_info.where > 0 {
			locID = sql.NullInt64{Int64: int64(b.x_loc_info.where), Valid: true}
		}

		if b.x_char != nil {
			ch := b.x_char
			health = int(ch.health)
			sick = int(ch.sick)

			if ch.loy_kind != 0 {
				loyKind = sql.NullInt64{Int64: int64(ch.loy_kind), Valid: true}
			}
			if ch.loy_rate != 0 {
				loyRate = sql.NullInt64{Int64: int64(ch.loy_rate), Valid: true}
			}
			if ch.unit_item != 0 {
				unitItem = sql.NullInt64{Int64: int64(ch.unit_item), Valid: true}
			}
			if ch.guard != 0 {
				guard = sql.NullInt64{Int64: int64(ch.guard), Valid: true}
			}
			if ch.npc_prog != 0 {
				npcProg = sql.NullInt64{Int64: int64(ch.npc_prog), Valid: true}
			}
			if ch.moving != 0 {
				movingSince = sql.NullInt64{Int64: int64(ch.moving), Valid: true}
			}
			if ch.unit_lord > 0 {
				playerID = sql.NullInt64{Int64: int64(ch.unit_lord), Valid: true}
			}
		}

		if _, err := stmt.Exec(id, playerID, locID, health, sick,
			loyKind, loyRate, unitItem, guard, npcProg,
			movingSince, goneFlag, isNPC, isDead); err != nil {
			return fmt.Errorf("insert character %d: %w", id, err)
		}
	}

	return nil
}

// saveCharMagic saves character magic data to the char_magic table.
func (e *Engine) saveCharMagic(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`
		INSERT INTO char_magic (char_id, pray, hide_self, vis_protect, hide_mage,
		                        cur_aura, max_aura, aura_reflect, pledge, auraculum,
		                        ability_shroud, fee, ferry_flag)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id := 1; id < MAX_BOXES; id++ {
		b := e.globals.bx[id]
		if b == nil || b.kind != T_char || b.x_char == nil || b.x_char.x_char_magic == nil {
			continue
		}

		m := b.x_char.x_char_magic

		var pledge, auraculum sql.NullInt64
		if m.pledge != 0 {
			pledge = sql.NullInt64{Int64: int64(m.pledge), Valid: true}
		}
		if m.auraculum != 0 {
			auraculum = sql.NullInt64{Int64: int64(m.auraculum), Valid: true}
		}

		if _, err := stmt.Exec(id, int(m.pray), int(m.hide_self), int(m.vis_protect),
			int(m.hide_mage), m.cur_aura, m.max_aura, int(m.aura_reflect),
			pledge, auraculum, int(m.ability_shroud), m.fee, int(m.ferry_flag)); err != nil {
			return fmt.Errorf("insert char_magic %d: %w", id, err)
		}
	}

	return nil
}

// savePlayers saves player data to the players table.
func (e *Engine) savePlayers(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`
		INSERT INTO players (id, code, name, subkind)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id := 1; id < MAX_BOXES; id++ {
		b := e.globals.bx[id]
		if b == nil || b.kind != T_player {
			continue
		}

		code := int_to_code(id)
		var name sql.NullString
		if n := e.globals.names[id]; n != "" {
			name = sql.NullString{String: n, Valid: true}
		}

		if _, err := stmt.Exec(id, code, name, int(b.skind)); err != nil {
			return fmt.Errorf("insert player %d: %w", id, err)
		}
	}

	return nil
}

// saveGates saves gate data to the gates table.
func (e *Engine) saveGates(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`
		INSERT INTO gates (id, from_loc_id, to_loc_id, road_hidden)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id := 1; id < MAX_BOXES; id++ {
		b := e.globals.bx[id]
		if b == nil || b.kind != T_gate {
			continue
		}

		fromLocID := b.x_loc_info.where
		toLocID := 0
		roadHidden := 0

		if b.x_gate != nil {
			toLocID = b.x_gate.to_loc
			roadHidden = int(b.x_gate.road_hidden)
		}

		if _, err := stmt.Exec(id, fromLocID, toLocID, roadHidden); err != nil {
			return fmt.Errorf("insert gate %d: %w", id, err)
		}
	}

	return nil
}

// saveStorms saves storm data to the storms table.
func (e *Engine) saveStorms(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`
		INSERT INTO storms (id, strength, moving_to, moving_since)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id := 1; id < MAX_BOXES; id++ {
		b := e.globals.bx[id]
		if b == nil || b.kind != T_storm {
			continue
		}

		strength := 0
		var movingTo, movingSince sql.NullInt64

		if b.x_misc != nil {
			strength = int(b.x_misc.storm_str)
			if b.x_misc.storm_move != 0 {
				movingTo = sql.NullInt64{Int64: int64(b.x_misc.storm_move), Valid: true}
			}
		}

		if _, err := stmt.Exec(id, strength, movingTo, movingSince); err != nil {
			return fmt.Errorf("insert storm %d: %w", id, err)
		}
	}

	return nil
}

// saveShips saves ship data to the ships table.
func (e *Engine) saveShips(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`
		INSERT INTO ships (id, loc_id, capacity, storm_bind, moving_since)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id := 1; id < MAX_BOXES; id++ {
		b := e.globals.bx[id]
		if b == nil || b.kind != T_ship {
			continue
		}

		var locID, capacity, stormBind, movingSince sql.NullInt64

		if b.x_loc_info.where > 0 {
			locID = sql.NullInt64{Int64: int64(b.x_loc_info.where), Valid: true}
		}

		if b.x_subloc != nil {
			if b.x_subloc.capacity != 0 {
				capacity = sql.NullInt64{Int64: int64(b.x_subloc.capacity), Valid: true}
			}
			if b.x_subloc.moving != 0 {
				movingSince = sql.NullInt64{Int64: int64(b.x_subloc.moving), Valid: true}
			}
		}

		if b.x_misc != nil && b.x_misc.bind_storm != 0 {
			stormBind = sql.NullInt64{Int64: int64(b.x_misc.bind_storm), Valid: true}
		}

		if _, err := stmt.Exec(id, locID, capacity, stormBind, movingSince); err != nil {
			return fmt.Errorf("insert ship %d: %w", id, err)
		}
	}

	return nil
}

// saveItemTypes saves item type data to the item_types table.
func (e *Engine) saveItemTypes(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`
		INSERT INTO item_types (id, subkind, name, weight, is_animal, prominent)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id := 1; id < MAX_BOXES; id++ {
		b := e.globals.bx[id]
		if b == nil || b.kind != T_item {
			continue
		}

		name := e.globals.names[id]
		weight, isAnimal, prominent := 0, 0, 0

		if b.x_item != nil {
			weight = int(b.x_item.weight)
			isAnimal = int(b.x_item.is_man_item)
			prominent = int(b.x_item.prominent)
		}

		if _, err := stmt.Exec(id, int(b.skind), name, weight, isAnimal, prominent); err != nil {
			return fmt.Errorf("insert item_type %d: %w", id, err)
		}
	}

	return nil
}

// saveSkills saves skill data to the skills table.
func (e *Engine) saveSkills(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`
		INSERT INTO skills (id, name, category, is_magic)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id := 1; id < MAX_BOXES; id++ {
		b := e.globals.bx[id]
		if b == nil || b.kind != T_skill {
			continue
		}

		name := e.globals.names[id]
		isMagic := 0
		if b.skind == sub_magic {
			isMagic = 1
		}

		var category sql.NullString

		if _, err := stmt.Exec(id, name, category, isMagic); err != nil {
			return fmt.Errorf("insert skill %d: %w", id, err)
		}
	}

	return nil
}

// saveCharSkills saves character skill data to the char_skills table.
func (e *Engine) saveCharSkills(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`
		INSERT INTO char_skills (char_id, skill_id, level, experience)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Save skills from the charSkills map
	for charID, skills := range e.globals.charSkills {
		for _, sk := range skills {
			if sk == nil {
				continue
			}
			if _, err := stmt.Exec(charID, sk.skill, sk.days_studied, int(sk.experience)); err != nil {
				return fmt.Errorf("insert char_skill %d/%d: %w", charID, sk.skill, err)
			}
		}
	}

	return nil
}
