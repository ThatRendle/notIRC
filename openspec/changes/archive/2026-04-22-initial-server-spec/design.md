# Design: notIRC Backend Server

## Technology choices

### Language: Go

Go compiles to a single static binary, has excellent built-in concurrency primitives for managing concurrent WebSocket connections, and produces a minimal container image. Its standard library and `github.com/coder/websocket` cover everything needed without pulling in a large framework.

**Rejected:** Node.js (higher memory baseline, runtime dependency in container), Python (slower, GIL complicates concurrent connection handling), Rust (correct choice for a production system at scale; over-engineered for this use case).

### State: in-memory

Nickname and connection state is held in a `sync.RWMutex`-protected map. No external store.

**Rejected:** Redis — adds cost (Railway bills per service), external dependency, and operational complexity that is not justified by the scale (20–40 connections, single instance, ephemeral session). If notIRC grows beyond the workshop context, Redis pub/sub for multi-instance coordination would be the right next step.

### Message format: JSON

All messages — both client-to-server and server-to-client — are JSON objects with a `type` string discriminator. JSON is consumable from every language and platform workshop participants might use.

### Authentication: query parameter token

The WebSocket upgrade URL carries the API token as a query parameter:

```
wss://<host>/ws?token=<token>
```

The server rejects the upgrade with HTTP 401 if the token is absent or incorrect. This approach works universally — browser WebSocket APIs do not support custom headers, so a query parameter is the only method that works across JavaScript, Python, Go, Swift, Kotlin, Rust, and .NET without special handling.

The token is a single shared secret configured via an environment variable (`NOTIRC_TOKEN`). Token distribution to workshop participants is out of scope for the server.

---

## Message protocol

### Client → Server

#### `join`

Sent as the first message after a successful WebSocket upgrade. The token has already been validated at the HTTP layer; this message registers the client's chosen nickname.

```json
{ "type": "join", "nick": "alice" }
```

#### `message`

Sent by a joined client to broadcast a message to the channel.

```json
{ "type": "message", "text": "hello everyone" }
```

---

### Server → Client

#### `join_ok`

Sent to the connecting client when their nickname is accepted. Includes the current user list at the moment of acceptance.

```json
{ "type": "join_ok", "users": ["bob", "carol"] }
```

#### `join_error`

Sent to the connecting client when their nickname is rejected.

```json
{ "type": "join_error", "reason": "nick_taken" }
```

After receiving `join_error`, the client remains connected and may send another `join` message with a different nickname.

#### `user_joined`

Broadcast to all already-connected clients when a new client's nickname is accepted.

```json
{ "type": "user_joined", "nick": "alice" }
```

#### `user_left`

Broadcast to all remaining clients when a client disconnects.

```json
{ "type": "user_left", "nick": "alice" }
```

#### `message`

Broadcast to all connected clients (including the sender) when a message is accepted.

```json
{ "type": "message", "nick": "alice", "text": "hello everyone" }
```

#### `message_error`

Sent only to the sending client when their message is rejected.

```json
{ "type": "message_error", "reason": "too_long" }
```

---

## Connection lifecycle

```
Client                          Server
  |                               |
  |-- WS upgrade (?token=...) --> |
  |                               | 401 if token invalid
  |<-- upgrade accepted --------- |
  |                               |
  |-- { type: "join", nick } ---> |
  |                               | if nick taken:
  |<-- { type: "join_error" } --- |   client may retry with new nick
  |                               |
  |                               | if nick accepted:
  |<-- { type: "join_ok",         |   send user list to new client
  |       users: [...] } -------- |   broadcast user_joined to others
  |                               |
  |-- { type: "message", text} -> | validate length
  |<-- { type: "message_error" }  |   if invalid, send error to sender only
  |                               |   if valid, broadcast to all clients
  |<-- { type: "message",         |
  |       nick, text } ---------- |
  |                               |
  |-- (disconnect) -------------> | free nickname immediately
  |                               | broadcast user_left to remaining clients
```

---

## Server structure

```
/
  main.go          — entry point, configuration, HTTP server setup
  server.go        — hub: manages connected clients, broadcasts
  client.go        — per-connection read/write pumps
  message.go       — message type definitions and JSON marshalling
  Dockerfile
  railway.toml
```

### Hub

A single `Hub` struct owns the authoritative state:

```go
type Hub struct {
    mu      sync.RWMutex
    clients map[string]*Client  // nick → client
}
```

All mutations (join, leave, broadcast) go through the Hub. The Hub does not hold channels — it calls client write methods directly under appropriate locking.

### Client

Each WebSocket connection is managed by a `Client` with two goroutines: a read pump (receives messages from the WebSocket) and a write pump (sends messages to the WebSocket via a buffered channel). This is the standard Go WebSocket pattern.

---

## Configuration

| Variable | Description | Default |
|---|---|---|
| `PORT` | HTTP listen port | `8080` |
| `NOTIRC_TOKEN` | Shared API token (required) | — |

Railway injects `PORT` automatically. `NOTIRC_TOKEN` must be set in the Railway environment.

---

## Deployment

A single `Dockerfile` produces the container image. Multi-stage build: Go builder → minimal `gcr.io/distroless/static` or `alpine` base. Target image size under 20MB.

`railway.toml` configures the start command and health check endpoint (`GET /healthz` returns 200).

---

## Architecture Decision Records

### ADR-001: In-memory state over Redis

**Context:** The server needs to maintain a list of connected nicknames and their associated WebSocket connections.

**Decision:** Use an in-memory `sync.RWMutex`-protected map. No external store.

**Rationale:** The scale is 20–40 concurrent connections for a single workshop session of a few hours. In-memory state is sufficient, free, and eliminates an external dependency. A server restart during a workshop is a recoverable event — users reconnect and re-register their nicknames. Redis would add Railway billing cost, operational complexity (connection pooling, error handling for external calls), and no meaningful benefit at this scale.

**If revisited:** If notIRC grows into a persistent multi-tenant service with multiple server instances, Redis pub/sub for inter-instance message routing and an external store for nickname persistence would be the appropriate next step.

### ADR-002: Query parameter for API token

**Context:** The server needs to authenticate clients before accepting WebSocket connections.

**Decision:** Clients pass the token as a query parameter in the WebSocket upgrade URL (`?token=<token>`).

**Rationale:** Browser WebSocket APIs (`new WebSocket(url)`) do not support custom headers. A query parameter is the only authentication mechanism that works without special handling across all client platforms: browsers, CLIs, mobile apps, and desktop apps. The token is a low-sensitivity shared secret (spam deterrent, not credential protection), so query parameter exposure is acceptable.
