package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln("Error connecting to server:", err)
	}

	log.Printf("Connected to server: %v\n", conn.RemoteAddr().String())
	fmt.Print(">> ")
	inputScanner := bufio.NewScanner(os.Stdin)

	for inputScanner.Scan() {
		_, err := conn.Write([]byte(inputScanner.Text() + "\n"))
		if err != nil {
			log.Printf("Error sending message to client: %v\n", err)
			return
		}

		fmt.Println("Sent!")

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("Error receiving message from client: %v\n", err)
			return
		}

		fmt.Println("Server:", strings.TrimSpace(response))
		fmt.Print(">> ")
	}
}