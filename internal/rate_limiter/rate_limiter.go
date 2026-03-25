package ratelimiter

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rajathjn/deny-by-default-as-a-service/internal/utils"
	"golang.org/x/time/rate"
)

var (
	rateLimit  float64 = 0.5 // 1 request every 2 seconds
	burstLimit int     = 3 // Allow short bursts of up to 3 requests
)

func init() {
	if v := os.Getenv("RATE_LIMIT"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil && parsed > 0 {
			rateLimit = parsed
		}
	}
	if v := os.Getenv("BURST_LIMIT"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			burstLimit = parsed
		}
	}

	go cleanupLimiters()
}

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu         sync.Mutex
	limiterMap = make(map[string]*clientLimiter)
)

func cleanupLimiters() {
	for {
		time.Sleep(time.Minute)
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Cleanup recovered from panic: %v", r)
				}
			}()

			mu.Lock()
			defer mu.Unlock()

			for ip, c := range limiterMap {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(limiterMap, ip)
				}
			}
		}()
	}
}

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	c, exists := limiterMap[ip]
	if !exists {
		c = &clientLimiter{
			limiter:  rate.NewLimiter(rate.Limit(rateLimit), burstLimit),
			lastSeen: time.Now(),
		}
		limiterMap[ip] = c
	}

	return c.limiter
}

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		limiter := getLimiter(clientIP)

		if !limiter.Allow() {
			msg := "Too many requests. Do you really need that many reasons!! Please slow down."
			if utils.WantsJSON(c) {
				c.JSON(http.StatusTooManyRequests, gin.H{"error": msg})
			} else {
				c.String(http.StatusTooManyRequests, msg)
			}
			c.Abort()
		}
	}
}
