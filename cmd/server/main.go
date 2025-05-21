package main

import (
	"net/http"

	"codeexec/internal/config"
	"codeexec/internal/metrics"
	"codeexec/internal/runner"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	config, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}

	metrics.InitMetrics()
	runner.PullAllImages()
	go runner.StartImageMonitor()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/run", runHandler)
	http.HandleFunc("/", indexHandler)

	log.Infof("Server on %s", config.PORT)
	log.Fatal(http.ListenAndServe(config.PORT, nil))
}
