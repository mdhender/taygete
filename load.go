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
	"database/sql"
	"fmt"
)

// LoadWorld loads the game world from the database into memory.
// It clears any existing world state and populates the bx array
// from the database tables.
func (e *Engine) LoadWorld() error {
	// Clear existing world state
	e.clearWorld()

	// Load entities (creates boxes with kind/subkind)
	if err := e.loadEntities(); err != nil {
		return fmt.Errorf("load entities: %w", err)
	}

	// Load locations
	if err := e.loadLocations(); err != nil {
		return fmt.Errorf("load locations: %w", err)
	}

	// Load characters
	if err := e.loadCharacters(); err != nil {
		return fmt.Errorf("load characters: %w", err)
	}

	// Load character magic
	if err := e.loadCharMagic(); err != nil {
		return fmt.Errorf("load char_magic: %w", err)
	}

	// Load character skills
	if err := e.loadCharSkills(); err != nil {
		return fmt.Errorf("load char_skills: %w", err)
	}

	// Load players
	if err := e.loadPlayers(); err != nil {
		return fmt.Errorf("load players: %w", err)
	}

	// Load item types
	if err := e.loadItemTypes(); err != nil {
		return fmt.Errorf("load item_types: %w", err)
	}

	// Load skills
	if err := e.loadSkills(); err != nil {
		return fmt.Errorf("load skills: %w", err)
	}

	// Load gates
	if err := e.loadGates(); err != nil {
		return fmt.Errorf("load gates: %w", err)
	}

	// Load storms
	if err := e.loadStorms(); err != nil {
		return fmt.Errorf("load storms: %w", err)
	}

	// Load ships
	if err := e.loadShips(); err != nil {
		return fmt.Errorf("load ships: %w", err)
	}

	// Load system config
	if err := e.loadSystemConfig(); err != nil {
		return fmt.Errorf("load system_config: %w", err)
	}

	return nil
}

// clearWorld resets the in-memory world state.
func (e *Engine) clearWorld() {
	for i := range e.globals.bx {
		e.globals.bx[i] = nil
	}
	for i := range e.globals.box_head {
		e.globals.box_head[i] = 0
	}
	for i := range e.globals.sub_head {
		e.globals.sub_head[i] = 0
	}
	e.globals.names = make(map[int]string)
	e.globals.banners = make(map[int]string)
	e.globals.pluralNames = make(map[int]string)
	e.globals.charSkills = make(map[int][]*skill_ent)
}

