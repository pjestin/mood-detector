package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/pjestin/mood-detector/io"
	"github.com/pjestin/mood-detector/model/binance"
)

func parseBinanceKlines(filePath string) ([]binance.Kline, error) {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var klines []binance.Kline

	err = json.Unmarshal(dat, &klines)
	if err != nil {
		return nil, err
	}

	return klines, nil
}

func main() {
	klines, err := parseBinanceKlines("data/binance_klines_ETHUSDT_1h_1502942400000.json")
	if err != nil {
		log.Fatalln("Error reading Binance klines")
	}

	redis_reddit := io.RedisClient{}
	redis_reddit.Init(1)
	redis_reddit.Get()
}
