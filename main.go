package main

import (
	"csu-import/controllers"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	// post 账号密码下载ics
	e.POST("/api/download", controllers.Login)
	e.Logger.Fatal(e.Start(":3000"))
}
