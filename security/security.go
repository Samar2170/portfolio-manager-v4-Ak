package security

type Security interface {
	create() error
	getOrCreate() (Security, error)
	GetLatestPrice() (float64, error)
}
