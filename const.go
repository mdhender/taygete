// Copyright 2026 Michael D Henderson. All rights reserved.

package taygete

const RAND = 1
const LAND = 1
const WATER = 2

const RANK_lord = 10
const RANK_knight = 20
const RANK_baron = 30
const RANK_count = 40
const RANK_earl = 50
const RANK_marquess = 60
const RANK_duke = 70
const RANK_king = 80

const MAX_BOXES = 150000

const item_gold = 1

const item_peasant = 10
const item_worker = 11
const item_soldier = 12
const item_archer = 13
const item_knight = 14
const item_elite_guard = 15
const item_pikeman = 16
const item_blessed_soldier = 17
const item_ghost_warrior = 18
const item_sailor = 19
const item_swordsman = 20
const item_crossbowman = 21
const item_elite_arch = 22
const item_angry_peasant = 23
const item_pirate = 24
const item_elf = 25
const item_spirit = 26

const item_corpse = 31
const item_savage = 32
const item_skeleton = 33
const item_barbarian = 34

const item_wild_horse = 51
const item_riding_horse = 52
const item_warmount = 53
const item_pegasus = 54
const item_nazgul = 55

const item_flotsam = 59
const item_battering_ram = 60
const item_catapult = 61
const item_siege_tower = 62
const item_ratspider_venom = 63
const item_lana_bark = 64
const item_avinia_leaf = 65
const item_spiny_root = 66
const item_farrenstone = 67
const item_yew = 68
const item_elfstone = 69
const item_mallorn_wood = 70
const item_pretus_bones = 71
const item_longbow = 72
const item_plate = 73
const item_longsword = 74
const item_pike = 75
const item_ox = 76
const item_lumber = 77
const item_stone = 78
const item_iron = 79
const item_leather = 80
const item_ratspider = 81
const item_mithril = 82
const item_gate_crystal = 83
const item_blank_scroll = 84
const item_crossbow = 85
const item_fish = 87
const item_opium = 93
const item_basket = 94 // woven basket
const item_pot = 95    // clay pot
const item_tax_cookie = 96
const item_drum = 98
const item_hide = 99
const item_mob_cookie = 101
const item_lead = 102

const item_glue = 261

const item_centaur = 271
const item_minotaur = 272
const item_undead_cookie = 273
const item_fog_cookie = 274
const item_wind_cookie = 275
const item_rain_cookie = 276
const item_mage_menial = 277 // mage menial labor cookie
const item_spider = 278      // giant spider
const item_rat = 279         // horde of rats
const item_lion = 280
const item_bird = 281 // giant bird
const item_lizard = 282
const item_bandit = 283
const item_chimera = 284
const item_harpie = 285
const item_dragon = 286
const item_orc = 287
const item_gorgon = 288
const item_wolf = 289
const item_orb = 290
const item_cyclops = 291
const item_giant = 292
const item_faery = 293
const item_petty_thief = 294
const item_hound = 295

const lore_skeleton_npc_token = 931
const lore_orc_npc_token = 932
const lore_undead_npc_token = 933
const lore_savage_npc_token = 934
const lore_barbarian_npc_token = 935
const lore_orb = 936
const lore_faery_stone = 937
const lore_barbarian_kill = 938
const lore_savage_kill = 939
const lore_undead_kill = 940
const lore_orc_kill = 941
const lore_skeleton_kill = 942

const sk_shipcraft = 600
const sk_pilot_ship = 601
const sk_shipbuilding = 602
const sk_fishing = 603

const sk_combat = 610
const sk_survive_fatal = 611
const sk_fight_to_death = 612
const sk_make_catapult = 613
const sk_defense = 614
const sk_archery = 615
const sk_swordplay = 616
const sk_weaponsmith = 617

const sk_stealth = 630
const sk_petty_thief = 631
const sk_spy_inv = 632    // determine char inventory
const sk_spy_skills = 633 // determine char skill
const sk_spy_lord = 634   // determine char's lord
const sk_hide_lord = 635
const sk_find_rich = 636
const sk_torture = 637
const sk_hide_self = 638
const sk_sneak_build = 639

