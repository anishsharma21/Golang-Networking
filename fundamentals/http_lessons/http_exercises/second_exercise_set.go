package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Base path init.\n")
	})

	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("id")
		if userId != "" {
			fmt.Fprintf(w, "User ID: %s\n", userId)
		} else {
			fmt.Fprintf(w, "Could not parse user id...\n")
		}
	})

	log.Println("Starting HTTP server on port 8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Printf("Error starting HTTP server on port 8080: %v\n", err)
		return
	}
}