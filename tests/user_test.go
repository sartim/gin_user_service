package main

import (
	"fmt"
	"gin-shop-api/app/controllers"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	Setup()
	// mockResponse := `{"data":[{"created_at":"0001-01-01T00:00:00Z","updated_at":"2023-01-26T15:11:42.478209+03:00","deleted_at":null,"id":"ab2fc765-cab7-4434-bef2-121bef275572","first_name":"Jane","last_name":"Doe","email":"janedoe@mail.com","is_active":true,"deleted":false}]}`
	r := SetUpRouter()
	r.GET("/api/v1/user", controllers.UserGetAll)
	req, _ := http.NewRequest("GET", "/api/v1/user", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)
	fmt.Println(string(responseData))
	// assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
	TearDown()
}
