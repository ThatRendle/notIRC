package main

import (
	"testing"
)

func TestRoom_NicknameUniqueness(t *testing.T) {
	room := newRoom("test")

	c1 := &Client{room: room, send: make(chan []byte, writeBufSize)}
	ok, users := room.join(c1, "alice")
	if !ok {
		t.Fatal("first join should succeed")
	}
	if len(users) != 0 {
		t.Fatalf("expected empty user list for first join, got %v", users)
	}

	c2 := &Client{room: room, send: make(chan []byte, writeBufSize)}
	ok, _ = room.join(c2, "alice")
	if ok {
		t.Fatal("duplicate nick should be rejected")
	}

	c3 := &Client{room: room, send: make(chan []byte, writeBufSize)}
	ok, users = room.join(c3, "bob")
	if !ok {
		t.Fatal("unique nick should be accepted")
	}
	if len(users) != 1 || users[0] != "alice" {
		t.Fatalf("expected [alice] in user list, got %v", users)
	}
}

func TestRoom_Leave_FreesNick(t *testing.T) {
	room := newRoom("test")

	c1 := &Client{room: room, send: make(chan []byte, writeBufSize)}
	room.join(c1, "alice") //nolint:errcheck

	room.leave(c1)

	c2 := &Client{room: room, send: make(chan []byte, writeBufSize)}
	ok, _ := room.join(c2, "alice")
	if !ok {
		t.Fatal("nick should be available after original client leaves")
	}
}

func TestRoomManager_TokenDeduplication(t *testing.T) {
	rm, err := newRoomManager("tok-abc,tok-abc,tok-xyz")
	if err != nil {
		t.Fatalf("newRoomManager: %v", err)
	}
	if len(rm.rooms) != 2 {
		t.Fatalf("expected 2 rooms after deduplication, got %d", len(rm.rooms))
	}
	if rm.roomFor("tok-abc") == nil {
		t.Error("tok-abc room should exist")
	}
	if rm.roomFor("tok-xyz") == nil {
		t.Error("tok-xyz room should exist")
	}
}

func TestRoomManager_EmptyTokensFatal(t *testing.T) {
	_, err := newRoomManager("")
	if err == nil {
		t.Fatal("expected error for empty tokens")
	}
	_, err = newRoomManager(",")
	if err == nil {
		t.Fatal("expected error for comma-only tokens")
	}
}

func TestRoomManager_WhitespaceTokens(t *testing.T) {
	rm, err := newRoomManager("  tok-a , , tok-b  ")
	if err != nil {
		t.Fatalf("newRoomManager: %v", err)
	}
	if len(rm.rooms) != 2 {
		t.Fatalf("expected 2 rooms, got %d", len(rm.rooms))
	}
	if rm.roomFor("tok-a") == nil {
		t.Error("tok-a room should exist after trimming whitespace")
	}
	if rm.roomFor("tok-b") == nil {
		t.Error("tok-b room should exist after trimming whitespace")
	}
}
