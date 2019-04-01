package view

import (
	"check-in-backend/controller"
	"github.com/labstack/echo"
)

func InitGroupView(group *echo.Group) {
	group.GET("/list", controller.GetGroups)
	group.POST("", controller.CreateGroup)
	//group.GET("", controller.GetGroup)
	group.POST("/delete", controller.DeleteGroups)
}
