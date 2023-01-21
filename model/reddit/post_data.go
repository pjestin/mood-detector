package reddit

type PostData struct {
	Name        string  `json:"name"`
	Subreddit   string  `json:"subreddit"`
	Title       string  `json:"title"`
	Selftext    string  `json:"selftext"`
	Downs       float64 `json:"downs"`
	Ups         float64 `json:"ups"`
	Score       float64 `json:"score"`
	UpvoteRatio float64 `json:"upvote_ratio"`
}
