package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/goodbye", goodbyeHandler)

	fmt.Println("Starting HTTP server on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("Error starting HTTP server: %v\n", err)
		return
	}
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, world!\n")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	name := queryParams.Get("name")
	if name != "" {
		fmt.Fprintf(w, "hello %s!\n", name)
	} else {
		fmt.Fprintf(w, "hello!\n")
	}
}

func goodbyeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "goodbye!\n")
}