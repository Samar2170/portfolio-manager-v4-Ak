package ets

// BIRET	Brookfield India Real Estate Trust	RR	16-Feb-21	275	200	INE0FDU25010	275
// MINDSPACE	Mindspace Business Parks REIT	RR	7-Aug-20	275	200	INE0CCU25019	275
// EMBASSY	Embassy Office Parks REIT	RR	1-Apr-19	300	400	INE041025011	300

var ETSList = []ETS{
	{Symbol: "EMBASSY", Name: "Embassy Office Parks REIT", SecurityCode: "INE041025011", Category: "REIT"},
	{Symbol: "MINDSPACE", Name: "Mindspace Business Parks REIT", SecurityCode: "INE0CCU25019", Category: "REIT"},
	{Symbol: "BIRET", Name: "Brookfield India Real Estate Trust", SecurityCode: "INE0FDU25010", Category: "REIT"},
	{Symbol: "EBBETF0430", Name: "Edelweiss ETF - Nifty Bank", SecurityCode: "INF204KB1QZ3", Category: "ETF"},
}

func LoadData() {
	for _, ets := range ETSList {
		ets.getOrCreate()
	}
}
