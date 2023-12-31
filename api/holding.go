package api

import (
	"github.com/labstack/echo/v4"
	"github.com/samar2170/portfolio-manager-v4-Ak/internal/models"
	"github.com/samar2170/portfolio-manager-v4-Ak/internal/portfolio"
	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/mw"
	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/response"
)

func listHoldings(c echo.Context) error {
	user := c.Get("user").(models.User)
	tf := mw.GetTradeFilters(&c)
	holdings, err := portfolio.GetHoldings(tf, user.UserCID)
	if err != nil {
		return c.JSON(response.BadRequestResponseEcho(err.Error()))
	}
	return c.JSON(response.JSONResponseEcho(holdings))
}
