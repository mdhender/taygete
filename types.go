// Copyright 2026 Michael D Henderson. All rights reserved.

package taygete

type exit_view struct {
	direction     int /* which direction does the exit go */
	destination   int /* where the exit goes */
	orig          int /* loc we're coming from */
	distance      int /* how far, in days */
	impassable    int /* set if not possible to go there */
	dest_hidden   int /* set if destination hidden */
	orig_hidden   int /* set if origination or road is hidden */
	hidden        int /* set if hidden destination unknown to us */
	inside        int /* different region destinion is in */
	road          int /* road entity number, if this is a road */
	water         int /* is a water link */
	in_transit    int /* is link to a ship that is moving? */
	magic_barrier int /* a magical barrier prevents travel */
	hades_cost    int /* Gate Spirit of Hades fee to enter */
}

type uchar = uint8
type schar = int8

type sparse *int
type ilist *int

type short = int

type olytime struct {
	day              short /* day of month */
	turn             short /* turn number */
	days_since_epoch int   /* days since game begin */
}

type loc_info struct {
	where     int
	here_list ilist
}

type char = int8

type box struct {
	kind  schar
	skind schar
	name  *char

	x_loc_info loc_info
	x_player   *entity_player
	x_char     *entity_char
	x_loc      *entity_loc
	x_subloc   *entity_subloc
	x_item     *entity_item
	x_skill    *entity_skill
	x_gate     *entity_gate
	x_misc     *entity_misc
	x_disp     *att_ent

	cmd    *command
	items  **item_ent /* ilist of items held */
	trades **trade    /* pending buys/sells */

	temp         int /* scratch space */
	output_order int /* for report ordering -- not saved */

	x_next_kind int /* link to next entity of same type */
	x_next_sub  int /* link to next of same subkind */
}

type entity_player struct {
	full_name       *char
	email           *char
	vis_email       *char /* address to put in player list */
	last_email      *char /* where did the last orders come from? */
	password        *char
	first_turn      int          /* which turn was their first? */
	last_order_turn int          /* last turn orders were submitted */
	orders          **order_list /* ilist of orders for units in */
	/* this faction */
	known sparse /* visited, lore seen, encountered */

	units    ilist   /* what units are in our faction? */
	admits   **admit /* admit permissions list */
	unformed ilist   /* nobles as yet unformed */

	split_lines int   /* split mail at this many lines */
	split_bytes int   /* split mail at this many bytes */
	fast_study  short /* instant study days available */

	noble_points short /* how many NP's the player has */

	format        schar /* turn report formatting control */
	notab         schar /* player can't tolerate tabs */
	first_tower   schar /* has player built first tower yet? */
	sent_orders   schar /* sent in orders this turn? */
	dont_remind   schar /* don't send a reminder */
	compuserve    schar /* get Times from CIS */
	broken_mailer schar /* quote begin lines */

	/* not saved: */

	public_turn     char   /* turn is public view */
	times_paid      char   /* only 1 Times credit per player */
	swear_this_turn schar  /* have we used SWEAR this turn? */
	cmd_count       short  /* count of cmds started this turn */
	np_gained       short  /* np's added this turn -- not saved */
	np_spent        short  /* np's lost this turn -- not saved */
	deliver_lore    ilist  /* show these to player -- not saved */
	weather_seen    sparse /* locs we've viewed the weather */
	output          sparse /* units with output -- not saved */
	locs            sparse /* locs we touched -- not saved */
}

type order_list struct {
	//     int unit;            /* unit orders are for */
	//     char **l;            /* ilist of orders for unit */
}

type accept_ent struct {
	//     int item;            /* 0 = any item */
	//     int from_who;            /* 0 = anyone, else char or player */
	//     int qty;            /* 0 = any qty */
}

type att_ent struct {
	//     ilist neutral;
	//     ilist hostile;
	//     ilist defend;
}

