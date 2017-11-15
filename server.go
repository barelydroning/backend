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
	Pitch    float64
	Roll     float64
	Azimuth  float64
	Altitude float64
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

func handleState(onChange chan sensorData, onUpdate chan sensorData) {
	for {
		newState := <-onUpdate
		log.Println("Updating state")
		onChange <- newState
	}
}

func handleInput(conn *websocket.Conn, onUpdate chan sensorData) {
	for {
		state := &sensorData{0, 0, 0, 0}
		err := conn.ReadJSON(state)
		if err != nil {
			log.Println(err)
		}
		onUpdate <- *state
	}
}

func handleOutput(conn *websocket.Conn, onChange chan sensorData) {
	state := sensorData{0, 0, 0, 0}
	ticker := time.NewTicker(2000 * time.Millisecond)
	for {
		select {
		case newState := <-onChange:
			state = newState
		case <-ticker.C:
			fmt.Println("Sending data...")
			conn.WriteJSON(state)
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	onChange := make(chan sensorData)
	onUpdate := make(chan sensorData)

	go handleState(onChange, onUpdate)
	go handleInput(conn, onUpdate)
	go handleOutput(conn, onChange)
}

func main() {
	fmt.Println("Hello, 世界")
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Println(err)
	}
}
