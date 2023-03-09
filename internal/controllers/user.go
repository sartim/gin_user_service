package controllers

import (
	"gin-shop-api/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	*BaseController
}

func NewUserController(db *gorm.DB) *UserController {
	var user models.User
	return &UserController{NewBaseController(db, user)}
}

func (ctrl *UserController) RegisterUserRoutes(router *gin.RouterGroup) {
	userRouter := router.Group("/users")
	{
		userRouter.GET("", ctrl.GetAll)
		userRouter.GET(":id", ctrl.Get)
		userRouter.POST("", ctrl.Create)
		userRouter.PUT(":id", ctrl.Update)
		userRouter.DELETE(":id", ctrl.Delete)
	}
}
