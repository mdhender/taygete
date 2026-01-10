// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

// glob.go ports glob.c - global variables and initialization

// Game configuration globals (from glob.c)
const (
	from_host        = "t3@taygete.example.com"
	reply_host       = "t3@taygete.example.com"
	gm_address       = "gm@taygete.example.com"
	game_title       = "Olympia T3"
	game_url         = "https://taygete.example.com/t3"
	rules_url        = "https://taygete.example.com/t3/node/11"
	times_url        = "https://taygete.example.com/t3/node/409"
	htpasswd_loc     = "/var/www/com.example.taygete.t3/lib/ht-passwords"
	garrison_pay     = 2
	army_slow_factor = 20
	auto_quit_turns  = 0
)

// Special player IDs
const (
	indep_player = 100 // independent player
	gm_player    = 200 // The Fates
	skill_player = 202 // skill listing
	eat_pl       = 203 // Order scanner
	npc_pl       = 206 // Subloc monster player
	garr_pl      = 207 // Garrison unit owner
)

const T_deleted = 0 // forget on save
const T_player = 1
const T_char = 2
const T_loc = 3
const T_item = 4
const T_skill = 5
const T_gate = 6
const T_road = 7
const T_deadchar = 8
const T_ship = 9
const T_post = 10
const T_storm = 11
const T_unform = 12 // unformed noble
const T_lore = 13
const T_MAX = T_lore + 1 // one past highest T_xxx define

// String tables for display
var kind_s = []string{ //  (treat as const)
	"deleted",  // T_deleted
	"player",   // T_player
	"char",     // T_char
	"loc",      // T_loc
	"item",     // T_item
	"skill",    // T_skill
	"gate",     // T_gate
	"road",     // T_road
	"deadchar", // T_deadchar
	"ship",     // T_ship
	"post",     // T_post
	"storm",    // T_storm
	"unform",   // T_unform
	"lore",     // T_lore
}

const sub_ocean = 1
const sub_forest = 2
const sub_plain = 3
const sub_mountain = 4
const sub_desert = 5
const sub_swamp = 6
const sub_under = 7      // underground
const sub_faery_hill = 8 // gateway to Faery
const sub_island = 9     // island subloc
const sub_stone_cir = 10 // ring of stones
const sub_mallorn_grove = 11
const sub_bog = 12
const sub_cave = 13
const sub_city = 14
const sub_lair = 15 // dragon lair
const sub_graveyard = 16
const sub_ruins = 17
const sub_battlefield = 18
const sub_ench_forest = 19 // enchanted forest
const sub_rocky_hill = 20
const sub_tree_circle = 21
const sub_pits = 22
const sub_pasture = 23
const sub_oasis = 24
const sub_yew_grove = 25
const sub_sand_pit = 26
const sub_sacred_grove = 27
const sub_poppy_field = 28
const sub_temple = 29
const sub_galley = 30
const sub_roundship = 31
const sub_castle = 32
const sub_galley_notdone = 33
const sub_roundship_notdone = 34
const sub_ghost_ship = 35
const sub_temple_notdone = 36
const sub_inn = 37
const sub_inn_notdone = 38
const sub_castle_notdone = 39
const sub_mine = 40
const sub_mine_notdone = 41
const sub_scroll = 42 // item is a scroll
const sub_magic = 43  // this skill is magical
const sub_palantir = 44
const sub_auraculum = 45
const sub_tower = 46
const sub_tower_notdone = 47
const sub_pl_system = 48  // system player
const sub_pl_regular = 49 // regular player
const sub_region = 50     // region wrapper loc
const sub_pl_savage = 51  // Savage King
const sub_pl_npc = 52
const sub_mine_collapsed = 53
const sub_ni = 54        // ni=noble_item
const sub_undead = 55    // undead lord
const sub_dead_body = 56 // dead noble's body
const sub_fog = 57
const sub_wind = 58
const sub_rain = 59
const sub_hades_pit = 60
const sub_artifact = 61
const sub_pl_silent = 62
const sub_npc_token = 63 // npc group control art
const sub_garrison = 64  // npc group control art
const sub_cloud = 65     // cloud terrain type
const sub_raft = 66      // raft made out of flotsam
const sub_raft_notdone = 67
const sub_suffuse_ring = 68
const sub_relic = 69 // 400 series artifacts
const sub_tunnel = 70
const sub_sewer = 71
const sub_chamber = 72
const sub_tradegood = 73
const SUB_MAX = sub_tradegood + 1 // one past highest sub_

