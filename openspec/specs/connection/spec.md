# Connection

## Purpose

Establish WebSocket connections and manage client lifecycle including nickname registration.

## Requirements

### Requirement: WebSocket connections and nicknames
The server SHALL accept client connections over WebSocket. The server SHALL require a client to submit a nickname before it is considered connected.

#### Scenario: WebSocket connection accepted
- **WHEN** a client opens a WebSocket connection to `/ws`
- **THEN** the server upgrades the connection to WebSocket

#### Scenario: Nickname must be submitted before messaging
- **WHEN** a client connects but has not submitted a valid nickname
- **THEN** the server ignores any `message` type messages from that client

### Requirement: Nickname uniqueness
The server SHALL reject a nickname that is already in use by a connected client and prompt the client to choose a different one. The server SHALL accept a nickname once it is confirmed to be unique.

#### Scenario: Nickname accepted
- **WHEN** a client connects and submits a nickname that is not in use
- **THEN** the server responds with a success confirmation and the current list of connected nicknames

#### Scenario: Nickname rejected
- **WHEN** a client connects and submits a nickname that is already in use
- **THEN** the server rejects the nickname and indicates the client should choose a different one

### Requirement: Handshake and join event
Upon accepting a nickname, the server SHALL send the connecting client the current list of connected nicknames. The server SHALL broadcast a join event to all other connected clients.

#### Scenario: Join event on accepted nickname
- **WHEN** a client's nickname is accepted
- **THEN** the server broadcasts a join event containing the new nickname to all other connected clients

### Requirement: Nickname lifecycle on disconnect
The server SHALL free a nickname immediately when its client disconnects. The server SHALL broadcast a leave event to all connected clients when a client disconnects.

#### Scenario: Leave event on disconnect
- **WHEN** a client disconnects for any reason (clean or dropped)
- **THEN** the server immediately frees the nickname
- **AND** broadcasts a leave event containing the nickname to all remaining connected clients
