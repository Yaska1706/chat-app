package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Yaska1706/chat-app/pkg/client"
)

var addr = flag.String("addr", ":8080", "http server address")

func main() {
	flag.Parse()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		client.ServeWS(w, r)
	})
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
