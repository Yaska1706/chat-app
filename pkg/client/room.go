package client

type Room struct {
	name       string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func NewRoom(name string) *Room {
	return &Room{
		name:       name,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (r *Room) Run() {
	select {
	case client := <-r.register:
		r.registerClientToRoom(client)
	case client := <-r.unregister:
		r.unregisterClientFromRoom(client)
	case message := <-r.broadcast:
		r.broadcastToRoom(message)
	}
}
func (r *Room) registerClientToRoom(client *Client) {
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