package main

import (
	"log"
	"net/http"
)

func home(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)

		return
	}

	writer.Write(([]byte("Hello from Snippetbox ...")))
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
	writer.Write([]byte("Display a specifi snippet ..."))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Printf("%s", "Starting Server on :8080")

	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
