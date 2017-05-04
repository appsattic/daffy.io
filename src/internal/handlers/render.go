package handlers

import (
	"bytes"
	"html/template"
	"net/http"
)

func render(w http.ResponseWriter, tmpl *template.Template, tmplName string, data interface{}) {
	buf := &bytes.Buffer{}
	err := tmpl.ExecuteTemplate(buf, tmplName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}
