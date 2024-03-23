package routes

type SearchResponse struct {
	DateTime string `json:"dateTime"`
	Count uint64 `json:"count"`
}