package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/ requested")
	fmt.Fprintf(w, "Hello, world!")
}

func main() {
	http.HandleFunc("/", helloHandler)
	fmt.Println("HTTP server starting on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting HTTP server on port 8080: %v\n", err)
		return
	}
}