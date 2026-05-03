## Purpose

Isolate message delivery per room so that clients connected with different tokens do not receive each other's messages.

## MODIFIED Requirements

### Requirement: Single shared channel
The server SHALL support multiple isolated channels, each associated with a configured token. All clients connecting with the same token participate in the same room. The server SHALL broadcast each message received from a client to all connected clients in the same room, including the sender. The server SHALL NOT deliver messages across room boundaries.

#### Scenario: Message broadcast within a room
- **WHEN** a connected client sends a message within the 1024 UTF-8 byte limit
- **THEN** the server broadcasts the message and the sender's nickname to all connected clients in the same room

#### Scenario: Message isolation across rooms
- **WHEN** a client in room A sends a message
- **THEN** clients in room B do NOT receive the message
