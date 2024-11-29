package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

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
		packet := make([]byte, messageLength)

		_, err = client.Read(packet)
		if err != nil {
			log.Printf("Error receiving message from %s: %v\n", client.RemoteAddr().String(), err)
		}

		fmt.Printf("%s: %s, %d\n", client.RemoteAddr().String(), string(packet), messageLength)

		responseMessage := []byte(fmt.Sprintf("%v\n", packet))
		responseMessageLength := uint16(len(responseMessage))

		var buf bytes.Buffer
		err = binary.Write(&buf, binary.BigEndian, responseMessageLength)
		if err != nil {
			log.Fatalf("Error writing length to buffer: %v\n", err)
		}
		buf.Write(responseMessage)

		_, err = client.Write(buf.Bytes())
		if err != nil {
			log.Fatalf("Error sending message to client %s: %v\n", client.RemoteAddr().String(), err)
		}
	}

	log.Printf("Disconnected from %s\n", client.RemoteAddr().String())
}