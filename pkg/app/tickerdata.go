package app

import (
	"adnan/binance-bot/pkg/config"
	"adnan/binance-bot/pkg/utils"
	"encoding/json"
	"fmt"
)

func GetTickerData() uint8 {
	ctx, rdb := utils.GetRedisClient()
	res, err := rdb.HGet(ctx, "binance", "symbols").Result()
	if err != nil {
		fmt.Printf("error while doing GET command in gredis : %v", err)
	}

	var symbolList []string
	json.Unmarshal([]byte(res), &symbolList)

	pipeline := rdb.Pipeline()
	for _, symbol := range symbolList[:2] {
		fmt.Println(symbol)
		limit := "10"
		interval := "1h"
		endpointUrl := fmt.Sprintf(
			"%s?interval=%s&limit=%s&symbol=%s",
			config.CandleStick.String(), interval, limit, symbol,
		)
		_, body := utils.GetData(endpointUrl)
		pipeline.HSet(ctx, symbol, "ticker", string(body))
	}
	pipeline.Exec(ctx)
	return 200
}
