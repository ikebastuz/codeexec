package main

import (
	"net/http"

	"codeexec/internal/runner"

	log "github.com/sirupsen/logrus"
)

// TODO: move to env
const PORT = ":1450"

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	runner.PullAllImages()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/run", runHandler)
	log.Infof("Server on %s", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
