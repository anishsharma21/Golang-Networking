package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
)

// TODO improve concurrency
// TODO tidy up code

var globalId int = 0
const NEWLINE_LENGTH = 1
const DEFAULT_PORT uint16 = 8080
var clients = make(map[net.Conn]int)
var mu sync.Mutex
const privateKey string = "tcpinit8989"

func main() {
	var serverReadyWg sync.WaitGroup
	var serverDownWg sync.WaitGroup

	serverReadyWg.Add(1)
	go startServer(DEFAULT_PORT, &serverReadyWg, &serverDownWg)
	serverReadyWg.Wait()

	log.Printf("TCP server started on port %d...\n", DEFAULT_PORT)

	serverDownWg.Add(1)
	serverDownWg.Wait()

	log.Printf("TCP server shutting down.")
}

func startServer(port uint16, serverReadyWg *sync.WaitGroup, serverDownWg *sync.WaitGroup) {
	defer serverDownWg.Done()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}
	defer listener.Close()

	serverReadyWg.Done()

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
	clients[client] = generateId()
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

		if _, err = client.Read(packet); err != nil {
			log.Printf("Error receiving message from %s: %v\n", client.RemoteAddr().String(), err)
			return
		}

		log.Printf("%s: %s, %d\n", client.RemoteAddr().String(), string(packet), messageLength)

		mu.Lock()
		clientId, ok := clients[client]
		if !ok {
			fmt.Printf("Error: client not saved globally and not found\n")
			return
		}
		mu.Unlock()

		clientIdStr := strconv.Itoa(clientId)
		responseMessageLength := uint16(len(packet) + NEWLINE_LENGTH + len(clientIdStr))

		var buf bytes.Buffer
		err = binary.Write(&buf, binary.BigEndian, responseMessageLength)
		if err != nil {
			log.Printf("Error writing length to buffer: %v\n", err)
			return
		}
		_, err = buf.Write([]byte(clientIdStr))
		if err != nil {
			log.Printf("Error writing client id '%s' to buffer: %v\n", clientIdStr, err)
			return
		}
		_, err = buf.Write([]byte("\n"))
		if err != nil {
			log.Printf("Error writing client id '%s' to buffer: %v\n", clientIdStr, err)
			return
		}
		_, err = buf.Write(packet)
		if err != nil {
			log.Printf("Error writing message to buffer: %v\n", err)
			return
		}
		broadcastMessage(buf.Bytes(), client)
	}
}

func broadcastMessage(message []byte, client net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	for conn, id := range clients {
		if client == conn {
			continue
		}
		if _, err := conn.Write(message); err != nil {
			log.Printf("Error sending message to client %d: %v\n", id, err)
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

func generateId() int {
	globalId++
	return globalId
}