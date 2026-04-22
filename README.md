# notIRC

A WebSocket-based chat server built for the Claude Code workshop. Everyone connects to the same channel, picks a nickname, and sends messages. The final workshop activity is to build a client — in any language or framework you like — that connects to this server.

## Connecting

Connect to the WebSocket endpoint with a valid API token:

```
wss://<host>/ws?token=<token>
```

The token will be provided by the workshop facilitator. Connections without a valid token receive HTTP 401.

## Protocol

All messages are JSON with a `type` field.

### Joining

After connecting, send a `join` message with your chosen nickname:

```json
{ "type": "join", "nick": "alice" }
```

If the nickname is available you'll receive a `join_ok` with the current user list:

```json
{ "type": "join_ok", "users": ["bob", "carol"] }
```

If the nickname is taken you'll receive a `join_error` — pick a different one and try again:

```json
{ "type": "join_error", "reason": "nick_taken" }
```

### Sending a message

```json
{ "type": "message", "text": "hello everyone" }
```

Messages must be 1024 UTF-8 bytes or fewer. Oversized messages are rejected with a `message_error` and your connection stays open:

```json
{ "type": "message_error", "reason": "too_long" }
```

### Receiving messages

**Broadcast message** (sent to all connected clients including the sender):

```json
{ "type": "message", "nick": "alice", "text": "hello everyone" }
```

**Someone joined:**

```json
{ "type": "user_joined", "nick": "alice" }
```

**Someone left:**

```json
{ "type": "user_left", "nick": "alice" }
```

### `@mentions`

`@nickname` in a message is a convention for clients to handle — the server passes it through unchanged.

## Health check

```
GET /healthz → 200 OK
```

## Running locally

```sh
NOTIRC_TOKEN=yourtoken PORT=8080 go run .
```

## Running with Docker

```sh
docker build -t notirc .
docker run -e NOTIRC_TOKEN=yourtoken -p 8080:8080 notirc
```

## Running the published image

Pre-built images are published to the GitHub Container Registry on every release tag:

```sh
docker run -e NOTIRC_TOKEN=yourtoken -p 8080:8080 ghcr.io/rendle/notirc:latest
```

To pin to a specific release:

```sh
docker run -e NOTIRC_TOKEN=yourtoken -p 8080:8080 ghcr.io/rendle/notirc:v1.0.0
```

## Environment variables

| Variable | Description | Required |
|---|---|---|
| `NOTIRC_TOKEN` | Shared API token for client authentication | Yes |
| `PORT` | HTTP listen port (Railway sets this automatically) | No (default: 8080) |

## Deployment

The server deploys automatically to Railway on push to `main`.
