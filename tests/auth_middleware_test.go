package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gin-shop-api/internal/middleware"
	"gin-shop-api/internal/models"

	"github.com/gin-gonic/gin"
)

func TestRequireAuthRejectsInvalidHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name   string
		header string
	}{
		{name: "missing"},
		{name: "scheme only", header: "Bearer"},
		{name: "wrong scheme", header: "Basic abc"},
		{name: "too many fields", header: "Bearer abc extra"},
		{name: "malformed token", header: "Bearer not-a-jwt"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nextCalled := false
			router := gin.New()
			router.Use(middleware.RequireAuth(nil, "test-secret-with-at-least-32-characters"))
			router.GET("/protected", func(c *gin.Context) {
				nextCalled = true
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if test.header != "" {
				req.Header.Set("Authorization", test.header)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)

			if response.Code != http.StatusUnauthorized {
				t.Fatalf("expected 401, got %d", response.Code)
			}
			if nextCalled {
				t.Fatal("protected handler was called")
			}
		})
	}
}

func TestRequireAdmin(t *testing.T) {
	tests := []struct {
		name       string
		user       models.User
		wantStatus int
	}{
		{name: "regular user", user: models.User{IsAdmin: false}, wantStatus: http.StatusForbidden},
		{name: "administrator", user: models.User{IsAdmin: true}, wantStatus: http.StatusOK},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("user", test.user)
				c.Next()
			})
			router.Use(middleware.RequireAdmin())
			router.GET("/admin", func(c *gin.Context) { c.Status(http.StatusOK) })

			response := performRequest(router, http.MethodGet, "/admin", "")
			if response.Code != test.wantStatus {
				t.Fatalf("expected %d, got %d", test.wantStatus, response.Code)
			}
		})
	}
}
