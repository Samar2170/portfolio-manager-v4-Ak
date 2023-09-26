package fetch

import (
	"encoding/json"
	"strings"
)

var SymbolSuffix map[string]string = map[string]string{
	"NSE": ".NSE",
	"BSE": ".BSE",
}

type AvResponse struct {
	MetaData struct {
		Information string `json:"1. Information"`
		Symbol      string `json:"2. Symbol"`
		LastRefresh string `json:"3. Last Refreshed"`
		OutputSize  string `json:"5. Output Size"`
		TimeZone    string `json:"6. Time Zone"`
	} `json:"Meta Data"`
	TimeSeries map[string]struct {
		Open   float64 `json:"1. open,string"`
		High   float64 `json:"2. high,string"`
		Low    float64 `json:"3. low,string"`
		Close  float64 `json:"4. close,string"`
		Volume float64 `json:"5. volume,string"`
	} `json:"Time Series (Daily)"`
}

func FetchAVStockData(symbol, exchange string) (AvResponse, error) {
	symbol = strings.ToUpper(symbol)
	if _, ok := SymbolSuffix[exchange]; ok {
		symbol += SymbolSuffix[exchange]
	}
	uri := AVBaseUrl + "function=" + UpdatePriceQueryAVFunc + "&symbol=" + symbol + "&apikey=" + AVApiKey
	req := NewBaseRequest("GET", uri, nil, nil, nil, nil, AVApiKey)
	resp, err := req.Execute(3, false)
	if err != nil {
		return AvResponse{}, err
	}
	defer resp.Body.Close()

	var avResp AvResponse
	err = json.NewDecoder(resp.Body).Decode(&avResp)
	if err != nil {
		return avResp, err
	}
	return avResp, nil
}
