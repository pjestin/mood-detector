package io

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pjestin/mood-detector/model/reddit"
)

func buildParamString(params map[string]string) string {
	paramString := ""
	for key, value := range params {
		if len(paramString) != 0 {
			paramString += "&"
		}
		paramString += key + "=" + value
	}
	return paramString
}

type RedditClient struct {
	authToken string
	client    http.Client
}

func (c *RedditClient) Init() error {
	log.Println("Requesting Reddit access token")

	user := os.Getenv("REDDIT_USER")
	password := os.Getenv("REDDIT_PASSWORD")
	client_id := os.Getenv("REDDIT_CLIENT_ID")
	client_secret := os.Getenv("REDDIT_CLIENT_SECRET")

	body_params := map[string]string{"grant_type": "password", "username": user, "password": password}
	body_bytes := []byte(buildParamString(body_params))

	req, err := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token", bytes.NewReader(body_bytes))
	if err != nil {
		return err
	}

	basic_auth := fmt.Sprintf("r%s:%s", client_id, client_secret)
	base64_auth := base64.StdEncoding.EncodeToString([]byte(basic_auth))

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64_auth))
	req.Header.Set("User-Agent", "MoodDetector/0.1.0 (by pjestin)")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unable to retrieve auth token; error status code: %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var body reddit.AuthTokenResponse
	err = json.Unmarshal(b, &body)
	if err != nil {
		return err
	}

	c.authToken = body.AccessToken
	log.Println("Access token retrieved")
	return nil
}

func (c *RedditClient) Get(path string, params map[string]string) ([]byte, error) {
	paramString := buildParamString(params)
	req, err := http.NewRequest("GET", fmt.Sprintf("https://oauth.reddit.com%s?%s", path, paramString), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))
	req.Header.Set("User-Agent", "MoodDetector/0.1.0 (by pjestin)")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error status code: %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	return b, err
}

func (c *RedditClient) GetHotPosts(subreddit string) ([]reddit.PostData, error) {
	b, err := c.Get(fmt.Sprintf("/%s/hot", subreddit), map[string]string{"limit": "100", "show": "all"})
	if err != nil {
		return nil, err
	}

	var body reddit.Listing
	err = json.Unmarshal(b, &body)
	if err != nil {
		return nil, err
	}

	var posts []reddit.PostData
	for _, post := range body.Data.Children {
		posts = append(posts, post.Data)
	}

	return posts, nil
}
