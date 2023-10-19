package fetch

import "github.com/spf13/viper"

const (
	UpdatePriceQueryAVFunc = "TIME_SERIES_DAILY"
	AVBaseUrl              = "https://www.alphavantage.co/query?"
	AvTimeFormat           = "2006-01-02"
)

var AVApiKey string

func init() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	AVApiKey = viper.GetString("AV_API_KEY")
}
