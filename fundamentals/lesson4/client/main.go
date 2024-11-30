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
	"time"
)

// TODO clear current line when other client messages broadcasted with this: fmt.Print("\r\033[K"), \r to go to start of line, other ANSI escape code to clear the line

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Error connecting to server: %v\n", err)
	}
	defer disconnect(conn)
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Private key: ")
	scanner.Scan()
	key := []byte(scanner.Text())
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from input: %v\n", err)
		return
	}
	length := uint16(len(key))
	var buf bytes.Buffer
	err = binary.Write(&buf, binary.BigEndian, length)
	if err != nil {
		log.Printf("Error writing length to buffer: %v\n", err)
		return
	}
	buf.Write(key)
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		log.Printf("Error sending authentication message to server: %v\n", err)
		return
	}

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	confirmation, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Printf("Error receiving confirmation message from server: %v\n", err)
		return
	}
	if strings.TrimSpace(confirmation) == "Invalid key." {
		log.Printf("Invalid key. Failed to connect to server.")
		return
	}

	fmt.Println("Connected to server on port 8080...")
	fmt.Print(">> ")
	// fmt.Print("\r\033[K")

	for scanner.Scan() {
		message := []byte(scanner.Text())
		length := uint16(len(message))

		var buf bytes.Buffer

		err := binary.Write(&buf, binary.BigEndian, length)
		if err != nil {
			log.Printf("Error writing length to buffer: %v\n", err)
			return
		}

		buf.Write(message)

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			log.Printf("Error sending message to server: %v\n", err)
			return
		}
		messageSentTime := time.Now()

		lengthBuff := make([]byte, 2)

		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		_, err = conn.Read(lengthBuff)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading in length of message from server: %v\n", err)
			}
			return
		}
		log.Printf("Response received in: %v\n", time.Since(messageSentTime))

		responseLength := uint16(lengthBuff[0]) << 8 + uint16(lengthBuff[1])
		responseBuff := make([]byte, responseLength)

		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		_, err = conn.Read(responseBuff)
		if err != nil {
			log.Printf("Error receiving response from server: %v\n", err)
			return
		}

		fmt.Printf("Server: %v\n", strings.TrimRight(string(responseBuff), "\r\n"))
		fmt.Print(">> ")
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from input: %v\n", err)
		return
	}
}

func disconnect(conn net.Conn) {
	log.Printf("Disconnected from %s\n", conn.RemoteAddr().String())
}