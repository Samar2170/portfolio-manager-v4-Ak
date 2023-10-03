package portfolio

import (
	"github.com/samar2170/portfolio-manager-v4/internal/models"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pbond"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pets"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pmutualfund"
	portfoliobase "github.com/samar2170/portfolio-manager-v4/internal/portfolio/portfolio-base"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pstock"
	"github.com/samar2170/portfolio-manager-v4/pkg/db"
	"github.com/samar2170/portfolio-manager-v4/pkg/utils"
	"golang.org/x/exp/slices"
)

// lets do it blockchain style
func RegisterTrade(td portfoliobase.TradeInterface) error {
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

	switch tf.SortBy {
	case "name":
		slices.SortFunc(trades, nameComp)
	case "date":
		slices.SortFunc(trades, dateComp)
	case "value":
		slices.SortFunc(trades, valueComp)
	default:
		slices.SortFunc(trades, valueComp)
	}
	// if tf.Page == 0 {
	// 	tf.Page = 1
	// }
	// aoffset := (tf.Page - 1) * limit
	// alimit := aoffset + limit
	// if alimit > len(trades) {
	// 	return []portfoliobase.TradeData{}, errors.New("thats all")
	// }
	// return trades[aoffset:alimit], nil
	return trades, nil
}

func GetHoldings(tf TradeFilters, userCID string) ([]portfoliobase.HoldingData, error) {
	var err error
	var holdings []portfoliobase.HoldingData
	var sh []pstock.StockHolding
	var bh []pbond.BondHolding
	var eh []pets.ETSHolding
	var mfh []pmutualfund.MutualFundHolding
	dematAccounts, err := models.GetDematAccountIDsByUserCID(userCID)
	if err != nil {
		return []portfoliobase.HoldingData{}, err
	}
	if utils.ArrayContains(tf.Security, "stock") || utils.ArrayContains(tf.Security, "all") {
		err = db.DB.Model(pstock.StockHolding{}).Preload("Stock").Where("account_id IN ?", dematAccounts).Find(&sh).Error
	}
	if utils.ArrayContains(tf.Security, "bond") || utils.ArrayContains(tf.Security, "all") {
		err = db.DB.Model(pbond.BondHolding{}).Preload("Bond").Where("account_id IN ?", dematAccounts).Find(&bh).Error
	}
	if utils.ArrayContains(tf.Security, "ets") || utils.ArrayContains(tf.Security, "all") {
		err = db.DB.Model(pets.ETSHolding{}).Preload("ETS").Where("account_id IN ?", dematAccounts).Find(&eh).Error
	}
	if utils.ArrayContains(tf.Security, "mutual-fund") || utils.ArrayContains(tf.Security, "all") {
		err = db.DB.Model(pmutualfund.MutualFundHolding{}).Preload("MutualFund").Where("account_id IN ?", dematAccounts).Find(&mfh).Error
	}
	for _, s := range sh {
		holdings = append(holdings, portfoliobase.HoldingData{
			Symbol: s.Stock.Symbol, Quantity: float64(s.Quantity), Price: s.BuyPrice,
			InvestedValue: s.GetInvestedValue(),
		})
	}
	for _, b := range bh {
		holdings = append(holdings, portfoliobase.HoldingData{
			Symbol: b.Bond.Symbol, Quantity: float64(b.Quantity), Price: b.BuyPrice,
			InvestedValue: b.GetInvestedValue(),
		})
	}
	for _, m := range mfh {
		holdings = append(holdings, portfoliobase.HoldingData{
			Symbol: m.MutualFund.SchemeNavName, Quantity: float64(m.Quantity), Price: m.BuyPrice,
			InvestedValue: m.GetInvestedValue(),
		})
	}
	for _, e := range eh {
		holdings = append(holdings, portfoliobase.HoldingData{
			Symbol: e.ETS.Symbol, Quantity: float64(e.Quantity), Price: e.BuyPrice,
			InvestedValue: e.GetInvestedValue(),
		})
	}

	// switch tf.SortBy {
	// case "name":
	// 	slices.SortFunc(holdings, nameComp)
	// case "date":
	// 	slices.SortFunc(holdings, dateComp)
	// case "value":
	// 	slices.SortFunc(holdings, valueComp)
	// default:
	// 	slices.SortFunc(holdings, valueComp)
	// }
	// aoffset := (tf.Page - 1) * limit
	// alimit := aoffset + limit
	// if alimit > len(trades) {
	// 	return []portfoliobase.TradeData{}, errors.New("thats all")
	// }
	return holdings, nil

}
