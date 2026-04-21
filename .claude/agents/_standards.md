# Coding Standards

## Language and framework

**Language:** Go
**WebSocket library:** `github.com/coder/websocket`
**Logging:** `log/slog` (stdlib, structured JSON output)

## Conventions

**Interfaces:** Defined at the usage site — in the package that consumes them, not alongside the implementing type.

**Error handling:** Idiomatic Go `(value, error)` returns throughout. No panics in library code; panics acceptable in `main()` for unrecoverable startup errors.

**Logging:** Structured via `log/slog`. Log: connections accepted/rejected (with reason), nickname joins and leaves, and errors. Do not log message content.

**Dependencies:** Pragmatic — use well-maintained libraries freely. No artificial stdlib-only constraint.

**Documentation:** Inline comments for non-obvious logic only. Comment the *why*, not the *what*. No comment blocks for self-evident code.

**Testing:** Pragmatic — TDD for complex logic (Hub state management, message validation), test-after for straightforward code. Target 80% coverage on business logic. Do not test framework boilerplate.

**Patterns:** Mixed / pragmatic. Use whatever fits the problem. No enforced paradigm.

**Naming:** Standard Go conventions — `CamelCase` for exported, `camelCase` for unexported, short variable names in small scopes (`c` for client, `h` for hub), descriptive names for package-level declarations.

## Tooling

*No linting or formatting config detected (greenfield). Assume `gofmt` and `go vet` as baseline.*
