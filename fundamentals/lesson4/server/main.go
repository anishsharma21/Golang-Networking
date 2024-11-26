package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

// TODO custom protocol with length prefixed data
// TODO broadcast functionality
// TODO timeouts for read operations

const defaultPort uint16 = 8080

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", defaultPort))
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}

	log.Printf("TCP server started on port %d...\n", defaultPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error connecting to client:", err)
			continue
		}
		log.Println("Connected to client:", conn.RemoteAddr().String())
		go handleClient(conn)
	}
}

func handleClient(client net.Conn) {
	defer client.Close()
	packet := make([]byte, 8)

	for {
		_, err := client.Read(packet)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error receiving message from %s: %v\n", client.RemoteAddr().String(), err)
			}
			break
		}

		fmt.Printf("%s: %v\n", client.RemoteAddr().String(), packet)

		_, err = client.Write(packet)
		if err != nil {
			log.Printf("Error sending message to %s: %v\n", client.RemoteAddr().String(), err)
			break
		}

		copy(packet, make([]byte, len(packet)))
	}

	log.Printf("Disconnected from %s\n", client.RemoteAddr().String())
}