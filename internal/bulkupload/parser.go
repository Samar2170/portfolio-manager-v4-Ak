package bulkupload

import (
	"errors"
	"fmt"
	"log"

	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pstock"
	"github.com/xuri/excelize/v2"
)

func ParseBulkUploadSheets() {
	bulkUploadSheetIds, err := getUnparsedBulkUploadSheets(5)
	if err != nil {
		log.Println(err)
	}
	for _, buId := range bulkUploadSheetIds {
		err = parseBulkUploadTradeSheet(buId)
		log.Println(buId)
		log.Println(err)
	}

}

// TODO: complete this
func parseBulkUploadTradeSheet(buId uint) error {
	var err error
	var stockRows [][]string
	var bondRows [][]string
	var mfRows [][]string
	var etsRows [][]string
	bus, err := GetBulkUploadSheetByID(buId)
	if err != nil {
		return err
	}
	sheet, err := excelize.OpenFile(bus.Path)
	if err != nil {
		return err
	}
	defer func() {
		if err := sheet.Close(); err != nil {
			log.Println(err)
		}
	}()
	stockRows, err = sheet.GetRows("Stock")
	if err != nil {
		return err
	}
	bondRows, err = sheet.GetRows("Bond")
	if err != nil {
		return err
	}
	mfRows, err = sheet.GetRows("MutualFund")
	if err != nil {
		return err
	}
	etsRows, err = sheet.GetRows("ETS")
	if err != nil {
		return err
	}
	fmt.Println(stockRows)
	NewStockTrades(stockRows)
	fmt.Println(bondRows)
	fmt.Print(mfRows)
	fmt.Print(etsRows)
	return nil
}

func NewStockTrades(rows [][]string) ([]pstock.StockTrade, error) {
	if len(rows) == 0 {
		return []pstock.StockTrade{}, errors.New("empty sheet")
	} else if len(rows) == 1 {
		return []pstock.StockTrade{}, errors.New("no data")
	}
	var stockTradesMaps []map[string]string
	keys := rows[0]

	for _, row := range rows[1:] {
		stockTradeMap := make(map[string]string)
		for indx, field := range row {
			stockTradeMap[keys[indx]] = field
		}
		stockTradesMaps = append(stockTradesMaps, stockTradeMap)
	}
	fmt.Println(stockTradesMaps)
	// var stockTrades []pstock.StockTrade
	// for _, mapRow := range stockTradesMaps {
	// 	nst, err := pstock.NewStockTrade(mapRow["Symbol"], mapRow["Quantity"], mapRow["Price"], mapRow["TradeDate"], mapRow["TradeType"])
	// }
	return []pstock.StockTrade{}, errors.New("no data")
}

// func NewTrades(rows [][]string, typeOfTrade string) ([]portfolio.TradeInterface, error) {
// 	if len(rows) == 0 {
// 		return []portfolio.TradeInterface{}, nil
// 	}
// 	switch typeOfTrade{
// 	case "Stock":

// 	}
// }
