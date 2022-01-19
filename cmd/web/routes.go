package main

import "net/http"

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir(app.cfg.StaticDir))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
