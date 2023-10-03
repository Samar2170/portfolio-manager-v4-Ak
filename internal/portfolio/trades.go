package portfolio

import (
	"crypto/sha1"
	"fmt"

	"github.com/samar2170/portfolio-manager-v4/internal/models"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pbond"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pets"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pmutualfund"
	portfoliobase "github.com/samar2170/portfolio-manager-v4/internal/portfolio/portfolio-base"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pstock"
	"github.com/samar2170/portfolio-manager-v4/pkg/db"
	"github.com/samar2170/portfolio-manager-v4/pkg/utils"
	"github.com/samar2170/portfolio-manager-v4/pkg/utils/structs"
	"gorm.io/gorm"
)

var (
	limit = 50
)

type TradeFilters struct {
	Security []string
	SortBy   string
	Page     int
}
type BlockHash struct {
	*gorm.Model
	UserCID      string
	Hash         string
	PreviousHash string
}

func (b *BlockHash) create() error {
	return db.DB.Create(b).Error
}
func GetLatestBlockHash(userCID string) (BlockHash, error) {
	var b BlockHash
	err := db.DB.Where("user_cid = ?", userCID).Last(&b).Error
	return b, err
}

type TradeInterface interface {
	GetAccount() models.DematAccount
	GetInvestedValue() float64
	GetTradeData() portfoliobase.TradeData
}

// lets do it blockchain style
func RegisterTrade(td TradeInterface) error {
	var err error
	switch td := td.(type) {
	case *pstock.StockTrade:
		err = pstock.RegisterStockTrade(td)
	case *pbond.BondTrade:
		err = pbond.RegisterBondTrade(td)
	case *pmutualfund.MutualFundTrade:
		err = pmutualfund.RegisterMutualFundTrade(td)
	case *pets.ETSTrade:
		err = pets.RegisterETSTrade(td)
	}
	// createHashBlockForTrade(&td)
	return err
}

func createHashBlockForTrade(td *TradeInterface) error {
	s := structs.New(td)
	m := fmt.Sprint(s.Map())
	accountCID := (*td).GetAccount().UserCID
	latestBlock, err := GetLatestBlockHash(accountCID)
	if err != nil {
		return err
	}
	hasher := sha1.New()
	hasher.Write([]byte(m))
	hash := hasher.Sum(nil)
	hashString := fmt.Sprintf("%x", hash)
	blockHash := BlockHash{
		UserCID:      accountCID,
		Hash:         hashString,
		PreviousHash: latestBlock.Hash,
	}
	err = blockHash.create()
	if err != nil {
		return err
	}
	return nil
}

func GetTrades(tf TradeFilters, userCID string) ([]portfoliobase.TradeData, error) {
	var err error
	var trades []portfoliobase.TradeData
	var st []pstock.StockTrade
	var bt []pbond.BondTrade
	var et []pets.ETSTrade
	var mft []pmutualfund.MutualFundTrade
	dematAccounts, err := models.GetDematAccountIDsByUserCID(userCID)
	if err != nil {
		return []portfoliobase.TradeData{}, err
	}
	if utils.ArrayContains(tf.Security, "stock") || utils.ArrayContains(tf.Security, "all") {
		err = db.DB.Model(pstock.StockTrade{}).Preload("Stock").Where("account_id IN ?", dematAccounts).Find(&st).Error
	}
	if utils.ArrayContains(tf.Security, "bond") || utils.ArrayContains(tf.Security, "all") {
		err = db.DB.Model(pbond.BondTrade{}).Preload("Bond").Where("account_id IN ?", dematAccounts).Find(&bt).Error
	}
	if utils.ArrayContains(tf.Security, "ets") || utils.ArrayContains(tf.Security, "all") {
		err = db.DB.Model(pets.ETSTrade{}).Preload("ETS").Where("account_id IN ?", dematAccounts).Find(&et).Error
	}
	if utils.ArrayContains(tf.Security, "mutual-fund") || utils.ArrayContains(tf.Security, "all") {
		err = db.DB.Model(pmutualfund.MutualFundTrade{}).Preload("MutualFund").Where("account_id IN ?", dematAccounts).Find(&mft).Error
	}
	for _, s := range st {
		trades = append(trades, s.GetTradeData())
	}
	for _, b := range bt {
		trades = append(trades, b.GetTradeData())
	}
	for _, m := range mft {
		trades = append(trades, m.GetTradeData())
	}
	for _, e := range et {
		trades = append(trades, e.GetTradeData())
	}
	return trades, nil
}
