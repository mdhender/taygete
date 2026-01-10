# Sprint 1: z.c Core Implementation Plan

## Overview
Port the core utilities from `src/z.c` (~1k LOC) to Go. Since we're using SQLite for persistence instead of flat files, several file I/O functions will be deprecated stubs.

---

## Functions to Port

### 1. Memory Helpers (SKIP - Go has GC)
These C functions are unnecessary in Go:

| C Function | Go Status | Notes |
|------------|-----------|-------|
| `my_malloc` | Skip | Go has garbage collection |
| `my_realloc` | Skip | Go has garbage collection |
| `my_free` | Skip | Go has garbage collection |
| `str_save` | Skip | Go strings are immutable, use `string` type |

---

### 2. ilist (Reallocating Integer Array)
**File**: `z.go`

Go equivalent: Use `[]int` slices. Create wrapper type for API compatibility.

| C Function | Go Function | Signature | Notes |
|------------|-------------|-----------|-------|
| `ilist_len` | `IListLen` | `func IListLen(l []int) int` | Returns `len(l)` |
| `ilist_append` | `IListAppend` | `func IListAppend(l *[]int, n int)` | `*l = append(*l, n)` |
| `ilist_prepend` | `IListPrepend` | `func IListPrepend(l *[]int, n int)` | Prepend to slice |
| `ilist_delete` | `IListDelete` | `func IListDelete(l *[]int, i int)` | Remove element at index |
| `ilist_clear` | `IListClear` | `func IListClear(l *[]int)` | `*l = (*l)[:0]` |
| `ilist_reclaim` | `IListReclaim` | `func IListReclaim(l *[]int)` | `*l = nil` |
| `ilist_lookup` | `IListLookup` | `func IListLookup(l []int, n int) int` | Find index or -1 |
| `ilist_rem_value` | `IListRemValue` | `func IListRemValue(l *[]int, n int)` | Remove all occurrences |
| `ilist_rem_value_uniq` | `IListRemValueUniq` | `func IListRemValueUniq(l *[]int, n int)` | Remove first occurrence |
| `ilist_copy` | `IListCopy` | `func IListCopy(l []int) []int` | Return copy of slice |
| `ilist_scramble` | `IListScramble` | `func IListScramble(l []int)` | Fisher-Yates shuffle |
| `ilist_insert` | `IListInsert` | `func IListInsert(l *[]int, pos, n int)` | Insert at position |

---

### 3. plist (Reallocating Pointer Array)
**File**: `z.go`

Go equivalent: Use `[]any` slices.

| C Function | Go Function | Signature | Notes |
|------------|-------------|-----------|-------|
| `plist_len` | `PListLen` | `func PListLen(l []any) int` | Returns `len(l)` |
| `plist_append` | `PListAppend` | `func PListAppend(l *[]any, n any)` | Append |
| `plist_prepend` | `PListPrepend` | `func PListPrepend(l *[]any, n any)` | Prepend |
| `plist_delete` | `PListDelete` | `func PListDelete(l *[]any, i int)` | Delete at index |
| `plist_clear` | `PListClear` | `func PListClear(l *[]any)` | Clear |
| `plist_reclaim` | `PListReclaim` | `func PListReclaim(l *[]any)` | Reclaim |
| `plist_lookup` | `PListLookup` | `func PListLookup(l []any, n any) int` | Find index |
| `plist_rem_value` | `PListRemValue` | `func PListRemValue(l *[]any, n any)` | Remove all |
| `plist_rem_value_uniq` | `PListRemValueUniq` | `func PListRemValueUniq(l *[]any, n any)` | Remove first |
| `plist_copy` | `PListCopy` | `func PListCopy(l []any) []any` | Copy |
| `plist_scramble` | `PListScramble` | `func PListScramble(l []any)` | Shuffle |
| `plist_insert` | `PListInsert` | `func PListInsert(l *[]any, pos int, n any)` | Insert |

---

### 4. String Comparators
**File**: `z.go`

| C Function | Go Function | Signature | Notes |
|------------|-------------|-----------|-------|
| `i_strcmp` | `i_strcmp` | `func i_strcmp(s, t string) int` | Case-insensitive compare |
| `i_strncmp` | `i_strncmp` | `func i_strncmp(s, t string, n int) int` | Case-insensitive prefix compare |
| `fuzzy_strcmp` | `fuzzy_strcmp` | `func fuzzy_strcmp(one, two string) bool` | Fuzzy match (typo-tolerant) |
| `lcase` | `lcase` | `func lcase(s string) string` | Lowercase (use `strings.ToLower`) |

Helper functions (internal):
- `fuzzy_transpose` - Check transposed chars
- `fuzzy_one_less` - Check one char missing
- `fuzzy_one_extra` - Check one char extra  
- `fuzzy_one_bad` - Check one wrong char

---

### 5. Assert
**File**: `z.go`

