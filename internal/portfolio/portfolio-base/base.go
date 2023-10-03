package portfoliobase

import (
	"time"

	"github.com/samar2170/portfolio-manager-v4/internal/models"
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

type HoldingData struct {
	Symbol        string
	Quantity      float64
	Price         float64
	InvestedValue float64
}
