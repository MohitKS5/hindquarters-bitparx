package bnb_socket

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "stream.binance.com:9443", "http service address")

func connect(path string) (client *websocket.Conn, err error) {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: *addr, Path: path}
	log.Printf("connecting to %s", u.String())

	client, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return
}

var stopC = make(chan struct{})
func GetMessages(path string) {
	client, err := connect(path)
	if err != nil {
		panic(err)
	}

	doneC := make(chan struct{})

	go func() {
		defer func() {
			clienterr := client.Close()
			if clienterr != nil {
				log.Println(clienterr)
			}
		}()
		defer close(doneC)
		for {
			select {
			case <-stopC:
				return
			default:
				_, message, err := client.ReadMessage()
				if err != nil {
					go errorHandler(err)
					return
				}
			go messageHandler(string(message))
			}
		}
	}()
}
func disconnect()  {
	close(stopC)
}

func errorHandler(err error)  {
	log.Println("read error:", err)
}

func messageHandler(message string)  {
	log.Printf("recv: %s", message)
}