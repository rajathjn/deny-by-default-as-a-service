package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rajathjn/deny-by-default-as-a-service/internal/help"
	"github.com/rajathjn/deny-by-default-as-a-service/internal/rate_limiter"
	"github.com/rajathjn/deny-by-default-as-a-service/internal/utils"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(cors.Default())
	router.Use(ratelimiter.RateLimiter())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/help", help.Help)
	router.GET("/", func(c *gin.Context) {
		respondWithReason(c, utils.GetNegativeReason(), "no")
	})
	router.GET("/no", func(c *gin.Context) {
		respondWithReason(c, utils.GetNegativeReason(), "no")
	})
	router.GET("/yes", func(c *gin.Context) {
		respondWithReason(c, utils.GetPositiveReason(), "yes")
	})
	router.NoRoute(func(c *gin.Context) {
		respondWithReason(c, utils.GetPositiveReason(), "yes")
	})
	return router
}

func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.RemoteAddr = "10.1.0.1:1234"
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", response.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("expected status 'ok', got '%s'", body["status"])
	}
}

func TestRootReturnsNegativeReason(t *testing.T) {
	router := setupTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.RemoteAddr = "10.1.0.2:1234"
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", response.Code)
	}
	if response.Body.Len() == 0 {
		t.Error("expected non-empty response body")
	}
}

func TestNoEndpoint(t *testing.T) {
	router := setupTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/no", nil)
	request.RemoteAddr = "10.1.0.3:1234"
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", response.Code)
	}
	if response.Body.Len() == 0 {
		t.Error("expected non-empty response body")
	}
}

func TestYesEndpoint(t *testing.T) {
	router := setupTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/yes", nil)
	request.RemoteAddr = "10.1.0.4:1234"
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", response.Code)
	}
	if response.Body.Len() == 0 {
		t.Error("expected non-empty response body")
	}
}

func TestJSONFormatQueryParam(t *testing.T) {
	router := setupTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/?format=json", nil)
	request.RemoteAddr = "10.1.0.5:1234"
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", response.Code)
	}

	var body jsonResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected JSON response, got error: %v", err)
	}
	if body.Reason == "" {
		t.Error("expected non-empty reason in JSON response")
	}
	if body.Type != "no" {
		t.Errorf("expected type 'no', got '%s'", body.Type)
	}
}

func TestJSONFormatAcceptHeader(t *testing.T) {
	router := setupTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/yes", nil)
	request.RemoteAddr = "10.1.0.6:1234"
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", response.Code)
	}

	var body jsonResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected JSON response, got error: %v", err)
	}
	if body.Type != "yes" {
		t.Errorf("expected type 'yes', got '%s'", body.Type)
	}
}

func TestNoRouteReturnsPositiveReason(t *testing.T) {
	router := setupTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/nonexistent-path", nil)
	request.RemoteAddr = "10.1.0.7:1234"
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", response.Code)
	}
	if response.Body.Len() == 0 {
		t.Error("expected non-empty response body")
	}
}

func TestHelpEndpoint(t *testing.T) {
	router := setupTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/help", nil)
	request.RemoteAddr = "10.1.0.8:1234"
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", response.Code)
	}
	if response.Body.Len() == 0 {
		t.Error("expected non-empty response body")
	}
}