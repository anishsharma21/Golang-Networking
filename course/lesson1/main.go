package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
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
	tcpserver(port)
}

func tcpserver(port uint16) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()
	fmt.Printf("TCP server started on port %s...\n", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Client connected:", conn.RemoteAddr().String())
		_, err = conn.Write([]byte("WELCOME!\n>> "))
		if err != nil {
			fmt.Println("Error sending welcome message to client:", err)
			conn.Close()
			continue
		}
		go handleClient(conn)
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