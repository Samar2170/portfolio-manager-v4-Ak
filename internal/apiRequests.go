package internal

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type BankAccountRequest struct {
	Bank      string `json:"bank"`
	AccountNo string `json:"account_no"`
}

type DematAccountRequest struct {
	AccountCode string `json:"account_code"`
	Broker      string `json:"broker"`
}

type StockTradeRequest struct {
	Symbol      string  `json:"symbol"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	TradeType   string  `json:"trade_type"`
	TradeDate   string  `json:"trade_date"`
	AccountCode string  `json:"account_code"`
}

// func (st *StockTradeRequest) BulkUploadTemplate() {
// 	fields := structs.Fields(st)
// }

type BondTradeRequest struct {
	Symbol      string  `json:"symbol"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	TradeType   string  `json:"trade_type"`
	TradeDate   string  `json:"trade_date"`
	AccountCode string  `json:"account_code"`
}

type MutualFundTradeRequest struct {
	MutualFundID int     `json:"mutual_fund_id"`
	Quantity     float64 `json:"quantity"`
	Price        float64 `json:"price"`
	TradeType    string  `json:"trade_type"`
	TradeDate    string  `json:"trade_date"`
	AccountCode  string  `json:"account_code"`
}

type ETSTradeRequest struct {
	Symbol      string  `json:"symbol"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	TradeType   string  `json:"trade_type"`
	TradeDate   string  `json:"trade_date"`
	AccountCode string  `json:"account_code"`
}
