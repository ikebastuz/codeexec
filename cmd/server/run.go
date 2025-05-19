package main

import (
	"codeexec/internal/runner"
	"fmt"
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

	out, err := runner.Run(lang, code)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(out))
}
