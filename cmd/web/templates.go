package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/rmcs87/snippetbox/pkg/forms"
	"github.com/rmcs87/snippetbox/pkg/models"
)

type templateData struct {
	CurrentYear       int
	Snippet           *models.Snippet
	Snippets          []*models.Snippet
	Form              *forms.Form
	Flash             string
	AuthenticatedUser int
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl.html"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
