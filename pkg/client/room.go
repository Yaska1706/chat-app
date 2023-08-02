package client

import "fmt"

const welcomeMessage = "%s joined the room"

type Room struct {
	name       string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

func NewRoom(name string) *Room {
	return &Room{
		name:       name,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

func (r *Room) Run() {
	select {
	case client := <-r.register:
		r.registerClientToRoom(client)
	case client := <-r.unregister:
		r.unregisterClientFromRoom(client)
	case message := <-r.broadcast:
		r.broadcastToRoom(message.Encode())
	}
}
func (r *Room) registerClientToRoom(client *Client) {
	r.notifyClientJoined(client)
	r.clients[client] = true
}

func (r *Room) unregisterClientFromRoom(client *Client) {
	delete(r.clients, client)
}

func (r *Room) broadcastToRoom(message []byte) {

	for client := range r.clients {
		client.send <- message
	}
}

func (r *Room) GetName() string {
	return r.name
}

func (r *Room) notifyClientJoined(client *Client) {
	message := &Message{
		Action:  SendMessageAction,
		Target:  r.name,
		Message: fmt.Sprintf(welcomeMessage, client.GetName()),
	}

	r.broadcastToRoom(message.Encode())
}
