package models

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetUser e.GET("/users/:id", GetUser)
func GetUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}
