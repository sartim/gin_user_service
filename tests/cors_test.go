package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gin-shop-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func TestCORSUsesConfiguredOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.CORS([]string{"https://app.example.com"}))
	router.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "https://app.example.com")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example.com" {
		t.Fatalf("unexpected allowed origin %q", got)
	}
}

func TestCORSRejectsUnknownPreflightOrigin(t *testing.T) {
	router := gin.New()
	router.Use(middleware.CORS([]string{"https://app.example.com"}))
	router.OPTIONS("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "https://evil.example")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", response.Code)
	}
}
