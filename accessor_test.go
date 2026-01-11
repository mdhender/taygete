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

import "testing"

func clearBx() {
	for i := range teg.globals.bx {
		teg.globals.bx[i] = nil
	}
}

func TestKindSubkindValidBox(t *testing.T) {
	defer clearBx()
	clearBx()

	if kind(0) != T_deleted {
		t.Errorf("kind(0) = %d, want T_deleted", kind(0))
	}
	if kind(-1) != T_deleted {
		t.Errorf("kind(-1) = %d, want T_deleted", kind(-1))
	}
	if kind(MAX_BOXES) != T_deleted {
		t.Errorf("kind(MAX_BOXES) = %d, want T_deleted", kind(MAX_BOXES))
	}
	if valid_box(100) {
		t.Error("valid_box(100) = true for nil box, want false")
	}

	teg.globals.bx[100] = &box{kind: T_char, skind: sub_ni}
	if kind(100) != T_char {
		t.Errorf("kind(100) = %d, want T_char", kind(100))
	}
	if subkind(100) != sub_ni {
		t.Errorf("subkind(100) = %d, want sub_ni", subkind(100))
	}
	if !valid_box(100) {
		t.Error("valid_box(100) = false for T_char box, want true")
	}

	teg.globals.bx[200] = &box{kind: T_loc, skind: sub_forest}
	if kind(200) != T_loc {
		t.Errorf("kind(200) = %d, want T_loc", kind(200))
	}
	if subkind(200) != sub_forest {
		t.Errorf("subkind(200) = %d, want sub_forest", subkind(200))
	}
}

func TestRpAccessorsNil(t *testing.T) {
	defer clearBx()
	clearBx()

	if rp_char(100) != nil {
		t.Error("rp_char(100) should be nil for nil box")
	}
	if rp_loc(100) != nil {
		t.Error("rp_loc(100) should be nil for nil box")
	}
	if rp_subloc(100) != nil {
		t.Error("rp_subloc(100) should be nil for nil box")
	}
	if rp_item(100) != nil {
		t.Error("rp_item(100) should be nil for nil box")
	}
	if rp_player(100) != nil {
		t.Error("rp_player(100) should be nil for nil box")
	}
	if rp_skill(100) != nil {
		t.Error("rp_skill(100) should be nil for nil box")
	}
	if rp_gate(100) != nil {
		t.Error("rp_gate(100) should be nil for nil box")
	}
	if rp_misc(100) != nil {
		t.Error("rp_misc(100) should be nil for nil box")
	}
	if rp_disp(100) != nil {
		t.Error("rp_disp(100) should be nil for nil box")
	}
	if rp_command(100) != nil {
		t.Error("rp_command(100) should be nil for nil box")
	}
	if rp_magic(100) != nil {
		t.Error("rp_magic(100) should be nil for nil box")
	}
	if rp_item_magic(100) != nil {
		t.Error("rp_item_magic(100) should be nil for nil box")
	}
	if rp_loc_info(100) != nil {
		t.Error("rp_loc_info(100) should be nil for nil box")
	}
}

