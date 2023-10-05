package internal

import (
	"github.com/samar2170/portfolio-manager-v4-Ak/client/cognitio/cauth"
	"github.com/spf13/viper"
)

var passwordDecryptionKey string
var signingKey []byte
var cognitioClient = new(cauth.AuthServiceClient)

const TradeDateFormat = "2006-01-02"

func init() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	passwordDecryptionKey = viper.GetString("PASSWORD_ENCRYPTION_KEY")
	signingKey = []byte(viper.GetString("SIGNING_KEY"))
}
