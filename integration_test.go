package main

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const testToken = "test-secret"

func startTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	hub := newHub()
	mux := newMux(hub, testToken)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func wsConnect(t *testing.T, srv *httptest.Server, token string) *websocket.Conn {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	if token != "" {
		url += "?token=" + token
	}
	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	t.Cleanup(func() { conn.CloseNow() })
	return conn
}

func send(t *testing.T, conn *websocket.Conn, v any) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := wsjson.Write(ctx, conn, v); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func recv(t *testing.T, conn *websocket.Conn) map[string]any {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var m map[string]any
	if err := wsjson.Read(ctx, conn, &m); err != nil {
		t.Fatalf("read: %v", err)
	}
	return m
}

func join(t *testing.T, conn *websocket.Conn, nick string) map[string]any {
	t.Helper()
	send(t, conn, map[string]any{"type": "join", "nick": nick})
	return recv(t, conn)
}

func TestIntegration_TokenValidation(t *testing.T) {
	srv := startTestServer(t)

	ctx := context.Background()
	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"

	// missing token
	_, resp, err := websocket.Dial(ctx, url, nil)
	if err == nil {
		t.Error("expected connection to fail with no token")
	}
	if resp != nil && resp.StatusCode != 401 {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}

	// wrong token
	_, resp, err = websocket.Dial(ctx, url+"?token=wrong", nil)
	if err == nil {
		t.Error("expected connection to fail with wrong token")
	}
	if resp != nil && resp.StatusCode != 401 {
		t.Errorf("expected 401 for wrong token, got %d", resp.StatusCode)
	}
}

func TestIntegration_FullHandshake(t *testing.T) {
	srv := startTestServer(t)
	conn := wsConnect(t, srv, testToken)

	m := join(t, conn, "alice")
	if m["type"] != "join_ok" {
		t.Fatalf("expected join_ok, got %v", m["type"])
	}
	users, ok := m["users"].([]any)
	if !ok {
		t.Fatal("users field missing or wrong type")
	}
	if len(users) != 0 {
		t.Errorf("expected empty user list for first join, got %v", users)
	}
}

func TestIntegration_JoinError_DuplicateNick(t *testing.T) {
	srv := startTestServer(t)

	c1 := wsConnect(t, srv, testToken)
	join(t, c1, "alice")

	c2 := wsConnect(t, srv, testToken)
	m := join(t, c2, "alice")
	if m["type"] != "join_error" {
		t.Fatalf("expected join_error, got %v", m["type"])
	}
	if m["reason"] != "nick_taken" {
		t.Errorf("expected nick_taken reason, got %v", m["reason"])
	}

	// retry with different nick succeeds
	m = join(t, c2, "bob")
	if m["type"] != "join_ok" {
		t.Fatalf("expected join_ok on retry, got %v", m["type"])
	}
}

func TestIntegration_UserJoined_Broadcast(t *testing.T) {
	srv := startTestServer(t)

	c1 := wsConnect(t, srv, testToken)
	join(t, c1, "alice")

	c2 := wsConnect(t, srv, testToken)
	join(t, c2, "bob")

	// c1 should receive user_joined for bob
	m := recv(t, c1)
	if m["type"] != "user_joined" {
		t.Fatalf("expected user_joined, got %v", m["type"])
	}
	if m["nick"] != "bob" {
		t.Errorf("expected nick=bob, got %v", m["nick"])
	}
}

func TestIntegration_MessageBroadcast(t *testing.T) {
	srv := startTestServer(t)

	c1 := wsConnect(t, srv, testToken)
	join(t, c1, "alice")

	c2 := wsConnect(t, srv, testToken)
	m2join := join(t, c2, "bob")
	_ = m2join

	// Drain the user_joined that c1 received for bob.
	recv(t, c1)

	send(t, c1, map[string]any{"type": "message", "text": "hello"})

	// Both c1 and c2 should receive the broadcast.
	m1 := recv(t, c1)
	m2 := recv(t, c2)

	for _, m := range []map[string]any{m1, m2} {
		if m["type"] != "message" {
			t.Errorf("expected message type, got %v", m["type"])
		}
		if m["nick"] != "alice" {
			t.Errorf("expected nick=alice, got %v", m["nick"])
		}
		if m["text"] != "hello" {
			t.Errorf("expected text=hello, got %v", m["text"])
		}
	}
}

func TestIntegration_UserLeft_Broadcast(t *testing.T) {
	srv := startTestServer(t)

	c1 := wsConnect(t, srv, testToken)
	join(t, c1, "alice")

	c2 := wsConnect(t, srv, testToken)
	join(t, c2, "bob")

	// Drain user_joined on c1.
	recv(t, c1)

	// Disconnect c2 cleanly.
	c2.Close(websocket.StatusNormalClosure, "bye")

	// c1 should receive user_left for bob.
	m := recv(t, c1)
	if m["type"] != "user_left" {
		t.Fatalf("expected user_left, got %v", m["type"])
	}
	if m["nick"] != "bob" {
		t.Errorf("expected nick=bob, got %v", m["nick"])
	}
}

func TestIntegration_MessageError_OversizedMessage(t *testing.T) {
	srv := startTestServer(t)

	conn := wsConnect(t, srv, testToken)
	join(t, conn, "alice")

	oversized := strings.Repeat("x", 1025)
	send(t, conn, map[string]any{"type": "message", "text": oversized})

	m := recv(t, conn)
	if m["type"] != "message_error" {
		t.Fatalf("expected message_error, got %v", m["type"])
	}
	if m["reason"] != "too_long" {
		t.Errorf("expected too_long reason, got %v", m["reason"])
	}

	// Connection should still be open — send another message successfully.
	send(t, conn, map[string]any{"type": "message", "text": "still alive"})
	m = recv(t, conn)
	if m["type"] != "message" {
		t.Fatalf("expected message after error, got %v", m["type"])
	}
}

