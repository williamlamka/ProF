package route

import (
	"new_project/handler"

	"github.com/labstack/echo/v4"
)

func AuthRouteInit(e *echo.Echo) {
	g := e.Group("/auth")
	g.POST("/register", handler.Register)
	g.POST("/login", handler.Login)
}
