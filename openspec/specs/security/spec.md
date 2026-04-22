# Security

## Requirements

- The server SHALL require a valid API token to accept a WebSocket connection.
- The server SHALL reject the WebSocket upgrade with HTTP 401 if the token is absent or incorrect.
- The server SHALL accept the API token as a query parameter named `token` in the WebSocket upgrade URL.
- The server SHALL read the expected token value from an environment variable (`NOTIRC_TOKEN`).
- The server SHALL NOT implement user authentication, user accounts, or per-user authorisation.

## Scenarios

WHEN a client attempts to upgrade to WebSocket with a correct `token` query parameter
THEN the server accepts the upgrade and the connection proceeds

WHEN a client attempts to upgrade to WebSocket with an incorrect `token` query parameter
THEN the server rejects the upgrade with HTTP 401 and closes the connection

WHEN a client attempts to upgrade to WebSocket without a `token` query parameter
THEN the server rejects the upgrade with HTTP 401 and closes the connection
