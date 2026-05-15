package middleware

import (
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	limitDuration = 60 * time.Second
	maxRequests   = 100
)

type visitor struct {
	count      int
	limitStart time.Time
	lastSeen   time.Time
}

var visitors = make(map[string]visitor)
var mutex sync.Mutex

func init() {
	go CleanupVisitors()
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
		if allowedOrigin == "" {
			allowedOrigin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Idempotency-Key")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, errSplitHost := net.SplitHostPort(r.RemoteAddr)

		if errSplitHost != nil {
			http.Error(w, "Error getting IP address", http.StatusInternalServerError)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		v, ok := visitors[host]

		if !ok {
			visitors[host] = visitor{
				count:      1,
				limitStart: time.Now(),
				lastSeen:   time.Now(),
			}
			next.ServeHTTP(w, r)

		} else {
			timeElapsed := time.Since(v.limitStart)
			countIncreased := v.count + 1

			if timeElapsed >= limitDuration {
				visitors[host] = visitor{
					count:      1,
					limitStart: time.Now(),
					lastSeen:   time.Now(),
				}

				next.ServeHTTP(w, r)

			} else if countIncreased >= maxRequests {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			} else if countIncreased < maxRequests {
				visitors[host] = visitor{
					count:      countIncreased,
					limitStart: v.limitStart,
					lastSeen:   time.Now(),
				}
				next.ServeHTTP(w, r)
			}
		}
	})
}

func CleanupVisitors() {
	for {
		time.Sleep(5 * time.Minute)

		mutex.Lock()

		for ip, v := range visitors {
			timeElapsed := time.Now().Sub(v.lastSeen)

			if timeElapsed.Minutes() >= 5 {
				delete(visitors, ip)
			}
		}

		mutex.Unlock()
	}
}
