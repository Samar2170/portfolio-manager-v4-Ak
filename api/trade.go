package api

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/samar2170/portfolio-manager-v4/internal"
	"github.com/samar2170/portfolio-manager-v4/internal/models"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pbond"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pets"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pmutualfund"
	portfoliobase "github.com/samar2170/portfolio-manager-v4/internal/portfolio/portfolio-base"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pstock"
	"github.com/samar2170/portfolio-manager-v4/pkg/response"
)

func registerTrade(c echo.Context) error {
	security := c.Param("security")
	user := c.Get("user").(models.User)
	var trade portfoliobase.TradeInterface
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

// trades can be filtered by security type and pagination

func getTradeFilters(c *echo.Context) portfolio.TradeFilters {
	var tf portfolio.TradeFilters
	securityParam := (*c).Param("security")
	securitiesSplit := strings.Split(securityParam, ",")
	var securities []string
	for _, s := range securitiesSplit {
		securities = append(securities, s)
	}
	tf.Security = securities
	tf.SortBy = (*c).QueryParam("sort_by")
	pageNo := (*c).QueryParam("page")
	pageNoParsed, err := strconv.ParseInt(pageNo, 10, 64)
	if err != nil {
		tf.Page = 1
	}
	if tf.Page == 0 {
		tf.Page = int(pageNoParsed)
	}
	return tf
}

func listTrades(c echo.Context) error {
	user := c.Get("user").(models.User)
	tf := getTradeFilters(&c)
	trades, err := portfolio.GetTrades(tf, user.UserCID)
	if err != nil {
		return c.JSON(response.BadRequestResponseEcho(err.Error()))
	}
	return c.JSON(response.JSONResponseEcho(trades))
}
