package main

import (
	"log"

	"github.com/anishsharma21/Golang-Networking/projects/todolist/internal/app"
)

func main() {
	err := app.RunApp(8080)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server shutting down...")
}