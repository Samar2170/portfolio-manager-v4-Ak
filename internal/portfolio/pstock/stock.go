package pstock

import (
	"errors"
	"strings"
	"time"

	"github.com/samar2170/portfolio-manager-v4-Ak/internal"
	"github.com/samar2170/portfolio-manager-v4-Ak/internal/models"
	portfoliobase "github.com/samar2170/portfolio-manager-v4-Ak/internal/portfolio/portfolio-base"
	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/db"
	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/utils"
	"github.com/samar2170/portfolio-manager-v4-Ak/security/stock"
	"gorm.io/gorm"
)

func init() {
	db.DB.AutoMigrate(&StockTrade{}, &StockHolding{})
}

type StockTrade struct {
	*gorm.Model
	ID        int
	Stock     stock.Stock `gorm:"foreignKey:StockID"`
	StockID   int
	Quantity  int
	Price     float64
	TradeType string
	TradeDate time.Time
	Account   models.DematAccount `gorm:"foreignKey:AccountID"`
	AccountID int
}

func NewStockTrade(symbol string, quantity int, price float64, tradeDate, tradeType, accountCode, userCID string) (*StockTrade, error) {
	stock, err := stock.GetStockBySymbol(symbol)
	if err != nil {
		return nil, err
	}
	if stock.Symbol == "" {
		return nil, errors.New("stock does not exist")
	}
	t, err := utils.ParseTime(tradeDate, internal.TradeDateFormat)
	if err != nil {
		return nil, err
	}
	account, err := models.GetDematAccountByCode(accountCode, userCID)
	if err != nil {
		return nil, err
	}

	return &StockTrade{
		StockID:   stock.ID,
		Quantity:  quantity,
		Price:     price,
		TradeType: tradeType,
		TradeDate: t,
		Account:   account,
		AccountID: account.ID,
	}, nil
}

type StockHolding struct {
	*gorm.Model
	ID        uint
	Stock     stock.Stock `gorm:"foreignKey:StockID"`
	StockID   int
	Quantity  int
	BuyPrice  float64
	Account   models.DematAccount `gorm:"foreignKey:AccountID"`
	AccountID int
}

func (s *StockTrade) create() error {
	return db.DB.Create(s).Error
}
func (s *StockTrade) GetAccount() models.DematAccount {
	return s.Account
}
func (s *StockTrade) GetInvestedValue() float64 {
	return s.Price * float64(s.Quantity)
}
func (s *StockTrade) GetTradeData() portfoliobase.TradeData {
	return portfoliobase.TradeData{
		Symbol:        s.Stock.Symbol,
		Quantity:      float64(s.Quantity),
		Price:         s.Price,
		InvestedValue: s.GetInvestedValue(),
		Date:          s.TradeDate,
	}
}

func (s *StockHolding) create() error {
	return db.DB.Create(s).Error
}
func (s *StockHolding) update() error {
	return db.DB.Save(s).Error
}
func (s *StockHolding) GetInvestedValue() float64 {
	return float64(s.Quantity) * s.BuyPrice
}

func RegisterStockTrade(s *StockTrade) error {
	err := s.create()
	if err != nil {
		return err
	}
	existingHolding := stockHoldingExists(s.StockID, s.AccountID)
	if existingHolding {
		holding, err := getStockHolding(s.StockID, s.AccountID)
		if err != nil {
			return err
		}
		if strings.ToLower(s.TradeType) == "buy" {
			holding.Quantity += s.Quantity
			holding.BuyPrice = (holding.GetInvestedValue() + s.GetInvestedValue()) / (float64(holding.Quantity) + float64(s.Quantity))
		} else {
			holding.Quantity -= s.Quantity
		}
		err = holding.update()
		if err != nil {
			return err
		}
	} else {
		if strings.ToLower(s.TradeType) == "sell" {
			return errors.New("cannot sell stock that you do not own")
		} else {
			holding := StockHolding{
				StockID:  s.StockID,
				Quantity: s.Quantity,
				BuyPrice: s.Price,
				Account:  s.Account,
			}
			err := holding.create()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getStockHolding(stockId int, accountID int) (StockHolding, error) {
	var stockHolding StockHolding
	err := db.DB.Where("stock_id = ? AND account_id = ?", stockId, accountID).First(&stockHolding).Error
	return stockHolding, err
}

func stockHoldingExists(stockId int, accountID int) bool {
	return db.DB.Where("stock_id = ? AND account_id = ?", stockId, accountID).First(&StockHolding{}).Error == nil
}
