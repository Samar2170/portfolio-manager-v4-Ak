package portfoliobase

import (
	"time"

	"github.com/samar2170/portfolio-manager-v4-Ak/internal/models"
)

type TradeInterface interface {
	GetAccount() models.DematAccount
	GetInvestedValue() float64
	GetTradeData() TradeData
}

type TradeData struct {
	Symbol        string
	Quantity      float64
	Price         float64
	InvestedValue float64
	Date          time.Time
}

func (td *TradeData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"symbol":         td.Symbol,
		"quantity":       td.Quantity,
		"price":          td.Price,
		"invested_value": td.InvestedValue,
		"date":           td.Date,
	}
}

type HoldingData struct {
	Symbol        string
	Quantity      float64
	Price         float64
	InvestedValue float64
}

func (hd *HoldingData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"symbol":         hd.Symbol,
		"quantity":       hd.Quantity,
		"price":          hd.Price,
		"invested_value": hd.InvestedValue,
	}
}
