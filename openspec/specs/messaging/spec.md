# Messaging

## Purpose

Deliver messages between connected clients within a shared chat channel.
## Requirements
### Requirement: Single shared channel
The server SHALL support multiple isolated channels, each associated with a configured token. All clients connecting with the same token participate in the same room. The server SHALL broadcast each message received from a client to all connected clients in the same room, including the sender. The server SHALL NOT deliver messages across room boundaries.

#### Scenario: Message broadcast within a room
- **WHEN** a connected client sends a message within the 1024 UTF-8 byte limit
- **THEN** the server broadcasts the message and the sender's nickname to all connected clients in the same room

#### Scenario: Message isolation across rooms
- **WHEN** a client in room A sends a message
- **THEN** clients in room B do NOT receive the message

### Requirement: Message byte limit
The server SHALL reject any message exceeding 1024 UTF-8 bytes in length.

#### Scenario: Oversized message rejected
- **WHEN** a connected client sends a message exceeding 1024 UTF-8 bytes
- **THEN** the server rejects the message and notifies the sending client

### Requirement: @mentions pass-through
The server SHALL NOT assign any special meaning to `@nickname` mentions; they are a text convention for clients to handle.

#### Scenario: @mention treated as plain text
- **WHEN** a message contains an `@nickname` mention
- **THEN** the server broadcasts it unchanged, treating it as plain text

### Requirement: No message history
The server SHALL NOT deliver messages sent before a client connected (no history).

#### Scenario: No prior messages on connect
- **WHEN** a client connects
- **THEN** it receives only messages sent after the moment it joined; no prior history is delivered

### Requirement: No rate limiting or content filtering
The server SHALL NOT enforce rate limiting or content filtering.

#### Scenario: No rate limiting applied
- **WHEN** a client sends many messages rapidly
- **THEN** the server delivers all messages without throttling or blocking

