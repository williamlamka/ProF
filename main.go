package main

import (
	"new_project/config"
	"new_project/handler"
	"new_project/redis"
	"new_project/route"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	//middlerware
	e.Use(middleware.Logger())
	err := godotenv.Load()
	if err != nil {
		e.Logger.Fatal("Error loading .env file")
	}
	e.HTTPErrorHandler = handler.CustomHTTPErrorHandler

	//db
	config.DBInit()
	redis.RedisInit()
	//route
	route.AuthRouteInit(e)
	route.UserRouteInit(e)
	route.TransactionRouteInit(e)
	route.BillRouteInit(e)

	port := ":" + os.Getenv("PORT")
	e.Logger.Fatal(e.Start(port))
}
