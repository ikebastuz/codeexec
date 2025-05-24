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

	result := runner.NewRunner(lang, code).Run()

	w.Header().Set("Content-Type", "application/json")

	metrics.ExecutionsCounter.WithLabelValues(string(lang)).Inc()
	if result.BuildDuration > 0 {
		metrics.Duration.WithLabelValues(string(lang), "build").Observe(result.BuildDuration)
	}
	metrics.Duration.WithLabelValues(string(lang), "exec").Observe(result.ExecDuration)

	if result.Stderr != "" {
		metrics.StdErrCounter.WithLabelValues(string(lang)).Inc()
	}

	if result.Error != "" {
		metrics.ErrorCounter.WithLabelValues(string(lang)).Inc()
		metrics.ErrorTypeCounter.WithLabelValues(result.Error).Inc()

		log.Errorf("Response error: %s", result.Error)
	}

	json.NewEncoder(w).Encode(result)
}
