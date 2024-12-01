package main

import (
	"log"
	"math/rand"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln("Error connecting to server:", err)
	}

	log.Println("Connected to server:", conn.RemoteAddr().String())
	packet := make([]byte, 64)

	for {
		for i := 0; i < len(packet); i++ {
			packet[i] = byte(rand.Intn(256))
		}
		log.Printf("Packet: %v\n", packet)

		n, err := conn.Write(packet)
		if err != nil {
			log.Fatalf("Error sending %d bytes to server: %v\n", n, err)
		}

		_, err = conn.Write(packet)
		if err != nil {
			log.Fatalf("Error receiving message from server: %v", err)
		}

		log.Printf("Server packet: %v\n", packet)

		time.Sleep(3 * time.Second)
	}
}