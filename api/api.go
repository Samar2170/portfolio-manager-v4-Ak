package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

var SigningKey string

func init() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	SigningKey = viper.GetString("SIGNING_KEY")
}
func StartServer() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use(mw.JwtMiddleware(SigningKey))
	subroute := e.Group("/api/v1")
	subroute.POST("/signup", signup)
	subroute.POST("/login", login)
	subroute.GET("/generate-api-key", generateApiKey)

	subroute.POST("/register-account/:accountType", registerAccount)
	subroute.GET("/download-trade-template", downloadTradeTemplate)
	subroute.POST("/upload-trade-sheet", uploadTradeSheet)

	e.Logger.Fatal(e.Start(":8080"))

}