func TestRpAccessorsWithSubstructures(t *testing.T) {
	defer clearBx()
	clearBx()

	teg.globals.bx[100] = &box{
		kind:  T_char,
		skind: sub_ni,
		x_char: &entity_char{
			health: 80,
			guard:  1,
			x_char_magic: &char_magic{
				cur_aura: 10,
				max_aura: 20,
			},
		},
	}

	c := rp_char(100)
	if c == nil {
		t.Fatal("rp_char(100) returned nil")
	}
	if c.health != 80 {
		t.Errorf("rp_char(100).health = %d, want 80", c.health)
	}

	m := rp_magic(100)
	if m == nil {
		t.Fatal("rp_magic(100) returned nil")
	}
	if m.cur_aura != 10 {
		t.Errorf("rp_magic(100).cur_aura = %d, want 10", m.cur_aura)
	}
	if m.max_aura != 20 {
		t.Errorf("rp_magic(100).max_aura = %d, want 20", m.max_aura)
	}

	teg.globals.bx[200] = &box{
		kind:  T_loc,
		skind: sub_forest,
		x_loc: &entity_loc{
			barrier: 5,
			shroud:  3,
			civ:     2,
		},
		x_subloc: &entity_subloc{
			prominence: 10,
			opium_econ: 50,
		},
	}

	l := rp_loc(200)
	if l == nil {
		t.Fatal("rp_loc(200) returned nil")
	}
	if l.barrier != 5 {
		t.Errorf("rp_loc(200).barrier = %d, want 5", l.barrier)
	}

	sl := rp_subloc(200)
	if sl == nil {
		t.Fatal("rp_subloc(200) returned nil")
	}
	if sl.prominence != 10 {
		t.Errorf("rp_subloc(200).prominence = %d, want 10", sl.prominence)
	}
}

func TestPAccessorsAllocate(t *testing.T) {
	defer clearBx()
	clearBx()

	c := p_char(100)
	if c == nil {
		t.Fatal("p_char(100) returned nil")
	}
	if teg.globals.bx[100] == nil {
		t.Error("p_char(100) did not allocate box")
	}
	if teg.globals.bx[100].x_char == nil {
		t.Error("p_char(100) did not allocate entity_char")
	}

	c.health = 75
	if rp_char(100).health != 75 {
		t.Error("p_char modification not reflected in rp_char")
	}

	m := p_magic(100)
	if m == nil {
		t.Fatal("p_magic(100) returned nil")
	}
	if teg.globals.bx[100].x_char.x_char_magic == nil {
		t.Error("p_magic(100) did not allocate char_magic")
	}
	m.cur_aura = 15
	if rp_magic(100).cur_aura != 15 {
		t.Error("p_magic modification not reflected in rp_magic")
	}

	l := p_loc(200)
	if l == nil {
		t.Fatal("p_loc(200) returned nil")
	}
	if teg.globals.bx[200] == nil || teg.globals.bx[200].x_loc == nil {
		t.Error("p_loc(200) did not allocate properly")
	}
	l.barrier = 7
	if rp_loc(200).barrier != 7 {
		t.Error("p_loc modification not reflected in rp_loc")
	}
}

