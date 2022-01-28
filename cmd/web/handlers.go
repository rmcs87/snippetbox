package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rmcs87/snippetbox/pkg/forms"
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
	app.render(writer, request, "create.page.tmpl.html", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	form := forms.New(request.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(writer, request, "create.page.tmpl.html", &templateData{Form: form})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(writer, err)
		return
	}

	app.session.Put(request, "flash", "Snippet successfully created!")

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

func (app *application) signupUserForm(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "signup.page.tmpl.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	form := forms.New(request.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(writer, request, "signup.page.tmpl.html", &templateData{Form: form})
		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.render(writer, request, "signup.page.tmpl.html", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(writer, err)
		return
	}

	app.session.Put(request, "flash", "Your signup was successful. Please log in")

	http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "login.page.tmpl.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	form := forms.New(request.Form)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect ")
		app.render(writer, request, "login.page.tmpl.html", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(writer, err)
		return
	}
	app.session.Put(request, "userID", id)

	http.Redirect(writer, request, "/snippet/create", http.StatusSeeOther)

}

func (app *application) logoutUser(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Logout the User")
}
