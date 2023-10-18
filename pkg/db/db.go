package db

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DBURI string
var environment string
var DB *gorm.DB

func connect() {
	var err error
	if environment == "prod" {
		DB, err = gorm.Open(sqlite.Open("db/prod.db"), &gorm.Config{})
		// DB, err = gorm.Open(postgres.Open(DBURI), &gorm.Config{
		// 	DisableForeignKeyConstraintWhenMigrating: true,
		// })
	} else {
		DB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	}
	if err != nil {
		panic(err)
	}
}

func loadEnvVars() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	DBURI = fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable",
		viper.GetString("POSTGRES_HOST"),
		viper.GetString("POSTGRES_USER"),
		viper.GetString("POSTGRES_DB"),
		viper.GetString("POSTGRES_PASSWORD"),
		viper.GetString("POSTGRES_PORT"),
	)
	log.Println(DBURI)
	environment = viper.GetString("ENVIRONMENT")
}

func init() {
	loadEnvVars()
	connect()
}
