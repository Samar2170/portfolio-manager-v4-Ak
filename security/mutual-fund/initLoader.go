package mutualfund

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var MfNavTimeLayout = "02-Jan-2006"

func LoadData() {
	loadDataWithNav()
}

func loadDataWithNav() {
	file, err := os.Open("assets/init/NAVHistoryReport.txt")
	var schemeCategory string
	var amcName string
	if err != nil {
		log.Println(err)
	}
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}
	defer file.Close()

	for indx, line := range fileLines {
		if indx == 0 {
			continue
		}
		fields := strings.Split(line, ";")
		if len(fields) == 1 {
			if strings.Contains(line, "(") && strings.Contains(line, ")") {
				schemeCategory = fields[0]
				log.Println(schemeCategory)
			} else if len(fields[0]) > 10 {
				amcName = fields[0]
				log.Println(amcName)
			}
		} else {
			mf, err := GetMutualFundBySchemeNavName(fields[1])
			log.Println("fields -----------------------------")
			log.Println(fields)
			log.Println(mf)
			log.Println(err)
			log.Println("--------------------------------------")
			if err != nil {
				mf = MutualFund{
					SchemeNavName: fields[1],
				}
				mf, err = mf.getOrCreate()
				if err != nil {
					log.Println("error in creating mf")
					log.Println(err)
				}
			}
			navFloat, err := strconv.ParseFloat(fields[4], 32)
			if err != nil {
				log.Println("nav is not float type")
			}
			parsedTime, err := time.Parse(MfNavTimeLayout, fields[7])
			if err != nil {
				log.Println("time is not proper format")
			}
			mfnh := MutualFundNavHistory{
				MutualFundID: mf.ID,
				Nav:          navFloat,
				Date:         parsedTime,
			}
			mfnh.create()
		}
	}

}
