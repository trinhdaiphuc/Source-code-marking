package middlewares

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
)

func GetUser(c echo.Context) models.User {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(interface{})
	user := claims.(models.User)
	return user
}

func IsStudent(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userToken := c.Get("user").(*jwt.Token)
		claims := userToken.Claims.(jwt.MapClaims)
		userRole := claims["role"].(string)
		if userRole == "STUDENT" {
			return next(c)
		}
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "Invalid role",
		}
	}
}

func IsTeacher(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userToken := c.Get("user").(*jwt.Token)
		claims := userToken.Claims.(jwt.MapClaims)
		userRole := claims["role"].(string)
		if userRole == "TEACHER" {
			return next(c)
		}
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "Invalid role",
		}
	}
}

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userToken := c.Get("user").(*jwt.Token)
		claims := userToken.Claims.(jwt.MapClaims)
		userRole := claims["role"].(string)
		if userRole == "ADMIN" {
			return next(c)
		}
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "Invalid role",
		}
	}
}

func IsAdminOrTeacher(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userToken := c.Get("user").(*jwt.Token)
		claims := userToken.Claims.(jwt.MapClaims)
		userRole := claims["role"].(string)
		if userRole == "ADMIN" || userRole == "TEACHER" {
			return next(c)
		}
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "Invalid role",
		}
	}
}

func IsAdminOrStudent(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userToken := c.Get("user").(*jwt.Token)
		claims := userToken.Claims.(jwt.MapClaims)
		userRole := claims["role"].(string)
		if userRole == "ADMIN" || userRole == "STUDENT" {
			return next(c)
		}
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "Invalid role",
		}
	}
}

func IsTeacherOrStudent(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userToken := c.Get("user").(*jwt.Token)
		claims := userToken.Claims.(jwt.MapClaims)
		userRole := claims["role"].(string)
		if userRole == "TEACHER" || userRole == "STUDENT" {
			return next(c)
		}
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "Invalid role",
		}
	}
}
