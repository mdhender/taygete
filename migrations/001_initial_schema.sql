--  taygete - a game engine for a game.
--  Copyright (c) 2026 Michael D Henderson.
--
--  This program is free software: you can redistribute it and/or modify
--  it under the terms of the GNU Affero General Public License as published by
--  the Free Software Foundation, either version 3 of the License, or
--  (at your option) any later version.
--
--  This program is distributed in the hope that it will be useful,
--  but WITHOUT ANY WARRANTY; without even the implied warranty of
--  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
--  GNU Affero General Public License for more details.
--
--  You should have received a copy of the GNU Affero General Public License
--  along with this program.  If not, see <https://www.gnu.org/licenses/>.
--
--  This program is free software: you can redistribute it and/or modify
--  it under the terms of the GNU Affero General Public License as published by
--  the Free Software Foundation, either version 3 of the License, or
--  (at your option) any later version.
--
--  This program is distributed in the hope that it will be useful,
--  but WITHOUT ANY WARRANTY; without even the implied warranty of
--  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
--  GNU Affero General Public License for more details.
--
--  You should have received a copy of the GNU Affero General Public License
--  along with this program.  If not, see <https://www.gnu.org/licenses/>.

-- Initial schema for Olympia game database

-- Meta & RNG
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

CREATE TABLE prng_state (
  name  TEXT PRIMARY KEY,
  state BLOB NOT NULL
);

CREATE TABLE passwords (
  key   TEXT PRIMARY KEY,
  value TEXT NOT NULL
);

-- Accounts & Players
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

-- Entities (Core)
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

-- Locations
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

-- Characters
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

-- Items & Inventory
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

-- Skills
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

-- Gates, Ships, Storms
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

-- Turns, Orders, Commands
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

-- Combat & Reports
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

-- Turn logs (from log/* files)
CREATE TABLE turn_logs (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  turn_number  INTEGER NOT NULL,
  player_id    INTEGER NOT NULL REFERENCES players(id),
  log_text     TEXT NOT NULL,
  created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_turn_logs_turn ON turn_logs(turn_number);
CREATE INDEX idx_turn_logs_player ON turn_logs(player_id);

-- Skill lore (from lore/* files)
CREATE TABLE skill_lore (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  skill_id     INTEGER NOT NULL REFERENCES skills(id),
  lore_text    TEXT NOT NULL,
  display_order INTEGER DEFAULT 0
);

CREATE INDEX idx_skill_lore_skill ON skill_lore(skill_id);

-- System configuration (from system file)
CREATE TABLE system_config (
  key          TEXT PRIMARY KEY,
  value        TEXT NOT NULL
);
