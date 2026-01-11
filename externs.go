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

// extern char *int_to_code(int l);
// extern int code_to_int(char *s);
// extern int scode(char *s);
//
// extern char *name(int n);
// extern void set_name(int n, char *s);
// extern void set_banner(int n, char *s);
//
// extern char *display_name(int n);
// extern char *display_kind(int n);
// extern char *box_code(int n);
// extern char *box_code_less(int n);
// extern char *box_name(int n);
// extern char *just_name(int n);
// extern char *just_name_qty(int item, int qty);
// extern char *plural_item_name(int item, int qty);
// extern char *plural_item_box(int item, int qty);
// extern char *box_name_qty(int item, int qty);
// extern char *box_name_kind(int n);

// extern void list_exits(int who, int where);
// extern struct exit_view **exits_from_loc(int who, int where);
// extern void determine_map_edges();
// extern void dir_assert();
// extern int exit_distance(int, int);
// extern void find_hidden_exit(int who, struct exit_view **l, int which);
// extern int count_hidden_exits(struct exit_view **l);
// extern void list_sailable_routes(int who, int ship);
// extern int hidden_count_to_index(int which, struct exit_view **l);
// extern int has_ocean_access(int where);
// extern int location_direction(int where, int dir);
//
// /*
//  *  From move.c:
//  */
//
// extern struct exit_view *
//     parse_exit_dir(struct command *c, int where, char *zero_arg);
//
//
// #define    DIR_NSEW(a)        ((a) >= DIR_N && (a) <= DIR_W)

// extern struct exit_view **exits_from_loc_nsew(int, int);
// extern struct exit_view **exits_from_loc_nsew_select(int, int, int, int);
//
// extern char *liner_desc(int n);
// extern char *display_with(int who);
// extern char *display_owner(int who);
// extern char *show_loc_header(int where);
// extern void show_chars_below(int who, int n);
// extern void show_owner_stack(int who, int n);
//
// extern int show_display_string;
// extern void turn_end_loc_reports();
// extern int any_chars_here(int where);
// extern char *loc_civ_s(int where);
//
// extern int max_eff_aura(int who);                /* art.c */
// extern int has_auraculum(int who);                /* art.c */
// extern void touch_loc_pl(int pl, int where);            /* day.c */
// extern void touch_loc(int who);                    /* day.c */
// extern struct trade *add_city_trade(int, int, int, int, int, int);    /* buy.c */
// extern int distance(int orig, int dest, int gate);        /* seed.c */
// extern void interrupt_order(int who);                /* input.c */
// extern int may_cookie_npc(int who, int where, int cookie);    /* npc.c */
// extern int do_cookie_npc(int, int, int, int);            /* npc.c */
// extern int weather_here(int where, int sk);            /* storm.c */
// extern void show_char_inventory(int who, int num);        /* report.c */
// extern char **parse_line(char **l, char *s);            /* input.c */
// extern int unit_maint_cost(int who);                /* day.c */
// extern char *wield_s(int who);                    /* combat.c */
// extern int can_see_weather_here(int who, int where);        /* storm.c */
// extern void queue_lore(int who, int num, int anyway);        /* lore.c */
// extern char *read_pw(char *type);                /* pw.c */

