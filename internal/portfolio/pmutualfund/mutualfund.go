package pmutualfund

import (
	"errors"
	"strconv"
	"time"

	"github.com/samar2170/portfolio-manager-v4/internal"
	"github.com/samar2170/portfolio-manager-v4/internal/models"
	"github.com/samar2170/portfolio-manager-v4/pkg/db"
	mutualfund "github.com/samar2170/portfolio-manager-v4/security/mutual-fund"
	"gorm.io/gorm"
)

func init() {
	db.DB.AutoMigrate(&MutualFundTrade{}, &MutualFundHolding{})
}

type MutualFundTrade struct {
	*gorm.Model
	ID           int
	MutualFundID int
	MutualFund   *mutualfund.MutualFund
	Quantity     float64
	Price        float64
	TradeType    string
	TradeDate    time.Time
	Account      models.DematAccount
}

func NewMutualFundTrade(mutualFundID int, quantity, price, tradeDate, tradeType string) (*MutualFundTrade, error) {
	mutualFund, err := mutualfund.GetMutualFundByID(mutualFundID)
	if err != nil {
		return nil, err
	}
	tradeDateParsed, err := time.Parse(tradeDate, internal.TradeDateFormat)
	if err != nil {
		return nil, errors.New("invalid trade date")
	}
	quantityParsed, err := strconv.ParseFloat(quantity, 64)
	if err != nil {
		return nil, err
	}
	priceParsed, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return nil, err
	}
	return &MutualFundTrade{
		MutualFundID: mutualFund.ID,
		Quantity:     quantityParsed,
		Price:        priceParsed,
		TradeType:    tradeType,
		TradeDate:    tradeDateParsed,
	}, nil
}

type MutualFundHolding struct {
	*gorm.Model
	MutualFundID int
	MutualFund   *mutualfund.MutualFund
	Quantity     float64
	BuyPrice     float64
	Account      models.DematAccount
}

func (m *MutualFundTrade) create() error {
	return db.DB.Create(m).Error
}
func (m *MutualFundTrade) GetAccount() models.DematAccount {
	return m.Account
}
func (m *MutualFundTrade) GetInvestedValue() float64 {
	return m.Price * float64(m.Quantity)
}
func (mf *MutualFundHolding) create() error {
	return db.DB.Create(mf).Error
}

func (mf *MutualFundHolding) update() error {
	return db.DB.Save(mf).Error
}
func (mf *MutualFundHolding) getInvestedValue() float64 {
	return float64(mf.Quantity) * mf.BuyPrice
}

func getMutualFundHolding(mfId int, userCID string) (MutualFundHolding, error) {
	var mfHolding MutualFundHolding
	dematAccounts, _ := models.GetDematAccountsByUserCID(userCID)
	dematIds := make([]int, len(dematAccounts))
	for i, account := range dematAccounts {
		dematIds[i] = account.ID
	}
	err := db.DB.Where("mutual_fund_id = ? AND account_id IN ?", mfId, dematIds).First(&mfHolding).Error
	return mfHolding, err
}

func mutualFundHoldingExists(mfId int, userCID string) bool {
	return db.DB.Where("mutual_fund_id = ? AND account_id IN ?", mfId, userCID).First(&MutualFundHolding{}).Error == nil
}
func RegisterMutualFundTrade(m *MutualFundTrade) error {
	err := m.create()
	if err != nil {
		return err
	}
	existingHolding := mutualFundHoldingExists(m.MutualFundID, m.Account.UserCID)
	if existingHolding {
		holding, err := getMutualFundHolding(m.MutualFundID, m.Account.UserCID)
		if err != nil {
			return err
		}
		if m.TradeType == "buy" {
			holding.Quantity += m.Quantity
			holding.BuyPrice = (holding.getInvestedValue() + m.GetInvestedValue()) / (float64(holding.Quantity) + float64(m.Quantity))
		} else {
			holding.Quantity -= m.Quantity
		}
		err = holding.update()
		if err != nil {
			return err
		}
	} else {
		if m.TradeType == "sell" {
			return errors.New("cannot sell mutual fund that you do not own")
		} else {
			holding := MutualFundHolding{
				MutualFundID: m.MutualFundID,
				Quantity:     m.Quantity,
				BuyPrice:     m.Price,
				Account:      m.Account,
			}
			err := holding.create()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
