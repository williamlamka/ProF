package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"new_project/config"
	"new_project/handler"
	"new_project/model"
	"new_project/redis"
	"new_project/utils"

	"github.com/go-playground/validator/v10"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type UserResponse struct {
	Email    string			`json:"email"`
	Username string			`json:"userName"`
	Category pq.StringArray	`json:"category"`
}

func UserRouteInit(e *echo.Echo) {
	g := e.Group("/user")
	g.Use(echojwt.WithConfig(config.JwtConfig()))

	g.GET("", func(c echo.Context) error {
		claims := config.ExtractJwt(c)
		key := fmt.Sprintf("user:%s", claims.ID)
		val, err := redis.Get(key)
		if err != nil {
			return err
		}
		var data UserResponse
		if val != nil {
			err = json.Unmarshal(val, &data)
		} else {
			user, err := handler.GetUser(claims.ID)
			if err != nil{
				return err
			}
			jsonstr, _ := json.Marshal(user)
			err = redis.Set(key, string(jsonstr))
			if err != nil {
				return err
			}
			json.Unmarshal(jsonstr, &data)
		}
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, utils.SuccessResponse{
			Success: true,
			Data:    data,
		})
	})

	g.GET("/category", func(c echo.Context) error {
		claims := config.ExtractJwt(c)
		var category *pq.StringArray
		key := fmt.Sprintf("user:%s:category", claims.ID)
		val, err := redis.Get(key); if err != nil {
			return err
		}
		if val != nil {
			json.Unmarshal(val, &category)
		} else {
			category = handler.GetCategory(claims.ID)
			jsonstr, _ := json.Marshal(category)
			err := redis.Set(key, string(jsonstr))
			if err != nil {
				return err
			}
		}
		return c.JSON(http.StatusOK, utils.SuccessResponse{
			Success: true,
			Data:    category,
		})
	})

	g.POST("/update", func(c echo.Context) error {
		claims := config.ExtractJwt(c)
		key := fmt.Sprintf("user:%s", claims.ID)
		var dto model.UpdateUserDto
		err := c.Bind(&dto)
		if err != nil {
			return err
		}
		err = validator.New().Struct(dto)
		if err != nil {
			return err
		}
		err = handler.UpdateUser(claims.ID, dto)
		if err != nil {
			return err
		}
		redis.Delete(key)
		return c.JSON(http.StatusOK, utils.SuccessResponse{
			Success: true,
		})
	})

	g.POST("/category/add", func(c echo.Context) error {
		claims := config.ExtractJwt(c)
		var dto model.CategoryDto
		err := c.Bind(&dto)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("user:%s", claims.ID)
		key2 := fmt.Sprintf("user:%s:category", claims.ID)
		handler.AddCategory(claims.ID, dto.Category)
		redis.Delete(key)
		redis.Delete(key2)
		return c.JSON(http.StatusOK, utils.SuccessResponse{
			Success: true,
		})
	})

	g.POST("/category/remove", func(c echo.Context) error {
		claims := config.ExtractJwt(c)
		var dto model.CategoryDto
		err := c.Bind(&dto)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("user:%s", claims.ID)
		key2 := fmt.Sprintf("user:%s:category", claims.ID)
		handler.RemoveCategroy(claims.ID, dto.Category)
		redis.Delete(key)
		redis.Delete(key2)
		return c.JSON(http.StatusOK, utils.SuccessResponse{
			Success: true,
		})
	})
}
