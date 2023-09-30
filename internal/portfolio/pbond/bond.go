package pbond

import (
	"errors"
	"strings"
	"time"

	"github.com/samar2170/portfolio-manager-v4/internal"
	"github.com/samar2170/portfolio-manager-v4/internal/models"
	"github.com/samar2170/portfolio-manager-v4/pkg/db"
	"github.com/samar2170/portfolio-manager-v4/pkg/utils"
	"github.com/samar2170/portfolio-manager-v4/security/bond"
	"gorm.io/gorm"
)

func init() {
	db.DB.AutoMigrate(&BondTrade{}, &BondHolding{})
}

type BondTrade struct {
	*gorm.Model
	ID        int
	BondID    int
	Bond      *bond.Bond
	Quantity  int
	Price     float64
	TradeType string
	TradeDate time.Time
	Account   models.DematAccount
	AccountID int
}

func NewBondTrade(symbol string, quantity int, price float64, tradeDate, tradeType, accountCode, userCID string) (*BondTrade, error) {
	bond, err := bond.GetBondBySymbol(symbol)
	if err != nil {
		return nil, err
	}
	t, err := utils.ParseTime(tradeDate, internal.TradeDateFormat)
	if err != nil {
		return nil, err
	}
	account, err := models.GetDematAccountByCode(accountCode, userCID)
	if err != nil {
		return nil, err
	}

	return &BondTrade{
		BondID:    bond.ID,
		Quantity:  quantity,
		Price:     price,
		TradeType: tradeType,
		TradeDate: t,
		Account:   account,
		AccountID: account.ID,
	}, nil
}

type BondHolding struct {
	*gorm.Model
	BondID   int
	Bond     *bond.Bond
	Quantity int
	BuyPrice float64
	Account  models.DematAccount
}

func (b *BondTrade) create() error {
	return db.DB.Create(b).Error
}
func (b *BondTrade) GetAccount() models.DematAccount {
	return b.Account
}
func (b *BondTrade) GetInvestedValue() float64 {
	return b.Price * float64(b.Quantity)
}
func (b *BondHolding) create() error {
	return db.DB.Create(b).Error
}

func (b *BondHolding) update() error {
	return db.DB.Save(b).Error
}

func (b *BondHolding) getInvestedValue() float64 {
	return float64(b.Quantity) * b.BuyPrice
}

func RegisterBondTrade(b *BondTrade) error {
	err := b.create()
	if err != nil {
		return err
	}
	existingHolding := bondHoldingExists(b.BondID, b.AccountID)
	if existingHolding {
		holding, err := getBondHolding(b.BondID, b.AccountID)
		if err != nil {
			return err
		}
		if strings.ToLower(b.TradeType) == "buy" {
			holding.Quantity += b.Quantity
			holding.BuyPrice = (holding.getInvestedValue() + b.GetInvestedValue()) / (float64(holding.Quantity) + float64(b.Quantity))
		} else {
			holding.Quantity -= b.Quantity
		}
		err = holding.update()
		if err != nil {
			return err
		}
	} else {
		if strings.ToLower(b.TradeType) == "sell" {
			return errors.New("cannot sell bond that you do not own")
		} else {
			holding := BondHolding{
				BondID:   b.BondID,
				Quantity: b.Quantity,
				BuyPrice: b.Price,
				Account:  b.Account,
			}
			err := holding.create()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getBondHolding(bondId int, accountID int) (BondHolding, error) {
	var bondHolding BondHolding
	err := db.DB.Where("bond_id = ? AND account_id = ?", bondId, accountID).First(&bondHolding).Error
	return bondHolding, err
}

func bondHoldingExists(bondId int, accountID int) bool {
	return db.DB.Where("bond_id = ? AND account_id = ?", bondId, accountID).First(&BondHolding{}).Error == nil
}