var subkind_s = []string{ //  (treat as const)
	"<no subkind>",
	"ocean",                 // sub_ocean
	"forest",                // sub_forest
	"plain",                 // sub_plain
	"mountain",              // sub_mountain
	"desert",                // sub_desert
	"swamp",                 // sub_swamp
	"underground",           // sub_under
	"faery hill",            // sub_faery_hill
	"island",                // sub_island
	"ring of stones",        // sub_stone_cir
	"mallorn grove",         // sub_mallorn_grove
	"bog",                   // sub_bog
	"cave",                  // sub_cave
	"city",                  // sub_city
	"lair",                  // sub_lair
	"graveyard",             // sub_graveyard
	"ruins",                 // sub_ruins
	"battlefield",           // sub_battlefield
	"enchanted forest",      // sub_ench_forest
	"rocky hill",            // sub_rocky_hill
	"circle of trees",       // sub_tree_cir
	"pits",                  // sub_pits
	"pasture",               // sub_pasture
	"oasis",                 // sub_oasis
	"yew grove",             // sub_yew_grove
	"sand pit",              // sub_sand_pit
	"sacred grove",          // sub_sacred_grove
	"poppy field",           // sub_poppy_field
	"temple",                // sub_temple
	"galley",                // sub_galley
	"roundship",             // sub_roundship
	"castle",                // sub_castle
	"galley-in-progress",    // sub_galley_notdone
	"roundship-in-progress", // sub_roundship_notdone
	"ghost ship",            // sub_ghost_ship
	"temple-in-progress",    // sub_temple_notdone
	"inn",                   // sub_inn
	"inn-in-progress",       // sub_inn_notdone
	"castle-in-progress",    // sub_castle_notdone
	"mine",                  // sub_mine
	"mine-in-progress",      // sub_mine_notdone
	"scroll",                // sub_scroll
	"magic",                 // sub_magic
	"palantir",              // sub_palantir
	"auraculum",             // sub_auraculum
	"tower",                 // sub_tower
	"tower-in-progress",     // sub_tower_notdone
	"pl_system",             // sub_pl_system
	"pl_regular",            // sub_pl_regular
	"region",                // sub_region
	"pl_savage",             // sub_pl_savage
	"pl_npc",                // sub_pl_npc
	"collapsed mine",        // sub_mine_collapsed
	"ni",                    // sub_ni
	"demon lord",            // sub_undead
	"dead body",             // sub_dead_body
	"fog",                   // sub_fog
	"wind",                  // sub_wind
	"rain",                  // sub_rain
	"pit",                   // sub_hades_pit
	"artifact",              // sub_artifact
	"pl_silent",             // sub_pl_silent
	"npc_token",             // sub_npc_token
	"garrison",              // sub_garrison
	"cloud",                 // sub_cloud
	"raft",                  // sub_raft
	"raft-in-progress",      // sub_raft_notdone
	"suffuse_ring",          // sub_suffuse_ring
	"relic",                 // sub_relic
	"tunnel",                // sub_tunnel
	"sewer",                 // sub_sewer
	"chamber",               // sub_chamber
	"tradegood",             // sub_tradegood
}

const DIR_N = 1
const DIR_E = 2
const DIR_S = 3
const DIR_W = 4
const DIR_UP = 5
const DIR_DOWN = 6
const DIR_IN = 7
const DIR_OUT = 8
const MAX_DIR = DIR_OUT + 1 // one past highest direction

var short_dir_s = []string{ //  (treat as const)
	"<no dir>",
	"n",
	"e",
	"s",
	"w",
	"u",
	"d",
	"i",
	"o",
}

var full_dir_s = []string{ //  (treat as const)
	"<no dir>",
	"north",
	"east",
	"south",
	"west",
	"up",
	"down",
	"in",
	"out",
}

var exit_opposite = []int{ //  (treat as const)
	0,
	DIR_S,
	DIR_W,
	DIR_N,
	DIR_E,
	DIR_OUT,
	DIR_IN,
	DIR_DOWN,
	DIR_UP,
	0,
}

