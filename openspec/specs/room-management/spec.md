# room-management Specification

## Purpose
TBD - created by archiving change multi-room-tokens. Update Purpose after archive.
## Requirements
### Requirement: Server creates a room for each configured token
The server SHALL parse the `NOTIRC_TOKENS` environment variable as a comma-separated list of tokens and create an isolated chat room for each unique token. Duplicate tokens SHALL be silently deduplicated. If the list is empty or the variable is not set, the server SHALL exit with a fatal error.

#### Scenario: Multiple unique tokens
- **WHEN** `NOTIRC_TOKENS` is set to `"tok-abc,tok-xyz,tok-123"`
- **THEN** the server creates three rooms, one for each token, and starts listening

#### Scenario: Duplicate tokens
- **WHEN** `NOTIRC_TOKENS` is set to `"tok-abc,tok-abc,tok-xyz"`
- **THEN** the server creates two rooms (`tok-abc` and `tok-xyz`), silently ignoring the duplicate

#### Scenario: Empty token list
- **WHEN** `NOTIRC_TOKENS` is set to `""` or is not present
- **THEN** the server exits with a fatal error and does not start

### Requirement: Rooms are isolated from each other
Each room SHALL maintain its own nickname namespace, presence list, and message broadcast scope. Clients in different rooms SHALL NOT receive each other's messages, join/leave events, or presence information.

#### Scenario: Same nickname in different rooms
- **WHEN** client A connects with token `tok-abc` and joins as `alice`
- **AND** client B connects with token `tok-xyz` and joins as `alice`
- **THEN** both joins succeed independently

#### Scenario: Message isolation
- **WHEN** a client in room `tok-abc` sends a message
- **THEN** only clients connected with token `tok-abc` receive the broadcast
- **AND** clients connected with token `tok-xyz` do NOT receive the message

#### Scenario: Presence isolation
- **WHEN** a client in room `tok-abc` joins or leaves
- **THEN** presence events are broadcast only to clients in room `tok-abc`

### Requirement: Server rejects connections with unrecognized tokens
The server SHALL reject WebSocket connections whose `token` query parameter does not match any configured token, returning HTTP 401.

#### Scenario: Unrecognized token
- **WHEN** a client connects with `?token=not-configured`
- **THEN** the server returns HTTP 401 and does not upgrade to WebSocket

### Requirement: All rooms use the single WebSocket endpoint
The server SHALL route all room connections through the existing `/ws` endpoint using the `token` query parameter to determine the target room. No path-based room routing is required.

#### Scenario: Room routing by token
- **WHEN** a client connects to `/ws?token=tok-abc`
- **THEN** the server routes the connection to the room associated with `tok-abc`

### Requirement: Room identifier is included in log output
The server SHALL include the room's token (truncated to 8 characters) in all log lines emitted by room-scoped operations.

#### Scenario: Join event logged with room context
- **WHEN** a client joins a room
- **THEN** the log line includes a truncated room identifier (e.g., `room=tok-abc` for token `tok-abc...`)

