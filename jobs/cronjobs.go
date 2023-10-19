package jobs

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/samar2170/portfolio-manager-v4-Ak/internal/portfolio/pstock"
	"github.com/samar2170/portfolio-manager-v4-Ak/security/stock"
)

var cronLogger *log.Logger

func init() {
	t := time.Now()

	file, err := os.OpenFile(fmt.Sprintf("logs/Cron_logs_%d-%d-%d", t.Day(), int(t.Month()), t.Year()), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	cronLogger = log.New(file, "[Cron]-", log.Ldate|log.Ltime|log.Lshortfile)
	cronLogger.Println("Cron Server Working buddy")
}

func StartCronServer() {
	t := time.Now()

	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Day().At("00:00").Do(func() {
		cronLogger.Println("Calculating Stock Prices" + t.String())
		pstock.CalculatePriceToBeUpdatedRank()
	})
	s.Every(2).Hour().Do(func() {
		cronLogger.Println("Updating Stock Prices" + t.String())
		stock.UpdateStockPrices()
	})
}
