package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/rmcs87/snippetbox/pkg/models"
)

func (app *application) home(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		app.notFound(writer)
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(writer, err)
	}

	data := &templateData{Snippets: s}

	files := []string{
		".\\ui\\html\\home.page.tmpl.html",
		".\\ui\\html\\base.layout.tmpl.html",
		".\\ui\\html\\footer.partial.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(writer, err)
		return
	}

	err = ts.Execute(writer, data)
	if err != nil {
		app.serverError(writer, err)
	}
}

func (app *application) createSnippet(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		writer.Header().Set("Allow", "POST")
		app.clientError(writer, http.StatusMethodNotAllowed)

		return
	}

	title := "0 snail"
	content := "A arte de amar."
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(writer, err)
		return
	}

	http.Redirect(writer, request, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

func (app *application) showSnippet(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(writer)
		return
	}

	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(writer)
		return
	} else if err != nil {
		app.serverError(writer, err)
		return
	}

	data := &templateData{Snippet: s}

	files := []string{
		".\\ui\\html\\show.page.tmpl.html",
		".\\ui\\html\\base.layout.tmpl.html",
		".\\ui\\html\\footer.partial.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(writer, err)
		return
	}

	err = ts.Execute(writer, data)

	if err != nil {
		app.serverError(writer, err)
	}
}