| C Function | Go Function | Signature | Notes |
|------------|-------------|-----------|-------|
| `asfail` | `asfail` | `func asfail(file string, line int, cond string)` | Panic with formatted message |

In Go, use `panic()` directly or create a helper that logs and panics.

---

### 6. Character Classification / Case Folding (Sprint 2)
These are marked for Sprint 2 but included here for reference:

| C Item | Go Status | Notes |
|--------|-----------|-------|
| `lower_array` | Use `strings.ToLower` | Go handles Unicode |
| `init_lower` | Skip | Not needed in Go |
| `iswhite` | `iswhite` | `func iswhite(c byte) bool` |
| `isalpha` | Use `unicode.IsLetter` | Standard library |
| `isdigit` | Use `unicode.IsDigit` | Standard library |

---

### 7. File I/O Functions (DEPRECATED)
**File**: `z.go`

These functions are replaced by SQLite persistence. Create deprecated stubs that panic if called.

| C Function | Go Function | Signature | Notes |
|------------|-------------|-----------|-------|
| `readfile` | `readfile` | `func readfile(path string) int` | **DEPRECATED**: panic |
| `readlin` | `readlin` | `func readlin() string` | **DEPRECATED**: panic |
| `readlin_ew` | `readlin_ew` | `func readlin_ew() string` | **DEPRECATED**: panic |
| `closefile` | `closefile` | `func closefile(path string)` | **DEPRECATED**: panic |
| `getlin` | `getlin` | `func getlin(fp *os.File) string` | **DEPRECATED**: panic |
| `getlin_ew` | `getlin_ew` | `func getlin_ew(fp *os.File) string` | **DEPRECATED**: panic |
| `copy_fp` | `copy_fp` | `func copy_fp(a, b *os.File)` | **DEPRECATED**: panic |

---

### 8. String Whitespace Trimming (DEPRECATED)
**File**: `z.go`

These are primarily used for file parsing, which is replaced by SQLite.

| C Function | Go Function | Signature | Notes |
|------------|-------------|-----------|-------|
| `eat_leading_trailing_whitespace` | `eat_leading_trailing_whitespace` | `func eat_leading_trailing_whitespace(s string) string` | **DEPRECATED**: panic |

---

### 9. Utility Functions
**File**: `z.go`

| C Function | Go Function | Signature | Notes |
|------------|-------------|-----------|-------|
| `int_comp` | `int_comp` | `func int_comp(a, b int) int` | Compare integers for sorting |

---

## Implementation Order

1. **z.go** - Create file with package header
2. **ilist functions** - Implement all IList* functions
3. **plist functions** - Implement all PList* functions  
4. **String comparators** - i_strcmp, i_strncmp, fuzzy_strcmp
5. **Assert helper** - asfail
6. **Character classification** - iswhite
7. **Deprecated stubs** - All file I/O functions with panic
8. **Tests** - z_test.go

---

## Test Plan

### z_test.go

```go
// Tests for Sprint 1 z.c port

func TestIListAppendPrepend(t *testing.T)
func TestIListDelete(t *testing.T)
func TestIListLookup(t *testing.T)
func TestIListRemValue(t *testing.T)
func TestIListCopy(t *testing.T)
func TestIListScramble(t *testing.T)

func TestPListAppendPrepend(t *testing.T)
func TestPListDelete(t *testing.T)
func TestPListLookup(t *testing.T)

func Test_i_strcmp(t *testing.T)
func Test_i_strncmp(t *testing.T)
func Test_fuzzy_strcmp(t *testing.T)
```

---

## Deprecated Function Template

```go
// Deprecated: readfile is not needed; game state is stored in SQLite.
// This function will panic if called.
func readfile(path string) int {
	panic("readfile: deprecated - use SQLite for game state persistence")
}
```

---

## File Structure After Sprint 1

```
taygete/
├── z.go           # Core utilities from z.c
├── z_test.go      # Tests for z.go
├── rnd.go         # Already ported (Sprint 3)
├── rnd_test.go    # Already exists
├── types.go       # Type definitions
├── ...
```

---

## Acceptance Criteria

- [x] All ilist functions implemented with Go slices
- [x] All plist functions implemented with Go slices
- [x] i_strcmp, i_strncmp return correct comparison results
- [x] fuzzy_strcmp correctly identifies typos (transpose, extra/missing char, wrong char)
- [x] Deprecated functions have proper `// Deprecated:` comments and panic
- [x] All tests pass: `go test ./...`
- [x] Build succeeds: `go build .`

---

## Notes

- Go slices replace the C ilist/plist implementation naturally
- The C ilist stores length in `array[-2]` and capacity in `array[-1]`; Go slices handle this automatically
- `ilist_scramble` and `plist_scramble` use the existing `rnd()` function from rnd.go
- String functions use Go's `strings` package where appropriate
- Deprecated functions must panic to catch any accidental usage during the port
