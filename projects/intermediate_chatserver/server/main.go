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

// TODO testing
// TODO broadcast functionality
// TODO improve concurrency

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
	defer disconnectClient(client)
	defer client.Close()
	lengthBuff := make([]byte, 2)

	log.Printf("Establishing connection to client %s...\n", client.RemoteAddr().String())
	err := establishConnection(client, lengthBuff)
	if err != nil {
		log.Printf("Error establishing connection to %s: %v\n", client.RemoteAddr().String(), err)
	}
	log.Printf("Connected to %s", client.RemoteAddr().String())

	mu.Lock()
	clients[client] = client.RemoteAddr().String()
	mu.Unlock()

	for {
		_, err := client.Read(lengthBuff)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error receiving length of message from %s: %v\n", client.RemoteAddr().String(), err)
			}
			return
		}

		messageLength := uint16(lengthBuff[0]) << 8 + uint16(lengthBuff[1])
		packet := make([]byte, messageLength)

		_, err = client.Read(packet)
		if err != nil {
			log.Printf("Error receiving message from %s: %v\n", client.RemoteAddr().String(), err)
			return
		}

		log.Printf("%s: %s, %d\n", client.RemoteAddr().String(), string(packet), messageLength)

		responseMessage := []byte(fmt.Sprintf("%v\n", packet))
		responseMessageLength := uint16(len(responseMessage))

		var buf bytes.Buffer
		err = binary.Write(&buf, binary.BigEndian, responseMessageLength)
		if err != nil {
			log.Printf("Error writing length to buffer: %v\n", err)
			return
		}
		_, err = buf.Write(responseMessage)
		if err != nil {
			log.Printf("Error writing message to buffer: %v\n", err)
			return
		}
		_, err = client.Write(buf.Bytes())
		if err != nil {
			log.Printf("Error sending message to client %s: %v\n", client.RemoteAddr().String(), err)
			return
		}
	}
}

func establishConnection(client net.Conn, lengthBuff []byte) error {
	_, err := client.Read(lengthBuff)
	if err != nil {
		if err != io.EOF {
			return fmt.Errorf("length of message not received from %s: %v", client.RemoteAddr().String(), err)
		}
		return fmt.Errorf("disconnected from %s", client.RemoteAddr().String())
	}

	messageLength := uint16(lengthBuff[0]) << 8 + uint16(lengthBuff[1])
	messageBuffer := make([]byte, messageLength)
	_, err = client.Read(messageBuffer)
	if err != nil {
		return fmt.Errorf("issue receiving message from %s: %v", client.RemoteAddr().String(), err)
	}

	if string(messageBuffer) != privateKey {
		client.Write([]byte("Invalid key.\n"))
		return fmt.Errorf("invalid connection key: '%s'", string(messageBuffer))
	}

	_, err = client.Write([]byte("AUTHENTICATED\n"))
	if err != nil {
		return fmt.Errorf("failed to send confirmation message to %s", client.RemoteAddr().String())
	}

	return nil
}

func disconnectClient(client net.Conn) {
	mu.Lock()
	delete(clients, client)
	mu.Unlock()

	log.Printf("Disconnected from %s\n", client.RemoteAddr().String())
}