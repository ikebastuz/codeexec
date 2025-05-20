package main

import (
	"codeexec/internal/runner"
	"encoding/json"
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

type PageData struct {
	Languages     []string
	SampleCodeMap string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	langs := []string{}
	sampleCodeMap := make(map[string]string)
	for lang, def := range runner.LangDefinitions {
		langs = append(langs, string(lang))
		sampleCodeMap[string(lang)] = def.SampleCode
	}
	sampleCodeJSON, _ := json.Marshal(sampleCodeMap)
	data := PageData{
		Languages:     langs,
		SampleCodeMap: string(sampleCodeJSON),
	}
	tmpl.Execute(w, data)
}
