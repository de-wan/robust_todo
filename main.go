package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/de-wan/robust_todo/db_sqlc"
	"github.com/de-wan/robust_todo/handlers"
	"github.com/joho/godotenv"
)

func loadEnv() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadEnv()

	serverPort := 3000
	serverPortString := os.Getenv("SERVER_PORT")
	if serverPortString != "" {
		serverPortParsed, err := strconv.ParseInt(serverPortString, 10, 64)
		if err != nil {
			log.Fatal("invalid variable SERVER_PORT in .env")
		}

		serverPort = int(serverPortParsed)
	}

	db_sqlc.Init()

	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("GET /", handlers.IndexHandler)
	http.HandleFunc("GET /add-todo", handlers.AddTodoViewHandler)
	http.HandleFunc("POST /add-todo", handlers.AddTodoCreateHandler)
	http.HandleFunc("PUT /toggle-todo/", handlers.ToggleTodoHandler)

	log.Printf("Starting server on port %d", serverPort)
	http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil)
}
