// You can edit this code!
// Click here and start typing.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type sensorData struct {
	pitch    float64
	roll     float64
	azimuth  float64
	altitude float64
}

var addr = flag.String("addr", ":8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// allow all connections
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func inputHandler(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		if messageType == websocket.TextMessage {
			fmt.Println(string(p))
		} else {
			log.Println("Byte message not supported")
		}
	}
}

func outputHandler(conn *websocket.Conn) {
	ticker := time.NewTicker(2000 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			fmt.Println("Sending data...")
			conn.WriteMessage(websocket.TextMessage, []byte("Hello!"))
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	go inputHandler(conn)
	go outputHandler(conn)
}

func main() {
	fmt.Println("Hello, 世界")
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Println(err)
	}
}
