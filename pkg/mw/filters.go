package mw

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio"
)

func GetTradeFilters(c *echo.Context) portfolio.TradeFilters {
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
