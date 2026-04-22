# Tasks: notIRC Backend Server

## Programmer

- [x] Initialise Go module (`github.com/rendle/notirc`)
- [x] Add `github.com/coder/websocket` dependency
- [x] Implement `message.go` — define all message types and JSON marshalling/unmarshalling
- [x] Implement `client.go` — per-connection struct with read pump and write pump goroutines
- [x] Implement `server.go` — Hub struct with `sync.RWMutex`-protected client map; join, leave, and broadcast methods
- [x] Implement connection handler — validate token from query parameter (HTTP 401 on failure), upgrade to WebSocket, hand off to Hub
- [x] Implement join flow — receive `join` message, validate nickname uniqueness, send `join_ok` + user list or `join_error`
- [x] Implement message flow — receive `message`, validate UTF-8 byte length (<= 1024), broadcast or send `message_error`
- [x] Implement disconnect handling — remove client from Hub, broadcast `user_left`
- [x] Implement `GET /healthz` endpoint returning 200
- [x] Implement `main.go` — read `PORT` and `NOTIRC_TOKEN` from environment, start HTTP server
- [x] Write `Dockerfile` — multi-stage build, distroless or alpine base, expose `PORT`
- [x] Write `railway.toml` — configure start command and health check

## Tester

- [x] Unit test: nickname uniqueness enforcement (accept unique, reject duplicate)
- [x] Unit test: message byte length validation (accept <= 1024, reject > 1024)
- [x] Unit test: JSON marshalling for all message types
- [x] Integration test: full connection handshake (token validation → join → join_ok with user list)
- [x] Integration test: join_error on duplicate nickname, retry with new nickname succeeds
- [x] Integration test: message broadcast reaches all connected clients including sender
- [x] Integration test: user_joined broadcast on new client join
- [x] Integration test: user_left broadcast on client disconnect
- [x] Integration test: message_error on oversized message, connection remains open
- [x] Integration test: HTTP 401 on missing or incorrect token
