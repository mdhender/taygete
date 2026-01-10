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
	unit int    /* unit orders are for */
	l    **char /* ilist of orders for unit */
}

type accept_ent struct {
	item     int /* 0 = any item */
	from_who int /* 0 = anyone, else char or player */
	qty      int /* 0 = any qty */
}

type att_ent struct {
	neutral ilist
	hostile ilist
	defend  ilist
}

type entity_char struct {
	unit_item schar /* unit is made of this kind of item */

	health schar
	sick   schar /* 1=character is getting worse */

	guard    schar /* character is guarding the loc */
	loy_kind schar /* LOY_xxx */
	loy_rate int   /* level with kind of loyalty */

	death_time olytime /* when was character killed */

	skills **skill_ent /* ilist of skills known by char */

	moving    int /* daystamp of beginning of movement */
	unit_lord int /* who is our owner? */
	prev_lord int /* who was our previous owner? */

	contact ilist /* who have we contacted, also, who has found us */

	x_char_magic *char_magic

	prisoner    schar /* is this character a prisoner? */
	behind      schar /* are we behind in combat? */
	time_flying schar /* time airborne over ocean */
	break_point schar /* break point when fighting */
	rank        schar /* noble peerage status */
	npc_prog    schar /* npc program */

	attack  short /* fighter attack rating */
	defense short /* fighter defense rating */
	missile short /* capable of missile attacks? */

	/* not saved: */

	melt_me    schar        /* in process of melting away */
	fresh_hire schar        /* don't erode loyalty */
	new_lord   schar        /* got a new lord this turn */
	studied    schar        /* num days we studied */
	accept     **accept_ent /* what we can be given */
}

type char_magic struct {
	max_aura  int /* maximum aura level for magician */
	cur_aura  int /* current aura level for magician */
	auraculum int /* char created an auraculum */

	visions sparse /* visions revealed */
	pledge  int    /* lands are pledged to another */
	token   int    /* we are controlled by this art */
	fee     int    /* gold/100 wt. to board this ship */

	project_cast   int   /* project next cast */
	quick_cast     short /* speed next cast */
	ability_shroud short

	hide_mage         schar /* hide magician status */
	hinder_meditation schar
	magician          schar /* is a magician */
	pray              schar /* have prayed */
	aura_reflect      schar /* reflect aura blast */
	hide_self         schar /* character is hidden */
	swear_on_release  schar /* swear to one who frees us */
	knows_weather     schar /* knows weather magic */
	vis_protect       schar /* vision protection level */
	default_garr      schar /* default initial garrison */

	/* not saved: */

	mage_worked   schar /* worked this month -- not saved */
	ferry_flag    schar /* ferry has tooted its horn -- ns */
	pledged_to_us ilist /* temp -- not saved */
}

type skill_ent struct {
	skill        int   /* skill id */
	days_studied int   /* days studied * TOUGH_NUM */
	experience   short /* experience level with skill */
	know         char  /* SKILL_xxx */

	/* not saved: */

	exp_this_month char /* flag for add_skill_experience() */
}

type item_ent struct {
	item int
	qty  int
}

type entity_loc struct {
	prov_dest      ilist /* province destinations */
	shroud         short /* magical scry shroud */
	barrier        short /* magical barrier */
	civ            schar /* civilization level (0 = wild) */
	hidden         schar /* is location hidden? */
	dist_from_gate schar
	sea_lane       schar /* fast ocean travel here, also "tracks" for npc ferries */
	next           int   /* temp loc list link */
}

type entity_subloc struct {
	teaches    ilist /* skills location offers */
	opium_econ int   /* addiction level of city */
	defense    int   /* defense rating of structure */

	loot        schar /* loot & pillage level */
	recent_loot schar /* pillaged this month -- not saved */
	damage      uchar /* 0=none, 100=fully destroyed */
	galley_ram  schar /* galley is fitted with a ram */
	shaft_depth short /* depth of mine shaft */
	castle_lev  schar /* level of castle improvement */

	build_materials int /* fifths of materials we've used */
	effort_required int /* not finished if nonzero */
	effort_given    int

	moving   int /* daystamp of beginning of movement */
	capacity int /* capacity of ship */

	near_cities  ilist /* cities rumored to be nearby */
	safe         schar /* safe haven */
	major        schar /* major city */
	prominence   schar /* prominence of city */
	uldim_flag   schar /* Uldim pass */
	summer_flag  schar /* Summerbridge */
	quest_late   schar /* quest decay counter */
	tunnel_level schar /* depth of tunnel */

	link_when    schar /* month link is open, -1 = never */
	link_open    schar /* link is open now */
	link_to      ilist /* where we are linked to */
	link_from    ilist /* where we are linked from */
	bound_storms ilist /* storms bound to this ship */
}

type entity_item struct {
	weight   short /* item weight */
	land_cap short /* land carrying capacity */
	ride_cap short /* ride carrying capacity */
	fly_cap  short /* fly carrying capacity */
	attack   short /* fighter attack rating */
	defense  short /* fighter defense rating */
	missile  short /* capable of missile attacks? */

	is_man_item schar /* unit is a character like thing */
	animal      schar /* unit is or contains a horse or an ox */
	prominent   schar /* big things that everyone sees */
	capturable  schar /* ni-char contents are capturable */

	plural_name *char /* plural name of item */
	base_price  int   /* base price of item for market seeding */
	who_has     int   /* who has this unique item */

	x_item_magic *item_magic
}

