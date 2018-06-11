package bit_sockets

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"flag"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options

func echo (w http.ResponseWriter, r *http.Request) {
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