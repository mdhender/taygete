# Refactoring ilist to Go Slices

**Date:** Sprint 18 Decision  
**Status:** Approved for Partial Implementation

## Background

The original C codebase uses `ilist` - a pointer-to-int type with hidden length/capacity stored at positions `[-2]` and `[-1]`:

```c
typedef int *ilist;  // z.h
// ilist[0] = length, ilist[1] = capacity, user data starts at &ilist[2]
```

The Go port initially preserved this as `type ilist *int` but has already migrated to slice-based helpers in `z.go`:

```go
func IListLen(l []int) int
func IListAppend(l *[]int, n int)
func IListPrepend(l *[]int, n int)
// ... etc
```

This creates a mismatch: the `IList*` helpers operate on `[]int`, but some struct fields still use the legacy `ilist` type.

## The Problem

`entity_skill` has two `ilist` fields that block full skill tree validation:

```go
type entity_skill struct {
    offered  ilist  // skills learnable after this one
    research ilist  // skills researchable with this one
    // ...
}
```

These fields cannot be used with the `IList*` helpers without type conversion gymnastics, and the C-style pointer encoding was never implemented in Go.

## Scope Analysis

### Narrow Scope: entity_skill Only

**Files affected:** 1 (`types.go`)  
**Effort:** S (<1 hour)  
**Risk:** Very Low

The `offered` and `research` fields are currently **dead data** - defined but never set or read in Go code. Changing them to `[]int` has no ripple effects.

### Full Scope: All ilist Fields

Many structs still use `ilist`:

- `entity_player.units`
- `entity_subloc.teaches`
- `item_magic.may_use`, `may_study`
- `entity_misc.garr_watch`, `garr_host`
- `entity_loc.prov_dest`
- And others...

**Effort:** L-XL (1-2 days)  
**Risk:** Moderate - touches many structs and requires careful testing

## Decision

### Sprint 18: Narrow Refactor (Approved)

Change `entity_skill.offered` and `entity_skill.research` from `ilist` to `[]int`:

```go
type entity_skill struct {
    time_to_learn  int
    required_skill int
    np_req         int
    offered        []int  // skills learnable after this one (was: ilist)
    research       []int  // skills researchable with this one (was: ilist)
    // ...
}
```

**Rationale:**
- Very small, localized change
- Aligns with established Go-side pattern using `[]int` + `IList*` helpers
- Unblocks skill tree validation tests
- Does not force full migration

### Future Sprint: Full ilist Migration (Deferred)

Defer the global `ilist` → `[]int` migration to a dedicated sprint when:
- Porting logic that manipulates player units, subloc teaches, item magic lists, etc.
- The number of workaround maps becomes unmanageable

The full migration will:
1. Change `type ilist` definition or remove it entirely
2. Update all struct fields using `ilist` to `[]int`
3. Eliminate ad-hoc maps by backing them onto slice fields in entities

## Implementation Steps (Sprint 18)

1. **Update types.go:**
   - Change `offered ilist` to `offered []int`
   - Change `research ilist` to `research []int`

2. **Optional convenience helpers:**
   ```go
   func addOfferedSkill(sk, child int) {
       s := p_skill(sk)
       IListAppend(&s.offered, child)
   }
   ```

3. **Future: Wire up from SQLite:**
   - Add `skill_offered` and `skill_research` tables (if needed)
   - Add loaders in `load.go`
   - Extend `checkSkills` validation

## Notes

- The `type ilist *int` alias remains for other unported fields
- Add comments in `types.go` noting that `entity_skill` has diverged from C's ilist representation by design
- When porting `check.c` skill graph logic, translate ilist operations to `IList*` slice operations
