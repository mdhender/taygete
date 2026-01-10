# Loop Helpers Decision Record

## Background

The original C source uses `#define` macros in `src/loop.h` to create iterator patterns:

```c
loop_here(where, i) {
    // use i
}
next_here;
```

These macros:
1. Copy lists before iterating (safe deletion during iteration)
2. Use linked-list traversal via `kind_first()`/`kind_next()`
3. Include assertions for debugging
4. Handle resource cleanup with `ilist_reclaim()`

## Decision

Replace `loop_*`/`next_*` macro pairs with **slice-returning helper functions** plus normal `for range` loops.

## Rationale

1. **SQLite backend**: The original macros worked around pointer invalidation and manual memory management. With SQLite, iterating a snapshot slice of IDs is naturally safe.

2. **Type safety**: Macros lose types; Go functions with typed parameters and return values improve correctness.

3. **Idiomatic Go**: `for _, v := range slice` is idiomatic and readable.

4. **Cleanup is automatic**: Go's GC handles slice reclamation; no `ilist_reclaim()` needed.

5. **Semantic match**: The C macros already copied lists before iterating—returning a slice achieves the same snapshot semantics.

## Pattern Mappings

| C Macro                                      | Go Replacement                            |
|----------------------------------------------|-------------------------------------------|
| `loop_here(where, i)`                        | `e.UnitsHere(where) ([]int, error)`       |
| `loop_all_here(where, i)`                    | `e.AllHere(where) ([]int, error)`         |
| `loop_char_here(where, i)`                   | `e.CharsHere(where) ([]int, error)`       |
| `loop_stack(who, i)`                         | `e.Stack(who) ([]int, error)`             |
| `loop_units(pl, i)`                          | `e.UnitsOfPlayer(pl) ([]int, error)`      |
| `loop_inv(who, e)`                           | `e.InventoryOf(who) ([]item_ent, error)`  |
| `loop_char_skill(who, e)`                    | `e.CharSkills(who) ([]*skill_ent, error)` |
| `loop_trade(who, e)`                         | `e.TradesOf(who) ([]*trade, error)`       |
| `loop_kind(T_char, i)` / `loop_char(i)`      | `e.Characters() ([]int, error)`           |
| `loop_kind(T_player, i)` / `loop_player(i)`  | `e.Players() ([]int, error)`              |
| `loop_kind(T_loc, i)` / `loop_loc(i)`        | `e.Locations() ([]int, error)`            |
| `loop_subkind(sub_city, i)` / `loop_city(i)` | `e.Cities() ([]int, error)`               |
| `loop_province(i)`                           | `e.Provinces() ([]int, error)`            |
| `loop_boxes(i)`                              | `e.Boxes() ([]int, error)`                |
| `loop_known(kn, i)`                          | `sort.Ints(kn); for _, i := range kn`     |

## Example Translation

### C (original)

```c
loop_here(where, i) {
    delete_something(i);
}
next_here;
```

### Go (ported)

```go
ids, err := e.UnitsHere(where)
if err != nil {
    return err
}
for _, id := range ids {
    e.deleteSomething(id)
}
```

## When to Use Range-Over-Func Iterators

Only if profiling shows slice allocations are a hotspot for large entity scans. Otherwise, slice helpers are simpler and equally idiomatic.

## Implementation Notes

1. Helper functions should query DB → fill slice → close cursor → return slice
2. Use typed identifiers (`BoxID`, `PlayerID`) where beneficial
3. Return `([]T, error)` to handle query failures gracefully
4. For in-memory iteration (like `box_head` chains), implement as methods on Engine
