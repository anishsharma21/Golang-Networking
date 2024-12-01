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
	"sync"
	"time"
)

// TODO channels for checking that messages with specific ID's have been responded too

const PROTOCOL_PARAMETERS_NUM int = 2
const DEFAULT_PORT uint16 = 8080
var printMu sync.Mutex

func main() {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", DEFAULT_PORT))
	if err != nil {
		log.Fatalf("Error connecting to server on port %d: %v\n", DEFAULT_PORT, err)
	}
	defer handleDisconnect(conn)
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	err = authenticateWithServer(conn, *scanner)
	if err != nil {
		log.Printf("Error authenticating with server: %v\n", err)
		return
	}

	fmt.Printf("Connected to server on port %d...\n", DEFAULT_PORT)
	fmt.Print(">> ")

	responseChan := make(chan string)
	defer close(responseChan)

	go listenForServerMessages(conn, responseChan)
	go printMessages(responseChan)

	for scanner.Scan() {
		message := []byte(scanner.Text())
		length := uint16(len(message))

		printMu.Lock()
		// clear current line in display
		fmt.Print("\033[A\r\033[K")
		fmt.Println("You:", string(message))
		printMu.Unlock()

		var buf bytes.Buffer

		err := binary.Write(&buf, binary.BigEndian, length)
		if err != nil {
			log.Printf("Error writing length to buffer: %v\n", err)
			return
		}

		if _, err = buf.Write(message); err != nil {
			log.Printf("Error writing message to buffer: %v\n", err)
			return
		}

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			log.Printf("Error sending message to server: %v\n", err)
			return
		}
		fmt.Print(">> ")
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from input: %v\n", err)
		return
	}
}

func listenForServerMessages(conn net.Conn, serverChannel chan<- string) {
	lengthBuff := make([]byte, 2)

	for {
		_, err := conn.Read(lengthBuff)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading in length of message from server: %v\n", err)
			} else {
				log.Println("Server closed.")
			}
			return
		}

		responseLength := uint16(lengthBuff[0]) << 8 + uint16(lengthBuff[1])
		responseBuff := make([]byte, responseLength)

		_, err = conn.Read(responseBuff)
		if err != nil {
			log.Printf("Error receiving response from server: %v\n", err)
			return
		}

		serverChannel <- strings.TrimRight(string(responseBuff), "\r\n")
	}
}

func printMessages(serverChannel <-chan string) {
	for {
		serverMessage := <- serverChannel
		serverMessageParts, err := parseServerMessage(serverMessage)
		if err != nil {
			log.Printf("Error parsing server message: %v\n", err)
			os.Exit(1)
			// TODO should have a wait group for client to close it down gracefully
		}
		printMu.Lock()
		fmt.Print("\r\033[K")
		fmt.Printf("Client %s: %s\n", serverMessageParts[0], serverMessageParts[1])
		fmt.Print(">> ")
		printMu.Unlock()
	}
}

func parseServerMessage(serverMessage string) ([]string, error) {
	serverMessageParts := strings.Split(serverMessage, "\n")
	if len(serverMessageParts) != PROTOCOL_PARAMETERS_NUM {
		return []string{}, fmt.Errorf("server message has wrong number of protocol parameters (should be %d): %s", PROTOCOL_PARAMETERS_NUM, serverMessage)
	}
	return serverMessageParts, nil
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
	conn.SetReadDeadline(time.Time{})
	if strings.TrimSpace(confirmation) == "Invalid key." {
		return fmt.Errorf("invalid key: '%s'", string(key))
	}

	return nil
}

func handleDisconnect(conn net.Conn) {
	log.Printf("Disconnected from %s\n", conn.RemoteAddr().String())
}