func TestConvenienceAccessors(t *testing.T) {
	defer clearBx()
	clearBx()

	if char_health(100) != 0 {
		t.Error("char_health(100) should return 0 for nil box")
	}
	if loc_barrier(100) != 0 {
		t.Error("loc_barrier(100) should return 0 for nil box")
	}

	teg.globals.bx[100] = &box{
		kind: T_char,
		x_char: &entity_char{
			health:      90,
			sick:        1,
			guard:       1,
			loy_kind:    LOY_oath,
			loy_rate:    5,
			behind:      1,
			attack:      10,
			defense:     15,
			missile:     5,
			rank:        RANK_knight,
			break_point: 3,
			x_char_magic: &char_magic{
				pray:           1,
				hide_self:      1,
				cur_aura:       20,
				max_aura:       30,
				magician:       1,
				knows_weather:  1,
				ability_shroud: 2,
			},
		},
	}

	if char_health(100) != 90 {
		t.Errorf("char_health(100) = %d, want 90", char_health(100))
	}
	if char_sick(100) != 1 {
		t.Errorf("char_sick(100) = %d, want 1", char_sick(100))
	}
	if char_guard(100) != 1 {
		t.Errorf("char_guard(100) = %d, want 1", char_guard(100))
	}
	if loyal_kind(100) != LOY_oath {
		t.Errorf("loyal_kind(100) = %d, want LOY_oath", loyal_kind(100))
	}
	if loyal_rate(100) != 5 {
		t.Errorf("loyal_rate(100) = %d, want 5", loyal_rate(100))
	}
	if char_behind(100) != 1 {
		t.Errorf("char_behind(100) = %d, want 1", char_behind(100))
	}
	if char_attack(100) != 10 {
		t.Errorf("char_attack(100) = %d, want 10", char_attack(100))
	}
	if char_defense(100) != 15 {
		t.Errorf("char_defense(100) = %d, want 15", char_defense(100))
	}
	if char_missile(100) != 5 {
		t.Errorf("char_missile(100) = %d, want 5", char_missile(100))
	}
	if char_rank(100) != RANK_knight {
		t.Errorf("char_rank(100) = %d, want RANK_knight", char_rank(100))
	}
	if char_break(100) != 3 {
		t.Errorf("char_break(100) = %d, want 3", char_break(100))
	}
	if char_pray(100) != 1 {
		t.Errorf("char_pray(100) = %d, want 1", char_pray(100))
	}
	if char_hidden(100) != 1 {
		t.Errorf("char_hidden(100) = %d, want 1", char_hidden(100))
	}
	if char_cur_aura(100) != 20 {
		t.Errorf("char_cur_aura(100) = %d, want 20", char_cur_aura(100))
	}
	if char_max_aura(100) != 30 {
		t.Errorf("char_max_aura(100) = %d, want 30", char_max_aura(100))
	}
	if is_magician(100) != 1 {
		t.Errorf("is_magician(100) = %d, want 1", is_magician(100))
	}
	if weather_mage(100) != 1 {
		t.Errorf("weather_mage(100) = %d, want 1", weather_mage(100))
	}
	if char_abil_shroud(100) != 2 {
		t.Errorf("char_abil_shroud(100) = %d, want 2", char_abil_shroud(100))
	}

	teg.globals.bx[200] = &box{
		kind:  T_loc,
		skind: sub_forest,
		x_loc: &entity_loc{
			barrier:        10,
			shroud:         5,
			civ:            3,
			sea_lane:       1,
			dist_from_gate: 2,
		},
		x_subloc: &entity_subloc{
			prominence: 7,
			opium_econ: 100,
			capacity:   500,
			defense:    50,
			loot:       2,
			safe:       1,
			major:      1,
		},
	}

	if loc_barrier(200) != 10 {
		t.Errorf("loc_barrier(200) = %d, want 10", loc_barrier(200))
	}
	if loc_shroud(200) != 5 {
		t.Errorf("loc_shroud(200) = %d, want 5", loc_shroud(200))
	}
	if loc_civ(200) != 3 {
		t.Errorf("loc_civ(200) = %d, want 3", loc_civ(200))
	}
	if loc_sea_lane(200) != 1 {
		t.Errorf("loc_sea_lane(200) = %d, want 1", loc_sea_lane(200))
	}
	if gate_dist(200) != 2 {
		t.Errorf("gate_dist(200) = %d, want 2", gate_dist(200))
	}
	if loc_prominence(200) != 7 {
		t.Errorf("loc_prominence(200) = %d, want 7", loc_prominence(200))
	}
	if loc_opium(200) != 100 {
		t.Errorf("loc_opium(200) = %d, want 100", loc_opium(200))
	}
	if ship_cap_raw(200) != 500 {
		t.Errorf("ship_cap_raw(200) = %d, want 500", ship_cap_raw(200))
	}
	if loc_defense(200) != 50 {
		t.Errorf("loc_defense(200) = %d, want 50", loc_defense(200))
	}
	if loc_pillage(200) != 2 {
		t.Errorf("loc_pillage(200) = %d, want 2", loc_pillage(200))
	}
	if safe_haven(200) != 1 {
		t.Errorf("safe_haven(200) = %d, want 1", safe_haven(200))
	}
	if major_city(200) != 1 {
		t.Errorf("major_city(200) = %d, want 1", major_city(200))
	}
}

