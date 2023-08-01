package client

import (
	"encoding/json"
	"log"
)

const (
	JoinAction        = "join-room"
	SendMessageAction = "send-message"
	LeaveAction       = "leave-room"
)

type Message struct {
	Action  string  `json:"action"`
	Message string  `json:"message"`
	Target  string  `json:"target"`
	Sender  *Client `json:"sender"`
}

func (m *Message) Encode() []byte {
	jsonValue, err := json.Marshal(m)
	if err != nil {
		log.Printf("Marshal() : %s", err)
		return nil
	}
	return jsonValue
}