// extern int top_ruler(int n);
// extern void garrison_gold();
// extern char *rank_s(int who);
// extern void touch_garrison_locs();
// extern void determine_noble_ranks();
// extern int may_rule_here(int who, int where);
// extern ilist players_who_rule_here(int where);
// extern int garrison_notices(int garr, int target);
//
// extern int loc_owner(int where);
// extern void mark_loc_stack_known(int stack, int where);
//
// extern void all_here(int who, ilist *l);
// extern void all_char_here(int who, ilist *l);
//
// extern int region(int who);
// extern int province(int who);
// extern int subloc(int who);
// extern int viewloc(int who);
//
// extern void add_to_here_list(int loc, int who);
// extern void remove_from_here_list(int loc, int who);
// extern void set_where(int who, int new_loc);
// extern int in_here_list(int loc, int who);
// extern int somewhere_inside(int a, int b);
// extern int in_safe_now(int who);
// extern int subloc_here(int where, int sk);
// extern int count_loc_structures(int where, int a, int b);
//
// #define    city_here(a)    subloc_here((a), sub_city)
//
//
// /*
//  *  loop.h -- abstracted loops
//  *
//  *  or, "Too bad this language doesn't have generators"
//  *
//  *
//  *  Q:  Why abstract loops?  These defines are really gross.
//  *
//  *  A:  We use abstract data types so that we can change the representation
//  *      without having to change all of the code that uses the types.
//  *
//  *      It is easy to abstract add, delete and fetch operations.  But the
//  *      most common operation is to iterate over the elements of a collection.
//  *
//  *      The method of iteration is almost certain to change when switching
//  *      implementations, too.  List-to-tree, tree-to-bit-array, etc.
//  *
//  *      Abstracting loops makes the code cleaner and easier to read, and
//  *      easier to change.
//  */
//
//
// /*
//  *  Loops below should be free of the "delete problem", i.e.
//  *
//  *    loop_something(i)
//  *    {
//  *        delete_something(i);
//  *    }
//  *    next_something;
//  *
//  *  should work, i.e. it shouldn't core dump because next(i) is no
//  *  longer defined, or go over an element twice, or miss an element.
//  */
//
// /*
//  *  break and continue should work inside these loops, but don't
//  *  return out of them.  Return will bypass the end-of-loop cleanup.
//  */
//
// /*
//  *  Thanks to the X Window system for going ahead of us and making sure
//  *  that most vendor's compilers can handle unreasonably large defines.
//  */
//
//
// #define    loop_all_here(where, i) \
// { int ll_i; \
//   int ll_check = 2; \
//   ilist ll_l = NULL; \
//     all_here(where, &ll_l); \
//     for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) \
//     { \
//         i = ll_l[ll_i];
//
// #define    next_all_here        } assert(ll_check == 2); ilist_reclaim(&ll_l); }
//
//
// #define    loop_char_here(where, i) \
// { int ll_i; \
//   int ll_check = 13; \
//   ilist ll_l = NULL; \
//     all_char_here(where, &ll_l); \
//     for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) \
//     { \
//         i = ll_l[ll_i];
//
// #define    next_char_here    } assert(ll_check == 13); ilist_reclaim(&ll_l); }
//
//
//
// #define    loop_stack(who, i) \
// { int ll_i; \
//   int ll_check = 20; \
//   ilist ll_l = NULL; \
//     all_stack(who, &ll_l); \
//     for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) \
//     { \
//         i = ll_l[ll_i];
//
// #define    next_stack    } assert(ll_check == 20); ilist_reclaim(&ll_l); }
//
//
// #define    loop_known(kn, i) \
// { int ll_i; \
//   int ll_check = 3; \
//   extern int int_comp(); \
//     qsort(kn, ilist_len(kn), sizeof(int), int_comp); \
//     for (ll_i = 0; ll_i < ilist_len(kn); ll_i++) { \
//         i = kn[ll_i];
//
// #define    next_known    } assert(ll_check == 3); }
//
//
// /*
//  *  Iterate over all valid boxes.  i is instantiated with the entity
//  *  numbers.
//  */
//
// #define    loop_boxes(i) \
// { int ll_i; \
//   int ll_check = 4; \
//     for (ll_i = 1; ll_i < MAX_BOXES; ll_i++) \
//         if (kind(ll_i) != T_deleted) \
//         { \
//             i = ll_i;
//
// #define    next_box    } assert(ll_check == 4); }
//
//
// #define    loop_kind(kind, i) \
// { int ll_i, ll_next; \
//   int ll_check = 5; \
//     ll_next = kind_first(kind); \
//     while ((ll_i = ll_next) > 0) { \
//         ll_next = kind_next(ll_i); \
//         i = ll_i;
//
// #define    next_kind    } assert(ll_check == 5); }
//
//
// #define    loop_char(i)    loop_kind(T_char, i)
// #define    next_char    next_kind
//
// #define    loop_player(i)    loop_kind(T_player, i)
// #define    next_player    next_kind
//
// #define    loop_loc(i)    loop_kind(T_loc, i)
// #define    next_loc    next_kind
//
// #define    loop_item(i)    loop_kind(T_item, i)
// #define    next_item    next_kind
//
// #define    loop_exit(i)    loop_kind(T_exit, i)
// #define    next_exit    next_kind
//
// #define    loop_skill(i)    loop_kind(T_skill, i)
// #define    next_skill    next_kind
//
// #define    loop_gate(i)    loop_kind(T_gate, i)
// #define    next_gate    next_kind
//
// #define    loop_ship(i)    loop_kind(T_ship, i)
// #define    next_ship    next_kind
//
// #define    loop_post(i)    loop_kind(T_post, i)
// #define    next_post    next_kind
//
// #define    loop_storm(i)    loop_kind(T_storm, i)
// #define    next_storm    next_kind
//
//
// #define    loop_subkind(sk, i) \
// { int ll_i, ll_next; \
//   int ll_check = 26; \
//     ll_next = sub_first(sk); \
//     while ((ll_i = ll_next) > 0) { \
//         ll_next = sub_next(ll_i); \
//         i = ll_i;
//
// #define    next_subkind    } assert(ll_check == 26); }
//
// #define    loop_garrison(i)    loop_subkind(sub_garrison, i)
// #define    next_garrison        next_subkind
//
// #define    loop_city(i)        loop_subkind(sub_city, i)
// #define    next_city        next_subkind
//
// #define    loop_mountain(i)    loop_subkind(sub_mountain, i)
// #define    next_mountain        next_subkind
//
// #define    loop_inn(i)        loop_subkind(sub_inn, i)
// #define    next_inn        next_subkind
//
// #define    loop_temple(i)        loop_subkind(sub_temple, i)
// #define    next_temple        next_subkind
//
// #define    loop_collapsed_mine(i)    loop_subkind(sub_mine_collapsed, i)
// #define    next_collapsed_mine    next_subkind
//
// #define    loop_dead_body(i)    loop_subkind(sub_dead_body, i)
// #define    next_dead_body        next_subkind
//
// #define    loop_castle(i)        loop_subkind(sub_castle, i)
// #define    next_castle        next_subkind
//
// #define    loop_pl_regular(i)    loop_subkind(sub_pl_regular, i)
// #define    next_pl_regular        next_subkind
//
//
// #define    loop_loc_or_ship(i) \
// { int ll_i; \
//   int ll_check = 17; \
//   int ll_state = 1; \
//     ll_i = kind_first(T_ship); \
//     if (ll_i <= 0) { ll_i = kind_first(T_loc); ll_state = 0; } \
//     while (ll_i > 0) { \
//         i = ll_i;
//
// #define    next_loc_or_ship \
//     ll_i = kind_next(ll_i); \
//     if (ll_i <= 0 && ll_state) \
//      { ll_i = kind_first(T_loc); ll_state = 0; } \
//     } assert(ll_check == 17); }
//
//
// #define    loop_province(i) \
// { int ll_i; \
//   int ll_check = 6; \
//     for (ll_i = kind_first(T_loc); ll_i > 0; ll_i = kind_next(ll_i)) \
//         if (loc_depth(ll_i) == LOC_province) { \
//             i = ll_i;
//
// #define    next_province    } assert(ll_check == 6); }
//
//
// #define    loop_loc_teach(where, i) \
// { int ll_i; \
//   int ll_check = 7; \
//   ilist ll_l = NULL; \
//     assert(valid_box(where)); \
//     if (rp_subloc(where)) \
//         ll_l = ilist_copy(rp_subloc(where)->teaches); \
//         for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) { \
//             i = ll_l[ll_i];
//
// #define    next_loc_teach    } assert(ll_check == 7); ilist_reclaim(&ll_l); }
//
//
// #define    loop_units(pl, i) \
// { int ll_i; \
//   ilist ll_l = NULL; \
//   int ll_check = 21; \
//     if (rp_player(pl)) \
//         ll_l = ilist_copy(rp_player(pl)->units); \
//         for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) { \
//             i = ll_l[ll_i];
//
// #define    next_unit    } assert(ll_check == 21); ilist_reclaim(&ll_l); }
//
//
// #define    loop_here(where, i) \
// { int ll_i; \
//   ilist ll_l = NULL; \
//   int ll_check = 8; \
//     assert(valid_box(where)); \
//     if (rp_loc_info(where)) \
//         ll_l = ilist_copy(rp_loc_info(where)->here_list); \
//         for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) { \
//             i = ll_l[ll_i];
//
// #define    next_here    } assert(ll_check == 8); ilist_reclaim(&ll_l); }
//
//
// #define    loop_gates_here(where, i) \
// { int ll_i; \
//   ilist ll_l = NULL; \
//   int ll_check = 18; \
//     assert(valid_box(where)); \
//     if (rp_loc_info(where)) \
//         ll_l = ilist_copy(rp_loc_info(where)->here_list); \
//         for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) \
//             if (kind(ll_l[ll_i]) == T_gate) { \
//                 i = ll_l[ll_i];
//
// #define    next_gate_here    } assert(ll_check == 18); ilist_reclaim(&ll_l); }
//
//
// #define    loop_exits_here(where, i) \
// { int ll_i; \
//   int ll_check = 10; \
//   ilist ll_l = NULL; \
//     assert(valid_box(where)); \
//     if (rp_loc(where)) \
//         ll_l = ilist_copy(rp_loc(where)->exits_here); \
//         for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) { \
//             i = ll_l[ll_i];
//
// #define    next_exit_here    } assert(ll_check == 10); ilist_reclaim(&ll_l); }
//
//
// /*
//  *  Iterate struct item_ent *e over who's inventory
//  */
//
// #define    loop_inv(who, e) \
// { int ll_i; \
//   int ll_check = 11; \
//   struct item_ent ll_copy; \
//   struct item_ent **ll_l = NULL; \
//     assert(valid_box(who)); \
//     ll_l = (struct item_ent **) plist_copy((plist) bx[who]->items); \
//     for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) \
//         if (valid_box(ll_l[ll_i]->item) && ll_l[ll_i]->qty > 0) { \
//             ll_copy = *ll_l[ll_i]; \
//             e = &ll_copy;
//
// #define    next_inv   } assert(ll_check == 11); ilist_reclaim((ilist *) &ll_l); }
//
//
//
// #define    loop_char_skill(who, e) \
// { int ll_i; \
//   int ll_check = 15; \
//   struct skill_ent **ll_l = NULL; \
//     assert(valid_box(who)); \
//     if (rp_char(who)) \
//        ll_l = (struct skill_ent **) plist_copy((plist) rp_char(who)->skills); \
//        for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) { \
//          e = ll_l[ll_i];
//
// #define    next_char_skill \
//         } assert(ll_check == 15); ilist_reclaim((ilist *) &ll_l); }
//
//
//
// #define    loop_char_skill_known(who, e) \
// { int ll_i; \
//   int ll_check = 16; \
//   struct skill_ent **ll_l = NULL; \
//     assert(valid_box(who)); \
//     if (rp_char(who)) \
//        ll_l = (struct skill_ent **) plist_copy((plist) rp_char(who)->skills); \
//        for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) \
//           if (ll_l[ll_i]->know == SKILL_know) { \
//         e = ll_l[ll_i];
//
// #define    next_char_skill_known \
//         } assert(ll_check == 16); ilist_reclaim((ilist *) &ll_l); }
//
//
// #define    loop_trade(who, e) \
// { int ll_i; \
//   int ll_check = 19; \
//   struct trade **ll_l = NULL; \
//     assert(valid_box(who)); \
//     ll_l = (struct trade **) plist_copy((plist) bx[who]->trades); \
//     for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) \
//         if (valid_box(ll_l[ll_i]->item) && ll_l[ll_i]->qty > 0) { \
//             e = ll_l[ll_i];
//
// #define    next_trade  } assert(ll_check == 19); ilist_reclaim((ilist *) &ll_l); }
//
//
// #define    loop_prov_dest(where, i) \
// { int ll_i; \
//   int ll_check = 23; \
//   struct entity_loc *ll_p; \
//     assert(loc_depth(where) == LOC_province); \
//     ll_p = rp_loc(where); \
//     if (ll_p) { \
//     for (ll_i = 0; ll_i < ilist_len(ll_p->prov_dest); ll_i++) \
//     { \
//         i = ll_p->prov_dest[ll_i];
//
// #define    next_prov_dest        } } assert(ll_check == 23); }
//
//
// #define    loop_loc_owner(where, i) \
// { int ll_i, ll_next; \
//   int ll_check = 25; \
//     ll_next = province_admin(where); \
//     while ((ll_i = ll_next) > 0) { \
//         ll_next = char_pledge(ll_i); \
//         i = ll_i;
//
// #define    next_loc_owner    } assert(ll_check == 25); }