const sk_beast = 650
const sk_bird_spy = 651
const sk_capture_beasts = 652
const sk_use_beasts = 653 // use beasts in battle
const sk_breed_beasts = 654
const sk_catch_horse = 655
const sk_train_wild = 656
const sk_train_warmount = 657
const sk_summon_savage = 658
const sk_keep_savage = 659
const sk_breed_hound = 661

const sk_persuasion = 670
const sk_bribe_noble = 671
const sk_persuade_oath = 672
const sk_raise_mob = 673
const sk_rally_mob = 674
const sk_incite_mob = 675
const sk_train_angry = 676

const sk_construction = 680
const sk_make_siege = 681
const sk_quarry_stone = 682

const sk_alchemy = 690
const sk_brew_heal = 691
const sk_record_skill = 692
const sk_extract_venom = 693 // from ratspider
const sk_brew_slave = 694    // potion of slavery
const sk_collect_elem = 695
const sk_brew_death = 696
const sk_lead_to_gold = 697

const sk_forestry = 700
const sk_make_ram = 701 // make battering ram
const sk_harvest_lumber = 702
const sk_harvest_yew = 703
const sk_collect_foliage = 704
const sk_harvest_mallorn = 705
const sk_harvest_opium = 706
const sk_improve_opium = 707

const sk_mining = 720
const sk_mine_iron = 721
const sk_mine_gold = 722
const sk_mine_mithril = 723

const sk_trade = 730
const sk_cloak_trade = 731
const sk_find_sell = 732
const sk_find_buy = 733

const sk_religion = 750
const sk_reveal_vision = 751
const sk_last_rites = 752
const sk_pray = 753
const sk_resurrect = 754
const sk_remove_bless = 755
const sk_vision_protect = 756

const sk_basic = 800
const sk_meditate = 801
const sk_mage_menial = 802 // menial labor for mages
const sk_appear_common = 803
const sk_view_aura = 804
const sk_heal = 805
const sk_write_basic = 806
const sk_reveal_mage = 807 // reveal abilities of mage
const sk_tap_health = 808
const sk_shroud_abil = 809 // ability shroud
const sk_detect_abil = 811 // detect ability scry
const sk_dispel_abil = 812 // dispel ability shroud
const sk_adv_med = 813     // advanced meditation
const sk_hinder_med = 814  // hinder meditation

const sk_weather = 820
const sk_fierce_wind = 821
const sk_bind_storm = 822
const sk_write_weather = 823
const sk_summon_wind = 824
const sk_summon_rain = 825
const sk_summon_fog = 826
const sk_direct_storm = 827
const sk_dissipate = 828
const sk_renew_storm = 829
const sk_lightning = 831
const sk_seize_storm = 832
const sk_death_fog = 833

const sk_scry = 840
const sk_scry_region = 841
const sk_write_scry = 842
const sk_shroud_region = 843
const sk_dispel_region = 844 // dispel region shroud
const sk_bar_loc = 845       // create location barrier
const sk_unbar_loc = 846
const sk_locate_char = 847
const sk_detect_scry = 848 // detect region scry
const sk_proj_cast = 849   // project next cast
const sk_save_proj = 851   // save projected cast
const sk_banish_corpses = 852

const sk_gate = 860
const sk_detect_gates = 861
const sk_jump_gate = 862
const sk_write_gate = 863
const sk_seal_gate = 864
const sk_unseal_gate = 865
const sk_notify_unseal = 866
const sk_rem_seal = 867 // forcefully unseal gate
const sk_reveal_key = 868
const sk_notify_jump = 869
const sk_teleport = 871
const sk_rev_jump = 872

const sk_artifact = 880
const sk_forge_aura = 881 // forge auraculum
const sk_write_art = 882
const sk_forge_weapon = 883
const sk_forge_armor = 884
const sk_forge_bow = 885
const sk_curse_noncreat = 886 // curse noncreator loyalty
const sk_show_art_creat = 887 // learn who created art
const sk_show_art_reg = 888   // learn where art created
const sk_destroy_art = 889
const sk_cloak_creat = 891
const sk_cloak_reg = 892
const sk_rem_art_cloak = 893 // dispel artifact cloaks
const sk_forge_palantir = 894

