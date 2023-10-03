package portfolio

import portfoliobase "github.com/samar2170/portfolio-manager-v4/internal/portfolio/portfolio-base"

var (
	limit  = 50
	offset = 0
)

type TradeFilters struct {
	Security []string
	SortBy   string
	Page     int
}

var nameComp func(portfoliobase.TradeData, portfoliobase.TradeData) int = func(a, b portfoliobase.TradeData) int {
	if a.Symbol < b.Symbol {
		return 1
	} else {
		return -1
	}
}

var dateComp func(portfoliobase.TradeData, portfoliobase.TradeData) int = func(a, b portfoliobase.TradeData) int {
	if a.Date.Before(b.Date) {
		return 1
	} else {
		return -1
	}
}
var valueComp func(portfoliobase.TradeData, portfoliobase.TradeData) int = func(a, b portfoliobase.TradeData) int {
	if a.InvestedValue < b.InvestedValue {
		return 1
	} else {
		return -1
	}
}
