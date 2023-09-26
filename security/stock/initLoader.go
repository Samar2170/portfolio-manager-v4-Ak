package stock

import (
	"log"

	"github.com/xuri/excelize/v2"
)

func loadNSEStocks() {
	f, err := excelize.OpenFile("assets/init/ind_nifty500list.xlsx")
	if err != nil {
		log.Println(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		log.Println(err)
		return
	}
	for i, row := range rows {
		if i == 0 {
			continue
		}
		s := Stock{
			Symbol:   row[2],
			Name:     row[0],
			Industry: row[1],
		}
		_, e := s.GetOrCreate()
		if e != nil {
			log.Println(e)
		}
	}
}

func loadBSEStocks() {
	f, err := excelize.OpenFile("assets/init/BSEEquities.xlsx")
	if err != nil {
		log.Println(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()
	rows, err := f.GetRows("Equity")
	if err != nil {
		log.Println(err)
		return
	}
	indexMap := make(map[string]int)
	for i, row := range rows {
		if i == 0 {
			for j, dp := range row {
				indexMap[dp] = j
			}
		}
		if row[indexMap["Active"]] == "Active" {
			s := Stock{
				Symbol:   row[indexMap["Security Id"]],
				Name:     row[indexMap["Issuer Name"]],
				Industry: row[indexMap["Sector Name"]],
			}
			_, e := s.GetOrCreate()
			if e != nil {
				log.Println(e)
			}
		}
	}
}

func LoadData() {
	loadNSEStocks()
	loadBSEStocks()
}
