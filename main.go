package main

import (
	"csu-import/controllers"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Static("/", "public")
	// post 账号密码下载ics
	e.POST("/api/download", controllers.Login)
	e.Logger.Fatal(e.Start(":9093"))
}
