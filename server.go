// You can edit this code!
// Click here and start typing.
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type sensorData struct {
	Pitch    float64 `json:"pitch"`
	Roll     float64 `json:"roll"`
	Azimuth  float64 `json:"azimuth"`
	Altitude float64 `json:"altitude"`
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
	ticker := time.NewTicker(1000 * time.Millisecond)
	state := sensorData{0, 0, 0, 0}
	for {
		select {
		case newState := <-onUpdate:
			log.Println("Updating state")
			state = newState
			onChange <- newState

		case <-ticker.C:
			state.Pitch += (rand.Float64() - 0.5) * 10
			state.Roll += (rand.Float64() - 0.5) * 10
			state.Azimuth += (rand.Float64() - 0.5) * 10
			state.Altitude += (rand.Float64() - 0.5) * 10
			onChange <- state
		}
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
	ticker := time.NewTicker(100 * time.Millisecond)
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
