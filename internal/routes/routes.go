package routes

import (
	"gin-shop-api/internal/controllers"
	"gin-shop-api/internal/helpers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/api/v1/generate-jwt", controllers.GenerateJWT)
}

func UserRoutes(r *gin.Engine) {
	r.GET("/api/v1/user", helpers.RequireAuth, controllers.UserGetAll)
	r.GET("/api/v1/user/:id", helpers.RequireAuth, controllers.UserGetByID)
	r.POST("/api/v1/user", helpers.RequireAuth, controllers.UserCreate)
	r.PUT("/api/v1/user/:id", helpers.RequireAuth, controllers.UserUpdate)
	r.DELETE("/api/v1/user/:id", helpers.RequireAuth, controllers.UserDelete)
}
