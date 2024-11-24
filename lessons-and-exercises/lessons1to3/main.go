package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

func main() {
	httpfetch()
}

func httpfetch() {
	response, err := http.Get("https://example.com")
	if err != nil {
		fmt.Println("Error fetching page data:", err)
		return
	}
	defer response.Body.Close()

	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
}

func tcpserver1(port int) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer ln.Close()

	fmt.Printf("Listening on port %d...\n", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		go tcpserver1_handleconnection(conn)
	}
}

func tcpserver1_handleconnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected.")

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Message from client:", message)
	conn.Write([]byte("Hello, client!\n"))
}

func tcpclient1() {
	conn, err := net.Dial("tcp", "example.com:80")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: example.com\r\n\r\n")

	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Server response:", response)
}