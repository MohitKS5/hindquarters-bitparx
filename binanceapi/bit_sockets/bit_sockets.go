package main

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"flag"
)

var addr = flag.String("addr", "localhost:4201", "http service address")
var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

func echo(w http.ResponseWriter, r *http.Request) {
	client, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer client.Close()
	for {
		mt, message, err := client.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = client.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