// #define    oly_month(a)    (((a).turn-1) % NUM_MONTHS)
// #define oly_year(a)    (((a).turn-1) / NUM_MONTHS)

// #define    numargs(c)    (ilist_len(c->parse) - 1)
//
// /*
//  *  How long a command has been running
//  */
//
// #define    command_days(c)        (c->days_executing)

// #define    if_malloc(p)        ((p) ? (p) : ((p) = my_malloc(sizeof(*(p)))))
//
// /*
//  *  malloc-on-demand substructure references
//  */
//
// #define    p_loc_info(n)        (&bx[n]->x_loc_info)
// #define    p_char(n)        if_malloc(bx[n]->x_char)
// #define    p_loc(n)        if_malloc(bx[n]->x_loc)
// #define    p_subloc(n)        if_malloc(bx[n]->x_subloc)
// #define    p_item(n)        if_malloc(bx[n]->x_item)
// #define    p_player(n)        if_malloc(bx[n]->x_player)
// #define    p_skill(n)        if_malloc(bx[n]->x_skill)
// #define    p_gate(n)        if_malloc(bx[n]->x_gate)
// #define    p_misc(n)        if_malloc(bx[n]->x_misc)
// #define    p_disp(n)        if_malloc(bx[n]->x_disp)
// #define    p_command(n)        if_malloc(bx[n]->cmd)
//
//
// /*
//  *  "raw" pointers to substructures, may be NULL
//  */
//
// #define    rp_loc_info(n)        (&bx[n]->x_loc_info)
// #define    rp_char(n)        (bx[n]->x_char)
// #define    rp_loc(n)        (bx[n]->x_loc)
// #define    rp_subloc(n)        (bx[n]->x_subloc)
// #define    rp_item(n)        (bx[n]->x_item)
// #define    rp_player(n)        (bx[n]->x_player)
// #define    rp_skill(n)        (bx[n]->x_skill)
// #define    rp_gate(n)        (bx[n]->x_gate)
// #define    rp_misc(n)        (bx[n]->x_misc)
// #define    rp_disp(n)        (bx[n]->x_disp)
// #define    rp_command(n)        (bx[n]->cmd)
//
//
// #define    rp_magic(n)        (rp_char(n) ? rp_char(n)->x_char_magic : NULL)
// #define    p_magic(n)        if_malloc(p_char(n)->x_char_magic)
//
// #define    rp_item_magic(n)    (rp_item(n) ? rp_item(n)->x_item_magic : NULL)
// #define    p_item_magic(n)        if_malloc(p_item(n)->x_item_magic)
//
// extern struct box **bx;
// extern int box_head[];            /* head of x_next_kind chain */
// extern int sub_head[];            /* head of x_next_sub chain */
//
//
// #define    kind(n)        (((n) > 0 && (n) < MAX_BOXES && bx[n]) ? bx[n]->kind : T_deleted)
//
// #define    subkind(n)    (bx[n] ? bx[n]->skind : 0)
// #define    valid_box(n)    (kind(n) != T_deleted)
//
// #define    kind_first(n)    (box_head[(n)])
// #define    kind_next(n)    (bx[(n)]->x_next_kind)
//
// #define    sub_first(n)    (sub_head[(n)])
// #define    sub_next(n)    (bx[(n)]->x_next_sub)
//
// #define    is_loc_or_ship(n)    (kind(n) == T_loc || kind(n) == T_ship)
// #define    is_ship(n)    (subkind(n) == sub_galley || subkind(n) == sub_roundship || subkind(n) == sub_raft)
// #define    is_ship_notdone(n)    (subkind(n) == sub_galley_notdone || subkind(n) == sub_roundship_notdone || subkind(n) == sub_raft_notdone)
// #define    is_ship_either(n)    (is_ship(n) || is_ship_notdone(n))
//
//
// /*
//  *    _moving indicates that the unit has initiated movement
//  *    _gone indicates that the unit has actually left the locations,
//  *    and should not be interacted with anymore.
//  *
//  *    The distinction allows zero time commands to interact with
//  *    the entity on the day movement is begun.
//  */
//
// #define    ship_moving(n)    (rp_subloc(n) ? rp_subloc(n)->moving : 0)
// #define    ship_gone(n)    (ship_moving(n) ? sysclock.days_since_epoch - ship_moving(n) + evening : 0)
//
// #define    char_moving(n)    (rp_char(n) ? rp_char(n)->moving : 0)
//
// #if 0
// #define    char_gone(n)    (char_moving(n) ? sysclock.days_since_epoch - char_moving(n) + evening : 0)
// #else
// #define    char_gone(n)    (char_moving(n) ? 1 : 0)
// #endif
//
// #define player_split_lines(n)    (rp_player(n) ? rp_player(n)->split_lines : 0)
// #define player_split_bytes(n)    (rp_player(n) ? rp_player(n)->split_bytes : 0)
// #define    player_email(n)        (rp_player(n) ? rp_player(n)->email : NULL)
// #define    times_paid(n)        (rp_player(n) ? rp_player(n)->times_paid : 0)
// #define    player_public_turn(n)    (rp_player(n) ? rp_player(n)->public_turn : 0)
// #define    player_format(n)    (rp_player(n) ? rp_player(n)->format : 0)
// #define    player_notab(n)        (rp_player(n) ? rp_player(n)->notab : 0)
// #define    player_compuserve(n)    (rp_player(n) ? rp_player(n)->compuserve : 0)
// #define    player_broken_mailer(n)    (rp_player(n) ? rp_player(n)->broken_mailer : 0)
// #define banner(n)        (rp_misc(n) ? rp_misc(n)->display : NULL)
// #define storm_bind(n)        (rp_misc(n) ? rp_misc(n)->bind_storm : 0)
// #define npc_program(n)        (rp_char(n) ? rp_char(n)->npc_prog : 0)
// #define char_studied(n)        (rp_char(n) ? rp_char(n)->studied : 0)
// #define char_guard(n)        (rp_char(n) ? rp_char(n)->guard : 0)
// #define char_health(n)        (rp_char(n) ? rp_char(n)->health : 0)
// #define char_sick(n)        (rp_char(n) ? rp_char(n)->sick : 0)
// #define loyal_kind(n)        (rp_char(n) ? rp_char(n)->loy_kind : 0)
// #define loyal_rate(n)        (rp_char(n) ? rp_char(n)->loy_rate : 0)
// #define noble_item(n)        (rp_char(n) ? rp_char(n)->unit_item : 0)
// #define char_new_lord(n)    (rp_char(n) ? rp_char(n)->new_lord : 0)
// #define char_melt_me(n)        (rp_char(n) ? rp_char(n)->melt_me : 0)
// #define char_behind(n)        (rp_char(n) ? rp_char(n)->behind : 0)
// #define char_pray(n)        (rp_magic(n) ? rp_magic(n)->pray : 0)
// #define char_hidden(n)        (rp_magic(n) ? rp_magic(n)->hide_self : 0)
// #define vision_protect(n)    (rp_magic(n) ? rp_magic(n)->vis_protect : 0)
// #define default_garrison(n)    (rp_magic(n) ? rp_magic(n)->default_garr : 0)
// #define char_hide_mage(n)    (rp_magic(n) ? rp_magic(n)->hide_mage : 0)
// #define char_cur_aura(n)    (rp_magic(n) ? rp_magic(n)->cur_aura : 0)
// #define char_max_aura(n)    (rp_magic(n) ? rp_magic(n)->max_aura : 0)
// #define reflect_blast(n)    (rp_magic(n) ? rp_magic(n)->aura_reflect : 0)
// #define char_pledge(n)        (rp_magic(n) ? rp_magic(n)->pledge : 0)
// #define char_auraculum(n)    (rp_magic(n) ? rp_magic(n)->auraculum : 0)
// #define char_abil_shroud(n)    (rp_magic(n) ? rp_magic(n)->ability_shroud : 0)
// #define board_fee(n)        (rp_magic(n) ? rp_magic(n)->fee : 0)
// #define ferry_horn(n)        (rp_magic(n) ? rp_magic(n)->ferry_flag : 0)
// #define loc_prominence(n)    (rp_subloc(n) ? rp_subloc(n)->prominence : 0)
// #define loc_opium(n)        (rp_subloc(n) ? rp_subloc(n)->opium_econ : 0)
// #define loc_barrier(n)        (rp_loc(n) ? rp_loc(n)->barrier : 0)
// #define loc_shroud(n)        (rp_loc(n) ? rp_loc(n)->shroud : 0)
// #define loc_civ(n)        (rp_loc(n) ? rp_loc(n)->civ : 0)
// #define loc_sea_lane(n)        (rp_loc(n) ? rp_loc(n)->sea_lane : 0)
// #define ship_cap_raw(n)        (rp_subloc(n) ? rp_subloc(n)->capacity : 0)
// #define body_old_lord(n)    (rp_misc(n) ? rp_misc(n)->old_lord : 0)
// #define gate_dist(n)        (rp_loc(n) ? rp_loc(n)->dist_from_gate : 0)
// #define    road_dest(n)        (rp_gate(n) ? rp_gate(n)->to_loc : 0)
// #define    road_hidden(n)        (rp_gate(n) ? rp_gate(n)->road_hidden : 0)
// #define    item_animal(n)        (rp_item(n) ? rp_item(n)->animal : 0)
// #define    item_prominent(n)    (rp_item(n) ? rp_item(n)->prominent : 0)
// #define    item_weight(n)        (rp_item(n) ? rp_item(n)->weight : 0)
// #define    item_price(n)        (rp_item(n) ? rp_item(n)->base_price : 0)
// #define    item_unique(n)        (rp_item(n) ? rp_item(n)->who_has : 0)
// #define    item_land_cap(n)    (rp_item(n) ? rp_item(n)->land_cap : 0)
// #define    item_ride_cap(n)    (rp_item(n) ? rp_item(n)->ride_cap : 0)
// #define    item_fly_cap(n)        (rp_item(n) ? rp_item(n)->fly_cap : 0)
// #define    req_skill(n)        (rp_skill(n) ? rp_skill(n)->required_skill : 0)
// #define char_persuaded(n)    (rp_char(n) ? rp_char(n)->persuaded : 0)
// #define char_rank(n)        (rp_char(n) ? rp_char(n)->rank : 0)
// #define    skill_produce(n)    (rp_skill(n) ? rp_skill(n)->produced : 0)
// #define    skill_no_exp(n)        (rp_skill(n) ? rp_skill(n)->no_exp : 0)
// #define    skill_np_req(n)        (rp_skill(n) ? rp_skill(n)->np_req : 0)
// #define    skill_aura_req(n)    (rp_skill(n) ? rp_skill(n)->aura_req : 0)
// #define    ship_has_ram(n)        (rp_subloc(n) ? rp_subloc(n)->galley_ram : 0)
// #define    loc_link_open(n)    (rp_subloc(n) ? rp_subloc(n)->link_open : 0)
// #define    loc_damage(n)        (rp_subloc(n) ? rp_subloc(n)->damage : 0)
// #define    loc_defense(n)        (rp_subloc(n) ? rp_subloc(n)->defense : 0)
// #define    loc_pillage(n)        (rp_subloc(n) ? rp_subloc(n)->loot : 0)
// #define    recent_pillage(n)    (rp_subloc(n) ? rp_subloc(n)->recent_loot : 0)
// #define    safe_haven(n)        (rp_subloc(n) ? rp_subloc(n)->safe: 0)
// #define    major_city(n)        (rp_subloc(n) ? rp_subloc(n)->major : 0)
// #define    uldim(n)        (rp_subloc(n) ? rp_subloc(n)->uldim_flag : 0)
// #define    summerbridge(n)        (rp_subloc(n) ? rp_subloc(n)->summer_flag : 0)
// #define    subloc_quest(n)        (rp_subloc(n) ? rp_subloc(n)->quest_late : 0)
// #define    tunnel_depth(n)        (rp_subloc(n) ? rp_subloc(n)->tunnel_level : 0)
// #define    learn_time(n)        (rp_skill(n) ? rp_skill(n)->time_to_learn : 0)
// #define    player_np(n)        (rp_player(n) ? rp_player(n)->noble_points : 0)
// #define    player_fast_study(n)    (rp_player(n) ? rp_player(n)->fast_study : 0)
// #define    gate_seal(n)        (rp_gate(n) ? rp_gate(n)->seal_key : 0)
// #define    gate_dest(n)        (rp_gate(n) ? rp_gate(n)->to_loc : 0)
// #define    char_proj_cast(n)    (rp_magic(n) ? rp_magic(n)->project_cast : 0)
// #define    char_quick_cast(n)    (rp_magic(n) ? rp_magic(n)->quick_cast : 0)
// #define    is_magician(n)        (rp_magic(n) ? rp_magic(n)->magician : 0)
// #define    weather_mage(n)        (rp_magic(n) ? rp_magic(n)->knows_weather : 0)
// #define    garrison_castle(n)    (rp_misc(n) ? rp_misc(n)->garr_castle : 0)
// #define    npc_last_dir(n)        (rp_misc(n) ? rp_misc(n)->npc_dir : 0)
// #define    restricted_control(n)    (rp_misc(n) ? rp_misc(n)->cmd_allow : 0)
// #define    item_capturable(n)    (rp_item(n) ? rp_item(n)->capturable : 0)
// #define    storm_strength(n)    (rp_misc(n) ? rp_misc(n)->storm_str : 0)
// #define    npc_summoner(n)        (rp_misc(n) ? rp_misc(n)->summoned_by : 0)
// #define    char_break(n)        (rp_char(n) ? rp_char(n)->break_point : 0)
// #define    only_defeatable(n)    (rp_misc(n) ? rp_misc(n)->only_vulnerable : 0)
//
// #define    mine_depth(n)      (rp_subloc(n) ? rp_subloc(n)->shaft_depth / 3 : 0)
// #define    release_swear(n)  (rp_magic(n) ? rp_magic(n)->swear_on_release : 0)
// #define    our_token(n)      (rp_magic(n) ? rp_magic(n)->token : 0)
// #define    castle_level(n)      (rp_subloc(n) ? rp_subloc(n)->castle_lev : 0)
//
// #define    item_token_num(n)  (rp_item_magic(n) ? rp_item_magic(n)->token_num : 0)
// #define    item_token_ni(n)   (rp_item_magic(n) ? rp_item_magic(n)->token_ni : 0)
//
// #define    item_lore(n)         (rp_item_magic(n) ? rp_item_magic(n)->lore : 0)
// #define    item_use_key(n)      (rp_item_magic(n) ? rp_item_magic(n)->use_key : 0)
// #define    item_creator(n)      (rp_item_magic(n) ? rp_item_magic(n)->creator : 0)
// #define    item_aura(n)         (rp_item_magic(n) ? rp_item_magic(n)->aura : 0)
// #define    item_creat_cloak(n)  (rp_item_magic(n) ? rp_item_magic(n)->cloak_creator : 0)
// #define    item_reg_cloak(n)    (rp_item_magic(n) ? rp_item_magic(n)->cloak_region : 0)
// #define    item_creat_loc(n)    (rp_item_magic(n) ? rp_item_magic(n)->region_created : 0)
// #define    item_curse_non(n)    (rp_item_magic(n) ? rp_item_magic(n)->curse_loyalty : 0)
//
// #define    item_attack(n)        (rp_item(n) ? rp_item(n)->attack : 0)
// #define    item_defense(n)        (rp_item(n) ? rp_item(n)->defense : 0)
// #define    item_missile(n)        (rp_item(n) ? rp_item(n)->missile : 0)
//
// #define    char_attack(n)        (rp_char(n) ? rp_char(n)->attack : 0)
// #define    char_defense(n)        (rp_char(n) ? rp_char(n)->defense : 0)
// #define    char_missile(n)        (rp_char(n) ? rp_char(n)->missile : 0)
//
// #define    item_attack_bonus(n)             (rp_item_magic(n) ? rp_item_magic(n)->attack_bonus : 0)
// #define    item_defense_bonus(n)             (rp_item_magic(n) ? rp_item_magic(n)->defense_bonus : 0)
// #define    item_missile_bonus(n)             (rp_item_magic(n) ? rp_item_magic(n)->missile_bonus : 0)
// #define    item_aura_bonus(n)             (rp_item_magic(n) ? rp_item_magic(n)->aura_bonus : 0)
//
// #define    item_relic_decay(n)             (rp_item_magic(n) ? rp_item_magic(n)->relic_decay : 0)
//
// #define    is_fighter(n)        (item_attack(n) || item_defense(n) || item_missile(n) || n == item_ghost_warrior)
//
// #define man_item(n)    (rp_item(n) ? rp_item(n)->is_man_item : 0)
// #define    is_priest(n)        has_skill((n), sk_religion)
//
// /*
//  *  Return exactly where a unit is
//  *  May point to another character, a structure, or a region
//  */
//
// #define    loc(n)        (rp_loc_info(n)->where)
//
//
// #include "loop.h"
//
// /*
//  *  Prototypes, defines and externs
//  */
//
// #include "code.h"
// #include "dir.h"
// #include "display.h"
// #include "etc.h"
// #include "garr.h"
// #include "loc.h"
// #include "order.h"
// #include "sout.h"
// #include "stack.h"
// #include "swear.h"
// #include "u.h"
// #include "use.h"
//
// extern char *libdir;
//
// /*
//  *  Saved in libdir/system:
//  */
//
// extern olytime sysclock;        /* current time in Olympia */
// extern int indep_player;        /* independent unit player */
// extern int gm_player;            /* The Fates */
// extern int eat_pl;            /* Order scanner */
// extern int skill_player;        /* Player for skill list report */
// extern int post_has_been_run;
// extern int tunnel_region;
// extern int under_region;
// extern int faery_region;
// extern int faery_player;
// extern int hades_region;
// extern int hades_pit;            /* Pit of Hades */
// extern int hades_player;
// extern int nowhere_region;
// extern int nowhere_loc;
// extern int cloud_region;
// extern int npc_pl;
// extern int garr_pl;            /* garrison player */
// extern int combat_pl;            /* combat log */
// extern int garrison_magic;
// extern int show_to_garrison;
// extern int mount_olympus;
// extern int pledge_backlinks;
//
// #define    in_faery(n)        (region(n) == faery_region)
// #define    in_hades(n)        (region(n) == hades_region)
// #define    in_clouds(n)        (region(n) == cloud_region)
//
// extern int immed_see_all;    /* override hidden-ness, for debugging */
//
// #define    see_all(n)        immed_see_all
//
//
// extern ilist trades_to_check;
// extern char *kind_s[];
// extern char *subkind_s[];
// extern char *dir_s[];
// extern char *loc_depth_s[];
// extern char *short_dir_s[];
// extern char *full_dir_s[];
// extern char *month_names[];
// extern char *entab(int);
// extern int exit_opposite[];
// extern int immediate;
// extern int win_flag;
// extern int indent;
// extern int show_day;
// extern struct cmd_tbl_ent cmd_tbl[];
// extern int evening;            /* are we in the evening phase? */
// extern char *from_host;
// extern char *reply_host;
// extern char *gm_address;
// extern char *game_title;
// extern char *game_url;
// extern char *rules_url;
// extern char *times_url;
// extern char *htpasswd_loc;
// extern int garrison_pay;
// extern int auto_quit_turns;
// extern int army_slow_factor;
// extern ilist new_players;        /* new players added this turn */
//
//
// #define        wait(n)        (rp_command(n) ? rp_command(n)->wait : 0)
// #define        is_prisoner(n)    (rp_char(n) ? rp_char(n)->prisoner : FALSE)
// #define        magic_skill(n)    (subkind(skill_school(n)) == sub_magic)
//
// #define        add_s(n)    ((n) == 1 ? "" : "s")
// #define        add_ds(n)    n, ((n) == 1 ? "" : "s")
//
// #define        alive(n)    (kind(n) == T_char)
//
// #define        is_npc(n)    (subkind(n) || loyal_kind(n) == LOY_npc || loyal_kind(n) == LOY_summon)
//
// #define    char_alone(n)    (stack_parent(n) == 0 && count_stack_any(n) == 1)
// #define    char_really_hidden(n)    (char_hidden(n) && char_alone(n))
//
// #define        CHAR_FIELD    6    /* field length for box_code_less */
//
//
// #define    MAX_POST 60    /* max line length for posts and messages */

