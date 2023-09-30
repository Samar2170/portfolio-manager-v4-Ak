package pets

import (
	"errors"
	"time"

	"github.com/samar2170/portfolio-manager-v4/internal"
	"github.com/samar2170/portfolio-manager-v4/internal/models"
	"github.com/samar2170/portfolio-manager-v4/pkg/db"
	"github.com/samar2170/portfolio-manager-v4/pkg/utils"
	"github.com/samar2170/portfolio-manager-v4/security/ets"
	"github.com/samar2170/portfolio-manager-v4/security/stock"
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
	ETSID    int
	ETS      *stock.Stock
	Quantity int
	BuyPrice float64
	Account  models.DematAccount
}

func (e *ETSHolding) create() error {
	return db.DB.Create(e).Error
}

func (e *ETSHolding) update() error {
	return db.DB.Save(e).Error
}
func (e *ETSHolding) getInvestedValue() float64 {
	return float64(e.Quantity) * e.BuyPrice
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
func getETSHolding(etsId int, userCID string) (ETSHolding, error) {
	var etsHolding ETSHolding
	dematAccounts, _ := models.GetDematAccountsByUserCID(userCID)
	dematIds := make([]int, len(dematAccounts))
	for i, account := range dematAccounts {
		dematIds[i] = account.ID
	}
	err := db.DB.Where("ets_id = ? AND account_id IN ?", etsId, dematIds).First(&etsHolding).Error
	return etsHolding, err
}

func etsHoldingExists(etsId int, userCID string) bool {
	return db.DB.Where("ets_id = ? AND account_id IN ?", etsId, userCID).First(&ETSHolding{}).Error == nil
}

func RegisterETSTrade(e *ETSTrade) error {
	err := e.create()
	if err != nil {
		return err
	}
	existingHolding := etsHoldingExists(e.ETSID, e.Account.UserCID)
	if existingHolding {
		holding, err := getETSHolding(e.ETSID, e.Account.UserCID)
		if err != nil {
			return err
		}
		if e.TradeType == "buy" {
			holding.Quantity += e.Quantity
			holding.BuyPrice = (holding.getInvestedValue() + e.GetInvestedValue()) / (float64(holding.Quantity) + float64(e.Quantity))
		} else {
			holding.Quantity -= e.Quantity
		}
		err = holding.update()
		if err != nil {
			return err
		}
	} else {
		if e.TradeType == "sell" {
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
