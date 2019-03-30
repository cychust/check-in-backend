package main

import (
	"check-in-backend/constant"
	"check-in-backend/view"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	v1 := e.Group(constant.APIPrefix)

	view.InitView(v1)
	e.Logger.Fatal(e.Start(":3000"))
}
