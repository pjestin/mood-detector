package main

import (
	"encoding/json"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/pjestin/mood-detector/io"
	"github.com/pjestin/mood-detector/model/binance"
	"github.com/pjestin/mood-detector/model/reddit"
	"github.com/pjestin/mood-detector/util"
)

type DataPoint struct {
	instant     time.Time
	price       string
	redditPosts []reddit.PostData
}

func parseBinanceKlines(filePath string) ([]binance.Kline, error) {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var klinesData [][]interface{}

	err = json.Unmarshal(dat, &klinesData)
	if err != nil {
		return nil, err
	}

	var klines []binance.Kline

	for _, klineData := range klinesData {
		klines = append(klines, binance.Kline{
			OpenTime:         klineData[0].(float64),
			Open:             klineData[1].(string),
			High:             klineData[2].(string),
			Low:              klineData[3].(string),
			Close:            klineData[4].(string),
			Volume:           klineData[5].(string),
			CloseTime:        klineData[6].(float64),
			QuoteAssetVolume: klineData[7].(string),
			NumberOfTrades:   klineData[8].(float64),
		})
	}

	return klines, nil
}

func getDataPoints() ([]DataPoint, error) {
	var dataPoints []DataPoint

	klines, err := parseBinanceKlines("data/binance_klines_ETHUSDT_1h_1502942400000.json")
	if err != nil {
		return nil, err
	}

	redis := io.RedisClient{}
	redis.Init(1)
	redisEntries, err := redis.GetEntries()
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(redisEntries))
	for redisKey := range redisEntries {
		keys = append(keys, redisKey)
	}

	sort.Strings(keys)

	klineIndex := 0
	keyIndex := 0

	for keyIndex < len(keys) {
		keyTime, err := time.Parse(time.RFC3339, keys[keyIndex])
		if err != nil {
			return nil, err
		}

		for int64(klines[klineIndex].CloseTime) < keyTime.UnixMilli() {
			klineIndex += 1

			if klineIndex >= len(klines) {
				return dataPoints, nil
			}
		}

		postDataString := redisEntries[keys[keyIndex]]
		var postData []reddit.PostData

		err = json.Unmarshal([]byte(postDataString), &postData)
		if err != nil {
			return nil, err
		}

		dataPoints = append(dataPoints, DataPoint{
			instant:     time.UnixMilli(int64(klines[klineIndex].OpenTime)),
			price:       klines[klineIndex].Open,
			redditPosts: postData,
		})

		keyIndex += 1
	}

	return dataPoints, nil
}

type State struct {
	isBought bool
	capital  float64
	buyPrice float64
}

const (
	MIN_SELL_MOOD = 160
	MAX_BUY_MOOD  = -160
)

func main() {
	dataPoints, err := getDataPoints()
	if err != nil {
		log.Fatalln("Error getting data points", err)
	}
	duration := dataPoints[len(dataPoints)-1].instant.Sub(dataPoints[0].instant).Hours() / 24
	log.Println("Number of data points:", len(dataPoints), "; duration (days):", duration)

	state := State{
		isBought: false,
		capital:  1.0,
		buyPrice: 0.0,
	}

	for _, dataPoint := range dataPoints {
		price, err := strconv.ParseFloat(dataPoint.price, 64)
		if err != nil {
			log.Fatalln("Error parsing price", err)
		}

		mood := util.ProcessPostMood(dataPoint.redditPosts)
		if mood >= MIN_SELL_MOOD && state.isBought {
			log.Println("Instant:", dataPoint.instant, "; selling at", price)
			state.isBought = false
			state.capital *= price / state.buyPrice
		} else if mood <= MAX_BUY_MOOD && !state.isBought {
			log.Println("Instant:", dataPoint.instant, "; buying at", price)
			state.isBought = true
			state.buyPrice = price
		}
	}

	log.Printf("Capital: %.4f; expected factor over 1 year: %.4f", state.capital, math.Pow(state.capital, 365/duration))
}
