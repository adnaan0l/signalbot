package main

import (
	"adnan/binance-bot/pkg/app"
	"fmt"
	"log"

	"github.com/sdcoffey/techan"
)

const (
	limit = 24
	// symbol   = "ETHUSDT"
	startIdx = 20
)

func main() {
	limit := 24

	symbolList, err := app.GetSymbolList()
	if err != nil {
		log.Printf("failed to get symbol list: %v\n", err)
	}

	if err := app.LoadTickerData(symbolList, limit); err != nil {
		log.Printf("failed to load ticker data: %v", err)
	}

	for _, symbol := range symbolList[:5] {
		fmt.Printf("Checking %s", symbol)
		series, err := app.GetTimeSeries(symbol)
		if err != nil {
			log.Fatalf("failed to get time series for %s: %v", symbol, err)
		}

		strategy, err := app.Strategy1(series)
		if err != nil {
			log.Fatalf("failed to create trading strategy: %v", err)
		}

		record := techan.NewTradingRecord()
		for i := 20; i < limit; i++ {
			if strategy.ShouldEnter(i, record) {
				fmt.Printf("Buy at %s", series.Candles[i].MinPrice)
				break
			} else if strategy.ShouldExit(i, record) {
				fmt.Printf("Sell at %s", series.Candles[i].MaxPrice)
				break
			}
		}
		fmt.Println("Do Nothing")
	}
}
