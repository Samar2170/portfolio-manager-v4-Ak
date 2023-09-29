package api

import (
	"log"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/samar2170/portfolio-manager-v4/internal"
	"github.com/samar2170/portfolio-manager-v4/pkg/response"
)

func signup(c echo.Context) error {
	signupRequest := new(internal.SignupRequest)

	if err := c.Bind(signupRequest); err != nil {
		return c.JSON(response.BadRequestResponseEcho(err.Error()))
	}
	err := internal.Signup(*signupRequest)
	if err != nil {
		return c.JSON(response.BadRequestResponseEcho(err.Error()))
	}
	return c.JSON(response.SuccessResponseEcho("User Created Successfully"))
}

func login(c echo.Context) error {
	loginRequest := new(internal.LoginRequest)
	if err := c.Bind(loginRequest); err != nil {
		return c.JSON(response.BadRequestResponseEcho(err.Error()))
	}
	if loginRequest.Username == "" || loginRequest.Password == "" {
		return c.JSON(response.BadRequestResponseEcho("Username or Password cannot be empty"))
	}
	token, err := internal.Login(*loginRequest)
	if err != nil {
		return c.JSON(response.BadRequestResponseEcho(err.Error()))
	}
	return c.JSON(response.SuccessResponseEcho(token))
}

func generateApiKey(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(response.UnauthorizedResponseEcho("Missing Authorization Header"))
	}
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	token := strings.Split(tokenString, " ")[1]
	log.Println("token-->", token)
	if tokenString == "" {
		return c.JSON(response.BadRequestResponseEcho("token is missing"))
	}
	apiKey, err := internal.GenerateApiKey(token)
	if err != nil {
		return c.JSON(response.InternalServerErrorResponseEcho(err.Error()))
	}
	return c.JSON(response.SuccessResponseEcho(apiKey))
}
