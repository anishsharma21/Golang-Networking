package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no argument given - choose from server (s) or client (c)")
		return
	}
	if os.Args[1] == "s" {
		tcpserver()
	} else if os.Args[1] == "c" {
		tcpclient()
	} else {
		fmt.Println("invalid argument - choose from server (s) or client (c)")
	}
}

func tcpserver() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("TCP server started on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Client connected:", conn.RemoteAddr().String())
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

		clientData := string(buffer[:n])
		fmt.Println("Received from client:", clientData)
		conn.Write([]byte("Echo: " + clientData))
	}
}

func tcpclient() {

}