// extern void style(int n);
//
// extern char * char_rep_location(int who);

// extern char *top_order(int player, int who);
// extern void pop_order(int player, int who);
// extern void queue_order(int player, int who, char *s);
// extern void prepend_order(int pl, int who, char *s);
// extern void queue_stop(int pl, int who);
// extern void load_orders();
// extern void save_orders();
//
// extern void delete_box(int n);
// extern void change_box_kind(int n, int kind);
// extern void alloc_box(int n, int kind, int sk);
// extern int new_ent(int kind, int sk);
// extern void list_pending_orders_sup(int who, int num, int show_empty);
// extern int stop_order(int pl, int who);

// extern void initialize_buffer();
// extern char *comma_append(char *s, char *t);
// extern char *sout();
//
// extern ilist out_vector;
//
// extern void wrap_set(int who);
// extern void wrap_done();
//
// extern void vector_char_here(int where);
// extern void vector_add(int who);
// extern void vector_stack(int who, int clear);
// extern void vector_clear();
// extern void lines(int who, char *s);
// extern void match_lines(int who, char *s);
// extern void restore_output_vector(ilist t);
// extern ilist save_output_vector();
//
// extern int out_path;            /* alternate sout directive */
// extern int out_alt_who;            /* used if path == MASTER */
//
// extern char *spaces;
// extern int spaces_len;
//
// void log_write(int k, const char *format, ...);
// void out(int who, const char *format, ...);
// #define wout out
//
// extern int here_preceeds(int a, int b);
// extern int stack_parent(int who);
// extern int stack_leader(int who);
// extern int stacked_beneath(int a, int b);
// extern void leave_stack(int who);
// extern void join_stack(int who, int stack);
// extern void drop_stack(int who, int to_drop);
// extern void promote_stack(int lower, int higher);
// extern void take_prisoner(int who, int target);
// extern int has_prisoner(int who, int pris);
// extern int give_prisoner(int who, int target, int pris);
// extern int move_prisoner(int who, int target, int pris);
// extern void promote(int who, int new_pos);
// extern int check_prisoner_escape(int who, int chance);
// extern void prisoner_escapes(int who);
//
// extern void set_lord(int who, int new_lord, int k, int lev);
// extern void swear_char(int who, int target);
// extern void unit_deserts(int who, int to_who, int loy_check, int k, int lev);
// extern void set_loyal(int who, int k, int lev);
// extern int is_unit(int pl, int v);
// extern int lord(int n);
// extern int player(int n);
//
// /*
//   - defines, prototypes and externs for u.c
//     */

