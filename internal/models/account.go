package models

import (
	"github.com/samar2170/portfolio-manager-v4/pkg/db"
	"gorm.io/gorm"
)

type BankAccount struct {
	*gorm.Model
	ID        int
	UserCID   string `gorm:"unique;index"`
	Bank      string
	AccountNo string `gorm:"unique;index"`
}

type DematAccount struct {
	*gorm.Model
	ID          int
	UserCID     string `gorm:"index"`
	AccountCode string `gorm:"unique;index"`
	Broker      string
}

type GeneralAccount struct {
	*gorm.Model
	ID          int
	UserCID     string `gorm:"unique;index"`
	AccountCode string `gorm:"unique;index"`
}

func (b *BankAccount) create() error {
	return db.DB.Create(b).Error
}
func (b *BankAccount) update() error {
	return db.DB.Save(b).Error
}

func (d *DematAccount) create() error {
	return db.DB.Create(d).Error
}

func (d *DematAccount) update() error {
	return db.DB.Save(d).Error
}

func (g *GeneralAccount) create() error {
	return db.DB.Create(g).Error
}

func (g *GeneralAccount) update() error {
	return db.DB.Save(g).Error
}

func GetAccountsByUserCID(userCID string) ([]BankAccount, []DematAccount, []GeneralAccount, error) {
	var bankAccounts []BankAccount
	var dematAccounts []DematAccount
	var generalAccounts []GeneralAccount
	err := db.DB.Where("user_c_id = ?", userCID).Find(&bankAccounts).Error
	if err != nil {
		return nil, nil, nil, err
	}
	err = db.DB.Where("user_c_id = ?", userCID).Find(&dematAccounts).Error
	if err != nil {
		return nil, nil, nil, err
	}
	err = db.DB.Where("user_c_id = ?", userCID).Find(&generalAccounts).Error
	if err != nil {
		return nil, nil, nil, err
	}
	return bankAccounts, dematAccounts, generalAccounts, nil
}

func GetDematAccountsByUserCID(userCID string) ([]DematAccount, error) {
	var dematAccounts []DematAccount
	err := db.DB.Where("user_c_id = ?", userCID).Find(&dematAccounts).Error
	if err != nil {
		return nil, err
	}
	return dematAccounts, nil
}
func GetDematAccountIDsByUserCID(userCID string) ([]int, error) {
	var dematAccountIDs []int
	err := db.DB.Model(&DematAccount{}).Where("user_c_id = ?", userCID).Select("id").Find(&dematAccountIDs).Error
	if err != nil {
		return nil, err
	}
	return dematAccountIDs, nil
}
func GetBankAccountsByUserCID(userCID string) ([]BankAccount, error) {
	var bankAccounts []BankAccount
	err := db.DB.Where("user_c_id = ?", userCID).Find(&bankAccounts).Error
	if err != nil {
		return nil, err
	}
	return bankAccounts, nil
}

func GetGeneralAccountsByUserCID(userCID string) ([]GeneralAccount, error) {
	var generalAccounts []GeneralAccount
	err := db.DB.Where("user_c_id = ?", userCID).Find(&generalAccounts).Error
	if err != nil {
		return nil, err
	}
	return generalAccounts, nil
}

func GetDematAccountByCode(accountCode, userCID string) (DematAccount, error) {
	var da DematAccount
	err := db.DB.Where("account_code = ? AND user_c_id = ?", accountCode, userCID).First(&da).Error
	return da, err
}
