package reddit

type ListingData struct {
	After	 string `json:"after"`
	Before	 string `json:"before"`
	Children []Post `json:"children"`
	Dist     int    `json:"dist"`
}
