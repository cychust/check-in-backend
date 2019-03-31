package main

import (
	"check-in-backend/constant"
	"check-in-backend/middleware"
	"check-in-backend/view"
	"flag"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var debug bool

func init() {
	flag.BoolVar(&debug, "debug", false, "debug api interface")
	flag.Parse()
}

func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	v1 := e.Group(constant.APIPrefix)

	if !debug {
		skipperPath := []string{
			"/api/v1/login",
		}
		v1.Use(middleware.JWTWithConfig(mid.CustomJWTConfig(skipperPath, "Bearer")))
	}
	view.InitView(v1)
	e.Logger.Fatal(e.Start(":3000"))
}
