package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Error connecting to server: %v\n", err)
	}

	fmt.Println("Connected to server on port 8080...")

	for range 5 {
		message := []byte("hi there")
		length := uint16(len(message))

		var buf bytes.Buffer

		err := binary.Write(&buf, binary.BigEndian, length)
		if err != nil {
			log.Fatalf("Error writing length to buffer: %v\n", err)
		}

		buf.Write(message)

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			log.Fatalf("Error sending message to server: %v\n", err)
		}
	}
}