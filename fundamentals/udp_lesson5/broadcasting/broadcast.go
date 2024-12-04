package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

const (
	BROADCAST_ADDRESS string = "255.255.255.255"
	DEFAULT_PORT uint16 = 8080
	WRITE_BUFFER int = 1024
)

func main() {
	broadcastReady := make(chan bool)
	var broadCastWg sync.WaitGroup

	fmt.Printf("Starting up broadcast...")
	broadCastWg.Add(1)
	go setupBroadcast(BROADCAST_ADDRESS, DEFAULT_PORT, broadcastReady, &broadCastWg)

	if !<-broadcastReady {
		fmt.Printf("Broadcast not started up :(")
	}
	close(broadcastReady)
	fmt.Printf("Broadcast started on port %d...", DEFAULT_PORT)

	broadCastWg.Wait()
	fmt.Printf("Broadcast shut down.")
}

func setupBroadcast(broadcastAddress string, defaultPort uint16, broadcastReady chan<- bool, broadCastWg *sync.WaitGroup) {
	defer broadCastWg.Done()
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", broadcastAddress, defaultPort))
	if err != nil {
		fmt.Printf("Error setting up broadcast: %v\n", err)
		broadcastReady <- false
		return
	}
	defer conn.Close()

	err = conn.(*net.UDPConn).SetWriteBuffer(WRITE_BUFFER)
	if err != nil {
		fmt.Printf("Error enabling broadcast: %v\n", err)
		broadcastReady <- false
		return
	}

	broadcastReady <- true

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Error broadcasting message %q: %v\n", message, err)
			return
		}
		fmt.Printf("Broadcast message sent.")
	}
}