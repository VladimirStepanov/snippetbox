package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"githib.com/VladimirStepanov/snippetbox/pkg/models"
	validation "github.com/go-ozzo/ozzo-validation"
)

type snippetForm struct {
	Title   string
	Content string
	Expire  string
	Type    string
}

type templateData struct {
	Snippets    []*models.Snippet
	Snippet     *models.Snippet
	User        *models.User
	FormUser    *models.User
	FormSnippet *snippetForm
	Errors      validation.Errors
	Flashes     []interface{}
	CSRFField   template.HTML
	Title       string
	Year        int
}

func getError(errMap validation.Errors, key string) string {
	if value, ok := errMap[key]; ok {
		return value.Error()
	}

	return ""
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006")
}

func (s *Server) addDefaultData(t *templateData) *templateData {
	if t == nil {
		t = &templateData{}
	}

	t.Year = time.Now().Year()
	return t
}

func (s *Server) render(w http.ResponseWriter, r *http.Request, templateName string, td *templateData) {
	var err error
	key := fmt.Sprintf("%s.page.html", templateName)
	val, ok := s.templateCache[key]
	if !ok {
		s.serverError(w, fmt.Errorf("Template  %s not found", key))
		return
	}

	td = s.addDefaultData(td)

	td.Flashes, err = s.getFlashes(w, r)
	td.User = getAuthUserFromRequest(r)

	if err != nil {
		s.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	err = val.ExecuteTemplate(buf, key, td)

	if err != nil {
		s.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {

	funcMap := template.FuncMap{
		"humanDate": humanDate,
		"getError":  getError,
	}

	res := map[string]*template.Template{}

	files, err := filepath.Glob(filepath.Join(dir, "*.page.html"))

	if err != nil {
		return nil, err
	}

	for _, f := range files {
		name := filepath.Base(f)
		tmpl, err := template.New(name).Funcs(funcMap).ParseFiles(f, filepath.Join(dir, "footer.partial.html"), filepath.Join(dir, "base.layout.html"))

		if err != nil {
			return nil, err
		}

		res[name] = tmpl
	}
	return res, nil
}
