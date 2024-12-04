package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const DEFAULT_PORT int = 8080
const DATAGRAM_BUFFER_SIZE = 1024

func main() {
	conn, err := net.Dial("udp", fmt.Sprintf("localhost:%d", DEFAULT_PORT))
	if err != nil {
		fmt.Printf("Error connecting to server on port %d: %v\n", DEFAULT_PORT, err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to UDP server on port %d...\n", DEFAULT_PORT)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Error sending message %q to server %s: %v\n", message, conn.RemoteAddr().String(), err)
			return
		}

		buffer := make([]byte, DATAGRAM_BUFFER_SIZE)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Error receiving response from server %s: %v\n", conn.RemoteAddr().String(), err)
			return
		}
		fmt.Printf("Server: %s\n", strings.TrimSpace(string(buffer[:n])))
	}

	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			fmt.Printf("Error scanning input: %v\n", err)
			return
		}
		fmt.Printf("End of file encountered: %v\n", err)
	}
}