// loadEntities loads all entities from the database.
func (e *Engine) loadEntities() error {
	rows, err := e.db.Query(`
		SELECT id, kind, subkind, name, display_name, parent_loc_id
		FROM entities
		WHERE is_deleted = 0
		ORDER BY id
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, kind, subkind int
		var name, displayName sql.NullString
		var parentLocID sql.NullInt64

		if err := rows.Scan(&id, &kind, &subkind, &name, &displayName, &parentLocID); err != nil {
			return fmt.Errorf("scan entity %d: %w", id, err)
		}

		if id <= 0 || id >= MAX_BOXES {
			continue
		}

		// Allocate the box
		e.globals.bx[id] = &box{
			kind:  schar(kind),
			skind: schar(subkind),
		}

		// Add to kind and subkind chains
		e.addToKindChain(id)
		e.addToSubkindChain(id)

		// Set name
		if name.Valid && name.String != "" {
			e.globals.names[id] = name.String
		}

		// Set parent location
		if parentLocID.Valid {
			e.globals.bx[id].x_loc_info.where = int(parentLocID.Int64)
		}
	}

	return rows.Err()
}

// addToKindChain adds entity n to the kind chain (sorted by ID).
func (e *Engine) addToKindChain(n int) {
	if e.globals.bx[n] == nil {
		return
	}
	k := int(e.globals.bx[n].kind)

	if e.globals.box_head[k] == 0 || n < e.globals.box_head[k] {
		e.globals.bx[n].x_next_kind = e.globals.box_head[k]
		e.globals.box_head[k] = n
		return
	}

	i := e.globals.box_head[k]
	for e.globals.bx[i].x_next_kind > 0 && e.globals.bx[i].x_next_kind < n {
		i = e.globals.bx[i].x_next_kind
	}
	e.globals.bx[n].x_next_kind = e.globals.bx[i].x_next_kind
	e.globals.bx[i].x_next_kind = n
}

// addToSubkindChain adds entity n to the subkind chain (sorted by ID).
func (e *Engine) addToSubkindChain(n int) {
	if e.globals.bx[n] == nil {
		return
	}
	sk := int(e.globals.bx[n].skind)

	if e.globals.sub_head[sk] == 0 || n < e.globals.sub_head[sk] {
		e.globals.bx[n].x_next_sub = e.globals.sub_head[sk]
		e.globals.sub_head[sk] = n
		return
	}

	i := e.globals.sub_head[sk]
	for e.globals.bx[i].x_next_sub > 0 && e.globals.bx[i].x_next_sub < n {
		i = e.globals.bx[i].x_next_sub
	}
	e.globals.bx[n].x_next_sub = e.globals.bx[i].x_next_sub
	e.globals.bx[i].x_next_sub = n
}

// loadLocations loads location data into entity_loc structs.
func (e *Engine) loadLocations() error {
	rows, err := e.db.Query(`
		SELECT id, region_id, province_id, parent_loc_id, terrain_subkind,
		       barrier, shroud, civ, sea_lane, is_safe_haven
		FROM locations
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var regionID, provinceID, parentLocID sql.NullInt64
		var terrainSubkind, barrier, shroud, civ, seaLane, safeHaven int

		if err := rows.Scan(&id, &regionID, &provinceID, &parentLocID,
			&terrainSubkind, &barrier, &shroud, &civ, &seaLane, &safeHaven); err != nil {
			return fmt.Errorf("scan location %d: %w", id, err)
		}

		if e.globals.bx[id] == nil {
			continue
		}

		// Ensure x_loc exists
		if e.globals.bx[id].x_loc == nil {
			e.globals.bx[id].x_loc = &entity_loc{}
		}
		loc := e.globals.bx[id].x_loc

		loc.barrier = short(barrier)
		loc.shroud = short(shroud)
		loc.civ = schar(civ)
		loc.sea_lane = schar(seaLane)

		// Set parent location in loc_info
		if parentLocID.Valid {
			e.globals.bx[id].x_loc_info.where = int(parentLocID.Int64)
		}
	}

	return rows.Err()
}

