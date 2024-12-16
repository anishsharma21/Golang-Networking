package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
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

	mux.HandleFunc("GET /data", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Here is some data")
	})

	mux.HandleFunc("POST /data", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Data saved")
	})

	mux.HandleFunc("GET /secure", secureHandler)

	mux.HandleFunc("POST /divide", handleDivide)

	mux.HandleFunc("GET /redirect", redirectHandler)

	loggedMux := LoggingMiddleware(mux)

	log.Println("Starting HTTP server on port 8080...")
	err := http.ListenAndServe(":8080", loggedMux)
	if err != nil {
		log.Printf("Error starting HTTP server on port 8080: %v\n", err)
		return
	}
}

type DivideData struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func handleDivide(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		bodyJson, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("handleDivide: error reading data from request body: %v\n", err)
			http.Error(w, "error reading request body", http.StatusInternalServerError)
			return
		}

		divideData := new(DivideData)
		err = json.Unmarshal(bodyJson, divideData)
		if err != nil {
			fmt.Printf("handleDivide: error parsing JSON data: %v\n", err)
		}

		if divideData.Y == 0 {
			log.Printf("handleDivide: cannot divide by 0\n")
			http.Error(w, "cannot divide by 0", http.StatusBadRequest)
			return
		}

		result := float64(divideData.X) / float64(divideData.Y)
		fmt.Fprintf(w, "%d / %d = %.4f\n", divideData.X, divideData.Y, result)
	}
}

const USERNAME string = "admin"
const PASSWORD string = "hehe"

func secureHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	if username == "" {
		log.Printf("secureHandler: username not provided\n")
		http.Error(w, "username not provided", http.StatusBadRequest)
		return
	} else if password == "" {
		log.Printf("secureHandler: password not provided\n")
		http.Error(w, "password not provided", http.StatusBadRequest)
		return
	}

	if username != USERNAME || password != PASSWORD {
		log.Printf("secureHanlder: username or password incorrect:\nusername given: %q, username required: %q\npassword given: %q, password required: %q\n", username, USERNAME, password, PASSWORD)
		http.Error(w, "username or password incorrect", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Authenticated! Welcome!\n")
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:8080/hello", http.StatusFound)
}