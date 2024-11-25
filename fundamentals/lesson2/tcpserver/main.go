package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var serverErrCount int = 0

func main() {
	var port uint16 = 8080
	if len(os.Args) > 2 {
		portInt, err := strconv.ParseInt(os.Args[1], 10, 16)
		if err != nil {
			fmt.Printf("Error parsing %q: %e\n", os.Args[1], err)
			return
		}
		port = uint16(portInt)
	}
	tcpserver(port)
}

func tcpserver(port uint16) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("TCP Server started on port %d...\n", port)

	for {
		if serverErrCount > 3 {
			fmt.Println("Too many server errors... shutting down.")
			return
		}
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error connecting to client:", err)
			serverErrCount++
			continue
		}
		fmt.Println("Connected to client:", conn.RemoteAddr().String())
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("%s disconnected\n", conn.RemoteAddr().String())
			return
		}

		fmt.Printf("Message from %s: %s\n", conn.RemoteAddr().String(), strings.TrimSpace(message))

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Error sending message to %s: %e\n", conn.RemoteAddr().String(), err)
			return
		}
	}
}