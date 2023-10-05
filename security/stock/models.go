package stock

import (
	"time"

	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/db"
	"gorm.io/gorm"
)

func init() {
	db.DB.AutoMigrate(&Stock{})
	db.DB.AutoMigrate(&StockPriceHistory{})
}

type Stock struct {
	*gorm.Model
	ID                   int
	Symbol               string `gorm:"index"`
	Name                 string
	Industry             string
	Exchange             string
	SecurityCode         string `gorm:"index;"`
	PriceToBeUpdated     bool
	PriceToBeUpdatedRank int
}

type StockPriceHistory struct {
	*gorm.Model
	ID      int
	Stock   Stock
	StockID int
	Price   float64
	Volume  float64
	Date    time.Time
	Source  string
}

func (s *StockPriceHistory) Create() error {
	err := db.DB.Create(&s).Error
	return err
}

func (s *Stock) Create() error {
	err := db.DB.Create(&s).Error
	return err
}

func (s *Stock) GetOrCreate() (Stock, error) {
	err := db.DB.FirstOrCreate(&s, Stock{Symbol: s.Symbol, Exchange: s.Exchange}).Error
	return *s, err
}

func (s *Stock) GetLatestPrice() (float64, error) {
	var sph StockPriceHistory
	err := db.DB.Find(&sph, "stock_id = ?", s.ID).Order("created_at desc").Limit(1).Error
	return sph.Price, err
}

func (s *Stock) GetLatestDate() (time.Time, error) {
	var sph StockPriceHistory
	err := db.DB.Find(&sph, "stock_id = ?", s.ID).Order("created_at desc").Limit(1).Error
	return sph.Date, err
}

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
	err := db.DB.Find(&stocks, "price_to_be_updated = ?", true).Order("price_to_be_updated desc").Error
	return stocks, err
}

func GetStockBySymbol(symbol string) (Stock, error) {
	var stock Stock
	err := db.DB.Find(&stock, "symbol = ?", symbol).Error
	return stock, err
}
