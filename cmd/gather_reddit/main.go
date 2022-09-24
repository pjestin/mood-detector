package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pjestin/mood-detector/io"
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

	posts_json, err := json.Marshal(&posts)
	if err != nil {
		log.Fatalln("Error when marshalling posts:", err)
	}

	redis := io.RedisClient{}
	redis.Init(1)
	now := time.Now().Format(time.RFC3339)
	log.Println("Adding posts to Redis")
	err = redis.Set(now, posts_json)
	if err != nil {
		log.Fatalln("Error when setting posts in Redis:", err)
	}
}
