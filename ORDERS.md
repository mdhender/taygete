# Orders

This document outlines the order processing system for the game engine.

## Order Types
1. Main commands
2. Meta
3. Debug/GM
4. Aliases

### Main Commands (Player/Character)

| Command    | Function       |
|------------|----------------|
| accept     | `v_accept`     |
| admit      | `v_admit`      |
| attack     | `v_attack`     |
| banner     | `v_banner`     |
| behind     | `v_behind`     |
| bind       | `v_bind_storm` |
| board      | `v_board`      |
| breed      | `v_breed`      |
| bribe      | `v_bribe`      |
| build      | `v_build`      |
| buy        | `v_buy`        |
| catch      | `v_catch`      |
| claim      | `v_claim`      |
| collect    | `v_collect`    |
| contact    | `v_contact`    |
| credit     | `v_credit`     |
| decree     | `v_decree`     |
| default    | `v_att_clear`  |
| defend     | `v_defend`     |
| die        | `v_die`        |
| discard    | `v_discard`    |
| drop       | `v_discard`    |
| emote      | `v_emote`      |
| execute    | `v_execute`    |
| explore    | `v_explore`    |
| fee        | `v_fee`        |
| ferry      | `v_ferry`      |
| fish       | `v_fish`       |
| flag       | `v_flag`       |
| fly        | `v_fly`        |
| forget     | `v_forget`     |
| form       | `v_form`       |
| format     | `v_format`     |
| garrison   | `v_garrison`   |
| get        | `v_get`        |
| give       | `v_give`       |
| go         | `v_move`       |
| guard      | `v_guard`      |
| hide       | `v_hide`       |
| honor      | `v_honor`      |
| honour     | `v_honor`      |
| hostile    | `v_hostile`    |
| improve    | `v_improve`    |
| incite     | `v_incite`     |
| make       | `v_make`       |
| mallorn    | `v_mallorn`    |
| message    | `v_message`    |
| move       | `v_move`       |
| name       | `v_name`       |
| neutral    | `v_neutral`    |
| notab      | `v_notab`      |
| oath       | `v_oath`       |
| opium      | `v_opium`      |
| pay        | `v_pay`        |
| pillage    | `v_pillage`    |
| pledge     | `v_pledge`     |
| plugh      | `v_plugh`      |
| post       | `v_post`       |
| press      | `v_press`      |
| promote    | `v_promote`    |
| public     | `v_public`     |
| quarry     | `v_quarry`     |
| quest      | `v_quest`      |
| quit       | `v_quit`       |
| raise      | `v_raise`      |
| rally      | `v_rally`      |
| raze       | `v_raze`       |
| realname   | `v_fullname`   |
| reclaim    | `v_reclaim`    |
| recruit    | `v_recruit`    |
| repair     | `v_repair`     |
| research   | `v_research`   |
| rumor      | `v_rumor`      |
| sail       | `v_sail`       |
| seek       | `v_seek`       |
| sell       | `v_sell`       |
| sneak      | `v_sneak`      |
| split      | `v_split`      |
| stack      | `v_stack`      |
| stone      | `v_quarry`     |
| study      | `v_study`      |
| surrender  | `v_surrender`  |
| swear      | `v_swear`      |
| take       | `v_get`        |
| terrorize  | `v_terrorize`  |
| times      | `v_times`      |
| torture    | `v_torture`    |
| train      | `v_make`       |
| trance     | `v_trance`     |
| ungarrison | `v_ungarrison` |
| unload     | `v_unload`     |
| unstack    | `v_unstack`    |
| use        | `v_use`        |
| wait       | `v_wait`       |
| wood       | `v_wood`       |
| xyzzy      | `v_xyzzy`      |
| yew        | `v_yew`        |

### Meta (Email parsing, not playable)

| Command   | Function |
|-----------|----------|
| begin     | (none)   |
| email     | (none)   |
| end       | (none)   |
| flush     | (none)   |
| lore      | (none)   |
| passwd    | (none)   |
| password  | (none)   |
| players   | (none)   |
| resend    | (none)   |
| stop      | `v_stop` |
| unit      | (none)   |
| vis_email | (none)   |

### Debug/GM (Internal only)

| Command    | Function       |
|------------|----------------|
| additem    | `v_add_item`   |
| be         | `v_be`         |
| ct         | `v_ct`         |
| dump       | `v_dump`       |
| fix        | `v_fix`        |
| fix2       | `v_fix2`       |
| h          | `v_listcmds`   |
| i          | `v_invent`     |
| kill       | `v_kill`       |
| know       | `v_know`       |
| l          | `v_look`       |
| look       | `v_look`       |
| los        | `v_los`        |
| makeloc    | `v_makeloc`    |
| poof       | `v_poof`       |
| postproc   | `v_postproc`   |
| relore     | `v_relore`     |
| remail     | `v_remail`     |
| save       | `v_save`       |
| seeall     | `v_see_all`    |
| seed       | `v_seed`       |
| seedmarket | `v_seedmarket` |
| sheet      | `v_lore`       |
| sk         | `v_skills`     |
| subitem    | `v_sub_item`   |
| tp         | `v_take_pris`  |

### Movement (Aliases)

Aliases that translate to `move` commands.

| Command | Function  | Arguments           | Calls       | Example       |
|---------|-----------|---------------------|-------------|---------------|
| e       | `v_east`  | `struct command *c` | `v_move(c)` | `e`           |
| east    | `v_east`  | `struct command *c` | `v_move(c)` | `east`        |
| enter   | `v_enter` | `struct command *c` | `v_move(c)` | `enter tower` |
| exit    | `v_exit`  | `struct command *c` | `v_move(c)` | `exit`        |
| in      | `v_enter` | `struct command *c` | `v_move(c)` | `in`          |
| n       | `v_north` | `struct command *c` | `v_move(c)` | `n`           |
| north   | `v_north` | `struct command *c` | `v_move(c)` | `north`       |
| out     | `v_exit`  | `struct command *c` | `v_move(c)` | `out`         |
| s       | `v_south` | `struct command *c` | `v_move(c)` | `s`           |
| south   | `v_south` | `struct command *c` | `v_move(c)` | `south`       |
| w       | `v_west`  | `struct command *c` | `v_move(c)` | `w`           |
| west    | `v_west`  | `struct command *c` | `v_move(c)` | `west`        |

