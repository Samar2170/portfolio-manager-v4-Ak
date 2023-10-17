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

func addSecurity(c echo.Context) error {
	security := c.Param("security")
	switch security {
	case "mutual-fund":
		rMf := new(mutualfund.MutualFund)
		if err := c.Bind(rMf); err != nil {
			return c.JSON(response.BadRequestResponseEcho("Invalid Request Body  " + err.Error()))
		}
		err := mutualfund.CreateMutualFund(*rMf)
		if err != nil {
			return c.JSON(response.BadRequestResponseEcho("While creating Mutual Fund  " + err.Error()))
		}
		return c.JSON(response.SuccessResponseEcho("Mutual Fund Created Successfully"))
	case "ets":
		rEts := new(ets.ETS)
		if err := c.Bind(rEts); err != nil {
			return c.JSON(response.BadRequestResponseEcho("Invalid Request Body  " + err.Error()))
		}
		err := ets.CreateETS(*rEts)
		if err != nil {
			return c.JSON(response.BadRequestResponseEcho("While creating ETS  " + err.Error()))
		}
		return c.JSON(response.SuccessResponseEcho("ETS Created Successfully"))
	case "bond":
		rbond := new(bond.Bond)
		if err := c.Bind(rbond); err != nil {
			return c.JSON(response.BadRequestResponseEcho("Invalid Request Body  " + err.Error()))
		}
		err := bond.CreateNewBond(*rbond)
		if err != nil {
			return c.JSON(response.BadRequestResponseEcho("While creating bond  " + err.Error()))
		}
		return c.JSON(response.SuccessResponseEcho("Bond Created Successfully"))
	default:
		return c.JSON(response.NotFoundResponseEcho("Security Not Found"))
	}
}
