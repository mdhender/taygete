# Sprint 25: u.c/c1.c/c2.c – Helpers & Easiest Commands

## Status: TODO

Sprint 25 ports foundational helpers from `u.c` and the simplest commands from `c1.c`/`c2.c`. Complex lifecycle, death,
WAIT/FLAG, ferries, and noble formation are deferred to Sprint 26.

---

## Scope

- **u.c**: Inventory/gold helpers, knowledge/bitsets, time/numeric formatters, weight/capacity, region/visibility
  helpers, NP/aura accounting
- **c1.c/c2.c**: Simple metadata commands (NAME, BANNER, LOOK), deprecated text-report commands (stubs), basic economy
  commands (DISCARD, QUIT)

---

## Tasks

### 25.1 – Inventory & Gold Helpers (`u.c`)

| Task                                                  | Status | Notes                                           |
|-------------------------------------------------------|--------|-------------------------------------------------|
| Port `has_item`, `gen_item`, `consume_item`           | ☑      | Basic inventory queries                         |
| Port `move_item`, `hack_unique_item`                  | ☑      | Item transfers                                  |
| Port `create_unique_item`, `create_unique_item_alloc` | ☑      | Unique item creation                            |
| Port `destroy_unique_item`                            | ☑      | Include NP grant for dead bodies                |
| Port `drop_item`                                      | ☑      | Handles unique items, province/ocean logic      |
| Port `can_pay`, `charge`                              | ☑      | Gold payment helpers                            |
| Port `stack_has_item`, `stack_sub_item`               | ☑      | Stack-level inventory                           |
| Port `has_use_key`, `stack_has_use_key`               | ☑      | Item use-key checks                             |
| Port `autocharge`                                     | ☑      | Auto-payment from stack                         |
| Unit tests for inventory helpers                      | ☑      | Unique/non-unique, stack borrowing, drop at sea |

---

### 25.2 – Knowledge & Bitsets (`u.c`)

| Task                                         | Status | Notes                                 |
|----------------------------------------------|--------|---------------------------------------|
| Port `test_bit`, `set_bit`, `clear_know_rec` | ☑      | Bitset operations                     |
| Port `test_known`, `set_known`               | ☑      | Player knowledge tracking             |
| Unit tests for knowledge helpers             | ☑      | Idempotent add/remove, clear on empty |

---

### 25.3 – Time & Numeric Helpers (`u.c`)

| Task                                            | Status | Notes                            |
|-------------------------------------------------|--------|----------------------------------|
| Port `olytime_increment`, `olytime_turn_change` | ☑      | Already in day.go                |
| Port `max`, `min`                               | ☑      | Using Go stdlib (Go 1.21+)       |
| Port `comma_num`, `nice_num`, `knum`, `ordinal` | ☑      | Number formatting in format.go   |
| Port `weeks`, `more_weeks`                      | ☑      | Duration formatting in format.go |
| Port `gold_s`, `loyal_s`, `cap`                 | ☑      | String formatters in format.go   |
| Unit tests for formatters                       | ☑      | format_test.go, code_test.go     |

---

### 25.4 – NP & Aura Accounting (`u.c`)

| Task                                            | Status | Notes                         |
|-------------------------------------------------|--------|-------------------------------|
| Port `deduct_np`, `add_np`                      | ☑      | Noble point accounting        |
| Port `deduct_aura`, `charge_aura`, `check_aura` | ☑      | Aura accounting               |
| Unit tests for NP/aura                          | ☑      | Sufficient/insufficient cases |

---

### 25.5 – Weight & Capacity Helpers (`u.c`)

| Task                           | Status | Notes                            |
|--------------------------------|--------|----------------------------------|
| Port `add_item_weight`         | ☑      | Single item weight calculation   |
| Port `determine_unit_weights`  | ☑      | Unit weight totals               |
| Port `determine_stack_weights` | ☑      | Stack weight/capacity totals     |
| Port `ship_weight`, `ship_cap` | ☑      | Ship capacity with damage factor |
| Unit tests for weights         | ☑      | Unit, stack, ship scenarios      |

---

### 25.6 – Region & Visibility Helpers (`u.c`)

| Task                                 | Status | Notes                       |
|--------------------------------------|--------|-----------------------------|
| Port `loc_depth`                     | ☑      | Map subkind → depth level   |
| Port `loc_hidden`                    | ☑      | Location hidden flag check  |
| Port `nprovinces`                    | ☑      | Province count (cached)     |
| Port `greater_region`, `diff_region` | ☑      | Region comparison           |
| Port `clear_temps`                   | ☑      | Clear temp fields           |
| Port `lookup`                        | ☑      | Table lookup helper         |
| Unit tests for region helpers        | ☑      | Depth mapping, hidden flags |

