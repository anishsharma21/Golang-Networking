package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("form.html"))

func main() {
	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/goodbye", goodbyeHandler)
	http.HandleFunc("/form", formHandler)

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

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html")
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}