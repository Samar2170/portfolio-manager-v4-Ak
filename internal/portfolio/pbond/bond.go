package pbond

import (
	"errors"
	"strconv"

	"github.com/samar2170/portfolio-manager-v4/internal/models"
	"github.com/samar2170/portfolio-manager-v4/pkg/db"
	"github.com/samar2170/portfolio-manager-v4/security/bond"
	"gorm.io/gorm"
)

type BondTrade struct {
	*gorm.Model
	ID        int
	BondID    int
	Bond      *bond.Bond
	Quantity  int
	Price     float64
	TradeType string
	TradeDate string
	Account   models.DematAccount
}

func NewBondTrade(symbol string, quantity, price, tradeDate, tradeType string) (*BondTrade, error) {
	bond, err := bond.GetBondBySymbol(symbol)
	if err != nil {
		return nil, err
	}
	quantityParsed, err := strconv.ParseInt(quantity, 10, 64)
	if err != nil {
		return nil, err
	}
	priceParsed, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return nil, err
	}
	return &BondTrade{
		BondID:    bond.ID,
		Quantity:  int(quantityParsed),
		Price:     priceParsed,
		TradeType: tradeType,
		TradeDate: tradeDate,
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

func getBondHolding(bondId int, userCID string) (BondHolding, error) {
	var bondHolding BondHolding
	dematAccounts, _ := models.GetDematAccountsByUserCID(userCID)
	dematIds := make([]int, len(dematAccounts))
	for i, account := range dematAccounts {
		dematIds[i] = account.ID
	}
	err := db.DB.Where("bond_id = ? AND account_id IN ?", bondId, dematIds).First(&bondHolding).Error
	return bondHolding, err
}

func bondHoldingExists(bondId int, userCID string) bool {
	return db.DB.Where("bond_id = ? AND account_id IN ?", bondId, userCID).First(&BondHolding{}).Error == nil
}

func RegisterBondTrade(b *BondTrade) error {
	err := b.create()
	if err != nil {
		return err
	}
	existingHolding := bondHoldingExists(b.BondID, b.Account.UserCID)
	if existingHolding {
		holding, err := getBondHolding(b.BondID, b.Account.UserCID)
		if err != nil {
			return err
		}
		if b.TradeType == "buy" {
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
		if b.TradeType == "sell" {
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
