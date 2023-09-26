package mw

import (
	"github.com/labstack/echo/v4"
	"github.com/samar2170/portfolio-manager-v4/pkg/utils"
)

var ExemptedPaths = []string{
	"/api/v1/generate-api-key",
	"/api/v1/signup",
	"/api/v1/login",
}

type Credentials struct {
	ApiKey string `json:"api_key"`
}

func ApiKeyMiddleware(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// var apiKey string
			if c.Request().Method == "OPTIONS" {
				return next(c)
			}
			if utils.ArrayContains(ExemptedPaths, c.Request().URL.Path) {
				return next(c)
			}
			// if c.Request().Method == "POST" {
			// 	apiKey = c.Bind(Credentials)
			// }

			// authHeader := c.Request().Header.Get("Authorization")
			// if authHeader == "" {
			// 	return c.JSON(response.UnauthorizedResponseEcho("Missing Authorization Header"))
			// }
			// tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
			// token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			// 	return []byte(secretKey), nil
			// })
			// if err != nil {
			// 	return c.JSON(response.UnauthorizedResponseEcho(err.Error()))
			// }
			// claims, ok := token.Claims.(*JwtCustomClaims)
			// if !ok {
			// 	return c.JSON(response.UnauthorizedResponseEcho("Invalid Token"))
			// }
			// log.Println(claims)
			// c.Set("user_cid", claims.UserCID)
			// c.Set("username", claims.Username)
			// return next(c)
			return next(c)
		}
	}
}
