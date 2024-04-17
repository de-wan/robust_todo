package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/de-wan/robust_todo/handlers"
)

func main() {
	port := 3000

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", handlers.IndexHandler)

	log.Printf("Starting server on port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
