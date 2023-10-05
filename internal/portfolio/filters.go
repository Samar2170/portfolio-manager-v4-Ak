package portfolio

import portfoliobase "github.com/samar2170/portfolio-manager-v4-Ak/internal/portfolio/portfolio-base"

var (
	limit  = 50
	offset = 0
)

type TradeFilters struct {
	Security []string
	SortBy   string
	Page     int
}

var dateComp func(portfoliobase.TradeData, portfoliobase.TradeData) int = func(a, b portfoliobase.TradeData) int {
	if a.Date.Before(b.Date) {
		return 1
	} else {
		return -1
	}
}

func nameComp[T any](a, b T) int {
	if _, ok := any(a).(portfoliobase.TradeData); ok {
		if any(a).(portfoliobase.TradeData).Symbol < any(b).(portfoliobase.TradeData).Symbol {
			return -1
		} else if any(a).(portfoliobase.TradeData).Symbol > any(b).(portfoliobase.TradeData).Symbol {
			return 1
		}
	}
	if _, ok := any(a).(portfoliobase.HoldingData); ok {
		if any(a).(portfoliobase.HoldingData).Symbol < any(b).(portfoliobase.HoldingData).Symbol {
			return -1
		} else if any(a).(portfoliobase.HoldingData).Symbol > any(b).(portfoliobase.HoldingData).Symbol {
			return 1
		}
	}
	return 0
}
func valueComp[T any](a, b T) int {
	if _, ok := any(a).(portfoliobase.TradeData); ok {
		if any(a).(portfoliobase.TradeData).InvestedValue < any(b).(portfoliobase.TradeData).InvestedValue {
			return 1
		} else {
			return -1
		}
	}
	if _, ok := any(a).(portfoliobase.HoldingData); ok {
		if any(a).(portfoliobase.HoldingData).InvestedValue < any(b).(portfoliobase.HoldingData).InvestedValue {
			return 1
		} else {
			return -1
		}
	}
	return 0
}

// var nameComp func(portfoliobase.TradeData, portfoliobase.TradeData) int = func(a, b portfoliobase.TradeData) int {
// 	if a.Symbol < b.Symbol {
// 		return 1
// 	} else {
// 		return -1
// 	}
// }

// var valueComp func(portfoliobase.TradeData, portfoliobase.TradeData) int = func(a, b portfoliobase.TradeData) int {
// 	if a.InvestedValue < b.InvestedValue {
// 		return 1
// 	} else {
// 		return -1
// 	}
// }
