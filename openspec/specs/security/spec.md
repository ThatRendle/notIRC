# Security

## Purpose

Authenticate clients via a shared API token and reject unauthorised connections.
## Requirements
### Requirement: Token authentication
The server SHALL require a valid API token to accept a WebSocket connection.

#### Scenario: Client connects with a valid token
- **WHEN** a client attempts to upgrade to WebSocket with a correct `token` query parameter
- **THEN** the server accepts the upgrade and the connection proceeds

### Requirement: Token validation rejection
The server SHALL reject the WebSocket upgrade with HTTP 401 if the token is absent or incorrect.

#### Scenario: Client connects with an incorrect token
- **WHEN** a client attempts to upgrade to WebSocket with an incorrect `token` query parameter
- **THEN** the server rejects the upgrade with HTTP 401 and closes the connection

#### Scenario: Client connects without a token
- **WHEN** a client attempts to upgrade to WebSocket without a `token` query parameter
- **THEN** the server rejects the upgrade with HTTP 401 and closes the connection

### Requirement: Token transport
The server SHALL accept the API token as a query parameter named `token` in the WebSocket upgrade URL.

#### Scenario: Token passed as query parameter
- **WHEN** a client connects to `/ws?token=my-secret`
- **THEN** the server reads the token from the `token` query parameter

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

### Requirement: No user authentication
The server SHALL NOT implement user authentication, user accounts, or per-user authorisation.

#### Scenario: No login required
- **WHEN** a client connects with a valid token
- **THEN** the server does not require additional login credentials or user registration

