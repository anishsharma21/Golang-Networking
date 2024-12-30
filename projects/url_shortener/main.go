package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"text/template"
	"time"
)

const port uint16 = 8080

var tmpl = template.Must(template.ParseFiles("index.html"))

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", baseHandler)
	mux.HandleFunc("POST /url-shorten", urlFormSubmitHandler)

	log.Printf("Server started on port %d...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

func urlFormSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		url := r.FormValue("url")
		fmt.Fprintf(w, "URL received: %s\n", url)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
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

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(result)
}