// extern void determine_unit_weights(int who, struct weights *w);
// extern void determine_stack_weights(int who, struct weights *w);
// extern int ship_weight(int ship);
//
// extern void olytime_increment(olytime *p);
// extern void olytime_turn_change(olytime *p);
//
// extern int lookup(char *table[], char *s);
// extern char *cap(char *s);
// extern char *nice_num(int n);
// extern char *ordinal(int n);
// extern char *comma_num(int n);
// extern char *knum(int n, int nozero);
// extern char *gold_s(int n);
// extern char *weeks(int n);
// extern char *more_weeks(int n);
//
// extern int has_item(int who, int item);
// extern void gen_item(int who, int item, int qty);
// extern int consume_item(int who, int item, int qty);
// extern int move_item(int from, int to, int item, int qty);
// extern int create_unique_item(int who, int sk);
// extern void destroy_unique_item(int who, int item);
// extern int drop_item(int who, int item, int qty);
//
// extern int charge(int who, int amount);
// extern int can_pay(int who, int amount);
//
// extern void set_known(int who, int i);
// extern int test_known(int who, int i);
// extern void clear_know_rec(sparse *kr);
//
// extern void clear_temps(int kind);
//
// extern int count_loc_char_item(int where, int item);
// extern int check_char_here(int who, int target);
// extern int check_char_gone(int who, int target);
// extern int check_still_here(int who, int target);
// extern int add_structure_damage(int fort, int damage, int can_destroy);
// extern void delta_loyalty(int who, int amount, int silent);
// extern void add_np(int pl, int num);
// extern int deduct_np(int pl, int num);
// extern int deduct_aura(int who, int amount);
// extern int charge_aura(int who, int amount);
// extern int count_man_items(int who);
// extern int count_stack_units(int who);
// extern int count_stack_figures(int who);
// extern int check_aura(int who, int amount);
// extern void get_rid_of_building(int fort);

