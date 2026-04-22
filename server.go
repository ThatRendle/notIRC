package main

import (
	"log/slog"
	"sync"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client
}

func newHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
	}
}

// join registers nick for client. Returns (true, userList) on success or
// (false, nil) if the nick is already taken.
func (h *Hub) join(c *Client, nick string) (bool, []string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.clients[nick]; exists {
		return false, nil
	}

	// Snapshot current users before adding the new one (join_ok excludes self).
	users := make([]string, 0, len(h.clients))
	for n := range h.clients {
		users = append(users, n)
	}

	c.nick = nick
	h.clients[nick] = c
	slog.Info("nick joined", "nick", nick, "total", len(h.clients))

	// Broadcast user_joined to all existing clients.
	b, _ := marshalUserJoined(nick)
	for n, existing := range h.clients {
		if n != nick {
			select {
			case existing.send <- b:
			default:
			}
		}
	}

	return true, users
}

// leave removes the client from the hub and broadcasts user_left.
func (h *Hub) leave(c *Client) {
	if c.nick == "" {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[c.nick] != c {
		return
	}
	delete(h.clients, c.nick)
	slog.Info("nick left", "nick", c.nick, "total", len(h.clients))

	b, _ := marshalUserLeft(c.nick)
	for _, existing := range h.clients {
		select {
		case existing.send <- b:
		default:
		}
	}
}

// broadcast sends a message from nick to all connected clients including sender.
func (h *Hub) broadcast(nick, text string) {
	b, _ := marshalBroadcast(nick, text)

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, c := range h.clients {
		select {
		case c.send <- b:
		default:
		}
	}
}
