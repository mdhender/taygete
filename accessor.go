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

func kind(n int) schar {
	if n > 0 && n < MAX_BOXES && teg.globals.bx[n] != nil {
		return teg.globals.bx[n].kind
	}
	return T_deleted
}

func subkind(n int) schar {
	if teg.globals.bx[n] != nil {
		return teg.globals.bx[n].skind
	}
	return 0
}

func valid_box(n int) bool {
	return kind(n) != T_deleted
}

func kind_first(n int) int {
	return teg.globals.box_head[n]
}

func kind_next(n int) int {
	return teg.globals.bx[n].x_next_kind
}

func sub_first(n int) int {
	return teg.globals.sub_head[n]
}

func sub_next(n int) int {
	return teg.globals.bx[n].x_next_sub
}

func rp_loc_info(n int) *loc_info {
	if teg.globals.bx[n] == nil {
		return nil
	}
	return &teg.globals.bx[n].x_loc_info
}

func rp_char(n int) *entity_char {
	if teg.globals.bx[n] == nil {
		return nil
	}
	return teg.globals.bx[n].x_char
}

func rp_loc(n int) *entity_loc {
	if teg.globals.bx[n] == nil {
		return nil
	}
	return teg.globals.bx[n].x_loc
}

func rp_subloc(n int) *entity_subloc {
	if teg.globals.bx[n] == nil {
		return nil
	}
	return teg.globals.bx[n].x_subloc
}

func rp_item(n int) *entity_item {
	if teg.globals.bx[n] == nil {
		return nil
	}
	return teg.globals.bx[n].x_item
}

func rp_player(n int) *entity_player {
	if teg.globals.bx[n] == nil {
		return nil
	}
	return teg.globals.bx[n].x_player
}

func rp_skill(n int) *entity_skill {
	if teg.globals.bx[n] == nil {
		return nil
	}
	return teg.globals.bx[n].x_skill
}

func rp_gate(n int) *entity_gate {
	if teg.globals.bx[n] == nil {
		return nil
	}
	return teg.globals.bx[n].x_gate
}

func rp_misc(n int) *entity_misc {
	if teg.globals.bx[n] == nil {
		return nil
	}
	return teg.globals.bx[n].x_misc
}

func rp_disp(n int) *att_ent {
	if teg.globals.bx[n] == nil {
		return nil
	}
	return teg.globals.bx[n].x_disp
}

func rp_command(n int) *command {
	if teg.globals.bx[n] == nil {
		return nil
	}
	return teg.globals.bx[n].cmd
}

func rp_magic(n int) *char_magic {
	c := rp_char(n)
	if c == nil {
		return nil
	}
	return c.x_char_magic
}

func rp_item_magic(n int) *item_magic {
	it := rp_item(n)
	if it == nil {
		return nil
	}
	return it.x_item_magic
}

func p_loc_info(n int) *loc_info {
	if teg.globals.bx[n] == nil {
		teg.globals.bx[n] = &box{}
	}
	return &teg.globals.bx[n].x_loc_info
}

func p_char(n int) *entity_char {
	if teg.globals.bx[n] == nil {
		teg.globals.bx[n] = &box{}
	}
	if teg.globals.bx[n].x_char == nil {
		teg.globals.bx[n].x_char = &entity_char{}
	}
	return teg.globals.bx[n].x_char
}

func p_loc(n int) *entity_loc {
	if teg.globals.bx[n] == nil {
		teg.globals.bx[n] = &box{}
	}
	if teg.globals.bx[n].x_loc == nil {
		teg.globals.bx[n].x_loc = &entity_loc{}
	}
	return teg.globals.bx[n].x_loc
}

func p_subloc(n int) *entity_subloc {
	if teg.globals.bx[n] == nil {
		teg.globals.bx[n] = &box{}
	}
	if teg.globals.bx[n].x_subloc == nil {
		teg.globals.bx[n].x_subloc = &entity_subloc{}
	}
	return teg.globals.bx[n].x_subloc
}

func p_item(n int) *entity_item {
	if teg.globals.bx[n] == nil {
		teg.globals.bx[n] = &box{}
	}
	if teg.globals.bx[n].x_item == nil {
		teg.globals.bx[n].x_item = &entity_item{}
	}
	return teg.globals.bx[n].x_item
}

