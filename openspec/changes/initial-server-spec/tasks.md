# Tasks: notIRC Backend Server

## Programmer

- [ ] Initialise Go module (`github.com/rendle/notirc`)
- [ ] Add `github.com/coder/websocket` dependency
- [ ] Implement `message.go` — define all message types and JSON marshalling/unmarshalling
- [ ] Implement `client.go` — per-connection struct with read pump and write pump goroutines
- [ ] Implement `server.go` — Hub struct with `sync.RWMutex`-protected client map; join, leave, and broadcast methods
- [ ] Implement connection handler — validate token from query parameter (HTTP 401 on failure), upgrade to WebSocket, hand off to Hub
- [ ] Implement join flow — receive `join` message, validate nickname uniqueness, send `join_ok` + user list or `join_error`
- [ ] Implement message flow — receive `message`, validate UTF-8 byte length (<= 1024), broadcast or send `message_error`
- [ ] Implement disconnect handling — remove client from Hub, broadcast `user_left`
- [ ] Implement `GET /healthz` endpoint returning 200
- [ ] Implement `main.go` — read `PORT` and `NOTIRC_TOKEN` from environment, start HTTP server
- [ ] Write `Dockerfile` — multi-stage build, distroless or alpine base, expose `PORT`
- [ ] Write `railway.toml` — configure start command and health check

## Tester

- [ ] Unit test: nickname uniqueness enforcement (accept unique, reject duplicate)
- [ ] Unit test: message byte length validation (accept <= 1024, reject > 1024)
- [ ] Unit test: JSON marshalling for all message types
- [ ] Integration test: full connection handshake (token validation → join → join_ok with user list)
- [ ] Integration test: join_error on duplicate nickname, retry with new nickname succeeds
- [ ] Integration test: message broadcast reaches all connected clients including sender
- [ ] Integration test: user_joined broadcast on new client join
- [ ] Integration test: user_left broadcast on client disconnect
- [ ] Integration test: message_error on oversized message, connection remains open
- [ ] Integration test: HTTP 401 on missing or incorrect token
