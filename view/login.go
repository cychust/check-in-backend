package view

import (
	"check-in-backend/controller"
	"github.com/labstack/echo"
)

func InitView(group *echo.Group) {
	group.POST("/login", controller.Login)

}