// extern void add_char_damage(int who, int amount, int inherit);
// extern void take_unit_items(int from, int inherit, int how_many);
// extern int stackmate_inheritor(int who);
// extern void kill_char(int who, int inherit);
// extern int first_char_here(int where);
// extern char *loyal_s(int who);
// extern int find_nearest_land(int where);
// extern int stack_has_item(int who, int item);
// extern int new_char(int, int, int, int, int, int, int, char *);
// extern void kill_stack_ocean(int who);
// extern void hack_unique_item(int item, int owner);
// extern int loc_hidden(int n);
// extern char *rest_name(struct command *c, int a);
// extern int nprovinces();
// extern int has_use_key(int who, int key);
// extern int stack_has_use_key(int who, int key);
// extern int my_prisoner(int who, int pris);
// extern int beast_capturable(int who);
//
// extern int greater_region(int who);
// extern int diff_region(int a, int b);
//
// extern void set_bit(sparse * kr, int i);
//
//
// extern struct skill_ent *p_skill_ent(int who, int skill);
// extern struct skill_ent *rp_skill_ent(int who, int skill);
// extern int has_skill(int who, int skill);
// extern void set_skill(int who, int skill, int know);
//
// extern int skill_school(int sk);
// extern void list_skills(int who, int num);
//
// extern void learn_skill(int who, int sk);
// extern char *exp_s(int level);
// extern int forget_skill(int who, int skill);

