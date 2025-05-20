package main

import (
	"codeexec/internal/metrics"
	"codeexec/internal/runner"
	"encoding/json"
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))
var pageData PageData

type PageData struct {
	Languages     []string
	SampleCodeMap string // JSON string [lang:sampleCode]
}

func init() {
	langs := runner.GetLangs()
	sampleCodeMap := make(map[string]string)
	for lang, def := range runner.LangDefinitions {
		sampleCodeMap[string(lang)] = def.SampleCode
	}
	sampleCodeJSON, err := json.Marshal(sampleCodeMap)
	if err != nil {
		log.Fatalf("Failed to marshal sample code map: %v", err)
	}
	pageData = PageData{
		Languages:     langs,
		SampleCodeMap: string(sampleCodeJSON),
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		metrics.IndexPageCounter.Inc()
	}
	tmpl.Execute(w, pageData)
}
