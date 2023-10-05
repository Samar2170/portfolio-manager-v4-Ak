package analytics

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samar2170/portfolio-manager-v4-Ak/internal/models"
	"github.com/samar2170/portfolio-manager-v4-Ak/internal/portfolio"
	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/mw"
)

func ViewTrades(c echo.Context) error {
	user := c.Get("user").(models.User)
	tf := mw.GetTradeFilters(&c)
	trades, err := portfolio.GetTrades(tf, user.UserCID)
	if err != nil {
		return c.Render(http.StatusBadRequest, "error.html", map[string]string{
			"error": err.Error(),
		})
	}
	tradeMaps := []map[string]interface{}{}
	for _, t := range trades {
		tradeMap := t.ToMap()
		tradeMaps = append(tradeMaps, tradeMap)
	}
	data := map[string]interface{}{
		"trades": tradeMaps,
	}
	return c.Render(http.StatusOK, "trades.html", data)
}

func ViewHoldings(c echo.Context) error {
	user := c.Get("user").(models.User)
	tf := mw.GetTradeFilters(&c)
	holdings, err := portfolio.GetHoldings(tf, user.UserCID)
	if err != nil {
		return c.Render(http.StatusBadRequest, "error.html", map[string]string{
			"error": err.Error(),
		})
	}
	holdingMaps := []map[string]interface{}{}
	for _, h := range holdings {
		holdingMap := h.ToMap()
		holdingMaps = append(holdingMaps, holdingMap)
	}
	data := map[string]interface{}{
		"holdings": holdingMaps,
	}
	return c.Render(http.StatusOK, "holdings.html", data)

}
