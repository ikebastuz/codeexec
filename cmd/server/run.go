package main

import (
	"codeexec/internal/config"
	"codeexec/internal/metrics"
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

	stdout, stderr, buildDuration, execDuration, err := runner.Run(lang, code)

	w.Header().Set("Content-Type", "application/json")
	resp := Response{
		Stdout:        stdout,
		Stderr:        stderr,
		ExecDuration:  execDuration,
		BuildDuration: buildDuration,
	}
	metrics.ExecutionsCounter.WithLabelValues(string(lang)).Inc()
	metrics.ExecutionsDuration.WithLabelValues(string(lang)).Observe(execDuration)

	if stderr != "" {
		metrics.StdErrCounter.WithLabelValues(string(lang)).Inc()
	}

	if err != nil {
		metrics.ErrorCounter.WithLabelValues(string(lang)).Inc()
		log.Errorf("Response error: %s", err)
		resp.Error = err.Error()
	}
	json.NewEncoder(w).Encode(resp)
}