func p_player(n int) *entity_player {
	if teg.globals.bx[n] == nil {
		teg.globals.bx[n] = &box{}
	}
	if teg.globals.bx[n].x_player == nil {
		teg.globals.bx[n].x_player = &entity_player{}
	}
	return teg.globals.bx[n].x_player
}

func p_skill(n int) *entity_skill {
	if teg.globals.bx[n] == nil {
		teg.globals.bx[n] = &box{}
	}
	if teg.globals.bx[n].x_skill == nil {
		teg.globals.bx[n].x_skill = &entity_skill{}
	}
	return teg.globals.bx[n].x_skill
}

func p_gate(n int) *entity_gate {
	if teg.globals.bx[n] == nil {
		teg.globals.bx[n] = &box{}
	}
	if teg.globals.bx[n].x_gate == nil {
		teg.globals.bx[n].x_gate = &entity_gate{}
	}
	return teg.globals.bx[n].x_gate
}

func p_misc(n int) *entity_misc {
	if teg.globals.bx[n] == nil {
		teg.globals.bx[n] = &box{}
	}
	if teg.globals.bx[n].x_misc == nil {
		teg.globals.bx[n].x_misc = &entity_misc{}
	}
	return teg.globals.bx[n].x_misc
}

func p_disp(n int) *att_ent {
	if teg.globals.bx[n] == nil {
		teg.globals.bx[n] = &box{}
	}
	if teg.globals.bx[n].x_disp == nil {
		teg.globals.bx[n].x_disp = &att_ent{}
	}
	return teg.globals.bx[n].x_disp
}

func p_command(n int) *command {
	if teg.globals.bx[n] == nil {
		teg.globals.bx[n] = &box{}
	}
	if teg.globals.bx[n].cmd == nil {
		teg.globals.bx[n].cmd = &command{}
	}
	return teg.globals.bx[n].cmd
}

func p_magic(n int) *char_magic {
	c := p_char(n)
	if c.x_char_magic == nil {
		c.x_char_magic = &char_magic{}
	}
	return c.x_char_magic
}

func p_item_magic(n int) *item_magic {
	it := p_item(n)
	if it.x_item_magic == nil {
		it.x_item_magic = &item_magic{}
	}
	return it.x_item_magic
}

func loc(n int) int {
	li := rp_loc_info(n)
	if li == nil {
		return 0
	}
	return li.where
}

func is_loc_or_ship(n int) bool {
	k := kind(n)
	return k == T_loc || k == T_ship
}

func is_ship(n int) bool {
	sk := subkind(n)
	return sk == sub_galley || sk == sub_roundship || sk == sub_raft
}

func is_ship_notdone(n int) bool {
	sk := subkind(n)
	return sk == sub_galley_notdone || sk == sub_roundship_notdone || sk == sub_raft_notdone
}

func is_ship_either(n int) bool {
	return is_ship(n) || is_ship_notdone(n)
}

func char_health(n int) schar {
	c := rp_char(n)
	if c == nil {
		return 0
	}
	return c.health
}

func char_sick(n int) schar {
	c := rp_char(n)
	if c == nil {
		return 0
	}
	return c.sick
}

func char_guard(n int) schar {
	c := rp_char(n)
	if c == nil {
		return 0
	}
	return c.guard
}

func loyal_kind(n int) schar {
	c := rp_char(n)
	if c == nil {
		return 0
	}
	return c.loy_kind
}

func loyal_rate(n int) int {
	c := rp_char(n)
	if c == nil {
		return 0
	}
	return c.loy_rate
}

func noble_item(n int) schar {
	c := rp_char(n)
	if c == nil {
		return 0
	}
	return c.unit_item
}

func char_behind(n int) schar {
	c := rp_char(n)
	if c == nil {
		return 0
	}
	return c.behind
}

func npc_program(n int) schar {
	c := rp_char(n)
	if c == nil {
		return 0
	}
	return c.npc_prog
}

func char_studied(n int) schar {
	c := rp_char(n)
	if c == nil {
		return 0
	}
	return c.studied
}

func char_moving(n int) int {
	c := rp_char(n)
	if c == nil {
		return 0
	}
	return c.moving
}

func char_attack(n int) short {
	c := rp_char(n)
	if c == nil {
		return 0
	}
	return c.attack
}

