package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
    BROADCAST_ADDRESS string = "255.255.255.255"
    DEFAULT_PORT uint16 = 8080
    WRITE_BUFFER int = 1024
)

func main() {
    broadcastReady := make(chan bool)
    var broadCastWg sync.WaitGroup
    ctx, cancel := context.WithCancel(context.Background())

    fmt.Printf("Starting up broadcast...\n")
    broadCastWg.Add(1)
    go setupBroadcast(BROADCAST_ADDRESS, DEFAULT_PORT, broadcastReady, &broadCastWg, ctx)

    if !<-broadcastReady {
        fmt.Printf("Broadcast not started up :(\n")
    }
    close(broadcastReady)
    fmt.Printf("Broadcast started on port %d...\n", DEFAULT_PORT)

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigChan
        println()
        fmt.Printf("Broadcast shutdown signal received...\n")
        cancel()
    }()

    broadCastWg.Wait()
    fmt.Printf("Broadcast shut down gracefully.\n")
}

func setupBroadcast(broadcastAddress string, defaultPort uint16, broadcastReady chan<- bool, broadCastWg *sync.WaitGroup, ctx context.Context) {
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

    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    reader := bufio.NewReader(os.Stdin)
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Context cancelled, shutting down broadcast...\n")
            return
        case <-ticker.C:
            if reader.Buffered() > 0 {
                message, err := reader.ReadString('\n')
                if err != nil {
                    fmt.Printf("Error reading message: %v\n", err)
                    return
                }
                _, err = conn.Write([]byte(message))
                if err != nil {
                    fmt.Printf("Error broadcasting message %q: %v\n", message, err)
                    return
                }
                fmt.Printf("Broadcast message sent.\n")
            }
        }
    }
}