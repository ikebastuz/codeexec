package main

import (
	"codeexec/internal/config"
	"codeexec/internal/metrics"
	"codeexec/internal/runner"
	"encoding/json"
	"net/http"

	"codeexec/internal/db"
	"database/sql"

	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

var semaphore = make(chan struct{}, config.MAX_PROCESSES)

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func nullFloatToFloat(nf sql.NullFloat64) float64 {
	if nf.Valid {
		return nf.Float64
	}
	return 0
}

func runHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(config.MAX_MEMORY); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	lang := runner.Lang(r.FormValue("language"))
	code := r.FormValue("code")

	encoded := runner.EncodeSource(code)

	cfg, err := config.GetConfig()
	if err != nil {
		http.Error(w, "server config error", http.StatusInternalServerError)
		return
	}
	dsn := "host=" + cfg.DBHost + " port=" + cfg.DBPort + " user=" + cfg.DBUser + " password=" + cfg.DBPass + " dbname=" + cfg.DBName + " sslmode=disable"
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		http.Error(w, "db connection error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()
	queries := db.New(dbConn)

	// Try to get cached result
	cached, err := queries.GetCodeExecutionResult(r.Context(), db.GetCodeExecutionResultParams{
		EncodedCode: encoded,
		Language:    string(lang),
	})
	if err == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"stdout":         nullStringToString(cached.Stdout),
			"stderr":         nullStringToString(cached.Stderr),
			"exec_duration":  nullFloatToFloat(cached.ExecDuration),
			"build_duration": nullFloatToFloat(cached.BuildDuration),
			"error":          nullStringToString(cached.Error),
		})
		return
	}

	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	result := runner.NewRunner(lang, code).Run()

	// Store result in DB
	_, _ = queries.InsertCodeExecutionResult(r.Context(), db.InsertCodeExecutionResultParams{
		Code:          code,
		Language:      string(lang),
		EncodedCode:   encoded,
		Stdout:        toNullString(result.Stdout),
		Stderr:        toNullString(result.Stderr),
		Error:         toNullString(result.Error),
		BuildDuration: sql.NullFloat64{Float64: result.BuildDuration, Valid: true},
		ExecDuration:  sql.NullFloat64{Float64: result.ExecDuration, Valid: true},
	})

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
