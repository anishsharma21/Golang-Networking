package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"text/template"
	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const port uint16 = 8080

var tmpl = template.Must(template.ParseGlob("templates/*.html"))
var db *sql.DB
var dbMu sync.Mutex

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./url_shortener.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = createTable()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", baseHandler)
	mux.HandleFunc("POST /shorten-url", urlFormSubmitHandler)
	mux.HandleFunc("GET /r/", redirectHandler)

	log.Printf("Server started on port %d...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

func createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		short_url TEXT NOT NULL,
		original_url TEXT NOT NULL
	);`
	_, err := db.Exec(query)
	return err
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		shortUrl := fmt.Sprintf("http://localhost:%d%s", port, r.URL.Path)
		var originalUrl string

		dbMu.Lock()
		err := db.QueryRow("SELECT original_url FROM urls WHERE short_url = ?", shortUrl).Scan(&originalUrl)
		dbMu.Unlock()

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Original URL not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		http.Redirect(w, r, originalUrl, http.StatusFound)
	}
}

type ShortUrlResponseData struct {
	ShortenedUrl string
}

func urlFormSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		url := r.FormValue("url")
		shortUrl := fmt.Sprintf("http://localhost:%d/r/%s", port, randomString(5)) 

		dbMu.Lock()
		_, err := db.Exec("INSERT INTO urls (short_url, original_url) VALUES (?, ?)", shortUrl, url)
		dbMu.Unlock()

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		data := ShortUrlResponseData{
			ShortenedUrl: shortUrl,
		}
		w.Header().Set("Content-Type", "text/html")
		if err := tmpl.ExecuteTemplate(w, "url-response.html", data); err != nil {
			http.Error(w, "error rendering template", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html")
		err := tmpl.ExecuteTemplate(w, "index.html", nil)
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