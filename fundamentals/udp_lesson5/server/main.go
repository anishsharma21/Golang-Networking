package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

const DEFAULT_PORT uint16 = 8080
const DATAGRAM_BUFFER_SIZE = 1024

func main() {
	serverReady := make(chan bool)
	var serverWg sync.WaitGroup
	serverWg.Add(1)
	go startServer(DEFAULT_PORT, serverReady, &serverWg)
	if !<-serverReady {
		log.Printf("Server not started :(")
		return
	}
	close(serverReady)
	fmt.Printf("UDP Server started on port %d...\n", DEFAULT_PORT)
	serverWg.Wait()
	fmt.Println("Server shutting down.")
}

func startServer(port uint16, ready chan<- bool, serverWg *sync.WaitGroup) {
	defer serverWg.Done()
	conn, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("Error starting UDP server on port %d: %v\n", port, err)
		ready <- false
		serverWg.Done()
		return
	}
	defer conn.Close()
	ready <- true
	buffer := make([]byte, DATAGRAM_BUFFER_SIZE)
	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Printf("Error reading from %s: %v\n", addr, err)
			return
		}
		log.Printf("%s: %q\n", addr, string(buffer[:n]))
		message := "Hello, " + string(buffer[:n])
		_, err = conn.WriteTo([]byte(message), addr)
		if err != nil {
			log.Printf("Error sending message to %s: %v\n", addr, err)
			return
		}
	}
}