package main

import (
	"codeexec/internal/runner"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const maxMemory = 10 * 1024 * 1024 // 10MB

func runHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	lang := runner.Lang(r.FormValue("language"))
	code := r.FormValue("code")

	log.Infof("\nLanguage: %s\nCode:\n%s\n", lang, code)

	stdout, stderr, err := runner.Run(lang, code)
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{
		"stdout": stdout,
		"stderr": stderr,
	}
	log.Infof("\nstdout:\n%s\nstderr:\n%s\n", stdout, stderr)
	if err != nil {
		log.Errorf("Response error: %s", err)
		resp["error"] = err.Error()
	}
	json.NewEncoder(w).Encode(resp)
}
