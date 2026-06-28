# AGENTS.md — new-repo

## Project

`new-repo` is a Go CLI that scaffolds an opencode-ready project directory:
creates `<name>/`, `.gitignore`, `SKILLS.md`, copies skills into
`.opencode/skills/`, and runs `git init`.

## Current state

Implemented and working. See `TDD-INIT-GO.md` for design decisions.

## Key decisions (from TDD)

- **Stdlib only** — no external Go deps (`os`, `os/exec`, `flag`, `embed`, `path/filepath`)
- **Entrypoint** — `cmd/new-repo/main.go`
- **Packages** — `internal/{scaffold,skill,git,config}/`
- **Skills origin** — embedded in binary via `//go:embed SKILLS.md`
- **Skills destination** — `.opencode/skills/<skill-name>/` inside target dir
- **No `--no-git` flag**

## Commands

```bash
go build -o new-repo ./cmd/new-repo
go test ./...
go vet ./...
```

## Testing

- Unit tests use standard `go test` with table-driven tests
- Integration tests create a temp dir via `t.TempDir()`, run the scaffold,
  and verify files + `git rev-parse`
- Git must be installed for integration tests

## Requirements

- Go 1.22+
- Git (runtime dependency for `git init`)
