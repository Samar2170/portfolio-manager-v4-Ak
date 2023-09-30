package api

import (
	"github.com/labstack/echo/v4"
	"github.com/samar2170/portfolio-manager-v4/internal"
	"github.com/samar2170/portfolio-manager-v4/internal/models"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pbond"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pets"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pmutualfund"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pstock"
	"github.com/samar2170/portfolio-manager-v4/pkg/response"
)

func registerTrade(c echo.Context) error {
	security := c.Param("security")
	user := c.Get("user").(models.User)
	var trade portfolio.TradeInterface
	var err error
	switch security {
	case "stock":
		var str internal.StockTradeRequest
		if err = c.Bind(&str); err != nil {
			return c.JSON(response.BadRequestResponseEcho(err.Error()))
		}
		trade, err = pstock.NewStockTrade(str.Symbol, str.Quantity, str.Price, str.TradeDate, str.TradeType, str.AccountCode, user.UserCID)
		if err != nil {
			return c.JSON(response.BadRequestResponseEcho(err.Error()))
		}
	case "bond":
		var btr internal.BondTradeRequest
		if err = c.Bind(&btr); err != nil {
			return c.JSON(response.BadRequestResponseEcho(err.Error()))
		}
		trade, err = pbond.NewBondTrade(btr.Symbol, btr.Quantity, btr.Price, btr.TradeDate, btr.TradeType, btr.AccountCode, user.UserCID)
		if err != nil {
			return c.JSON(response.BadRequestResponseEcho(err.Error()))
		}
	case "mutual-fund":
		var mftr internal.MutualFundTradeRequest
		if err = c.Bind(&mftr); err != nil {
			return c.JSON(response.BadRequestResponseEcho(err.Error()))
		}
		trade, err = pmutualfund.NewMutualFundTrade(mftr.MutualFundID, mftr.Quantity, mftr.Price, mftr.TradeDate, mftr.TradeType, mftr.AccountCode, user.UserCID)
		if err != nil {
			return c.JSON(response.BadRequestResponseEcho(err.Error()))
		}
	case "ets":
		var etr internal.ETSTradeRequest
		if err = c.Bind(&etr); err != nil {
			return c.JSON(response.BadRequestResponseEcho(err.Error()))
		}
		trade, err = pets.NewETSTrade(etr.Symbol, etr.Quantity, etr.Price, etr.TradeDate, etr.TradeType, etr.AccountCode, user.UserCID)
		if err != nil {
			return c.JSON(response.BadRequestResponseEcho(err.Error()))
		}
	}
	err = portfolio.RegisterTrade(trade)
	if err != nil {
		return c.JSON(response.BadRequestResponseEcho(err.Error()))
	}
	return c.JSON(response.SuccessResponseEcho("trade registered"))
}
