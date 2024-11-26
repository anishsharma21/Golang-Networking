package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

const defaultPort uint16 = 8080

// TCP server that reads in packets rather than strings
func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", defaultPort))
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}
	defer listener.Close()

	fmt.Printf("Starting server on port %d...\n", defaultPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Error accepting connection:", err)
		}
		fmt.Println("Connected to client:", conn.RemoteAddr().String())
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	packet := make([]byte, 64)

	for {
		_, err := reader.Read(packet)
		if err != nil {
			if err == io.EOF {
				log.Printf("%s disconnected\n", conn.RemoteAddr().String())
			} else {
				log.Printf("Error reading from %s: %e\n", conn.RemoteAddr().String(), err)
			}
			return
		}

		fmt.Printf("Packet from %s: %x\n", conn.RemoteAddr().String(), packet)

		_, err = conn.Write(packet)
		if err != nil {
			log.Printf("Error sending packet to %s: %e\n", conn.RemoteAddr().String(), err)
		}
	}
}