package reddit

type ListingData struct {
	Dist     int    `json:"dist"`
	Children []Post `json:"children"`
}
