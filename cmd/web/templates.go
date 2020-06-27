package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

type templateData struct {
	Title string
	Year  int
}

func addDefaultData(t *templateData) *templateData {
	if t == nil {
		t = &templateData{}
	}

	if t.Title == "" {
		t.Title = "Deafault"
	}

	t.Year = time.Now().Year()
	return t
}

func (s *Server) render(w http.ResponseWriter, templateName string, td *templateData) {

	if td == nil || td.Title == "" {
		s.serverError(w, fmt.Errorf("td.Title is empty"))
		return
	}

	key := fmt.Sprintf("%s.page.html", templateName)
	val, ok := s.templateCache[key]
	if !ok {
		s.serverError(w, fmt.Errorf("Template  %s not found", key))
		return
	}

	buf := new(bytes.Buffer)
	err := val.ExecuteTemplate(buf, key, addDefaultData(td))

	if err != nil {
		s.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	res := map[string]*template.Template{}

	files, err := filepath.Glob(filepath.Join(dir, "*.page.html"))

	if err != nil {
		return nil, err
	}

	for _, f := range files {
		tmpl, err := template.ParseFiles(f, filepath.Join(dir, "footer.partial.html"), filepath.Join(dir, "base.layout.html"))

		if err != nil {
			return nil, err
		}

		res[filepath.Base(f)] = tmpl
	}
	return res, nil
}
