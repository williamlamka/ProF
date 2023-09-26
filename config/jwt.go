package config

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	ID  	uuid.UUID 	`json:"id"`
	Email 	string   	`json:"email"`
	jwt.RegisteredClaims
}

var jwtConfig = echojwt.Config {
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		SigningKey: []byte(os.Getenv("JWT")),
}

func JwtConfig() echojwt.Config{
	return jwtConfig
}

func ExtractJwt(c echo.Context) *JwtCustomClaims {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*JwtCustomClaims)
	return claims
}
