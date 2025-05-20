package main

import (
	"codeexec/internal/config"
	"codeexec/internal/runner"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var semaphore = make(chan struct{}, config.MAX_PROCESSES)

func runHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(config.MAX_MEMORY); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	lang := runner.Lang(r.FormValue("language"))
	code := r.FormValue("code")

	log.Infof("\nLanguage: %s\nCode:\n%s\n", lang, code)

	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	stdout, stderr, err := runner.Run(lang, code)

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{
		"stdout": stdout,
		"stderr": stderr,
	}
	if err != nil {
		log.Errorf("Response error: %s", err)
		resp["error"] = err.Error()
	}
	json.NewEncoder(w).Encode(resp)
}
