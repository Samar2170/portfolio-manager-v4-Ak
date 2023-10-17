package api

import (
	"github.com/labstack/echo/v4"
	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/response"
	"github.com/samar2170/portfolio-manager-v4-Ak/security/bond"
	"github.com/samar2170/portfolio-manager-v4-Ak/security/ets"
	mutualfund "github.com/samar2170/portfolio-manager-v4-Ak/security/mutual-fund"
	"github.com/samar2170/portfolio-manager-v4-Ak/security/stock"
)

func getSecurity(c echo.Context) error {
	security := c.Param("security")
	query := c.QueryParam("query")
	switch security {
	case "mutual-fund":
		mfs, err := mutualfund.SearchMutualFund(query)
		if err != nil {
			return c.JSON(response.NotFoundResponseEcho("Mutual Fund Not Found  " + err.Error()))
		}
		return c.JSON(response.JSONResponseEcho(mfs))
	case "stock":
		stocks, err := stock.SearchStock(query)
		if err != nil {
			return c.JSON(response.NotFoundResponseEcho("Stock Not Found  " + err.Error()))
		}
		return c.JSON(response.JSONResponseEcho(stocks))
	case "ets":
		ets, err := ets.SearchETS(query)
		if err != nil {
			return c.JSON(response.NotFoundResponseEcho("ETS Not Found  " + err.Error()))
		}
		return c.JSON(response.JSONResponseEcho(ets))
	case "bond":
		bonds, err := bond.SearchBond(query)
		if err != nil {
			return c.JSON(response.NotFoundResponseEcho("Bond Not Found  " + err.Error()))
		}
		return c.JSON(response.JSONResponseEcho(bonds))
	default:
		return c.JSON(response.NotFoundResponseEcho("Security Not Found"))
	}

}
