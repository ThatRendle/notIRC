# Messaging

## Requirements

- The server SHALL support a single shared channel that all connected clients participate in.
- The server SHALL broadcast each message received from a client to all connected clients, including the sender.
- The server SHALL reject any message exceeding 1024 UTF-8 bytes in length.
- The server SHALL attach the sender's nickname to each broadcast message.
- The server SHALL NOT assign any special meaning to `@nickname` mentions; they are a text convention for clients to handle.
- The server SHALL NOT deliver messages sent before a client connected (no history).
- The server SHALL NOT enforce rate limiting or content filtering.

## Scenarios

WHEN a connected client sends a message within the 1024 UTF-8 byte limit
THEN the server broadcasts the message and the sender's nickname to all connected clients

WHEN a connected client sends a message exceeding 1024 UTF-8 bytes
THEN the server rejects the message and notifies the sending client

WHEN a message contains an `@nickname` mention
THEN the server broadcasts it unchanged, treating it as plain text

WHEN a client connects
THEN it receives only messages sent after the moment it joined; no prior history is delivered
