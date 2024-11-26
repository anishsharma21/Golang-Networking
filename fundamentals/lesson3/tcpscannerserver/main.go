package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

const port uint16 = 8080

func main() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}
	defer ln.Close()

	log.Printf("Server started on port %d...\n", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		log.Println("Connected to client:", conn.RemoteAddr().String())
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := strings.TrimSpace(scanner.Text())
		log.Printf("%s: %s\n", conn.RemoteAddr().String(), message)
		log.Printf("%s bytes representation: %v\n", conn.RemoteAddr().String(), []byte(message))
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			log.Printf("Error sending message to %s: %v\n", conn.RemoteAddr().String(), err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading message from %s: %v\n", conn.RemoteAddr().String(), err)
	} else {
		log.Printf("%s disconnected.", conn.RemoteAddr().String())
	}
}