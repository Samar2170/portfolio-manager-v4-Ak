package stock

import (
	"log"
	"time"

	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/db"
	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/utils"
	"github.com/samar2170/portfolio-manager-v4-Ak/security/fetch"
)

func UpdateStockPrices() error {
	stocks, err := getPriceUpdatePendingStocks()
	if err != nil {
		return err
	}
	for _, stock := range stocks {
		updateStockPrice(stock)
	}
	return nil
}

func updateStockPrice(stock Stock) error {
	data, err := fetch.FetchAVStockData(stock.Symbol, "BSE")
	if err != nil {
		log.Println(err, "Error fetching stock data for ", stock.Symbol)
	}
	var latestDate time.Time
	latest, latestErr := stock.GetLatestDate()
	if latestErr != nil {
		latestDate = time.Now().AddDate(-7, 0, 0)
	} else {
		latestDate = latest
	}

	for ts, dp := range data.TimeSeries {
		t, err := utils.ParseTime(ts, fetch.AvTimeFormat)
		if err != nil {
			log.Println(err, "Error parsing time for ", stock.Symbol)
		}
		if t.After(latestDate) {
			stockPriceHistory := StockPriceHistory{
				Stock:  stock,
				Price:  dp.Close,
				Volume: dp.Volume,
				Date:   t,
			}
			stockPriceHistory.Create()
			log.Println("Created stock price history for ", stock.Symbol, " on ", t)
		}
	}
	stock.PriceToBeUpdated = false
	stock.PriceToBeUpdatedRank = 0
	db.DB.Save(&stock)
	return nil
}
