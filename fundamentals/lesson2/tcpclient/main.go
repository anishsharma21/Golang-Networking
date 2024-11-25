package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var port uint16 = 8080
	if len(os.Args) > 1 {
		portInt, err := strconv.ParseInt(os.Args[1], 10, 16)
		if err != nil {
			fmt.Printf("Error parsing %q: %e\n", os.Args[1], err)
			return
		}
		port = uint16(portInt)
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}

	fmt.Println("Connected to server!")

	for {
		fmt.Print(">> ")
		reader := bufio.NewReader(os.Stdin)

		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message:", err)
			continue
		}

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message to server:", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		responseChan := make(chan string)
		errorChan := make(chan error)

		go func() {
			response, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				errorChan <- err
				return
			}
			responseChan <- response
		}()

		select {
		case response := <- responseChan:
			fmt.Println("Server:", strings.TrimSpace(response))
		case err := <- errorChan:
			fmt.Println("Error reading response from server:", err)
			return
		case <-ctx.Done():
			fmt.Println("Response time out.")
			return
		}
	}
}