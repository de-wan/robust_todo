package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type TodoItem struct {
	IsDone bool
	Value  string
}

type IndexData struct {
	Todos []TodoItem
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index.html").ParseFiles("templates/pages/index.html", "templates/layouts/base.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	indexData := IndexData{
		Todos: []TodoItem{
			{
				Value: "First Todo",
			}, {
				Value: "Second Todo",
			},
		},
	}

	tmpl.ExecuteTemplate(w, "base", indexData)
}

func ToggleTodoHandler(w http.ResponseWriter, r *http.Request) {
	// todo: implement this
}
