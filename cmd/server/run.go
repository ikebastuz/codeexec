package main

import (
	"codeexec/internal/runner"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const maxMemory = 10 * 1024 * 1024 // 10MB

var dockerSemaphore = make(chan struct{}, 5)

func runHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	lang := runner.Lang(r.FormValue("language"))
	code := r.FormValue("code")

	log.Infof("\nLanguage: %s\nCode:\n%s\n", lang, code)

	dockerSemaphore <- struct{}{}
	c := make(chan runner.Job)
	go func() {
		defer func() { <-dockerSemaphore }() // Release the slot when done
		runner.Run(lang, code, c)
	}()
	res := <-c

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{
		"stdout": res.Stdout,
		"stderr": res.Stderr,
	}
	log.Infof("\nstdout:\n%s\nstderr:\n%s\n", res.Stdout, res.Stderr)
	if res.Error != nil {
		log.Errorf("Response error: %s", res.Error)
		resp["error"] = res.Error.Error()
	}
	json.NewEncoder(w).Encode(resp)
}
