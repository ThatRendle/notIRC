package main

import "encoding/json"

// Incoming message types (client → server)

type IncomingMessage struct {
	Type string `json:"type"`
}

type JoinMessage struct {
	Type string `json:"type"`
	Nick string `json:"nick"`
}

type SendMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Outgoing message types (server → client)

type JoinOkMessage struct {
	Type  string   `json:"type"`
	Users []string `json:"users"`
}

type JoinErrorMessage struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type UserJoinedMessage struct {
	Type string `json:"type"`
	Nick string `json:"nick"`
}

type UserLeftMessage struct {
	Type string `json:"type"`
	Nick string `json:"nick"`
}

type BroadcastMessage struct {
	Type string `json:"type"`
	Nick string `json:"nick"`
	Text string `json:"text"`
}

type MessageErrorMessage struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

func marshalJoinOk(users []string) ([]byte, error) {
	if users == nil {
		users = []string{}
	}
	return json.Marshal(JoinOkMessage{Type: "join_ok", Users: users})
}

func marshalJoinError(reason string) ([]byte, error) {
	return json.Marshal(JoinErrorMessage{Type: "join_error", Reason: reason})
}

func marshalUserJoined(nick string) ([]byte, error) {
	return json.Marshal(UserJoinedMessage{Type: "user_joined", Nick: nick})
}

func marshalUserLeft(nick string) ([]byte, error) {
	return json.Marshal(UserLeftMessage{Type: "user_left", Nick: nick})
}

func marshalBroadcast(nick, text string) ([]byte, error) {
	return json.Marshal(BroadcastMessage{Type: "message", Nick: nick, Text: text})
}

func marshalMessageError(reason string) ([]byte, error) {
	return json.Marshal(MessageErrorMessage{Type: "message_error", Reason: reason})
}