// #define    abs(n)        ((n) < 0 ? ((n) * -1) : (n))
//
// #define    isalpha(c)    (((c)>='a' && (c)<='z') || ((c)>='A' && (c)<='Z'))
// #define    isdigit(c)    ((c) >= '0' && (c) <= '9')
// #define    iswhite(c)    ((c) == ' ' || (c) == '\t')
//
// #if 1
// #define    tolower(c)    (lower_array[c])
// extern char lower_array[];
// #else
// #define    tolower(c)    (((c) >= 'A' && (c) <= 'Z') ? ((c) - 'A' + 'a') : (c))
// #endif
//
// #define    toupper(c)    (((c) >= 'a' && (c) <= 'z') ? ((c) - 'a' + 'A') : (c))
//
// extern void *my_malloc(unsigned size);
// extern void *my_realloc(void *ptr, unsigned size);
// extern void my_free(void *ptr);
// extern char *str_save(char *);
//
// extern char *getlin(FILE *);
// extern char *getlin_ew(FILE *);
// extern int i_strncmp(char *s, char *t, int n);
// extern int i_strcmp(char *s, char *t);
// extern int fuzzy_strcmp(char *, char *);
// extern int rnd(int low, int high);
//
// /*
//  *  Assertion verifier
//  */
//
// extern void asfail(char *file, int line, char *cond);
//
// #ifdef __STDC__
// #define    assert(p)    if(!(p)) asfail(__FILE__, __LINE__, #p);
// #else
// #define    assert(p)    if(!(p)) asfail(__FILE__, __LINE__, "p");
// #endif
//
//
// /*
//  *  'ilist' reallocing array definitions
//  */