const LOC_region = 1   // top most continent/island group
const LOC_province = 2 // main location area
const LOC_subloc = 3   // inner sublocation
const LOC_build = 4    // building, structure, etc.

var loc_depth_s = []string{ //  (treat as const)
	"<no depth>",
	"region",
	"province",
	"subloc",
}

const NUM_MONTHS = 8
const MONTH_DAYS = 30

var month_names = []string{ //  (treat as const)
	"Fierce winds",     // 0
	"Snowmelt",         // 1
	"Blossom bloom",    // 2
	"Sunsear",          // 3
	"Thunder and rain", // 4
	"Harvest",          // 5
	"Waning days",      // 6
	"Dark night",       // 7
}

// glob_init initializes global game state.
// In Go, the bx array is pre-allocated in Engine.globals.
// This function resets the box_head and sub_head chains.
func glob_init() {
	for i := 0; i < T_MAX; i++ {
		teg.globals.box_head[i] = 0
	}
	for i := 0; i < SUB_MAX; i++ {
		teg.globals.sub_head[i] = 0
	}
}

// GlobInit initializes global game state on the Engine.
func (e *Engine) GlobInit() {
	for i := 0; i < T_MAX; i++ {
		e.globals.box_head[i] = 0
	}
	for i := 0; i < SUB_MAX; i++ {
		e.globals.sub_head[i] = 0
	}
}

// Sysclock returns the current game time.
func (e *Engine) Sysclock() olytime {
	return e.globals.sysclock
}

// SetSysclock sets the current game time.
func (e *Engine) SetSysclock(t olytime) {
	e.globals.sysclock = t
}

// KindFirst returns the first entity of the given kind.
func (e *Engine) KindFirst(k int) int {
	if k < 0 || k >= T_MAX {
		return 0
	}
	return e.globals.box_head[k]
}

// KindNext returns the next entity of the same kind.
func (e *Engine) KindNext(id int) int {
	if id <= 0 || id >= MAX_BOXES {
		return 0
	}
	b := e.globals.bx[id]
	if b == nil {
		return 0
	}
	return b.x_next_kind
}

// SubFirst returns the first entity of the given subkind.
func (e *Engine) SubFirst(sk int) int {
	if sk < 0 || sk >= SUB_MAX {
		return 0
	}
	return e.globals.sub_head[sk]
}

// SubNext returns the next entity of the same subkind.
func (e *Engine) SubNext(id int) int {
	if id <= 0 || id >= MAX_BOXES {
		return 0
	}
	b := e.globals.bx[id]
	if b == nil {
		return 0
	}
	return b.x_next_sub
}

// Loop helper functions - return slices for idiomatic Go iteration.
// See sprints/LOOPS.md for the design rationale.

// Characters returns all character entity IDs.
func (e *Engine) Characters() []int {
	var result []int
	for id := e.KindFirst(T_char); id > 0; id = e.KindNext(id) {
		result = append(result, id)
	}
	return result
}

// Players returns all player entity IDs.
func (e *Engine) Players() []int {
	var result []int
	for id := e.KindFirst(T_player); id > 0; id = e.KindNext(id) {
		result = append(result, id)
	}
	return result
}

// Locations returns all location entity IDs.
func (e *Engine) Locations() []int {
	var result []int
	for id := e.KindFirst(T_loc); id > 0; id = e.KindNext(id) {
		result = append(result, id)
	}
	return result
}

// Items returns all item entity IDs.
func (e *Engine) Items() []int {
	var result []int
	for id := e.KindFirst(T_item); id > 0; id = e.KindNext(id) {
		result = append(result, id)
	}
	return result
}

// Skills returns all skill entity IDs.
func (e *Engine) Skills() []int {
	var result []int
	for id := e.KindFirst(T_skill); id > 0; id = e.KindNext(id) {
		result = append(result, id)
	}
	return result
}

// Gates returns all gate entity IDs.
func (e *Engine) Gates() []int {
	var result []int
	for id := e.KindFirst(T_gate); id > 0; id = e.KindNext(id) {
		result = append(result, id)
	}
	return result
}

// Ships returns all ship entity IDs.
func (e *Engine) Ships() []int {
	var result []int
	for id := e.KindFirst(T_ship); id > 0; id = e.KindNext(id) {
		result = append(result, id)
	}
	return result
}

