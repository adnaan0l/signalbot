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

	log.Println("Getting Symbol List")
	symbolList, err := app.GetSymbolList()
	fmt.Println(symbolList)
	if err != nil {
		log.Fatalf("failed to get symbol list: %v\n", err)
	}
	if err := app.LoadTickerData(symbolList, limit); err != nil {
		log.Fatalf("failed to load ticker data: %v", err)
	}
	for _, symbol := range symbolList[:5] {
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
				log.Printf("Buy %s at %s", symbol, series.Candles[i].MinPrice)
				break
			} else if strategy.ShouldExit(i, record) {
				log.Printf("Sell at %s", series.Candles[i].MaxPrice)
				break
			}
		}
	}
}
