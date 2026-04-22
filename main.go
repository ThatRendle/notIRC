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
	token := os.Getenv("NOTIRC_TOKEN")
	if token == "" {
		log.Fatal("NOTIRC_TOKEN environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	hub := newHub()
	mux := newMux(hub, token)

	addr := ":" + port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func newMux(hub *Hub, token string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsHandler(hub, token))
	mux.HandleFunc("/healthz", healthzHandler)
	return mux
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func wsHandler(hub *Hub, token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		remote := r.RemoteAddr
		if r.URL.Query().Get("token") != token {
			slog.Warn("connection rejected: bad token", "remote", remote)
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
		slog.Info("connection accepted", "remote", remote)

		client := newClient(conn, hub)
		ctx, cancel := context.WithCancel(r.Context())

		go client.writePump(ctx)
		client.readPump(ctx)
		cancel()
	}
}
