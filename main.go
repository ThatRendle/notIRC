package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/coder/websocket"
)

func main() {
	tokensEnv := os.Getenv("NOTIRC_TOKENS")

	rm, err := newRoomManager(tokensEnv)
	if err != nil {
		log.Fatalf("NOTIRC_TOKENS: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := newMux(rm)

	addr := ":" + port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func newMux(rm *RoomManager) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsHandler(rm))
	mux.HandleFunc("/healthz", healthzHandler)
	return mux
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func wsHandler(rm *RoomManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		remote := r.RemoteAddr
		token := r.URL.Query().Get("token")

		room := rm.roomFor(token)
		if room == nil {
			slog.Warn("connection rejected: bad token", "token", token, "remote", remote)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// InsecureSkipVerify disables origin checking so clients from any
		// origin (browser, CLI, mobile) can connect — intentional for a
		// multi-platform workshop server (see design.md ADR-002).
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})
		if err != nil {
			slog.Error("websocket upgrade failed", "remote", remote, "err", err)
			return
		}
		slog.Info("connection accepted", "room", room.logToken, "remote", remote)

		client := newClient(conn, room)
		ctx, cancel := context.WithCancel(r.Context())

		go client.writePump(ctx)
		client.readPump(ctx)
		cancel()
	}
}
