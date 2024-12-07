package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	addr := net.JoinHostPort("127.0.0.1", "8080")
	fmt.Println(addr)
	addrs, err := net.LookupHost("example.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, addr := range addrs {
		fmt.Printf("Address and %d: %s\n", i+1, addr)
	}
	netip, pIp, err := net.ParseCIDR("192.0.2.0/24")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%v is type %T and %v is type %T\n", netip, netip, pIp, pIp)

	conn1, conn2 := net.Pipe()

	go func() {
		fmt.Fprintf(conn1, "hello from conn1")
	}()

	go func() {
		buffer := make([]byte, 64)
		n, err := conn2.Read(buffer)
		if err != nil {
			fmt.Printf("Error reading from conn2: %v\n", err)
			return
		}
		fmt.Printf("Recieved in conn2: %s\n", string(buffer[:n]))
	}()

	time.Sleep(1 * time.Second)
}