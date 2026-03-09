# AGENTS Guide for `searxng-cli`
This file is the operating guide for agentic coding assistants in this repository.
Follow it as the source of truth for local commands and coding conventions.

## 1) Project Snapshot
- Language: Go (`go 1.26.1` in `go.mod`)
- Module: `searxng-cli`
- CLI framework: Cobra (`github.com/spf13/cobra`)
- Entry point: `main.go`
- CLI command wiring: `cmd/`
- Core domain logic: `internal/`
- Key domains:
  - `internal/search`: SearXNG query and decode
  - `internal/read`: web fetch, robots policy, extraction, sanitize
  - `internal/config`: config load + precedence resolution
  - `internal/auth`: auth header injection and base64 decode
  - `internal/apperr`: error taxonomy, classification, stderr rendering
  - `internal/render`: Markdown table rendering

## 2) Build / Run / Lint / Test Commands

### Setup
```bash
go mod download
```

### Build
```bash
go build ./...
go build -o searxng-cli .
```

### Run locally
```bash
go run . --help
go run . search --help
go run . read --help
go run . config --help
```

### Format + static checks ("lint" baseline)
```bash
gofmt -w $(rg --files -g '*.go')
go vet ./...
```

### Full test suite
```bash
go test ./...
```

### Package-level tests
```bash
go test ./cmd -v
go test ./internal/search -v
go test ./internal/read -v
go test ./internal/config -v
go test ./internal/auth -v
go test ./internal/apperr -v
go test ./internal/render -v
```

### Single test (important)
Use package + anchored `-run` regex:
```bash
go test ./cmd -run '^TestReadHelp$' -v
go test ./cmd -run '^TestHelpMentionsScopeForRootAndSearch$' -v
go test ./internal/search -run '^TestFetch_AppliesAPIKeyAuth$' -v
go test ./internal/read -run '^TestFetchHTML$' -v
go test ./internal/apperr -run '^TestClassifyVariants$' -v
```

### Subset tests by pattern
```bash
go test ./cmd -run 'TestRead|TestHelp' -v
go test ./internal/read -run 'TestFetch|TestExtract|TestSanitize' -v
```

### Optional race check (slower)
```bash
go test -race ./...
```

## 3) Suggested Agent Workflow
1. Read related files first (`cmd/*`, then matching `internal/*`).
2. Keep changes minimal and scoped to the requested behavior.
3. Update or add tests for behavior changes in the same package.
4. Run `gofmt`, then `go test` for affected package(s), then `go test ./...`.
5. If CLI help/output behavior changes, update tests and `SPEC.md`.

## 4) Code Style and Architecture Rules

### Imports and formatting
- Always run `gofmt` before finishing.
- Keep imports in standard Go grouping/order (stdlib, module-local, third-party).
- Do not hand-align or manually format Go code.

### Package boundaries
- `cmd/` contains Cobra command definitions and flag plumbing.
- `internal/*` contains reusable business logic.
- Avoid placing business rules directly in Cobra handlers.

### Types and configuration
- Prefer typed structs over generic maps for options/config.
- Keep config precedence stable: `flags > env > config > defaults`.
- Extend existing config structs instead of introducing parallel config paths.

### Naming conventions
- Exported identifiers: `PascalCase`.
- Internal identifiers: `camelCase`.
- Prefer descriptive names over abbreviations; keep acronyms readable (`URL`, `HTTP`, `APIKey`).
- Test names should describe behavior, not implementation detail.

### Error handling
- Never swallow errors silently.
- Wrap lower-level errors with context (`%w`) where useful.
- For user-facing failures, return typed/classifiable errors through `internal/apperr`.
- Keep `code`, `retryable`, `message`, and `hint` semantics stable.
- Never include secrets (API keys, credentials) in error text or metadata.

### CLI output contract
- Normal results must go to `stdout`.
- Failures must render structured key/value output to `stderr` (see `cmd/root.go`, `internal/apperr`).
- Keep help text aligned with actual command behavior.
- If output schema changes, update tests and `SPEC.md` together.

### Testing conventions
- Prefer table-driven tests for pure logic.
- Use `httptest` for HTTP/network behavior.
- Test both success and failure paths.
- For error taxonomy changes, add/adjust classifier tests in `internal/apperr/*_test.go`.
- For `read` flows, assert sanitized output quality (no noisy HTML leftovers).

### Dependency policy
- Default to Go stdlib and existing dependencies first.
- Add a dependency only when benefit is clear and test-backed.
- Avoid heavyweight/runtime browser dependencies unless explicitly requested.

## 5) Spec and Documentation Sync
When behavior changes, update all relevant artifacts:
- Command help text (`cmd/*.go`)
- Tests asserting behavior/help/output
- `SPEC.md` contract descriptions
- This `AGENTS.md` when workflows/conventions change

## 6) Cursor / Copilot Rules Integration
Checked in this repository:
- `.cursor/rules/`: not present
- `.cursorrules`: not present
- `.github/copilot-instructions.md`: not present

If any of these files are added later, treat them as repository-level instructions and merge their constraints into this guide.

## 7) Git and Safety Rules for Agents
- Do not run destructive git commands unless explicitly requested.
- Never revert unrelated local changes in a dirty worktree.
- Keep commits focused and small.
- Prefer clear commit messages that explain intent and behavior impact.

## 8) Release Notes
- Release workflow: `.github/workflows/release.yml`
- Trigger: pushing a `v*` tag
- Current target assets include Linux (`amd64`, `arm64`) and macOS (`arm64`) binaries plus checksum file.

Release example:
```bash
git tag v0.1.0
git push origin v0.1.0
```
