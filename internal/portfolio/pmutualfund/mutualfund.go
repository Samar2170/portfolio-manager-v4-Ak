package pmutualfund

import (
	"errors"
	"strings"
	"time"

	"github.com/samar2170/portfolio-manager-v4-Ak/internal"
	"github.com/samar2170/portfolio-manager-v4-Ak/internal/models"
	portfoliobase "github.com/samar2170/portfolio-manager-v4-Ak/internal/portfolio/portfolio-base"
	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/db"
	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/utils"
	mutualfund "github.com/samar2170/portfolio-manager-v4-Ak/security/mutual-fund"
	"gorm.io/gorm"
)

func init() {
	db.DB.AutoMigrate(&MutualFundTrade{}, &MutualFundHolding{})
}

type MutualFundTrade struct {
	*gorm.Model
	ID           int
	MutualFund   mutualfund.MutualFund
	MutualFundID int
	Quantity     float64
	Price        float64
	TradeType    string
	TradeDate    time.Time
	Account      models.DematAccount
	AccountID    int
}

func NewMutualFundTrade(mutualFundID int, quantity float64, price float64, tradeDate, tradeType, accountCode, userCID string) (*MutualFundTrade, error) {
	mutualFund, err := mutualfund.GetMutualFundByID(mutualFundID)
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

	return &MutualFundTrade{
		MutualFundID: mutualFund.ID,
		Quantity:     float64(quantity),
		Price:        price,
		TradeType:    tradeType,
		TradeDate:    tradeDateParsed,
		Account:      account,
		AccountID:    account.ID,
	}, nil
}

type MutualFundHolding struct {
	*gorm.Model
	ID           int
	MutualFund   mutualfund.MutualFund
	MutualFundID int
	Quantity     float64
	BuyPrice     float64
	Account      models.DematAccount
	AccountID    int
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
func (m *MutualFundTrade) GetTradeData() portfoliobase.TradeData {
	return portfoliobase.TradeData{
		Symbol:        m.MutualFund.SchemeNavName,
		Quantity:      m.Quantity,
		Price:         m.Price,
		InvestedValue: m.GetInvestedValue(),
		Date:          m.TradeDate,
	}

}

func (mf *MutualFundHolding) create() error {
	return db.DB.Create(mf).Error
}

func (mf *MutualFundHolding) update() error {
	return db.DB.Save(mf).Error
}
func (mf *MutualFundHolding) GetInvestedValue() float64 {
	return float64(mf.Quantity) * mf.BuyPrice
}

func RegisterMutualFundTrade(m *MutualFundTrade) error {
	err := m.create()
	if err != nil {
		return err
	}
	existingHolding := mutualFundHoldingExists(m.MutualFundID, m.AccountID)
	if existingHolding {
		holding, err := getMutualFundHolding(m.MutualFundID, m.AccountID)
		if err != nil {
			return err
		}
		if strings.ToLower(m.TradeType) == "buy" {
			holding.Quantity += m.Quantity
			holding.BuyPrice = (holding.GetInvestedValue() + m.GetInvestedValue()) / (float64(holding.Quantity) + float64(m.Quantity))
		} else {
			holding.Quantity -= m.Quantity
		}
		err = holding.update()
		if err != nil {
			return err
		}
	} else {
		if strings.ToLower(m.TradeType) == "sell" {
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

func getMutualFundHolding(mfId int, accountID int) (MutualFundHolding, error) {
	var mfHolding MutualFundHolding
	err := db.DB.Where("mutual_fund_id = ? AND account_id = ?", mfId, accountID).First(&mfHolding).Error
	return mfHolding, err
}

func mutualFundHoldingExists(mfId int, accountID int) bool {
	return db.DB.Where("mutual_fund_id = ? AND account_id = ?", mfId, accountID).First(&MutualFundHolding{}).Error == nil
}
