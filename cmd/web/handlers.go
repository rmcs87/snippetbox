package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rmcs87/snippetbox/pkg/models"
)

func (app *application) home(writer http.ResponseWriter, request *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(writer, err)
	}

	app.render(writer, request, "home.page.tmpl.html", &templateData{
		Snippets: s,
	})
}

func (app *application) createSnippetForm(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Create a new snippet..."))
}

func (app *application) createSnippet(writer http.ResponseWriter, request *http.Request) {
	title := "0 snail"
	content := "A arte de amar."
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(writer, err)
		return
	}

	http.Redirect(writer, request, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) showSnippet(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get(":id"))
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

	app.render(writer, request, "show.page.tmpl.html", data)
}