type entity_char struct {
	//     int unit_item;            /* unit is made of this kind of item */
	//
	//         schar health;
	//     schar sick;            /* 1=character is getting worse */
	//
	//     schar guard;            /* character is guarding the loc */
	//     schar loy_kind;            /* LOY_xxx */
	//     int loy_rate;            /* level with kind of loyalty */
	//
	//     olytime death_time;        /* when was character killed */
	//
	//     struct skill_ent **skills;    /* ilist of skills known by char */
	//
	//     int moving;            /* daystamp of beginning of movement */
	//     int unit_lord;            /* who is our owner? */
	//     int prev_lord;            /* who was our previous owner? */
	//
	//     ilist contact;            /* who have we contacted, also, who */
	//                     /* has found us */
	//
	//     struct char_magic *x_char_magic;
	//
	//     schar prisoner;            /* is this character a prisoner? */
	//     schar behind;            /* are we behind in combat? */
	//     schar time_flying;        /* time airborne over ocean */
	//     schar break_point;        /* break point when fighting */
	//     schar rank;            /* noble peerage status */
	//     schar npc_prog;            /* npc program */
	//
	//     short attack;            /* fighter attack rating */
	//     short defense;            /* fighter defense rating */
	//     short missile;            /* capable of missile attacks? */
	//
	// /*
	//  *  The following are not saved by io.c:
	//  */
	//
	//     schar melt_me;            /* in process of melting away */
	//     schar fresh_hire;        /* don't erode loyalty */
	//     schar new_lord;            /* got a new lord this turn */
	//     schar studied;            /* num days we studied */
	//     struct accept_ent **accept;    /* what we can be given */
}

type char_magic struct {
	//     int max_aura;            /* maximum aura level for magician */
	//     int cur_aura;            /* current aura level for magician */
	//     int auraculum;            /* char created an auraculum */
	//
	//     sparse visions;            /* visions revealed */
	//     int pledge;            /* lands are pledged to another */
	//     int token;            /* we are controlled by this art */
	//     int fee;            /* gold/100 wt. to board this ship */
	//
	//     int project_cast;        /* project next cast */
	//     short quick_cast;        /* speed next cast */
	//     short ability_shroud;
	//
	//     schar hide_mage;        /* hide magician status */
	//     schar hinder_meditation;
	//     schar magician;            /* is a magician */
	//     schar pray;            /* have prayed */
	//     schar aura_reflect;        /* reflect aura blast */
	//     schar hide_self;        /* character is hidden */
	//     schar swear_on_release;        /* swear to one who frees us */
	//     schar knows_weather;        /* knows weather magic */
	//     schar vis_protect;        /* vision protection level */
	//     schar default_garr;        /* default initial garrison */
	//
	//     schar mage_worked;        /* worked this month -- not saved */
	//     schar ferry_flag;        /* ferry has tooted its horn -- ns */
	//     ilist pledged_to_us;        /* temp -- not saved */
}

type skill_ent struct {
	//     int skill;
	//     int days_studied;        /* days studied * TOUGH_NUM */
	//     short experience;        /* experience level with skill */
	//     char know;            /* SKILL_xxx */
	//
	// /*
	//  *  Not saved:
	//  */
	//
	//     char exp_this_month;        /* flag for add_skill_experience() */
}

type item_ent struct {
	//     int item;
	//     int qty;
}

type entity_loc struct {
	//     ilist prov_dest;
	//     short shroud;            /* magical scry shroud */
	//     short barrier;            /* magical barrier */
	//     schar civ;                /* civilization level (0 = wild) */
	//     schar hidden;            /* is location hidden? */
	//     schar dist_from_gate;
	//     schar sea_lane;            /* fast ocean travel here */
	//                             /* also "tracks" for npc ferries */
	//     int next;                /* temp loc list link */
}

