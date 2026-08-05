package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gin-shop-api/internal/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type testRecord struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

func TestBaseControllerPaginationAndCRUD(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.AutoMigrate(&testRecord{}); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}
	for _, record := range []testRecord{{Name: "one"}, {Name: "two"}, {Name: "three"}} {
		if err := db.Create(&record).Error; err != nil {
			t.Fatalf("seed test database: %v", err)
		}
	}

	controller := controllers.NewBaseController[testRecord](db, map[string]string{"name": "name"})
	router := gin.New()
	router.GET("/records", controller.GetAll)
	router.GET("/records/:id", controller.Get)
	router.PATCH("/records/:id", controller.Update)
	router.DELETE("/records/:id", controller.Delete)

	response := performRequest(router, http.MethodGet, "/records?page=2&limit=2", "")
	if response.Code != http.StatusOK {
		t.Fatalf("expected list status 200, got %d", response.Code)
	}
	var list struct {
		Data       []testRecord `json:"data"`
		Pagination struct {
			Total int64 `json:"total"`
		} `json:"pagination"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &list); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(list.Data) != 1 || list.Pagination.Total != 3 {
		t.Fatalf("unexpected pagination response: %+v", list)
	}

	response = performRequest(router, http.MethodPatch, "/records/1", `{"name":"updated"}`)
	if response.Code != http.StatusOK {
		t.Fatalf("expected update status 200, got %d", response.Code)
	}

	response = performRequest(router, http.MethodPatch, "/records/1", `{"id":99}`)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected protected-field status 400, got %d", response.Code)
	}

	response = performRequest(router, http.MethodDelete, "/records/1", "")
	if response.Code != http.StatusNoContent {
		t.Fatalf("expected delete status 204, got %d", response.Code)
	}

	response = performRequest(router, http.MethodGet, "/records/1", "")
	if response.Code != http.StatusNotFound {
		t.Fatalf("expected missing record status 404, got %d", response.Code)
	}
}

func performRequest(router http.Handler, method, target, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	return response
}
