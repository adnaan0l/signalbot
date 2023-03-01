package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func ChunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len((slice))
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

func GetData(endpointURL string) (int, []byte) {
	baseURL := "https://data.binance.com"
	url := baseURL + endpointURL
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
	}
	return res.StatusCode, resBody
}

func GetRedisClient() (context.Context, *redis.Client) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return ctx, rdb
}