// loadCharacters loads character data into entity_char structs.
func (e *Engine) loadCharacters() error {
	rows, err := e.db.Query(`
		SELECT id, player_id, loc_id, health, sick, loy_kind, loy_rate,
		       unit_item, guard, npc_prog, moving_since, gone_flag,
		       is_npc, is_dead
		FROM characters
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var playerID, locID sql.NullInt64
		var health, sick int
		var loyKind, loyRate, unitItem, guard sql.NullInt64
		var npcProg, movingSince, goneFlag sql.NullInt64
		var isNPC, isDead int

		if err := rows.Scan(&id, &playerID, &locID, &health, &sick,
			&loyKind, &loyRate, &unitItem, &guard, &npcProg,
			&movingSince, &goneFlag, &isNPC, &isDead); err != nil {
			return fmt.Errorf("scan character %d: %w", id, err)
		}

		if e.globals.bx[id] == nil {
			continue
		}

		// Ensure x_char exists
		if e.globals.bx[id].x_char == nil {
			e.globals.bx[id].x_char = &entity_char{}
		}
		ch := e.globals.bx[id].x_char

		ch.health = schar(health)
		ch.sick = schar(sick)

		if loyKind.Valid {
			ch.loy_kind = schar(loyKind.Int64)
		}
		if loyRate.Valid {
			ch.loy_rate = int(loyRate.Int64)
		}
		if unitItem.Valid {
			ch.unit_item = schar(unitItem.Int64)
		}
		if guard.Valid {
			ch.guard = schar(guard.Int64)
		}
		if npcProg.Valid {
			ch.npc_prog = schar(npcProg.Int64)
		}
		if movingSince.Valid {
			ch.moving = int(movingSince.Int64)
		}

		// Set player as unit_lord for top-level units
		if playerID.Valid {
			ch.unit_lord = int(playerID.Int64)
		}

		// Set location
		if locID.Valid {
			e.globals.bx[id].x_loc_info.where = int(locID.Int64)
		}
	}

	return rows.Err()
}

// loadCharMagic loads character magic data.
func (e *Engine) loadCharMagic() error {
	rows, err := e.db.Query(`
		SELECT char_id, pray, hide_self, vis_protect, hide_mage,
		       cur_aura, max_aura, aura_reflect, pledge, auraculum,
		       ability_shroud, fee, ferry_flag
		FROM char_magic
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var charID int
		var pray, hideSelf, visProtect, hideMage int
		var curAura, maxAura, auraReflect int
		var pledge, auraculum sql.NullInt64
		var abilityShroud, fee, ferryFlag int

		if err := rows.Scan(&charID, &pray, &hideSelf, &visProtect, &hideMage,
			&curAura, &maxAura, &auraReflect, &pledge, &auraculum,
			&abilityShroud, &fee, &ferryFlag); err != nil {
			return fmt.Errorf("scan char_magic %d: %w", charID, err)
		}

		if e.globals.bx[charID] == nil {
			continue
		}

		// Ensure x_char and x_char_magic exist
		if e.globals.bx[charID].x_char == nil {
			e.globals.bx[charID].x_char = &entity_char{}
		}
		if e.globals.bx[charID].x_char.x_char_magic == nil {
			e.globals.bx[charID].x_char.x_char_magic = &char_magic{}
		}
		m := e.globals.bx[charID].x_char.x_char_magic

		m.pray = schar(pray)
		m.hide_self = schar(hideSelf)
		m.vis_protect = schar(visProtect)
		m.hide_mage = schar(hideMage)
		m.cur_aura = curAura
		m.max_aura = maxAura
		m.aura_reflect = schar(auraReflect)
		m.ability_shroud = short(abilityShroud)

		if pledge.Valid {
			m.pledge = int(pledge.Int64)
		}
		if auraculum.Valid {
			m.auraculum = int(auraculum.Int64)
		}
	}

	return rows.Err()
}

// loadPlayers loads player data into entity_player structs.
func (e *Engine) loadPlayers() error {
	rows, err := e.db.Query(`
		SELECT id, code, name, subkind
		FROM players
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var code string
		var name sql.NullString
		var subkind int

		if err := rows.Scan(&id, &code, &name, &subkind); err != nil {
			return fmt.Errorf("scan player %d: %w", id, err)
		}

		if id <= 0 || id >= MAX_BOXES {
			continue
		}

		// Create box if it doesn't exist (players may not be in entities table)
		if e.globals.bx[id] == nil {
			e.globals.bx[id] = &box{
				kind:  T_player,
				skind: schar(subkind),
			}
			e.addToKindChain(id)
			e.addToSubkindChain(id)
		}

		// Ensure x_player exists
		if e.globals.bx[id].x_player == nil {
			e.globals.bx[id].x_player = &entity_player{}
		}

		// Set name
		if name.Valid && name.String != "" {
			e.globals.names[id] = name.String
		}
	}

	return rows.Err()
}

// loadGates loads gate data into entity_gate structs.
func (e *Engine) loadGates() error {
	rows, err := e.db.Query(`
		SELECT id, from_loc_id, to_loc_id, road_hidden
		FROM gates
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, fromLocID, toLocID, roadHidden int

		if err := rows.Scan(&id, &fromLocID, &toLocID, &roadHidden); err != nil {
			return fmt.Errorf("scan gate %d: %w", id, err)
		}

		if e.globals.bx[id] == nil {
			continue
		}

		// Ensure x_gate exists
		if e.globals.bx[id].x_gate == nil {
			e.globals.bx[id].x_gate = &entity_gate{}
		}
		g := e.globals.bx[id].x_gate

		g.to_loc = toLocID
		g.road_hidden = schar(roadHidden)

		// Set location (from_loc_id)
		e.globals.bx[id].x_loc_info.where = fromLocID
	}

	return rows.Err()
}

