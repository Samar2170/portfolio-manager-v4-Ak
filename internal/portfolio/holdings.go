package portfolio

type HoldingInterface interface {
	create() error
	GetInvestedValue() float64
	// getCurrentValue() float64
	// getProfit() float64
	// getProfitPercentage() float64
}
