package handlers

import (
	"context"
	"net/http"

	"github.com/anishsharma21/Golang-Networking/projects/todolist/public/templates"
	"github.com/anishsharma21/Golang-Networking/projects/todolist/types"
)

func BasePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		todos := []types.Todo{}
		basePage := templates.Base(todos)
		err := basePage.Render(context.Background(), w)
		if err != nil {
			http.Error(w, "Unable to load the page", http.StatusInternalServerError)
		}
	}
}