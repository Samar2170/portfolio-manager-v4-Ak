package pstock

import (
	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/db"
	"github.com/samar2170/portfolio-manager-v4-Ak/security/stock"
)

type PriceToBeUpdatedResult struct {
	StockID       int
	TotalQuantity int
}

func CalculatePriceToBeUpdatedRank() error {
	var stockHoldings []PriceToBeUpdatedResult
	var err error
	err = db.DB.Model(&StockHolding{}).Where("stock_id > 0").Group("stock_id").Select("stock_id, sum(quantity) as total_quantity").Scan(&stockHoldings).Error
	if err != nil {
		return err
	}
	for _, stockHolding := range stockHoldings {
		var stock stock.Stock
		err = db.DB.Where("id = ?", stockHolding.StockID).First(&stock).Error
		if err != nil {
			return err
		}
		stock.PriceToBeUpdatedRank = stockHolding.TotalQuantity
		stock.PriceToBeUpdated = true
		db.DB.Save(&stock)
	}
	return nil
}
