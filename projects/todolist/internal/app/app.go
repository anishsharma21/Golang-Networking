package app

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/anishsharma21/Golang-Networking/projects/todolist/internal/handlers"
)

func RunApp(port uint16) error {
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