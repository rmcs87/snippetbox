package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Addr      string
	StaticDir string
}

func main() {
	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":8080", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", ".\\ui\\static\\", "Path to static Files")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(http.Dir(cfg.StaticDir))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	infoLog.Printf("Starting server on %s", cfg.Addr)

	srv := &http.Server{
		Addr:     cfg.Addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
