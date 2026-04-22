# Proposal: notIRC Backend Server

## What

A WebSocket-based chat server supporting a single shared channel. Clients connect, claim a nickname, and exchange messages in real time. The server maintains authoritative presence state in memory and broadcasts join, leave, and message events to all connected clients.

## Why

This server is the centrepiece of a Claude Code workshop. The final activity asks 20–40 participants to build their own chat client — in any language or framework they choose — and connect to this shared backend. The server must be:

- **Stable** enough to run for a full workshop session without intervention
- **Simple** enough that a participant can understand the full API from a short spec document
- **Accessible** enough that a client can be built against it in any language

## What it is not

- A full IRC server. No IRC protocol, no server-to-server federation, no modes, no operators.
- A persistent service. No message history, no user accounts, no database.
- A scalable multi-tenant platform. Single instance, single channel, ephemeral state.

## Scope

**In scope:**
- WebSocket server with API token authentication (query parameter)
- Nickname registration and uniqueness enforcement
- Real-time broadcast of messages, join events, and leave events
- Presence: current user list delivered on join
- Message validation (1024 UTF-8 byte limit)
- Deployment to Railway as a single container

**Out of scope:**
- Client implementations (built by workshop participants)
- Message history or persistence
- Multiple channels
- Rate limiting or content moderation
- Horizontal scaling or high availability
