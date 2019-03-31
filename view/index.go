package view

import (
	"check-in-backend/controller"
	"github.com/labstack/echo"
)

func InitViewV1(group *echo.Group) {
	group.POST("/login", controller.Login)
	group.POST("/login-web", controller.LoginWeb)
}
