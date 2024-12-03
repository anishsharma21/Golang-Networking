package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const DEFAULT_PORT uint16 = 8080
const DATAGRAM_BUFFER_SIZE = 1024

func main() {
	serverReady := make(chan bool)
	var serverWg sync.WaitGroup
	serverWg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	go startServer(ctx, DEFAULT_PORT, serverReady, &serverWg)

	if !<-serverReady {
		log.Printf("Server not started :(")
		return
	}
	close(serverReady)
	log.Printf("UDP Server started on port %d...\n", DEFAULT_PORT)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	fmt.Println()
	log.Println("Server shutdown signal received...")
	cancel()

	serverWg.Wait()
	log.Println("Server shut down gracefully.")
}

func startServer(ctx context.Context, port uint16, ready chan<- bool, serverWg *sync.WaitGroup) {
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
		select {
		case <-ctx.Done():
			log.Println("Server is shutting down...")
			return
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				log.Printf("Error reading from %s: %v\n", addr, err)
				continue
			}
			go handlePacket(conn, addr, buffer[:n])
		}
	}
}

func handlePacket(conn net.PacketConn, addr net.Addr, data []byte) {
	log.Printf("%s: %q\n", addr, string(data))
	message := fmt.Sprintf("Length of message: %d\n", len(data))
	_, err := conn.WriteTo([]byte(message), addr)
	if err != nil {
		log.Printf("Error sending message to %s: %v\n", addr, err)
	}
}
