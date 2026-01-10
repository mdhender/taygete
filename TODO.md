# Olympia C→Go Port Plan

## Overview
- **Source**: 55 C files, ~52k LOC in `src/`
- **Target**: Single Go package (refactor after port)
- **Storage**: SQLite3 (one DB per game)
- **Frontend**: Next.js + Tailwind "Oatmeal" UI
- **Sprints**: ~50 (reduced from 52 by skipping I/O parsing/reporting)

Note that parsing functions like `readfile`, `eat_leading_trailing_whitespace` may not be needed since we will be saving orders and game state to the database instead of flat files. Deprecated functions should be created with the expected signatures so we can track progress. Deprecated functions must include the "// Deprecated: ...." comment. They should panic with a status message if called.

---

## Phase 0 – Setup (S0)
- [x] Go module structure, `db` layer scaffolding, test harnesses

---

## Phase 1 – Core Utilities (S1–S4)

### Sprint 1: z.c core (~1k LOC)
- [x] Port memory helpers (skipped - Go GC), `ilist`/`plist`, `assert`, `readfile` (deprecated), string trim/comparators
- [x] Tests for `ilist`/`plist`, `i_strcmp`, `fuzzy_strcmp`

### Sprint 2: z.c remaining
- [x] `lower_array`/`init_lower`, character classification
- [x] Tests for case folding & fuzzy comparisons

### Sprint 3: rnd.c
- [x] Port MD5-based RNG (`rnd.go`)
- [x] Tests comparing to C RNG sequence (`rnd_test.go`)

### Sprint 4: sout.c
- [x] SKIP – string formatting for text reports; Next.js renders from DB

---

## Phase 2 – Types & Globals (S5–S8)

### Sprint 5: types.go
- [x] All `T_*`, `sub_*`, `item_*`, `sk_*` constants from `oly.h`
- [x] Basic `Box` structure and substructure skeletons

### Sprint 6: accessor functions
- [x] `kind`, `subkind`, `valid_box`, `rp_*`/`p_*` analogs
- [x] Tests for accessor correctness

### Sprint 7: glob.c
- [x] Port `glob_init`, `sysclock`, `boxHead`, `subHead`, `loop_*` helpers
- [x] Tests for initialization

### Sprint 8: code.c
- [x] Port `int_to_code`, `code_to_int`, naming functions
- [x] Tests for ID encoding/decoding

---

## Phase 3 – Spatial Model (S9–S13)

### Sprint 9: loc.c (ownership)
- [x] Port `loc_owner`, `region`, `province`, `subloc`
- [x] Minimal tests with mock world

### Sprint 10: loc.c (here-lists)
- [x] Port `all_here`, `in_safe_now`, `subloc_here`, `count_loc_structures`
- [x] Tests for here-list correctness

### Sprint 11: stack.c
- [x] Port stacking logic, movement of stacks
- [x] Tests for nested location/stack layout

### Sprint 12: gate.c
- [x] Port gates & roads
- [x] Tests for `road_dest`, `road_hidden`

### Sprint 13: storm.c/cloud.c
- [x] Port `ship_moving`, `ship_gone`, `char_moving`, `char_gone`
- [x] Tests for movement timing

---

## Phase 4 – Persistence (S14–S18)

### Sprint 14: DB schema
- [x] Implement schema migrations in Go
- [x] `OpenGameDB`, `BeginTurn`, `CommitTurn` helpers

### Sprint 15: LoadWorld
- [x] Map DB rows → `bx`/substructures
- [x] Tests with small world fixtures

### Sprint 16: SaveWorld
- [x] Map `bx` → DB
- [x] Round-trip tests

### Sprint 17: io.go
- [x] Item types loader/saver (`loadItemTypes`, `saveItemTypes`)
- [x] Skills loader/saver (`loadSkills`, `saveSkills`)
- [x] Character skills loader/saver (`loadCharSkills`, `saveCharSkills`)
- [x] System config loader (`loadSystemConfig`)
- [x] New schema: `turn_logs`, `skill_lore`, `system_config` tables
- [x] Tests for new loaders

### Sprint 18: check_db

NOTE: Full skill tree validation (offered/research lists) requires refactoring entity_skill to use Go slices instead of C-style ilist pointers.

- [x] Implement consistency checks
- [x] Tests for integrity issues
- [x] Refactor entity_skill to use Go slices instead of C-style ilist pointers.

---

