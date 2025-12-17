package middleware

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type client struct {
	limiter *rate.Limiter
}

var (
	mu      sync.Mutex
	clients = make(map[string]*client)
)

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if c, exists := clients[ip]; exists {
		return c.limiter
	}

	limiter := rate.NewLimiter(1, 5) // 1 req/sec, burst 5
	clients[ip] = &client{limiter: limiter}
	return limiter
}

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
