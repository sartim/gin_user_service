package main

import (
	"fmt"
	"gin-shop-api/internal/controllers"
	"gin-shop-api/internal/repository"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	Setup()
	r := SetUpRouter()
	ctrl := controllers.NewUserController(repository.DB)
	r.GET("/api/v1/users", ctrl.GetAll)
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	req.Header.Set("Authorization", "1234")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)
	fmt.Println(string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
	TearDown()
}
