package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/pjestin/mood-detector/io"
	"github.com/pjestin/mood-detector/util"
)

func main() {
	godotenv.Load()

	subreddit := os.Getenv("SUBREDDIT")
	if len(subreddit) == 0 {
		log.Fatalln("SUBREDDIT variable not set")
	}

	reddit := io.RedditClient{}
	err := reddit.Init()
	if err != nil {
		log.Fatalln("Error while retrieving access token:", err)
	}

	posts, err := reddit.GetHotPosts(subreddit)
	if err != nil {
		log.Fatalln("Error when getting hot posts from Reddit:", err)
	}

	mood := util.ProcessPostMood(posts)
	log.Printf("Total mood: %d", mood)

	symbol := os.Getenv("BINANCE_SYMBOL")
	quantity := os.Getenv("BINANCE_QUANTITY")

	mood_upper_bound, err := strconv.Atoi(os.Getenv("MOOD_UPPER_BOUND"))
	if err != nil {
		log.Fatalln("Unable to parse MOOD_UPPER_BOUND:", err)
	}

	mood_lower_bound, err := strconv.Atoi(os.Getenv("MOOD_LOWER_BOUND"))
	if err != nil {
		log.Fatalln("Unable to parse MOOD_LOWER_BOUND:", err)
	}

	binance := io.BinanceClient{}
	binance.Init()

	isInBoughtState := false

	// Get last trade
	lastTrade, err := binance.GetLastTrade(symbol)
	if err != nil {
		log.Println("Unable to get last trade:", err)
	} else {
		log.Println("Last trade:", lastTrade)
		isInBoughtState = lastTrade.IsBuyer
	}

	log.Println("In bought state?", isInBoughtState)

	// Trade on Binance
	if isInBoughtState && mood >= mood_upper_bound {
		log.Println("Posting sell order")
		order, err := binance.PostOrder(symbol, "MARKET", "SELL", quantity)
		if err != nil {
			log.Fatalln("Unable to post sell order:", err)
		}
		log.Printf("Sell order successful: %#v", order)
		return
	}

	if !isInBoughtState && mood <= mood_lower_bound {
		log.Println("Posting buy order")
		order, err := binance.PostOrder(symbol, "MARKET", "BUY", quantity)
		if err != nil {
			log.Fatalln("Unable to post buy order:", err)
		}
		log.Printf("Buy order successful: %#v", order)
		return
	}

	log.Println("No trade to be done")
}
