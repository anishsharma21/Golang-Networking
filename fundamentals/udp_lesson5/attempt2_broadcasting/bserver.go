package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
    go startBroadcastResponseListener(8080)
    startBroadcastingClient()
}

func startBroadcastingClient() {
    broadcastAddr, err := net.ResolveUDPAddr("udp", "192.168.20.255:8080")
    if err != nil {
        fmt.Printf("Error resolving broadcast address: %v\n", err)
        return
    }

    localAddr, err := net.ResolveUDPAddr("udp", "192.168.20.210:0")
    if err != nil {
        fmt.Printf("Error resolving local address: %v\n", err)
        return
    }

    conn, err := net.DialUDP("udp", localAddr, broadcastAddr)
    if err != nil {
        fmt.Printf("Error creating UDP connection to %s: %v\n", broadcastAddr, err)
        return
    }
    defer conn.Close()

    message := []byte("Hello, subnet!")
    for {
        _, err := conn.Write(message)
        if err != nil {
            fmt.Printf("Error sending message to subnet: %v\n", err)
            return
        }
        log.Printf("Message sent to broadcast address: %s\n", broadcastAddr)
        time.Sleep(5 * time.Second)
    }
}

func startBroadcastResponseListener(port uint16) {
    addr := net.UDPAddr{
        Port: int(port),
        IP:   net.IPv4zero,
    }
    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Printf("Error setting up listener for broadcast responses: %v\n", err)
        return
    }
    defer conn.Close()

    buffer := make([]byte, 1024)
    for {
        n, remoteAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Printf("Error reading response: %v\n", err)
            return
        }
        fmt.Printf("Received response from %s: %s\n", remoteAddr, string(buffer[:n]))
    }
}