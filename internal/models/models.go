package models

import (
	"errors"

	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/db"
	"gorm.io/gorm"
)

const TradeDateFormat = "2006-01-02"

func init() {
	db.DB.AutoMigrate(&User{})
	db.DB.AutoMigrate(&BankAccount{})
	db.DB.AutoMigrate(&DematAccount{})
	db.DB.AutoMigrate(&GeneralAccount{})
}

type DBModel interface {
	create() error
	update() error
}

func CreateModelInstance(model DBModel) error {
	err := model.create()
	switch err {
	case gorm.ErrDuplicatedKey:
		return errors.New("user already exists")
	default:
		return err
	}
}

func updateModelInstance(model DBModel) error {
	err := model.update()
	switch err {
	case gorm.ErrRecordNotFound:
		return errors.New("user not found")
	default:
		return err
	}
}