// Storms returns all storm entity IDs.
func (e *Engine) Storms() []int {
	var result []int
	for id := e.KindFirst(T_storm); id > 0; id = e.KindNext(id) {
		result = append(result, id)
	}
	return result
}

// Cities returns all city sublocation IDs.
func (e *Engine) Cities() []int {
	var result []int
	for id := e.SubFirst(sub_city); id > 0; id = e.SubNext(id) {
		result = append(result, id)
	}
	return result
}

// Garrisons returns all garrison entity IDs.
func (e *Engine) Garrisons() []int {
	var result []int
	for id := e.SubFirst(sub_garrison); id > 0; id = e.SubNext(id) {
		result = append(result, id)
	}
	return result
}

// Castles returns all castle sublocation IDs.
func (e *Engine) Castles() []int {
	var result []int
	for id := e.SubFirst(sub_castle); id > 0; id = e.SubNext(id) {
		result = append(result, id)
	}
	return result
}

// Mountains returns all mountain sublocation IDs.
func (e *Engine) Mountains() []int {
	var result []int
	for id := e.SubFirst(sub_mountain); id > 0; id = e.SubNext(id) {
		result = append(result, id)
	}
	return result
}

// Inns returns all inn sublocation IDs.
func (e *Engine) Inns() []int {
	var result []int
	for id := e.SubFirst(sub_inn); id > 0; id = e.SubNext(id) {
		result = append(result, id)
	}
	return result
}

// Temples returns all temple sublocation IDs.
func (e *Engine) Temples() []int {
	var result []int
	for id := e.SubFirst(sub_temple); id > 0; id = e.SubNext(id) {
		result = append(result, id)
	}
	return result
}

// CollapsedMines returns all collapsed mine sublocation IDs.
func (e *Engine) CollapsedMines() []int {
	var result []int
	for id := e.SubFirst(sub_mine_collapsed); id > 0; id = e.SubNext(id) {
		result = append(result, id)
	}
	return result
}

// DeadBodies returns all dead body sublocation IDs.
func (e *Engine) DeadBodies() []int {
	var result []int
	for id := e.SubFirst(sub_dead_body); id > 0; id = e.SubNext(id) {
		result = append(result, id)
	}
	return result
}

// RegularPlayers returns all regular player IDs.
func (e *Engine) RegularPlayers() []int {
	var result []int
	for id := e.SubFirst(sub_pl_regular); id > 0; id = e.SubNext(id) {
		result = append(result, id)
	}
	return result
}

// Provinces returns all province location IDs.
func (e *Engine) Provinces() []int {
	var result []int
	for id := e.KindFirst(T_loc); id > 0; id = e.KindNext(id) {
		if e.LocDepth(id) == LOC_province {
			result = append(result, id)
		}
	}
	return result
}

// LocsAndShips returns all location and ship entity IDs.
func (e *Engine) LocsAndShips() []int {
	var result []int
	for id := e.KindFirst(T_ship); id > 0; id = e.KindNext(id) {
		result = append(result, id)
	}
	for id := e.KindFirst(T_loc); id > 0; id = e.KindNext(id) {
		result = append(result, id)
	}
	return result
}

// Boxes returns all valid (non-deleted) box IDs.
func (e *Engine) Boxes() []int {
	var result []int
	for id := 1; id < MAX_BOXES; id++ {
		if e.Kind(id) != T_deleted {
			result = append(result, id)
		}
	}
	return result
}

// Kind returns the kind of an entity.
// Returns T_deleted if the entity doesn't exist.
func (e *Engine) Kind(id int) schar {
	if id > 0 && id < MAX_BOXES && e.globals.bx[id] != nil {
		return e.globals.bx[id].kind
	}
	return T_deleted
}

// Subkind returns the subkind of an entity.
// Returns 0 if the entity doesn't exist.
func (e *Engine) Subkind(id int) schar {
	if id > 0 && id < MAX_BOXES && e.globals.bx[id] != nil {
		return e.globals.bx[id].skind
	}
	return 0
}

// ValidBox returns true if the entity exists (is not deleted).
func (e *Engine) ValidBox(id int) bool {
	return e.Kind(id) != T_deleted
}

// LocDepth returns the depth of a location (region, province, subloc, build).
// Returns 0 if the entity is not a location.
func (e *Engine) LocDepth(id int) int {
	return loc_depth(id)
}
