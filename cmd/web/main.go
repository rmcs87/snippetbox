package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/rmcs87/snippetbox/pkg/models"
	"github.com/rmcs87/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql" // New import
)

type Config struct {
	Addr      string
	StaticDir string
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	cfg      *Config
	snippets interface {
		Insert(string, string, string) (int, error)
		Get(int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	}
	templateCache map[string]*template.Template
	session       *sessions.Session
	users         interface {
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
}

func main() {
	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":8080", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", ".\\ui\\static\\", "Path to static Files")

	dsn := flag.String("dsn", "root:123456@/snippetbox?parseTime=true", "Mysql Connection String")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Cryp Secret")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)

	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache(".\\ui\\html")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 16 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		cfg:           cfg,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		session:       session,
		users:         &mysql.UserModel{DB: db},
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         cfg.Addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", cfg.Addr)

	err = srv.ListenAndServeTLS(".\\tls\\cert.pem", ".\\tls\\key.pem")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
