# Sprint 17: Flat File to Database Mappings

## Status: COMPLETE

Sprint 17 implements database loaders/savers for all flat file data that was previously read from `.upstream/Dist/lib/` files.

## Implementation Summary

### ✅ Priority 1: Add Missing Loaders (schema exists)

| Task | Status | Notes |
|------|--------|-------|
| `loadItemTypes()` | ✅ Done | Loads from `item_types` table into `entity_item` structs |
| `loadSkills()` | ✅ Done | Loads from `skills` table into `entity_skill` structs |
| `loadCharSkills()` | ✅ Done | Loads from `char_skills` table into character skill lists |

### ✅ Priority 2: System Config Loading

| Task | Status | Notes |
|------|--------|-------|
| `loadSystemConfig()` | ✅ Done | Loads from `game_meta` table, sets sysclock.turn |
| System config table | ✅ Done | Added `system_config` key-value table |

### ✅ Priority 3: New Schema

| Task | Status | Notes |
|------|--------|-------|
| `turn_logs` table | ✅ Done | For `log/*` files |
| `skill_lore` table | ✅ Done | For `lore/*` sheets |
| `system_config` table | ✅ Done | For `system` file key-value pairs |

### ✅ Save Functions

| Task | Status | Notes |
|------|--------|-------|
| `saveItemTypes()` | ✅ Done | Saves T_item entities to `item_types` |
| `saveSkills()` | ✅ Done | Saves T_skill entities to `skills` |
| `saveCharSkills()` | ✅ Done | Saves character skills from `charSkills` map |

---

## Flat File Coverage (Updated)

### ✅ Complete (schema + loader + saver)

| File | DB Table | Implementation |
|------|----------|----------------|
| `loc` | `entities` + `locations` | `loadEntities()` + `loadLocations()` |
| `gate` | `gates` | `loadGates()` + `saveGates()` |
| `road` | `gates` | Same as gate (road_hidden=1) |
| `ship` | `ships` | `loadShips()` + `saveShips()` |
| `players` | `players` | `loadPlayers()` + `savePlayers()` |
| `item` | `item_types` | `loadItemTypes()` + `saveItemTypes()` |
| `skill` | `skills` | `loadSkills()` + `saveSkills()` |
| `system` | `game_meta` + `system_config` | `loadSystemConfig()` |
| `randseed` | `rng_state`/`prng_state` | Existing RNG infrastructure |

### Schema Ready (loaders deferred)

| File/Directory | DB Table | Notes |
|----------------|----------|-------|
| `log/*` | `turn_logs` | Schema added, loader TBD |
| `lore/*` | `skill_lore` | Schema added, loader TBD |
| `fact/*` | `players` | Faction data in players table |
| `misc` (lore stubs) | `entities` | Can use T_lore entities |
| `email` | `accounts.email` | Covered by accounts table |
| `forward` | `players.email` | Covered by players table |
| `unform` | `characters.is_unformed` | Already in schema |
| `times_*` | `reports` | Already in schema |

---

## Technical Notes

### Character Skills Storage

The original C code uses `**skill_ent` (pointer to growable array) for character skills. To avoid unsafe pointer manipulation during the port, character skills are stored in a separate map:

```go
e.globals.charSkills map[int][]*skill_ent
```

Helper functions:
- `e.appendCharSkill(charID, skillEnt)` - Add skill to character
- `e.getCharSkills(charID)` - Get skills for character

### Schema Changes

Added to `migrations/001_initial_schema.sql`:
- `turn_logs` table with indexes
- `skill_lore` table with indexes  
- `system_config` key-value table

---

## Files Changed

- `load.go` - Added loaders + helper functions
- `save.go` - Added savers
- `engine.go` - Added `charSkills` map to globals
- `migrations/001_initial_schema.sql` - Added new tables
- `load_test.go` - Added tests for new loaders
