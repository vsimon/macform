# macform

`macform` is a Go CLI tool for declaratively managing macOS system settings via a YAML spec file.

@docs/spec.md

## Working in this codebase

> Think carefully and implement the most concise solution that changes as little code as possible. Follow existing patterns in the codebase.

- **Language**: Go only, no external scripting runtimes
- **Package layout**: `cmd/` for CLI commands, `internal/` for all business logic
- **Settings registry**: all settings defined in `internal/registry/` — adding a new setting requires a registry entry + an entry in `examples/macform.yaml`
- **Provider abstraction**: `internal/provider/` — new settings providers implement the `Provider` interface
- **Build tooling**: managed via `mise` — run `mise install` before building; use `mise build`
- **Worktree directory**: `.claude/worktrees` (project-local, hidden)
