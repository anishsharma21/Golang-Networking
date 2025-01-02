package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/anishsharma21/Golang-Networking/projects/todolist/internal/handlers"
)

var db *sql.DB
var dbMu sync.Mutex

func RunApp(port uint16) error {
	var err error
	db, err := sql.Open("sqlite3", "./db/todolist.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = createTable()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// IMPORTANT FOR CSS
	staticDir := http.Dir(filepath.Join("public", "css"))
	mux.Handle("GET /css/", http.StripPrefix("/css/", http.FileServer(staticDir)))

	mux.HandleFunc("GET /", handlers.BasePageHandler)

	srv := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	log.Printf("Server starting on port %d...\n", port)
	return srv.ListenAndServe()
}

func createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS todo (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		done BOOLEAN NOT NULL
	);`

	_, err := db.Exec(query)
	return err
}