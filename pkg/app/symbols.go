package app

import (
	"adnan/binance-bot/pkg/config"
	"adnan/binance-bot/pkg/models"
	"adnan/binance-bot/pkg/utils"
	"encoding/json"
	"fmt"
	"os"
)

func UpdateSymbolList() uint8 {
	var symbolList []string

	_, resBody := utils.GetData(config.ExchangeInfo.String())

	var parsed models.ExchangeInfoResponse
	if err := json.Unmarshal(resBody, &parsed); err != nil {
		os.Exit(1)
	}

	for _, symbol := range parsed.Symbols {
		if symbol.QuoteAsset == "USDT" {
			symbolList = append(symbolList, symbol.Symbol)
		}
	}

	j, err := json.Marshal(symbolList)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	ctx, rdb := utils.GetRedisClient()
	hsetErr := rdb.HSet(ctx, "binance", "symbols", j).Err()
	if hsetErr != nil {
		panic(hsetErr)
	}

	return 200
}