type item_magic struct {
	creator        int /* who created this item */
	region_created int /* where was it created */
	lore           int /* deliver this lore for the item */

	curse_loyalty schar /* curse noncreator loyalty */
	cloak_region  schar
	cloak_creator schar
	use_key       schar /* special use action */

	may_use   ilist /* list of usable skills via this */
	may_study ilist /* list of skills studying from this */

	project_cast int   /* stored projected cast */
	token_ni     int   /* ni for controlled npc units */
	quick_cast   short /* stored quick cast */

	aura_bonus  short /* aura bonus */
	aura        short /* auraculum aura */
	relic_decay short /* countdown timer */

	attack_bonus  schar /* attack bonus */
	defense_bonus schar /* defense bonus */
	missile_bonus schar /* missile bonus */

	token_num    schar /* how many token controlled units */
	orb_use_ount schar /* how many uses left in the orb */

	/* not saved: */

	one_turn_use schar /* flag for one use per turn */
}

type entity_skill struct {
	time_to_learn  int   /* days of study req'd to learn skill */
	required_skill int   /* skill required to learn this skill */
	np_req         int   /* noble points required to learn */
	offered        ilist /* skills learnable after this one */
	research       ilist /* skills researchable with this one */

	req      **req_ent /* ilist of items required for use or cast */
	produced int       /* simple production skill result */

	no_exp int /* this skill not rated for experience */

	/* not saved: */

	use_count    int /* times skill used during turn */
	last_use_who int /* who last used the skill (this turn) */
}

type req_ent struct {
	item    int   /* item required to use */
	qty     int   /* quantity required */
	consume schar /* REQ_xx */
}

type entity_gate struct {
	to_loc        int   /* destination of gate */
	notify_jumps  int   /* whom to notify */
	notify_unseal int   /* whom to notify */
	seal_key      short /* numeric gate password */
	road_hidden   schar /* this is a hidden road or passage */
}

type entity_misc struct {
	display     *char  /* entity display banner */
	npc_created int    /* turn peasant mob created */
	npc_home    int    /* where npc was created */
	npc_cookie  int    /* allocation cookie item for us */
	summoned_by int    /* who summoned us? */
	save_name   *char  /* orig name of noble for dead bodies */
	old_lord    int    /* who did this dead body used to belong to */
	npc_memory  sparse /* npc memory */
	only_vuln   int    /* only defeatable with this rare artifact */
	garr_castle int    /* castle which owns this garrison */
	bind_storm  int    /* storm bound to this ship */

	storm_str  short /* storm strength */
	npc_dir    schar /* last direction npc moved */
	mine_delay schar /* time until collapsed mine vanishes */
	cmd_allow  char  /* unit under restricted control */

	/* not saved: */

	opium_double schar  /* improved opium production -- not saved */
	post_txt     **char /* text of posted sign -- not saved */
	storm_move   int    /* next loc storm will move to -- not saved */
	garr_watch   ilist  /* units garrison watches for -- not saved */
	garr_host    ilist  /* units garrison will attack -- not saved */
	garr_tax     int    /* garrison taxes collected -- not saved */
	garr_forward int    /* garrison taxes forwarded -- not saved */
}

type wait_arg struct {
	tag  int
	a1   int
	a2   int
	flag *char
}

type command struct {
	who  int /* entity this is under (redundant) */
	wait int /* time until completion */
	cmd  int /* index into cmd_tbl */

	use_skill      int /* skill we are using, if any */
	use_ent        int /* index into use_tbl[] for skill usage */
	use_exp        int /* experience level at using this skill */
	days_executing int /* how long has this command been running */

	a int /* command arguments */
	b int
	c int
	d int
	e int
	f int
	g int
	h int

	line        *char  /* original command line */
	parsed_line *char  /* cut-up line, pointed to by parse */
	parse       **char /* ilist of parsed arguments */

	state          schar /* STATE_LOAD, STATE_RUN, STATE_ERROR, STATE_DONE */
	status         schar /* success or failure */
	poll           schar /* call finish routine each day? */
	pri            schar /* command priority or precedence */
	conditional    schar /* 0=none 1=last succeeded 2=last failed */
	inhibit_finish schar /* don't call d_xxx */

	/* not saved: */

	fuzzy       schar      /* command matched fuzzy -- not saved */
	second_wait schar      /* delay resulting from auto attacks -- ns */
	wait_parse  **wait_arg /* not saved */
	debug       schar      /* debugging check -- not saved */
}

type commandFunction func(*command) int

type cmd_tbl_ent struct {
	allow *char /* who may execute the command */
	name  *char /* name of command */

	start     commandFunction /* initiator */
	finish    commandFunction /* conclusion */
	interrupt commandFunction /* interrupted order */

	time int /* how long command takes */
	poll int /* call finish each day, not just at end */
	pri  int /* command priority or precedence */
}

type trade struct {
	kind       int /* BUY or SELL */
	item       int
	qty        int
	cost       int
	cloak      int /* don't reveal identity of trader */
	have_left  int
	month_prod int /* month city produces item */
	expire     int /* countdown timer for tradegoods */
	who        int /* redundant -- not saved */
	sort       int /* temp key for sorting -- not saved */
}

type admit struct {
	targ  int /* char or loc admit is declared for */
	sense int /* 0=default no, 1=all but.. */
	l     ilist

	flag int /* first time set this turn -- not saved */
}

type weights struct {
	animals int

	total_weight int /* total weight of unit or stack */

	land_cap    int /* carryable weight on land */
	land_weight int

	ride_cap    int /* carryable weight on horseback */
	ride_weight int

	fly_cap    int
	fly_weight int
}

type void any
type plist **any
