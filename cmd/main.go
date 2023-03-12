package main

import (
	"adnan/binance-bot/pkg/app"
	"fmt"

	"github.com/sdcoffey/techan"
)

func main() {
	limit := 24
	app.GetTickerData(limit)
	series := app.GetTimeSeries("ETHUSDT")
	strategy := app.Strategy1(series)

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
