## 1. Core Implementation

- [x] 1.1 Rename `Hub` type to `Room` in server.go, update all method receivers and references
- [x] 1.2 Create `RoomManager` struct with `rooms map[string]*Room` and token lookup method
- [x] 1.3 Implement `NOTIRC_TOKENS` parsing: comma-split, deduplicate, fatal on empty, pre-create all rooms
- [x] 1.4 Update `newMux` and `wsHandler` to accept `*RoomManager` instead of `*Hub` + `token`
- [x] 1.5 Add room identifier (truncated token) to all `Room` method log calls
- [x] 1.6 Update `main.go` to use `NOTIRC_TOKENS` and `RoomManager`

## 2. Tests

- [x] 2.1 Rename `Hub` → `Room` in `server_test.go` and `integration_test.go`
- [x] 2.2 Update `startTestServer` helper to create a `RoomManager` with test tokens
- [x] 2.3 Add `TestRoomManager_TokenDeduplication` — verify duplicate tokens produce one room
- [x] 2.4 Add `TestRoomManager_UnrecognizedToken` — verify 401 for unconfigured tokens
- [x] 2.5 Add `TestIntegration_RoomIsolation` — verify clients on different tokens don't see each other
- [x] 2.6 Add `TestRoomManager_EmptyTokensFatal` — verify server exits on empty `NOTIRC_TOKENS`

## 3. Documentation

- [x] 3.1 Update `README.md` — replace `NOTIRC_TOKEN` examples with `NOTIRC_TOKENS`
- [x] 3.2 Update `docs/BUILD_A_CLIENT.md` — update environment variable references
- [x] 3.3 Update `docker-compose.yml` — change environment variable name
