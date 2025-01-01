package handlers

import (
	"context"
	"net/http"

	"github.com/anishsharma21/Golang-Networking/projects/todolist/public/templates"
)

func BasePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		basePage := templates.Base("hi")
		err := basePage.Render(context.Background(), w)
		if err != nil {
			http.Error(w, "Unable to load the page", http.StatusInternalServerError)
		}
	}
}