# Sprint 26: u.c/c1.c/c2.c – Remaining Helpers & Complex Commands

## Status: TODO

Sprint 26 completes the port of `u.c`, `c1.c`, and `c2.c` by implementing lifecycle/death logic, ship/building
destruction, and all remaining complex commands (exploration, inventory transfer, WAIT/FLAG, FORM, ferries, training).

---

## Scope

- **u.c**: Character lifecycle & death, ship/building destruction, char visibility & skill checks, remaining misc
  utilities
- **c1.c**: Exploration, inventory transfer (GIVE/GET/CLAIM), noble formation (FORM), WAIT/FLAG system
- **c2.c**: Ships & ferries, combat training, opium improvement, v_die

---

## Dependencies

Sprint 26 depends on Sprint 25 helpers:

- Inventory & gold helpers
- Weight & capacity system
- NP/aura accounting
- Region/visibility helpers

---

## Tasks

### 26.1 – Character Lifecycle & Death (`u.c`)

| Task                          | Status | Notes                                                       |
|-------------------------------|--------|-------------------------------------------------------------|
| Port `survive_fatal`          | ☐      | Check sk_survive_fatal, forget skill, restore health        |
| Port `char_reclaim`           | ☐      | Mark for melting, trigger kill_char                         |
| Port `kill_stack_ocean`       | ☐      | Kill stack at sea, survivors swim to nearest land           |
| Port `stackmate_inheritor`    | ☐      | Find inheritor for items                                    |
| Port `take_unit_items`        | ☐      | Transfer inventory on death                                 |
| Port `add_char_damage`        | ☐      | Damage accumulation                                         |
| Port `dead_char_body`         | ☐      | Create dead body item                                       |
| Port `restore_dead_body`      | ☐      | Revive char from dead body                                  |
| Port `put_back_cookie`        | ☐      | Handle tokens                                               |
| Port `kill_char`              | ☐      | Main death pipeline                                         |
| Unit tests for death pipeline | ☐      | With/without survive_fatal, dead body creation, NPC cleanup |

---

### 26.2 – Ship & Building Destruction (`u.c`)

| Task                             | Status | Notes                                                        |
|----------------------------------|--------|--------------------------------------------------------------|
| Port `sink_ship`                 | ☐      | Ship destruction, unbind storms                              |
| Port `get_rid_of_collapsed_mine` | ☐      | Mine collapse cleanup                                        |
| Port `building_collapses`        | ☐      | Building destruction                                         |
| Port `add_structure_damage`      | ☐      | Damage with can_destroy flag                                 |
| Port `find_nearest_land`         | ☐      | Find closest land province from ocean                        |
| Unit tests for destruction       | ☐      | Damage accumulation, collapse triggers, item/char relocation |

---

### 26.3 – Character Visibility & Skill Checks (`u.c`)

| Task                                       | Status | Notes                                          |
|--------------------------------------------|--------|------------------------------------------------|
| Port `contacted`                           | ☐      | Check if contacted by                          |
| Port `char_where`, `char_here`             | ☐      | Character location checks                      |
| Port `check_char_where`, `check_char_here` | ☐      | Validation with error messages                 |
| Port `check_char_gone`, `check_still_here` | ☐      | Movement state checks                          |
| Port `check_skill`                         | ☐      | Skill requirement check with message           |
| Unit tests for visibility checks           | ☐      | Hidden chars, contacted, garrison, error paths |

---

### 26.4 – Exploration Commands (`c1.c`)

| Task                       | Status | Notes                                 |
|----------------------------|--------|---------------------------------------|
| Port `find_lost_items`     | ☐      | Find unique items in location         |
| Port `v_explore`           | ☐      | Exploration start (trivial)           |
| Port `d_explore`           | ☐      | Exploration execution with RNG        |
| Unit tests for exploration | ☐      | No features, hidden exits, lost items |

---

### 26.5 – Inventory Transfer Commands (`c1.c`)

| Task                                  | Status | Notes                                        |
|---------------------------------------|--------|----------------------------------------------|
| Port `v_accept`                       | ☐      | Set accept rules                             |
| Port `will_accept_sup`, `will_accept` | ☐      | Check accept rules                           |
| Port `how_many`                       | ☐      | Quantity calculation helper                  |
| Port `v_give`                         | ☐      | Give items to target                         |
| Port `v_pay`                          | ☐      | Pay gold (wrapper around give)               |
| Port `may_take`, `v_get`              | ☐      | Take items from target                       |
| Port `v_claim`                        | ☐      | Claim gold from faction with auto-correction |
| Unit tests for transfers              | ☐      | Accept rules, faction/non-faction, garrisons |

---

### 26.6 – Noble Formation (`c1.c`)

| Task                       | Status | Notes                                  |
|----------------------------|--------|----------------------------------------|
| Port `next_np_turn`        | ☐      | Calculate next NP grant turn           |
| Port `print_hiring_status` | ☐      | Display NP status                      |
| Port `print_unformed`      | ☐      | List unformed nobles                   |
| Port `form_new_noble`      | ☐      | Core noble creation logic              |
| Port `v_form`              | ☐      | FORM command start                     |
| Port `d_form`              | ☐      | FORM command execution                 |
| Unit tests for FORM        | ☐      | NP consumption, unformed IDs, stacking |

---

### 26.7 – WAIT/FLAG System (`c1.c`)

