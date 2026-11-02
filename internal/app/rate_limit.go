package app

import (
	"log"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"homework-reminder-bot/internal/ratelimit"
)

type rateLimitPolicy struct {
	key    string
	limit  int64
	window time.Duration
}

func requestIP(r *http.Request, trustProxyHeaders bool) string {
	if trustProxyHeaders {
		if forwardedIP := net.ParseIP(strings.TrimSpace(r.Header.Get("X-Real-IP"))); forwardedIP != nil {
			return forwardedIP.String()
		}
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		if ip := net.ParseIP(host); ip != nil {
			return ip.String()
		}
	}
	if ip := net.ParseIP(r.RemoteAddr); ip != nil {
		return ip.String()
	}
	return "unknown"
}

func enforceRateLimits(limiter ratelimit.Limiter, w http.ResponseWriter, r *http.Request, policies ...rateLimitPolicy) bool {
	for _, policy := range policies {
		allowed, retryAfter, err := limiter.Allow(r.Context(), policy.key, policy.limit, policy.window)
		if err != nil {
			log.Printf("Redis rate limiter unavailable: %v", err)
			http.Error(w, "rate limiter unavailable", http.StatusServiceUnavailable)
			return false
		}
		if allowed {
			continue
		}
		seconds := int(math.Ceil(retryAfter.Seconds()))
		if seconds < 1 {
			seconds = 1
		}
		w.Header().Set("Retry-After", strconv.Itoa(seconds))
		http.Error(w, "too many requests", http.StatusTooManyRequests)
		return false
	}
	return true
}
