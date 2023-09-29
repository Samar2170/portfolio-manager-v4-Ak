package models

import (
	"github.com/samar2170/portfolio-manager-v4/pkg/db"
	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	ID       int
	Username string `gorm:"uniqueIndex"`
	UserCID  string `gorm:"uniqueIndex"`

	Password string
	Email    string
	ApiKey   string `gorm:"uniqueIndex"`
}

// User signs up/login, we call rpc, save CID and username,
func (u *User) create() error {
	return db.DB.Create(u).Error
}

func (u *User) update() error {
	return db.DB.Save(u).Error
}

func (u *User) UpdateApiKey(hashedNewApiKey string) error {
	return db.DB.Model(&u).Update("api_key", hashedNewApiKey).Error
}

func GetUserByCID(userCID string) (User, error) {
	user := User{}
	err := db.DB.Where("user_c_id = ?", userCID).First(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func GetUserByUsername(username string) (User, error) {
	var u User
	err := db.DB.Where("username = ?", username).First(&u).Error
	return u, err
}
func GetUserByApiKey(hashedapiKey string) (User, error) {
	var u User
	err := db.DB.Where("api_key = ?", hashedapiKey).Find(&u).Error
	return u, err
}

type Watchlist struct {
	*gorm.Model
	UserCID string `gorm:"index"`
	Name    string
}
