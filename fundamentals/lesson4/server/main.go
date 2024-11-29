package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

// TODO create TCP client which implements 2-byte data size prefix custom protocol
// TODO broadcast functionality
// TODO timeouts for read operations
// TODO first message is to check that client understands protocol
// TODO improve concurrency

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
	lengthBuff := make([]byte, 2)

	for {
		_, err := client.Read(lengthBuff)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error receiving length of message from %s: %v\n", client.RemoteAddr().String(), err)
			}
			break
		}

		messageLength := uint16(lengthBuff[0]) << 8 + uint16(lengthBuff[1])
		fmt.Println("Size of message:", messageLength)

		packet := make([]byte, messageLength)
		_, err = client.Read(packet)
		if err != nil {
			log.Printf("Error receiving message from %s: %v\n", client.RemoteAddr().String(), err)
		}

		fmt.Printf("%s: %s\n", client.RemoteAddr().String(), string(packet))

		_, err = client.Write([]byte("8 bytes:" + string(packet)))
		if err != nil {
			log.Printf("Error sending message to %s: %v\n", client.RemoteAddr().String(), err)
			break
		}

		copy(packet, make([]byte, len(packet)))
	}

	log.Printf("Disconnected from %s\n", client.RemoteAddr().String())
}