package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

// TODO broadcast functionality
// TODO improve concurrency
// TODO improve memory efficiency
// TODO tidy up code, plan this out a lil

const defaultPort uint16 = 8080
var clients = make(map[net.Conn]string)
var mu sync.Mutex
const privateKey string = "tcpinit8989"

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
		go handleClient(conn)
	}
}

func handleClient(client net.Conn) {
	defer client.Close()
	defer disconnectClient(client)

	log.Printf("Establishing connection to client %s...\n", client.RemoteAddr().String())

	lengthBuff := make([]byte, 2)
	_, err := client.Read(lengthBuff)
	if err != nil {
		if err != io.EOF {
			log.Printf("Error receiving length of message from %s: %v\n", client.RemoteAddr().String(), err)
			return
		}
		log.Printf("Disconnected from %s\n", client.RemoteAddr().String())
		return
	}

	messageLength := uint16(lengthBuff[0]) << 8 + uint16(lengthBuff[1])
	messageBuffer := make([]byte, messageLength)
	_, err = client.Read(messageBuffer)
	if err != nil {
		log.Printf("Error receiving message from %s: %v\n", client.RemoteAddr().String(), err)
		return
	}

	if string(messageBuffer) != privateKey {
		log.Printf("Failed to establish connection with %s: invalid connection key: '%s'\n", client.RemoteAddr().String(), string(messageBuffer))
		client.Write([]byte("Invalid key. Failed to authenticate with server."))
		return
	}

	log.Printf("Connected to %s.", client.RemoteAddr().String())
	_, err = client.Write([]byte("Connected to server."))
	if err != nil {
		log.Printf("Failed to send confirmation message to %s\n", client.RemoteAddr().String())
		return
	}

	mu.Lock()
	clients[client] = client.RemoteAddr().String()
	mu.Unlock()

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
			break
		}

		fmt.Printf("%s: %s, %d\n", client.RemoteAddr().String(), string(packet), messageLength)

		responseMessage := []byte(fmt.Sprintf("%v\n", packet))
		responseMessageLength := uint16(len(responseMessage))

		var buf bytes.Buffer
		err = binary.Write(&buf, binary.BigEndian, responseMessageLength)
		if err != nil {
			log.Printf("Error writing length to buffer: %v\n", err)
			break
		}
		buf.Write(responseMessage)

		_, err = client.Write(buf.Bytes())
		if err != nil {
			log.Printf("Error sending message to client %s: %v\n", client.RemoteAddr().String(), err)
			break
		}
	}
}

func disconnectClient(client net.Conn) {
	mu.Lock()
	delete(clients, client)
	mu.Unlock()

	log.Printf("Disconnected from %s\n", client.RemoteAddr().String())
}