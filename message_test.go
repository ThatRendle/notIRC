package main

import (
	"encoding/json"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestMessageByteLength(t *testing.T) {
	exactly1024 := strings.Repeat("a", 1024)
	over1024 := strings.Repeat("a", 1025)
	multibyte := strings.Repeat("é", 513) // 513 × 2 bytes = 1026 bytes

	cases := []struct {
		name    string
		text    string
		allowed bool
	}{
		{"empty", "", true},
		{"1024 ascii bytes", exactly1024, true},
		{"1025 ascii bytes", over1024, false},
		{"1026 bytes via multibyte chars", multibyte, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			byteLen := len([]byte(tc.text))
			tooLong := byteLen > maxMessageBytes
			if tc.allowed && tooLong {
				t.Errorf("expected %q (%d bytes) to be allowed", tc.name, byteLen)
			}
			if !tc.allowed && !tooLong {
				t.Errorf("expected %q (%d bytes) to be rejected", tc.name, byteLen)
			}
		})
	}
}

func TestMarshalJoinOk(t *testing.T) {
	b, err := marshalJoinOk([]string{"alice", "bob"})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m["type"] != "join_ok" {
		t.Errorf("expected type=join_ok, got %v", m["type"])
	}
}

func TestMarshalJoinOk_NilUsers(t *testing.T) {
	b, err := marshalJoinOk(nil)
	if err != nil {
		t.Fatal(err)
	}
	var m JoinOkMessage
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m.Users == nil {
		t.Error("users should be empty slice, not null")
	}
}

func TestMarshalJoinError(t *testing.T) {
	b, err := marshalJoinError("nick_taken")
	if err != nil {
		t.Fatal(err)
	}
	var m JoinErrorMessage
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m.Type != "join_error" || m.Reason != "nick_taken" {
		t.Errorf("unexpected join_error payload: %+v", m)
	}
}

func TestMarshalUserJoined(t *testing.T) {
	b, err := marshalUserJoined("alice")
	if err != nil {
		t.Fatal(err)
	}
	var m UserJoinedMessage
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m.Type != "user_joined" || m.Nick != "alice" {
		t.Errorf("unexpected user_joined payload: %+v", m)
	}
}

func TestMarshalUserLeft(t *testing.T) {
	b, err := marshalUserLeft("alice")
	if err != nil {
		t.Fatal(err)
	}
	var m UserLeftMessage
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m.Type != "user_left" || m.Nick != "alice" {
		t.Errorf("unexpected user_left payload: %+v", m)
	}
}

func TestMarshalBroadcast(t *testing.T) {
	b, err := marshalBroadcast("alice", "hello")
	if err != nil {
		t.Fatal(err)
	}
	var m BroadcastMessage
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m.Type != "message" || m.Nick != "alice" || m.Text != "hello" {
		t.Errorf("unexpected broadcast payload: %+v", m)
	}
}

func TestMarshalMessageError(t *testing.T) {
	b, err := marshalMessageError("too_long")
	if err != nil {
		t.Fatal(err)
	}
	var m MessageErrorMessage
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m.Type != "message_error" || m.Reason != "too_long" {
		t.Errorf("unexpected message_error payload: %+v", m)
	}
}

func TestUTF8ByteLength(t *testing.T) {
	s := strings.Repeat("é", 512)
	if utf8.RuneCountInString(s) != 512 {
		t.Errorf("expected 512 runes")
	}
	if len([]byte(s)) != 1024 {
		t.Errorf("expected 1024 bytes, got %d", len([]byte(s)))
	}
}
