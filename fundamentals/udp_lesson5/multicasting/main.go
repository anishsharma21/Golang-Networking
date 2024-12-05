package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	group := "224.0.0.1:8080"

	addr, err := net.ResolveUDPAddr("udp", group)
	if err != nil {
		fmt.Printf("Error resolving multicast address %q: %v\n", group, err)
		return
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		fmt.Printf("Error joining multicast group: %v\n", err)
		return
	}
	defer conn.Close()
	fmt.Printf("Joined multicast group on port %s...\n", strings.Split(group, ":")[1])

	buffer := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading multicast message: %v\n", err)
			continue
		}
		fmt.Printf("Received multicast message from %s: %s\n", remoteAddr, string(buffer[:n]))
	}
}