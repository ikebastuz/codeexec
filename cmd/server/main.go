package main

import (
	"net/http"

	"codeexec/internal/metrics"
	"codeexec/internal/runner"

	"github.com/prometheus/client_golang/prometheus/promhttp"

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

	metrics.InitMetrics()
	runner.PullAllImages()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/run", runHandler)
	http.HandleFunc("/", indexHandler)

	log.Infof("Server on %s", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
