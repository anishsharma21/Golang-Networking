package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	var port uint16 = 8080
	if len(os.Args) >= 2 {
		portInt, err := strconv.ParseInt(os.Args[1], 10, 16)
		if err != nil {
			fmt.Printf("Error parsing %s: %e\n", os.Args[1], err)
			return
		}
		port = uint16(portInt)
	}
	fmt.Println("Port:", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	stopChan := make(chan struct{})

	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal, gracefully shutting down...")
		close(stopChan)
	}()

	tcpserver(port, stopChan)
}

func tcpserver(port uint16, stopChan <-chan struct{}) {
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        fmt.Println("Error starting TCP server:", err)
        return
    }
    defer listener.Close()
    fmt.Printf("TCP server started on port %s...\n", listener.Addr().String())

    for {
        select {
        case <-stopChan:
            fmt.Println("Server is shutting down...")
            return
        default:
			listener.(*net.TCPListener).SetDeadline(time.Now().Add(1 * time.Second))
            conn, err := listener.Accept()
            if err != nil {
                if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
                    continue
                }
                fmt.Println("Error accepting connection:", err)
                return
            }
            go handleClient(conn)
        }
    }
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Connection closed by client:", err)
			return
		}
		clientData := string(buffer[:n-1])
		fmt.Println("Received from client:", clientData)
		conn.Write([]byte("Echo: " + reverse(&clientData) + "\n>> "))
	}
}

func reverse(s *string) string {
	runes := []rune(*s)
	for i, j := 0, len(*s) - 1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}