package handler

import (
	"errors"
	"net/http"
	"new_project/config"
	"new_project/utils"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserRegisterDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	UserName string `json:"username" validate:"required"`
}

type UserLoginDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func Register(c echo.Context) error {
	var dto UserRegisterDto
	err := c.Bind(&dto)
	if err != nil {
		return err
	}
	err = validator.New().Struct(dto)
	if err != nil {
		return err
	}
	_, err = GetUserByEmail(dto.Email)
	if err == nil {
		return errors.New("email is used")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	dto.Password = string(hashedPassword)
	newUser, _ := CreateUser(dto)
	return c.JSON(http.StatusOK, utils.SuccessResponse{
		Success: true,
		Data:    newUser,
	})
}

func Login(c echo.Context) error {
	var dto UserLoginDto
	err := c.Bind(&dto)
	if err != nil {
		return err
	}
	err = validator.New().Struct(dto)
	if err != nil {
		return err
	}
	user, err := GetUserByEmail(dto.Email)
	if err != nil {
		return errors.New("no such user")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		return errors.New("password incorrect")
	}
	claims := &config.JwtCustomClaims{
		ID:    user.ID,
		Email: user.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString([]byte(os.Getenv("JWT")))
	Subscribe(user.ID)
	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
		"token":   t,
	})
}
