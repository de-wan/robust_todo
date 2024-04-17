package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/de-wan/robust_todo/handlers"
)

func main() {
	port := 3000

	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("GET /", handlers.IndexHandler)

	http.HandleFunc("PUT /toggle-todo", handlers.ToggleTodoHandler)

	log.Printf("Starting server on port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
