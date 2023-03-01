package app

import (
	"adnan/binance-bot/pkg/utils"
	"fmt"
)

func GetTickerData() {
	ctx, rdb := utils.GetRedisClient()
	symbolList := rdb.HGet(ctx, "binance", "symbols")
	fmt.Println(symbolList)
}
