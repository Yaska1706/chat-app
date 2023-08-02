package client

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 10000
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

var (
	newline = []byte{'\n'}
	// space   = []byte{' '}
)

type Client struct {
	conn     *websocket.Conn
	wsServer *WsServer

	send  chan []byte
	rooms map[*Room]bool

	Name string `json:"name"`
}

func newClient(conn *websocket.Conn, wsServer *WsServer, name string) *Client {
	return &Client{
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte),
		rooms:    make(map[*Room]bool),
		Name:     name,
	}
}

func ServeWS(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["name"]
	if !ok || len(name[0]) < 1 {
		log.Println("URL Param 'name' is missing")
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	client := newClient(conn, wsServer, name[0])

	go client.readPump()
	go client.writePump()
}

func (c *Client) disconnect() {
	c.wsServer.unregister <- c
	for room := range c.rooms {
		room.unregister <- c
	}
	close(c.send)
	c.conn.Close()
}

func (c *Client) readPump() {
	defer func() {
		c.disconnect()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, jsonMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error : %s", err)
			}
			break
		}
		c.handleNewMessage(jsonMessage)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, message)
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("NewWriter() : %s", err)
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				log.Printf("Close(): %s", err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("%s", err)
				return
			}
		}
	}
}

func (c *Client) handleNewMessage(jsonMessage []byte) {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Unmarsha() : %s", err)
		return
	}

	message.Sender = c

	switch message.Action {
	case SendMessageAction:
		roomName := message.Target
		if room := c.wsServer.findRoomByName(roomName); room != nil {
			room.broadcast <- &message
		}
	case JoinAction:
		c.handleJoinRoomMessage(message)
	case LeaveAction:
		c.handleLeaveRoomMessage(message)
	}
}

func (c *Client) handleJoinRoomMessage(message Message) {
	roomName := message.Target

	room := c.wsServer.findRoomByName(roomName)
	if room == nil {
		room = c.wsServer.createRoom(roomName)
	}
	c.rooms[room] = true

	room.register <- c
}
func (c *Client) handleLeaveRoomMessage(message Message) {
	room := c.wsServer.findRoomByName(message.Target)
	delete(c.rooms, room)
	room.unregister <- c
}

func (c *Client) GetName() string {
	return c.Name
}
