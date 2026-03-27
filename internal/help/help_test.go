package help

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupHelpRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/help", Help)
	return router
}

func TestHelp_PlainText_ReturnsOK(t *testing.T) {
	router := setupHelpRouter()
	request := httptest.NewRequest(http.MethodGet, "/help", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", response.Code)
	}
}

func TestHelp_PlainText_ContainsName(t *testing.T) {
	router := setupHelpRouter()
	request := httptest.NewRequest(http.MethodGet, "/help", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	body := response.Body.String()
	if !strings.Contains(body, "Deny By Default as a Service") {
		t.Errorf("expected plain text to contain service name, got: %s", body)
	}
}

func TestHelp_PlainText_ContainsEndpoints(t *testing.T) {
	router := setupHelpRouter()
	request := httptest.NewRequest(http.MethodGet, "/help", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	body := response.Body.String()
	for _, path := range []string{"/no", "/yes", "/random", "/health", "/help"} {
		if !strings.Contains(body, path) {
			t.Errorf("expected plain text to contain endpoint %s", path)
		}
	}
}

func TestHelp_PlainText_IsNotJSON(t *testing.T) {
	router := setupHelpRouter()
	request := httptest.NewRequest(http.MethodGet, "/help", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	var js map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &js); err == nil {
		t.Error("expected plain text response, but got valid JSON")
	}
}

func TestHelp_JSON_ViaQueryParam(t *testing.T) {
	router := setupHelpRouter()
	request := httptest.NewRequest(http.MethodGet, "/help?format=json", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", response.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected valid JSON, got error: %v", err)
	}

	if body["name"] == nil {
		t.Error("expected 'name' field in JSON response")
	}
	if body["endpoints"] == nil {
		t.Error("expected 'endpoints' field in JSON response")
	}
}

func TestHelp_JSON_ViaAcceptHeader(t *testing.T) {
	router := setupHelpRouter()
	request := httptest.NewRequest(http.MethodGet, "/help", nil)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", response.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected valid JSON, got error: %v", err)
	}
}

func TestHelp_JSON_ViaContentTypeHeader(t *testing.T) {
	router := setupHelpRouter()
	request := httptest.NewRequest(http.MethodGet, "/help", nil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	var body map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected valid JSON, got error: %v", err)
	}
}

func TestHelp_JSON_EndpointHasRequiredFields(t *testing.T) {
	router := setupHelpRouter()
	request := httptest.NewRequest(http.MethodGet, "/help?format=json", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	var body map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	endpoints, ok := body["endpoints"].([]interface{})
	if !ok {
		t.Fatal("expected 'endpoints' array in response")
	}

	for i, ep := range endpoints {
		epMap, ok := ep.(map[string]interface{})
		if !ok {
			t.Fatalf("expected endpoint %d to be an object", i)
		}

		if epMap["path"] == nil || epMap["path"] == "" {
			t.Errorf("endpoint %d: expected non-empty path", i)
		}

		if epMap["method"] == nil || epMap["method"] == "" {
			t.Errorf("endpoint %d: expected non-empty method", i)
		}

		if epMap["description"] == nil || epMap["description"] == "" {
			t.Errorf("endpoint %d: expected non-empty description", i)
		}
	}
}

func TestHelp_JSON_ContainsContentNegotiation(t *testing.T) {
	router := setupHelpRouter()
	request := httptest.NewRequest(http.MethodGet, "/help?format=json", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	var body map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	cn, ok := body["content_negotiation"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'content_negotiation' object in response")
	}

	if cn["default"] == nil || cn["default"] == "" {
		t.Error("expected 'default' field in content_negotiation")
	}

	if cn["json_methods"] == nil {
		t.Error("expected 'json_methods' field in content_negotiation")
	}
}
