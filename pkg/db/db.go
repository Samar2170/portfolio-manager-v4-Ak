package db

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DBURI string
var environment string
var DB *gorm.DB

func connect() {
	var err error
	if environment == "prod" {
		// DB, err = gorm.Open(sqlite.Open("database/prod.db"), &gorm.Config{})
		DB, err = gorm.Open(postgres.Open(DBURI), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
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
	environment = viper.GetString("ENVIRONMENT")
	pgHost, pgPort := viper.GetString("POSTGRES_HOST"), viper.GetInt("POSTGRES_PORT")
	pgUser, pgPassWord := viper.GetString("POSTGRES_USER"), viper.GetString("POSTGRES_PASSWORD")
	pgDB := viper.GetString("POSTGRES_DB")
	DBURI = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Kolkata", pgHost, pgUser, pgPassWord, pgDB, pgPort)
	log.Println(DBURI)
}

func init() {
	loadEnvVars()
	connect()
}
