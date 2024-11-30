package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

const port uint16 = 8080
var clients = make(map[net.Conn]string)
var mu sync.Mutex

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Chat server started on port %d...\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error connecting to client:", err)
			continue
		}
		go handleClientConnection(conn)
	}
}

func handleClientConnection(conn net.Conn) {
    defer conn.Close()
    fmt.Printf("Connected to client %s\n", conn.RemoteAddr().String())

    mu.Lock()
    clients[conn] = conn.RemoteAddr().String()
    mu.Unlock()

    for {
        message, err := bufio.NewReader(conn).ReadString('\n')
        if err != nil {
            fmt.Println("Client disconnected:", conn.RemoteAddr())
            break
        }

		fmt.Printf("Message from client %s: %s\n", conn.RemoteAddr().String(), strings.TrimSpace(message))

        broadcastMessage(conn, message)
    }

    mu.Lock()
    delete(clients, conn)
    mu.Unlock()
}

func broadcastMessage(sender net.Conn, message string) {
    mu.Lock()
    defer mu.Unlock()
    for client := range clients {
        if sender != client {
            fmt.Printf("Broadcasting message to %s\n", client.RemoteAddr().String())
            writer := bufio.NewWriter(client)
			message = fmt.Sprintf("%s: %s\n", client.RemoteAddr().String(), message)
            _, err := writer.WriteString(message)
            if err != nil {
                fmt.Printf("Error sending message from client %s to client %s: %v\n", sender.RemoteAddr().String(), client.RemoteAddr().String(), err)
                continue
            }
            err = writer.Flush()
            if err != nil {
                fmt.Printf("Error flushing message to client %s: %v\n", client.RemoteAddr().String(), err)
            }
        }
    }
}