| Task                                       | Status | Notes                                             |
|--------------------------------------------|--------|---------------------------------------------------|
| Port `struct flag_ent`, `flag_raised`      | ☐      | Flag signaling                                    |
| Port `v_flag`                              | ☐      | Raise a flag                                      |
| Port `wait_tags` table                     | ☐      | WAIT condition keywords                           |
| Port `clear_wait_parse`, `parse_wait_args` | ☐      | WAIT argument parsing                             |
| Port `check_wait_conditions`               | ☐      | Evaluate all WAIT conditions                      |
| Port `wait_list` global                    | ☐      | Track waiting units                               |
| Port `v_wait`, `d_wait`, `i_wait`          | ☐      | WAIT command lifecycle                            |
| Unit tests for WAIT                        | ☐      | Time, gold, item, unit, loc, ship, weather, flags |

---

### 26.8 – Ships & Ferries (`c2.c`)

| Task                   | Status | Notes                                          |
|------------------------|--------|------------------------------------------------|
| Port `v_fee`           | ☐      | Set boarding fee                               |
| Port `board_message`   | ☐      | Boarding announcement                          |
| Port `v_board`         | ☐      | Board a ferry                                  |
| Port `unboard_message` | ☐      | Disembark announcement                         |
| Port `v_unload`        | ☐      | Unload passengers                              |
| Port `v_ferry`         | ☐      | Sound ferry horn                               |
| Unit tests for ferries | ☐      | Fee, capacity, boarding, unloading, ferry flag |

---

### 26.9 – Combat Training (`c2.c`)

| Task                              | Status | Notes                              |
|-----------------------------------|--------|------------------------------------|
| Port `v_archery`, `d_archery`     | ☐      | Archery training                   |
| Port `v_defense`, `d_defense`     | ☐      | Defense training                   |
| Port `v_swordplay`, `d_swordplay` | ☐      | Swordplay training                 |
| Port `v_fight_to_death`           | ☐      | Set breakpoint to 0                |
| Unit tests for training           | ☐      | Stat increases, breakpoint changes |

---

### 26.10 – Opium & Misc Commands (`c2.c`)

| Task                         | Status | Notes                                           |
|------------------------------|--------|-------------------------------------------------|
| Port `v_improve_opium`       | ☐      | Start opium improvement                         |
| Port `d_improve_opium`       | ☐      | Execute opium improvement                       |
| Port `v_die`                 | ☐      | Suicide command (calls kill_char)               |
| Unit tests for misc commands | ☐      | Opium in poppy_field only, v_die triggers death |

---

### 26.11 – Remaining Utilities (`u.c`)

| Task                                   | Status | Notes                            |
|----------------------------------------|--------|----------------------------------|
| Port `first_char_here`                 | ☐      | Find first character at location |
| Port `bark_dogs`                       | ☐      | Hound barking alerts             |
| Port or deprecate `print_dot`, `stage` | ☐      | Debug/progress helpers           |
| Minimal tests                          | ☐      | Verify basic functionality       |

---

## Files to Create/Modify

- `lifecycle.go` – Character death & revival
- `destruction.go` – Ship/building destruction
- `visibility.go` – Char visibility & skill checks (or add to existing)
- `cmd_explore.go` – Exploration commands
- `cmd_transfer.go` – Inventory transfer commands
- `cmd_form.go` – Noble formation
- `cmd_wait.go` – WAIT/FLAG system
- `cmd_ferry.go` – Ferry commands
- `cmd_training.go` – Combat training commands
- `*_test.go` – Unit tests for each module

---

## Technical Notes

### Death Pipeline (`kill_char`)

The death pipeline is complex with many branches:

1. Check `survive_fatal` (sk_survive_fatal skill)
2. Handle tokens (`our_token`, `put_back_cookie`)
3. Create dead body or destroy (NPC, at sea, melt_me flag)
4. Transfer items to stackmate inheritor
5. Handle prisoners (release or transfer)
6. Flush orders, interrupt current command
7. Clear from garrison, wait_list, contact lists
8. Possibly move to Hades (sk_transcend_death)

Test each branch independently with minimal fixtures.

### WAIT Condition Tags

From `wait_tags[]` in C:

- `time`, `day`, `turn` – temporal
- `gold`, `item`, `have` – inventory
- `unit`, `char` – character presence
- `loc`, `at` – location
- `ship` – ship presence
- `rain`, `fog`, `wind`, `clear` – weather
- `ferry` – ferry horn signal
- `flag` – player flag
- `top`, `owner`, `stack` – stack position
- `not` – negation modifier

### Ferry Integration

`v_ferry` sets `p_magic(ship)->ferry_flag = TRUE`, which:

- Wakes any WAIT FERRY conditions in the same location
- Is checked by `ferry_horn()` helper

Ensure WAIT system handles this during condition checking.

---

## Risks

- **Death pipeline complexity**: Many interdependent systems (skills, tokens, items, orders). Test incrementally.
- **WAIT conditions**: 20+ condition types with `not` modifier. Use deterministic test scenarios.
- **Ship/ferry integration**: Depends on movement system partially ported in earlier sprints.

---

## Verification

After Sprint 26, the following commands should be fully functional:

- EXPLORE
- GIVE, GET, PAY, CLAIM, ACCEPT
- FORM
- WAIT, FLAG
- BOARD, UNLOAD, FERRY, FEE
- ARCHERY, DEFENSE, SWORDPLAY, FIGHT
- DIE
- Implicit: character death from combat/damage
