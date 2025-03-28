package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/silasburger/lenslocked/context"
	"github.com/silasburger/lenslocked/models"
)

type Template struct {
	htmlTpl *template.Template
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(patterns[0])
	tpl = tpl.Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return `<-- CSRF field placeholder until template is executed -->`, fmt.Errorf("csrfField not implemented")
		},
		"currentUser": func() (*models.User, error) {
			return nil, fmt.Errorf("currentUser not implemented")
		},
	},
	)
	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing embedded template: %w", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}
	tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
		"currentUser": func() *models.User {
			return context.User(r.Context())
		},
	})
	w.Header().Set("Content-Type", "text/html")

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func Must(template Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return template
}

// Earlier ways of parsing templates. Now it is done on startup and embedded with views.ParseFS

// func Parse(filePath string) (Template, error) {
// 	tpl, err := template.ParseFiles(filePath)
// 	if err != nil {
// 		return Template{}, fmt.Errorf("parsing template: %w", err)
// 	}
// 	return Template{
// 		htmlTpl: tpl,
// 	}, nil
// }

// func (t Template) Parse(w http.ResponseWriter, data interface{}) {
// 	w.Header().Set("Content-Type", "text/html")
// 	tpl, err := template.ParseFiles(t.htmlTpl)
// 	if err != nil {
// 		log.Printf("parsing template: %v", err)
// 		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
// 		return
// 	}
// }
