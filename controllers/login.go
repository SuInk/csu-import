package controllers

import (
	"csu-import/models"
	"fmt"
	"github.com/labstack/echo/v4"
)

// Login e.GET("/users/:id/:pwd", GetUser)
func Login(c echo.Context) error {
	// User ID from path `users/:id`

	fmt.Println(c.FormValue("id") + "登录成功")

	user := &models.User{
		Id:  c.FormValue("id"),
		Pwd: c.FormValue("pwd"),
	}
	client, err := models.Login(user)
	if err != nil {
		fmt.Println(err)
		return err
	}
	body, _ := models.GetCourse(client)
	course, err := models.CourseParser(body)
	_, err = models.GetIcs(course)
	fmt.Println(c.FormValue("id") + "登陆成功")
	if err != nil {
		return err
	}
	return c.File("temp.ics")
}
