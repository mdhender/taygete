# AGENTS.md

## Build Commands
- **Build**: `go build .` (builds the main package)
- **Run tests**: `go test ./...`
- **Single test**: `go test -run TestName .`

## Architecture
Olympia is a play-by-email strategy game being ported from 32-bit C to Go.

- **src/**: Original C source (read-only reference, do not modify)
- **Root package**: Go port (single package, refactor after port complete)
- **Storage**: SQLite3 replaces legacy flat files
- **Frontend**: Next.js + Tailwind "Oatmeal" UI Kit (future)

## Porting Approach
- Port C global values to Engine globals
- Port C files to Go with matching names: `src/rnd.c` → `rnd.go`
- Port C functions to Go functions with matching names
  - Keep C function signatures as close as possible to the original
  - In ported functions, forward to the Go Engine implementation
  - Use Go idioms where appropriate
- Port C structs to Go structs with matching names
- Port C enums to Go enums with matching names
- Write unit tests as we port each file
- Use SQLite3 for data persistence (no flat file parsing)
- Web frontend replaces email reports

Key C headers for reference: `src/oly.h`, `src/code.h`, `src/loc.h`

## Code Style
- **Go**: Use `log/slog` for logging, standard library preferred
- **Errors**: Go returns `error`
- **Testing**: Write internal tests alongside each ported file
- **Naming**: Translate C snake_case names:
  - Go functions and variables mirror C names
  - Go enum type names mirror C names
  - Go struct type names mirror C names
  - New Go types, functions, and variables use camelCase
