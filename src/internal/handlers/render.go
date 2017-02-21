package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

var tmpl *template.Template

func init() {
	// load up all templates
	var err error
	tmpl, err = template.New("").ParseGlob("./templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
}

func render(w http.ResponseWriter, tmpl *template.Template, tmplName string, data interface{}) {
	buf := &bytes.Buffer{}
	err := tmpl.ExecuteTemplate(buf, tmplName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}
