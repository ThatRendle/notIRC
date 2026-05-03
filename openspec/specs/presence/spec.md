# Presence

## Purpose

Track connected client nicknames and broadcast join/leave events.

## Requirements

### Requirement: Authoritative presence list
The server SHALL maintain an authoritative list of currently connected nicknames.

#### Scenario: Presence list is server-controlled
- **WHEN** clients connect and disconnect
- **THEN** the server alone determines the contents of the nickname list

### Requirement: Presence list in handshake
The server SHALL include the full nickname list in the connection handshake response sent to a newly joined client.

#### Scenario: Nickname list on join
- **WHEN** a client successfully joins
- **THEN** the server sends that client the complete list of currently connected nicknames before any other messages

### Requirement: Immediate presence updates
The server SHALL update the nickname list immediately when a client joins or leaves.

#### Scenario: List updated on join
- **WHEN** a client joins
- **THEN** their nickname appears in the presence list without delay

#### Scenario: List updated on disconnect
- **WHEN** a client disconnects
- **THEN** their nickname is removed from the presence list without delay

### Requirement: Join event broadcast
The server SHALL broadcast a join event when a client's nickname is accepted.

#### Scenario: Join event sent to all
- **WHEN** a client joins
- **THEN** all previously connected clients receive a join event containing the new client's nickname

### Requirement: Leave event broadcast
The server SHALL broadcast a leave event when a client disconnects.

#### Scenario: Leave event sent to remaining
- **WHEN** a client disconnects
- **THEN** all remaining connected clients receive a leave event containing the departed client's nickname

#### Scenario: Nickname removed on disconnect
- **WHEN** a client disconnects
- **THEN** that client's nickname is removed from the server's nickname list immediately
