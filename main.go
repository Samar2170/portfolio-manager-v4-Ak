package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/samar2170/portfolio-manager-v4-Ak/api"
	"github.com/samar2170/portfolio-manager-v4-Ak/internal/bulkupload"
	"github.com/samar2170/portfolio-manager-v4-Ak/jobs"

	"github.com/samar2170/portfolio-manager-v4-Ak/security/bond"
	"github.com/samar2170/portfolio-manager-v4-Ak/security/ets"
	mutualfund "github.com/samar2170/portfolio-manager-v4-Ak/security/mutual-fund"
	"github.com/samar2170/portfolio-manager-v4-Ak/security/stock"
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
		start()
	}
	fmt.Println(time.Since(t))
}

func start() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		api.StartServer()
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		jobs.StartCronServer()
		wg.Done()
	}()

	wg.Wait()

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
