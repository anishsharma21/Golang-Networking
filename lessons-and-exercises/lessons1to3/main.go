package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	ready := make(chan bool)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		tcpserver1(8080, ready)
	}()

	<-ready

	tcpclient2()

	wg.Wait()
}

func tcpclient2() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}

	defer conn.Close()
	fmt.Println("Connected to server. Type your message:")

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		message, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Error reading message:", err)
			continue
		}

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from server response:", err)
			break
		}

		fmt.Println("Server response:", response)
	}
}

func tcpserver1(port int, ready chan<- bool) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer ln.Close()

	fmt.Printf("Listening on port %d...\n", port)
	ready <- true

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

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message from client:", err)
			return
		}

		fmt.Println("Message from client:", message)
		conn.Write([]byte("Hello, client!\n"))
	}
}

func httpfetch() {
	start := time.Now()
	response, err := http.Get("https://example.com")
	if err != nil {
		fmt.Println("Error fetching page data:", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Webpage content:")
	fmt.Println(string(body))
	fmt.Printf("Time elapsed: %.3fs\n", time.Since(start).Seconds())
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