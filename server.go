package main

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
)

type Room struct {
	mu       sync.RWMutex
	clients  map[string]*Client
	logToken string
}

// RoomManager holds a fixed set of rooms keyed by token. The rooms map is
// populated at startup and never mutated thereafter — all roomFor() calls
// are safe without synchronization.
type RoomManager struct {
	rooms map[string]*Room
}

func newRoomManager(tokensEnv string) (*RoomManager, error) {
	if tokensEnv == "" {
		return nil, fmt.Errorf("NOTIRC_TOKENS environment variable is required")
	}

	seen := make(map[string]struct{})
	rooms := make(map[string]*Room)

	for _, token := range strings.Split(tokensEnv, ",") {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}
		if _, exists := seen[token]; exists {
			continue
		}
		seen[token] = struct{}{}
		rooms[token] = newRoom(token)
	}

	if len(rooms) == 0 {
		return nil, fmt.Errorf("NOTIRC_TOKENS must contain at least one token")
	}

	return &RoomManager{rooms: rooms}, nil
}

func (rm *RoomManager) roomFor(token string) *Room {
	return rm.rooms[token]
}

func newRoom(token string) *Room {
	logToken := token
	if len(logToken) > 8 {
		logToken = logToken[:8]
	}
	return &Room{
		clients:  make(map[string]*Client),
		logToken: logToken,
	}
}

// join registers nick for client. Returns (true, userList) on success or
// (false, nil) if the nick is already taken.
func (r *Room) join(c *Client, nick string) (bool, []string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.clients[nick]; exists {
		return false, nil
	}

	// Snapshot current users before adding the new one (join_ok excludes self).
	users := make([]string, 0, len(r.clients))
	for n := range r.clients {
		users = append(users, n)
	}

	c.nick = nick
	r.clients[nick] = c
	slog.Info("nick joined", "room", r.logToken, "nick", nick, "total", len(r.clients))

	// Broadcast user_joined to all existing clients.
	b, _ := marshalUserJoined(nick)
	for n, existing := range r.clients {
		if n != nick {
			select {
			case existing.send <- b:
			default:
			}
		}
	}

	return true, users
}

// leave removes the client from the room and broadcasts user_left.
func (r *Room) leave(c *Client) {
	if c.nick == "" {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.clients[c.nick] != c {
		return
	}
	delete(r.clients, c.nick)
	slog.Info("nick left", "room", r.logToken, "nick", c.nick, "total", len(r.clients))

	b, _ := marshalUserLeft(c.nick)
	for _, existing := range r.clients {
		select {
		case existing.send <- b:
		default:
		}
	}
}

// broadcast sends a message from nick to all connected clients including sender.
func (r *Room) broadcast(nick, text string) {
	b, _ := marshalBroadcast(nick, text)

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, c := range r.clients {
		select {
		case c.send <- b:
		default:
		}
	}
}