func TestItemAccessors(t *testing.T) {
	defer clearBx()
	clearBx()

	teg.globals.bx[300] = &box{
		kind:  T_item,
		skind: sub_artifact,
		x_item: &entity_item{
			weight:      10,
			land_cap:    100,
			ride_cap:    200,
			fly_cap:     300,
			attack:      5,
			defense:     3,
			missile:     2,
			is_man_item: 1,
			animal:      0,
			prominent:   1,
			capturable:  1,
			base_price:  500,
			who_has:     1234,
			x_item_magic: &item_magic{
				creator:        1000,
				region_created: 2000,
				lore:           lore_orb,
				use_key:        use_orb,
				attack_bonus:   2,
				defense_bonus:  3,
				missile_bonus:  1,
				aura_bonus:     5,
				token_num:      3,
				token_ni:       4000,
			},
		},
	}

	if item_weight(300) != 10 {
		t.Errorf("item_weight(300) = %d, want 10", item_weight(300))
	}
	if item_land_cap(300) != 100 {
		t.Errorf("item_land_cap(300) = %d, want 100", item_land_cap(300))
	}
	if item_ride_cap(300) != 200 {
		t.Errorf("item_ride_cap(300) = %d, want 200", item_ride_cap(300))
	}
	if item_fly_cap(300) != 300 {
		t.Errorf("item_fly_cap(300) = %d, want 300", item_fly_cap(300))
	}
	if item_attack(300) != 5 {
		t.Errorf("item_attack(300) = %d, want 5", item_attack(300))
	}
	if item_defense(300) != 3 {
		t.Errorf("item_defense(300) = %d, want 3", item_defense(300))
	}
	if item_missile(300) != 2 {
		t.Errorf("item_missile(300) = %d, want 2", item_missile(300))
	}
	if man_item(300) != 1 {
		t.Errorf("man_item(300) = %d, want 1", man_item(300))
	}
	if item_prominent(300) != 1 {
		t.Errorf("item_prominent(300) = %d, want 1", item_prominent(300))
	}
	if item_capturable(300) != 1 {
		t.Errorf("item_capturable(300) = %d, want 1", item_capturable(300))
	}
	if item_price(300) != 500 {
		t.Errorf("item_price(300) = %d, want 500", item_price(300))
	}
	if item_unique(300) != 1234 {
		t.Errorf("item_unique(300) = %d, want 1234", item_unique(300))
	}
	if item_creator(300) != 1000 {
		t.Errorf("item_creator(300) = %d, want 1000", item_creator(300))
	}
	if item_creat_loc(300) != 2000 {
		t.Errorf("item_creat_loc(300) = %d, want 2000", item_creat_loc(300))
	}
	if item_lore(300) != lore_orb {
		t.Errorf("item_lore(300) = %d, want lore_orb", item_lore(300))
	}
	if item_use_key(300) != use_orb {
		t.Errorf("item_use_key(300) = %d, want use_orb", item_use_key(300))
	}
	if item_attack_bonus(300) != 2 {
		t.Errorf("item_attack_bonus(300) = %d, want 2", item_attack_bonus(300))
	}
	if item_defense_bonus(300) != 3 {
		t.Errorf("item_defense_bonus(300) = %d, want 3", item_defense_bonus(300))
	}
	if item_missile_bonus(300) != 1 {
		t.Errorf("item_missile_bonus(300) = %d, want 1", item_missile_bonus(300))
	}
	if item_aura_bonus(300) != 5 {
		t.Errorf("item_aura_bonus(300) = %d, want 5", item_aura_bonus(300))
	}
	if item_token_num(300) != 3 {
		t.Errorf("item_token_num(300) = %d, want 3", item_token_num(300))
	}
	if item_token_ni(300) != 4000 {
		t.Errorf("item_token_ni(300) = %d, want 4000", item_token_ni(300))
	}
}

