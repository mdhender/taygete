# AGENTS.md

## Build Commands
- **Go CLI**: `go build ./cmd/xlat` (builds the translator tool)
- **C Game Engine**: `cd olympia && make` (builds 32-bit Olympia game server)
- **Clean**: `cd olympia && make clean`
- **Run tests**: `go test ./...` (no tests exist yet)
- **Single test**: `go test -run TestName ./path/to/package`

## Architecture
This is Olympia, a play-by-email strategy game. The codebase has two parts:
- **olympia/**: Legacy C game engine with combat, magic, NPCs, locations, items, skills
- **cmd/xlat/**: Go CLI tool (cobra) to translate Olympia data files to JSON

Key C headers: `olympia/oly.h` (main types/constants), `olympia/code.h`, `olympia/loc.h`

## Code Style
- **Go**: Use `log/slog` for logging, `cobra` for CLI, standard library preferred
- **C**: 32-bit compilation (`-m32`), K&R style, macros for accessors (see oly.h)
- **Naming**: Go uses camelCase; C uses snake_case with prefixes (sk_, item_, sub_)
- **Errors**: Go returns `error`; C uses return codes
