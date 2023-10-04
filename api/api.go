package api

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samar2170/portfolio-manager-v4/api/analytics"
	"github.com/samar2170/portfolio-manager-v4/pkg/mw"
	"github.com/spf13/viper"
)

var SigningKey string

func init() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	SigningKey = viper.GetString("SIGNING_KEY")
}

type TemplateRenderer struct {
	template *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	return t.template.ExecuteTemplate(w, name, data)
}

func StartServer() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(mw.ApiKeyMiddleware())

	subroute := e.Group("/api/v1")

	subroute.POST("/signup", signup)
	subroute.POST("/login", login)
	subroute.GET("/generate-api-key", generateApiKey)

	subroute.POST("/register-account/:accountType", registerAccount)
	subroute.GET("/account/:accountType/", listAccounts)

	subroute.POST("/register-trade/:security/", registerTrade)
	subroute.GET("/list-trade/:security/", listTrades)

	subroute.GET("/list-holding/:security/", listHoldings)

	subroute.GET("/download-trade-template", downloadTradeTemplate)
	subroute.POST("/upload-trade-sheet", uploadTradeSheet)

	renderer := &TemplateRenderer{
		template: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = renderer

	analyticSubroute := e.Group("/api/v1/analytics")
	analyticSubroute.Use(mw.ApiKeyMiddleware())
	analyticSubroute.GET("/view-trades/:security/", analytics.ViewTrades).Name = "view-trades"
	analyticSubroute.GET("/view-holdings/:security/", analytics.ViewHoldings).Name = "view-holdings"
	e.Logger.Fatal(e.Start(":8080"))

}
