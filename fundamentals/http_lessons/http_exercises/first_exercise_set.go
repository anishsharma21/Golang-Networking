package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("form.html"))

func first() {
	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/goodbye", goodbyeHandler)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/greet", greetHandler)
	http.HandleFunc("/users/{id}", idHandler)

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

type HelloMessage struct {
	Message string `json:"message"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	jsonData, err := json.Marshal(HelloMessage{Message: "hello there"})
	if err != nil {
		log.Printf("Error encoding JSON: %v\n", err)
		return
	}

	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(jsonData)
		if err != nil {
			fmt.Printf("Error sending json data: %v\n%v\n", err, jsonData)
			return
		}
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

func greetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}
		name := r.FormValue("name")
		if name != "" {
			fmt.Fprintf(w, "hello %s!\n", name)
		} else {
			fmt.Fprintf(w, "hello!\n")
		}
	}
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Fragment)
	fmt.Println(r.URL.Path)
	fmt.Println(r.URL.Query())
	fmt.Println(r.URL.EscapedPath())
}