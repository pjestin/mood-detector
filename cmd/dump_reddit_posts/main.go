package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/pjestin/mood-detector/io"
	"github.com/pjestin/mood-detector/model/reddit"
)

const (
	LIMIT = 100000000
	SUBREDDIT = "r/CryptoCurrency"
)

func main() {
	godotenv.Load()

	redditClient := io.RedditClient{}
	err := redditClient.Init()
	if err != nil {
		log.Fatalln("Error while retrieving access token:", err)
	}

	f, err := os.Create("data/posts.json")
	if err != nil {
		log.Fatalln("Error creating data file:", err)
	}

	defer f.Close()

	var after string
	var count int64
	var dist int64

	for after == "" || (count < LIMIT && dist > 0) {
		b, err := redditClient.Get(fmt.Sprintf("/%s/new", SUBREDDIT), map[string]string{"limit": "100", "show": "all", "after": after, "count": strconv.FormatInt(count, 10)})
		if err != nil {
			log.Fatalln("Unable to retrieve posts from Reddit:", err)
		}

		var body reddit.Listing
		err = json.Unmarshal(b, &body)
		if err != nil {
			log.Fatalln("Unable to unmarshall posts from Reddit:", err)
		}

		dist = int64(body.Data.Dist)
		count += dist
		after = body.Data.After

		var postData []reddit.PostData
		for _, post := range body.Data.Children {
			postData = append(postData, post.Data)
		}

		posts_json, err := json.Marshal(&postData)
		if err != nil {
			log.Fatalln("Error when marshalling posts:", err)
		}
	
		_, err = f.Write(posts_json)
		if err != nil {
			log.Fatalln("Error writing posts to file:", err)
		}
	
		_, err = f.WriteString("\n")
		if err != nil {
			log.Fatalln("Error separating posts in file:", err)
		}

		log.Println("Post data retrieved")

		time.Sleep(time.Second)
	}
}
