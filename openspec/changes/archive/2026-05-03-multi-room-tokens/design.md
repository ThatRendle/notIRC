## Context

NotIRC is a WebSocket-based group chat server used for workshops. Currently it accepts a single `NOTIRC_TOKEN` environment variable and places all clients into one global chat room (the `Hub`). Running multiple concurrent workshops requires deploying separate server instances.

The server is a single-binary Go application with no external dependencies beyond `github.com/coder/websocket`. It has no database, no persistent state, and no authentication beyond a shared token.

## Goals / Non-Goals

**Goals:**
- Support multiple isolated chat rooms in a single server process
- Route clients to rooms based on the token they present at connection time
- Keep the WebSocket wire protocol unchanged — clients built for the current protocol work without modification
- Minimal operational overhead: configure rooms with a single environment variable
- Backward-incompatible: replace `NOTIRC_TOKEN` with `NOTIRC_TOKENS` (no dual-mode support to keep code simple)

**Non-Goals:**
- Human-readable room names or labels (tokens are opaque identifiers)
- Dynamic room creation at runtime (rooms are pre-created from config at startup)
- Room lifecycle management (cleanup of empty rooms, room deletion)
- Admin endpoints for inspecting active rooms
- Per-room message history or persistence
- Token rotation or revocation

## Decisions

### 1. RoomManager pattern

A `RoomManager` struct holds the token-to-room mapping. At startup, `NOTIRC_TOKENS` is parsed, deduplicated, and a `Room` is created for each unique token.

```
RoomManager {
    rooms:  map[string]*Room   // token → Room
}

Room (formerly Hub) {
    mu:      sync.RWMutex
    clients: map[string]*Client   // nick → Client
}
```

**Alternatives considered:**
- *Lazy room creation on first connect* — adds race conditions and complexity for no benefit since tokens are known at startup.
- *Single map of token→Room without a manager struct* — works but makes it harder to add future room-level operations (stats, listing).

### 2. Environment variable format

`NOTIRC_TOKENS` is a comma-separated list of tokens. No key-value pairing, no labels.

```
NOTIRC_TOKENS=tok-abc,tok-xyz,tok-123
```

**Alternatives considered:**
- *`label=token` format* — adds room naming but the user explicitly chose opaque tokens for simplicity.
- *JSON/YAML config file* — adds a file dependency and deployment complexity.
- *Multiple env vars (`NOTIRC_TOKEN_1`, `NOTIRC_TOKEN_2`)* — fragile, hard to enumerate.

### 3. Empty config is fatal

If `NOTIRC_TOKENS` is empty or not set, the server logs a fatal error and exits. This matches the existing behavior for `NOTIRC_TOKEN` and avoids silently running a server with no rooms.

### 4. No protocol changes

The WebSocket message types (`join`, `join_ok`, `message`, `user_joined`, `user_left`) are unchanged. The token-to-room routing is transparent to clients. A client that works against the current single-room server works identically against a multi-room server — it just won't see clients in other rooms.

### 5. Hub → Room rename

The existing `Hub` struct and all references are renamed to `Room`. This is a pure rename with no behavioral changes. The `Room` struct's methods (`join`, `leave`, `broadcast`) are unchanged.

### 6. Logging room context

All log lines emitted from `Room` methods include the room's token (or a prefix) so operators can distinguish which room a log line belongs to. The token is truncated to first 8 characters in logs to avoid leaking full tokens into log output.

## Risks / Trade-offs

- **Duplicate token config** — silently deduplicated. An operator who accidentally types the same token twice gets one room, not an error. Accepted risk for simplicity.
- **No room-level metrics** — operators can't query which rooms are active or how many clients per room without parsing logs. Acceptable for a workshop tool; can be added later.
- **Token in URL** — tokens remain visible in WebSocket URLs (query parameter). This is the existing pattern and is not changed by this design.
