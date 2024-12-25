package main

import (
	"fmt"
	"log"
	"net/http"
)

const port uint16 = 8080

type Content struct {
	Title string
	Heading string
}

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, world!")
	})

	fmt.Println("Starting HTTP server on port 8080...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}