package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}

	files := []string{
		".\\ui\\html\\home.page.tmpl.html",
		".\\ui\\html\\base.layout.tmpl.html",
		".\\ui\\html\\footer.partial.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(writer, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(writer, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(writer, "Internal Server Error", 500)
	}

	return
}

func createSnippet(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		writer.Header().Set("Allow", "POST")
		http.Error(writer, "Método não Permitido", 405)

		return
	}

	writer.Write([]byte("Create a new snippet ..."))
}

func showSnippet(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(writer, request)
		return
	}
	fmt.Fprintf(writer, "Exibindo Snippet de ID: %d", id)
}