func char_defense(n int) short {
	c := rp_char(n)
	if c == nil {
		return 0
	}
	return c.defense
}

func char_missile(n int) short {
	c := rp_char(n)
	if c == nil {
		return 0
	}
	return c.missile
}

func char_rank(n int) schar {
	c := rp_char(n)
	if c == nil {
		return 0
	}
	return c.rank
}

func char_break(n int) schar {
	c := rp_char(n)
	if c == nil {
		return 0
	}
	return c.break_point
}

func char_pray(n int) schar {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.pray
}

func char_hidden(n int) schar {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.hide_self
}

func vision_protect(n int) schar {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.vis_protect
}

func default_garrison(n int) schar {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.default_garr
}

func char_hide_mage(n int) schar {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.hide_mage
}

func char_cur_aura(n int) int {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.cur_aura
}

func char_max_aura(n int) int {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.max_aura
}

func reflect_blast(n int) schar {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.aura_reflect
}

func char_pledge(n int) int {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.pledge
}

func char_auraculum(n int) int {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.auraculum
}

func char_abil_shroud(n int) short {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.ability_shroud
}

func board_fee(n int) int {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.fee
}

func ferry_horn(n int) schar {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.ferry_flag
}

func char_proj_cast(n int) int {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.project_cast
}

func char_quick_cast(n int) short {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.quick_cast
}

func is_magician(n int) schar {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.magician
}

func weather_mage(n int) schar {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.knows_weather
}

func loc_barrier(n int) short {
	l := rp_loc(n)
	if l == nil {
		return 0
	}
	return l.barrier
}

func loc_shroud(n int) short {
	l := rp_loc(n)
	if l == nil {
		return 0
	}
	return l.shroud
}

func loc_civ(n int) schar {
	l := rp_loc(n)
	if l == nil {
		return 0
	}
	return l.civ
}

func loc_sea_lane(n int) schar {
	l := rp_loc(n)
	if l == nil {
		return 0
	}
	return l.sea_lane
}

func gate_dist(n int) schar {
	l := rp_loc(n)
	if l == nil {
		return 0
	}
	return l.dist_from_gate
}

func loc_prominence(n int) schar {
	sl := rp_subloc(n)
	if sl == nil {
		return 0
	}
	return sl.prominence
}

func loc_opium(n int) int {
	sl := rp_subloc(n)
	if sl == nil {
		return 0
	}
	return sl.opium_econ
}

func ship_cap_raw(n int) int {
	sl := rp_subloc(n)
	if sl == nil {
		return 0
	}
	return sl.capacity
}

func ship_moving(n int) int {
	sl := rp_subloc(n)
	if sl == nil {
		return 0
	}
	return sl.moving
}

func loc_damage(n int) uchar {
	sl := rp_subloc(n)
	if sl == nil {
		return 0
	}
	return sl.damage
}

func loc_defense(n int) int {
	sl := rp_subloc(n)
	if sl == nil {
		return 0
	}
	return sl.defense
}

func loc_pillage(n int) schar {
	sl := rp_subloc(n)
	if sl == nil {
		return 0
	}
	return sl.loot
}

func recent_pillage(n int) schar {
	sl := rp_subloc(n)
	if sl == nil {
		return 0
	}
	return sl.recent_loot
}

func safe_haven(n int) schar {
	sl := rp_subloc(n)
	if sl == nil {
		return 0
	}
	return sl.safe
}

func major_city(n int) schar {
	sl := rp_subloc(n)
	if sl == nil {
		return 0
	}
	return sl.major
}

func uldim(n int) schar {
	sl := rp_subloc(n)
	if sl == nil {
		return 0
	}
	return sl.uldim_flag
}

func summerbridge(n int) schar {
	sl := rp_subloc(n)
	if sl == nil {
		return 0
	}
	return sl.summer_flag
}

func subloc_quest(n int) schar {
	sl := rp_subloc(n)
	if sl == nil {
		return 0
	}
	return sl.quest_late
}

func tunnel_depth(n int) schar {
	sl := rp_subloc(n)
	if sl == nil {
		return 0
	}
	return sl.tunnel_level
}

func mine_depth(n int) short {
	sl := rp_subloc(n)
	if sl == nil {
		return 0
	}
	return sl.shaft_depth / 3
}

