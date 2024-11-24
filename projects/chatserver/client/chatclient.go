package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const port uint16 = 8080
// TODO client struct so that usernames can be presented in terminal instead

func main() {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server!")

	go func() {
		for {
			response, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println("Error receiving message from server:", err)
				return // NOTE might need to make this a continue instead
			}
			fmt.Printf("%s\n", strings.TrimSpace(response))
			fmt.Print(">> ")
		}
	}()

	for {
		fmt.Printf(">> ")
		reader := bufio.NewReader(os.Stdin)

		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading client message:", err)
			continue
		}

		message = strings.TrimSpace(message)

		if message == "quit" {
			break
		}

		_, err = conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("Error sending message to server:", err)
			break
		}
	}
}