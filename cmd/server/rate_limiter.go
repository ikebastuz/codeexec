package main

import (
	"codeexec/internal/config"
	"codeexec/internal/metrics"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(config.RATE_LIMIT_REQ_S, config.RATE_LIMIT_BURST_S)
		visitors[ip] = limiter
	}

	return limiter
}

func rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := getVisitor(ip)

		if !limiter.Allow() {
			log.Warnf("Too Many Requests from ip: %s", ip)
			metrics.RateLimitCounter.WithLabelValues(ip).Inc()
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
