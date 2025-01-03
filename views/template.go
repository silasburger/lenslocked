package views

import (
	"html/template"
	"log"
	"net/http"
)

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Excute(w http.ResponseWriter, data interface{}) {
	err = tpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}

func (t Template) Parse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	tpl, err := template.ParseFiles(t.htmlTpl.)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}
}