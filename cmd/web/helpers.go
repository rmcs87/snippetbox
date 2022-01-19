package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(writer http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(writer http.ResponseWriter, status int) {
	http.Error(writer, http.StatusText(status), status)
}

func (app *application) notFound(writer http.ResponseWriter) {
	app.clientError(writer, http.StatusNotFound)
}

func (app *application) render(writer http.ResponseWriter, request *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(writer, fmt.Errorf("Teh tamplate %s does note exist", name))
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, td)
	if err != nil {
		app.serverError(writer, err)
		return
	}

	buf.WriteTo(writer)
}
