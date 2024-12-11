package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, world!")
}

func greetHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Greetings from the server!")
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
    response := map[string]string{"message": "Hello, JSON!"}
    jsonData, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Error generating JSON", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)
}

func first() {
    http.HandleFunc("/", helloHandler)
    http.HandleFunc("/greet", greetHandler)
    http.HandleFunc("/json", jsonHandler)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    fmt.Println("HTTP server starting on port 8080...")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Printf("Error starting HTTP server on port 8080: %v\n", err)
        return
    }
}