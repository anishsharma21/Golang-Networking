package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	tcpclient1()
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