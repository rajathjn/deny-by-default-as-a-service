package ratelimiter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimiter())
	router.GET(
		"/",
		func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		},
	)
	return router
}

func TestRateLimiterAllowsNormalTraffic(t *testing.T) {
	router := setupRouter()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.RemoteAddr = "10.0.0.1:1234"
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", response.Code)
	}
}

func TestRateLimiterBlocksExcessiveTraffic(t *testing.T) {
	router := setupRouter()

	var lastCode int
	for i := 0; i < 20; i++ {
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.RemoteAddr = "10.0.0.2:1234"
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		lastCode = response.Code
	}

	if lastCode != http.StatusTooManyRequests {
		t.Errorf("expected 429 after burst, got %d", lastCode)
	}
}

func TestRateLimiterPerIP(t *testing.T) {
	router := setupRouter()

	// Exhaust limiter for IP1
	for i := 0; i < 20; i++ {
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.RemoteAddr = "10.0.0.3:1234"
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
	}

	// IP2 should still be allowed
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.RemoteAddr = "10.0.0.4:1234"
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected 200 for different IP, got %d", response.Code)
	}
}
