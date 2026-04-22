package main

import (
	"testing"
)

func TestHub_NicknameUniqueness(t *testing.T) {
	hub := newHub()

	c1 := &Client{hub: hub, send: make(chan []byte, writeBufSize)}
	ok, users := hub.join(c1, "alice")
	if !ok {
		t.Fatal("first join should succeed")
	}
	if len(users) != 0 {
		t.Fatalf("expected empty user list for first join, got %v", users)
	}

	c2 := &Client{hub: hub, send: make(chan []byte, writeBufSize)}
	ok, _ = hub.join(c2, "alice")
	if ok {
		t.Fatal("duplicate nick should be rejected")
	}

	c3 := &Client{hub: hub, send: make(chan []byte, writeBufSize)}
	ok, users = hub.join(c3, "bob")
	if !ok {
		t.Fatal("unique nick should be accepted")
	}
	if len(users) != 1 || users[0] != "alice" {
		t.Fatalf("expected [alice] in user list, got %v", users)
	}
}

func TestHub_Leave_FreesNick(t *testing.T) {
	hub := newHub()

	c1 := &Client{hub: hub, send: make(chan []byte, writeBufSize)}
	hub.join(c1, "alice") //nolint:errcheck

	hub.leave(c1)

	c2 := &Client{hub: hub, send: make(chan []byte, writeBufSize)}
	ok, _ := hub.join(c2, "alice")
	if !ok {
		t.Fatal("nick should be available after original client leaves")
	}
}