func TestGateAccessors(t *testing.T) {
	defer clearBx()
	clearBx()

	teg.globals.bx[400] = &box{
		kind:  T_gate,
		skind: 0,
		x_gate: &entity_gate{
			to_loc:      5000,
			seal_key:    123,
			road_hidden: 1,
		},
	}

	if gate_dest(400) != 5000 {
		t.Errorf("gate_dest(400) = %d, want 5000", gate_dest(400))
	}
	if road_dest(400) != 5000 {
		t.Errorf("road_dest(400) = %d, want 5000", road_dest(400))
	}
	if gate_seal(400) != 123 {
		t.Errorf("gate_seal(400) = %d, want 123", gate_seal(400))
	}
	if road_hidden(400) != 1 {
		t.Errorf("road_hidden(400) = %d, want 1", road_hidden(400))
	}
}

func TestShipPredicates(t *testing.T) {
	defer clearBx()
	clearBx()

	teg.globals.bx[100] = &box{kind: T_ship, skind: sub_galley}
	teg.globals.bx[200] = &box{kind: T_ship, skind: sub_roundship}
	teg.globals.bx[300] = &box{kind: T_ship, skind: sub_raft}
	teg.globals.bx[400] = &box{kind: T_ship, skind: sub_galley_notdone}
	teg.globals.bx[500] = &box{kind: T_loc, skind: sub_forest}

	if !is_ship(100) {
		t.Error("is_ship(100) should be true for galley")
	}
	if !is_ship(200) {
		t.Error("is_ship(200) should be true for roundship")
	}
	if !is_ship(300) {
		t.Error("is_ship(300) should be true for raft")
	}
	if is_ship(400) {
		t.Error("is_ship(400) should be false for galley_notdone")
	}
	if !is_ship_notdone(400) {
		t.Error("is_ship_notdone(400) should be true for galley_notdone")
	}
	if !is_ship_either(100) {
		t.Error("is_ship_either(100) should be true for galley")
	}
	if !is_ship_either(400) {
		t.Error("is_ship_either(400) should be true for galley_notdone")
	}
	if is_ship(500) {
		t.Error("is_ship(500) should be false for forest")
	}

	if !is_loc_or_ship(100) {
		t.Error("is_loc_or_ship(100) should be true for ship")
	}
	if !is_loc_or_ship(500) {
		t.Error("is_loc_or_ship(500) should be true for loc")
	}
}

func TestIsFighter(t *testing.T) {
	defer clearBx()
	clearBx()

	teg.globals.bx[100] = &box{kind: T_item, x_item: &entity_item{attack: 5}}
	teg.globals.bx[200] = &box{kind: T_item, x_item: &entity_item{defense: 3}}
	teg.globals.bx[300] = &box{kind: T_item, x_item: &entity_item{missile: 2}}
	teg.globals.bx[400] = &box{kind: T_item, x_item: &entity_item{}}

	if !is_fighter(100) {
		t.Error("is_fighter(100) should be true for item with attack")
	}
	if !is_fighter(200) {
		t.Error("is_fighter(200) should be true for item with defense")
	}
	if !is_fighter(300) {
		t.Error("is_fighter(300) should be true for item with missile")
	}
	if is_fighter(400) {
		t.Error("is_fighter(400) should be false for item with no combat stats")
	}
	if !is_fighter(item_ghost_warrior) {
		t.Error("is_fighter(item_ghost_warrior) should be true")
	}
}

func TestAliveAndIsPrisoner(t *testing.T) {
	defer clearBx()
	clearBx()

	teg.globals.bx[100] = &box{kind: T_char, x_char: &entity_char{prisoner: 0}}
	teg.globals.bx[200] = &box{kind: T_char, x_char: &entity_char{prisoner: 1}}
	teg.globals.bx[300] = &box{kind: T_loc}

	if !alive(100) {
		t.Error("alive(100) should be true for T_char")
	}
	if alive(300) {
		t.Error("alive(300) should be false for T_loc")
	}
	if is_prisoner(100) {
		t.Error("is_prisoner(100) should be false")
	}
	if !is_prisoner(200) {
		t.Error("is_prisoner(200) should be true")
	}
}

