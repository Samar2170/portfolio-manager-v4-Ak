package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/samar2170/portfolio-manager-v4/api"
	"github.com/samar2170/portfolio-manager-v4/internal/bulkupload"

	"github.com/samar2170/portfolio-manager-v4/security/bond"
	"github.com/samar2170/portfolio-manager-v4/security/ets"
	mutualfund "github.com/samar2170/portfolio-manager-v4/security/mutual-fund"
	"github.com/samar2170/portfolio-manager-v4/security/stock"
)

func main() {
	t := time.Now()
	arg := os.Args[1]
	fmt.Println(arg)
	switch arg {
	case "setup":
		fmt.Println("setting up")
		setup()
	case "dev":
		dev()
	default:
		fmt.Println("starting")
		api.StartServer()
	}
	fmt.Println(time.Since(t))
}

func setup() {
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		stock.LoadData()
		wg.Done()
	}()
	go func() {
		bond.LoadData()
		wg.Done()
	}()
	go func() {
		mutualfund.LoadData()
		wg.Done()
	}()
	go func() {
		ets.LoadData()
		wg.Done()
	}()
	wg.Wait()
}

func dev() {
	bulkupload.ParseBulkUploadSheets()
}
