package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := 3000
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Printf("Starting server on port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