func castle_level(n int) schar {
	sl := rp_subloc(n)
	if sl == nil {
		return 0
	}
	return sl.castle_lev
}

func ship_has_ram(n int) schar {
	sl := rp_subloc(n)
	if sl == nil {
		return 0
	}
	return sl.galley_ram
}

func loc_link_open(n int) schar {
	sl := rp_subloc(n)
	if sl == nil {
		return 0
	}
	return sl.link_open
}

func road_dest(n int) int {
	g := rp_gate(n)
	if g == nil {
		return 0
	}
	return g.to_loc
}

func road_hidden(n int) schar {
	g := rp_gate(n)
	if g == nil {
		return 0
	}
	return g.road_hidden
}

func gate_seal(n int) short {
	g := rp_gate(n)
	if g == nil {
		return 0
	}
	return g.seal_key
}

func gate_dest(n int) int {
	g := rp_gate(n)
	if g == nil {
		return 0
	}
	return g.to_loc
}

func item_animal(n int) schar {
	it := rp_item(n)
	if it == nil {
		return 0
	}
	return it.animal
}

func item_prominent(n int) schar {
	it := rp_item(n)
	if it == nil {
		return 0
	}
	return it.prominent
}

func item_weight(n int) short {
	it := rp_item(n)
	if it == nil {
		return 0
	}
	return it.weight
}

func item_price(n int) int {
	it := rp_item(n)
	if it == nil {
		return 0
	}
	return it.base_price
}

func item_unique(n int) int {
	it := rp_item(n)
	if it == nil {
		return 0
	}
	return it.who_has
}

func item_land_cap(n int) short {
	it := rp_item(n)
	if it == nil {
		return 0
	}
	return it.land_cap
}

func item_ride_cap(n int) short {
	it := rp_item(n)
	if it == nil {
		return 0
	}
	return it.ride_cap
}

func item_fly_cap(n int) short {
	it := rp_item(n)
	if it == nil {
		return 0
	}
	return it.fly_cap
}

func item_attack(n int) short {
	it := rp_item(n)
	if it == nil {
		return 0
	}
	return it.attack
}

func item_defense(n int) short {
	it := rp_item(n)
	if it == nil {
		return 0
	}
	return it.defense
}

func item_missile(n int) short {
	it := rp_item(n)
	if it == nil {
		return 0
	}
	return it.missile
}

func man_item(n int) schar {
	it := rp_item(n)
	if it == nil {
		return 0
	}
	return it.is_man_item
}

func item_capturable(n int) schar {
	it := rp_item(n)
	if it == nil {
		return 0
	}
	return it.capturable
}

func player_np(n int) short {
	p := rp_player(n)
	if p == nil {
		return 0
	}
	return p.noble_points
}

func player_fast_study(n int) short {
	p := rp_player(n)
	if p == nil {
		return 0
	}
	return p.fast_study
}

func player_email(n int) *char {
	p := rp_player(n)
	if p == nil {
		return nil
	}
	return p.email
}

func times_paid(n int) char {
	p := rp_player(n)
	if p == nil {
		return 0
	}
	return p.times_paid
}

func player_public_turn(n int) char {
	p := rp_player(n)
	if p == nil {
		return 0
	}
	return p.public_turn
}

func player_format(n int) schar {
	p := rp_player(n)
	if p == nil {
		return 0
	}
	return p.format
}

func player_notab(n int) schar {
	p := rp_player(n)
	if p == nil {
		return 0
	}
	return p.notab
}

func req_skill(n int) int {
	sk := rp_skill(n)
	if sk == nil {
		return 0
	}
	return sk.required_skill
}

func skill_produce(n int) int {
	sk := rp_skill(n)
	if sk == nil {
		return 0
	}
	return sk.produced
}

func skill_no_exp(n int) int {
	sk := rp_skill(n)
	if sk == nil {
		return 0
	}
	return sk.no_exp
}

func skill_np_req(n int) int {
	sk := rp_skill(n)
	if sk == nil {
		return 0
	}
	return sk.np_req
}

func learn_time(n int) int {
	sk := rp_skill(n)
	if sk == nil {
		return 0
	}
	return sk.time_to_learn
}

func banner(n int) *char {
	m := rp_misc(n)
	if m == nil {
		return nil
	}
	return m.display
}

