# Presence

## Requirements

- The server SHALL maintain an authoritative list of currently connected nicknames.
- The server SHALL include the full nickname list in the connection handshake response sent to a newly joined client.
- The server SHALL update the nickname list immediately when a client joins or leaves.
- The server SHALL broadcast a join event when a client's nickname is accepted.
- The server SHALL broadcast a leave event when a client disconnects.

## Scenarios

WHEN a client successfully joins
THEN the server sends that client the complete list of currently connected nicknames before any other messages

WHEN a client joins
THEN all previously connected clients receive a join event containing the new client's nickname

WHEN a client disconnects
THEN all remaining connected clients receive a leave event containing the departed client's nickname

WHEN a client disconnects
THEN that client's nickname is removed from the server's nickname list immediately