## Phase 5 – Turn Engine & Orders (S19–S26)

### Sprint 19: input.c
- [x] SKIP – browser submits structured orders directly to DB

### Sprint 20: order.c
- [ ] In-memory order representation, scheduling (structure only, no text parsing)
- [ ] Load orders from DB `orders` table

### Sprint 21: day.c skeleton
- [ ] `process_orders`, `post_month` with stubbed handlers
- [ ] Tests for no-op turn processing

### Sprint 22: Command lifecycle
- [ ] `STATE_LOAD`/`RUN`/`ERROR`/`DONE` mapping
- [ ] Tests for scheduling/priorities

### Sprint 23: immed.c
- [ ] Immediate commands
- [ ] Tests for immediate operations

### Sprint 24: check.c
- [ ] Full port, integrate into turn end

### Sprint 25–26: u.c, c1.c/c2.c subset
- [ ] Helper modules and easiest command implementations
- [ ] Unit tests for those commands

---

## Phase 6 – Gameplay Subsystems (S27–S44)

### Movement & World (S27–S30)
- [ ] S27: `move.c` core movement
- [ ] S28: `dir.c` region/path utilities
- [ ] S29: `faery.c`, `hades.c` special regions
- [ ] S30: `tunnel.c` finishing edge cases

### Economy & Construction (S31–S34)
- [ ] S31: `basic.c` economic foundations
- [ ] S32: `build.c` building creation/ownership
- [ ] S33: `buy.c` trade interactions
- [ ] S34: `produce.c`, `make.c` crafting/production

### Combat & Stealth (S35–S38)
- [ ] S35: `combat.c` core battle resolution
- [ ] S36: `beast.c`, `savage.c` special combat/mobs
- [ ] S37: `stealth.c`, `scry.c` mechanics
- [ ] S38: `garr.c`, `npc.c` garrison & NPC AI

### Magic & Special (S39–S42)
- [ ] S39: `alchem.c` alchemy & items
- [ ] S40: `necro.c` necromancy
- [ ] S41: `lore.c` knowledge tracking
- [ ] S42: `relig.c`, `art.c`, `quest.c` religion/artifacts/quests

### GM & Meta (S43–S44)
- [ ] S43: `gm.c`, `perm.c` GM tools
- [ ] S44: `add.c`, `pw.c` player onboarding & accounts

---

## Phase 7 – Web API (S45–S48)

### Sprint 45–46: HTTP API core
- [ ] Go HTTP server setup
- [ ] Player login/session endpoints (from `accounts`/`players` tables)
- [ ] Order submission endpoint (writes to `orders` table)

### Sprint 47–48: Game data endpoints
- [ ] Turn results/game state queries for Next.js
- [ ] Player-specific data (what they can see)
- [ ] Integrate with Next.js "Oatmeal" frontend

*Note: display.c, report.c, summary.c SKIPPED – Next.js renders reports from DB*

---

## Phase 8 – CLI & Cleanup (S49–S50)

### Sprint 49: cmd/olympia/main.go
- [ ] CLI for running turns, DB management
- [ ] Tests for CLI

### Sprint 50: Final cleanup
- [ ] Refactor global state to cleaner Go patterns
- [ ] Documentation and integration tests

---

## SQLite Schema

### Meta & RNG
```sql
CREATE TABLE game_meta (
  id            INTEGER PRIMARY KEY CHECK (id = 1),
  game_name     TEXT NOT NULL,
  created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  current_turn  INTEGER NOT NULL DEFAULT 0,
  rules_version TEXT,
  options_json  TEXT
);

CREATE TABLE rng_state (
  id        INTEGER PRIMARY KEY CHECK (id = 1),
  seed_blob BLOB NOT NULL
);
```

### Accounts & Players
```sql
CREATE TABLE accounts (
  id              INTEGER PRIMARY KEY AUTOINCREMENT,
  email           TEXT NOT NULL UNIQUE,
  password_hash   TEXT NOT NULL,
  full_name       TEXT,
  status          TEXT NOT NULL DEFAULT 'active',
  created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_login_at   DATETIME
);

CREATE TABLE players (
  id                INTEGER PRIMARY KEY,
  account_id        INTEGER REFERENCES accounts(id),
  code              TEXT NOT NULL,
  name              TEXT,
  banner            TEXT,
  subkind           INTEGER NOT NULL,
  email             TEXT,
  vis_email         TEXT,
  public_turn       INTEGER DEFAULT 0,
  report_format     INTEGER DEFAULT 0,
  notab             INTEGER DEFAULT 0,
  times_paid        INTEGER DEFAULT 0,
  is_system         INTEGER NOT NULL DEFAULT 0,
  acct_balance      INTEGER NOT NULL DEFAULT 0,
  acct_status       TEXT NOT NULL DEFAULT 'ok',
  settings_json     TEXT
);
```