func storm_bind(n int) int {
	m := rp_misc(n)
	if m == nil {
		return 0
	}
	return m.bind_storm
}

func storm_strength(n int) short {
	m := rp_misc(n)
	if m == nil {
		return 0
	}
	return m.storm_str
}

func npc_summoner(n int) int {
	m := rp_misc(n)
	if m == nil {
		return 0
	}
	return m.summoned_by
}

func garrison_castle(n int) int {
	m := rp_misc(n)
	if m == nil {
		return 0
	}
	return m.garr_castle
}

func npc_last_dir(n int) schar {
	m := rp_misc(n)
	if m == nil {
		return 0
	}
	return m.npc_dir
}

func restricted_control(n int) char {
	m := rp_misc(n)
	if m == nil {
		return 0
	}
	return m.cmd_allow
}

func body_old_lord(n int) int {
	m := rp_misc(n)
	if m == nil {
		return 0
	}
	return m.old_lord
}

func only_defeatable(n int) int {
	m := rp_misc(n)
	if m == nil {
		return 0
	}
	return m.only_vuln
}

func item_lore(n int) int {
	im := rp_item_magic(n)
	if im == nil {
		return 0
	}
	return im.lore
}

func item_use_key(n int) schar {
	im := rp_item_magic(n)
	if im == nil {
		return 0
	}
	return im.use_key
}

func item_creator(n int) int {
	im := rp_item_magic(n)
	if im == nil {
		return 0
	}
	return im.creator
}

func item_aura(n int) short {
	im := rp_item_magic(n)
	if im == nil {
		return 0
	}
	return im.aura
}

func item_creat_cloak(n int) schar {
	im := rp_item_magic(n)
	if im == nil {
		return 0
	}
	return im.cloak_creator
}

func item_reg_cloak(n int) schar {
	im := rp_item_magic(n)
	if im == nil {
		return 0
	}
	return im.cloak_region
}

func item_creat_loc(n int) int {
	im := rp_item_magic(n)
	if im == nil {
		return 0
	}
	return im.region_created
}

func item_curse_non(n int) schar {
	im := rp_item_magic(n)
	if im == nil {
		return 0
	}
	return im.curse_loyalty
}

func item_attack_bonus(n int) schar {
	im := rp_item_magic(n)
	if im == nil {
		return 0
	}
	return im.attack_bonus
}

func item_defense_bonus(n int) schar {
	im := rp_item_magic(n)
	if im == nil {
		return 0
	}
	return im.defense_bonus
}

func item_missile_bonus(n int) schar {
	im := rp_item_magic(n)
	if im == nil {
		return 0
	}
	return im.missile_bonus
}

func item_aura_bonus(n int) short {
	im := rp_item_magic(n)
	if im == nil {
		return 0
	}
	return im.aura_bonus
}

func item_relic_decay(n int) short {
	im := rp_item_magic(n)
	if im == nil {
		return 0
	}
	return im.relic_decay
}

func item_token_num(n int) schar {
	im := rp_item_magic(n)
	if im == nil {
		return 0
	}
	return im.token_num
}

func item_token_ni(n int) int {
	im := rp_item_magic(n)
	if im == nil {
		return 0
	}
	return im.token_ni
}

func release_swear(n int) schar {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.swear_on_release
}

func our_token(n int) int {
	m := rp_magic(n)
	if m == nil {
		return 0
	}
	return m.token
}

func is_fighter(n int) bool {
	return item_attack(n) != 0 || item_defense(n) != 0 || item_missile(n) != 0 || n == item_ghost_warrior
}

func alive(n int) bool {
	return kind(n) == T_char
}

func wait_time(n int) int {
	c := rp_command(n)
	if c == nil {
		return 0
	}
	return c.wait
}

func is_prisoner(n int) bool {
	c := rp_char(n)
	if c == nil {
		return false
	}
	return c.prisoner != 0
}

// player returns the owning player for a character.
// Walks up the unit_lord chain until it finds a player.
// Ported from src/u.c.
func player(n int) int {
	count := 0
	for n > 0 && kind(n) != T_player {
		c := rp_char(n)
		if c == nil {
			return 0
		}
		n = c.unit_lord
		count++
		if count >= 1000 {
			panic("player: infinite loop detected")
		}
	}
	return n
}
