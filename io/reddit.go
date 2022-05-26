package io

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/pjestin/mood-detector/model"
)

type RedditClient struct {
	authToken string
	client    http.Client
}

func (c *RedditClient) Init(authToken string) {
	c.authToken = authToken
	c.client = http.Client{}
}

func (c *RedditClient) buildParamString(params map[string]string) string {
	paramString := ""
	for key, value := range params {
		if len(paramString) == 0 {
			paramString += "?"
		} else {
			paramString += "&"
		}
		paramString += key + "=" + value
	}
	return paramString
}

func (c *RedditClient) Get(path string, params map[string]string) ([]byte, error) {
	paramString := c.buildParamString(params)
	req, err := http.NewRequest("GET", "https://oauth.reddit.com"+path+paramString, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.authToken)
	req.Header.Set("User-Agent", "MoodDetector/0.1.0 (by pjestin)")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(fmt.Sprintf("Error status code: %d", resp.StatusCode))
	}

	b, err := io.ReadAll(resp.Body)
	return b, err
}

func (c *RedditClient) GetHotPosts(subreddit string) ([]model.PostData, error) {
	b, err := c.Get("/r/CryptoCurrency/hot", map[string]string{"limit": "100", "show": "all"})
	if err != nil {
		return nil, err
	}

	var body model.Listing
	json.Unmarshal(b, &body)

	var posts []model.PostData
	for _, post := range body.Data.Children {
		posts = append(posts, post.Data)
	}

	return posts, nil
}
