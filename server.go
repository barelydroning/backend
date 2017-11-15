// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	s := string(buf[:reqLen])
	fmt.Println("received: ", s)
	// Send a response back to person contacting us.
	conn.Write([]byte("Message received."))
	// Close the connection when you're done with it.
	conn.Close()
}

func main() {
	fmt.Println("Hello, 世界")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error in connection")
			continue
		}
		fmt.Println("New connection!")
		go handleConnection(conn)
	}
}
