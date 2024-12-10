package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")

	userAgent := r.UserAgent()

	fmt.Fprintf(w, "Hello %s\n", name)
	fmt.Fprintf(w, "Hello %s\n", userAgent)

	response := Response{Message: "hello there!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	json.NewEncoder(w).Encode(response)
}