# Sprint 9: loc.c (Ownership & Spatial Functions)

## Objective
Port the core ownership and spatial navigation functions from `src/loc.c` to Go. These functions determine where entities are in the world hierarchy (region → province → subloc → building).

## Prerequisites (Already Complete)
- ✅ `kind()`, `subkind()`, `valid_box()` in `accessor.go`
- ✅ `loc()` function in `accessor.go` (line 255)
- ✅ `safe_haven()` function in `accessor.go` (line 642)
- ✅ `rp_loc_info()`, `p_loc_info()` accessors
- ✅ `LOC_region`, `LOC_province`, `LOC_subloc`, `LOC_build` constants in `glob.go`
- ✅ `ilist` operations in `z.go`

## Dependencies to Port First
Before the main `loc.c` functions, we need:

### Task 1: Port `loc_depth()` from `src/u.c`
The `loc_depth()` function (lines 150-220 in u.c) determines location depth based on subkind. This is used by `region()`, `province()`, and `viewloc()`.

**Location**: Currently stubbed in `glob.go:578` - needs proper implementation.

**Subkinds required** (verify in `const.go`):
- `sub_region`, `sub_ocean`, `sub_forest`, `sub_plain`, `sub_mountain`, `sub_desert`, `sub_swamp`
- `sub_under`, `sub_cloud`, `sub_tunnel`, `sub_chamber`
- `sub_island`, `sub_stone_cir`, `sub_mallorn_grove`, `sub_bog`, `sub_cave`, `sub_city`
- `sub_lair`, `sub_graveyard`, `sub_ruins`, `sub_battlefield`, `sub_ench_forest`
- `sub_rocky_hill`, `sub_tree_circle`, `sub_pits`, `sub_pasture`, `sub_oasis`
- `sub_yew_grove`, `sub_sand_pit`, `sub_sacred_grove`, `sub_poppy_field`
- `sub_faery_hill`, `sub_hades_pit`
- `sub_temple`, `sub_galley`, `sub_roundship`, `sub_castle`, etc.

---

## Functions to Port (loc.go)

### Task 2: Core Navigation Functions

| C Function | Go Function | Description |
|------------|-------------|-------------|
| `region(who)` | `region(who int) int` | Return ultimate region containing entity |
| `province(who)` | `province(who int) int` | Return ultimate province containing entity |
| `subloc(who)` | `subloc(who int) int` | Return immediate location (ignoring stacks) |
| `viewloc(who)` | `viewloc(who int) int` | Return location for visibility purposes |

### Task 3: Safety & Containment Functions

| C Function | Go Function | Description |
|------------|-------------|-------------|
| `in_safe_now(who)` | `in_safe_now(who int) bool` | Check if entity is in a safe haven |
| `somewhere_inside(a, b)` | `somewhere_inside(a, b int) bool` | Check if b is nested inside a |

### Task 4: Building Owner (deferred to Sprint 10)
`building_owner()` and related here-list functions are Sprint 10 scope.

---

## Implementation Plan

### File: `loc.go`

```go
// loc.go - Spatial model functions ported from src/loc.c

package taygete

// loc_depth returns the depth level of a location based on its subkind.
// Returns LOC_region, LOC_province, LOC_subloc, or LOC_build.
func loc_depth(n int) int {
    switch subkind(n) {
    case sub_region:
        return LOC_region
    case sub_ocean, sub_forest, sub_plain, sub_mountain, sub_desert,
         sub_swamp, sub_under, sub_cloud, sub_tunnel, sub_chamber:
        return LOC_province
    case sub_island, sub_stone_cir, sub_mallorn_grove, sub_bog, sub_cave,
         sub_city, sub_lair, sub_graveyard, sub_ruins, sub_battlefield,
         sub_ench_forest, sub_rocky_hill, sub_tree_circle, sub_pits,
         sub_pasture, sub_oasis, sub_yew_grove, sub_sand_pit,
         sub_sacred_grove, sub_poppy_field, sub_faery_hill, sub_hades_pit:
        return LOC_subloc
    case sub_temple, sub_galley, sub_roundship, sub_castle,
         sub_galley_notdone, sub_roundship_notdone, sub_ghost_ship,
         sub_temple_notdone, sub_inn, sub_inn_notdone, sub_castle_notdone,
         sub_mine, sub_mine_notdone, sub_mine_collapsed, sub_tower,
         sub_tower_notdone, sub_sewer:
        return LOC_build
    default:
        return 0
    }
}

// region returns the ultimate region containing who.
func region(who int) int {
    for who > 0 && (kind(who) != T_loc || loc_depth(who) != LOC_region) {
        who = loc(who)
    }
    return who
}

// province returns the ultimate province containing who.
func province(who int) int {
    for who > 0 && (kind(who) != T_loc || loc_depth(who) != LOC_province) {
        who = loc(who)
    }
    return who
}

// subloc returns the immediate location (T_loc or T_ship) containing who.
func subloc(who int) int {
    for {
        who = loc(who)
        if who <= 0 || kind(who) == T_loc || kind(who) == T_ship {
            break
        }
    }
    return who
}

// viewloc returns the location to use for visibility calculations.
func viewloc(who int) int {
    for who > 0 &&
        loc_depth(who) != LOC_province &&
        subkind(who) != sub_city &&
        subkind(who) != sub_graveyard &&
        subkind(who) != sub_sewer &&
        subkind(who) != sub_faery_hill {
        who = loc(who)
    }
    return who
}

// in_safe_now returns true if who is anywhere inside a safe haven.
func in_safe_now(who int) bool {
    for {
        if safe_haven(who) != 0 {
            return true
        }
        who = loc(who)
        if who <= 0 {
            break
        }
    }
    return false
}

// somewhere_inside returns true if b is nested somewhere inside a.
func somewhere_inside(a, b int) bool {
    if a == b {
        return false
    }
    for b > 0 {
        b = loc(b)
        if a == b {
            return true
        }
    }
    return false
}
```

---

## Testing Plan

### File: `loc_test.go`

Create tests with a mock world hierarchy:

```
Region (aa00)
  └── Province (aa01) 
        ├── City (1001)
        │     └── Castle (2001)
        │           └── Character (3001)
        └── Character (3002)
```

**Test Cases:**
1. `TestRegion` - verify `region()` walks up to region level
2. `TestProvince` - verify `province()` stops at province level
3. `TestSubloc` - verify `subloc()` finds immediate T_loc/T_ship
4. `TestViewloc` - verify special handling for city/graveyard/sewer/faery_hill
5. `TestInSafeNow` - verify safe haven detection through nesting
6. `TestSomewhereInside` - verify containment checks

---

## Verification Steps

1. Run `go build .` - ensure no compile errors
2. Run `go test -run TestLoc ./...` - all loc tests pass
3. Run `go test ./...` - no regressions in existing tests

---

## Notes

- The `loc_depth()` function logically belongs in `loc.go` even though it's defined in `src/u.c`
- Update the stub in `glob.go:578` to call the real `loc_depth()` after implementation
- Sprint 10 will add here-list functions (`all_here`, `subloc_here`, `count_loc_structures`)