// loadStorms loads storm data into entity_misc structs.
func (e *Engine) loadStorms() error {
	rows, err := e.db.Query(`
		SELECT id, strength, moving_to, moving_since
		FROM storms
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, strength int
		var movingTo, movingSince sql.NullInt64

		if err := rows.Scan(&id, &strength, &movingTo, &movingSince); err != nil {
			return fmt.Errorf("scan storm %d: %w", id, err)
		}

		if e.globals.bx[id] == nil {
			continue
		}

		// Ensure x_misc exists
		if e.globals.bx[id].x_misc == nil {
			e.globals.bx[id].x_misc = &entity_misc{}
		}
		m := e.globals.bx[id].x_misc

		m.storm_str = short(strength)

		if movingTo.Valid {
			m.storm_move = int(movingTo.Int64)
		}
	}

	return rows.Err()
}

// loadShips loads ship data into entity_subloc structs.
func (e *Engine) loadShips() error {
	rows, err := e.db.Query(`
		SELECT id, loc_id, capacity, storm_bind, moving_since
		FROM ships
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var locID, capacity sql.NullInt64
		var stormBind, movingSince sql.NullInt64

		if err := rows.Scan(&id, &locID, &capacity, &stormBind, &movingSince); err != nil {
			return fmt.Errorf("scan ship %d: %w", id, err)
		}

		if e.globals.bx[id] == nil {
			continue
		}

		// Ensure x_subloc exists
		if e.globals.bx[id].x_subloc == nil {
			e.globals.bx[id].x_subloc = &entity_subloc{}
		}
		s := e.globals.bx[id].x_subloc

		if capacity.Valid {
			s.capacity = int(capacity.Int64)
		}
		if movingSince.Valid {
			s.moving = int(movingSince.Int64)
		}

		// Set location
		if locID.Valid {
			e.globals.bx[id].x_loc_info.where = int(locID.Int64)
		}

		// Storm binding goes in x_misc
		if stormBind.Valid {
			if e.globals.bx[id].x_misc == nil {
				e.globals.bx[id].x_misc = &entity_misc{}
			}
			e.globals.bx[id].x_misc.bind_storm = int(stormBind.Int64)
		}
	}

	return rows.Err()
}

