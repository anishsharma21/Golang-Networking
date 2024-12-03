package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func startClient() {
	conn, err := net.Dial("udp", fmt.Sprintf(":%d", DEFAULT_PORT))
	if err != nil {
		fmt.Printf("Error connecting to server on port %d: %v\n", DEFAULT_PORT, err)
		return
	}
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Error sending message to server %s: %v\n", conn.RemoteAddr().String(), err)
			return
		}
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("Error receiving response from server %s: %v\n", conn.RemoteAddr().String(), err)
			return
		}
		fmt.Printf("Server response: %s\n", response)
	}
}
