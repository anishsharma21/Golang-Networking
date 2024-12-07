package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
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

	time.Sleep(100 * time.Millisecond)

	getAddresses()

	ip := net.ParseIP("127.0.0.1")
	fmt.Println("Parsed IP:", ip)
	if ip.To4() != nil {
		fmt.Println("IPv4 address")
	} else if ip.To16() != nil {
		fmt.Println("IPv6 address")
	} else {
		fmt.Println("Invalid IP type")
	}

	IPDetails(ip)
}

func getAddresses() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("Error getting system interface addresses: %v\n", err)
		return
	}
	for i, addr := range addrs {
		fmt.Printf("Address %d: %v\n", i+1, addr)
	}
}

func IPDetails(ip net.IP) {
	if ip.To4() != nil {
		fmt.Println("IPv4 address")
	} else if ip.To16() != nil {
		fmt.Println("IPv6 address")
	} else {
		fmt.Println("Invalid IP type")
	}
	ipSplit := strings.Split(ip.String(), ".")
	for _, str := range ipSplit {
		val, err := strconv.ParseInt(str, 10, 8)
		if err != nil {
			fmt.Printf("Error parsing %q: %v\n", str, err)
			return
		}
		fmt.Printf("%08b", val)
	}
	println()
}