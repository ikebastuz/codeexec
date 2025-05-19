package main

import (
	"codeexec/internal/runner"
	"encoding/json"
	"log"
	"net/http"
)

func runHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	lang := runner.Lang(r.FormValue("language"))
	code := r.FormValue("code")

	log.Printf("Language: %s\nCode:\n%s\n", lang, code)

	stdout, stderr, err := runner.Run(lang, code)
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{
		"stdout": stdout,
		"stderr": stderr,
	}
	if err != nil {
		resp["error"] = err.Error()
	}
	json.NewEncoder(w).Encode(resp)
}
