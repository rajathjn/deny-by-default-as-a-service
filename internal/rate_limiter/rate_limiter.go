package ratelimiter

import (
	"net/http"
	"time"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type client_limiter struct {
	limiter 	*rate.Limiter
	last_seen	time.Time
}

var  (
	mutex_lock sync.Mutex
	ratelimiter = make(map[string]*client_limiter)
)

func init() {
	go cleanup_ratelimiters()
}

func cleanup_ratelimiters() {
	for {
		time.Sleep(time.Minute)
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Cleanup recovered from panic: %v", r)
				}
			}()

			mutex_lock.Lock()
			defer mutex_lock.Unlock()

			for ip, c := range ratelimiter {
				if time.Since(c.last_seen) > 3*time.Minute {
					delete(ratelimiter, ip)
				}
			}
		}()
	}
}

func get_limiter(ip string) *rate.Limiter {

	mutex_lock.Lock()
	defer mutex_lock.Unlock()

	c, exists := ratelimiter[ip]
	if !exists {
		c = &client_limiter{
			limiter: rate.NewLimiter(rate.Limit(2),5),
			last_seen: time.Now(),
		}
		ratelimiter[ip] = c
	}

	return c.limiter
}

func Ratelimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		client_ip := c.ClientIP()
		limiter := get_limiter(client_ip)

		if !limiter.Allow() {
			c.String(
				http.StatusTooManyRequests,
				"Too many requests. Do you really need that many reasons!! Please slow down.",
			)
			c.Abort()
		}
	}
}