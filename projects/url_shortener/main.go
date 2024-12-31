package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"text/template"
	"time"
)

const port uint16 = 8080

var tmpl = template.Must(template.ParseGlob("templates/*.html"))
var urlMap = make(map[string]string)
var urlMapMu sync.Mutex

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", baseHandler)
	mux.HandleFunc("POST /shorten-url", urlFormSubmitHandler)
	mux.HandleFunc("GET /r/", redirectHandler)

	log.Printf("Server started on port %d...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		urlMapMu.Lock()
		if url, ok := urlMap[fmt.Sprintf("http://localhost:%d%s", port, r.URL.Path)]; ok {
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			log.Printf("Error (redirectHandler): http://localhost:%d%s not found in map", port, r.URL.Path)
			http.Error(w, "Original URL not found", http.StatusNotFound)
		}
		urlMapMu.Unlock()
	}
}

type ShortUrlResponseData struct {
	ShortenedUrl string
}

func urlFormSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		url := r.FormValue("url")
		shortUrl := fmt.Sprintf("http://localhost:%d/r/%s", port, randomString(5)) 
		urlMapMu.Lock()
		urlMap[shortUrl] = url
		urlMapMu.Unlock()

		data := ShortUrlResponseData{
			ShortenedUrl: shortUrl,
		}
		w.Header().Set("Content-Type", "text/html")
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "error rendering template", http.StatusInternalServerError)
		}

		fmt.Fprintf(w, "Shortened URL: %s\n", shortUrl)
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