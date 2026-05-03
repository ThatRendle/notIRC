## Why

NotIRC currently supports a single global chat room shared by all clients. Running multiple workshops requires deploying a separate server instance per workshop. This change allows a single server to host multiple isolated chat rooms, each gated by its own token — one deployment, many workshops.

## What Changes

- **BREAKING**: Replace `NOTIRC_TOKEN` environment variable with `NOTIRC_TOKENS`, a comma-separated list of valid tokens
- The server creates an isolated chat room for each unique token in the list
- Clients connecting with different tokens are placed in separate rooms and cannot see each other's messages or presence
- Nickname uniqueness is scoped per-room (not globally)
- Duplicate tokens in the list are silently deduplicated
- The server refuses to start if `NOTIRC_TOKENS` is empty or not set
- Rename internal `Hub` type to `Room` for clarity
- Include room identifier in log output

## Capabilities

### New Capabilities
- `room-management`: Token-based room isolation — each token creates an independent chat room with its own nickname namespace, presence list, and message broadcast scope

### Modified Capabilities
- `security`: Environment variable changes from `NOTIRC_TOKEN` (single value) to `NOTIRC_TOKENS` (comma-separated list)
- `messaging`: Channel scope changes from a single shared global channel to isolated per-room channels

## Impact

- `main.go` — startup config, mux wiring
- `server.go` — rename `Hub` to `Room`, add `RoomManager`
- `client.go` — no functional changes (still references `hub`/room field)
- `message.go` — no changes (wire protocol unchanged)
- `integration_test.go` — test setup adapts to multi-token config
- `server_test.go` — rename `Hub` → `Room` in test code
- `README.md` and `docs/BUILD_A_CLIENT.md` — update token env var documentation
- `docker-compose.yml` — update environment variable name
