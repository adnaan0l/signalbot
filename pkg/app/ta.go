package app

import (
	"adnan/binance-bot/pkg/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

func GetTimeSeries(symbol string) *techan.TimeSeries {

	ctx, rdb := utils.GetRedisClient()

	series := techan.NewTimeSeries()

	result, err := rdb.HGet(ctx, symbol, "ticker").Result()
	if err != nil {
		fmt.Printf("error while doing HGET command in gredis : %v", err)
	}

	var ticker [][]interface{}
	json.Unmarshal([]byte(result), &ticker)

	for idx := range ticker {
		start := ticker[idx][0].(float64)
		period := techan.NewTimePeriod(time.Unix(int64(start)/1000, 0), time.Second*3599)

		candle := techan.NewCandle(period)
		candle.OpenPrice = big.NewFromString(ticker[idx][1].(string))
		candle.MaxPrice = big.NewFromString(ticker[idx][2].(string))
		candle.MinPrice = big.NewFromString(ticker[idx][3].(string))
		candle.ClosePrice = big.NewFromString(ticker[idx][4].(string))
		candle.Volume = big.NewFromString(ticker[idx][5].(string))

		series.AddCandle(candle)
	}

	return series
}

func getBolBands(series *techan.TimeSeries) (techan.Indicator, techan.Indicator) {
	closePrices := techan.NewClosePriceIndicator(series)
	bolHigh := techan.NewBollingerUpperBandIndicator(closePrices, 21, 2.0)
	bolLow := techan.NewBollingerLowerBandIndicator(closePrices, 21, 2.0)

	return bolHigh, bolLow
}

func Strategy1(series *techan.TimeSeries) techan.RuleStrategy {

	highPrices := techan.NewHighPriceIndicator(series)
	lowPrices := techan.NewLowPriceIndicator(series)

	bolHigh, bolLow := getBolBands(series)

	entryRule := techan.And(
		techan.NewCrossDownIndicatorRule(lowPrices, bolLow),
		techan.PositionNewRule{})

	exitRule := techan.And(
		techan.NewCrossUpIndicatorRule(highPrices, bolHigh),
		techan.PositionNewRule{})

	strategy := techan.RuleStrategy{
		UnstablePeriod: 0,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
	return strategy
}