### Entities (Core)
```sql
CREATE TABLE entities (
  id              INTEGER PRIMARY KEY,
  kind            INTEGER NOT NULL,
  subkind         INTEGER NOT NULL,
  name            TEXT,
  display_name    TEXT,
  created_turn    INTEGER NOT NULL DEFAULT 0,
  owner_player_id INTEGER REFERENCES players(id),
  parent_loc_id   INTEGER,
  is_deleted      INTEGER NOT NULL DEFAULT 0,
  extra           TEXT
);
```

### Locations
```sql
CREATE TABLE locations (
  id               INTEGER PRIMARY KEY REFERENCES entities(id),
  region_id        INTEGER,
  province_id      INTEGER,
  parent_loc_id    INTEGER,
  terrain_subkind  INTEGER NOT NULL,
  barrier          INTEGER DEFAULT 0,
  shroud           INTEGER DEFAULT 0,
  civ              INTEGER DEFAULT 0,
  sea_lane         INTEGER DEFAULT 0,
  prominence       INTEGER DEFAULT 0,
  opium_econ       INTEGER DEFAULT 0,
  gate_dist        INTEGER DEFAULT 0,
  is_safe_haven    INTEGER DEFAULT 0,
  is_start_loc     INTEGER DEFAULT 0,
  extra            TEXT
);

CREATE INDEX idx_locations_region ON locations(region_id);
CREATE INDEX idx_locations_parent ON locations(parent_loc_id);
```

### Characters
```sql
CREATE TABLE characters (
  id              INTEGER PRIMARY KEY REFERENCES entities(id),
  player_id       INTEGER REFERENCES players(id),
  loc_id          INTEGER REFERENCES locations(id),
  health          INTEGER NOT NULL DEFAULT 100,
  sick            INTEGER NOT NULL DEFAULT 0,
  loy_kind        INTEGER,
  loy_rate        INTEGER,
  unit_item       INTEGER,
  guard           INTEGER DEFAULT 0,
  npc_prog        INTEGER,
  studied         INTEGER,
  moving_since    INTEGER,
  gone_flag       INTEGER DEFAULT 0,
  is_npc          INTEGER NOT NULL DEFAULT 0,
  is_unformed     INTEGER NOT NULL DEFAULT 0,
  is_dead         INTEGER NOT NULL DEFAULT 0,
  extra           TEXT
);

CREATE TABLE char_magic (
  char_id        INTEGER PRIMARY KEY REFERENCES characters(id),
  pray           INTEGER DEFAULT 0,
  hide_self      INTEGER DEFAULT 0,
  vis_protect    INTEGER DEFAULT 0,
  default_garr   INTEGER,
  hide_mage      INTEGER DEFAULT 0,
  cur_aura       INTEGER DEFAULT 0,
  max_aura       INTEGER DEFAULT 0,
  aura_reflect   INTEGER DEFAULT 0,
  pledge         INTEGER,
  auraculum      INTEGER,
  ability_shroud INTEGER DEFAULT 0,
  fee            INTEGER DEFAULT 0,
  ferry_flag     INTEGER DEFAULT 0,
  extra          TEXT
);
```

### Items & Inventory
```sql
CREATE TABLE item_types (
  id             INTEGER PRIMARY KEY,
  subkind        INTEGER NOT NULL,
  name           TEXT NOT NULL,
  weight         INTEGER DEFAULT 0,
  is_animal      INTEGER DEFAULT 0,
  prominent      INTEGER DEFAULT 0,
  extra          TEXT
);

CREATE TABLE inventories (
  owner_entity_id INTEGER NOT NULL REFERENCES entities(id),
  item_id         INTEGER NOT NULL REFERENCES item_types(id),
  qty             INTEGER NOT NULL,
  PRIMARY KEY (owner_entity_id, item_id)
);
```

