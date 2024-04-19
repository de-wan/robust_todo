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
	AlertMsg   string
	AlertClass string
	Todos      []db_sqlc.ListTodosRow
}

func renderTodoList(w http.ResponseWriter, templateName string) {
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

	tmpl.ExecuteTemplate(w, templateName, indexData)
}

type AlertData struct {
	AlertMsg   string
	AlertClass string
}

// alert types:
//  0. success
//  1. info
//  2. warning
//  3. error
func renderAlert(w http.ResponseWriter, alert string, layout string, alertType int) (err error) {
	if alertType > 3 {
		alertType = 1
	}

	alertClass := ""
	if alertType == 0 {
		// success
		alertClass = "bg-green-100 border border-green-400 text-green-700 "
	} else if alertType == 1 {
		// info
		alertClass = "bg-blue-100 border border-blue-400 text-blue-700 "
	} else if alertType == 2 {
		// warning
		alertClass = "bg-yellow-100 border border-yellow-400 text-yellow-700 "
	} else {
		// error
		alertClass = "bg-red-100 border border-red-400 text-red-700 "
	}

	alertData := AlertData{
		AlertMsg:   alert,
		AlertClass: alertClass,
	}

	if layout == "" {
		layout = "base"
	}

	var tmpl *template.Template
	if layout == "base" {
		tmpl, err = template.New("layout.html").ParseFiles("templates/layouts/base.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return err
		}
	}

	w.Header().Set("Hx-Retarget", "#alert")
	tmpl.ExecuteTemplate(w, "alert", alertData)

	return nil
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// determining whether to render whole page or just the content
	templateName := "base"
	incomingTarget := r.Header.Get("HX-Target")

	if incomingTarget == "content" {
		templateName = "content" // sets the render to content
	}
	renderTodoList(w, templateName)
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

	// determining whether to render whole page or just the content
	templateName := "base"
	incomingTarget := r.Header.Get("HX-Target")

	if incomingTarget == "content" {
		templateName = "content" // sets the render to content
	}

	tmpl.ExecuteTemplate(w, templateName, AddTodoData{})
}

func AddTodoActionHandler(w http.ResponseWriter, r *http.Request) {
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

	// select and execute index template
	renderTodoList(w, "content")
}

type EditTodoFormData struct {
	Uuid         string
	Todo         string
	ContinueEdit string
}

type EditTodoFormErrors struct {
	Todo []string
}

type EditTodoData struct {
	AlertMsg   string
	AlertClass string
	FormData   EditTodoFormData
	FormErrors EditTodoFormErrors
}

func EditTodoViewHandler(w http.ResponseWriter, r *http.Request) {
	// get uuid from url
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	uuid := parts[2]

	tmpl, err := template.New("edit-todo.html").ParseFiles("templates/pages/edit-todo.html", "templates/layouts/base.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// determining whether to render whole page or just the content
	templateName := "base"
	incomingTarget := r.Header.Get("HX-Target")

	if incomingTarget == "content" {
		templateName = "content" // sets the render to content
	}

	// preparing db query
	c := context.Background()
	queries := db_sqlc.New(db_sqlc.DB)

	todoItem, err := queries.GetTodo(c, uuid)
	if err != nil {
		renderAlert(w, "Whoops!...Something went wrong", "base", 3)
		log.Println(err)
		return
	}

	editTodoData := EditTodoData{
		FormData: EditTodoFormData{
			Uuid: uuid,
			Todo: todoItem.Value,
		},
	}

	tmpl.ExecuteTemplate(w, templateName, editTodoData)
}

func EditTodoActionHandler(w http.ResponseWriter, r *http.Request) {
	// get uuid from url
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	uuid := parts[2]

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}

	// create form template
	formTmpl, err := template.New("edit-todo.html").ParseFiles("templates/pages/edit-todo.html", "templates/layouts/base.html")
	if err != nil {
		editTodoData := EditTodoData{
			AlertMsg:   "Whoops... Something went wrong",
			AlertClass: "bg-red-100 border border-red-400 text-red-700 ",
		}
		formTmpl.ExecuteTemplate(w, "content", editTodoData)
		return
	}

	// check for input error
	formData := EditTodoFormData{
		Todo: r.FormValue("value"),
	}

	editTodoData := EditTodoData{}

	if strings.TrimSpace(formData.Todo) == "" {
		editTodoData.AlertMsg = "Check your form and try again"
		editTodoData.AlertClass = "bg-red-100 border border-red-400 text-red-700 "
		editTodoData.FormErrors.Todo = append(editTodoData.FormErrors.Todo, "This field is required")
	}

	if editTodoData.AlertMsg != "" {
		formTmpl.ExecuteTemplate(w, "content", editTodoData)
		return
	}

	// prepare and run update query
	c := context.Background()
	queries := db_sqlc.New(db_sqlc.DB)

	err = queries.EditTodo(c, db_sqlc.EditTodoParams{
		Uuid:  uuid,
		Value: formData.Todo,
	})
	if err != nil {
		editTodoData.AlertMsg = "Whoops!... Unable to update todo record."
		editTodoData.AlertClass = "bg-red-100 border border-red-400 text-red-700 "
		formTmpl.ExecuteTemplate(w, "content", editTodoData)
		return
	}

	// select and execute index template
	renderTodoList(w, "content")
}

func ToggleTodoHandler(w http.ResponseWriter, r *http.Request) {
	// get uuid from url
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	uuid := parts[2]

	// preparing db query
	c := context.Background()
	queries := db_sqlc.New(db_sqlc.DB)

	// toggling todo
	err := queries.ToggleTodo(c, uuid)
	if err != nil {
		renderAlert(w, "Whoops!...Something went wrong", "base", 3)
		log.Println(err)
		return
	}

	// retrieving todo to display
	todoItem, err := queries.GetTodo(c, uuid)
	if err != nil {
		renderAlert(w, "Whoops!...Something went wrong", "base", 3)
		log.Println(err)
		return
	}

	tmpl, err := template.New("todo-item.html").ParseFiles("templates/pages/index.html")
	if err != nil {
		renderAlert(w, "Whoops!...Something went wrong", "base", 3)
		log.Println(err)
		return
	}

	tmpl.ExecuteTemplate(w, "todo-item", todoItem)
}

func ArchiveTodoActionHandler(w http.ResponseWriter, r *http.Request) {
	// get uuid from url
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	uuid := parts[2]

	// preparing db query
	c := context.Background()
	queries := db_sqlc.New(db_sqlc.DB)

	// archiving todo
	err := queries.ArchiveTodo(c, uuid)
	if err != nil {
		renderAlert(w, "Whoops!...Something went wrong", "base", 3)
		log.Println(err)
		return
	}

	// render nothing to remove todo item from list
}
