package mutualfund

import (
	"time"

	"github.com/samar2170/portfolio-manager-v4/pkg/db"
	"gorm.io/gorm"
)

func init() {
	db.DB.AutoMigrate(MutualFund{})
	db.DB.AutoMigrate(MutualFundNavHistory{})
}

type MutualFund struct {
	*gorm.Model
	ID                   int
	SchemeName           string `gorm:"index"`
	SchemeCategory       string `gorm:"index"`
	SchemeNavName        string `gorm:"uniqueIndex"`
	ParentSchemeCategory string `gorm:"index"`
	PriceToBeUpdated     bool
}

type MutualFundNavHistory struct {
	*gorm.Model
	ID           int
	MutualFund   MutualFund
	MutualFundID int
	Nav          float64
	Date         time.Time
	Source       string
}

func (m *MutualFundNavHistory) create() error {
	err := db.DB.Create(&m).Error
	return err
}

func (m *MutualFund) create() error {
	err := db.DB.Create(&m).Error
	return err
}

func (m *MutualFund) getOrCreate() (MutualFund, error) {
	err := db.DB.FirstOrCreate(&m, "scheme_nav_name = ?", m.SchemeNavName).Error
	return *m, err
}

func GetMutualFundBySchemeNavName(schemeNavName string) (MutualFund, error) {
	var mf MutualFund
	err := db.DB.First(&mf, "scheme_nav_name = ?", schemeNavName).Error
	return mf, err
}

func SearchMutualFund(query string) MutualFund {
	var mf MutualFund
	db.DB.Where("scheme_nav_name ILIKE ?", query).Find(&mf)
	return mf
}

func GetMutualFundByID(id int) (MutualFund, error) {
	var mf MutualFund
	err := db.DB.First(&mf, id).Error
	return mf, err
}
