package views

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Template struct {
	htmlTpl *template.Template
}

func Parse(filePath string) (Template, error) {
	tpl, err := template.ParseFiles(filePath)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil
}

func ParseFS(fs embed.FS, patterns ...string) (Template, error) {
	tpl, err := template.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing embedded template: %w", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	err := t.htmlTpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}

func Must(template Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return template
}

// func (t Template) Parse(w http.ResponseWriter, data interface{}) {
// 	w.Header().Set("Content-Type", "text/html")
// 	tpl, err := template.ParseFiles(t.htmlTpl)
// 	if err != nil {
// 		log.Printf("parsing template: %v", err)
// 		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
// 		return
// 	}
// }
