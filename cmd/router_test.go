package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gin-shop-api/internal/config"
	"gin-shop-api/internal/helpers/crypto"
	"gin-shop-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const testSecret = "test-secret-with-at-least-32-characters"

func testRouter(t *testing.T) (*gorm.DB, http.Handler) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:router-test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.Exec(`CREATE TABLE user (
		id text PRIMARY KEY, created_at datetime, updated_at datetime, deleted_at datetime,
		first_name text, last_name text, email text UNIQUE, password text,
		is_active numeric, is_admin numeric, deleted numeric
	)`).Error; err != nil {
		t.Fatalf("create test user table: %v", err)
	}
	cfg := config.App{Environment: "test", Port: "8000", SecretKey: testSecret, AccessTokenTTL: time.Hour, AllowedOrigins: []string{"http://localhost:3000"}}
	return db, newRouter(db, cfg)
}

func TestRouterHealthAndAuthentication(t *testing.T) {
	db, router := testRouter(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	response := routerRequest(router, http.MethodGet, "/health/live", nil, "")
	if response.Code != http.StatusOK {
		t.Fatalf("expected live health 200, got %d", response.Code)
	}

	response = routerRequest(router, http.MethodGet, "/api/v1/users", nil, "")
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected protected route 401, got %d", response.Code)
	}

	response = routerRequest(router, http.MethodPost, "/api/v1/auth/token", []byte(`{"email":"missing@example.com","password":"password"}`), "")
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected invalid login 401, got %d", response.Code)
	}
}

func TestRouterMetrics(t *testing.T) {
	db, router := testRouter(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	response := routerRequest(router, http.MethodGet, "/health/live", nil, "")
	if response.Code != http.StatusOK {
		t.Fatalf("expected live health 200, got %d", response.Code)
	}

	response = routerRequest(router, http.MethodGet, "/metrics", nil, "")
	if response.Code != http.StatusOK {
		t.Fatalf("expected metrics 200, got %d", response.Code)
	}
	if !bytes.Contains(response.Body.Bytes(), []byte("http_requests_total")) {
		t.Fatal("expected HTTP request counter in metrics output")
	}
}

func TestRouterAdminCanAccessUsers(t *testing.T) {
	db, router := testRouter(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	hash, err := crypto.HashPassword("correct-password")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	user := models.User{ID: uuid.New(), FirstName: "Admin", LastName: "User", Email: "admin@example.com", Password: hash, IsActive: true, IsAdmin: true}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create admin: %v", err)
	}

	response := routerRequest(router, http.MethodPost, "/api/v1/auth/token", []byte(`{"email":"admin@example.com","password":"correct-password"}`), "")
	if response.Code != http.StatusOK {
		t.Fatalf("expected login 200, got %d: %s", response.Code, response.Body.String())
	}
	var token struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &token); err != nil {
		t.Fatalf("decode token response: %v", err)
	}
	if token.AccessToken == "" {
		t.Fatal("expected access token")
	}

	response = routerRequest(router, http.MethodGet, "/api/v1/users?limit=10", nil, token.AccessToken)
	if response.Code != http.StatusOK {
		t.Fatalf("expected users 200, got %d: %s", response.Code, response.Body.String())
	}
}

func routerRequest(router http.Handler, method, path string, body []byte, token string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, bytes.NewReader(body))
	if len(body) > 0 {
		request.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}