### Skills
```sql
CREATE TABLE skills (
  id           INTEGER PRIMARY KEY,
  name         TEXT NOT NULL,
  category     TEXT,
  is_magic     INTEGER DEFAULT 0,
  extra        TEXT
);

CREATE TABLE char_skills (
  char_id        INTEGER NOT NULL REFERENCES characters(id),
  skill_id       INTEGER NOT NULL REFERENCES skills(id),
  level          INTEGER NOT NULL DEFAULT 0,
  experience     INTEGER NOT NULL DEFAULT 0,
  last_studied   INTEGER,
  PRIMARY KEY (char_id, skill_id)
);
```

### Gates, Ships, Storms
```sql
CREATE TABLE gates (
  id           INTEGER PRIMARY KEY REFERENCES entities(id),
  from_loc_id  INTEGER NOT NULL REFERENCES locations(id),
  to_loc_id    INTEGER NOT NULL REFERENCES locations(id),
  road_hidden  INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE ships (
  id           INTEGER PRIMARY KEY REFERENCES entities(id),
  loc_id       INTEGER REFERENCES locations(id),
  capacity     INTEGER DEFAULT 0,
  storm_bind   INTEGER REFERENCES entities(id),
  moving_since INTEGER,
  extra        TEXT
);

CREATE TABLE storms (
  id           INTEGER PRIMARY KEY REFERENCES entities(id),
  strength     INTEGER DEFAULT 0,
  moving_to    INTEGER REFERENCES locations(id),
  moving_since INTEGER,
  extra        TEXT
);
```

### Turns, Orders, Commands
```sql
CREATE TABLE turns (
  turn_number   INTEGER PRIMARY KEY,
  started_at    DATETIME,
  finished_at   DATETIME,
  status        TEXT NOT NULL DEFAULT 'pending'
);

CREATE TABLE orders (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  turn_number   INTEGER NOT NULL REFERENCES turns(turn_number),
  player_id     INTEGER NOT NULL REFERENCES players(id),
  source_char_id INTEGER REFERENCES characters(id),
  raw_text      TEXT NOT NULL,
  received_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  source_channel TEXT,
  extra         TEXT
);

CREATE TABLE commands (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  turn_number   INTEGER NOT NULL REFERENCES turns(turn_number),
  who_id        INTEGER NOT NULL REFERENCES entities(id),
  cmd_code      INTEGER NOT NULL,
  use_skill     INTEGER,
  use_ent       INTEGER,
  use_exp       INTEGER,
  days_executing INTEGER NOT NULL DEFAULT 0,
  state         INTEGER NOT NULL,
  status        INTEGER NOT NULL,
  poll          INTEGER NOT NULL DEFAULT 0,
  pri           INTEGER NOT NULL DEFAULT 0,
  conditional   INTEGER NOT NULL DEFAULT 0,
  inhibit_finish INTEGER NOT NULL DEFAULT 0,
  args_json     TEXT,
  raw_line      TEXT,
  created_at    DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Combat & Reports
```sql
CREATE TABLE combats (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  turn_number  INTEGER NOT NULL REFERENCES turns(turn_number),
  loc_id       INTEGER NOT NULL REFERENCES locations(id),
  started_at   DATETIME,
  winner_side  TEXT,
  summary      TEXT,
  log_text     TEXT
);

CREATE TABLE combat_participants (
  combat_id    INTEGER NOT NULL REFERENCES combats(id),
  char_id      INTEGER REFERENCES characters(id),
  side         TEXT NOT NULL,
  casualties   INTEGER DEFAULT 0,
  survived     INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (combat_id, char_id)
);

CREATE TABLE reports (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  turn_number  INTEGER NOT NULL REFERENCES turns(turn_number),
  player_id    INTEGER NOT NULL REFERENCES players(id),
  format       TEXT NOT NULL DEFAULT 'text',
  body         TEXT NOT NULL,
  created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

---

## Notes

### Porting Order (Dependencies)
1. `z.c` → utilities (everything depends on this)
2. `rnd.c` → RNG (used everywhere)
3. `oly.h` types → `types.go`
4. `glob.c` → globals initialization
5. `code.c` → ID encoding/naming
6. `loc.c`, `stack.c` → spatial model
7. DB layer → persistence
8. Turn engine → `input.c`, `order.c`, `day.c`
9. Gameplay subsystems → by domain
10. Reporting → web API

### Key C Headers
- `src/oly.h` – main types/constants
- `src/code.h` – ID encoding
- `src/loc.h` – location functions
- `src/z.h` – utilities
