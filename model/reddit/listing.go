package reddit

type Listing struct {
	Kind string      `json:"kind"`
	Data ListingData `json:"data"`
}