// loadItemTypes loads item type definitions into entity_item structs.
func (e *Engine) loadItemTypes() error {
	rows, err := e.db.Query(`
		SELECT id, subkind, name, weight, is_animal, prominent
		FROM item_types
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, subkind int
		var name string
		var weight, isAnimal, prominent int

		if err := rows.Scan(&id, &subkind, &name, &weight, &isAnimal, &prominent); err != nil {
			return fmt.Errorf("scan item_type %d: %w", id, err)
		}

		if id <= 0 || id >= MAX_BOXES {
			continue
		}

		// Create box if it doesn't exist
		if e.globals.bx[id] == nil {
			e.globals.bx[id] = &box{
				kind:  T_item,
				skind: schar(subkind),
			}
			e.addToKindChain(id)
			e.addToSubkindChain(id)
		}

		// Ensure x_item exists
		if e.globals.bx[id].x_item == nil {
			e.globals.bx[id].x_item = &entity_item{}
		}
		it := e.globals.bx[id].x_item

		it.weight = short(weight)
		it.is_man_item = schar(isAnimal)
		it.prominent = schar(prominent)

		// Set name
		if name != "" {
			e.globals.names[id] = name
		}
	}

	return rows.Err()
}

// loadSkills loads skill definitions into entity_skill structs.
func (e *Engine) loadSkills() error {
	rows, err := e.db.Query(`
		SELECT id, name, category, is_magic
		FROM skills
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var category sql.NullString
		var isMagic int

		if err := rows.Scan(&id, &name, &category, &isMagic); err != nil {
			return fmt.Errorf("scan skill %d: %w", id, err)
		}

		if id <= 0 || id >= MAX_BOXES {
			continue
		}

		// Determine subkind based on is_magic
		subkind := schar(0)
		if isMagic != 0 {
			subkind = sub_magic
		}

		// Create box if it doesn't exist
		if e.globals.bx[id] == nil {
			e.globals.bx[id] = &box{
				kind:  T_skill,
				skind: subkind,
			}
			e.addToKindChain(id)
			e.addToSubkindChain(id)
		}

		// Ensure x_skill exists
		if e.globals.bx[id].x_skill == nil {
			e.globals.bx[id].x_skill = &entity_skill{}
		}

		// Set name
		if name != "" {
			e.globals.names[id] = name
		}
	}

	return rows.Err()
}

// loadCharSkills loads character skill data.
func (e *Engine) loadCharSkills() error {
	rows, err := e.db.Query(`
		SELECT char_id, skill_id, level, experience
		FROM char_skills
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var charID, skillID, level, experience int

		if err := rows.Scan(&charID, &skillID, &level, &experience); err != nil {
			return fmt.Errorf("scan char_skill: %w", err)
		}

		if e.globals.bx[charID] == nil {
			continue
		}

		// Ensure x_char exists
		if e.globals.bx[charID].x_char == nil {
			e.globals.bx[charID].x_char = &entity_char{}
		}

		// Add skill to character's skill list
		skillEnt := &skill_ent{
			skill:        skillID,
			days_studied: level,
			experience:   short(experience),
			know:         SKILL_know,
		}

		// Append to skills list using the C-style plist pattern
		e.appendCharSkill(charID, skillEnt)
	}

	return rows.Err()
}

// loadSystemConfig loads system configuration from game_meta.
func (e *Engine) loadSystemConfig() error {
	row := e.db.QueryRow(`
		SELECT game_name, current_turn, options_json
		FROM game_meta
		WHERE id = 1
	`)

	var gameName string
	var currentTurn int
	var optionsJSON sql.NullString

	err := row.Scan(&gameName, &currentTurn, &optionsJSON)
	if err == sql.ErrNoRows {
		// No config yet, use defaults
		return nil
	}
	if err != nil {
		return fmt.Errorf("scan game_meta: %w", err)
	}

	// Set the game clock
	e.globals.sysclock.turn = short(currentTurn)

	return nil
}

// appendCharSkill appends a skill_ent to a character's skills list.
// This handles the C-style **skill_ent (plist) pattern.
func (e *Engine) appendCharSkill(charID int, sk *skill_ent) {
	ch := e.globals.bx[charID].x_char
	if ch == nil {
		return
	}

	// The skills field is **skill_ent, which in C represents a growable array.
	// We use unsafe to cast between **skill_ent and *[]*skill_ent.
	// For simplicity during the port, we store skills in a separate map.
	if e.globals.charSkills == nil {
		e.globals.charSkills = make(map[int][]*skill_ent)
	}
	e.globals.charSkills[charID] = append(e.globals.charSkills[charID], sk)
}

// getCharSkills returns the skills for a character.
func (e *Engine) getCharSkills(charID int) []*skill_ent {
	if e.globals.charSkills == nil {
		return nil
	}
	return e.globals.charSkills[charID]
}
