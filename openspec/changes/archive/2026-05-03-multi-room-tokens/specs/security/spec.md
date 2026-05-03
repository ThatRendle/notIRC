## Purpose

Update token authentication to support multiple tokens, each granting access to an isolated chat room.

## MODIFIED Requirements

### Requirement: Token environment variable
The server SHALL read the expected token values from a comma-separated list in the `NOTIRC_TOKENS` environment variable. The server SHALL accept any token present in the configured list. The server SHALL exit with a fatal error if `NOTIRC_TOKENS` is empty or not set.

#### Scenario: Client connects with a valid token
- **WHEN** a client attempts to upgrade to WebSocket with a `token` query parameter that matches one of the configured tokens
- **THEN** the server accepts the upgrade and the connection proceeds

#### Scenario: Client connects with an incorrect token
- **WHEN** a client attempts to upgrade to WebSocket with a `token` query parameter that does not match any configured token
- **THEN** the server rejects the upgrade with HTTP 401 and closes the connection

#### Scenario: Client connects without a token
- **WHEN** a client attempts to upgrade to WebSocket without a `token` query parameter
- **THEN** the server rejects the upgrade with HTTP 401 and closes the connection