---

### 25.7 – Simple Metadata Commands (`c1.c`)

| Task                             | Status | Notes                               |
|----------------------------------|--------|-------------------------------------|
| Port `may_name` helper           | ☑      | Permission check for naming         |
| Port `v_name`                    | ☑      | Rename entities with length check   |
| Port `v_fullname`                | ☑      | Set player full_name                |
| Port `v_banner`                  | ☑      | Set banner text                     |
| Port `v_public`                  | ☑      | Set public_turn, award 100 gold     |
| Port `v_look`                    | ☑      | Basic show_loc wrapper              |
| Port `v_emote`                   | ☑      | Send formatted message to target    |
| Port `v_stop`                    | ☑      | No-op success command               |
| Unit tests for metadata commands | ☑      | Permissions, length limits, effects |

---

### 25.8 – Simple Economy Commands (`c2.c`)

| Task                            | Status | Notes                                             |
|---------------------------------|--------|---------------------------------------------------|
| Port `v_discard`                | ☑      | Drop item (uses how_many, drop_item)              |
| Port `drop_player`              | ☑      | Thin wrapper, skip shell calls                    |
| Port `v_quit`                   | ☑      | Quit command, calls drop_player                   |
| Unit tests for economy commands | ☑      | Invalid item, insufficient qty, permission errors |

---

### 25.9 – Deprecated Text-Report Commands (Stubs)

Per TODO.md guidance: create stubs with `// Deprecated:` comments that panic if called.

| Task                                                          | Status | Notes                              |
|---------------------------------------------------------------|--------|------------------------------------|
| Stub `v_split`                                                | ☑      | Already "no longer supported" in C |
| Stub `v_format`                                               | ☑      | Report format (legacy)             |
| Stub `v_notab`                                                | ☑      | No-TABs option (legacy)            |
| Stub `v_times`                                                | ☑      | Times subscription (legacy)        |
| Stub `open_times`, `times_masthead`, `close_times`            | ☑      | Times file I/O (legacy)            |
| Stub `v_rumor`, `v_press`                                     | ☑      | Times submissions (legacy)         |
| Stub `text_list_free`, `line_length_check`, `parse_text_list` | ☑      | Text list parsing (legacy)         |
| Stub `v_post`, `v_message`                                    | ☑      | In-game posting (legacy)           |
| Stub `v_tell`                                                 | ☑      | Already `#if 0` in C               |
| Unit tests verify stubs panic                                 | ☑      | Subset of stubs tested             |

---

## Out of Scope (Deferred to Sprint 26)

- `kill_char` and death/lifecycle helpers
- `sink_ship`, `building_collapses`, `add_structure_damage`
- `find_nearest_land`
- `v_explore`, `d_explore`
- WAIT/FLAG system
- Inventory transfer commands (GIVE/GET/PAY/ACCEPT/CLAIM)
- Noble formation (FORM, `form_new_noble`)
- Combat training commands
- Ferry/boarding commands
- `v_improve_opium`, `v_die`

---

## Files to Create/Modify

- `inventory.go` – Inventory & gold helpers
- `knowledge.go` – Bitset & knowledge helpers
- `time.go` or `format.go` – Time & numeric formatters
- `weights.go` – Weight & capacity system
- `region.go` – Region/visibility helpers (or add to `loc.go`)
- `cmd_meta.go` – Metadata commands (NAME, BANNER, LOOK, etc.)
- `cmd_economy.go` – Economy commands (DISCARD, QUIT)
- `cmd_deprecated.go` – Deprecated command stubs
- `*_test.go` – Unit tests for each module

---

## Technical Notes

### Deprecated Command Pattern

```go
// v_split is deprecated: legacy report splitting. Not used in Go/DB version.
// Deprecated: use web frontend for report viewing.
func (e *Engine) v_split(c *command) int {
panic("Deprecated: v_split (report splitting) is not supported in the Go/DB version")
}
```

### Inventory Helper Dependencies

Many inventory helpers depend on:

- `bx[]` entity accessors
- `item_*` metadata functions
- `ilist`/`plist` operations (already ported in Sprint 1-2)
- `province()`, `find_nearest_land()` for drop-at-sea logic

Port in order: pure queries → mutations → stack-level → payment/charge.