func TestLocAccessor(t *testing.T) {
	defer clearBx()
	clearBx()

	if loc(100) != 0 {
		t.Error("loc(100) should return 0 for nil box")
	}

	teg.globals.bx[100] = &box{
		kind:       T_char,
		x_loc_info: loc_info{where: 5000},
	}

	if loc(100) != 5000 {
		t.Errorf("loc(100) = %d, want 5000", loc(100))
	}
}

func TestPlayerAccessors(t *testing.T) {
	defer clearBx()
	clearBx()

	teg.globals.bx[100] = &box{
		kind:  T_player,
		skind: sub_pl_regular,
		x_player: &entity_player{
			noble_points: 10,
			fast_study:   5,
			format:       1,
			notab:        1,
		},
	}

	if player_np(100) != 10 {
		t.Errorf("player_np(100) = %d, want 10", player_np(100))
	}
	if player_fast_study(100) != 5 {
		t.Errorf("player_fast_study(100) = %d, want 5", player_fast_study(100))
	}
	if player_format(100) != 1 {
		t.Errorf("player_format(100) = %d, want 1", player_format(100))
	}
	if player_notab(100) != 1 {
		t.Errorf("player_notab(100) = %d, want 1", player_notab(100))
	}
}

func TestSkillAccessors(t *testing.T) {
	defer clearBx()
	clearBx()

	teg.globals.bx[100] = &box{
		kind:  T_skill,
		skind: sub_magic,
		x_skill: &entity_skill{
			time_to_learn:  14,
			required_skill: sk_basic,
			np_req:         2,
			produced:       item_gold,
			no_exp:         1,
		},
	}

	if learn_time(100) != 14 {
		t.Errorf("learn_time(100) = %d, want 14", learn_time(100))
	}
	if req_skill(100) != sk_basic {
		t.Errorf("req_skill(100) = %d, want sk_basic", req_skill(100))
	}
	if skill_np_req(100) != 2 {
		t.Errorf("skill_np_req(100) = %d, want 2", skill_np_req(100))
	}
	if skill_produce(100) != item_gold {
		t.Errorf("skill_produce(100) = %d, want item_gold", skill_produce(100))
	}
	if skill_no_exp(100) != 1 {
		t.Errorf("skill_no_exp(100) = %d, want 1", skill_no_exp(100))
	}
}

func TestMiscAccessors(t *testing.T) {
	defer clearBx()
	clearBx()

	teg.globals.bx[100] = &box{
		kind: T_storm,
		x_misc: &entity_misc{
			storm_str:   50,
			bind_storm:  200,
			summoned_by: 300,
			garr_castle: 400,
			npc_dir:     DIR_N,
			old_lord:    500,
			only_vuln:   600,
		},
	}

	if storm_strength(100) != 50 {
		t.Errorf("storm_strength(100) = %d, want 50", storm_strength(100))
	}
	if storm_bind(100) != 200 {
		t.Errorf("storm_bind(100) = %d, want 200", storm_bind(100))
	}
	if npc_summoner(100) != 300 {
		t.Errorf("npc_summoner(100) = %d, want 300", npc_summoner(100))
	}
	if garrison_castle(100) != 400 {
		t.Errorf("garrison_castle(100) = %d, want 400", garrison_castle(100))
	}
	if npc_last_dir(100) != DIR_N {
		t.Errorf("npc_last_dir(100) = %d, want DIR_N", npc_last_dir(100))
	}
	if body_old_lord(100) != 500 {
		t.Errorf("body_old_lord(100) = %d, want 500", body_old_lord(100))
	}
	if only_defeatable(100) != 600 {
		t.Errorf("only_defeatable(100) = %d, want 600", only_defeatable(100))
	}
}

func TestWaitTime(t *testing.T) {
	defer clearBx()
	clearBx()

	if wait_time(100) != 0 {
		t.Error("wait_time(100) should return 0 for nil box")
	}

	teg.globals.bx[100] = &box{
		kind: T_char,
		cmd:  &command{wait: 5},
	}

	if wait_time(100) != 5 {
		t.Errorf("wait_time(100) = %d, want 5", wait_time(100))
	}
}
