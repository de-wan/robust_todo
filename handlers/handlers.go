package handlers

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/de-wan/robust_todo/db_sqlc"
	"github.com/de-wan/robust_todo/utils"
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

type AddTodoFormData struct {
	Todo string
}

type AddTodoFormErrors struct {
	Todo []string
}

type AddTodoData struct {
	ErrorMsg   string
	FormData   AddTodoFormData
	FormErrors AddTodoFormErrors
}

func AddTodoViewHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("add-todo.html").ParseFiles("templates/pages/add-todo.html", "templates/layouts/base.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// check for input errors
	templateName := "base"
	incomingTarget := r.Header.Get("HX-Target")

	if incomingTarget == "content" {
		templateName = "content"
	}

	tmpl.ExecuteTemplate(w, templateName, AddTodoData{})
}

func AddTodoCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}

	// create form template
	formTmpl, err := template.New("add-todo.html").ParseFiles("templates/pages/add-todo.html", "templates/layouts/base.html")
	if err != nil {
		addTodoData := AddTodoData{
			ErrorMsg: "Whoops... Something went wrong",
		}
		formTmpl.ExecuteTemplate(w, "content", addTodoData)
		return
	}

	// check for input error
	formData := AddTodoFormData{
		Todo: r.FormValue("value"),
	}

	addTodoData := AddTodoData{}

	if strings.TrimSpace(formData.Todo) == "" {
		addTodoData.ErrorMsg = "Check your form and try again"
		addTodoData.FormErrors.Todo = append(addTodoData.FormErrors.Todo, "This field is required")
	}

	if addTodoData.ErrorMsg != "" {
		formTmpl.ExecuteTemplate(w, "content", addTodoData)
		return
	}

	c := context.Background()

	// prepare and run insert query
	queries := db_sqlc.New(db_sqlc.DB)

	uuid := utils.GenerateUUID()
	err = queries.AddTodo(c, db_sqlc.AddTodoParams{
		Uuid:  uuid,
		Value: formData.Todo,
	})
	if err != nil {
		addTodoData.ErrorMsg = "Whoops!... Unable to create todo record."
		formTmpl.ExecuteTemplate(w, "content", addTodoData)
		return
	}

	// return redirect to todo list
	w.Header().Set("HX-Redirect", "/")
}

func ToggleTodoHandler(w http.ResponseWriter, r *http.Request) {
	// todo: implement this
}
