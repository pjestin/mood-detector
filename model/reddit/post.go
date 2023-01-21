package reddit

type Post struct {
	Kind string   `json:"kind"`
	Data PostData `json:"data"`
}
