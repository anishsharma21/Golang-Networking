package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/anishsharma21/Golang-Networking/projects/todolist/internal/handlers"
)

func RunApp(port uint16) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlers.BasePageHandler)

	srv := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	log.Printf("Server starting on port %d...", port)
	return srv.ListenAndServe()
}