package main

import (
	"log"
	"os"
	"time"

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
	log.Println("Total mood:", mood)

	redis := io.RedisClient{}
	redis.Init()
	now := time.Now().Format(time.RFC3339)
	log.Println("Adding mood to Redis")
	err = redis.Set(now, mood)
	if err != nil {
		log.Fatalln("Error when setting mood in Redis:", err)
	}
}
