package client

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

type Client struct {
	conn *websocket.Conn
}

func newClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

func ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	client := newClient(conn)

	fmt.Println("A new connection joined")
	fmt.Print(client)
}
