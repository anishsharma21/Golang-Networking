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

const defaultPort uint16 = 8080

func main() {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", defaultPort))
	if err != nil {
		log.Fatalf("Error connecting to server: %v\n", err)
	}
	defer handleDisconnect(conn)
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	err = authenticateWithServer(conn, *scanner)
	if err != nil {
		log.Printf("Error authenticating with server: %v\n", err)
		return
	}

	fmt.Printf("Connected to server on port %d...\n", defaultPort)
	fmt.Print(">> ")

	for scanner.Scan() {
		message := []byte(scanner.Text())
		length := uint16(len(message))

		var buf bytes.Buffer

		err := binary.Write(&buf, binary.BigEndian, length)
		if err != nil {
			log.Printf("Error writing length to buffer: %v\n", err)
			return
		}

		_, err = buf.Write(message)
		if err != nil {
			log.Printf("Error writing message to buffer: %v\n", err)
			return
		}

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			log.Printf("Error sending message to server: %v\n", err)
			return
		}
		messageSentTime := time.Now()

		lengthBuff := make([]byte, 2)

		err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			fmt.Printf("Error setting read deadline: %v\n", err)
			return
		}
		_, err = conn.Read(lengthBuff)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading in length of message from server: %v\n", err)
			} else {
				log.Println("Server closed.")
			}
			return
		}
		log.Printf("Response received in: %v\n", time.Since(messageSentTime))

		responseLength := uint16(lengthBuff[0]) << 8 + uint16(lengthBuff[1])
		responseBuff := make([]byte, responseLength)

		err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			fmt.Printf("Error setting read deadline: %v\n", err)
			return
		}
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

func authenticateWithServer(conn net.Conn, scanner bufio.Scanner) error {
	fmt.Print("Private key: ")
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading key from input: %v", err)
		}
		return fmt.Errorf("no input received")
	}
	key := []byte(scanner.Text())
	length := uint16(len(key))
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, length)
	if err != nil {
		return fmt.Errorf("problem writing length to buffer: %v", err)
	}
	buf.Write(key)
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("problem sending authentication message to server: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	confirmation, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return fmt.Errorf("problem receiving confirmation message from server: %v", err)
	}
	if strings.TrimSpace(confirmation) == "Invalid key." {
		return fmt.Errorf("invalid key: '%s'", string(key))
	}

	return nil
}

func handleDisconnect(conn net.Conn) {
	log.Printf("Disconnected from %s\n", conn.RemoteAddr().String())
}