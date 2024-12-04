package main

import (
	"fmt"
	"net"
)

const (
    DEFAULT_PORT uint16 = 8080
    BUFFER_SIZE  int    = 1024
)

func main() {
    addr := net.UDPAddr{
        Port: int(DEFAULT_PORT),
        IP:   net.IPv4(0, 0, 0, 0),
    }

    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Printf("Error starting UDP server: %v\n", err)
        return
    }
    defer conn.Close()

    fmt.Printf("Listening for broadcast messages on port %d...\n", DEFAULT_PORT)

    buffer := make([]byte, BUFFER_SIZE)
    for {
        n, remoteAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Printf("Error reading from UDP connection: %v\n", err)
            continue
        }
        message := string(buffer[:n])
        fmt.Printf("Received message from %s: %s\n", remoteAddr, message)
    }
}