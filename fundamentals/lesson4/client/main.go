package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Error connecting to server: %v\n", err)
	}
	defer conn.Close()

	fmt.Println("Connected to server on port 8080...")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">> ")

	for scanner.Scan() {
		message := []byte(scanner.Text())
		length := uint16(len(message))

		var buf bytes.Buffer

		err := binary.Write(&buf, binary.BigEndian, length)
		if err != nil {
			log.Fatalf("Error writing length to buffer: %v\n", err)
		}

		buf.Write(message)

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			log.Fatalf("Error sending message to server: %v\n", err)
		}

		lengthBuff := make([]byte, 2)
		_, err = conn.Read(lengthBuff)
		if err != nil {
			if err != io.EOF {
				log.Fatalf("Error reading in length of message from server: %v\n", err)
			}
			break
		}

		responseLength := uint16(lengthBuff[0]) << 8 + uint16(lengthBuff[1])
		responseBuff := make([]byte, responseLength)
		_, err = conn.Read(responseBuff)
		if err != nil {
			log.Fatalf("Error receiving response from server: %v\n", err)
		}

		fmt.Printf("Server: %v\n", strings.TrimRight(string(responseBuff), "\r\n"))
		fmt.Print(">> ")
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from input: %v\n", err)
	}

	log.Printf("Disconnected from %s\n", conn.RemoteAddr().String())
}