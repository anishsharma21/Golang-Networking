package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

const port uint16 = 8080

var tmpl = template.Must(template.ParseFiles("index.html"))

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", baseHandler)

	log.Printf("Server started on port %d...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html")
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
