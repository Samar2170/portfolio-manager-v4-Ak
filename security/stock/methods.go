package stock

import "github.com/samar2170/portfolio-manager-v4-Ak/pkg/db"

func GetAllStocks() ([]Stock, error) {
	var stocks []Stock
	err := db.DB.Find(&stocks).Error
	if err != nil {
		return stocks, err
	}
	return stocks, nil
}

func getPriceUpdatePendingStocks() ([]Stock, error) {
	var stocks []Stock
	err := db.DB.Find(&stocks, "price_to_be_updated = ?", true).Order("price_to_be_updated_rank desc").Error
	return stocks, err
}

func GetStockBySymbol(symbol string) (Stock, error) {
	var stock Stock
	err := db.DB.Find(&stock, "symbol = ?", symbol).Error
	return stock, err
}

func SearchStock(query string) ([]Stock, error) {
	var stocks []Stock
	var stocks2 []Stock
	err := db.DB.Where("symbol LIKE ?", "%"+query+"%").Find(&stocks).Error
	if err != nil {
		return nil, err
	}
	err = db.DB.Where("name LIKE ?", "%"+query+"%").Find(&stocks2).Error
	stocks = append(stocks, stocks2...)
	return stocks, nil
}
