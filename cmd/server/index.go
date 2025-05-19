package main

import (
	"codeexec/internal/runner"
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

type PageData struct {
	Languages []string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Languages: runner.GetLangs(),
	}
	tmpl.Execute(w, data)
}
