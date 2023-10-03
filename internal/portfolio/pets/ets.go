package pets

import (
	"errors"
	"strings"
	"time"

	"github.com/samar2170/portfolio-manager-v4/internal"
	"github.com/samar2170/portfolio-manager-v4/internal/models"
	portfoliobase "github.com/samar2170/portfolio-manager-v4/internal/portfolio/portfolio-base"
	"github.com/samar2170/portfolio-manager-v4/pkg/db"
	"github.com/samar2170/portfolio-manager-v4/pkg/utils"
	"github.com/samar2170/portfolio-manager-v4/security/ets"
	"gorm.io/gorm"
)

func init() {
	db.DB.AutoMigrate(&ETSTrade{}, &ETSHolding{})
}

type ETSTrade struct {
	*gorm.Model
	ID        int
	ETS       ets.ETS
	ETSID     int
	Quantity  int
	Price     float64
	TradeType string
	TradeDate time.Time
	Account   models.DematAccount
	AccountID int
}

func NewETSTrade(symbol string, quantity int, price float64, tradeDate, tradeType, accountCode, userCID string) (*ETSTrade, error) {
	ets, err := ets.GetETSBySymbol(symbol)
	if err != nil {
		return nil, err
	}
	tradeDateParsed, err := utils.ParseTime(tradeDate, internal.TradeDateFormat)
	if err != nil {
		return nil, errors.New("invalid trade date")
	}
	account, err := models.GetDematAccountByCode(accountCode, userCID)
	if err != nil {
		return nil, err
	}
	return &ETSTrade{
		ETSID:     ets.ID,
		Quantity:  quantity,
		Price:     price,
		TradeType: tradeType,
		TradeDate: tradeDateParsed,
		Account:   account,
		AccountID: account.ID,
	}, nil
}

type ETSHolding struct {
	*gorm.Model
	ID        int
	ETS       ets.ETS
	ETSID     int
	Quantity  int
	BuyPrice  float64
	Account   models.DematAccount
	AccountID int
}

func (e *ETSHolding) create() error {
	return db.DB.Create(e).Error
}

func (e *ETSHolding) update() error {
	return db.DB.Save(e).Error
}
func (e *ETSHolding) GetInvestedValue() float64 {
	return float64(e.Quantity) * e.BuyPrice
}
func (e *ETSTrade) GetTradeData() portfoliobase.TradeData {
	return portfoliobase.TradeData{
		Symbol:        e.ETS.Name,
		Quantity:      float64(e.Quantity),
		Price:         e.Price,
		InvestedValue: e.GetInvestedValue(),
		Date:          e.TradeDate,
	}
}

func (e *ETSTrade) create() error {
	return db.DB.Create(e).Error
}
func (e *ETSTrade) GetAccount() models.DematAccount {
	return e.Account
}
func (e *ETSTrade) GetInvestedValue() float64 {
	return e.Price * float64(e.Quantity)
}

func RegisterETSTrade(e *ETSTrade) error {
	err := e.create()
	if err != nil {
		return err
	}
	existingHolding := etsHoldingExists(e.ETSID, e.AccountID)
	if existingHolding {
		holding, err := getETSHolding(e.ETSID, e.AccountID)
		if err != nil {
			return err
		}
		if strings.ToLower(e.TradeType) == "buy" {
			holding.Quantity += e.Quantity
			holding.BuyPrice = (holding.GetInvestedValue() + e.GetInvestedValue()) / (float64(holding.Quantity) + float64(e.Quantity))
		} else {
			holding.Quantity -= e.Quantity
		}
		err = holding.update()
		if err != nil {
			return err
		}
	} else {
		if strings.ToLower(e.TradeType) == "sell" {
			return errors.New("cannot sell ets that you do not own")
		} else {
			holding := ETSHolding{
				ETSID:    e.ETSID,
				Quantity: e.Quantity,
				BuyPrice: e.Price,
				Account:  e.Account,
			}
			err := holding.create()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func getETSHolding(etsId int, accountID int) (ETSHolding, error) {
	var etsHolding ETSHolding
	err := db.DB.Where("ets_id = ? AND account_id = ?", etsId, accountID).First(&etsHolding).Error
	return etsHolding, err
}

func etsHoldingExists(etsId int, accountID int) bool {
	return db.DB.Where("ets_id = ? AND account_id = ?", etsId, accountID).First(&ETSHolding{}).Error == nil
}
