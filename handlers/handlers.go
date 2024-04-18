package handlers

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/de-wan/robust_todo/db_sqlc"
)

type TodoItem struct {
	Uuid  string
	Value string
}

type IndexData struct {
	Todos []db_sqlc.ListTodosRow
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index.html").ParseFiles("templates/pages/index.html", "templates/layouts/base.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	c := context.Background()
	queries := db_sqlc.New(db_sqlc.DB)
	todoItems, err := queries.ListTodos(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	indexData := IndexData{
		Todos: todoItems,
	}

	tmpl.ExecuteTemplate(w, "base", indexData)
}

func ToggleTodoHandler(w http.ResponseWriter, r *http.Request) {
	// todo: implement this
}