type entity_subloc struct {
	//     ilist teaches;            /* skills location offers */
	//     int opium_econ;            /* addiction level of city */
	//     int defense;            /* defense rating of structure */
	//
	//     schar loot;            /* loot & pillage level */
	//     schar recent_loot;        /* pillaged this month -- not saved */
	//     uchar damage;            /* 0=none, 100=fully destroyed */
	//     schar galley_ram;        /* galley is fitted with a ram */
	//     short shaft_depth;        /* depth of mine shaft */
	//     schar castle_lev;        /* level of castle improvement */
	//
	//     schar build_materials;        /* fifths of materials we've used */
	//     int effort_required;        /* not finished if nonzero */
	//     int effort_given;
	//
	//     int moving;            /* daystamp of beginning of movement */
	//     int capacity;            /* capacity of ship */
	//
	//     ilist near_cities;        /* cities rumored to be nearby */
	//     schar safe;            /* safe haven */
	//     schar major;            /* major city */
	//     schar prominence;        /* prominence of city */
	//     schar uldim_flag;        /* Uldim pass */
	//     schar summer_flag;        /* Summerbridge */
	//     schar quest_late;        /* quest decay counter */
	//     schar tunnel_level;        /* depth of tunnel */
	//
	//     schar link_when;        /* month link is open, -1 = never */
	//     schar link_open;        /* link is open now */
	//     ilist link_to;            /* where we are linked to */
	//     ilist link_from;        /* where we are linked from */
	//     ilist bound_storms;        /* storms bound to this ship */
}

type entity_item struct {
	//         short weight;
	//         short land_cap;
	//         short ride_cap;
	//         short fly_cap;
	//     short attack;        /* fighter attack rating */
	//     short defense;        /* fighter defense rating */
	//     short missile;        /* capable of missile attacks? */
	//
	//     schar is_man_item;    /* unit is a character like thing */
	//     schar animal;        /* unit is or contains a horse or an ox */
	//     schar prominent;    /* big things that everyone sees */
	//     schar capturable;    /* ni-char contents are capturable */
	//
	//     char *plural_name;
	//     int base_price;        /* base price of item for market seeding */
	//     int who_has;        /* who has this unique item */
	//
	//     struct item_magic    *x_item_magic;
}

type item_magic struct {
	//     int creator;
	//     int region_created;
	//     int lore;            /* deliver this lore for the item */
	//
	//     schar curse_loyalty;        /* curse noncreator loyalty */
	//     schar cloak_region;
	//     schar cloak_creator;
	//     schar use_key;            /* special use action */
	//
	//     ilist may_use;            /* list of usable skills via this */
	//     ilist may_study;        /* list of skills studying from this */
	//
	//     int project_cast;        /* stored projected cast */
	//     int token_ni;            /* ni for controlled npc units */
	//     short quick_cast;        /* stored quick cast */
	//
	//     short aura_bonus;
	//     short aura;            /* auraculum aura */
	//     short relic_decay;        /* countdown timer */
	//
	//     schar attack_bonus;
	//     schar defense_bonus;
	//     schar missile_bonus;
	//
	//     schar token_num;        /* how many token controlled units */
	//     schar orb_use_count;        /* how many uses left in the orb */
	//
	// /*
	//  *  Not saved:
	//  */
	//
	//     schar one_turn_use;        /* flag for one use per turn */
}

type entity_skill struct {
	//     int time_to_learn;    /* days of study req'd to learn skill */
	//     int required_skill;    /* skill required to learn this skill */
	//     int np_req;        /* noble points required to learn */
	//     ilist offered;        /* skills learnable after this one */
	//     ilist research;        /* skills researable with this one */
	//
	//     struct req_ent **req;    /* ilist of items required for use or cast */
	//     int produced;        /* simple production skill result */
	//
	//     int no_exp;        /* this skill not rated for experience */
	//
	// /* not saved */
	//
	//     int use_count;        /* times skill used during turn */
	//     int last_use_who;    /* who last used the skill (this turn) */
}

type req_ent struct {
	//     int item;        /* item required to use */
	//     int qty;        /* quantity required */
	//     schar consume;        /* REQ_xx */
}

type entity_gate struct {
	//     int to_loc;        /* destination of gate */
	//     int notify_jumps;    /* whom to notify */
	//     int notify_unseal;    /* whom to notify */
	//     short seal_key;        /* numeric gate password */
	//     schar road_hidden;    /* this is a hidden road or passage */
}

