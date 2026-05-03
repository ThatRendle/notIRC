package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"unicode/utf8"

	"github.com/coder/websocket"
)

const (
	maxMessageBytes = 1024
	writeBufSize    = 32
)

type Client struct {
	nick string
	conn *websocket.Conn
	room *Room
	send chan []byte
}

func newClient(conn *websocket.Conn, room *Room) *Client {
	return &Client{
		conn: conn,
		room: room,
		send: make(chan []byte, writeBufSize),
	}
}

func (c *Client) writePump(ctx context.Context) {
	defer c.conn.CloseNow()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			if err := c.conn.Write(ctx, websocket.MessageText, msg); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *Client) readPump(ctx context.Context) {
	defer func() {
		// leave must run before close(send) so the Room cannot send to a closed channel.
		c.room.leave(c)
		close(c.send)
	}()

	joined := false

	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			return
		}

		var incoming IncomingMessage
		if err := json.Unmarshal(data, &incoming); err != nil {
			continue
		}

		switch incoming.Type {
		case "join":
			if joined {
				continue
			}
			var msg JoinMessage
			if err := json.Unmarshal(data, &msg); err != nil || msg.Nick == "" {
				continue
			}
			ok, users := c.room.join(c, msg.Nick)
			if !ok {
				b, _ := marshalJoinError("nick_taken")
				c.send <- b
				continue
			}
			joined = true
			b, _ := marshalJoinOk(users)
			c.send <- b
		case "message":
			if !joined {
				continue
			}
			var msg SendMessage
			if err := json.Unmarshal(data, &msg); err != nil {
				continue
			}
			if utf8.RuneCountInString(msg.Text) == 0 {
				continue
			}
			if len([]byte(msg.Text)) > maxMessageBytes {
				b, _ := marshalMessageError("too_long")
				c.send <- b
				continue
			}
			c.room.broadcast(c.nick, msg.Text)
		default:
			slog.Warn("unknown message type", "type", incoming.Type)
		}
	}
}