const sk_necromancy = 900
const sk_raise_corpses = 901 // summon undead corpses
const sk_summon_ghost = 902  // summon ghost warriors
const sk_write_necro = 903
const sk_undead_lord = 904 // summon undead unit
const sk_renew_undead = 905
const sk_banish_undead = 906
const sk_eat_dead = 907
const sk_aura_blast = 908
const sk_absorb_blast = 909
const sk_transcend_death = 911

const sk_adv_sorcery = 920
const sk_trance = 921
const sk_teleport_item = 922

// dead skills
const sk_quick_cast = 999 // speed next cast
const sk_save_quick = 998 // save speeded cast
const sk_add_ram = 997    // add ram to galley

const PROG_bandit = 1 // wilderness spice
const PROG_subloc_monster = 2
const PROG_npc_token = 3
const PROG_faery_bandit = 4
const PROG_hades_bandit = 5

const use_death_potion = 1
const use_heal_potion = 2
const use_slave_potion = 3
const use_palantir = 4
const use_proj_cast = 5   // stored projected cast
const use_quick_cast = 6  // stored cast speedup
const use_drum = 7        // beat savage's drum
const use_faery_stone = 8 // Faery gate opener
const use_orb = 9         // crystal orb
const use_barbarian_kill = 10
const use_savage_kill = 11
const use_corpse_kill = 12
const use_orc_kill = 13
const use_skeleton_kill = 14
const use_bta_skull = 15

const LOY_UNCHANGED = (-1)
const LOY_unsworn = 0
const LOY_contract = 1
const LOY_oath = 2
const LOY_fear = 3
const LOY_npc = 4
const LOY_summon = 5

const exp_novice = 1 // apprentice
const exp_journeyman = 2
const exp_teacher = 3
const exp_master = 4
const exp_grand = 5 // grand master

const ATT_NONE = 0 // no attitude -- default
const NEUTRAL = 1  // explicitly neutral
const HOSTILE = 2  // attack on sight
const DEFEND = 3   // defend if attacked

const SKILL_dont = 0     // don't know the skill
const SKILL_learning = 1 // in the process of learning it
const SKILL_know = 2     // know it

const REQ_NO = 0  // don't consume item
const REQ_YES = 1 // consume item
const REQ_OR = 2  // or with next

// In-process command structure

const STATE_DONE = 0
const STATE_LOAD = 1
const STATE_RUN = 2
const STATE_ERROR = 3

const BUY = 1
const SELL = 2
const PRODUCE = 3
const CONSUME = 4

// style() tags:

// default style is 0 (regular)

const STYLE_TEXT = 1
const STYLE_HTML = 2
const STYLE_PREV = (-1)

const RELIC_THRONE = 401
const RELIC_CROWN = 402
const RELIC_BTA_SKULL = 403

// Possible destinations of output:

const VECT = (-1) // vector of recipients
// n >= 0: output to entity event log
const MASTER = (-2)

const OUT_SUMMARY = 0
const OUT_BANNER = 0
const OUT_INCLUDE = 1
const OUT_LORE = 2
const OUT_NEW = 3        // new player listing
const OUT_LOC = 4        // location descriptions
const OUT_TEMPLATE = 5   // order template
const OUT_GARR = 6       // garrison log
const OUT_SHOW_POSTS = 7 // show what press and rumor look like
const OUT_HTML_INDEX = 8

// tags for log()

const LOG_CODE = 10    // Code alerts
const LOG_SPECIAL = 11 // Special events
const LOG_DEATH = 12   // Character deaths
const LOG_MISC = 13    // Other junk
const LOG_DROP = 14    // Player drops

// tags for eat.c

const EAT_ERR = 20     // Errors in orders submitted
const EAT_WARN = 21    // Warnings in orders submitted
const EAT_QUEUE = 22   // Current order queues
const EAT_HEADERS = 23 // Email headers bounced back
const EAT_OKAY = 24    // Regular (non-error) output for scanner
const EAT_PLAYERS = 25 // Player list

const MATES = (-1)
const MATES_SILENT = (-2)
const TAKE_ALL = 1
const TAKE_SOME = 2
const TAKE_NI = 3 // noble item: wrapper adds one

const TRUE = 1
const FALSE = 0

const LEN = 2048 // generic string max length
