package main

import (
	"log"
	"time"

	"github.com/pjestin/mood-detector/io"
	"github.com/pjestin/mood-detector/util"
)

const AUTH_TOKEN = "487021076163-bQt9efTKJKpFSxWsno3V2QeIP1pJXw"

func main() {
	reddit := io.RedditClient{}
	reddit.Init(AUTH_TOKEN)
	posts, err := reddit.GetHotPosts("r/CryptoCurrrency")
	if err != nil {
		log.Fatalln("Error when getting hot posts from Reddit:", err)
	}

	mood := util.ProcessPostMood(posts)
	log.Println("Total mood:", mood)

	redis := io.RedisClient{}
	redis.Init()
	now := time.Now().Format(time.RFC3339)
	err = redis.Set(now, mood)
	if err != nil {
		log.Fatalln("Error when setting mood in Redis:", err)
	}
}