// #define ilist_len(a)        (((int *)(a)) == NULL ? 0 : ((int *)(a))[-2])
//
// extern void ilist_append(ilist *l, int n);
// extern void ilist_prepend(ilist *l, int n);
// extern void ilist_delete(ilist *l, int i);
// extern void ilist_clear(ilist *l);
// extern void ilist_reclaim(ilist *l);
// extern int ilist_lookup(ilist l, int n);
// extern void ilist_rem_value(ilist *l, int n);
// extern void ilist_rem_value_uniq(ilist *l, int n);
// extern ilist ilist_copy(ilist l);
// extern void ilist_scramble(ilist l);
// extern void ilist_insert(ilist *l, int pos, int n);
//
// /*
//  *  'plist' reallocing array definitions
//  *  (because a pointer doesn't necessarily fit in an int!)
//  */

// #define plist_len(a)        (((int *)(a)) == NULL ? 0 : ((int *)(a))[-2])
//
// extern void plist_append(plist *l, void *p);
// extern void plist_prepend(plist *l, void *p);
// extern void plist_delete(plist *l, int i);
// extern void plist_clear(plist *l);
// extern void plist_reclaim(plist *l);
// extern int plist_lookup(plist l, void *p);
// extern void plist_rem_value(plist *l, void *p);
// extern void plist_rem_value_uniq(plist *l, void *p);
// extern plist plist_copy(plist l);
// extern void plist_scramble(plist l);
// extern void plist_insert(plist *l, int pos, void *p);
//
// extern int readfile(char *path);
// extern char *readlin();
// extern char *readlin_ew();
// extern char *eat_leading_trailing_whitespace(char *s);
//
// extern int int_comp(void * a, void * b);
