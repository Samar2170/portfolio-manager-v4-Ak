package bond

import (
	"time"

	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/db"
)

func init() {
	db.DB.AutoMigrate(&Bond{})
	db.DB.AutoMigrate(&BondPriceHistory{})
}

type Bond struct {
	ID        int
	Symbol    string `gorm:"index"`
	Name      string
	IpFreq    string
	IpRate    float64
	IpDate    time.Time
	MtDate    time.Time
	FaceValue float64
	MtValue   float64
	IssueDate time.Time
}

type BondPriceHistory struct {
	ID     uint
	Bond   Bond
	BondID uint
	Price  float64
	Date   time.Time
}

func (b *Bond) create() error {
	err := db.DB.Create(&b).Error
	return err
}

func (b *Bond) getOrCreate() (Bond, error) {
	err := db.DB.FirstOrCreate(&b, Bond{Symbol: b.Symbol}).Error
	return *b, err
}

func (b *Bond) getLatestDate() (time.Time, error) {
	var bph BondPriceHistory
	err := db.DB.Order("date desc").First(&bph, "symbol = ?", b.Symbol).Error
	return bph.Date, err
}

func (b *BondPriceHistory) create() error {
	err := db.DB.Create(&b).Error
	return err
}

func GetBondBySymbol(symbol string) (Bond, error) {
	var b Bond
	err := db.DB.First(&b, "symbol = ?", symbol).Error
	return b, err
}