type entity_misc struct {
	//     char *display;        /* entity display banner */
	//     int npc_created;    /* turn peasant mob created */
	//     int npc_home;        /* where npc was created */
	//     int npc_cookie;        /* allocation cookie item for us */
	//     int summoned_by;    /* who summoned us? */
	//     char *save_name;    /* orig name of noble for dead bodies */
	//     int old_lord;        /* who did this dead body used to belong to */
	//     sparse npc_memory;
	//     int only_vulnerable;    /* only defeatable with this rare artifact */
	//     int garr_castle;    /* castle which owns this garrison */
	//     int bind_storm;        /* storm bound to this ship */
	//
	//     short storm_str;    /* storm strength */
	//     schar npc_dir;        /* last direction npc moved */
	//     schar mine_delay;    /* time until collapsed mine vanishes */
	//     char cmd_allow;        /* unit under restricted control */
	//
	//     schar opium_double;    /* improved opium production -- not saved */
	//     char **post_txt;    /* text of posted sign -- not saved */
	//     int storm_move;        /* next loc storm will move to -- not saved */
	//     ilist garr_watch;    /* units garrison watches for -- not saved */
	//     ilist garr_host;    /* units garrison will attack -- not saved */
	//     int garr_tax;        /* garrison taxes collected -- not saved */
	//     int garr_forward;    /* garrison taxes forwarded -- not saved */
}

type wait_arg struct {
	//     int tag;
	//     int a1, a2;
	//     char *flag;
}

type command struct {
	//     int who;        /* entity this is under (redundant) */
	//     int wait;        /* time until completion */
	//     int cmd;        /* index into cmd_tbl */
	//
	//     int use_skill;        /* skill we are using, if any */
	//     int use_ent;        /* index into use_tbl[] for skill usage */
	//     int use_exp;        /* experience level at using this skill */
	//     int days_executing;    /* how long has this command been running */
	//
	//     int a,b,c,d,e,f,g,h;    /* command arguments */
	//
	//     char *line;        /* original command line */
	//     char *parsed_line;    /* cut-up line, pointed to by parse */
	//     char **parse;        /* ilist of parsed arguments */
	//
	//     schar state;        /* STATE_LOAD, STATE_RUN, STATE_ERROR, STATE_DONE */
	//     schar status;        /* success or failure */
	//     schar poll;        /* call finish routine each day? */
	//     schar pri;        /* command priority or precedence */
	//     schar conditional;    /* 0=none 1=last succeeded 2=last failed */
	//     schar inhibit_finish;    /* don't call d_xxx */
	//
	//     schar fuzzy;        /* command matched fuzzy -- not saved */
	//     schar second_wait;    /* delay resulting from auto attacks -- ns */
	//     struct wait_arg **wait_parse;    /* not saved */
	//     schar debug;        /* debugging check -- not saved */
}

type cmd_tbl_ent struct {
	//     char *allow;        /* who may execute the command */
	//     char *name;        /* name of command */
	//
	//     int (*start)(struct command *);        /* initiator */
	//     int (*finish)(struct command *);    /* conclusion */
	//     int (*interrupt)(struct command *);    /* interrupted order */
	//
	//     int time;        /* how long command takes */
	//     int poll;        /* call finish each day, not just at end */
	//     int pri;        /* command priority or precedence */
}

type trade struct {
	//     int kind;        /* BUY or SELL */
	//     int item;
	//     int qty;
	//     int cost;
	//     int cloak;        /* don't reveal identity of trader */
	//     int have_left;
	//     int month_prod;        /* month city produces item */
	//     int expire;        /* countdown timer for tradegoods */
	//     int who;        /* redundant -- not saved */
	//     int sort;        /* temp key for sorting -- not saved */
}

type admit struct {
	//     int targ;        /* char or loc admit is declared for */
	//     int sense;        /* 0=default no, 1=all but.. */
	//     ilist l;
	//
	//     int flag;        /* first time set this turn -- not saved */
}

type weights struct {
	//     int animals;
	//
	//     int total_weight;    /* total weight of unit or stack */
	//
	//     int land_cap;        /* carryable weight on land */
	//     int land_weight;
	//
	//     int ride_cap;        /* carryable weight on horseback */
	//     int ride_weight;
	//
	//     int fly_cap;
	//     int fly_weight;
}

type void any
type plist **any
