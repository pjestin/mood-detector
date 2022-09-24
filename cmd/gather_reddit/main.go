package main

import (
	"encoding/json"
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

	posts_json, err := json.Marshal(&posts)
	if err != nil {
		log.Fatalln("Error when marshalling posts:", err)
	}

	mood := util.ProcessPostMood(posts)
	log.Println("Total mood:", mood)

	now := time.Now().Format(time.RFC3339)

	log.Println("Adding mood to Redis")
	redis_mood := io.RedisClient{}
	redis_mood.Init(0)
	err = redis_mood.Set(now, mood)
	if err != nil {
		log.Fatalln("Error when adding mood to Redis:", err)
	}

	log.Println("Adding posts to Redis")
	redis_reddit := io.RedisClient{}
	redis_reddit.Init(1)
	err = redis_reddit.Set(now, posts_json)
	if err != nil {
		log.Fatalln("Error when setting posts to Redis:", err)
	}
}
