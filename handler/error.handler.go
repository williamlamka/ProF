package handler

import (
	"net/http"
	"new_project/utils"

	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.JSON(code, utils.ErrorResponse{
		Success: false,
		Message: err.Error(),
	})
}
