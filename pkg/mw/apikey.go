package mw

import (
	"github.com/labstack/echo/v4"
	"github.com/samar2170/portfolio-manager-v4-Ak/internal"
	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/response"
	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/utils"
)

var ExemptedPaths = []string{
	"/api/v1/generate-api-key",
	"/api/v1/signup",
	"/api/v1/login",
}

type Credentials struct {
	ApiKey string `json:"api_key"`
}

var credentials map[string]string = map[string]string{}

func ApiKeyMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var creds Credentials
			if c.Request().Method == "OPTIONS" {
				return next(c)
			}
			if utils.ArrayContains(ExemptedPaths, c.Request().URL.Path) {
				return next(c)
			}
			if c.Request().Method == "POST" || c.Request().Method == "PUT" {
				api_key := c.Request().Header.Get("api_key")
				creds.ApiKey = api_key
			}
			if c.Request().Method == "GET" {
				creds.ApiKey = c.QueryParam("api_key")
			}
			if creds == (Credentials{}) || creds.ApiKey == "" {
				return c.JSON(response.UnauthorizedResponseEcho("Missing Api Key"))
			}

			user, err := internal.GetUserByApiKey(creds.ApiKey)
			if err != nil {
				return c.JSON(response.InternalServerErrorResponseEcho(err.Error()))
			}
			c.Set("user", user)
			return next(c)
		}
	}
}
